package main

import (
	"group/core"
	"group/kitex_gen/car/starstation/starservice"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

func main() {
	core.Nacos()
	core.Mysql()
	addr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:8881")

	opts := []server.Option{
		server.WithServiceAddr(addr),
		// 可选：添加日志、限流等其他配置
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "starstation-service"}),
	}

	svr := starservice.NewServer(new(StarStationServiceImpl), opts...)

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
