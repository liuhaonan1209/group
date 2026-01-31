// Package dao 数据访问层 - 退票管理模块
// 本文件包含所有退票相关的数据库操作函数
// 遵循DAO设计模式，将所有数据库操作封装在此层，供上层RPC Handler调用
package dao

import (
	"group/global"
	"group/handler/model"
	"time"

	"gorm.io/gorm"
)

// ==================== Refund 退票申请相关操作 ====================

// CreateRefund 创建退票申请记录
// 参数:
//   - data: 退票申请数据对象
//
// 返回:
//   - error: 创建失败时返回错误信息
func CreateRefund(data *model.RefundManage) error {
	return global.DB.Create(data).Error
}

// GetRefundByID 根据退票ID查询退票申请
// 参数:
//   - id: 退票申请ID
//
// 返回:
//   - refund: 退票申请对象
//   - err: 查询失败或记录不存在时返回错误
func GetRefundByID(id int64) (refund model.RefundManage, err error) {
	err = global.DB.Where("id = ?", id).First(&refund).Error
	return
}

// GetRefundByIDAndUserID 根据退票ID和用户ID查询退票申请（带权限验证）
// 用于确保用户只能查询自己的退票申请
// 参数:
//   - id: 退票申请ID
//   - userId: 用户ID
//
// 返回:
//   - refund: 退票申请对象
//   - err: 查询失败或无权限时返回错误
func GetRefundByIDAndUserID(id int64, userId int64) (refund model.RefundManage, err error) {
	err = global.DB.Where("id = ? AND user_id = ?", id, userId).First(&refund).Error
	return
}

// GetRefundByOrderIDAndUserID 根据订单ID和用户ID查询退票申请
// 用于检查某订单是否已有退票申请（防止重复申请）
// 参数:
//   - orderId: 订单ID
//   - userId: 用户ID
//   - statuses: 要查询的状态列表（如：pending, approved）
//
// 返回:
//   - refund: 退票申请对象
//   - err: 查询失败或记录不存在时返回错误
func GetRefundByOrderIDAndUserID(orderId int64, userId int64, statuses []string) (refund model.RefundManage, err error) {
	err = global.DB.Where("order_id = ? AND user_id = ? AND status IN ?", orderId, userId, statuses).
		First(&refund).Error
	return
}

// UpdateRefund 更新退票申请信息
// 参数:
//   - id: 退票申请ID
//   - updates: 需要更新的字段map，key为字段名，value为新值
//
// 返回:
//   - error: 更新失败时返回错误信息
func UpdateRefund(id int64, updates map[string]interface{}) error {
	return global.DB.Model(&model.RefundManage{}).Where("id = ?", id).Updates(updates).Error
}

