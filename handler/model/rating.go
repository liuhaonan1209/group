// Package model 数据模型层 - 评价管理模块
// 本文件定义评价管理、订单历史相关的数据库表结构
package model

import (
	"time"

	"gorm.io/gorm"
)

// RatingManage 评价管理表
// 用于存储乘客或司机对订单的评价信息，包括评分、标签、文字评价等
type RatingManage struct {
	gorm.Model                                                                                 // GORM基础模型（包含ID、CreatedAt、UpdatedAt、DeletedAt）
	OrderID    int64   `gorm:"not null;index;comment:订单ID" json:"order_id"`                   // 关联的订单ID
	UserID     int64   `gorm:"not null;index;comment:用户ID" json:"user_id"`                    // 评价人用户ID（乘客或司机）
	UserName   string  `gorm:"type:varchar(50);comment:用户姓名" json:"user_name"`                // 评价人姓名
	DriverID   *int64  `gorm:"index;comment:司机ID" json:"driver_id"`                           // 被评价的司机ID（如果是乘客评价司机）
	DriverName string  `gorm:"type:varchar(50);comment:司机姓名" json:"driver_name"`              // 被评价的司机姓名
	Score      float64 `gorm:"type:decimal(3,2);not null;comment:评分(1-5分)" json:"score"`     // 评分（1.0-5.0分，支持小数）
	Tags       string  `gorm:"type:varchar(500);comment:标签(JSON数组)" json:"tags"`             // 评价标签（JSON数组格式，如：["服务好","车辆干净","准点"]）
	Comment    string  `gorm:"type:text;comment:文字意见" json:"comment"`                         // 文字评价内容
	RaterType  string  `gorm:"type:varchar(20);not null;comment:评价人类型(passenger/driver)" json:"rater_type"` // 评价人类型：passenger-乘客评价，driver-司机评价
}

// TableName 指定数据库表名
// GORM会使用此方法返回的名称作为表名
func (RatingManage) TableName() string {
	return "rating_manage"
}

// OrderHistory 订单历史变更表
// 用于记录订单状态的每次变更，形成完整的订单生命周期追踪
type OrderHistory struct {
	gorm.Model                                                                   // GORM基础模型（包含ID、CreatedAt、UpdatedAt、DeletedAt）
	OrderID      int64     `gorm:"not null;index;comment:订单ID" json:"order_id"`  // 关联的订单ID
	Status       string    `gorm:"type:varchar(20);not null;comment:状态" json:"status"` // 订单状态（如：pending-待接单，accepted-已接单，completed-已完成，cancelled-已取消）
	Reason       string    `gorm:"type:text;comment:原因" json:"reason"`            // 状态变更原因（如：司机已接单、乘客取消等）
	OperatorID   *int64    `gorm:"index;comment:操作人ID" json:"operator_id"`       // 操作人ID（执行状态变更的用户）
	OperatorName string    `gorm:"type:varchar(50);comment:操作人姓名" json:"operator_name"` // 操作人姓名
	ChangedAt    time.Time `gorm:"type:datetime;not null;comment:变更时间" json:"changed_at"` // 状态变更时间
}

// TableName 指定数据库表名
// GORM会使用此方法返回的名称作为表名
func (OrderHistory) TableName() string {
	return "order_history"
}
