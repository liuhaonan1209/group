package main

import (
	"group/core"
	routeservice "group/kitex_gen/car/route/routeservice" // 导入生成的路由服务包
	"log"

	"net"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

// 假设你的 RouteServiceImpl 实现了 RouteService 接口
// 这里先定义一个空结构体（实际需在 handler.go 中实现接口方法）

func main() {
	core.Nacos()
	core.Mysql()

	addr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:8892")
	// 1. 创建 Kitex 服务器的配置选项，指定监听地址和端口
	opts := []server.Option{
		server.WithServiceAddr(addr),
		// 可选：添加日志、限流等其他配置
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "route-service"}),
	}

	// 2. 创建路由服务实例，传入实现了接口的结构体和配置选项
	svr := routeservice.NewServer(new(RouteServiceImpl), opts...)

	// 3. 启动RPC服务器，开始监听客户端请求
	log.Println("Kitex RPC server starting on 0.0.0.0:8082...")
	err := svr.Run()

	// 如果服务器启动失败，记录错误日志
	if err != nil {
		log.Fatalf("RPC server failed to start: %v", err)
	}
}
