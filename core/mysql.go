package core

import (
	"fmt"
	"group/global"
	"group/handler/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Mysql() {

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	data := global.AppConfig.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		data.User,
		data.Password,
		data.Host,
		data.Port,
		data.Database)
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败")
	} else {
		fmt.Println("数据库连接成功")
	}
	err = global.DB.AutoMigrate(
		&model.Driver{},
		&model.Passenger{})
	if err != nil {
		panic("数据库迁移失败")
	} else {
		fmt.Println("数据库迁移成功")
	}
}
