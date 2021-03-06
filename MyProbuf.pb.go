// Code generated by protoc-gen-go.
// source: MyProbuf.proto
// DO NOT EDIT!

/*
Package MyProbuf is a generated protocol buffer package.

It is generated from these files:
	MyProbuf.proto

It has these top-level messages:
	Msg
	User
	UserDetail
	UserCustomAttr
	FriendList
	SearchInfo
	Group
	GroupNumber
	GroupMsg
	PersonalMsg
	GroupNotice
*/
package main

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Msg struct {
	Sign             *string        `protobuf:"bytes,1,req,name=sign" json:"sign,omitempty"`
	DataLen          *string        `protobuf:"bytes,2,req,name=dataLen" json:"dataLen,omitempty"`
	UserOpt          *int32         `protobuf:"varint,3,req,name=userOpt" json:"userOpt,omitempty"`
	OptResult        *string        `protobuf:"bytes,4,opt,name=optResult" json:"optResult,omitempty"`
	ReceiveResult    *string        `protobuf:"bytes,5,opt,name=receiveResult" json:"receiveResult,omitempty"`
	User             *User          `protobuf:"bytes,6,opt,name=user" json:"user,omitempty"`
	Friends          []*User        `protobuf:"bytes,7,rep,name=friends" json:"friends,omitempty"`
	FriendLists      []*FriendList  `protobuf:"bytes,8,rep,name=friendLists" json:"friendLists,omitempty"`
	NewFriendRequest []*User        `protobuf:"bytes,9,rep,name=newFriendRequest" json:"newFriendRequest,omitempty"`
	Groups           []*Group       `protobuf:"bytes,10,rep,name=groups" json:"groups,omitempty"`
	GroupMsg         []*GroupMsg    `protobuf:"bytes,11,rep,name=groupMsg" json:"groupMsg,omitempty"`
	OfflineMsg       []*PersonalMsg `protobuf:"bytes,12,rep,name=OfflineMsg,json=offlineMsg" json:"OfflineMsg,omitempty"`
	PersonalMsg      *PersonalMsg   `protobuf:"bytes,13,opt,name=personalMsg" json:"personalMsg,omitempty"`
	SrchInfo         *SearchInfo    `protobuf:"bytes,14,opt,name=srchInfo" json:"srchInfo,omitempty"`
	ErrMsg           *string        `protobuf:"bytes,15,opt,name=ErrMsg,json=errMsg" json:"ErrMsg,omitempty"`
	Token            *string        `protobuf:"bytes,16,opt,name=token" json:"token,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *Msg) Reset()                    { *m = Msg{} }
func (m *Msg) String() string            { return proto.CompactTextString(m) }
func (*Msg) ProtoMessage()               {}
func (*Msg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Msg) GetSign() string {
	if m != nil && m.Sign != nil {
		return *m.Sign
	}
	return ""
}

func (m *Msg) GetDataLen() string {
	if m != nil && m.DataLen != nil {
		return *m.DataLen
	}
	return ""
}

func (m *Msg) GetUserOpt() int32 {
	if m != nil && m.UserOpt != nil {
		return *m.UserOpt
	}
	return 0
}

func (m *Msg) GetOptResult() string {
	if m != nil && m.OptResult != nil {
		return *m.OptResult
	}
	return ""
}

func (m *Msg) GetReceiveResult() string {
	if m != nil && m.ReceiveResult != nil {
		return *m.ReceiveResult
	}
	return ""
}

func (m *Msg) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *Msg) GetFriends() []*User {
	if m != nil {
		return m.Friends
	}
	return nil
}

func (m *Msg) GetFriendLists() []*FriendList {
	if m != nil {
		return m.FriendLists
	}
	return nil
}

func (m *Msg) GetNewFriendRequest() []*User {
	if m != nil {
		return m.NewFriendRequest
	}
	return nil
}

func (m *Msg) GetGroups() []*Group {
	if m != nil {
		return m.Groups
	}
	return nil
}

func (m *Msg) GetGroupMsg() []*GroupMsg {
	if m != nil {
		return m.GroupMsg
	}
	return nil
}

func (m *Msg) GetOfflineMsg() []*PersonalMsg {
	if m != nil {
		return m.OfflineMsg
	}
	return nil
}

func (m *Msg) GetPersonalMsg() *PersonalMsg {
	if m != nil {
		return m.PersonalMsg
	}
	return nil
}

func (m *Msg) GetSrchInfo() *SearchInfo {
	if m != nil {
		return m.SrchInfo
	}
	return nil
}

func (m *Msg) GetErrMsg() string {
	if m != nil && m.ErrMsg != nil {
		return *m.ErrMsg
	}
	return ""
}

func (m *Msg) GetToken() string {
	if m != nil && m.Token != nil {
		return *m.Token
	}
	return ""
}

// 用户
type User struct {
	UserID           *int32      `protobuf:"varint,1,req,name=userID" json:"userID,omitempty"`
	UserPwd          *string     `protobuf:"bytes,2,opt,name=userPwd" json:"userPwd,omitempty"`
	NickName         *string     `protobuf:"bytes,3,opt,name=nickName" json:"nickName,omitempty"`
	Icon             []byte      `protobuf:"bytes,4,opt,name=icon" json:"icon,omitempty"`
	IconName         *string     `protobuf:"bytes,5,opt,name=iconName" json:"iconName,omitempty"`
	UserDetail       *UserDetail `protobuf:"bytes,6,opt,name=userDetail" json:"userDetail,omitempty"`
	IsOnline         *bool       `protobuf:"varint,7,opt,name=isOnline" json:"isOnline,omitempty"`
	UesrIntro        *string     `protobuf:"bytes,8,opt,name=uesrIntro" json:"uesrIntro,omitempty"`
	UserRoleID       *int32      `protobuf:"varint,9,opt,name=userRoleID" json:"userRoleID,omitempty"`
	Remark           *string     `protobuf:"bytes,10,opt,name=remark" json:"remark,omitempty"`
	ListNO           *int32      `protobuf:"varint,11,opt,name=listNO" json:"listNO,omitempty"`
	IsAgree          *bool       `protobuf:"varint,12,opt,name=isAgree" json:"isAgree,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *User) GetUserID() int32 {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return 0
}

