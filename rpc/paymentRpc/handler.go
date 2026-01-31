package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"group/handler/dao"
	"group/handler/model"
	"log"
	"time"

	"gorm.io/gorm"
)

// PaymentServiceImpl 支付服务实现
type PaymentServiceImpl struct{}

// Payment 支付接口
func (s *PaymentServiceImpl) Payment(ctx context.Context, req *PaymentReq) (*PaymentResp, error) {
	log.Printf("支付请求: 订单ID=%d, 用户ID=%d, 支付方式=%s", req.OrderId, req.UserId, req.PaymentMethod)

	// 1. 查询订单信息
	orderData, err := dao.GetOrderByID(req.OrderId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &PaymentResp{
				Success: false,
				Message: "订单不存在",
			}, nil
		}
		log.Printf("查询订单失败: %v", err)
		return &PaymentResp{
			Success: false,
			Message: "查询订单失败",
		}, nil
	}

	// 2. 检查订单状态
	if orderData.Status == "cancelled" {
		return &PaymentResp{
			Success: false,
			Message: "订单已取消，无法支付",
		}, nil
	}

	// 3. 检查是否已支付
	if orderData.PaymentStatus == "success" {
		return &PaymentResp{
			Success: false,
			Message: "订单已支付，请勿重复支付",
		}, nil
	}

	// 4. 生成支付单号
	paymentNo := generatePaymentNo()

	// 5. 模拟支付处理（实际应该调用第三方支付接口）
	paymentStatus := "success" // 模拟支付成功
	var failReason string
	now := time.Now()

	// 加密传输：支付相关信息通过SSL/TLS协议传输，保障数据安全
	cardInfo := ""
	if req.CardInfo != nil {
		cardInfo = encryptCardInfo(*req.CardInfo)
	}

	// 6. 创建支付记录
	paymentData := &model.PaymentManage{
		PaymentNo:     paymentNo,
		OrderID:       req.OrderId,
		UserID:        req.UserId,
		PaymentMethod: req.PaymentMethod,
		Amount:        orderData.Price,
		Status:        paymentStatus,
		CardInfo:      cardInfo,
		FailReason:    failReason,
		PaymentTime:   &now,
	}

	// 7. 使用事务处理支付
	err = dao.PaymentTransaction(paymentData, req.OrderId)
	if err != nil {
		log.Printf("支付处理失败: %v", err)
		return &PaymentResp{
			Success: false,
			Message: "支付处理失败",
		}, nil
	}

	log.Printf("支付成功: 支付ID=%d, 支付单号=%s", paymentData.ID, paymentNo)

	return &PaymentResp{
		PaymentId: int64(paymentData.ID),
		PaymentNo: paymentNo,
		Status:    paymentStatus,
		Success:   true,
		Message:   "支付成功",
	}, nil
}

// PaymentIssue 支付问题处理接口
func (s *PaymentServiceImpl) PaymentIssue(ctx context.Context, req *PaymentIssueReq) (*PaymentIssueResp, error) {
	log.Printf("支付问题处理请求: 订单ID=%d, 用户ID=%d, 问题类型=%s", req.OrderId, req.UserId, req.IssueType)

	// 1. 查询订单信息
	orderData, err := dao.GetOrderByID(req.OrderId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &PaymentIssueResp{
				Success: false,
				Message: "订单不存在",
			}, nil
		}
		log.Printf("查询订单失败: %v", err)
		return &PaymentIssueResp{
			Success: false,
			Message: "查询订单失败",
		}, nil
	}

	// 2. 验证用户权限
	if orderData.PassengerID != req.UserId && (orderData.DriverID == nil || *orderData.DriverID != req.UserId) {
		return &PaymentIssueResp{
			Success: false,
			Message: "无权处理此订单的支付问题",
		}, nil
	}

	// 3. 根据问题类型生成解决方案
	solution := generateSolution(req.IssueType, orderData)

	// 4. 创建问题记录
	issueData := &model.PaymentIssue{
		OrderID:     req.OrderId,
		UserID:      req.UserId,
		IssueType:   req.IssueType,
		Description: req.Description,
		Status:      "pending",
	}

	// 5. 使用事务处理问题
	err = dao.PaymentIssueTransaction(issueData, solution)
	if err != nil {
		log.Printf("创建支付问题记录失败: %v", err)
		return &PaymentIssueResp{
			Success: false,
			Message: "创建支付问题记录失败",
		}, nil
	}

	log.Printf("支付问题记录创建成功: 问题ID=%d", issueData.ID)

	return &PaymentIssueResp{
		IssueId:  int64(issueData.ID),
		Solution: solution,
		Success:  true,
		Message:  "问题已记录，客服将尽快处理",
	}, nil
}

