package model

import (
	"time"

	"gorm.io/gorm"
)

// Passenger 乘客表模型
type Passenger struct {
	gorm.Model
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	Tel          string    `gorm:"type:varchar(20);not null;uniqueIndex;comment:乘客电话" json:"tel"`
	IDCard       string    `gorm:"type:varchar(18);not null;uniqueIndex;comment:乘客身份证" json:"id_card"`
	RegisterDate time.Time `gorm:"type:datetime;not null;comment:乘客注册时间" json:"register_date"`
}

func (Passenger) TableName() string {
	return "passengers"
}
