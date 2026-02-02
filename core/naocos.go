// Package core 核心初始化包
// 提供Nacos配置中心的连接和配置读取功能
package core

import (
	"group/global"
	"log"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

// Nacos 初始化Nacos配置中心连接
// 从Nacos配置中心读取应用配置信息，并解析到全局配置变量中
// 功能：
//   1. 创建Nacos配置客户端
//   2. 从Nacos获取配置内容
//   3. 使用Viper解析YAML格式的配置
//   4. 将配置映射到全局配置结构体
// 如果连接失败或配置读取失败，程序会panic终止
func Nacos() {
	// 创建Nacos配置客户端
	// serverConfigs: Nacos服务器地址和端口
	// clientConfig: 客户端配置，包括命名空间ID
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": []constant.ServerConfig{{IpAddr: "14.103.167.35", Port: 8848}}, // Nacos服务器地址
		"clientConfig":  []constant.ClientConfig{{NamespaceId: "pubilc"}},                // 命名空间ID（注意：这里拼写为pubilc）
	})
	if err != nil {
		panic("nacos 链接失败")
	}
	
	// 从Nacos获取配置内容
	// DataId: 配置文件的ID
	// Group: 配置分组，默认为DEFAULT_GROUP
	config, err := client.GetConfig(vo.ConfigParam{
		DataId: "document",
		Group:  "DEFAULT_GROUP",
	})
	if err != nil {
		panic("nacos 配置信息读取失败")
	}

	// 设置Viper配置类型为YAML
	viper.SetConfigType("yaml")

	// 从字符串读取配置内容
	viper.ReadConfig(strings.NewReader(config))

	// 将配置解析并映射到全局配置结构体
	viper.Unmarshal(&global.AppConfig)

	log.Println("配置信息读取成功")
}
