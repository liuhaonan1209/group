package main

import (
	"group/core"
	ordermanage "group/kitex_gen/car/order_manage/ordermanageservice"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
)

func main() {
	// 初始化配置和数据库
	core.Nacos()
	core.Mysql()

	log.Println("Order service starting...")
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8999")
	// 创建订单服务RPC服务器
	svr := ordermanage.NewServer(new(OrderManageServiceImpl), server.WithServiceAddr(addr))
	log.Println("Passenger service listening on :8999")
	// 启动服务器
	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
