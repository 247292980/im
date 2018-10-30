package main

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	_ "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

//addFriendRequest Listener
type addFriendRequest struct {
	conn *myConn
}

//******************************************addFriendRequest*************************************************************
func (this addFriendRequest) onProcess(requestMsg *Msg) {
	if int(requestMsg.GetUserOpt()) != Conf.AddFriendRequestOP {
		return
	}
	responseMsg := &Msg{
		UserOpt: proto.Int32(requestMsg.GetUserOpt()),
	}
	if requestMsg.User.GetUserID() <= 0 {
		this.conn.errReturn(responseMsg, "你是谁？", nil, glogInfo)
		return
	}
	if requestMsg.Friends[0].GetUserID() <= 0 {
		this.conn.errReturn(responseMsg, "你要加谁？", nil, glogInfo)
		return
	}
	ID := int(requestMsg.User.GetUserID())
	FrdID := int(requestMsg.Friends[0].GetUserID())
	if ID == FrdID {
		this.conn.errReturn(responseMsg, "不能加自己啊", nil, glogInfo)
		return
	}
	isOk := checkToken(ID, requestMsg.GetToken())
	if isOk == false {
		responseMsg.UserOpt = proto.Int(Conf.LoginOP)
		this.conn.errReturn(responseMsg, "用户已失效，请重新登录", nil, glogInfo)
		return
	}
	friendInfo, _ := GetUserFromRedis(strconv.Itoa(FrdID)) //从redis取出详细信息
	if friendInfo == nil {
		this.conn.errReturn(responseMsg, "该用户不存在", nil, glogInfo)
		return
	}
	myNick := responseMsg.User.GetNickName()
	friend := Friend{}
	FriendCollect.Find(bson.M{"ID": ID, "FrdID": FrdID}).One(&friend)
	if friend.ID > 0 {
		this.conn.errReturn(responseMsg, "对方已经是您的好友", nil, glogInfo)
		return
	} //不仅要查看是否已经是好友 还要在离线请求里面看  是否已经有这条请求记录   避免重复插入出错
	friendRequest := FriendRequest{}
	FrdRequesCollection.Find(bson.M{"SrcID": ID, "DstID": FrdID}).One(&friendRequest)
	if friendRequest.SrcID > 0 {
		this.conn.errReturn(responseMsg, "已经向对方发送好友添加请求，请您耐心等待", nil, glogInfo)
		return
	} //到这里判断结束  可以向对方发送好友申请了
	FriendConn := GetOnlineUserByID(FrdID)
	RequestTime := time.Now()
	if FriendConn == nil { //不在线-存在请求数据库
		err := FrdRequesCollection.Insert(bson.M{"SrcID": ID, "DstID": FrdID, "SrcName": myNick, "RequestTime": RequestTime})
		if err != nil {
			this.conn.errReturn(responseMsg, "服务器繁忙，请稍后", err, glogErr)
			return
		}
	} else { //在线-直接向请求发送请求
		newMsg := &Msg{}
		newMsg.UserOpt = requestMsg.UserOpt
		newMsg.Friends = requestMsg.Friends
		newMsg.User = requestMsg.User
		FriendConn.write(newMsg)
	}

	responseMsg.OptResult = proto.String(SUCESS_RESULT)
	this.conn.write(responseMsg)
	fmt.Println("增加好友请求：", responseMsg)
}

//addFriendResponse Listener
type addFriendResponse struct {
	conn *myConn
}

