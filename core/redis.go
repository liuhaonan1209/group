// Package core Redis初始化包
// 提供Redis连接和客户端管理功能
package core

import (
	"context"
	"fmt"
	"group/global"
	"log"

	"github.com/redis/go-redis/v9"
)

// Redis 初始化Redis连接
// 从全局配置中读取Redis配置信息，创建Redis客户端
// 功能：
//  1. 读取Redis配置（地址、端口、密码）
//  2. 创建Redis客户端
//  3. 测试连接是否成功
//  4. 将客户端保存到全局变量
func Redis() {
	// 构建Redis地址
	addr := fmt.Sprintf("%s:%d", global.AppConfig.Redis.Host, global.AppConfig.Redis.Port)

	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,                            // Redis地址
		Password: global.AppConfig.Redis.Password, // Redis密码（如果没有密码则为空字符串）
		DB:       0,                               // 使用默认DB 0
	})

	// 测试连接
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis连接失败: %v，验证码功能将使用内存存储", err)
		global.RedisClient = nil
		return
	}

	// 保存到全局变量
	global.RedisClient = rdb
	log.Println("Redis连接成功")
}
