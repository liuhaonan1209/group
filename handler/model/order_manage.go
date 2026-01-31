package model

import (
	"time"

	"gorm.io/gorm"
)

// OrderManage 订单表
type OrderManage struct {
	gorm.Model
	OrderNo       string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_no"`                   // 订单号
	TripID        int64      `gorm:"index;not null" json:"trip_id"`                                           // 行程ID
	OrderType     string     `gorm:"type:varchar(20);index;not null" json:"order_type"`                       // 订单类型（bus/carpool）
	PassengerID   int64      `gorm:"index;not null" json:"passenger_id"`                                      // 乘客ID
	PassengerName string     `gorm:"type:varchar(50);not null" json:"passenger_name"`                         // 乘客姓名
	PassengerTel  string     `gorm:"type:varchar(20);index;not null" json:"passenger_tel"`                    // 乘客手机号
	StartPoint    string     `gorm:"type:varchar(200);not null" json:"start_point"`                           // 起点
	EndPoint      string     `gorm:"type:varchar(200);not null" json:"end_point"`                             // 终点
	DepartureTime time.Time  `gorm:"index;not null" json:"departure_time"`                                    // 出发时间
	Status        string     `gorm:"type:varchar(20);index;not null;default:'pending'" json:"status"`         // 订单状态
	Price         float64    `gorm:"type:decimal(10,2);not null" json:"price"`                                // 价格
	PaymentStatus string     `gorm:"type:varchar(20);index;not null;default:'pending'" json:"payment_status"` // 支付状态
	PaymentMethod string     `gorm:"type:varchar(20)" json:"payment_method"`                                  // 支付方式
	PaymentTime   *time.Time `gorm:"type:datetime" json:"payment_time"`                                       // 支付时间
	DriverID      *int64     `gorm:"index" json:"driver_id"`                                                  // 司机ID
	DriverName    string     `gorm:"type:varchar(50)" json:"driver_name"`                                     // 司机姓名
	DriverTel     string     `gorm:"type:varchar(20)" json:"driver_tel"`                                      // 司机手机号
	VehicleInfo   string     `gorm:"type:varchar(200)" json:"vehicle_info"`                                   // 车辆信息
	RouteID       string     `gorm:"type:varchar(50);index" json:"route_id"`                                  // 线路ID
	SeatCount     int        `gorm:"not null;default:1" json:"seat_count"`                                    // 座位数量
	SpecialNeeds  string     `gorm:"type:text" json:"special_needs"`                                          // 特殊需求
	AcceptedAt    *time.Time `gorm:"type:datetime" json:"accepted_at"`                                        // 接单时间
	CompletedAt   *time.Time `gorm:"type:datetime" json:"completed_at"`                                       // 完成时间
	CancelledAt   *time.Time `gorm:"type:datetime" json:"cancelled_at"`                                       // 取消时间
}

// TableName 指定表名
func (OrderManage) TableName() string {
	return "order_manage"
}

// OrderStatusHistory 订单状态历史表
type OrderStatusHistory struct {
	gorm.Model
	OrderID      int64  `gorm:"index;not null" json:"order_id"`          // 订单ID
	Status       string `gorm:"type:varchar(20);not null" json:"status"` // 状态
	Reason       string `gorm:"type:text" json:"reason"`                 // 变更原因
	OperatorID   *int64 `gorm:"index" json:"operator_id"`                // 操作人ID
	OperatorName string `gorm:"type:varchar(50)" json:"operator_name"`   // 操作人姓名
	OperatorType string `gorm:"type:varchar(20)" json:"operator_type"`   // 操作人类型
}

// TableName 指定表名
func (OrderStatusHistory) TableName() string {
	return "order_status_history"
}