//******************************************addFriendResponse*************************************************************
func (this addFriendResponse) onProcess(requestMsg *Msg) {
	if int(requestMsg.GetUserOpt()) != Conf.AddFriendResponseOP {
		return
	}
	responseMsg := &Msg{
		UserOpt: proto.Int32(requestMsg.GetUserOpt()),
	}
	if requestMsg.User.GetUserID() <= 0 {
		this.conn.errReturn(responseMsg, "你是谁？", nil, glogInfo)
		return
	}
	if requestMsg.Friends[0].GetUserID() <= 0 {
		this.conn.errReturn(responseMsg, "谁要加你？", nil, glogInfo)
		return
	}
	ID := int(requestMsg.User.GetUserID())
	FrdID := int(requestMsg.Friends[0].GetUserID())
	IDStr := strconv.Itoa(ID)
	FrdIDStr := strconv.Itoa(FrdID)
	if ID == FrdID {
		this.conn.errReturn(responseMsg, "不能加自己啊", nil, glogInfo)
		return
	}
	friend := Friend{}
	FriendCollect.Find(bson.M{"ID": ID, "FrdID": FrdID}).One(&friend)
	if friend.ID > 0 {
		this.conn.errReturn(responseMsg, "对方已经是您的好友", nil, glogInfo)
		return
	}
	isOk := checkToken(ID, requestMsg.GetToken())
	if isOk == false {
		responseMsg.UserOpt = proto.Int(Conf.LoginOP)
		this.conn.errReturn(responseMsg, "用户已失效，请重新登录", nil, glogInfo)
		return
	}
	FriendConn := GetOnlineUserByID(FrdID)
	if requestMsg.Friends[0].GetIsAgree() == false { //对方拒绝
		if FriendConn != nil { //在线-发送拒绝信息
			newMsg := &Msg{}
			newMsg.UserOpt = requestMsg.UserOpt
			newMsg.Friends = requestMsg.Friends
			newMsg.User = requestMsg.User
			FriendConn.write(newMsg) //发送拒绝给申请者
		}
		responseMsg.User = requestMsg.User
		responseMsg.Friends = requestMsg.Friends
		responseMsg.OptResult = proto.String(SUCESS_RESULT)
		this.conn.write(responseMsg)
		return
	} //运行到这里表示 同意添加好友
	myInfo, _ := GetUserFromRedis(IDStr)
	friendInfo, err := GetUserFromRedis(FrdIDStr)
	if err != nil {
		this.conn.errReturn(responseMsg, "该用户不存在", nil, glogInfo)
		return
	}
	addTime := time.Now()
	err = FriendCollect.Insert(bson.M{"ID": ID, "FrdID": FrdID, "Remark": myInfo.Name, "AddTime": addTime, "Group": 0},
		bson.M{"ID": FrdID, "FrdID": ID, "Remark": friendInfo.Name, "AddTime": addTime, "Group": 0})
	if err != nil {
		this.conn.errReturn(responseMsg, "服务器繁忙，请稍后", err, glogErr)
		return
	} //到这里增加好友成功   一起发送好友信息  默认是在分组0 我的好友
	friend_item := &User{
		UserID:   proto.Int(FrdID),
		NickName: proto.String(friendInfo.Name),
		Remark:   proto.String(friendInfo.Name),
		ListNO:   proto.Int(0),
		IsAgree:  proto.Bool(true),
	}
	my_item := &User{
		UserID:   proto.Int(ID),
		NickName: proto.String(myInfo.Name),
		Remark:   proto.String(myInfo.Name),
		ListNO:   proto.Int(0),
		IsAgree:  proto.Bool(true),
	}
	if FriendConn != nil { //在线-发送同意信息
		newMsg := &Msg{}
		newMsg.UserOpt = requestMsg.UserOpt
		newMsg.Friends = append(newMsg.Friends, my_item)
		FriendConn.write(newMsg) //发送同意信息给申请者
	}
	responseMsg.Friends = append(responseMsg.Friends, friend_item) //返回给同意者
	responseMsg.OptResult = proto.String(SUCESS_RESULT)
	this.conn.write(responseMsg)
	fmt.Println("增加好友：", responseMsg)
}

//deltFriend Listener
type deltFriend struct {
	conn *myConn
}

