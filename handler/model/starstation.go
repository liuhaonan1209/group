package model

import (
	"time"

	"gorm.io/gorm"
)

// 站点
type StartStation struct {
	gorm.Model
	Name          string    `gorm:"type:varchar(255);unique_index;not null" json:"name"` // 站点名称（唯一）
	Latitude      float64   `gorm:"type:float;comment:纬度" json:"latitude"`             //纬度
	Longitude     float64   `gorm:"type:float;comment:经度" json:"longitude"`            //经度
	TimeFromStart time.Time `gorm:"type:datetime;default:null" json:"time_from_start"`
	Address       string    `gorm:"type:varchar(255);not null;comment:详细地址" json:"address"`                   // 详细地址
	IsActive      bool      `gorm:"type:tinyint(1);not null;default:true;index;comment:是否启用" json:"isActive"` // 是否启用：true-启用，false-禁用（默认启用，加普通索引）
}
