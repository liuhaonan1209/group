package model

import (
	"time"

	"gorm.io/gorm"
)

// BalanceSheet 收支对账表
type BalanceSheet struct {
	gorm.Model
	SettleDate   time.Time `gorm:"type:date;index" json:"settleDate"`        // 对账日期
	OrderIncome  float64   `gorm:"type:decimal(12,2)" json:"orderIncome"`    // 订单实收
	OrderRefund  float64   `gorm:"type:decimal(12,2)" json:"orderRefund"`    // 订单退款
	OrderSettle  float64   `gorm:"type:decimal(12,2)" json:"orderSettle"`    // 订单结算
	Status       int       `gorm:"type:tinyint;default:1" json:"status"`     // 1-正常，2-异常
}

// IncomeSheet 收入对账表
type IncomeSheet struct {
	gorm.Model
	SettleDate            time.Time `gorm:"type:date;index" json:"settleDate"`
	RouteId               uint      `gorm:"type:int;index" json:"routeId"`
	RouteNo               string    `gorm:"type:varchar(50)" json:"routeNo"`
	RouteName             string    `gorm:"type:varchar(100)" json:"routeName"`
	Fleet                 string    `gorm:"type:varchar(50)" json:"fleet"`
	TicketPrice           float64   `gorm:"type:decimal(10,2)" json:"ticketPrice"`
	DiscountAmount        float64   `gorm:"type:decimal(10,2)" json:"discountAmount"`
	CheckedIncome         float64   `gorm:"type:decimal(12,2)" json:"checkedIncome"`
	CheckedTickets        int       `gorm:"type:int" json:"checkedTickets"`
	UncheckedIncome       float64   `gorm:"type:decimal(12,2)" json:"uncheckedIncome"`
	UncheckedTickets      int       `gorm:"type:int" json:"uncheckedTickets"`
	CheckedRefund         float64   `gorm:"type:decimal(12,2)" json:"checkedRefund"`
	CheckedRefundTickets  int       `gorm:"type:int" json:"checkedRefundTickets"`
	UncheckedRefund       float64   `gorm:"type:decimal(12,2)" json:"uncheckedRefund"`
	UncheckedRefundTickets int      `gorm:"type:int" json:"uncheckedRefundTickets"`
	Status                int       `gorm:"type:tinyint;default:1" json:"status"`
}

// RouteSettle 线路结算表
type RouteSettle struct {
	gorm.Model
	SettleDate            time.Time `gorm:"type:date;index" json:"settleDate"`
	RouteId               uint      `gorm:"type:int;index" json:"routeId"`
	RouteNo               string    `gorm:"type:varchar(50)" json:"routeNo"`
	RouteName             string    `gorm:"type:varchar(100)" json:"routeName"`
	Fleet                 string    `gorm:"type:varchar(50)" json:"fleet"`
	TicketPrice           float64   `gorm:"type:decimal(10,2)" json:"ticketPrice"`
	DiscountAmount        float64   `gorm:"type:decimal(10,2)" json:"discountAmount"`
	CheckedIncome         float64   `gorm:"type:decimal(12,2)" json:"checkedIncome"`
	CheckedTickets        int       `gorm:"type:int" json:"checkedTickets"`
	UncheckedIncome       float64   `gorm:"type:decimal(12,2)" json:"uncheckedIncome"`
	UncheckedTickets      int       `gorm:"type:int" json:"uncheckedTickets"`
	CheckedRefund         float64   `gorm:"type:decimal(12,2)" json:"checkedRefund"`
	CheckedRefundTickets  int       `gorm:"type:int" json:"checkedRefundTickets"`
	UncheckedRefund       float64   `gorm:"type:decimal(12,2)" json:"uncheckedRefund"`
	UncheckedRefundTickets int      `gorm:"type:int" json:"uncheckedRefundTickets"`
}

