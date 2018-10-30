package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	_ "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

//chat Listener
type chat struct {
	conn *myConn
}

//******************************************chat**********************************************************************************
func (this chat) onProcess(requestMsg *Msg) {
	if int(requestMsg.GetUserOpt()) != Conf.ChatOP {
		return
	}
	responseMsg := &Msg{
		UserOpt: proto.Int(Conf.ChatOP),
		PersonalMsg: &PersonalMsg{
			SenderID: proto.Int32(requestMsg.PersonalMsg.GetSenderID()),
			RecverID: []int32{0}, //数组简单初始化
		},
	}
	SendTime := time.Now()
	SendTimeStr := SendTime.Format(timeFormat)
	responseMsg.PersonalMsg.SendTime = proto.String(SendTimeStr)
	SrcID := int(requestMsg.PersonalMsg.GetSenderID())
	SrcIDToString := strconv.Itoa(SrcID)
	_, isOk := checkWord(SrcIDToString)
	if isOk == false {
		this.conn.errReturn(responseMsg, "谁发的？", nil, glogInfo)
		return
	}
	isOk = checkToken(SrcID, requestMsg.GetToken())
	if isOk == false {
		responseMsg.UserOpt = proto.Int(Conf.LoginOP)
		this.conn.errReturn(responseMsg, "用户已失效，请重新登录", nil, glogInfo)
		return
	}

	MsgType := int(requestMsg.PersonalMsg.GetMsgType())
	Content := requestMsg.PersonalMsg.GetContent()
	if MsgType == 1 && Content == "" {
		this.conn.errReturn(responseMsg, "你要发什么？不能为空", nil, glogInfo)
		return
	}

	FileName := requestMsg.PersonalMsg.GetFileName()
	if (MsgType != 1) && len(FileName) <= 0 {
		this.conn.errReturn(responseMsg, "你要发什么文件？文件名不能为空", nil, glogInfo)
		return
	}
	FileMD5 := requestMsg.PersonalMsg.GetFileMd5()
	if (MsgType != 1) && len(FileMD5) <= 0 {
		this.conn.errReturn(responseMsg, "你要发的文件标识不能为空", nil, glogInfo)
		return
	}
	FileLen := requestMsg.PersonalMsg.GetFileLen()
	if (MsgType != 1) && len(FileMD5) <= 0 {
		this.conn.errReturn(responseMsg, "你要发的文件大小不能小于0", nil, glogInfo)
		return
	}
	if (MsgType == 1) && len(Content) >= 300 {
		this.conn.errReturn(responseMsg, "一次最多只能发300个字", nil, glogInfo)
		return
	}
	RecverID := requestMsg.PersonalMsg.GetRecverID()
	if len(RecverID) <= 0 {
		this.conn.errReturn(responseMsg, "你要发给谁？", nil, glogInfo)
		return
	}
	if len(RecverID) > 200 {
		this.conn.errReturn(responseMsg, "最多只能给200用户群发", nil, glogInfo)
		return
	}
	responseMsg.PersonalMsg = requestMsg.PersonalMsg
	if MsgType != 1 {
		chatMsg := ChatMsg{} //查看是否已存在这个文件，若存在，则发送，不需要再次上传
		err := ChatCollect.Find(bson.M{"FileMD5": FileMD5}).One(&chatMsg)
		if err != nil {
			responseMsg.PersonalMsg.IsNeedUpload = proto.Bool(true)
			this.conn.errReturn(responseMsg, "请先上传你要发的文件", err, glogInfo)
			return
		}
	}
	for _, v := range RecverID {
		vToString := strconv.Itoa(int(v))
		_, isOk := checkWord(vToString)
		if isOk == false {
			this.conn.errReturn(responseMsg, "目标ID不合法", nil, glogErr)
			return
		}
		/* 1文字word
		   2文件file名、文件MD5
		   插入数据库
		*/
		switch MsgType {
		case 1:
			err := ChatCollect.Insert(&ChatMsg{Src: SrcID, Dst: int(v), SendTime: SendTime, Content: Content, MsgType: 1})
			if err != nil {
				this.conn.errReturn(responseMsg, "服务器繁忙，请稍后", err, glogErr)
				fmt.Println("11111116")
				continue
			}
		default:
			err := ChatCollect.Insert(&ChatMsg{Src: SrcID, Dst: int(v), SendTime: SendTime, MsgType: MsgType, FileName: FileName, FileMD5: FileMD5, FileLen: FileLen})
			if err != nil {
				this.conn.errReturn(responseMsg, "服务器繁忙，请稍后", err, glogErr)
				return
			}
		}
		FriendConn := GetOnlineUserByID(int(v))
		if FriendConn == nil { //对方不在线，直接跳过下面的步骤
			continue
		}
		//运行到这里说明目标在线  直接发送
		newSendMsg := &Msg{
			UserOpt: proto.Int32(requestMsg.GetUserOpt()),
			PersonalMsg: &PersonalMsg{
				SenderID: proto.Int32(requestMsg.PersonalMsg.GetSenderID()),
				RecverID: []int32{v},
				SendTime: proto.String(SendTimeStr),
				Content:  proto.String(Content),
				MsgType:  proto.Int32(requestMsg.PersonalMsg.GetMsgType()),
				FileName: proto.String(FileName),
				FileMd5:  proto.String(FileMD5),
				FileLen:  proto.String(FileLen),
			},
		}
		FriendConn.write(newSendMsg)
	}
	responseMsg.OptResult = proto.String(SUCESS_RESULT)
	this.conn.write(responseMsg)
	fmt.Println("chat")
}

