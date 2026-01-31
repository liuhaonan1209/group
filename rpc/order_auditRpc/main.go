package main

import (
	"group/core"
	"group/global"
	"group/handler/model"
	orderaudit "group/kitex_gen/car/order_audit/orderauditservice"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
)

func main() {
	// 初始化配置和数据库
	core.Nacos()
	core.Mysql()

	// 自动迁移审计相关表
	if err := global.DB.AutoMigrate(
		&model.OrderAuditLog{},
		&model.OrderAnomaly{},
		&model.AuditLogExport{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("审计表迁移成功")

	log.Println("Order Audit service starting...")
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8893")
	// 创建订单审计服务RPC服务器
	svr := orderaudit.NewServer(new(OrderAuditServiceImpl), server.WithServiceAddr(addr))
	log.Println("Order Audit service listening on :8893")
	// 启动服务器
	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
