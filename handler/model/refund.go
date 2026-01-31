// Package model 数据模型层 - 退票管理模块
// 本文件定义退票管理相关的数据库表结构
package model

import (
	"time"

	"gorm.io/gorm"
)

// RefundManage 退票管理表
// 用于存储用户的退票申请信息，包括申请、审核、退款等全流程数据
type RefundManage struct {
	gorm.Model                                                                                                                          // GORM基础模型（包含ID、CreatedAt、UpdatedAt、DeletedAt）
	OrderID       int64      `gorm:"not null;index" json:"order_id"`                                                                     // 订单ID（关联订单表）
	UserID        int64      `gorm:"not null;index" json:"user_id"`                                                                      // 申请用户ID（乘客或司机）
	UserType      string     `gorm:"type:varchar(20);not null;comment:用户类型(passenger/driver)" json:"user_type"`                          // 用户类型：passenger-乘客，driver-司机
	Reason        string     `gorm:"type:text;not null;comment:退票原因" json:"reason"`                                                      // 退票原因（用户填写）
	ContactWay    string     `gorm:"type:varchar(100);comment:联系方式" json:"contact_way"`                                                  // 联系方式（手机号或邮箱）
	Status        string     `gorm:"type:varchar(20);default:'pending';comment:退票状态(pending/approved/rejected/completed)" json:"status"` // 退票状态：pending-待审核，approved-已通过，rejected-已拒绝，completed-已完成
	AuditUserID   *int64     `gorm:"comment:审核人ID" json:"audit_user_id"`                                                                 // 审核人ID（管理员或客服）
	AuditReason   string     `gorm:"type:text;comment:审核意见" json:"audit_reason"`                                                         // 审核意见（通过或拒绝的理由）
	RefundAmount  *float64   `gorm:"type:decimal(10,2);comment:退款金额" json:"refund_amount"`                                               // 实际退款金额（可能扣除手续费）
	ApplyTime     time.Time  `gorm:"type:datetime;not null;comment:申请时间" json:"apply_time"`                                              // 退票申请提交时间
	AuditTime     *time.Time `gorm:"type:datetime;comment:审核时间" json:"audit_time"`                                                       // 审核完成时间
	CompletedTime *time.Time `gorm:"type:datetime;comment:完成时间" json:"completed_time"`                                                   // 退款完成时间（实际到账时间）
}

// TableName 指定数据库表名
// GORM会使用此方法返回的名称作为表名
func (RefundManage) TableName() string {
	return "refund_manage"
}
