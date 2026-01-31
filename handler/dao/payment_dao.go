// Package dao 数据访问层 - 支付管理模块
// 本文件包含所有支付相关的数据库操作函数
// 遵循DAO设计模式，将所有数据库操作封装在此层，供上层RPC Handler调用
package dao

import (
	"group/global"
	"group/handler/model"
	"time"

	"gorm.io/gorm"
)

// ==================== Payment 支付记录相关操作 ====================

// CreatePayment 创建支付记录
// 参数:
//   - data: 支付记录数据对象
// 返回:
//   - error: 创建失败时返回错误信息
func CreatePayment(data *model.PaymentManage) error {
	return global.DB.Create(data).Error
}

// GetPaymentByID 根据支付ID查询支付记录
// 参数:
//   - id: 支付记录ID
// 返回:
//   - payment: 支付记录对象
//   - err: 查询失败或记录不存在时返回错误
func GetPaymentByID(id int64) (payment model.PaymentManage, err error) {
	err = global.DB.Where("id = ?", id).First(&payment).Error
	return
}

// GetPaymentByIDAndUserID 根据支付ID和用户ID查询支付记录（带权限验证）
// 用于确保用户只能查询自己的支付记录
// 参数:
//   - id: 支付记录ID
//   - userId: 用户ID
// 返回:
//   - payment: 支付记录对象
//   - err: 查询失败或无权限时返回错误
func GetPaymentByIDAndUserID(id int64, userId int64) (payment model.PaymentManage, err error) {
	err = global.DB.Where("id = ? AND user_id = ?", id, userId).First(&payment).Error
	return
}

// GetPaymentByOrderID 根据订单ID查询支付记录
// 按创建时间倒序，返回最新的支付记录
// 参数:
//   - orderId: 订单ID
// 返回:
//   - payment: 支付记录对象
//   - err: 查询失败或记录不存在时返回错误
func GetPaymentByOrderID(orderId int64) (payment model.PaymentManage, err error) {
	err = global.DB.Where("order_id = ?", orderId).Order("created_at DESC").First(&payment).Error
	return
}

// UpdatePayment 更新支付记录
// 参数:
//   - id: 支付记录ID
//   - updates: 需要更新的字段map，key为字段名，value为新值
// 返回:
//   - error: 更新失败时返回错误信息
func UpdatePayment(id int64, updates map[string]interface{}) error {
	return global.DB.Model(&model.PaymentManage{}).Where("id = ?", id).Updates(updates).Error
}

// QueryPayments 多条件查询支付记录列表（支持分页）
// 参数:
//   - userId: 用户ID（可选）
//   - orderId: 订单ID（可选）
//   - status: 支付状态（可选，如: success, failed, pending）
//   - startDate: 开始日期（可选）
//   - endDate: 结束日期（可选）
//   - page: 页码（从1开始）
//   - pageSize: 每页记录数
// 返回:
//   - payments: 支付记录列表
//   - total: 符合条件的总记录数
//   - err: 查询失败时返回错误
func QueryPayments(userId *int64, orderId *int64, status *string, startDate *string, endDate *string, page, pageSize int32) (payments []model.PaymentManage, total int64, err error) {
	query := global.DB.Model(&model.PaymentManage{})

	// 动态构建查询条件
	if userId != nil {
		query = query.Where("user_id = ?", *userId)
	}
	if orderId != nil {
		query = query.Where("order_id = ?", *orderId)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	// 获取符合条件的总记录数
	err = query.Count(&total).Error
	if err != nil {
		return
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&payments).Error
	return
}

// ==================== PaymentIssue 支付问题相关操作 ====================

// CreatePaymentIssue 创建支付问题记录
// 用于记录支付失败、订单问题等异常情况
// 参数:
//   - data: 支付问题数据对象
// 返回:
//   - error: 创建失败时返回错误信息
func CreatePaymentIssue(data *model.PaymentIssue) error {
	return global.DB.Create(data).Error
}

// GetPaymentIssueByID 根据问题ID查询支付问题记录
// 参数:
//   - id: 问题记录ID
// 返回:
//   - issue: 支付问题对象
//   - err: 查询失败或记录不存在时返回错误
func GetPaymentIssueByID(id int64) (issue model.PaymentIssue, err error) {
	err = global.DB.Where("id = ?", id).First(&issue).Error
	return
}

// UpdatePaymentIssue 更新支付问题记录
// 用于更新问题处理状态、解决方案等
// 参数:
//   - id: 问题记录ID
//   - updates: 需要更新的字段map
// 返回:
//   - error: 更新失败时返回错误信息
func UpdatePaymentIssue(id int64, updates map[string]interface{}) error {
	return global.DB.Model(&model.PaymentIssue{}).Where("id = ?", id).Updates(updates).Error
}

// GetPaymentIssuesByUserID 根据用户ID查询支付问题列表
// 按创建时间倒序返回该用户的所有问题记录
// 参数:
//   - userId: 用户ID
// 返回:
//   - issues: 支付问题列表
//   - err: 查询失败时返回错误
func GetPaymentIssuesByUserID(userId int64) (issues []model.PaymentIssue, err error) {
	err = global.DB.Where("user_id = ?", userId).Order("created_at DESC").Find(&issues).Error
	return
}

// ==================== 支付事务操作 ====================

// PaymentTransaction 支付事务处理（原子操作）
// 功能：创建支付记录并同步更新订单的支付状态
// 使用数据库事务确保数据一致性，任一步骤失败则全部回滚
// 参数:
//   - paymentData: 支付记录数据
//   - orderId: 关联的订单ID
// 返回:
//   - error: 事务执行失败时返回错误信息
func PaymentTransaction(paymentData *model.PaymentManage, orderId int64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 步骤1: 创建支付记录
		err := tx.Create(paymentData).Error
		if err != nil {
			return err
		}

		// 步骤2: 更新订单的支付状态和支付方式
		now := time.Now()
		updates := map[string]interface{}{
			"payment_status": paymentData.Status,
			"payment_method": paymentData.PaymentMethod,
		}

		// 如果支付成功，记录支付时间
		if paymentData.Status == "success" {
			updates["payment_time"] = &now
		}

		err = tx.Model(&model.OrderManage{}).Where("id = ?", orderId).Updates(updates).Error
		if err != nil {
			return err
		}

		return nil
	})
}

// PaymentIssueTransaction 支付问题处理事务
// 功能：创建问题记录，如果有解决方案则立即更新状态
// 参数:
//   - issueData: 支付问题数据
//   - solution: 解决方案（可为空）
// 返回:
//   - error: 事务执行失败时返回错误信息
func PaymentIssueTransaction(issueData *model.PaymentIssue, solution string) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 步骤1: 创建问题记录
		err := tx.Create(issueData).Error
		if err != nil {
			return err
		}

		// 步骤2: 如果提供了解决方案，立即更新问题状态
		if solution != "" {
			now := time.Now()
			err = tx.Model(&model.PaymentIssue{}).Where("id = ?", issueData.ID).
				Updates(map[string]interface{}{
					"solution":    solution,
					"status":      "processing", // 标记为处理中
					"resolved_at": &now,
				}).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}
