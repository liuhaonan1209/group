package main

import (
	"group/core"
	refund "group/kitex_gen/car/refund/refundservice"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
)

func main() {
	// 初始化配置和数据库
	core.Nacos()
	core.Mysql()
	log.Println("Refund service starting...")
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8990")
	// 创建退票服务RPC服务器
	svr := refund.NewServer(new(RefundServiceImpl), server.WithServiceAddr(addr))
	log.Println("Refund service listening on :8990")
	// 启动服务器
	err := svr.Run()

	if err != nil {
		log.Println(err.Error())

	}
}
