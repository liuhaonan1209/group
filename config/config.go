// Package config 配置包
// 定义应用程序的配置结构，包括数据库、缓存等配置项
package config

// Config 应用程序总配置结构
// 包含所有子配置项，如MySQL、Redis等
type Config struct {
	Mysql Mysql // MySQL数据库配置
	Redis Redis // Redis缓存配置
}

// Mysql MySQL数据库配置结构
// 定义连接MySQL数据库所需的参数
type Mysql struct {
	User     string // 数据库用户名
	Password string // 数据库密码
	Host     string // 数据库主机地址
	Port     int    // 数据库端口号
	Database string // 数据库名称
}

// Redis Redis缓存配置结构
// 定义连接Redis服务所需的参数
type Redis struct {
	Password string // Redis密码
	Host     string // Redis主机地址
	Port     int    // Redis端口号
}
