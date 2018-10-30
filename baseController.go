package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	_ "gopkg.in/mgo.v2"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

//选择以哪种方式打印到3级日志文件   错误日志、 警告日志、  信息日志
func (conn *myConn) errReturn(responseMsg *Msg, errMsg string, err error, choiceGlog int) {
	responseMsg.ErrMsg = proto.String(errMsg) //返回给顾客的可以是  服务器繁忙等等   不一定非要具体的err
	responseMsg.OptResult = proto.String(FAILED_RESULT)
	fmt.Println("err:", err, "返回给客户端的errMsg：", errMsg)
	switch choiceGlog {
	case glogErr:
		glog.Error(err)
	case glogwarn:
		glog.Warning(err)
	case glogInfo:
		glog.Info(err)
	}
	conn.write(responseMsg)
}

//Login Listener
type login struct {
	conn *myConn
}

//******************************************login***************************************************************************************
func (this login) onProcess(requestMsg *Msg) {
	if int(requestMsg.GetUserOpt()) != Conf.LoginOP {
		return
	}
	responseMsg := &Msg{
		UserOpt: proto.Int32(requestMsg.GetUserOpt()),
	}

	requestID := int(*requestMsg.User.UserID)
	requestIDStr := strconv.Itoa(requestID)
	_, isOk := checkWord(requestIDStr)
	if isOk == false {
		this.conn.errReturn(responseMsg, "账号格式不对", nil, glogInfo)
		return
	}
	_, isOk = checkWord(requestMsg.User.GetUserPwd())
	if isOk == false {
		this.conn.errReturn(responseMsg, "密码格式不对", nil, glogInfo)
		return
	}
	user, err := GetUserFromRedis(requestIDStr)
	if err != nil || user == nil {
		this.conn.errReturn(responseMsg, "该用户不存在", err, glogInfo)
		return
	}
	if user.Pwd != requestMsg.User.GetUserPwd() {
		this.conn.errReturn(responseMsg, "密码不对", nil, glogInfo)
		return
	}

	responseMsg.OptResult = proto.String(SUCESS_RESULT)
	oldConn := GetOnlineUserByID(requestID)
	if oldConn != nil {
		newMsg := &Msg{
			UserOpt:   proto.Int(Conf.LoginOP),
			OptResult: proto.String(FAILED_RESULT),
			ErrMsg:    proto.String("用户在其他地方登录"),
		}
		oldConn.write(newMsg) //断开旧连接
		delete(userList, oldConn.id)
	}
	this.conn.lock.Lock()
	this.conn.id = requestID
	token := getToken(requestIDStr, user.Pwd)
	// addListener(this.conn) //登陆后可以正常 聊天 之类的    测试时候注释掉本行  ！！！！
	this.conn.token = token
	responseMsg.Token = &token
	addListener(this.conn)
	this.conn.lock.Unlock()
	userList[requestID] = *this.conn //保存新连接
	userInfo := &User{
		UserID:     proto.Int32(int32(user.User_id)),
		NickName:   proto.String(user.Name),
		UesrIntro:  proto.String(user.Intro),
		UserRoleID: proto.Int32(int32(user.Role_id)),
		UserDetail: &UserDetail{
			Age: proto.Int32(int32(user.Age)),
		},
	}
	responseMsg.User = userInfo
	userList[requestID] = *this.conn
	fmt.Println("用户" + requestIDStr + "已登录")
	this.conn.write(responseMsg)
	showOnLines(userList)

}

func getToken(id, pwd string) string {
	t := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(t, 10)+id+pwd)
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token[0:11]
}

func checkToken(id int, token string) bool {
	if len(token) <= 0 {
		fmt.Println("token is nil")
		return false
	}
	if token != userList[id].token {
		return false
	}
	return true
}

//regis Listener
type regis struct {
	conn *myConn
}

//******************************************regis****************************************************************************************
func (this regis) onProcess(requestMsg *Msg) {
	if int(requestMsg.GetUserOpt()) != Conf.RegisOP {
		return
	}
	responseMsg := &Msg{
		UserOpt: proto.Int32(requestMsg.GetUserOpt()),
	}

	_, isOk := checkWord(requestMsg.User.GetUserPwd())
	if isOk == false {
		this.conn.errReturn(responseMsg, "密码格式不对", nil, glogInfo)
		return
	}

	if requestMsg.User.GetUserRoleID() > 3 || requestMsg.User.GetUserRoleID() < 1 {
		requestMsg.User.UserRoleID = proto.Int(1) //权限号默认为1
	}

	if requestMsg.User.GetUesrIntro() == "" {
		DefaultIntro := "该用户太懒，什么都没有下"
		requestMsg.User.UesrIntro = &DefaultIntro
	}

	lastUserID = lastUserID + 1
	user := UserStruct{lastUserID, requestMsg.User.GetNickName(), requestMsg.User.GetUserPwd(), 1, 20, requestMsg.User.GetUesrIntro(), int(requestMsg.User.GetUserRoleID())}
	Bool, err := UserToRedis(&user)
	if Bool == false || err != nil {
		this.conn.errReturn(responseMsg, "由于服务器问题，注册失败", err, glogInfo)
		return
	}

	responseMsg.OptResult = proto.String(SUCESS_RESULT)
	lastUserID32 := int32(lastUserID)
	*responseMsg.User.UserID = lastUserID32
	IDStr := strconv.Itoa(lastUserID)
	fmt.Println("新用户：" + IDStr + "注册成功")
	token := getToken(IDStr, user.Pwd) //注册成功即登录
	this.conn.lock.Lock()
	addListener(this.conn) //登陆后可以正常 聊天 之类的
	this.conn.id = lastUserID
	this.conn.token = token
	this.conn.lock.Unlock()
	responseMsg.Token = &token
	userInfo := &User{
		UserID:     proto.Int32(int32(user.User_id)),
		NickName:   proto.String(user.Name),
		UesrIntro:  proto.String(user.Intro),
		UserRoleID: proto.Int32(int32(user.Role_id)),
		UserDetail: &UserDetail{
			Age: proto.Int32(int32(user.Age)),
		},
	}
	responseMsg.User = userInfo
	userList[lastUserID] = *this.conn
	fmt.Println("用户" + IDStr + "已登录")
	showOnLines(userList)
	addFriendGroup(lastUserID, 0, "我的好友")
	addFriendGroup(lastUserID, 0, "我的家人")
	addListener(this.conn)
	this.conn.write(responseMsg)
}

//*****************************聊天缓存的Key：源ID++目标ID+时间************************************************************************
func getKey(SrcID, DstID int, SendTime string) string {
	SrcIDString := strconv.Itoa(SrcID)
	DstIDString := strconv.Itoa(DstID)
	newKeyString := SrcIDString + "142857" + DstIDString + "142857" + SendTime
	return newKeyString
}

//**********************************   从键盘读取广播信息    \n 结尾**************************************************************************
func getMsg() (msg string) {
	reader := bufio.NewReader(os.Stdin)
	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
		glog.Error(" reader.ReadString err" + err.Error())
		return
	}
	msg = strings.Replace(msg, "\n", "", -1)
	broadcastChannel <- msg
	return
}
