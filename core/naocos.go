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

func Nacos() {
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": []constant.ServerConfig{{IpAddr: "14.103.167.35", Port: 8848}},
		"clientConfig":  []constant.ClientConfig{{NamespaceId: "pubilc"}},
	})
	if err != nil {
		panic("nacos 链接失败")
	}
	config, err := client.GetConfig(vo.ConfigParam{
		DataId: "document",
		Group:  "DEFAULT_GROUP",
	})
	if err != nil {
		panic("nacos 配置信息读取失败")
	}

	viper.SetConfigType("yaml")

	viper.ReadConfig(strings.NewReader(config))

	viper.Unmarshal(&global.AppConfig)

	log.Println("配置信息读取成功")
}