//readed Listener
type readed struct {
	conn *myConn
}

//******************************************chat--readed**********************************************************************************
func (this readed) onProcess(requestMsg *Msg) {
	if int(requestMsg.GetUserOpt()) != Conf.ReadedOP {
		return
	}
	responseMsg := &Msg{
		UserOpt: proto.Int32(requestMsg.GetUserOpt()),
		PersonalMsg: &PersonalMsg{
			RecverID: []int32{0}, //数组简单初始化
		},
	}
	readTime := time.Now()
	readTimeStr := readTime.Format(timeFormat)
	fmt.Println(readTimeStr)
	responseMsg.PersonalMsg.ReadTime = &readTimeStr
	SrcID := int(requestMsg.PersonalMsg.GetSenderID())
	SrcIDToString := strconv.Itoa(SrcID)
	myID := int(requestMsg.User.GetUserID())

	isOk := checkToken(myID, requestMsg.GetToken())
	if isOk == false {
		responseMsg.UserOpt = proto.Int(Conf.LoginOP)
		this.conn.errReturn(responseMsg, "用户已失效，请重新登录", nil, glogInfo)
		return
	}
	_, isOk = checkWord(SrcIDToString)
	if isOk == false {
		this.conn.errReturn(responseMsg, "谁发给你的这条信息，ID不合法", nil, glogInfo)
		return
	}

	_, err := UserCollect.Upsert(bson.M{"ID": int(requestMsg.PersonalMsg.RecverID[0])}, &UserReadTime{ID: int(requestMsg.PersonalMsg.RecverID[0]), ReadTime: readTime})
	if err != nil {
		this.conn.errReturn(responseMsg, "服务器繁忙，请稍后", err, glogErr)
		return
	}
	responseMsg.OptResult = proto.String(SUCESS_RESULT)
	responseMsg.PersonalMsg.SenderID = proto.Int(SrcID)
	responseMsg.PersonalMsg.RecverID = requestMsg.PersonalMsg.RecverID
	this.conn.write(responseMsg)
	FriendConn := GetOnlineUserByID(SrcID)
	fmt.Println("readed")
	if FriendConn == nil {
		return
	}
	newSendMsg := &Msg{ //运行到这里表明源发送者在线，发送已阅读  给源
		UserOpt: proto.Int32(requestMsg.GetUserOpt()),
		PersonalMsg: &PersonalMsg{
			SenderID: proto.Int32(requestMsg.PersonalMsg.GetSenderID()),
			RecverID: []int32{requestMsg.PersonalMsg.RecverID[0]},
			SendTime: proto.String(requestMsg.PersonalMsg.GetSendTime()),
			ReadTime: proto.String(responseMsg.PersonalMsg.GetReadTime()),
			IsReaded: proto.Bool(requestMsg.PersonalMsg.GetIsReaded()),
		},
	}
	FriendConn.write(newSendMsg)
}

//offLineMsg Listener
type offLineMsg struct {
	conn *myConn
}

//******************************************offLineMsg************************************************************************
func (this offLineMsg) onProcess(requestMsg *Msg) {
	if int(requestMsg.GetUserOpt()) != Conf.OffLineMsgOP {
		return
	}
	responseMsg := &Msg{
		UserOpt: proto.Int32(requestMsg.GetUserOpt()),
	}
	if requestMsg.User.GetUserID() <= 0 {
		*responseMsg.ErrMsg = "你是谁？"
		glog.Warning("你是谁？")
		this.conn.write(responseMsg)
		return
	}
	ID := int(requestMsg.User.GetUserID())
	isOk := checkToken(ID, requestMsg.GetToken())
	if isOk == false {
		responseMsg.UserOpt = proto.Int(Conf.LoginOP)
		this.conn.errReturn(responseMsg, "用户已失效，请重新登录", nil, glogInfo)
		return
	}
	userReadTime := UserReadTime{}
	err := UserCollect.Find(bson.M{"ID": ID}).One(&userReadTime)
	if err != nil || userReadTime.ID == 0 {
		fmt.Println("该用户在新客户端登录，无版本号")
	}

	lastReadTime := userReadTime.ReadTime //查出大于自己最后一次读的时间  的所有聊天
	chatMsg := ChatMsg{}
	iter := ChatCollect.Find(bson.M{"Dst": ID, "SendTime": bson.M{"$gt": lastReadTime}}).Iter()
	for iter.Next(&chatMsg) {
		newMsg := &PersonalMsg{
			SenderID: proto.Int32(int32(chatMsg.Src)),
			RecverID: []int32{int32(chatMsg.Dst)},
			SendTime: proto.String(chatMsg.SendTime.Format(timeFormat)),
		}

		switch chatMsg.MsgType {
		case 1:
			newMsg.MsgType = proto.Int32(1)
			newMsg.Content = proto.String(chatMsg.Content)
		default:
			newMsg.MsgType = proto.Int(chatMsg.MsgType)
			newMsg.FileName = proto.String(chatMsg.FileName)
			newMsg.FileMd5 = proto.String(chatMsg.FileMD5)
			newMsg.FileLen = proto.String(chatMsg.FileLen)
		}
		responseMsg.OfflineMsg = append(responseMsg.OfflineMsg, newMsg)
	}
	fmt.Println("用户登录:" + strconv.Itoa(ID) + "，发送离线消息成功")
	responseMsg.OptResult = proto.String(SUCESS_RESULT)
	this.conn.write(responseMsg)
}
