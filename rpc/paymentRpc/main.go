package main

import (
	"group/core"
	"log"
	"net"
)

func main() {
	// 初始化配置和数据库
	core.Nacos()
	core.Mysql()

	log.Println("Payment service starting...")

	// 创建TCP监听器
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:7997")
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Payment service listening on :7997")

	// 创建服务实现
	service := &PaymentServiceImpl{}

	// 简单的服务循环（实际应该使用Kitex框架）
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go handleConnection(conn, service)
	}
}

func handleConnection(conn net.Conn, service *PaymentServiceImpl) {
	defer conn.Close()
	// 这里应该实现具体的RPC协议处理
	// 由于没有生成完整的Kitex代码，这里只是占位
	log.Printf("New connection from %s", conn.RemoteAddr())
}
