package main

import (
	"context"
	"fmt"
	"group/handler/dao"
	"group/handler/model"
	"group/kitex_gen/car/refund"
	"log"
	"time"

	"gorm.io/gorm"
)

// RefundServiceImpl 退票服务实现
type RefundServiceImpl struct{}

// RefundApply 退票申请接口
func (s *RefundServiceImpl) RefundApply(ctx context.Context, req *refund.RefundApplyReq) (*refund.RefundApplyResp, error) {
	// 1. 验证订单是否存在且属于该用户（调用DAO层）
	orderData, err := dao.GetOrderByIDAndUserID(req.OrderId, req.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &refund.RefundApplyResp{
				Success: false,
				Message: "订单不存在或无权限",
			}, nil
		}
		log.Printf("查询订单失败: %v", err)
		return &refund.RefundApplyResp{
			Success: false,
			Message: "查询订单失败",
		}, nil
	}

	// 2. 检查订单状态是否允许退票
	if orderData.Status == "completed" || orderData.Status == "cancelled" {
		return &refund.RefundApplyResp{
			Success: false,
			Message: "订单已完成或已取消，无法申请退票",
		}, nil
	}

	// 3. 检查是否已有退票申请（调用DAO层）
	statuses := []string{"pending", "approved"}
	_, err = dao.GetRefundByOrderIDAndUserID(req.OrderId, req.UserId, statuses)
	if err == nil {
		return &refund.RefundApplyResp{
			Success: false,
			Message: "该订单已有退票申请正在处理中",
		}, nil
	}

	// 4. 创建退票申请（调用DAO层）
	refundData := model.RefundManage{
		OrderID:    req.OrderId,
		UserID:     req.UserId,
		UserType:   req.UserType,
		Reason:     req.Reason,
		ContactWay: getStringValue(req.ContactWay),
		Status:     "pending",
		ApplyTime:  time.Now(),
	}

	err = dao.CreateRefund(&refundData)
	if err != nil {
		log.Printf("创建退票申请失败: %v", err)
		return &refund.RefundApplyResp{
			Success: false,
			Message: "创建退票申请失败",
		}, nil
	}

	log.Printf("退票申请创建成功: 退票ID=%d", refundData.ID)
	return &refund.RefundApplyResp{
		RefundId: int64(refundData.ID),
		Success:  true,
		Message:  "退票申请提交成功",
	}, nil
}

// RefundAudit 退票审核接口
func (s *RefundServiceImpl) RefundAudit(ctx context.Context, req *refund.RefundAuditReq) (*refund.RefundAuditResp, error) {
	log.Printf("退票审核请求: 退票ID=%d, 审核人ID=%d, 审核结果=%s",
		req.RefundId, req.AuditUserId, req.AuditResult_)

	// 1. 查询退票申请（调用DAO层）
	refundData, err := dao.GetRefundByID(req.RefundId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &refund.RefundAuditResp{
				Success: false,
				Message: "退票申请不存在",
			}, nil
		}
		log.Printf("查询退票申请失败: %v", err)
		return &refund.RefundAuditResp{
			Success: false,
			Message: "查询退票申请失败",
		}, nil
	}

	// 2. 检查退票申请状态
	if refundData.Status != "pending" {
		return &refund.RefundAuditResp{
			Success: false,
			Message: "该退票申请已处理，无法重复审核",
		}, nil
	}

	// 3. 查询原订单信息（调用DAO层）
	orderData, err := dao.GetOrderByID(refundData.OrderID)
	if err != nil {
		log.Printf("查询原订单失败: %v", err)
		return &refund.RefundAuditResp{
			Success: false,
			Message: "查询原订单失败",
		}, nil
	}

	// 4. 计算退款金额
	var refundAmount *float64
	if req.AuditResult_ == "approved" {
		// 审核通过，计算退款金额（这里简化为全额退款）
		amount := orderData.Price
		refundAmount = &amount
	}

	// 5. 调用DAO层事务处理审核
	auditReason := getStringValue(req.AuditReason)
	err = dao.RefundAuditTransaction(req.RefundId, req.AuditUserId, req.AuditResult_, auditReason, refundAmount, refundData.OrderID)
	if err != nil {
		log.Printf("退票审核事务处理失败: %v", err)
		return &refund.RefundAuditResp{
			Success: false,
			Message: "退票审核失败",
		}, nil
	}

	// 6. 构造响应消息
	message := "退票审核完成"
	if req.AuditResult_ == "approved" {
		message = fmt.Sprintf("退票审核通过，退款金额: %.2f", *refundAmount)
	} else {
		message = "退票申请被拒绝"
	}

	resp := &refund.RefundAuditResp{
		Success: true,
		Message: message,
	}
	if refundAmount != nil {
		refundAmountStr := fmt.Sprintf("%.2f", *refundAmount)
		resp.RefundAmount = &refundAmountStr
	}

	log.Printf("退票审核完成: 退票ID=%d, 结果=%s", req.RefundId, req.AuditResult_)
	return resp, nil
}

