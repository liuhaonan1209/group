package main

import (
	"group/core"
	order "group/kitex_gen/car/order/orderservice"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
)

func main() {
	// 初始化配置和数据库
	core.Nacos()
	core.Mysql()

	log.Println("Order service starting...")
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8899")
	// 创建订单服务RPC服务器
	svr := order.NewServer(new(OrderServiceImpl), server.WithServiceAddr(addr))
	log.Println("Passenger service listening on :8899")
	// 启动服务器
	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
	
}
