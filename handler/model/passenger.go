// Package model 数据模型层
// 定义乘客表的数据结构和数据库映射
package model

import (
	"time"

	"gorm.io/gorm"
)

// Passenger 乘客表模型
// 存储乘客的基本信息，包括姓名、联系方式、身份证号等
// 使用GORM进行ORM映射
type Passenger struct {
	gorm.Model                                                                       // 嵌入GORM基础模型（包含ID、CreatedAt、UpdatedAt、DeletedAt字段）
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`         // 乘客姓名
	Tel          string    `gorm:"type:varchar(20);not null;uniqueIndex;comment:乘客电话" json:"tel"` // 乘客手机号（唯一索引）
	IDCard       string    `gorm:"type:varchar(18);not null;uniqueIndex;comment:乘客身份证" json:"id_card"` // 乘客身份证号（唯一索引）
	RegisterDate time.Time `gorm:"type:datetime;not null;comment:乘客注册时间" json:"register_date"` // 注册时间
}

// TableName 指定表名
// GORM会调用此方法获取表名，而不是使用默认的复数形式
// 返回:
//   - string: 数据库表名 "passengers"
func (Passenger) TableName() string {
	return "passengers"
}
