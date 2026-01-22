package model

import (
	"time"

	"gorm.io/gorm"
)

// 司机结算表
type DriverSettlement struct {
	gorm.Model
	SettlementCycle     string    `gorm:"type:varchar(10);not null;comment:'结算周期（可选值：日/周/月，工单要求自动生成每日/每周/每月账单）'"`
	SettlementStartDate time.Time `gorm:"type:date;not null;comment:'结算周期开始日期'"`
	SettlementEndDate   time.Time `gorm:"type:date;not null;comment:'结算周期结束日期'"`
	DriverID            string    `gorm:"type:varchar(32);not null;comment:'司机ID（关联司机表，工单要求司机可核对自己的账单）'"`
	DriverName          string    `gorm:"type:varchar(50);not null;comment:'司机姓名'"`
	TotalOrderAmount    float64   `gorm:"type:decimal(10,2);not null;comment:'周期内总订单金额'"`
	PlatformFee         float64   `gorm:"type:decimal(10,2);not null;comment:'平台服务费（抽成金额）'"`
	DriverIncome        float64   `gorm:"type:decimal(10,2);not null;comment:'司机实际收入（计算公式：总订单金额-平台服务费-总退款金额）'"`
	RefundTotal         float64   `gorm:"type:decimal(10,2);not null;comment:'周期内总退款金额'"`
	AccountStatus       string    `gorm:"type:varchar(10);not null;comment:'账目状态（可选值：正常/异常/待核对，工单要求司机核对账单）'"`
	CheckStatus         string    `gorm:"type:varchar(10);not null;default:'未核对';comment:'司机核对状态（可选值：未核对/已核对/有异议，支撑工单争议解决流程）'"`
	CreateTime          time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:'结算单生成时间'"`
	UpdateTime          time.Time `gorm:"type:datetime;not null;default:current_timestamp;autoUpdateTime;comment:'更新时间（司机核对/异议处理时自动更新）'"`
}
