package model

import (
	"time"

	"gorm.io/gorm"
)

// 3. 站点结算表
type StationSettlement struct {
	gorm.Model
	SettlementDate       time.Time `gorm:"type:date;not null;comment:'结算日期（与线路结算表一致）'"`
	RouteID              string    `gorm:"type:varchar(32);not null;comment:'关联线路ID（关联线路表）'"`
	RouteName            string    `gorm:"type:varchar(100);not null;comment:'线路名称'"`
	StationName          string    `gorm:"type:varchar(50);not null;comment:'上车站点（对应工单原型“上车站点”字段，如：丽水云泉酒店）'"`
	FleetID              string    `gorm:"type:varchar(32);not null;comment:'所属车队ID（关联车队表）'"`
	TotalTicketPrice     float64   `gorm:"type:decimal(10,2);not null;comment:'当日站点订单总额（保留两位小数）'"`
	TotalDiscount        float64   `gorm:"type:decimal(10,2);not null;comment:'总优惠金额'"`
	ReceivedChecked      float64   `gorm:"type:decimal(10,2);not null;comment:'实收金额（已检票）'"`
	TicketCountChecked   int       `gorm:"type:int;not null;comment:'已检票票数'"`
	ReceivedUnchecked    float64   `gorm:"type:decimal(10,2);not null;comment:'实收金额（未检票）'"`
	TicketCountUnchecked int       `gorm:"type:int;not null;comment:'未检票票数'"`
	RefundChecked        float64   `gorm:"type:decimal(10,2);not null;comment:'退款金额（已检票）'"`
	RefundUnchecked      float64   `gorm:"type:decimal(10,2);not null;comment:'退款金额（未检票）'"`
	AccountStatus        string    `gorm:"type:varchar(10);not null;comment:'账目状态（可选值：正常/异常）'"`
	CreateTime           time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:'结算表生成时间'"`
}