//******************************************deltFriend*************************************************************
func (this deltFriend) onProcess(requestMsg *Msg) {
	if int(requestMsg.GetUserOpt()) != Conf.DeltFriendOP {
		return
	}
	responseMsg := &Msg{
		UserOpt: proto.Int32(requestMsg.GetUserOpt()),
	}
	if requestMsg.User.GetUserID() <= 0 {
		this.conn.errReturn(responseMsg, "你是谁？", nil, glogInfo)
		return
	}
	if requestMsg.Friends[0].GetUserID() <= 0 {
		this.conn.errReturn(responseMsg, "你要删谁？", nil, glogInfo)
		return
	}
	ID := int(requestMsg.User.GetUserID())
	FrdID := int(requestMsg.Friends[0].GetUserID())
	if ID == FrdID {
		this.conn.errReturn(responseMsg, "不能删自己啊", nil, glogInfo)
		return
	}
	isOk := checkToken(ID, requestMsg.GetToken())
	if isOk == false {
		responseMsg.UserOpt = proto.Int(Conf.LoginOP)
		this.conn.errReturn(responseMsg, "用户已失效，请重新登录", nil, glogInfo)
		return
	}
	friend := Friend{}
	err := FriendCollect.Find(bson.M{"ID": ID, "FrdID": FrdID}).One(&friend)
	if err != nil {
		this.conn.errReturn(responseMsg, "对方不是您的好友", err, glogInfo)
		return
	}
	err = FriendCollect.Remove(bson.M{"ID": ID, "FrdID": FrdID})
	err = FriendCollect.Remove(bson.M{"ID": FrdID, "FrdID": ID})
	if err != nil {
		this.conn.errReturn(responseMsg, "服务器繁忙，请稍后", err, glogErr)
		return
	} //到这里删除成功
	FriendConn := GetOnlineUserByID(FrdID)
	if FriendConn != nil { //在线-发送被删除 同步信息
		newMsg := &Msg{}
		newMsg.UserOpt = requestMsg.UserOpt
		newMsg.Friends = append(responseMsg.Friends, requestMsg.User)
		FriendConn.write(newMsg)
	}
	responseMsg.OptResult = proto.String(SUCESS_RESULT)
	this.conn.write(responseMsg)
	fmt.Println("删除好友：", responseMsg)
}

//getFriends Listener
type getFriends struct {
	conn *myConn
}

//******************************************getFriends*************************************************************
func (this getFriends) onProcess(requestMsg *Msg) {
	if int(requestMsg.GetUserOpt()) != Conf.GetFriendListOP {
		return
	}
	responseMsg := &Msg{
		UserOpt: proto.Int32(requestMsg.GetUserOpt()),
	}
	if requestMsg.User.GetUserID() <= 0 {
		this.conn.errReturn(responseMsg, "你是谁？", nil, glogInfo)
		return
	}
	ID := int(requestMsg.User.GetUserID())
	isOk := checkToken(ID, requestMsg.GetToken())
	if isOk == false {
		responseMsg.UserOpt = proto.Int(Conf.LoginOP)
		this.conn.errReturn(responseMsg, "用户已失效，请重新登录", nil, glogInfo)
		return
	} //到这里判断结束  可以取出好友信息了
	friendGroup := FriendGroup{} //从mongoDB取出好友分组
	iterGroup := FrdGroupCollect.Find(bson.M{"UID": ID}).Iter()
	for iterGroup.Next(&friendGroup) {
		friendList := &FriendList{
			ListNO:   proto.Int(friendGroup.ListNO),
			ListName: proto.String(friendGroup.ListName),
		}
		responseMsg.FriendLists = append(responseMsg.FriendLists, friendList)
	}
	defer iterGroup.Close()
	friend := Friend{} //从mongoDB取出好友
	iter := FriendCollect.Find(bson.M{"ID": ID}).Iter()
	for iter.Next(&friend) {
		user, _ := GetUserFromRedis(strconv.Itoa(friend.FrdID)) //从redis取出详细信息
		friendDetail := &User{
			UserID:   proto.Int32(int32(user.User_id)),
			NickName: proto.String(user.Name),
			Remark:   proto.String(friend.Remark),
			ListNO:   proto.Int(friend.Group),
		}
		responseMsg.Friends = append(responseMsg.Friends, friendDetail)
	}
	defer iter.Close()
	friendRequest := FriendRequest{} //从mongoDB取出离线好友添加请求
	iterRequest := FrdRequesCollection.Find(bson.M{"DstID": ID}).Iter()
	for iterRequest.Next(&friendRequest) {
		requestfriend_item := &User{
			UserID:   proto.Int(friendRequest.SrcID),
			NickName: proto.String(friendRequest.SrcName),
			IsAgree:  proto.Bool(false),
		}
		responseMsg.NewFriendRequest = append(responseMsg.NewFriendRequest, requestfriend_item)
	} //到这里已经全部取出  取出后从离线好友申请表中删除记录
	FrdRequesCollection.RemoveAll(bson.M{"DstID": ID})
	defer iterRequest.Close()
	responseMsg.OptResult = proto.String(SUCESS_RESULT)
	this.conn.write(responseMsg)
	fmt.Println("发送好友列表：", responseMsg)
}

type remarkFrindListener struct {
	conn *myConn
}

