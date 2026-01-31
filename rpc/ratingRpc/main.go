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

	log.Println("Rating service starting...")

	// 创建TCP监听器
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Rating service listening on :9999")

	// 创建服务实现
	service := &RatingServiceImpl{}

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

func handleConnection(conn net.Conn, service *RatingServiceImpl) {
	defer conn.Close()
	// 这里应该实现具体的RPC协议处理
	log.Printf("New connection from %s", conn.RemoteAddr())
}
