package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
	"strconv"
)

//******************************建立redis连接池*******************************************************************************
func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(Conf.RedisNetWork, Conf.RedisAddress)
			if err != nil {
				fmt.Println("请检查redis服务器有没有打开/地址是否正确", err.Error())
				glog.Error("请检查redis服务器有没有打开/地址是否正确" + err.Error())

			}
			return c, err
		},
	}

}

//******************************入用户到redis**********************************************************************************
func UserToRedis(user *UserStruct) (bool, error) {
	if user == nil || user.User_id <= 0 {
		return false, errors.New("the user is invalid")
	}
	userStr, err := UserToGob64(user)
	if err != nil {
		return false, err
	}
	reply, err := redisConn.Do("SET", strconv.Itoa(user.User_id), userStr)
	if err != nil || (reply.(string)) != "OK" {
		return false, err
	}
	if (reply.(string)) == "OK" {
		return true, nil
	}
	return false, nil
}

//******************************根据用户名获得一个用户**************************************************************************
func GetUserFromRedis(userID string) (*UserStruct, error) {

	if len(userID) <= 0 {
		return nil, errors.New("the userID is invalid")
	}
	ID, err := strconv.Atoi(userID)
	if ID > lastUserID || err != nil {
		return nil, errors.New("该用户还未注册")
	}
	reply, err := redis.String(redisConn.Do("GET", userID))
	if err != nil {
		fmt.Println("redis.String() err", err.Error())
		glog.Error("redis.String() err" + err.Error())
		return nil, err
	}
	user, err := Gob64ToUser(reply)
	if err != nil {
		fmt.Println("Gob64ToUser err", err.Error())
		glog.Error("Gob64ToUser err" + err.Error())
		return nil, err
	}
	return user, nil
}