// BatchRefund 批量退票接口
func (s *RefundServiceImpl) BatchRefund(ctx context.Context, req *refund.BatchRefundReq) (*refund.BatchRefundResp, error) {
	log.Printf("批量退票请求: 订单数量=%d, 操作人ID=%d, 退票类型=%s",
		len(req.OrderIds), req.OperatorId, req.RefundType)

	var results []*refund.BatchRefundResultItem
	successCount := int32(0)
	failCount := int32(0)

	// 处理每个订单
	for _, orderID := range req.OrderIds {
		result := s.processSingleRefund(orderID, req.OperatorId, req.RefundType, req.Reason, req.RefundRatio)
		results = append(results, result)

		if result.Success {
			successCount++
		} else {
			failCount++
		}
	}

	message := fmt.Sprintf("批量退票处理完成，成功: %d, 失败: %d", successCount, failCount)

	return &refund.BatchRefundResp{
		Results:      results,
		SuccessCount: successCount,
		FailCount:    failCount,
		Message:      message,
		Success:      successCount > 0,
	}, nil
}

// processSingleRefund 处理单个订单退票
func (s *RefundServiceImpl) processSingleRefund(orderID int64, operatorID int64, refundType string, reason string, refundRatio *float64) *refund.BatchRefundResultItem {
	// 1. 查询订单（调用DAO层）
	orderData, err := dao.GetOrderByID(orderID)
	if err != nil {
		return &refund.BatchRefundResultItem{
			OrderId: orderID,
			Success: false,
			Message: "订单不存在",
		}
	}

	// 2. 检查订单状态
	if orderData.Status == "cancelled" {
		return &refund.BatchRefundResultItem{
			OrderId: orderID,
			Success: false,
			Message: "订单已取消",
		}
	}

	// 3. 计算退款金额
	var refundAmount float64
	if refundType == "full" {
		refundAmount = orderData.Price
	} else if refundType == "partial" && refundRatio != nil {
		refundAmount = orderData.Price * (*refundRatio)
	} else {
		return &refund.BatchRefundResultItem{
			OrderId: orderID,
			Success: false,
			Message: "退款类型或比例参数错误",
		}
	}

	// 4. 创建退票记录（调用DAO层事务处理）
	now := time.Now()
	refundData := model.RefundManage{
		OrderID:       orderID,
		UserID:        orderData.PassengerID, // 默认使用乘客ID
		UserType:      "system",              // 系统批量操作
		Reason:        reason,
		Status:        "approved", // 批量退票直接通过
		AuditUserID:   &operatorID,
		AuditReason:   "系统批量退票",
		RefundAmount:  &refundAmount,
		ApplyTime:     now,
		AuditTime:     &now,
		CompletedTime: &now,
	}

	err = dao.BatchRefundTransaction(&refundData, orderID)
	if err != nil {
		return &refund.BatchRefundResultItem{
			OrderId: orderID,
			Success: false,
			Message: "退票处理失败",
		}
	}

	return &refund.BatchRefundResultItem{
		OrderId:      orderID,
		Success:      true,
		Message:      "退票成功",
		RefundAmount: &refundAmount,
	}
}

