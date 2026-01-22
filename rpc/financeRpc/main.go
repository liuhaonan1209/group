package main

import (
	"group/core"
	"log"
	"net"

	financeservice "group/kitex_gen/car/finance/financeservice"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

func main() {
	core.Nacos()
	core.Mysql()

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8894")
	// 1. 创建 Kitex 服务器的配置选项，指定监听地址和端口
	opts := []server.Option{
		server.WithServiceAddr(addr),
		// 可选：添加日志、限流等其他配置
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "finance-service"}),
	}
	// 2. 创建路由服务实例，传入实现了接口的结构体和配置选项
	svr := financeservice.NewServer(new(FinanceServiceImpl), opts...)
	// 3. 启动RPC服务器，开始监听客户端请求
	err := svr.Run()
	// 如果服务器启动失败，记录错误日志
	if err != nil {
		log.Println(err.Error())
	}
}