//******************************使用gob编码，把UserStruct结构体化为string存进redis中*********************************************
func UserToGob64(user *UserStruct) (string, error) {
	buf := bytes.Buffer{}
	userToGob64 := gob.NewEncoder(&buf)
	err := userToGob64.Encode(user)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

//******************************使用gob编码，把string转化为UserStruct结构体******************************************************
func Gob64ToUser(gob64 string) (*UserStruct, error) {
	user := UserStruct{}
	gob64ToUser, err := base64.StdEncoding.DecodeString(gob64)
	if err != nil {
		fmt.Println("UserToGob64 err", err.Error())
		glog.Error("UserToGob64 err" + err.Error())
		return nil, err
	}
	buf := bytes.Buffer{}
	buf.Write(gob64ToUser)
	d := gob.NewDecoder(&buf)
	err = d.Decode(&user)
	return &user, err
}

//-----------------------------------------连接数据库--------------------------------------------------------------------------
func setupDB(root, password, ip, port, database string) error {
	var err error
	db, err = sql.Open("mysql", root+":"+password+"@tcp("+ip+":"+port+")/"+"mysql"+"?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println("%v", err.Error())
		return err
	}

	crDb := "CREATE DATABASE IF NOT EXISTS  " + database + " DEFAULT CHARSET utf8 COLLATE utf8_general_ci;"
	crStmt, err := db.Prepare(crDb)
	if err != nil {
		fmt.Println("%v", err.Error())
		return err
	}
	_, err = crStmt.Exec()
	if err != nil {
		fmt.Println("%v", err.Error())
		return err
	}
	crStmt.Close()

	db, err = sql.Open("mysql", root+":"+password+"@tcp("+ip+":"+port+")/"+database+"?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println("%v", err.Error())
		return err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("%v", err.Error())
		return err
	}

	SQL := `CREATE TABLE IF NOT EXISTS user_struct
      (
      user_id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
      name varchar(20),
      pwd varchar(20),
      sex int,
      age int,
      intro varchar(100),
      role_id int NOT NULL
      );`

	stmt, err := db.Prepare(SQL)
	if err != nil {
		fmt.Println("%v", err.Error())
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println("%v", err.Error())
		return err
	}
	stmt.Close()

	SQL = "select user_id from user_struct order by user_id desc limit 1" //LIMIT 子句可以被用于强制 SELECT 语句返回指定的记录数。LIMIT 接受一个或两个数字参数。参数必须是一个整数常量。如果给定两个参数，第一个参数指定第一个返回记录行的偏移量，第二个参数指定返回记录行的最大数目。初始记录行的偏移量是 0(而不是 1)
	rows, rowsErr := db.Query(SQL)
	if rowsErr != nil {
		fmt.Println(rowsErr.Error())
		return rowsErr
	}
	for rows.Next() {
		err = rows.Scan(&lastUserID)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	defer rows.Close()
	defer db.Close()

	Xorm, err = xorm.NewEngine("mysql", root+":"+password+"@/"+database+"?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println(err)
		return err
	}

	Xorm.SetMaxIdleConns(500) //set the max number of the connection to the database
	Xorm.SetMaxOpenConns(2000)

	err = Xorm.Ping()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("连接数据库成功")
	return err
}

//*****************************加载数据到缓存************************************************************************************
func dbToCache() {
	var err error
	fmt.Println("执行  dbToCache")
	user := new(UserStruct)
	rows, err := Xorm.Rows(&UserStruct{})
	if err != nil {
		fmt.Println(err)
		glog.Error(" dbToCache()  Xorm.Rows err" + err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() { //遍历数据库，并逐条保存到map中  cacheMap[user.User_id]=*user
		err = rows.Scan(user)
		if err != nil {
			glog.Error(" dbToCache()  rows.Next() err" + err.Error())
			return
		}

		isOk, err2 := UserToRedis(user)
		if isOk == false || err != nil {
			fmt.Println(err2.Error())
			glog.Warning(err2.Error())
			continue
		}
	}

}

//*****************************用户redis跟新到数据库******************************************************************************
func cacheToDB() {
	fmt.Println("执行  cacheToDB")

	redisConn := pool.Get()
	iter := 0 // 迭代器
	var keys []string
	for {
		arr, err := redis.MultiBulk(redisConn.Do("SCAN", iter))
		if err != nil || arr == nil {
			fmt.Println("redisConn.Do SCAN  err", err.Error())
			glog.Error("redisConn.Do SCAN  err" + err.Error())
			break
		}
		iter, err = redis.Int(arr[0], nil)
		if err != nil {
			fmt.Println("redis.Int(arr[0], nil)  err", err)
			glog.Error("redis.Int(arr[0], nil) err", err)
			break
		}
		keys, err = redis.Strings(arr[1], nil)
		if err != nil || keys == nil {
			fmt.Println("redis.Strings(arr[1], nil)  err", err.Error())
			glog.Error("redis.Strings(arr[1], nil) err" + err.Error())
			break
		}

		session := Xorm.NewSession() // add Begin() before any action
		err = session.Begin()
		if err != nil {
			glog.Error(" cacheToDB() session.Begin() err" + err.Error())
			return
		}

		for _, v := range keys {
			user, err := GetUserFromRedis(v)
			if err != nil || user == nil {
				fmt.Println("GetUserFromRedis err", err.Error())
				glog.Error("GetUserFromRedis err" + err.Error())
				continue
			}

			//replace into  1. 首先判断数据是否存在； 2. 如果不存在，则插入；3.如果存在，则更新
			strSQL := "replace into user_struct(user_id,name,pwd,sex,age,intro,role_id) values"
			strSQL = strSQL + "('" + strconv.Itoa(user.User_id) + "','" + user.Name + "','" + user.Pwd + "','" + strconv.Itoa(user.Sex) + "','" + strconv.Itoa(user.Age) + "','" + user.Intro + "','" + strconv.Itoa(user.Role_id) + "')"
			_, err = session.Exec(strSQL)
			if err != nil {
				glog.Error(" cacheToDB() session.Exec err" + err.Error())
				session.Rollback()
				glog.Error(" cacheToDB() session.Exec err" + err.Error())
				continue
			}
		}

		err = session.Commit()
		if err != nil {
			glog.Error("session.Commit() err" + err.Error())
			return
		}
		session.Close()
		if iter == 0 {
			break
		} //迭代遍历完成

	}
	return
}

func GetOnlineUserByID(ID int) *myConn {
	OnlineConn, ok := userList[ID]
	if !ok {
		return nil //运行到这里表明对方不在线
	}
	myConn := &OnlineConn
	return myConn
}
