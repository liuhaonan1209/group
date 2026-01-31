// Package model 数据模型层
// 定义司机表的数据结构和数据库映射
package model

import (
	"time"

	"gorm.io/gorm"
)

// Driver 司机表模型
// 存储司机的基本信息，包括姓名、联系方式、身份证号、驾驶证号、评分等
// 使用GORM进行ORM映射
type Driver struct {
	gorm.Model                                                                       // 嵌入GORM基础模型（包含ID、CreatedAt、UpdatedAt、DeletedAt字段）
	Name         string    `gorm:"type:varchar(100);not null;comment:司机姓名" json:"name"` // 司机姓名
	Tel          string    `gorm:"type:varchar(20);not null;uniqueIndex;comment:司机联系电话" json:"tel"` // 司机手机号（唯一索引）
	IDCard       string    `gorm:"type:varchar(18);not null;uniqueIndex;comment:司机身份证" json:"id_card"` // 司机身份证号（唯一索引）
	License      string    `gorm:"type:varchar(50);not null;comment:司机驾驶证" json:"license"` // 司机驾驶证号
	RegisterDate time.Time `gorm:"type:datetime;not null;comment:司机注册时间" json:"register_date"` // 注册时间
	Rating       float64   `gorm:"type:decimal(3,2);default:5.00;comment:司机评分" json:"rating"` // 司机评分（默认5.00分）
}

// TableName 指定表名
// GORM会调用此方法获取表名，而不是使用默认的复数形式
// 返回:
//   - string: 数据库表名 "drivers"
func (Driver) TableName() string {
	return "drivers"
}