// RefundStatusQuery 退票状态查询接口
func (s *RefundServiceImpl) RefundStatusQuery(ctx context.Context, req *refund.RefundStatusQueryReq) (*refund.RefundStatusQueryResp, error) {
	log.Printf("退票状态查询请求: 退票ID=%d, 用户ID=%d", req.RefundId, req.UserId)

	// 查询退票申请（调用DAO层，带权限验证）
	refundData, err := dao.GetRefundByIDAndUserID(req.RefundId, req.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &refund.RefundStatusQueryResp{
				Success: false,
				Message: "退票申请不存在或无权限查看",
			}, nil
		}
		log.Printf("查询退票申请失败: %v", err)
		return &refund.RefundStatusQueryResp{
			Success: false,
			Message: "查询退票申请失败",
		}, nil
	}

	// 构建返回信息
	refundInfo := &refund.RefundStatusInfo{
		RefundId:  int64(refundData.ID),
		OrderId:   refundData.OrderID,
		UserId:    refundData.UserID,
		UserType:  refundData.UserType,
		Reason:    refundData.Reason,
		Status:    refundData.Status,
		ApplyTime: refundData.ApplyTime.Format("2006-01-02 15:04:05"),
	}

	if refundData.AuditReason != "" {
		refundInfo.AuditReason = &refundData.AuditReason
	}
	if refundData.RefundAmount != nil {
		refundInfo.RefundAmount = refundData.RefundAmount
	}
	if refundData.AuditTime != nil {
		auditTime := refundData.AuditTime.Format("2006-01-02 15:04:05")
		refundInfo.AuditTime = &auditTime
	}
	if refundData.CompletedTime != nil {
		completedTime := refundData.CompletedTime.Format("2006-01-02 15:04:05")
		refundInfo.CompletedTime = &completedTime
	}

	return &refund.RefundStatusQueryResp{
		RefundInfo: refundInfo,
		Success:    true,
		Message:    "查询成功",
	}, nil
}

// RefundRecordQuery 退票记录查询接口
func (s *RefundServiceImpl) RefundRecordQuery(ctx context.Context, req *refund.RefundRecordQueryReq) (*refund.RefundRecordQueryResp, error) {
	log.Printf("退票记录查询请求: 用户ID=%v, 用户类型=%v, 状态=%v",
		req.UserId, req.UserType, req.Status)

	// 分页参数
	page := int32(1)
	pageSize := int32(10)
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.PageSize != nil && *req.PageSize > 0 {
		pageSize = *req.PageSize
	}

	// 调用DAO层查询退票记录
	refundList, total, err := dao.QueryRefunds(req.UserId, req.UserType, req.Status, req.StartDate, req.EndDate, page, pageSize)
	if err != nil {
		log.Printf("查询退票记录失败: %v", err)
		return &refund.RefundRecordQueryResp{
			Success: false,
			Message: "查询失败",
		}, nil
	}

	// 构建返回数据
	var records []*refund.RefundStatusInfo
	for _, refundItem := range refundList {
		record := &refund.RefundStatusInfo{
			RefundId:  int64(refundItem.ID),
			OrderId:   refundItem.OrderID,
			UserId:    refundItem.UserID,
			UserType:  refundItem.UserType,
			Reason:    refundItem.Reason,
			Status:    refundItem.Status,
			ApplyTime: refundItem.ApplyTime.Format("2006-01-02 15:04:05"),
		}

		if refundItem.AuditReason != "" {
			record.AuditReason = &refundItem.AuditReason
		}
		if refundItem.RefundAmount != nil {
			record.RefundAmount = refundItem.RefundAmount
		}
		if refundItem.AuditTime != nil {
			auditTime := refundItem.AuditTime.Format("2006-01-02 15:04:05")
			record.AuditTime = &auditTime
		}
		if refundItem.CompletedTime != nil {
			completedTime := refundItem.CompletedTime.Format("2006-01-02 15:04:05")
			record.CompletedTime = &completedTime
		}

		records = append(records, record)
	}

	return &refund.RefundRecordQueryResp{
		Records:  records,
		Total:    int32(total),
		Page:     page,
		PageSize: pageSize,
		Success:  true,
		Message:  "查询成功",
	}, nil
}

// 辅助函数：获取字符串指针的值
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
