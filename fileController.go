package main

import (
	"fmt"
	_ "github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"gopkg.in/mgo.v2/bson"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

/*
Http定义了与服务器交互的不同方法，最基本的方法有4种，分别是GET，POST，PUT，DELETE。
HTTP中的GET，POST，PUT，DELETE就对应着对这个资源的查，改，增，删4个操作。
GET一般用于获取/查询资源信息，而POST一般用于更新资源信息。
GET 请求一般不应产生副作用，它仅仅是获取资源信息，不会修改，增加数据，不会影响资源的状态。
*/

//**************************************上传文件**************************************
func HandleFile(w http.ResponseWriter, r *http.Request) {
	sign := r.Header.Get("sign")
	if sign != "142857" {
		fmt.Fprintf(w, "what's sign?")
		return
	}
	if r.Method != "POST" {
		fmt.Fprintf(w, "which method do you upload")
	}
	r.ParseForm()
	r.ParseMultipartForm(32 << 20)
	var fileMD5, fileName, fileLen, Src, msgType, token string
	var Dst []string
	if r.MultipartForm != nil {
		fileMD5 = r.MultipartForm.Value["fileMd5"][0]
		fileName = r.MultipartForm.Value["fileName"][0]
		fileLen = r.MultipartForm.Value["fileLen"][0]
		Src = r.MultipartForm.Value["Src"][0]
		Dst = r.MultipartForm.Value["Dst"]
		msgType = r.MultipartForm.Value["msgType"][0]
		token = r.MultipartForm.Value["token"][0]
	}
	MsgTypeInt, _ := strconv.Atoi(msgType)
	SrcInt, _ := strconv.Atoi(Src)
	if SrcInt <= 0 {
		fmt.Fprintf(w, "who are you ?")
		return
	}
	if MsgTypeInt <= 1 {
		fmt.Fprintf(w, "你要传什么类型的数据")
		return
	}
	isOk := checkToken(SrcInt, token)
	if isOk == false {
		fmt.Fprintf(w, "用户已失效，请重新登录")
		return
	}
	if len(fileMD5) <= 0 {
		fmt.Fprintf(w, "what's file MD5 ?")
		return
	}
	fmt.Println(fileMD5, fileName, fileLen, Src, Dst, msgType)
	chatMsg := ChatMsg{} //查看是否已存在这个文件，若存在，则发送，不需要再次上传
	err := ChatCollect.Find(bson.M{"FileMD5": fileMD5}).One(&chatMsg)
	if err == nil {
		fmt.Fprintf(w, "file already exist in server ")
		return
	}

	// 到这里说明验证通过了
	file, handle, err := r.FormFile("file")
	if err != nil {
		w.Write([]byte("lack of 'file' key"))
		fmt.Printf("%s\n", err.Error())
		return
	}
	defer file.Close()
	os.Mkdir(FILE_PATH, os.ModePerm)
	fmt.Printf("接收到的文件名是：%s\n", handle.Filename)
	f, err := os.OpenFile(FILE_PATH+"/"+fileMD5, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintf(w, "file upload failed ")
		fmt.Printf("%s\n", err.Error())
		return
	}
	//上传成功
	io.Copy(f, file)
	f.Close()
	SendTime := time.Now()
	SendTimeStr := SendTime.Format(timeFormat)
	for k, v := range Dst {
		fmt.Println("k:", k, "v:", v)
		DstInt, _ := strconv.Atoi(v)
		err := ChatCollect.Insert(&ChatMsg{Src: SrcInt, Dst: DstInt, SendTime: SendTime, MsgType: MsgTypeInt, FileName: fileName, FileMD5: fileMD5, FileLen: fileLen})
		if err != nil {
			fmt.Fprintf(w, "send failed")
			fmt.Printf("%s\n", err.Error())
			return
		}
		FriendConn := GetOnlineUserByID(DstInt)
		if FriendConn == nil { //对方不在线，直接跳过下面的步骤
			continue
		}
		//运行到这里说明目标在线  直接发送
		newSendMsg := &Msg{
			UserOpt: proto.Int(Conf.ChatOP),
			PersonalMsg: &PersonalMsg{
				SenderID: proto.Int(SrcInt),
				RecverID: []int32{int32(DstInt)},
				SendTime: proto.String(SendTimeStr),
				FileName: proto.String(fileName),
				MsgType:  proto.Int(MsgTypeInt),
				FileMd5:  proto.String(fileMD5),
				FileLen:  proto.String(fileLen),
			},
		}
		FriendConn.write(newSendMsg)
	}
	fmt.Fprintf(w, "ok")
	fmt.Println("send sucess")
}

//**************************************下载文件**************************************
func DownLoadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DownLoadFile ")
	sign := r.Header.Get("sign")
	if sign != "142857" {
		fmt.Fprintf(w, "what's sign?")
		return
	}
	if r.Method != "GET" {
		fmt.Fprintf(w, "which method do you request")
		return
	}
	r.ParseForm()
	SrcID, err := strconv.Atoi(r.FormValue("userID"))
	if err != nil {
		fmt.Fprintf(w, "你是谁")
		return
	}
	token := r.FormValue("token")
	isOk := checkToken(SrcID, token)
	if isOk == false {
		fmt.Fprintf(w, "用户已失效，请重新登录")
		return
	} //到这里判断结束  开始从服务器取文件
	fileMD5 := r.FormValue("fileMD5")
	file, err := ioutil.ReadFile(FILE_PATH + "/" + fileMD5)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "file not exist in server")
		return
	}
	w.Write(file)
	fmt.Println("DownLoadFile ")
}
