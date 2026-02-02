package main

import (
	"log"
	"net"

	upload "group/kitex_gen/car/upload/uploadservice"

	"github.com/cloudwego/kitex/server"
)

func main() {
	log.Println("Upload service starting...")
	
	// 设置服务监听地址
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8912")
	
	// 创建上传服务RPC服务器
	svr := upload.NewServer(new(UploadServiceImpl), server.WithServiceAddr(addr))
	
	log.Println("Upload service listening on :8912")
	
	// 启动服务器
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