//******************************************remarkFrind*************************************************************
func (this remarkFrindListener) onProcess(requestMsg *Msg) {
	if int(requestMsg.GetUserOpt()) != Conf.RemarkFriendOP {
		return
	}
	responseMsg := &Msg{
		UserOpt: proto.Int32(requestMsg.GetUserOpt()),
	}
	if requestMsg.User.GetUserID() <= 0 {
		this.conn.errReturn(responseMsg, "你是谁？", nil, glogInfo)
		return
	}
	if requestMsg.Friends[0].GetUserID() <= 0 {
		this.conn.errReturn(responseMsg, "你要改谁的备注？", nil, glogInfo)
		return
	}
	ID := int(requestMsg.User.GetUserID())
	FrdID := int(requestMsg.Friends[0].GetUserID())
	Remark := requestMsg.Friends[0].GetRemark()
	if ID == FrdID {
		this.conn.errReturn(responseMsg, "不能备注自己啊", nil, glogInfo)
		return
	}
	isOk := checkToken(ID, requestMsg.GetToken())
	if isOk == false {
		responseMsg.UserOpt = proto.Int(Conf.LoginOP)
		this.conn.errReturn(responseMsg, "用户已失效，请重新登录", nil, glogInfo)
		return
	}
	friend := Friend{}
	err := FriendCollect.Find(bson.M{"ID": ID, "FrdID": FrdID}).One(&friend)
	if err != nil {
		this.conn.errReturn(responseMsg, "对方不是您的好友", err, glogInfo)
		return
	}
	if friend.Remark == Remark {
		this.conn.errReturn(responseMsg, "跟以前备注一样，未修改", nil, glogInfo)
		return
	} //判断到这里   可以修改备注了
	err = FriendCollect.Update(bson.M{"ID": ID, "FrdID": FrdID},
		bson.M{"$set": bson.M{"Remark": Remark}})
	if err != nil {
		this.conn.errReturn(responseMsg, "服务器繁忙，请稍后", err, glogErr)
		return
	}
	responseMsg.Friends = requestMsg.Friends
	responseMsg.OptResult = proto.String(SUCESS_RESULT)
	this.conn.write(responseMsg)
	fmt.Println("修改好友备注：", responseMsg)
}

type addFriendGroupListener struct {
	conn *myConn
}

//******************************************addFriendGroup*************************************************************
func (this addFriendGroupListener) onProcess(requestMsg *Msg) {
	if int(requestMsg.GetUserOpt()) != Conf.AddFriendGroupOP {
		return
	}
	responseMsg := &Msg{
		UserOpt: proto.Int32(requestMsg.GetUserOpt()),
	}
	ID := int(requestMsg.User.GetUserID())
	if ID <= 0 {
		this.conn.errReturn(responseMsg, "你是谁？", nil, glogInfo)
		return
	}
	isOk := checkToken(ID, requestMsg.GetToken())
	if isOk == false {
		responseMsg.UserOpt = proto.Int(Conf.LoginOP)
		this.conn.errReturn(responseMsg, "用户已失效，请重新登录", nil, glogInfo)
		return
	}
	ListNO := int(requestMsg.FriendLists[0].GetListNO())
	ListName := requestMsg.FriendLists[0].GetListName()
	err := addFriendGroup(ID, ListNO, ListName)
	if err != nil {
		this.conn.errReturn(responseMsg, err.Error(), nil, glogInfo)
	}
	responseMsg.OptResult = proto.String(SUCESS_RESULT)
	responseMsg.FriendLists = requestMsg.FriendLists
	this.conn.write(responseMsg)
	fmt.Println("增加好友分组：", responseMsg)
}
func addFriendGroup(ID, ListNO int, ListName string) error {
	if ID <= 0 {
		return errors.New("the userID cannot less than 0")
	}
	if ListNO < 0 {
		return errors.New("你要加的组是什么？")
	}
	if len(ListName) <= 0 {
		return errors.New("你要加的组名不能为空？")
	}
	friendGroup := FriendGroup{}
	err := FrdGroupCollect.Find(bson.M{"UID": ID, "ListNO": ListNO}).One(&friendGroup)
	if err == nil {
		return errors.New("该分组已存在！")
	}
	err = FrdGroupCollect.Insert(bson.M{"UID": ID, "ListNO": ListNO, "ListName": ListName})
	if err != nil {
		return err
	}
	return nil
}

type remarkFrindGroupListener struct {
	conn *myConn
}

