package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson" //mongodb
	"net"
	"sync"
	"time"
)

var (
	FILE_PATH = "./upload_file"
	Xorm      *xorm.Engine
	db        *sql.DB
	pool      *redis.Pool
	redisConn redis.Conn
	session   *mgo.Session

	ChatCollect         *mgo.Collection
	UserCollect         *mgo.Collection
	FrdRequesCollection *mgo.Collection //离线好友请求存到这个表格，用户上线再发送给用户
	FrdGroupCollect     *mgo.Collection //好友分组
	FriendCollect       *mgo.Collection //好友
	lastUserID          int             //注册使用  得到新的id++

	Conf Config //配置文件 结构体

	packHead           string = ` "142857""0x12345678"` //包头
	headSize           int    = len(packHead)           //包头大小
	signBeginPos       int    = 2                       //包签名开始位置
	signEndnPos        int    = 8                       //包签名结束位置
	dataLengthBeginPos int    = 10                      //包长开始位置
	dataLengthEndnPos  int    = 20                      //包长结束位置

	timeFormat string = "2006-01-02 15:04:05.8581" //时间格式

	clnOffLineChannel chan *myConn
	broadcastChannel  chan string

	userList      map[int]myConn = make(map[int]myConn) //在线用户列表  "id":myConn
	glogErr       int            = 3                    //日志输出等级 3错误
	glogwarn      int            = 2                    //日志输出等级 2警告
	glogInfo      int            = 1                    //日志输出等级 1信息
	SUCESS_RESULT string         = "1"                  //操作结果  1 成功
	FAILED_RESULT string         = "0"                  //操作结果  2 失败
)

type myConn struct {
	//对conn (net.Conn) 包装  每个用户资源
	id         int
	token      string
	conn       net.Conn
	listenList []iListener
	lock       sync.RWMutex
}

//*************************配置文件信息结构体

type Config struct {
	SrvAddress  string
	FileSrvAddr string // http文件服务器的地址

	DBuser       string
	DBpsw        string
	DBAddress    string
	DBport       string
	DBname       string
	RedisNetWork string
	RedisAddress string
	MongoDBAddrs string
	MongoDBuser  string
	MongoDBpsw   string

	Sign        string
	PackMaxSize int

	LoginOP             int
	RegisOP             int
	ChatOP              int
	ReadedOP            int
	OffLineMsgOP        int
	GetFriendListOP     int
	AddFriendRequestOP  int
	AddFriendResponseOP int
	DeltFriendOP        int
	RemarkFriendOP      int
	AddFriendGroupOP    int
	DetFriendGroupOP    int
	RemarFriendkGroupOP int

	CacheToDBTime int
	DBToCacheTime int
}

//*************************聊天信息
type ChatMsg struct {
	Src      int       `bson:"Src"`
	Dst      int       `bson:"Dst"`
	SendTime time.Time `bson:"SendTime"`
	Content  string    `bson:"Content"`
	MsgType  int       `bson:"MsgType"`
	ReadTime time.Time `bson:"ReadTime"`
	FileName string    `bson:"FileName"`
	FileMD5  string    `bson:"FileMD5"`
	FileLen  string    `bson:"FileLen"`
}

//*************************记录每个用户最后一次读取的时间
type UserReadTime struct {
	ID       int       `bson:"ID"`
	ReadTime time.Time `bson:"ReadTime"`
}

//*************************用户信息结构体
type UserStruct struct {
	User_id int    `xorm:"index unique"`
	Name    string `xorm:"varchar(20)"`
	Pwd     string `xorm:"varchar(20) NOT NULL"`
	Sex     int    `xorm:"int"`
	Age     int    `xorm:"int"`
	Intro   string `xorm:"varchar(100)"`
	Role_id int    `xorm:"int NOT NULL"`
}

//***************************好友
type Friend struct {
	ID      int       `bson:"ID"`
	FrdID   int       `bson:"FrdID"`
	Remark  string    `bson:"Remark" default ""` //备注
	AddTime time.Time `bson:"AddTime"`
	Group   int       `bson:"Group" default 0` //所在分组
}
type FriendRequest struct {
	//离线好友请求存到这个表格，用户上线再发送给用户
	SrcID       int       `bson:"SrcID"`
	DstID       int       `bson:"DstID"`
	SrcName     string    `bson:"SrcName"`
	RequestTime time.Time `bson:"RequestTime"`
}
type FriendGroup struct {
	//好友分组信息
	UID      int    `bson:"UID"`
	ListNO   int    `bson:"ListNO"`
	ListName string `bson:"ListName"`
}

//-----------------------------------------初始化-------------------------------------------------------------------
func init() {

	flag.Parse()
	glog.Infof("info %d", 0)
	glog.Warningf("warning %d", 1)
	glog.Errorf("error %d", 2)

	// 退出时调用，确保日志写入文件中
	glog.Flush()
	err := loadConfig()
	if err != nil {
		fmt.Println("请检查配置文件myConfig.json", err.Error())
		glog.Error("请检查配置文件myConfig.json" + err.Error())
	}

	err2 := setupDB(Conf.DBuser, Conf.DBpsw, Conf.DBAddress, Conf.DBport, Conf.DBname)
	if err2 != nil {
		fmt.Println("请检查mysql配置", err2.Error())
		glog.Error("请检查mysql配置" + err2.Error())
	}

	err3 := initMongoDB()
	if err3 != nil {
		fmt.Println("请检查MongoDB配置", err3.Error())
		glog.Error("请检查MongoDB配置" + err3.Error())
	}

	pool = newPool() //redis缓存  建立连接池

	redisConn = pool.Get()
	// dbToCache() // 加载数据到缓存

	//go timer() //启动定时器
	glog.V(2).Infoln("服务器启动成功")
}

