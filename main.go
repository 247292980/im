package main

import (
	"fmt"
	"github.com/golang/glog"
	"net"
	"net/http"
)

func main() {

	go handleChannel() //处理通道

	listen, err := net.Listen("tcp", Conf.SrvAddress)
	if err != nil {
		glog.Error("net.Listen err" + err.Error())
		return
	}
	fmt.Println("监听中")

	go func() {
		http.HandleFunc("/downloadFile", DownLoadFile) // 下载文件
		http.HandleFunc("/file", HandleFile)           //上传文件
		err := http.ListenAndServe(Conf.FileSrvAddr, nil)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	}()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			glog.Info("新连接:" + conn.RemoteAddr().String())
			continue
		}
		fmt.Println("新连接:", conn.RemoteAddr().String())
		newConn := &myConn{
			conn:       conn,
			listenList: make([]iListener, 0),
		}
		go handleConn(newConn)
	}

}
