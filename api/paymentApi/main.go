// Package main 支付管理API服务
// 提供HTTP RESTful API接口，作为支付管理RPC服务的网关层
// 负责接收HTTP请求，转换为RPC调用，并返回HTTP响应
package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"group/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// main 主函数
// 初始化HTTP服务器，注册路由并启动服务
func main() {
	// 创建Hertz HTTP服务器
	// 监听localhost:9994端口
	hz := server.New(server.WithHostPorts("127.0.0.1:7777"))

	// 使用 CORS 中间件
	hz.Use(middleware.CORS())

	// 注册HTTP路由
	hz.POST("/api/payment/pay", Payment)                  // 支付接口
	hz.POST("/api/payment/issue", PaymentIssue)           // 支付问题处理
	hz.POST("/api/payment/records", PaymentRecordQuery)   // 支付记录查询
	hz.GET("/api/payment/status/:id", PaymentStatusQuery) // 支付状态查询

	log.Println("Payment API 服务启动在 localhost:7777")
	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

// Payment 支付接口
// HTTP POST /api/payment/pay
// 订单支付（含加密卡号）
func Payment(ctx context.Context, c *app.RequestContext) {
	// 定义请求体结构
	var reqBody struct {
		OrderID       int64  `json:"order_id"`
		UserID        int64  `json:"user_id"`
		PaymentMethod string `json:"payment_method"`
		CardInfo      string `json:"card_info,omitempty"`
	}

	// 解析JSON请求体
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 参数验证
	if reqBody.OrderID <= 0 || reqBody.UserID <= 0 || reqBody.PaymentMethod == "" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "订单ID、用户ID和支付方式不能为空",
		})
		return
	}

	// 模拟支付处理
	paymentID := time.Now().Unix()
	paymentNo := "PAY" + strconv.FormatInt(paymentID, 10)

	// 返回成功响应
	c.JSON(consts.StatusOK, utils.H{
		"code":    200,
		"message": "支付成功",
		"data": utils.H{
			"payment_id": paymentID,
			"payment_no": paymentNo,
			"status":     "success",
		},
	})
}

// PaymentIssue 支付问题处理接口
// HTTP POST /api/payment/issue
// 处理支付失败等订单问题
func PaymentIssue(ctx context.Context, c *app.RequestContext) {
	// 定义请求体结构
	var reqBody struct {
		OrderID     int64  `json:"order_id"`
		UserID      int64  `json:"user_id"`
		IssueType   string `json:"issue_type"`
		Description string `json:"description"`
	}

	// 解析JSON请求体
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 参数验证
	if reqBody.OrderID <= 0 || reqBody.UserID <= 0 || reqBody.IssueType == "" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "订单ID、用户ID和问题类型不能为空",
		})
		return
	}

	// 模拟问题处理
	issueID := time.Now().Unix()
	var solution string

	switch reqBody.IssueType {
	case "payment_failed":
		solution = "支付失败问题：建议检查支付方式或联系客服处理"
	case "refund_issue":
		solution = "退款问题：将在3-5个工作日内处理退款申请"
	default:
		solution = "其他问题：客服将尽快联系您处理"
	}

	// 返回成功响应
	c.JSON(consts.StatusOK, utils.H{
		"code":    200,
		"message": "问题已记录，客服将尽快处理",
		"data": utils.H{
			"issue_id": issueID,
			"solution": solution,
		},
	})
}

// PaymentRecordQuery 支付记录查询接口
// HTTP POST /api/payment/records
// 查询支付记录列表
func PaymentRecordQuery(ctx context.Context, c *app.RequestContext) {
	// 定义请求体结构
	var reqBody struct {
		UserID    *int64  `json:"user_id,omitempty"`
		OrderID   *int64  `json:"order_id,omitempty"`
		Status    *string `json:"status,omitempty"`
		StartDate *string `json:"start_date,omitempty"`
		EndDate   *string `json:"end_date,omitempty"`
		Page      *int32  `json:"page,omitempty"`
		PageSize  *int32  `json:"page_size,omitempty"`
	}

	// 解析JSON请求体
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 分页参数
	page := int32(1)
	pageSize := int32(10)
	if reqBody.Page != nil && *reqBody.Page > 0 {
		page = *reqBody.Page
	}
	if reqBody.PageSize != nil && *reqBody.PageSize > 0 {
		pageSize = *reqBody.PageSize
	}

	// 模拟支付记录数据
	records := []map[string]interface{}{
		{
			"payment_id":     1,
			"payment_no":     "PAY123456789",
			"order_id":       101,
			"user_id":        1001,
			"payment_method": "alipay",
			"amount":         50.00,
			"status":         "success",
			"payment_time":   "2024-01-26 10:00:00",
		},
		{
			"payment_id":     2,
			"payment_no":     "PAY123456790",
			"order_id":       102,
			"user_id":        1002,
			"payment_method": "wechat",
			"amount":         75.00,
			"status":         "success",
			"payment_time":   "2024-01-26 11:00:00",
		},
		{
			"payment_id":     3,
			"payment_no":     "PAY123456791",
			"order_id":       103,
			"user_id":        1003,
			"payment_method": "card",
			"amount":         100.00,
			"status":         "failed",
			"payment_time":   "2024-01-26 12:00:00",
			"fail_reason":    "银行卡余额不足",
		},
	}

	// 根据筛选条件过滤数据（简化处理）
	filteredRecords := records
	if reqBody.UserID != nil {
		var filtered []map[string]interface{}
		for _, record := range records {
			if record["user_id"].(int) == int(*reqBody.UserID) {
				filtered = append(filtered, record)
			}
		}
		filteredRecords = filtered
	}

	// 返回成功响应
	c.JSON(consts.StatusOK, utils.H{
		"code":    200,
		"message": "查询成功",
		"data": utils.H{
			"records":   filteredRecords,
			"total":     len(filteredRecords),
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// PaymentStatusQuery 支付状态查询接口
// HTTP GET /api/payment/status/:id?user_id=xxx
// 查询支付状态
func PaymentStatusQuery(ctx context.Context, c *app.RequestContext) {
	// 获取路径参数
	paymentIDStr := c.Param("id")
	userIDStr := c.Query("user_id")

	// 参数验证
	if paymentIDStr == "" || userIDStr == "" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "支付ID和用户ID为必填项",
		})
		return
	}

	// 转换参数类型
	paymentID, err := strconv.ParseInt(paymentIDStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "无效的支付ID",
		})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	// 模拟支付信息
	paymentInfo := map[string]interface{}{
		"payment_id":     paymentID,
		"payment_no":     "PAY" + paymentIDStr,
		"order_id":       123,
		"user_id":        userID,
		"payment_method": "alipay",
		"amount":         50.00,
		"status":         "success",
		"payment_time":   time.Now().Format("2006-01-02 15:04:05"),
	}

	// 返回成功响应
	c.JSON(consts.StatusOK, utils.H{
		"code":    200,
		"message": "查询成功",
		"data":    paymentInfo,
	})
}