// QueryRefunds 多条件查询退票记录列表（支持分页）
// 参数:
//   - userId: 用户ID（可选）
//   - userType: 用户类型（可选，如: passenger, driver）
//   - status: 退票状态（可选，如: pending, approved, rejected）
//   - startDate: 开始日期（可选）
//   - endDate: 结束日期（可选）
//   - page: 页码（从1开始）
//   - pageSize: 每页记录数
//
// 返回:
//   - refunds: 退票记录列表
//   - total: 符合条件的总记录数
//   - err: 查询失败时返回错误
func QueryRefunds(userId *int64, userType *string, status *string, startDate *string, endDate *string, page, pageSize int32) (refunds []model.RefundManage, total int64, err error) {
	query := global.DB.Model(&model.RefundManage{})

	// 动态构建查询条件
	if userId != nil {
		query = query.Where("user_id = ?", *userId)
	}
	if userType != nil {
		query = query.Where("user_type = ?", *userType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if startDate != nil {
		query = query.Where("apply_time >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("apply_time <= ?", *endDate)
	}

	// 获取符合条件的总记录数
	err = query.Count(&total).Error
	if err != nil {
		return
	}

	// 分页查询，按申请时间倒序
	offset := (page - 1) * pageSize
	err = query.Order("apply_time DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&refunds).Error
	return
}

// ==================== Order 订单相关操作（退票模块需要） ====================

// GetOrderByID 根据订单ID查询订单信息
// 参数:
//   - id: 订单ID
//
// 返回:
//   - order: 订单对象
//   - err: 查询失败或记录不存在时返回错误
func GetOrderByID(id int64) (order model.OrderManage, err error) {
	err = global.DB.Where("id = ?", id).First(&order).Error
	return
}

// GetOrderByIDAndUserID 根据订单ID和用户ID查询订单（带权限验证）
// 验证订单是否属于该用户（乘客或司机）
// 参数:
//   - orderId: 订单ID
//   - userId: 用户ID
//
// 返回:
//   - order: 订单对象
//   - err: 查询失败或无权限时返回错误
func GetOrderByIDAndUserID(orderId int64, userId int64) (order model.OrderManage, err error) {
	err = global.DB.Where("id = ? AND (passenger_id = ? OR driver_id = ?)", orderId, userId, userId).
		First(&order).Error
	return
}

// UpdateOrderStatus 更新订单状态
// 参数:
//   - id: 订单ID
//   - status: 新状态（如: cancelled）
//   - cancelledAt: 取消时间（可选）
//
// 返回:
//   - error: 更新失败时返回错误信息
func UpdateOrderStatus(id int64, status string, cancelledAt *time.Time) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if cancelledAt != nil {
		updates["cancelled_at"] = cancelledAt
	}
	return global.DB.Model(&model.OrderManage{}).Where("id = ?", id).Updates(updates).Error
}

// ==================== 退票审核事务操作 ====================

// RefundAuditTransaction 退票审核事务处理（原子操作）
// 功能：更新退票申请的审核结果，如果审核通过则同时取消订单
// 使用数据库事务确保数据一致性，任一步骤失败则全部回滚
// 参数:
//   - refundId: 退票申请ID
//   - auditUserId: 审核人ID
//   - auditResult: 审核结果（approved-通过, rejected-拒绝）
//   - auditReason: 审核理由
//   - refundAmount: 退款金额（审核通过时必填）
//   - orderId: 关联的订单ID
//
// 返回:
//   - error: 事务执行失败时返回错误信息
func RefundAuditTransaction(refundId int64, auditUserId int64, auditResult string, auditReason string, refundAmount *float64, orderId int64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// 步骤1: 更新退票申请的审核信息
		updates := map[string]interface{}{
			"audit_user_id": auditUserId,
			"audit_reason":  auditReason,
			"audit_time":    &now,
			"status":        auditResult,
		}

		// 如果审核通过，记录退款金额
		if auditResult == "approved" && refundAmount != nil {
			updates["refund_amount"] = *refundAmount
		}

		err := tx.Model(&model.RefundManage{}).Where("id = ?", refundId).Updates(updates).Error
		if err != nil {
			return err
		}

		// 步骤2: 如果审核通过，将订单状态更新为已取消
		if auditResult == "approved" {
			err = tx.Model(&model.OrderManage{}).Where("id = ?", orderId).
				Updates(map[string]interface{}{
					"status":       "cancelled",
					"cancelled_at": &now,
				}).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// ==================== 批量退票事务操作 ====================

// BatchRefundTransaction 批量退票事务处理（原子操作）
// 功能：创建退票记录并同时取消订单
// 用于批量处理退票场景（如线路取消、车辆故障等）
// 参数:
//   - refundData: 退票记录数据
//   - orderId: 关联的订单ID
//
// 返回:
//   - error: 事务执行失败时返回错误信息
func BatchRefundTransaction(refundData *model.RefundManage, orderId int64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// 步骤1: 创建退票记录
		err := tx.Create(refundData).Error
		if err != nil {
			return err
		}

		// 步骤2: 更新订单状态为已取消
		err = tx.Model(&model.OrderManage{}).Where("id = ?", orderId).
			Updates(map[string]interface{}{
				"status":       "cancelled",
				"cancelled_at": &now,
			}).Error
		if err != nil {
			return err
		}

		return nil
	})
}
