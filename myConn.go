package main

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"regexp" //正则表达式
	"strconv"
	_ "time"
)

//***************************对客户端的连接进行处理*************************
func handleConn(newConn *myConn) {

	basicListener(newConn)
	for newConn.conn != nil {
		msg, err := newConn.read()
		if err != nil {
			fmt.Println("Read error: ", err)
			glog.V(0).Infoln(err)
			clnOffLineChannel <- newConn
			break
			//读取错误，说明客户端连接已断开
		}

		for _, l := range newConn.listenList {
			l.onProcess(msg)
		}
	}
}

func (this *myConn) read() (*Msg, error) {
	requestProBuf := &Msg{}           //请求buf
	headBuf := make([]byte, headSize) //包头  39    签名+头
	subDataLen := 0

	//读取数据
	dataLen, err := this.conn.Read(headBuf)
	if err != nil {
		clnOffLineChannel <- this
		return nil, err
	}
	// err = this.conn.SetReadDeadline(time.Now().Add(time.Second * 1)) //设置读取限制时间
	// if err != nil {
	// 	return nil, err
	// }
	for dataLen != headSize {
		subDataLen, err = this.conn.Read(headBuf[dataLen:])
		dataLen += subDataLen
		if err != nil {
			return nil, err
		}
	}
	receivedData := make([]byte, 1024)
	copy(receivedData, headBuf[:dataLen])
	correctPack, err := this.handlePackage(receivedData)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(correctPack, requestProBuf)
	if err != nil {
		// "解析数据时出错,你传过来的是什么鬼东西"
		return nil, err
	}
	return requestProBuf, nil
}

//***************************处理黏包碎包*************************************************************************************
func (this *myConn) handlePackage(receivedData []byte) ([]byte, error) {

	subDataLen := 0                                                 //再次读取的长度
	curDataLen := headSize                                          //收到的实际长度
	if string(receivedData[signBeginPos:signEndnPos]) == "142857" { //判断是否是 包头：签名
		dataLength, err := strconv.ParseInt((string(receivedData[dataLengthBeginPos:dataLengthEndnPos])), 0, 64) //取出包长
		if err != nil {
			clnOffLineChannel <- this
			return nil, err
		}

		if dataLength > int64(Conf.PackMaxSize) {
			fmt.Fprint(this.conn, [57]byte{10, 6, 49, 52, 50, 56, 53, 55, 18, 10, 48, 88, 48, 48, 48, 48, 48, 48, 51, 57, 24, 0, 34, 1, 48, 98, 30, 232, 175, 165, 229, 140, 133, 229, 164, 170, 229, 164, 167, 228, 186, 134, 239, 188, 140, 232, 175, 183, 229, 136, 134, 229, 140, 133, 229, 143, 145})
			clnOffLineChannel <- this
			return nil, errors.New("该包太大了，请分包发")

		}
		correctPack := make([]byte, dataLength) //复制到 最终包correctPack
		copy(correctPack, receivedData)
		// err = this.conn.SetReadDeadline(time.Now().Add(time.Second * 1)) //设置读取限制时间
		// if err != nil {
		// 	//读取超时，断开连接   不管是攻击者还是网络问题什么鬼的
		// 	clnOffLineChannel <- this
		// 	fmt.Println("Read error: ", err)
		// 	glog.V(0).Infoln(err)
		// }
		for int64(curDataLen) != dataLength { //读指定长度DataLen
			subDataLen, err = this.conn.Read(correctPack[curDataLen:])
			curDataLen += subDataLen
			if err != nil {
				clnOffLineChannel <- this
				return nil, err
			}
		}
		return correctPack, nil //执行if表示正常包、黏包  ， 执行else表示碎包、丢包
	} else {
		return nil, errors.New("丢包了，请重发...........")

	}
	return nil, nil
}

//***************************处理通道*******************************************************************************************
func handleChannel() {
	clnOffLineChannel = make(chan *myConn) //离线通道
	broadcastChannel = make(chan string)   //广播通道

	for {
		go getMsg() //广播通道          从键盘获取广播信息

		select {

		case myConn := <-clnOffLineChannel:
			if myConn.id != 0 { //运行到这里是登录过得   现在下线了
				if myConn.conn == userList[myConn.id].conn {
					delete(userList, myConn.id)
					showOnLines(userList)
				}
				break
			} //运行到这里 说明还未登陆就断开连接的
			myConn.lock.Lock()
			myConn.conn.Close()
			// myConn.conn=nil
			myConn.lock.Unlock()
		case msg := <-broadcastChannel:
			bMsg := []byte(msg)
			for _, v := range userList {
				_, err := v.conn.Write(bMsg)
				if err != nil {
					fmt.Println(err)
					clnOffLineChannel <- &v
				}
			}

		} //select结束

	} //for 结束
}

//*****************************打包数据*********************************************************************************************
func Pack(data Msg) ([]byte, error) {

	sign := Conf.Sign
	data.Sign = &sign

	dataLen := "0x12345678"
	data.DataLen = &dataLen
	dataBuf, err := proto.Marshal(&data)
	if err != nil {
		fmt.Println(err.Error(), "压缩proto.Marshal出错")
		return nil, err
	}

	if len(dataBuf) <= 2 {
		fmt.Println("Msg is empty!")
		return nil, errors.New("Msg is empty!")
	}
	dataLen = fmt.Sprintf("%#08X", len(dataBuf))
	data.DataLen = &dataLen
	dataBuf, err = proto.Marshal(&data)
	if err != nil {
		fmt.Println(err.Error(), "压缩proto.Marshal出错")
		return nil, err
	}
	return dataBuf, err
}

//***************************send数据到客户端******************************************************************************************
func (this *myConn) write(arg_data *Msg) {
	dataBuf, err := Pack(*arg_data)
	if err != nil {
		fmt.Println(err.Error(), "压缩到proto.Marshal包出错")
		glog.Error("压缩到proto.Marshal包出错:" + err.Error())
		return
	}
	_, err = this.conn.Write(dataBuf)
	if err != nil {
		fmt.Println(err.Error(), "发送消息失败，该用户不在线")
		glog.V(0).Info(err.Error(), "发送消息失败，该用户不在线")
		return
	}

}

//******************************当前在线人数******************************************************************************************
func showOnLines(arg_conns map[int]myConn) {
	fmt.Println("在线人数: " + strconv.Itoa(len(arg_conns)))
}

//检查字符的合法性
func checkWord(args string) (int, bool) {
	if args == "" {
		return 0, false
	}

	lenth := len(args)
	if lenth <= 0 {
		return 0, false
	}

	boolean, err := regexp.MatchString("[A-Za-z0-9_]", args)
	if err != nil || boolean == false {
		return lenth, false
	}
	return lenth, true

	// 允许;  英文 数字   下划线
}
