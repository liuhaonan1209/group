package main

import (
	"group/core"
	driver "group/kitex_gen/car/driver/driverservice"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
)

func main() {
	core.Nacos()
	core.Mysql()
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9988")
	svr := driver.NewServer(new(DriverServiceImpl), server.WithServiceAddr(addr))
	log.Println("Passenger service listening on :9988")
	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
