// Package main 乘客服务RPC主程序
// 启动乘客服务的RPC服务器，监听端口8988
package main

import (
	"group/core"
	passenger "group/kitex_gen/car/passenger/passengerservice"
	"group/pkg/verifycode"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
)

// main 主函数
// 程序入口，负责初始化配置、数据库连接，并启动RPC服务器
// 执行流程：
//  1. 初始化Nacos配置中心连接，读取配置
//  2. 初始化MySQL数据库连接
//  3. 创建并启动Kitex RPC服务器
//  4. 监听8988端口，等待客户端请求
func main() {
	// 初始化Nacos配置中心
	core.Nacos()

	// 初始化MySQL数据库连接
	core.Mysql()

	// 初始化Redis连接
	core.Redis()

	// 启动验证码清理任务（用于内存存储的备用方案）
	verifycode.CleanExpiredCodes()
	log.Println("验证码清理任务已启动")

	log.Println("Passenger service starting...")

	// 解析TCP地址，监听所有网络接口的8988端口
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8988")

	// 创建Kitex RPC服务器
	// PassengerServiceImpl: 服务实现类
	// server.WithServiceAddr: 指定服务监听地址
	svr := passenger.NewServer(new(PassengerServiceImpl), server.WithServiceAddr(addr))

	log.Println("Passenger service listening on :8988")

	// 启动服务器，阻塞等待请求
	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
