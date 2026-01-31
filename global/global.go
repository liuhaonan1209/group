// Package global 全局变量包
// 定义项目中使用的全局变量，包括配置信息和数据库连接
package global

import (
	"group/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	// AppConfig 应用配置信息
	// 存储从Nacos配置中心读取的配置数据，包括MySQL、Redis等配置
	AppConfig config.Config
	
	// DB 数据库连接实例
	// GORM数据库连接对象，用于执行所有数据库操作
	DB *gorm.DB
	
	// RedisClient Redis客户端实例
	// Redis连接对象，用于缓存操作（如验证码存储）
	RedisClient *redis.Client
)
