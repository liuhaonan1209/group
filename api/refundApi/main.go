// Package main 退票管理API服务
// 提供HTTP RESTful API接口，作为退票管理RPC服务的网关层
// 负责接收HTTP请求，转换为RPC调用，并返回HTTP响应
package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"group/kitex_gen/car/refund"
	"group/kitex_gen/car/refund/refundservice"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

var (
	// cli 退票管理服务RPC客户端
	// 用于调用退票管理RPC服务的接口
	cli refundservice.Client
)

// main 主函数
// 初始化RPC客户端和HTTP服务器，注册路由并启动服务
func main() {
	// 创建退票管理服务RPC客户端
	// 连接到127.0.0.1:8996端口的RPC服务
	c, err := refundservice.NewClient("car.refund", client.WithHostPorts("127.0.0.1:8996"))
	if err != nil {
		log.Fatal(err)
	}
	cli = c

	// 创建Hertz HTTP服务器
	// 监听localhost:9993端口
	hz := server.New(server.WithHostPorts("127.0.0.1:9993"))

	// 注册HTTP路由
	hz.POST("/api/refund/apply", RefundApply)                // 退票申请
	hz.POST("/api/refund/audit", RefundAudit)                // 退票审核
	hz.POST("/api/refund/batch", BatchRefund)                // 批量退票
	hz.GET("/api/refund/status/:id", RefundStatusQuery)      // 退票状态查询
	hz.POST("/api/refund/records", RefundRecordQuery)        // 退票记录查询

	log.Println("Refund API 服务启动在 localhost:9993")
	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

// RefundApply 退票申请接口
// HTTP POST /api/refund/apply
// 用户提交退票申请
// 请求体:
//   - order_id: 订单ID
//   - user_id: 用户ID
//   - user_type: 用户类型（passenger/driver）
//   - reason: 退票原因
//   - contact_way: 联系方式（可选）
//
// 返回:
//   - 200: 申请成功
//   - 400: 参数错误或申请失败
//   - 500: 服务器错误
func RefundApply(ctx context.Context, c *app.RequestContext) {
	// 定义请求体结构
	var reqBody struct {
		OrderID    int64  `json:"order_id"`
		UserID     int64  `json:"user_id"`
		UserType   string `json:"user_type"`
		Reason     string `json:"reason"`
		ContactWay string `json:"contact_way,omitempty"`
	}

	// 解析JSON请求体
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 构造RPC请求
	req := refund.NewRefundApplyReq()
	req.OrderId = reqBody.OrderID
	req.UserId = reqBody.UserID
	req.UserType = reqBody.UserType
	req.Reason = reqBody.Reason

	if reqBody.ContactWay != "" {
		req.ContactWay = &reqBody.ContactWay
	}

	// 调用RPC服务提交退票申请
	resp, err := cli.RefundApply(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "提交退票申请失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查申请结果
	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	// 返回成功响应
	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data": map[string]interface{}{
			"refund_id": resp.RefundId,
		},
	})
}

// RefundAudit 退票审核接口
// HTTP POST /api/refund/audit
// 审核人员审核退票申请
// 请求体:
//   - refund_id: 退票申请ID
//   - audit_user_id: 审核人ID
//   - audit_result: 审核结果（approved/rejected）
//   - audit_reason: 审核意见（可选）
//
// 返回:
//   - 200: 审核成功
//   - 400: 参数错误或审核失败
//   - 500: 服务器错误
func RefundAudit(ctx context.Context, c *app.RequestContext) {
	// 定义请求体结构
	var reqBody struct {
		RefundID     int64  `json:"refund_id"`
		AuditUserID  int64  `json:"audit_user_id"`
		AuditResult  string `json:"audit_result"`
		AuditReason  string `json:"audit_reason,omitempty"`
	}

	// 解析JSON请求体
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 构造RPC请求
	req := refund.NewRefundAuditReq()
	req.RefundId = reqBody.RefundID
	req.AuditUserId = reqBody.AuditUserID
	req.AuditResult_ = reqBody.AuditResult

	if reqBody.AuditReason != "" {
		req.AuditReason = &reqBody.AuditReason
	}

	// 调用RPC服务审核退票申请
	resp, err := cli.RefundAudit(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "审核退票申请失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查审核结果
	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	// 返回成功响应
	result := map[string]interface{}{
		"code":    200,
		"message": resp.Message,
	}

	if resp.RefundAmount != nil {
		result["data"] = map[string]interface{}{
			"refund_amount": *resp.RefundAmount,
		}
	}

	c.JSON(200, result)
}

// BatchRefund 批量退票接口
// HTTP POST /api/refund/batch
// 系统批量处理退票
// 请求体:
//   - order_ids: 订单ID列表
//   - operator_id: 操作人ID
//   - refund_type: 退票类型（full/partial）
//   - reason: 批量退票原因
//   - refund_ratio: 退款比例（部分退款时使用，可选）
//
// 返回:
//   - 200: 处理完成
//   - 400: 参数错误
//   - 500: 服务器错误
func BatchRefund(ctx context.Context, c *app.RequestContext) {
	// 定义请求体结构
	var reqBody struct {
		OrderIDs    []int64  `json:"order_ids"`
		OperatorID  int64    `json:"operator_id"`
		RefundType  string   `json:"refund_type"`
		Reason      string   `json:"reason"`
		RefundRatio *float64 `json:"refund_ratio,omitempty"`
	}

	// 解析JSON请求体
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 构造RPC请求
	req := refund.NewBatchRefundReq()
	req.OrderIds = reqBody.OrderIDs
	req.OperatorId = reqBody.OperatorID
	req.RefundType = reqBody.RefundType
	req.Reason = reqBody.Reason

	if reqBody.RefundRatio != nil {
		req.RefundRatio = reqBody.RefundRatio
	}

	// 调用RPC服务批量退票
	resp, err := cli.BatchRefund(context.Background(), req, callopt.WithRPCTimeout(10*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "批量退票失败",
			"error":   err.Error(),
		})
		return
	}

	// 转换结果列表
	results := make([]map[string]interface{}, 0, len(resp.Results))
	for _, item := range resp.Results {
		resultMap := map[string]interface{}{
			"order_id": item.OrderId,
			"success":  item.Success,
			"message":  item.Message,
		}
		if item.RefundAmount != nil {
			resultMap["refund_amount"] = *item.RefundAmount
		}
		results = append(results, resultMap)
	}

	// 返回成功响应
	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data": map[string]interface{}{
			"results":       results,
			"success_count": resp.SuccessCount,
			"fail_count":    resp.FailCount,
		},
	})
}