//******************************************remarkFrindGroup*************************************************************
func (this remarkFrindGroupListener) onProcess(requestMsg *Msg) {
	if int(requestMsg.GetUserOpt()) != Conf.RemarFriendkGroupOP {
		return
	}
	responseMsg := &Msg{
		UserOpt: proto.Int32(requestMsg.GetUserOpt()),
	}
	ID := int(requestMsg.User.GetUserID())
	if ID <= 0 {
		this.conn.errReturn(responseMsg, "你是谁？", nil, glogInfo)
		return
	}
	ListNO := int(requestMsg.FriendLists[0].GetListNO())
	newListName := requestMsg.FriendLists[0].GetListName()
	if ListNO < 0 {
		this.conn.errReturn(responseMsg, "你要备注哪个组？", nil, glogInfo)
		return
	}
	isOk := checkToken(ID, requestMsg.GetToken())
	if isOk == false {
		responseMsg.UserOpt = proto.Int(Conf.LoginOP)
		this.conn.errReturn(responseMsg, "用户已失效，请重新登录", nil, glogInfo)
		return
	}
	friendGroup := FriendGroup{}
	err := FrdGroupCollect.Find(bson.M{"UID": ID, "ListNO": ListNO}).One(&friendGroup)
	if err != nil {
		this.conn.errReturn(responseMsg, "您没有该分组，你是要修改哪个分组名呢？", err, glogInfo)
		return
	}
	if friendGroup.ListName == newListName {
		this.conn.errReturn(responseMsg, "跟以前分组名一样，未修改", nil, glogInfo)
		return
	} //判断到这里   可以修改备注了

	err = FrdGroupCollect.Update(bson.M{"UID": ID, "ListNO": ListNO},
		bson.M{"$set": bson.M{"ListName": newListName}})
	if err != nil {
		this.conn.errReturn(responseMsg, "服务器繁忙，请稍后", err, glogErr)
		return
	}
	responseMsg.FriendLists = requestMsg.FriendLists
	responseMsg.OptResult = proto.String(SUCESS_RESULT)
	this.conn.write(responseMsg)
	fmt.Println("修改分组名：", responseMsg)
}

type deltFrindGroupListener struct {
	conn *myConn
}

//******************************************deltFrindGroup*************************************************************
func (this deltFrindGroupListener) onProcess(requestMsg *Msg) {
	if int(requestMsg.GetUserOpt()) != Conf.DetFriendGroupOP {
		return
	}
	responseMsg := &Msg{
		UserOpt: proto.Int32(requestMsg.GetUserOpt()),
	}
	ID := int(requestMsg.User.GetUserID())
	if ID <= 0 {
		this.conn.errReturn(responseMsg, "你是谁？", nil, glogInfo)
		return
	}
	isOk := checkToken(ID, requestMsg.GetToken())
	if isOk == false {
		responseMsg.UserOpt = proto.Int(Conf.LoginOP)
		this.conn.errReturn(responseMsg, "用户已失效，请重新登录", nil, glogInfo)
		return
	}
	ListNO := int(requestMsg.FriendLists[0].GetListNO())
	if ListNO < 0 {
		this.conn.errReturn(responseMsg, "你要删哪个组？", nil, glogInfo)
		return
	}
	if ListNO == 0 || ListNO == 1 {
		this.conn.errReturn(responseMsg, "这是基本组，不能删哦亲", nil, glogInfo)
		return
	}
	friendGroup := FriendGroup{}
	err := FrdGroupCollect.Find(bson.M{"UID": ID, "ListNO": ListNO}).One(&friendGroup)
	if err != nil {
		this.conn.errReturn(responseMsg, "删除失败，不存在该分组", err, glogInfo)
		return
	} //到这里判断结束 删除分组，该分组所有好友移到0：我的好友
	_, err = FriendCollect.UpdateAll(bson.M{"ID": ID, "Group": ListNO},
		bson.M{"$set": bson.M{"Group": 0}})
	if err != nil {
		this.conn.errReturn(responseMsg, "服务器繁忙，请稍后", err, glogErr)
		return
	}
	err = FrdGroupCollect.Remove(bson.M{"UID": ID, "ListNO": ListNO})
	if err != nil {
		this.conn.errReturn(responseMsg, "服务器繁忙，请稍后", err, glogErr)
		return
	}
	responseMsg.OptResult = proto.String(SUCESS_RESULT)
	this.conn.write(responseMsg)
	fmt.Println("删除好友分组：", responseMsg)
}