func (m *User) GetUserPwd() string {
	if m != nil && m.UserPwd != nil {
		return *m.UserPwd
	}
	return ""
}

func (m *User) GetNickName() string {
	if m != nil && m.NickName != nil {
		return *m.NickName
	}
	return ""
}

func (m *User) GetIcon() []byte {
	if m != nil {
		return m.Icon
	}
	return nil
}

func (m *User) GetIconName() string {
	if m != nil && m.IconName != nil {
		return *m.IconName
	}
	return ""
}

func (m *User) GetUserDetail() *UserDetail {
	if m != nil {
		return m.UserDetail
	}
	return nil
}

func (m *User) GetIsOnline() bool {
	if m != nil && m.IsOnline != nil {
		return *m.IsOnline
	}
	return false
}

func (m *User) GetUesrIntro() string {
	if m != nil && m.UesrIntro != nil {
		return *m.UesrIntro
	}
	return ""
}

func (m *User) GetUserRoleID() int32 {
	if m != nil && m.UserRoleID != nil {
		return *m.UserRoleID
	}
	return 0
}

func (m *User) GetRemark() string {
	if m != nil && m.Remark != nil {
		return *m.Remark
	}
	return ""
}

func (m *User) GetListNO() int32 {
	if m != nil && m.ListNO != nil {
		return *m.ListNO
	}
	return 0
}

func (m *User) GetIsAgree() bool {
	if m != nil && m.IsAgree != nil {
		return *m.IsAgree
	}
	return false
}

// 用户详情信息
type UserDetail struct {
	UserID           *int32            `protobuf:"varint,1,opt,name=userID" json:"userID,omitempty"`
	Phone            *string           `protobuf:"bytes,2,opt,name=phone" json:"phone,omitempty"`
	Addr             *string           `protobuf:"bytes,3,opt,name=addr" json:"addr,omitempty"`
	QQ               *string           `protobuf:"bytes,4,opt,name=QQ,json=qQ" json:"QQ,omitempty"`
	Wechat           *string           `protobuf:"bytes,5,opt,name=wechat" json:"wechat,omitempty"`
	Sex              *string           `protobuf:"bytes,6,opt,name=sex" json:"sex,omitempty"`
	Age              *int32            `protobuf:"varint,7,opt,name=age" json:"age,omitempty"`
	IdNo             *string           `protobuf:"bytes,8,opt,name=idNo" json:"idNo,omitempty"`
	Email            *string           `protobuf:"bytes,9,opt,name=email" json:"email,omitempty"`
	CcNo             *string           `protobuf:"bytes,10,opt,name=ccNo" json:"ccNo,omitempty"`
	ScNo             *string           `protobuf:"bytes,11,opt,name=scNo" json:"scNo,omitempty"`
	StdNo            *string           `protobuf:"bytes,12,opt,name=stdNo" json:"stdNo,omitempty"`
	CusteomAttr      []*UserCustomAttr `protobuf:"bytes,13,rep,name=custeomAttr" json:"custeomAttr,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *UserDetail) Reset()                    { *m = UserDetail{} }
func (m *UserDetail) String() string            { return proto.CompactTextString(m) }
func (*UserDetail) ProtoMessage()               {}
func (*UserDetail) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *UserDetail) GetUserID() int32 {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return 0
}

func (m *UserDetail) GetPhone() string {
	if m != nil && m.Phone != nil {
		return *m.Phone
	}
	return ""
}

func (m *UserDetail) GetAddr() string {
	if m != nil && m.Addr != nil {
		return *m.Addr
	}
	return ""
}

func (m *UserDetail) GetQQ() string {
	if m != nil && m.QQ != nil {
		return *m.QQ
	}
	return ""
}

func (m *UserDetail) GetWechat() string {
	if m != nil && m.Wechat != nil {
		return *m.Wechat
	}
	return ""
}

func (m *UserDetail) GetSex() string {
	if m != nil && m.Sex != nil {
		return *m.Sex
	}
	return ""
}

func (m *UserDetail) GetAge() int32 {
	if m != nil && m.Age != nil {
		return *m.Age
	}
	return 0
}

func (m *UserDetail) GetIdNo() string {
	if m != nil && m.IdNo != nil {
		return *m.IdNo
	}
	return ""
}

func (m *UserDetail) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *UserDetail) GetCcNo() string {
	if m != nil && m.CcNo != nil {
		return *m.CcNo
	}
	return ""
}

func (m *UserDetail) GetScNo() string {
	if m != nil && m.ScNo != nil {
		return *m.ScNo
	}
	return ""
}

func (m *UserDetail) GetStdNo() string {
	if m != nil && m.StdNo != nil {
		return *m.StdNo
	}
	return ""
}

func (m *UserDetail) GetCusteomAttr() []*UserCustomAttr {
	if m != nil {
		return m.CusteomAttr
	}
	return nil
}

// 用户自定义属性
type UserCustomAttr struct {
	UserName         *string  `protobuf:"bytes,1,req,name=userName" json:"userName,omitempty"`
	AttrName         []string `protobuf:"bytes,2,rep,name=attrName" json:"attrName,omitempty"`
	AttrContent      []string `protobuf:"bytes,3,rep,name=attrContent" json:"attrContent,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *UserCustomAttr) Reset()                    { *m = UserCustomAttr{} }
func (m *UserCustomAttr) String() string            { return proto.CompactTextString(m) }
func (*UserCustomAttr) ProtoMessage()               {}
func (*UserCustomAttr) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *UserCustomAttr) GetUserName() string {
	if m != nil && m.UserName != nil {
		return *m.UserName
	}
	return ""
}

