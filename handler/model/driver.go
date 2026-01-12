package model

import (
	"time"

	"gorm.io/gorm"
)

// Driver 司机表模型
type Driver struct {
	gorm.Model
	Name         string    `gorm:"type:varchar(100);not null;comment:司机姓名" json:"name"`
	Tel          string    `gorm:"type:varchar(20);not null;uniqueIndex;comment:司机联系电话" json:"tel"`
	IDCard       string    `gorm:"type:varchar(18);not null;uniqueIndex;comment:司机身份证" json:"id_card"`
	License      string    `gorm:"type:varchar(50);not null;comment:司机驾驶证" json:"license"`
	RegisterDate time.Time `gorm:"type:datetime;not null;comment:司机注册时间" json:"register_date"`
	Rating       float64   `gorm:"type:decimal(3,2);default:5.00;comment:司机评分" json:"rating"`
}

func (Driver) TableName() string {
	return "drivers"
}