// ScheduleSettle 班次结算表
type ScheduleSettle struct {
	gorm.Model
	SettleDate            time.Time `gorm:"type:date;index" json:"settleDate"`
	RouteId               uint      `gorm:"type:int;index" json:"routeId"`
	ScheduleId            uint      `gorm:"type:int;index" json:"scheduleId"`
	RouteNo               string    `gorm:"type:varchar(50)" json:"routeNo"`
	RouteName             string    `gorm:"type:varchar(100)" json:"routeName"`
	DepartureTime         string    `gorm:"type:varchar(10)" json:"departureTime"`
	TicketPrice           float64   `gorm:"type:decimal(10,2)" json:"ticketPrice"`
	DiscountAmount        float64   `gorm:"type:decimal(10,2)" json:"discountAmount"`
	CheckedIncome         float64   `gorm:"type:decimal(12,2)" json:"checkedIncome"`
	CheckedTickets        int       `gorm:"type:int" json:"checkedTickets"`
	UncheckedIncome       float64   `gorm:"type:decimal(12,2)" json:"uncheckedIncome"`
	UncheckedTickets      int       `gorm:"type:int" json:"uncheckedTickets"`
	CheckedRefund         float64   `gorm:"type:decimal(12,2)" json:"checkedRefund"`
	CheckedRefundTickets  int       `gorm:"type:int" json:"checkedRefundTickets"`
	UncheckedRefund       float64   `gorm:"type:decimal(12,2)" json:"uncheckedRefund"`
	UncheckedRefundTickets int      `gorm:"type:int" json:"uncheckedRefundTickets"`
}

// StationSettle 站点结算表
type StationSettle struct {
	gorm.Model
	SettleDate            time.Time `gorm:"type:date;index" json:"settleDate"`
	StationId             uint      `gorm:"type:int;index" json:"stationId"`
	StationName           string    `gorm:"type:varchar(100)" json:"stationName"`
	TicketPrice           float64   `gorm:"type:decimal(10,2)" json:"ticketPrice"`
	DiscountAmount        float64   `gorm:"type:decimal(10,2)" json:"discountAmount"`
	CheckedIncome         float64   `gorm:"type:decimal(12,2)" json:"checkedIncome"`
	CheckedTickets        int       `gorm:"type:int" json:"checkedTickets"`
	UncheckedIncome       float64   `gorm:"type:decimal(12,2)" json:"uncheckedIncome"`
	UncheckedTickets      int       `gorm:"type:int" json:"uncheckedTickets"`
	CheckedRefund         float64   `gorm:"type:decimal(12,2)" json:"checkedRefund"`
	CheckedRefundTickets  int       `gorm:"type:int" json:"checkedRefundTickets"`
	UncheckedRefund       float64   `gorm:"type:decimal(12,2)" json:"uncheckedRefund"`
	UncheckedRefundTickets int      `gorm:"type:int" json:"uncheckedRefundTickets"`
}

// DriverSettle 司机结算表
type DriverSettle struct {
	gorm.Model
	SettleDate  time.Time `gorm:"type:date;index" json:"settleDate"`
	DriverId    uint      `gorm:"type:int;index" json:"driverId"`
	DriverName  string    `gorm:"type:varchar(50)" json:"driverName"`
	Phone       string    `gorm:"type:varchar(20)" json:"phone"`
	TripCount   int       `gorm:"type:int" json:"tripCount"`
	TotalIncome float64   `gorm:"type:decimal(12,2)" json:"totalIncome"`
	Commission  float64   `gorm:"type:decimal(12,2)" json:"commission"`
	NetIncome   float64   `gorm:"type:decimal(12,2)" json:"netIncome"`
	Status      int       `gorm:"type:tinyint;default:1" json:"status"` // 1-待结算，2-已结算
}

// Transaction 交易明细表
type Transaction struct {
	gorm.Model
	TransDate     time.Time `gorm:"type:date;index" json:"transDate"`
	TransTime     time.Time `gorm:"type:datetime" json:"transTime"`
	OrderNo       string    `gorm:"type:varchar(50);index" json:"orderNo"`
	DriverId      uint      `gorm:"type:int" json:"driverId"`
	DriverName    string    `gorm:"type:varchar(50)" json:"driverName"`
	PassengerId   uint      `gorm:"type:int" json:"passengerId"`
	PassengerName string    `gorm:"type:varchar(50)" json:"passengerName"`
	StartStation  string    `gorm:"type:varchar(100)" json:"startStation"`
	EndStation    string    `gorm:"type:varchar(100)" json:"endStation"`
	Amount        float64   `gorm:"type:decimal(10,2)" json:"amount"`
	PaymentMethod int       `gorm:"type:tinyint" json:"paymentMethod"` // 1-微信，2-支付宝，3-银行卡
	TransType     int       `gorm:"type:tinyint" json:"transType"`     // 1-收入，2-支出，3-退款
	Remark        string    `gorm:"type:varchar(255)" json:"remark"`
}