func (m *UserCustomAttr) GetAttrName() []string {
	if m != nil {
		return m.AttrName
	}
	return nil
}

func (m *UserCustomAttr) GetAttrContent() []string {
	if m != nil {
		return m.AttrContent
	}
	return nil
}

// 好友分组
type FriendList struct {
	ListNO           *int32  `protobuf:"varint,1,req,name=listNO" json:"listNO,omitempty"`
	ListName         *string `protobuf:"bytes,2,opt,name=listName" json:"listName,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *FriendList) Reset()                    { *m = FriendList{} }
func (m *FriendList) String() string            { return proto.CompactTextString(m) }
func (*FriendList) ProtoMessage()               {}
func (*FriendList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *FriendList) GetListNO() int32 {
	if m != nil && m.ListNO != nil {
		return *m.ListNO
	}
	return 0
}

func (m *FriendList) GetListName() string {
	if m != nil && m.ListName != nil {
		return *m.ListName
	}
	return ""
}

// 搜索请求包含的信息
type SearchInfo struct {
	SearchType       *int32  `protobuf:"varint,1,req,name=searchType" json:"searchType,omitempty"`
	SrchName         *string `protobuf:"bytes,2,opt,name=srchName" json:"srchName,omitempty"`
	SrchAttrb        *string `protobuf:"bytes,3,opt,name=srchAttrb" json:"srchAttrb,omitempty"`
	OnlyOnline       *bool   `protobuf:"varint,4,opt,name=onlyOnline" json:"onlyOnline,omitempty"`
	AgeLow           *int32  `protobuf:"varint,5,opt,name=ageLow" json:"ageLow,omitempty"`
	AgeHigh          *int32  `protobuf:"varint,6,opt,name=ageHigh" json:"ageHigh,omitempty"`
	SelectMale       *bool   `protobuf:"varint,7,opt,name=selectMale" json:"selectMale,omitempty"`
	SelectFemale     *bool   `protobuf:"varint,8,opt,name=selectFemale" json:"selectFemale,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *SearchInfo) Reset()                    { *m = SearchInfo{} }
func (m *SearchInfo) String() string            { return proto.CompactTextString(m) }
func (*SearchInfo) ProtoMessage()               {}
func (*SearchInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *SearchInfo) GetSearchType() int32 {
	if m != nil && m.SearchType != nil {
		return *m.SearchType
	}
	return 0
}

func (m *SearchInfo) GetSrchName() string {
	if m != nil && m.SrchName != nil {
		return *m.SrchName
	}
	return ""
}

func (m *SearchInfo) GetSrchAttrb() string {
	if m != nil && m.SrchAttrb != nil {
		return *m.SrchAttrb
	}
	return ""
}

func (m *SearchInfo) GetOnlyOnline() bool {
	if m != nil && m.OnlyOnline != nil {
		return *m.OnlyOnline
	}
	return false
}

func (m *SearchInfo) GetAgeLow() int32 {
	if m != nil && m.AgeLow != nil {
		return *m.AgeLow
	}
	return 0
}

func (m *SearchInfo) GetAgeHigh() int32 {
	if m != nil && m.AgeHigh != nil {
		return *m.AgeHigh
	}
	return 0
}

func (m *SearchInfo) GetSelectMale() bool {
	if m != nil && m.SelectMale != nil {
		return *m.SelectMale
	}
	return false
}

func (m *SearchInfo) GetSelectFemale() bool {
	if m != nil && m.SelectFemale != nil {
		return *m.SelectFemale
	}
	return false
}

// 群
type Group struct {
	GroupID          *int32         `protobuf:"varint,1,req,name=groupID" json:"groupID,omitempty"`
	GroupName        *string        `protobuf:"bytes,2,opt,name=groupName" json:"groupName,omitempty"`
	GroupIntro       *string        `protobuf:"bytes,3,opt,name=groupIntro" json:"groupIntro,omitempty"`
	CreateTime       *string        `protobuf:"bytes,4,opt,name=createTime" json:"createTime,omitempty"`
	GroupNumber      []*GroupNumber `protobuf:"bytes,5,rep,name=groupNumber" json:"groupNumber,omitempty"`
	Notices          []*GroupNotice `protobuf:"bytes,6,rep,name=notices" json:"notices,omitempty"`
	Rank             *int32         `protobuf:"varint,7,opt,name=rank" json:"rank,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *Group) Reset()                    { *m = Group{} }
func (m *Group) String() string            { return proto.CompactTextString(m) }
func (*Group) ProtoMessage()               {}
func (*Group) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Group) GetGroupID() int32 {
	if m != nil && m.GroupID != nil {
		return *m.GroupID
	}
	return 0
}

func (m *Group) GetGroupName() string {
	if m != nil && m.GroupName != nil {
		return *m.GroupName
	}
	return ""
}

func (m *Group) GetGroupIntro() string {
	if m != nil && m.GroupIntro != nil {
		return *m.GroupIntro
	}
	return ""
}

func (m *Group) GetCreateTime() string {
	if m != nil && m.CreateTime != nil {
		return *m.CreateTime
	}
	return ""
}

func (m *Group) GetGroupNumber() []*GroupNumber {
	if m != nil {
		return m.GroupNumber
	}
	return nil
}

func (m *Group) GetNotices() []*GroupNotice {
	if m != nil {
		return m.Notices
	}
	return nil
}

func (m *Group) GetRank() int32 {
	if m != nil && m.Rank != nil {
		return *m.Rank
	}
	return 0
}

// 群成员
type GroupNumber struct {
	GroupID          *int32  `protobuf:"varint,1,req,name=groupID" json:"groupID,omitempty"`
	NumberID         *int32  `protobuf:"varint,2,opt,name=numberID" json:"numberID,omitempty"`
	Remark           *string `protobuf:"bytes,3,opt,name=remark" json:"remark,omitempty"`
	Identity         *int32  `protobuf:"varint,4,opt,name=identity" json:"identity,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *GroupNumber) Reset()                    { *m = GroupNumber{} }
func (m *GroupNumber) String() string            { return proto.CompactTextString(m) }
func (*GroupNumber) ProtoMessage()               {}
func (*GroupNumber) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *GroupNumber) GetGroupID() int32 {
	if m != nil && m.GroupID != nil {
		return *m.GroupID
	}
	return 0
}

func (m *GroupNumber) GetNumberID() int32 {
	if m != nil && m.NumberID != nil {
		return *m.NumberID
	}
	return 0
}

func (m *GroupNumber) GetRemark() string {
	if m != nil && m.Remark != nil {
		return *m.Remark
	}
	return ""
}

func (m *GroupNumber) GetIdentity() int32 {
	if m != nil && m.Identity != nil {
		return *m.Identity
	}
	return 0
}

// 群消息
type GroupMsg struct {
	GroupID          *int32  `protobuf:"varint,1,req,name=groupID" json:"groupID,omitempty"`
	SendTime         *string `protobuf:"bytes,2,opt,name=sendTime" json:"sendTime,omitempty"`
	SenderID         *int32  `protobuf:"varint,3,opt,name=senderID" json:"senderID,omitempty"`
	Content          *string `protobuf:"bytes,4,opt,name=content" json:"content,omitempty"`
	ReadedTime       *int32  `protobuf:"varint,5,opt,name=readedTime" json:"readedTime,omitempty"`
	MsgType          *int32  `protobuf:"varint,6,opt,name=msgType" json:"msgType,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *GroupMsg) Reset()                    { *m = GroupMsg{} }
func (m *GroupMsg) String() string            { return proto.CompactTextString(m) }
func (*GroupMsg) ProtoMessage()               {}
func (*GroupMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *GroupMsg) GetGroupID() int32 {
	if m != nil && m.GroupID != nil {
		return *m.GroupID
	}
	return 0
}

func (m *GroupMsg) GetSendTime() string {
	if m != nil && m.SendTime != nil {
		return *m.SendTime
	}
	return ""
}

func (m *GroupMsg) GetSenderID() int32 {
	if m != nil && m.SenderID != nil {
		return *m.SenderID
	}
	return 0
}

func (m *GroupMsg) GetContent() string {
	if m != nil && m.Content != nil {
		return *m.Content
	}
	return ""
}

func (m *GroupMsg) GetReadedTime() int32 {
	if m != nil && m.ReadedTime != nil {
		return *m.ReadedTime
	}
	return 0
}

func (m *GroupMsg) GetMsgType() int32 {
	if m != nil && m.MsgType != nil {
		return *m.MsgType
	}
	return 0
}

// 个人消息
type PersonalMsg struct {
	SenderID         *int32  `protobuf:"varint,1,opt,name=senderID" json:"senderID,omitempty"`
	RecverID         []int32 `protobuf:"varint,2,rep,name=recverID" json:"recverID,omitempty"`
	SendTime         *string `protobuf:"bytes,3,opt,name=sendTime" json:"sendTime,omitempty"`
	ReadTime         *string `protobuf:"bytes,4,opt,name=readTime" json:"readTime,omitempty"`
	Content          *string `protobuf:"bytes,5,opt,name=content" json:"content,omitempty"`
	MsgType          *int32  `protobuf:"varint,6,opt,name=msgType" json:"msgType,omitempty"`
	IsReaded         *bool   `protobuf:"varint,7,opt,name=isReaded" json:"isReaded,omitempty"`
	MsgID            *string `protobuf:"bytes,8,opt,name=msgID" json:"msgID,omitempty"`
	FileName         *string `protobuf:"bytes,10,opt,name=fileName" json:"fileName,omitempty"`
	FileMd5          *string `protobuf:"bytes,11,opt,name=FileMd5,json=fileMd5" json:"FileMd5,omitempty"`
	FileLen          *string `protobuf:"bytes,13,opt,name=fileLen" json:"fileLen,omitempty"`
	IsNeedUpload     *bool   `protobuf:"varint,15,opt,name=isNeedUpload" json:"isNeedUpload,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PersonalMsg) Reset()                    { *m = PersonalMsg{} }
func (m *PersonalMsg) String() string            { return proto.CompactTextString(m) }
func (*PersonalMsg) ProtoMessage()               {}
func (*PersonalMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *PersonalMsg) GetSenderID() int32 {
	if m != nil && m.SenderID != nil {
		return *m.SenderID
	}
	return 0
}

func (m *PersonalMsg) GetRecverID() []int32 {
	if m != nil {
		return m.RecverID
	}
	return nil
}

func (m *PersonalMsg) GetSendTime() string {
	if m != nil && m.SendTime != nil {
		return *m.SendTime
	}
	return ""
}

func (m *PersonalMsg) GetReadTime() string {
	if m != nil && m.ReadTime != nil {
		return *m.ReadTime
	}
	return ""
}

func (m *PersonalMsg) GetContent() string {
	if m != nil && m.Content != nil {
		return *m.Content
	}
	return ""
}

func (m *PersonalMsg) GetMsgType() int32 {
	if m != nil && m.MsgType != nil {
		return *m.MsgType
	}
	return 0
}

func (m *PersonalMsg) GetIsReaded() bool {
	if m != nil && m.IsReaded != nil {
		return *m.IsReaded
	}
	return false
}

func (m *PersonalMsg) GetMsgID() string {
	if m != nil && m.MsgID != nil {
		return *m.MsgID
	}
	return ""
}

func (m *PersonalMsg) GetFileName() string {
	if m != nil && m.FileName != nil {
		return *m.FileName
	}
	return ""
}

func (m *PersonalMsg) GetFileMd5() string {
	if m != nil && m.FileMd5 != nil {
		return *m.FileMd5
	}
	return ""
}

func (m *PersonalMsg) GetFileLen() string {
	if m != nil && m.FileLen != nil {
		return *m.FileLen
	}
	return ""
}

func (m *PersonalMsg) GetIsNeedUpload() bool {
	if m != nil && m.IsNeedUpload != nil {
		return *m.IsNeedUpload
	}
	return false
}

// 群公告
type GroupNotice struct {
	GroupID          *int32  `protobuf:"varint,1,req,name=groupID" json:"groupID,omitempty"`
	CreateTime       *string `protobuf:"bytes,2,opt,name=createTime" json:"createTime,omitempty"`
	CreateID         *int32  `protobuf:"varint,3,opt,name=createID" json:"createID,omitempty"`
	Title            *string `protobuf:"bytes,4,opt,name=title" json:"title,omitempty"`
	Content          *string `protobuf:"bytes,5,opt,name=content" json:"content,omitempty"`
	ModifyTime       *string `protobuf:"bytes,6,opt,name=modifyTime" json:"modifyTime,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *GroupNotice) Reset()                    { *m = GroupNotice{} }
func (m *GroupNotice) String() string            { return proto.CompactTextString(m) }
func (*GroupNotice) ProtoMessage()               {}
func (*GroupNotice) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *GroupNotice) GetGroupID() int32 {
	if m != nil && m.GroupID != nil {
		return *m.GroupID
	}
	return 0
}

func (m *GroupNotice) GetCreateTime() string {
	if m != nil && m.CreateTime != nil {
		return *m.CreateTime
	}
	return ""
}

func (m *GroupNotice) GetCreateID() int32 {
	if m != nil && m.CreateID != nil {
		return *m.CreateID
	}
	return 0
}

func (m *GroupNotice) GetTitle() string {
	if m != nil && m.Title != nil {
		return *m.Title
	}
	return ""
}

func (m *GroupNotice) GetContent() string {
	if m != nil && m.Content != nil {
		return *m.Content
	}
	return ""
}

func (m *GroupNotice) GetModifyTime() string {
	if m != nil && m.ModifyTime != nil {
		return *m.ModifyTime
	}
	return ""
}

func init() {
	proto.RegisterType((*Msg)(nil), "Msg")
	proto.RegisterType((*User)(nil), "User")
	proto.RegisterType((*UserDetail)(nil), "UserDetail")
	proto.RegisterType((*UserCustomAttr)(nil), "UserCustomAttr")
	proto.RegisterType((*FriendList)(nil), "FriendList")
	proto.RegisterType((*SearchInfo)(nil), "SearchInfo")
	proto.RegisterType((*Group)(nil), "Group")
	proto.RegisterType((*GroupNumber)(nil), "GroupNumber")
	proto.RegisterType((*GroupMsg)(nil), "GroupMsg")
	proto.RegisterType((*PersonalMsg)(nil), "PersonalMsg")
	proto.RegisterType((*GroupNotice)(nil), "GroupNotice")
}

var fileDescriptor0 = []byte{
	// 1120 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x56, 0xdb, 0x6e, 0x24, 0x35,
	0x10, 0xd5, 0xdc, 0x92, 0x99, 0xea, 0x6c, 0x36, 0xb2, 0x10, 0x6a, 0x10, 0x82, 0xa8, 0xc5, 0x65,
	0x25, 0x60, 0x24, 0x90, 0x78, 0x67, 0xb5, 0x4b, 0x96, 0x48, 0xb9, 0x9a, 0xdd, 0x0f, 0xe8, 0x74,
	0x7b, 0x26, 0x4d, 0x7a, 0xda, 0xb3, 0x6d, 0xcf, 0x86, 0xf9, 0x26, 0x5e, 0x90, 0x10, 0x9f, 0x82,
	0xf8, 0x08, 0x7e, 0x81, 0x07, 0xaa, 0xca, 0x76, 0xb7, 0x87, 0x88, 0x3c, 0x4d, 0x9d, 0x53, 0xe5,
	0xb2, 0x5d, 0x75, 0x5c, 0x3d, 0x70, 0x78, 0xbe, 0xbd, 0x6a, 0xf5, 0xcd, 0x66, 0x31, 0x5f, 0xb7,
	0xda, 0xea, 0xec, 0xd7, 0x31, 0x8c, 0xce, 0xcd, 0x52, 0x08, 0x18, 0x9b, 0x6a, 0xd9, 0xa4, 0x83,
	0xe3, 0xe1, 0xb3, 0x99, 0x64, 0x5b, 0xa4, 0xb0, 0x5f, 0xe6, 0x36, 0x3f, 0x53, 0x4d, 0x3a, 0x64,
	0x3a, 0x40, 0xf2, 0x6c, 0x8c, 0x6a, 0x2f, 0xd7, 0x36, 0x1d, 0xa1, 0x67, 0x22, 0x03, 0x14, 0x1f,
	0xc1, 0x4c, 0xaf, 0xad, 0x54, 0x66, 0x53, 0xdb, 0x74, 0x7c, 0x3c, 0xc0, 0x55, 0x3d, 0x21, 0x3e,
	0x85, 0x27, 0xad, 0x2a, 0x54, 0xf5, 0x4e, 0xf9, 0x88, 0x09, 0x47, 0xec, 0x92, 0xe2, 0x03, 0x18,
	0x53, 0xba, 0x74, 0x0f, 0x9d, 0xc9, 0xb7, 0x93, 0xf9, 0x1b, 0x04, 0x92, 0x29, 0xf1, 0x09, 0xec,
	0x2f, 0xda, 0x4a, 0x35, 0xa5, 0x49, 0xf7, 0x8f, 0x47, 0xbd, 0x37, 0xb0, 0xe2, 0x6b, 0x48, 0x9c,
	0x79, 0x56, 0x19, 0x6b, 0xd2, 0x29, 0x07, 0x25, 0xf3, 0x93, 0x8e, 0x93, 0xb1, 0x5f, 0x7c, 0x03,
	0x47, 0x8d, 0xba, 0x77, 0x5e, 0xa9, 0xde, 0x6e, 0x94, 0xb1, 0xe9, 0x2c, 0x4e, 0xfc, 0xc0, 0x2d,
	0x3e, 0x86, 0xbd, 0x65, 0xab, 0x37, 0x6b, 0x93, 0x02, 0x07, 0xee, 0xcd, 0x5f, 0x11, 0x94, 0x9e,
	0x15, 0x9f, 0xc1, 0x94, 0x2d, 0xac, 0x6a, 0x9a, 0x70, 0xc4, 0xcc, 0x45, 0x20, 0x21, 0x3b, 0x97,
	0xf8, 0x0a, 0xe0, 0x72, 0xb1, 0xa8, 0xab, 0x46, 0x51, 0xe0, 0x01, 0x07, 0x1e, 0xcc, 0xaf, 0x54,
	0x6b, 0x74, 0x93, 0xd7, 0x14, 0x0b, 0xba, 0xf3, 0x8b, 0x39, 0x24, 0xeb, 0xde, 0x95, 0x3e, 0xe1,
	0xca, 0xec, 0x86, 0xc7, 0x01, 0xe2, 0x0b, 0x98, 0x9a, 0xb6, 0xb8, 0x3d, 0x6d, 0x16, 0x3a, 0x3d,
	0xe4, 0xe0, 0x64, 0xfe, 0x93, 0xca, 0x3d, 0x25, 0x3b, 0xa7, 0x78, 0x1f, 0xf6, 0x7e, 0x68, 0x5b,
	0xca, 0xf9, 0x94, 0x5b, 0xb1, 0xa7, 0x18, 0x89, 0xf7, 0x60, 0x62, 0xf5, 0x1d, 0x76, 0xfe, 0x88,
	0x69, 0x07, 0xb2, 0x3f, 0x87, 0x30, 0xa6, 0xb2, 0xd0, 0x32, 0xea, 0xc7, 0xe9, 0x4b, 0x16, 0xcc,
	0x44, 0x7a, 0x14, 0x84, 0x71, 0x75, 0x5f, 0xa2, 0x64, 0x68, 0x61, 0x80, 0xe2, 0x43, 0x98, 0x36,
	0x55, 0x71, 0x77, 0x91, 0xaf, 0x14, 0x6a, 0x86, 0x5c, 0x1d, 0x26, 0xf1, 0x55, 0x85, 0x6e, 0x58,
	0x2f, 0x07, 0x92, 0x6d, 0x8a, 0xa7, 0x5f, 0x8e, 0x77, 0x2a, 0xe9, 0xb0, 0xf8, 0x12, 0x80, 0xd2,
	0xbe, 0x54, 0x36, 0xaf, 0x6a, 0x2f, 0x93, 0x84, 0xfb, 0xe5, 0x28, 0x19, 0xb9, 0x39, 0x91, 0xb9,
	0x6c, 0xa8, 0x92, 0xa8, 0x99, 0xc1, 0xb3, 0xa9, 0xec, 0x30, 0xa9, 0x15, 0x7b, 0xda, 0x9e, 0x36,
	0xb6, 0xd5, 0xa8, 0x15, 0x56, 0x6b, 0x47, 0x60, 0xa7, 0x39, 0x8f, 0xd4, 0xb5, 0xc2, 0x8b, 0xce,
	0xd0, 0x3d, 0x91, 0x11, 0x43, 0x45, 0x68, 0xd5, 0x2a, 0x6f, 0xef, 0x50, 0x09, 0x5c, 0x3b, 0x87,
	0x88, 0xaf, 0x51, 0x5d, 0x17, 0x97, 0xd8, 0x7f, 0x5a, 0xe3, 0x11, 0x15, 0xa7, 0x32, 0xcf, 0x97,
	0xad, 0x52, 0xd8, 0x6f, 0x3a, 0x48, 0x80, 0xd9, 0x1f, 0x43, 0x80, 0xfe, 0xf8, 0x3b, 0xd5, 0x1d,
	0x44, 0xd5, 0xc5, 0xa6, 0xac, 0x6f, 0x35, 0xde, 0xc3, 0xd5, 0xd6, 0x01, 0xaa, 0x5e, 0x5e, 0x96,
	0xad, 0xaf, 0x2a, 0xdb, 0xe2, 0x10, 0x86, 0xd7, 0xd7, 0xfe, 0xfd, 0x0d, 0xdf, 0x5e, 0x53, 0xc6,
	0x7b, 0x55, 0xdc, 0xe6, 0xe1, 0xc5, 0x79, 0x24, 0x8e, 0x60, 0x64, 0xd4, 0x2f, 0x5c, 0xc2, 0x99,
	0x24, 0x93, 0x98, 0x7c, 0xe9, 0x2a, 0x35, 0x91, 0x64, 0x72, 0x77, 0xca, 0x8b, 0x50, 0x1f, 0xb6,
	0xe9, 0x24, 0x78, 0x57, 0x2c, 0xfe, 0xcc, 0x9d, 0x84, 0x01, 0x45, 0x16, 0x05, 0x46, 0xba, 0x72,
	0xb0, 0xcd, 0x83, 0x85, 0xb8, 0xc4, 0x71, 0x64, 0xd3, 0x6a, 0x63, 0x29, 0xe5, 0x81, 0x5b, 0xcd,
	0x00, 0xdf, 0x62, 0x52, 0x6c, 0x8c, 0x55, 0x7a, 0xf5, 0xdc, 0xda, 0x16, 0x35, 0x4e, 0x4f, 0xe2,
	0x29, 0xb7, 0xf5, 0x05, 0xf2, 0x8e, 0x96, 0x71, 0x4c, 0xf6, 0x33, 0x1c, 0xee, 0xba, 0xa9, 0xdb,
	0x54, 0x2c, 0x96, 0x8d, 0x9b, 0x65, 0x1d, 0x26, 0x5f, 0x8e, 0x31, 0xec, 0x1b, 0x62, 0x76, 0xf4,
	0x05, 0x2c, 0x8e, 0x21, 0x21, 0xfb, 0x85, 0x6e, 0xac, 0x6a, 0x68, 0xaa, 0x91, 0x3b, 0xa6, 0xb2,
	0xef, 0x01, 0xfa, 0x29, 0x12, 0xf5, 0xd8, 0x3f, 0x00, 0xdf, 0x63, 0xdc, 0x83, 0x2d, 0xb7, 0x07,
	0xcb, 0x36, 0xe0, 0xec, 0x9f, 0x01, 0x40, 0xff, 0x08, 0x49, 0x5e, 0x86, 0xd1, 0xeb, 0xed, 0x5a,
	0xf9, 0x34, 0x11, 0x43, 0xa9, 0xe8, 0x99, 0xc6, 0xa9, 0x02, 0x26, 0xe1, 0x92, 0x4d, 0x57, 0xbe,
	0xf1, 0x8d, 0xef, 0x09, 0xca, 0xac, 0x9b, 0x7a, 0xeb, 0x45, 0x3f, 0x66, 0xad, 0x45, 0x0c, 0x1d,
	0x1e, 0x1b, 0x7b, 0xa6, 0xef, 0x59, 0x0d, 0x78, 0x78, 0x87, 0x48, 0xa0, 0x68, 0xfd, 0x58, 0x2d,
	0x6f, 0x59, 0x11, 0x38, 0xd6, 0x3d, 0x74, 0x67, 0xad, 0x55, 0x61, 0xcf, 0xf3, 0x3a, 0x3c, 0xa3,
	0x88, 0x11, 0x19, 0x1c, 0x38, 0x74, 0x82, 0x4a, 0xc0, 0x88, 0x29, 0x47, 0xec, 0x70, 0xd9, 0xdf,
	0x03, 0x98, 0xf0, 0x20, 0xa4, 0x7d, 0x78, 0x0e, 0x76, 0xe3, 0x23, 0x40, 0xba, 0x17, 0x9b, 0xd1,
	0xa5, 0x7b, 0x82, 0x4e, 0xe1, 0x02, 0xf9, 0xbd, 0xba, 0x6b, 0x47, 0x0c, 0xf9, 0x8b, 0x56, 0xe5,
	0x56, 0xbd, 0xae, 0x56, 0xca, 0xab, 0x3f, 0x62, 0x68, 0x8a, 0xba, 0x64, 0x9b, 0xd5, 0x0d, 0x7e,
	0x5f, 0x26, 0x7e, 0xe8, 0xbe, 0xea, 0x39, 0x19, 0x07, 0x88, 0xcf, 0x61, 0xbf, 0xd1, 0xb6, 0x2a,
	0x94, 0xc1, 0x7a, 0xc4, 0xb1, 0x4c, 0xca, 0xe0, 0x24, 0x8d, 0xb7, 0x79, 0x73, 0xe7, 0x1f, 0x0d,
	0xdb, 0xd9, 0x3d, 0x24, 0x51, 0xde, 0x47, 0xae, 0x4c, 0x83, 0x91, 0x63, 0xd0, 0x35, 0xe4, 0x04,
	0x1d, 0x8e, 0x26, 0xcc, 0x68, 0x67, 0xc2, 0xd0, 0x4c, 0x2b, 0x51, 0x94, 0x95, 0xdd, 0xf2, 0x35,
	0x71, 0x4d, 0xc0, 0xd9, 0x6f, 0x03, 0x98, 0x86, 0xef, 0xcd, 0xe3, 0xdb, 0x1a, 0x14, 0x33, 0x57,
	0x2a, 0xa8, 0xcb, 0xe3, 0xe0, 0xe3, 0x23, 0x8d, 0x5c, 0xfa, 0x80, 0x29, 0x63, 0xe1, 0x1f, 0x89,
	0x2b, 0x70, 0x80, 0x54, 0x7d, 0x2c, 0x75, 0xa9, 0x5c, 0x4e, 0xa7, 0xac, 0x88, 0xa1, 0x95, 0x2b,
	0xb3, 0x64, 0xb1, 0x7b, 0x75, 0x79, 0x98, 0xfd, 0x35, 0x84, 0x24, 0xfa, 0x94, 0xed, 0xec, 0x3f,
	0xf8, 0xcf, 0xfe, 0xe8, 0xc3, 0x7f, 0x0b, 0xef, 0x7c, 0xb9, 0x46, 0xe4, 0x0b, 0x78, 0xe7, 0x4e,
	0xa3, 0x87, 0x77, 0xa2, 0xb3, 0x44, 0xca, 0xe8, 0x70, 0x7c, 0xa7, 0xc9, 0xee, 0x9d, 0xfe, 0xf7,
	0xcc, 0xee, 0xb3, 0x22, 0xf9, 0x76, 0xfd, 0x67, 0xc5, 0x61, 0x9a, 0x6f, 0x18, 0x86, 0x07, 0x74,
	0x23, 0xd3, 0x01, 0x5a, 0xb1, 0xa8, 0x6a, 0xc5, 0xd2, 0x76, 0x13, 0xb2, 0xc3, 0xb4, 0xcf, 0x09,
	0xda, 0xe7, 0xe5, 0x77, 0x7e, 0x50, 0xee, 0x2f, 0x1c, 0x24, 0x0f, 0x99, 0xf4, 0x27, 0xec, 0x49,
	0xef, 0xa1, 0x3f, 0x61, 0xf8, 0xe6, 0x2a, 0x73, 0xa1, 0x54, 0xf9, 0x66, 0x5d, 0xeb, 0xbc, 0xe4,
	0x0f, 0x38, 0xbe, 0xb9, 0x98, 0xcb, 0x7e, 0x1f, 0x04, 0x19, 0xb2, 0x54, 0x1f, 0xd1, 0xc3, 0xee,
	0xdb, 0x19, 0x3e, 0x78, 0x3b, 0x78, 0x7a, 0x87, 0x7a, 0x4d, 0x04, 0xcc, 0x7f, 0x16, 0x2a, 0x5b,
	0x87, 0xc2, 0x3a, 0xf0, 0x48, 0x55, 0x71, 0xaf, 0x95, 0x2e, 0xab, 0xc5, 0x96, 0xf7, 0x72, 0x1f,
	0x9f, 0x88, 0xf9, 0x37, 0x00, 0x00, 0xff, 0xff, 0x5d, 0xad, 0xa3, 0x60, 0xa5, 0x0a, 0x00, 0x00,
}
