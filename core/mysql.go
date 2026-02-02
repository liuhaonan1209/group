// Package core 核心初始化包
// 提供MySQL数据库的初始化和连接功能
package core

import (
	"fmt"
	"group/global"
	"group/handler/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Mysql 初始化MySQL数据库连接
// 从全局配置中读取MySQL配置信息，建立数据库连接，并执行数据库迁移
// 功能：
//  1. 构建MySQL DSN连接字符串
//  2. 使用GORM连接MySQL数据库
//  3. 自动迁移数据表结构（Driver、Passenger表）
//
// 如果连接失败或迁移失败，程序会panic终止
func Mysql() {

	// 从全局配置中获取MySQL配置信息
	data := global.AppConfig.Mysql

	// 构建DSN（Data Source Name）连接字符串
	// 格式：用户名:密码@tcp(主机:端口)/数据库名?参数
	// charset=utf8mb4: 使用UTF-8字符集，支持emoji等特殊字符
	// parseTime=True: 自动解析数据库中的时间类型为Go的time.Time
	// loc=Local: 使用本地时区
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		data.User,
		data.Password,
		data.Host,
		data.Port,
		data.Database)

	var err error
	// 使用GORM连接MySQL数据库
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败")
	} else {
		fmt.Println("数据库连接成功")
	}

	// 自动迁移数据表结构
	// AutoMigrate会根据模型定义自动创建或更新表结构
	// 注意：只会添加新字段和索引，不会删除已存在的字段
	err = global.DB.AutoMigrate(
		&model.Driver{},             // 司机表
		&model.Passenger{},          // 乘客表
		&model.Trip{},               // 行程表
		&model.TripShare{},          // 行程分享记录表
		&model.PassengerHelp{},      // 乘客求助记录表
		&model.OrderManage{},        // 订单表
		&model.OrderStatusHistory{}, // 订单状态历史表
		&model.RefundManage{},       // 退票管理表
		&model.AuditLog{},           // 审计日志表
		&model.ExportLog{},          // 导出日志表
		&model.SystemLog{},          // 系统日志表
		&model.AccessLog{},          // 访问日志表
	)
	if err != nil {
		panic("数据库迁移失败")
	} else {
		fmt.Println("数据库迁移成功")
	}
}