// AbnormalTransaction 异常交易表
type AbnormalTransaction struct {
	gorm.Model
	TransDate      time.Time  `gorm:"type:date;index" json:"transDate"`
	OrderNo        string     `gorm:"type:varchar(50);index" json:"orderNo"`
	AbnormalType   int        `gorm:"type:tinyint" json:"abnormalType"` // 1-重复扣款，2-金额不符，3-其他
	ExpectedAmount float64    `gorm:"type:decimal(10,2)" json:"expectedAmount"`
	ActualAmount   float64    `gorm:"type:decimal(10,2)" json:"actualAmount"`
	DiffAmount     float64    `gorm:"type:decimal(10,2)" json:"diffAmount"`
	Status         int        `gorm:"type:tinyint;default:1" json:"status"` // 1-待处理，2-已处理
	Handler        string     `gorm:"type:varchar(50)" json:"handler"`
	HandleTime     *time.Time `gorm:"type:datetime" json:"handleTime"`
	HandleRemark   string     `gorm:"type:varchar(255)" json:"handleRemark"`
	HandleType     int        `gorm:"type:tinyint" json:"handleType"` // 1-退款，2-补收，3-忽略
}

// Bill 账单表
type Bill struct {
	gorm.Model
	BillNo      string    `gorm:"type:varchar(50);unique" json:"billNo"`
	BillType    int       `gorm:"type:tinyint" json:"billType"` // 1-日账单，2-周账单，3-月账单
	StartDate   time.Time `gorm:"type:date" json:"startDate"`
	EndDate     time.Time `gorm:"type:date" json:"endDate"`
	TotalIncome float64   `gorm:"type:decimal(12,2)" json:"totalIncome"`
	TotalRefund float64   `gorm:"type:decimal(12,2)" json:"totalRefund"`
	TotalSettle float64   `gorm:"type:decimal(12,2)" json:"totalSettle"`
	Status      int       `gorm:"type:tinyint;default:1" json:"status"`
}

// 常量定义
const (
	// 支付方式
	PaymentMethodWechat  = 1
	PaymentMethodAlipay  = 2
	PaymentMethodBankCard = 3

	// 交易类型
	TransTypeIncome  = 1
	TransTypeExpense = 2
	TransTypeRefund  = 3

	// 异常类型
	AbnormalTypeDuplicate = 1
	AbnormalTypeAmountErr = 2
	AbnormalTypeOther     = 3

	// 处理方式
	HandleTypeRefund = 1
	HandleTypeCharge = 2
	HandleTypeIgnore = 3

	// 账单类型
	BillTypeDaily   = 1
	BillTypeWeekly  = 2
	BillTypeMonthly = 3

	// 结算状态
	SettleStatusPending = 1
	SettleStatusDone    = 2

	// 账目状态
	BalanceStatusNormal   = 1
	BalanceStatusAbnormal = 2
)

// GetPaymentMethodText 获取支付方式文本
func GetPaymentMethodText(method int) string {
	switch method {
	case PaymentMethodWechat:
		return "微信支付"
	case PaymentMethodAlipay:
		return "支付宝"
	case PaymentMethodBankCard:
		return "银行卡"
	default:
		return "未知"
	}
}

// GetTransTypeText 获取交易类型文本
func GetTransTypeText(transType int) string {
	switch transType {
	case TransTypeIncome:
		return "收入"
	case TransTypeExpense:
		return "支出"
	case TransTypeRefund:
		return "退款"
	default:
		return "未知"
	}
}

// GetAbnormalTypeText 获取异常类型文本
func GetAbnormalTypeText(abnormalType int) string {
	switch abnormalType {
	case AbnormalTypeDuplicate:
		return "重复扣款"
	case AbnormalTypeAmountErr:
		return "金额不符"
	case AbnormalTypeOther:
		return "其他"
	default:
		return "未知"
	}
}

// GetBalanceStatusText 获取账目状态文本
func GetBalanceStatusText(status int) string {
	if status == BalanceStatusNormal {
		return "正常"
	}
	return "异常"
}

// GetSettleStatusText 获取结算状态文本
func GetSettleStatusText(status int) string {
	if status == SettleStatusPending {
		return "待结算"
	}
	return "已结算"
}