// PaymentRecordQuery 支付记录查询接口
func (s *PaymentServiceImpl) PaymentRecordQuery(ctx context.Context, req *PaymentRecordQueryReq) (*PaymentRecordQueryResp, error) {
	log.Printf("支付记录查询请求: 用户ID=%v, 订单ID=%v, 状态=%v", req.UserId, req.OrderId, req.Status)

	// 分页参数
	page := int32(1)
	pageSize := int32(10)
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.PageSize != nil && *req.PageSize > 0 {
		pageSize = *req.PageSize
	}

	// 查询支付记录
	paymentList, total, err := dao.QueryPayments(
		req.UserId,
		req.OrderId,
		req.Status,
		req.StartDate,
		req.EndDate,
		page,
		pageSize,
	)
	if err != nil {
		log.Printf("查询支付记录失败: %v", err)
		return &PaymentRecordQueryResp{
			Success: false,
			Message: "查询失败",
		}, nil
	}

	// 构建返回数据
	var records []*PaymentRecordInfo
	for _, payment := range paymentList {
		record := &PaymentRecordInfo{
			PaymentId:     int64(payment.ID),
			PaymentNo:     payment.PaymentNo,
			OrderId:       payment.OrderID,
			UserId:        payment.UserID,
			PaymentMethod: payment.PaymentMethod,
			Amount:        payment.Amount,
			Status:        payment.Status,
			PaymentTime:   payment.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if payment.FailReason != "" {
			record.FailReason = &payment.FailReason
		}

		records = append(records, record)
	}

	return &PaymentRecordQueryResp{
		Records:  records,
		Total:    int32(total),
		Page:     page,
		PageSize: pageSize,
		Success:  true,
		Message:  "查询成功",
	}, nil
}

// PaymentStatusQuery 支付状态查询接口
func (s *PaymentServiceImpl) PaymentStatusQuery(ctx context.Context, req *PaymentStatusQueryReq) (*PaymentStatusQueryResp, error) {
	log.Printf("支付状态查询请求: 支付ID=%d, 用户ID=%d", req.PaymentId, req.UserId)

	// 查询支付记录（带权限验证）
	paymentData, err := dao.GetPaymentByIDAndUserID(req.PaymentId, req.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &PaymentStatusQueryResp{
				Success: false,
				Message: "支付记录不存在或无权限查看",
			}, nil
		}
		log.Printf("查询支付记录失败: %v", err)
		return &PaymentStatusQueryResp{
			Success: false,
			Message: "查询支付记录失败",
		}, nil
	}

	// 构建返回信息
	paymentInfo := &PaymentRecordInfo{
		PaymentId:     int64(paymentData.ID),
		PaymentNo:     paymentData.PaymentNo,
		OrderId:       paymentData.OrderID,
		UserId:        paymentData.UserID,
		PaymentMethod: paymentData.PaymentMethod,
		Amount:        paymentData.Amount,
		Status:        paymentData.Status,
		PaymentTime:   paymentData.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if paymentData.FailReason != "" {
		paymentInfo.FailReason = &paymentData.FailReason
	}

	return &PaymentStatusQueryResp{
		PaymentInfo: paymentInfo,
		Success:     true,
		Message:     "查询成功",
	}, nil
}

// 辅助函数：生成支付单号
func generatePaymentNo() string {
	return fmt.Sprintf("PAY%d", time.Now().UnixNano()/1000)
}

// 辅助函数：加密卡号信息（简化版，实际应使用更安全的加密算法）
func encryptCardInfo(cardInfo string) string {
	hash := md5.Sum([]byte(cardInfo))
	return hex.EncodeToString(hash[:])
}

// 辅助函数：生成解决方案
func generateSolution(issueType string, orderData model.OrderManage) string {
	switch issueType {
	case "payment_failed":
		return "针对支付失败：司机未按时到达等订单问题提供退款决方案。请联系客服处理退款事宜，预计1-3个工作日内完成退款。"
	case "refund_issue":
		return "针对退款问题：我们将尽快处理您的退款申请，预计3-5个工作日内退款到账。如有疑问请联系客服。"
	default:
		return "您的问题已记录，客服人员将在24小时内与您联系，为您提供解决方案。"
	}
}

// 临时定义结构体（实际应该从生成的kitex代码中导入）
type PaymentReq struct {
	OrderId       int64
	UserId        int64
	PaymentMethod string
	CardInfo      *string
}

type PaymentResp struct {
	PaymentId int64
	PaymentNo string
	Status    string
	Message   string
	Success   bool
}

type PaymentIssueReq struct {
	OrderId     int64
	UserId      int64
	IssueType   string
	Description string
}

type PaymentIssueResp struct {
	IssueId  int64
	Solution string
	Message  string
	Success  bool
}

type PaymentRecordQueryReq struct {
	UserId    *int64
	OrderId   *int64
	Status    *string
	StartDate *string
	EndDate   *string
	Page      *int32
	PageSize  *int32
}

type PaymentRecordInfo struct {
	PaymentId     int64
	PaymentNo     string
	OrderId       int64
	UserId        int64
	PaymentMethod string
	Amount        float64
	Status        string
	PaymentTime   string
	FailReason    *string
}

type PaymentRecordQueryResp struct {
	Records  []*PaymentRecordInfo
	Total    int32
	Page     int32
	PageSize int32
	Message  string
	Success  bool
}

type PaymentStatusQueryReq struct {
	PaymentId int64
	UserId    int64
}

type PaymentStatusQueryResp struct {
	PaymentInfo *PaymentRecordInfo
	Message     string
	Success     bool
}
