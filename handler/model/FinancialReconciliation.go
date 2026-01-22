package model

import (
	"time"

	"gorm.io/gorm"
)

// 收支对账主表
type FinancialReconciliation struct {
	gorm.Model
	TransactionDate time.Time              `gorm:"type:datetime;not null;comment:'交易日期时间戳（工单要求：所有交易信息实时同步，保证数据一致性）'"`
	OrderID         string                 `gorm:"type:varchar(32);not null;comment:'关联订单表的订单编号（工单要求：记录相关订单号）'"`
	DriverID        string                 `gorm:"type:varchar(32);not null;comment:'司机ID（工单要求：精确记录每笔交易的司机信息）'"`
	PassengerID     string                 `gorm:"type:varchar(32);not null;comment:'乘客ID（工单要求：精确记录每笔交易的乘客信息）'"`
	TripStart       string                 `gorm:"type:varchar(100);not null;comment:'行程起点（工单要求：记录行程明细）'"`
	TripEnd         string                 `gorm:"type:varchar(100);not null;comment:'行程终点（工单要求：记录行程明细）'"`
	TotalAmount     float64                `gorm:"type:decimal(10,2);not null;comment:'订单总金额（保留两位小数，工单要求所有数字保留两位小数）'"`
	PaymentMethod   string                 `gorm:"type:varchar(20);not null;comment:'支付方式（支持信用卡/电子钱包/现金等，工单要求区分不同支付方式）'"`
	DiscountAmount  float64                `gorm:"type:decimal(10,2);not null;default:0.00;comment:'优惠金额（优惠券/折扣抵扣，默认0元）'"`
	ActualReceived  float64                `gorm:"type:decimal(10,2);not null;comment:'订单实收金额（对应工单原型“订单实收”字段）'"`
	RefundAmount    float64                `gorm:"type:decimal(10,2);not null;default:0.00;comment:'订单退款金额（对应工单原型“订单退款”字段，默认0元）'"`
	SettledAmount   float64                `gorm:"type:decimal(10,2);not null;comment:'订单结算金额（计算公式：实收金额-退款金额，对应工单原型“订单结算”字段）'"`
	FeeDetails      map[string]interface{} `gorm:"type:json;comment:'费用明细（JSON格式存储，如：基础费、里程费、时长费等，工单要求记录费用明细）'"`
	AccountStatus   string                 `gorm:"type:varchar(10);not null;comment:'账目状态（可选值：正常/异常/争议，对应工单原型“账目状态”字段，支撑异常检测功能）'"`
	AbnormalReason  string                 `gorm:"type:varchar(200);comment:'异常原因（仅账目状态为“异常/争议”时填写，工单要求错误检测与修正）'"`
	DisputeStatus   string                 `gorm:"type:varchar(10);not null;default:'未处理';comment:'争议处理状态（可选值：未处理/处理中/已解决，工单要求设立争议解决流程）'"`
	CreateTime      time.Time              `gorm:"type:datetime;not null;default:current_timestamp;comment:'记录创建时间（默认当前时间）'"`
	UpdateTime      time.Time              `gorm:"type:datetime;not null;default:current_timestamp;autoUpdateTime;comment:'记录更新时间（异常处理/争议解决时自动更新）'"`
}