//-----------------------------------------加载配置文件---------------------------------------------------------------
func loadConfig() error {
	viper.SetConfigName("Config") // name of config file (without extension)配置文件名（无扩展名）
	viper.AddConfigPath(".")      // path to look for the config file in
	viper.AddConfigPath("json")   // more path to look for the config files
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	Conf.LoginOP = viper.GetInt("LoginOP")
	Conf.RegisOP = viper.GetInt("RegisOP")
	Conf.ChatOP = viper.GetInt("ChatOP")
	Conf.ReadedOP = viper.GetInt("ReadedOP")
	Conf.OffLineMsgOP = viper.GetInt("OffLineMsgOP")
	Conf.GetFriendListOP = viper.GetInt("GetFriendListOP")
	Conf.AddFriendRequestOP = viper.GetInt("AddFriendRequestOP")
	Conf.AddFriendResponseOP = viper.GetInt("AddFriendResponseOP")
	Conf.DeltFriendOP = viper.GetInt("DeltFriendOP")
	Conf.RemarkFriendOP = viper.GetInt("RemarkFriendOP")
	Conf.AddFriendGroupOP = viper.GetInt("AddFriendGroupOP")
	Conf.DetFriendGroupOP = viper.GetInt("DetFriendGroupOP")
	Conf.RemarFriendkGroupOP = viper.GetInt("RemarFriendkGroupOP")

	Conf.Sign = viper.GetString("Sign")
	Conf.PackMaxSize = viper.GetInt("PackMaxSize")

	Conf.CacheToDBTime = viper.GetInt("CacheToDBTime")
	Conf.DBToCacheTime = viper.GetInt("DBTocacheTime")

	Conf.FileSrvAddr = viper.GetString("FileSrvAddr")
	Conf.SrvAddress = viper.GetString("SrvAddress")
	Conf.RedisNetWork = viper.GetString("RedisNetWork")
	Conf.RedisAddress = viper.GetString("RedisAddress")
	Conf.MongoDBuser = viper.GetString("MongoDBuser")
	Conf.MongoDBpsw = viper.GetString("MongoDBpsw")
	Conf.DBuser = viper.GetString("DBuser")
	Conf.DBpsw = viper.GetString("DBpsw")
	Conf.DBAddress = viper.GetString("DBAddress")
	Conf.DBport = viper.GetString("DBport")
	Conf.DBname = viper.GetString("DBname")

	fmt.Println("加载配置文件成功")
	glog.Error("加载配置文件成功")
	return err
}

func timer() {
	CacheToDBTime := time.NewTimer(time.Second * time.Duration(Conf.CacheToDBTime))
	DBToCacheTime := time.NewTimer(time.Second * time.Duration(Conf.DBToCacheTime))
	for {
		select {

		case <-CacheToDBTime.C:
			cacheToDB() // 缓存跟新到数据库
			CacheToDBTime.Reset(time.Second * time.Duration(Conf.CacheToDBTime))

		case <-DBToCacheTime.C:
			dbToCache() // 加载数据到缓存
			DBToCacheTime.Reset(time.Second * time.Duration(Conf.DBToCacheTime))
		}
	}
}

func initMongoDB() error {
	var err error
	dialInfo := &mgo.DialInfo{
		Addrs:  []string{"127.0.0.1"},
		Direct: false, Timeout: time.Second * 1,
		Database:  "login",
		Source:    "admin",
		Username:  Conf.MongoDBuser,
		Password:  Conf.MongoDBpsw,
		PoolLimit: 4096,
	}

	session, err = mgo.DialWithInfo(dialInfo) //连接MongoDB
	//session, err = mgo.Dial(Conf.MongoDBAddrs) //连接MongoDB
	if err != nil {
		fmt.Println("请检查MongoDB是否打开", err.Error())
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	ChatCollect = session.DB("login").C("ChatMsg")               //聊天记录集合
	UserCollect = session.DB("login").C("UserReadTime")          //用户最后一次读取时间集合
	FriendCollect = session.DB("login").C("Friend")              //好友
	FrdRequesCollection = session.DB("login").C("FriendRequest") ////离线好友请求存到这个表格，用户上线再发送给用户
	FrdGroupCollect = session.DB("login").C("FrdGroupCollect")   //好友分组
	err = ChatCollect.EnsureIndex(mgo.Index{Key: []string{"Src", "Dst", "-SendTime"}, DropDups: true})
	if err != nil {
		fmt.Println("请检查MongoDB配置", err.Error())
		return err
	}
	err = ChatCollect.EnsureIndexKey("FileMD5")
	if err != nil {
		fmt.Println("请检查MongoDB配置", err.Error())
		return err
	}
	err = UserCollect.EnsureIndex(mgo.Index{Key: []string{"ID"}, Unique: true})
	if err != nil {
		fmt.Println("请检查MongoDB配置", err.Error())
		return err
	}
	err = FriendCollect.EnsureIndex(mgo.Index{Key: []string{"ID", "FrdId"}, DropDups: true})
	if err != nil {
		fmt.Println("请检查MongoDB配置", err.Error())
		return err
	}
	err = FrdRequesCollection.EnsureIndex(mgo.Index{Key: []string{"SrcID", "DstID"}, Unique: true})
	if err != nil {
		fmt.Println("请检查MongoDB配置", err.Error())
		return err
	}
	err = FrdGroupCollect.EnsureIndex(mgo.Index{Key: []string{"UID", "ListNO"}, Unique: true})
	if err != nil {
		fmt.Println("请检查MongoDB配置", err.Error())
		return err
	}
	return nil
}
