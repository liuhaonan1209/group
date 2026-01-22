package model

import (
	"time"

	"gorm.io/gorm"
)

// 2. 线路结算表
type RouteSettlement struct {
	gorm.Model
	SettlementDate       time.Time `gorm:"type:date;not null;comment:'结算日期（工单要求：班车发车后开始结算，行程结束生成结算表）'"`
	RouteID              string    `gorm:"type:varchar(32);not null;comment:'线路编号（关联线路表，工单要求一条线路对应一个车队）'"`
	RouteName            string    `gorm:"type:varchar(100);not null;comment:'线路名称（对应工单原型“线路名称”字段）'"`
	FleetID              string    `gorm:"type:varchar(32);not null;comment:'所属车队ID（关联车队表，工单要求一条线路对应一个车队）'"`
	FleetName            string    `gorm:"type:varchar(50);not null;comment:'所属车队名称（对应工单原型“所属车队”字段）'"`
	TotalTicketPrice     float64   `gorm:"type:decimal(10,2);not null;comment:'当日线路订单总额（即工单原型“票价”字段，保留两位小数）'"`
	TotalDiscount        float64   `gorm:"type:decimal(10,2);not null;comment:'总优惠金额（优惠券扣除金额，对应工单原型“优惠金额”字段）'"`
	ReceivedChecked      float64   `gorm:"type:decimal(10,2);not null;comment:'实收金额（已检票，对应工单原型“实收金额(已检票)”字段）'"`
	TicketCountChecked   int       `gorm:"type:int;not null;comment:'已检票票数（对应工单原型“实收金额(已检票)/票数”中的票数）'"`
	ReceivedUnchecked    float64   `gorm:"type:decimal(10,2);not null;comment:'实收金额（未检票，对应工单原型“实收金额(未检票)”字段）'"`
	TicketCountUnchecked int       `gorm:"type:int;not null;comment:'未检票票数（对应工单原型“实收金额(未检票)/票数”中的票数）'"`
	RefundChecked        float64   `gorm:"type:decimal(10,2);not null;comment:'退款金额（已检票，对应工单原型“退款金额(已检票)”字段）'"`
	RefundCountChecked   int       `gorm:"type:int;not null;comment:'已检票退款票数（对应工单原型“退款金额(已检票)/票数”中的票数）'"`
	RefundUnchecked      float64   `gorm:"type:decimal(10,2);not null;comment:'退款金额（未检票，对应工单原型“退款金额(未检票)”字段）'"`
	RefundCountUnchecked int       `gorm:"type:int;not null;comment:'未检票退款票数（对应工单原型“退款金额(未检票)/票数”中的票数）'"`
	AccountStatus        string    `gorm:"type:varchar(10);not null;comment:'账目状态（可选值：正常/异常，对应工单原型“账目状态”字段）'"`
	CreateTime           time.Time `gorm:"type:datetime;not null;default:current_timestamp;comment:'结算表生成时间（行程结束后自动生成）'"`
}
