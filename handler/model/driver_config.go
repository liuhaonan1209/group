// Package model 数据模型层 - 司机配置管理模块
// 本文件定义司机配置相关的数据库表结构
package model

import (
	"gorm.io/gorm"
)

// DriverConfig 司机配置表
// 用于存储司机的各项配置信息，包括功能开关、有效期设置、计费规则等
type DriverConfig struct {
	gorm.Model              // GORM基础模型（包含ID、CreatedAt、UpdatedAt、DeletedAt）
	DriverID        int64   `gorm:"not null;uniqueIndex;comment:司机ID" json:"driver_id"`                  // 司机ID（唯一索引，一个司机只有一条配置）
	CanPublishTrip  bool    `gorm:"default:true;comment:是否可以发布行程（功能开关）" json:"can_publish_trip"`         // 功能开关：支持关闭司机发布行程功能
	AutoNotify      bool    `gorm:"default:true;comment:是否开启自动通知（用户购票时触发）" json:"auto_notify"`           // 功能开关：开启后用户购票时自动短信通知司机
	MaxPassengers   int     `gorm:"default:-1;comment:有效期设置：最多可配置乘客数（-1表示无限制）" json:"max_passengers"`    // 有效期设置：司机发布行程有效天数可配置（-1表示无限制）
	StartPrice      float64 `gorm:"type:decimal(10,2);default:0;comment:计费规则：车辆起步价" json:"start_price"`  // 计费规则：可设置车辆起步价（按公里计费）
	PricePerKm      float64 `gorm:"type:decimal(10,2);default:0;comment:计费规则：每公里单价" json:"price_per_km"` // 计费规则：订单折扣（支持按乘车人数打折）
	VehicleNumber   string  `gorm:"type:varchar(20);comment:车辆合规性：车牌号" json:"vehicle_number"`            // 车辆合规性：车牌号
	VehicleType     string  `gorm:"type:varchar(50);comment:车辆合规性：车型" json:"vehicle_type"`               // 车辆合规性：车型（如：丰田凯美瑞）
	InsuranceExpiry string  `gorm:"type:varchar(20);comment:车辆合规性：保险到期日期" json:"insurance_expiry"`       // 车辆合规性：保险到期日期（格式：2024-12-31）
	IsCompliant     bool    `gorm:"default:false;comment:车辆是否合规" json:"is_compliant"`                    // 车辆是否合规（需要验证车牌号、保险等）
}

// TableName 指定数据库表名
// GORM会使用此方法返回的名称作为表名
func (DriverConfig) TableName() string {
	return "driver_configs"
}