// RefundStatusQuery 退票状态查询接口
// HTTP GET /api/refund/status/:id?user_id=xxx
// 查询退票申请状态
// 参数:
//   - id: 路径参数，退票申请ID
//   - user_id: 查询参数，用户ID
//
// 返回:
//   - 200: 查询成功
//   - 400: 参数错误或查询失败
//   - 500: 服务器错误
func RefundStatusQuery(ctx context.Context, c *app.RequestContext) {
	// 获取路径参数
	refundIDStr := c.Param("id")
	userIDStr := c.Query("user_id")

	// 参数验证
	if refundIDStr == "" || userIDStr == "" {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "退票ID和用户ID为必填项",
		})
		return
	}

	// 转换参数类型
	refundID, err := strconv.ParseInt(refundIDStr, 10, 64)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "无效的退票ID",
		})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	// 构造RPC请求
	req := refund.NewRefundStatusQueryReq()
	req.RefundId = refundID
	req.UserId = userID

	// 调用RPC服务查询退票状态
	resp, err := cli.RefundStatusQuery(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "查询退票状态失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查查询结果
	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	// 构建退票信息响应
	info := resp.RefundInfo
	refundInfo := map[string]interface{}{
		"refund_id":  info.RefundId,
		"order_id":   info.OrderId,
		"user_id":    info.UserId,
		"user_type":  info.UserType,
		"reason":     info.Reason,
		"status":     info.Status,
		"apply_time": info.ApplyTime,
	}

	if info.AuditReason != nil {
		refundInfo["audit_reason"] = *info.AuditReason
	}
	if info.RefundAmount != nil {
		refundInfo["refund_amount"] = *info.RefundAmount
	}
	if info.AuditTime != nil {
		refundInfo["audit_time"] = *info.AuditTime
	}
	if info.CompletedTime != nil {
		refundInfo["completed_time"] = *info.CompletedTime
	}

	// 返回成功响应
	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data":    refundInfo,
	})
}

// RefundRecordQuery 退票记录查询接口
// HTTP POST /api/refund/records
// 查询退票记录列表
// 请求体:
//   - user_id: 用户ID（可选）
//   - user_type: 用户类型（可选）
//   - status: 状态筛选（可选）
//   - start_date: 开始日期（可选）
//   - end_date: 结束日期（可选）
//   - page: 页码（可选，默认1）
//   - page_size: 每页数量（可选，默认10）
//
// 返回:
//   - 200: 查询成功
//   - 400: 参数错误或查询失败
//   - 500: 服务器错误
func RefundRecordQuery(ctx context.Context, c *app.RequestContext) {
	// 定义请求体结构
	var reqBody struct {
		UserID    *int64  `json:"user_id,omitempty"`
		UserType  *string `json:"user_type,omitempty"`
		Status    *string `json:"status,omitempty"`
		StartDate *string `json:"start_date,omitempty"`
		EndDate   *string `json:"end_date,omitempty"`
		Page      *int32  `json:"page,omitempty"`
		PageSize  *int32  `json:"page_size,omitempty"`
	}

	// 解析JSON请求体
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 构造RPC请求
	req := refund.NewRefundRecordQueryReq()
	req.UserId = reqBody.UserID
	req.UserType = reqBody.UserType
	req.Status = reqBody.Status
	req.StartDate = reqBody.StartDate
	req.EndDate = reqBody.EndDate
	req.Page = reqBody.Page
	req.PageSize = reqBody.PageSize

	// 调用RPC服务查询退票记录
	resp, err := cli.RefundRecordQuery(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "查询退票记录失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查查询结果
	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	// 转换记录列表
	records := make([]map[string]interface{}, 0, len(resp.Records))
	for _, record := range resp.Records {
		recordMap := map[string]interface{}{
			"refund_id":  record.RefundId,
			"order_id":   record.OrderId,
			"user_id":    record.UserId,
			"user_type":  record.UserType,
			"reason":     record.Reason,
			"status":     record.Status,
			"apply_time": record.ApplyTime,
		}

		if record.AuditReason != nil {
			recordMap["audit_reason"] = *record.AuditReason
		}
		if record.RefundAmount != nil {
			recordMap["refund_amount"] = *record.RefundAmount
		}
		if record.AuditTime != nil {
			recordMap["audit_time"] = *record.AuditTime
		}
		if record.CompletedTime != nil {
			recordMap["completed_time"] = *record.CompletedTime
		}

		records = append(records, recordMap)
	}

	// 返回成功响应
	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data": map[string]interface{}{
			"records":   records,
			"total":     resp.Total,
			"page":      resp.Page,
			"page_size": resp.PageSize,
		},
	})
}
