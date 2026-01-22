package model

import (
	"time"

	"gorm.io/gorm"
)

// 班次结算表
type ShiftSettlement struct {
	gorm.Model
	SettlementDate     time.Time `gorm:"type:date;not null;comment:'结算日期'"`
	ShiftID            string    `gorm:"type:varchar(32);not null;comment:'班次编号（关联班次表）'"`
	RouteID            string    `gorm:"type:varchar(32);not null;comment:'关联线路ID（关联线路表）'"`
	RouteName          string    `gorm:"type:varchar(100);not null;comment:'线路名称'"`
	FleetID            string    `gorm:"type:varchar(32);not null;comment:'所属车队ID（关联车队表）'"`
	DepartureTime      time.Time `gorm:"type:datetime;not null;comment:'发车时间（对应工单原型“结算日期(发车时间)”字段）'"`
	TotalTicketPrice   float64   `gorm:"type:decimal(10,2);not null;comment:'本班次订单总额（保留两位小数）'"`
	TotalDiscount      float64   `gorm:"type:decimal(10,2);not null;comment:'本班次总优惠金额'"`
	ReceivedChecked    float64   `gorm:"type:decimal(10,2);not null;comment:'本班次实收金额（已检票）'"`
	TicketCountChecked int       `gorm:"type:int;not null;comment:'本班次已检票票数'"`
	RefundAmount       float64   `gorm:"type:decimal(10,2);not null;comment:'本班次总退款金额'"`
	AccountStatus      string    `gorm:"type:varchar(10);not null;comment:'账目状态（可选值：正常/异常）'"`
	CreateTime         time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:'结算表生成时间（行程结束后自动生成）'"`
}
