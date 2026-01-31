// Package model 数据模型层 - 支付管理模块
// 本文件定义支付管理相关的数据库表结构
package model

import (
	"time"

	"gorm.io/gorm"
)

// PaymentManage 支付管理表
// 用于存储订单的支付记录，包括支付方式、金额、状态等信息
type PaymentManage struct {
	gorm.Model                                                                                                  // GORM基础模型（包含ID、CreatedAt、UpdatedAt、DeletedAt）
	PaymentNo     string     `gorm:"type:varchar(50);uniqueIndex;not null;comment:支付单号" json:"payment_no"`      // 支付单号（唯一标识，格式：PAY+时间戳）
	OrderID       int64      `gorm:"not null;index;comment:订单ID" json:"order_id"`                                // 关联的订单ID
	UserID        int64      `gorm:"not null;index;comment:用户ID" json:"user_id"`                                 // 支付用户ID（通常是乘客）
	PaymentMethod string     `gorm:"type:varchar(20);not null;comment:支付方式(alipay/wechat/card)" json:"payment_method"` // 支付方式：alipay-支付宝，wechat-微信，card-银行卡
	Amount        float64    `gorm:"type:decimal(10,2);not null;comment:支付金额" json:"amount"`                     // 支付金额（单位：元）
	Status        string     `gorm:"type:varchar(20);default:'pending';comment:支付状态(pending/success/failed)" json:"status"` // 支付状态：pending-待支付，success-支付成功，failed-支付失败
	CardInfo      string     `gorm:"type:varchar(200);comment:加密的卡号信息" json:"card_info"`                          // 加密的银行卡号信息（使用MD5加密，保护用户隐私）
	FailReason    string     `gorm:"type:text;comment:失败原因" json:"fail_reason"`                                   // 支付失败原因（如：余额不足、银行卡过期等）
	PaymentTime   *time.Time `gorm:"type:datetime;comment:支付时间" json:"payment_time"`                             // 实际支付完成时间
}

// TableName 指定数据库表名
// GORM会使用此方法返回的名称作为表名
func (PaymentManage) TableName() string {
	return "payment_manage"
}

// PaymentIssue 支付问题表
// 用于记录支付过程中出现的问题和投诉，如支付失败、订单异常等
type PaymentIssue struct {
	gorm.Model                                                                                      // GORM基础模型（包含ID、CreatedAt、UpdatedAt、DeletedAt）
	OrderID     int64      `gorm:"not null;index;comment:订单ID" json:"order_id"`                    // 关联的订单ID
	UserID      int64      `gorm:"not null;index;comment:用户ID" json:"user_id"`                     // 提交问题的用户ID
	IssueType   string     `gorm:"type:varchar(50);not null;comment:问题类型" json:"issue_type"`      // 问题类型（如：payment_failed-支付失败，order_issue-订单问题）
	Description string     `gorm:"type:text;not null;comment:问题描述" json:"description"`             // 问题详细描述（用户填写）
	Solution    string     `gorm:"type:text;comment:解决方案" json:"solution"`                         // 客服提供的解决方案
	Status      string     `gorm:"type:varchar(20);default:'pending';comment:处理状态(pending/processing/resolved)" json:"status"` // 处理状态：pending-待处理，processing-处理中，resolved-已解决
	ResolvedAt  *time.Time `gorm:"type:datetime;comment:解决时间" json:"resolved_at"`                 // 问题解决时间
}

// TableName 指定数据库表名
// GORM会使用此方法返回的名称作为表名
func (PaymentIssue) TableName() string {
	return "payment_issue"
}
