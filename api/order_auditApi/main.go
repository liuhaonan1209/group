package main

import (
	"context"
	"log"
	"time"

	"group/kitex_gen/car/order_audit"
	"group/kitex_gen/car/order_audit/orderauditservice"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

var (
	cli orderauditservice.Client
)

func main() {
	// 创建 RPC 客户端
	c, err := orderauditservice.NewClient("car.order_audit", client.WithHostPorts("127.0.0.1:8893"))
	if err != nil {
		log.Fatal(err)
	}
	cli = c

	// 创建 HTTP 服务器
	hz := server.New(server.WithHostPorts("127.0.0.1:55555"))

	// 注册路由
	// 审计日志相关
	hz.POST("/api/audit/log/create", CreateAuditLog)
	hz.POST("/api/audit/log/query", QueryAuditLog)
	hz.POST("/api/audit/log/export", ExportAuditLog)

	// 异常检测相关
	hz.POST("/api/audit/anomaly/detect", DetectAnomaly)
	hz.POST("/api/audit/anomaly/handle", HandleAnomaly)

	// 监控统计
	hz.GET("/api/audit/monitor/stats", GetMonitorStats)

	log.Println("订单审计API服务启动在 localhost:5555")
	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

// CreateAuditLog 创建审计日志
func CreateAuditLog(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		OrderNo      string  `json:"order_no"`
		Action       string  `json:"action"`
		OperatorID   int64   `json:"operator_id"`
		OperatorType string  `json:"operator_type"`
		OperatorName string  `json:"operator_name"`
		OperatorIP   string  `json:"operator_ip"`
		BeforeData   *string `json:"before_data,omitempty"`
		AfterData    *string `json:"after_data,omitempty"`
		Reason       *string `json:"reason,omitempty"`
		Remark       *string `json:"remark,omitempty"`
		UserAgent    string  `json:"user_agent"`
		RequestID    string  `json:"request_id"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	req := order_audit.NewAuditLogCreateReq()
	req.OrderNo = reqBody.OrderNo
	req.Action = reqBody.Action
	req.OperatorId = reqBody.OperatorID
	req.OperatorType = reqBody.OperatorType
	req.OperatorName = reqBody.OperatorName
	req.OperatorIp = reqBody.OperatorIP
	req.BeforeData = reqBody.BeforeData
	req.AfterData = reqBody.AfterData
	req.Reason = reqBody.Reason
	req.Remark = reqBody.Remark
	req.UserAgent = reqBody.UserAgent
	req.RequestId = reqBody.RequestID

	resp, err := cli.AuditLogCreate(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "创建审计日志失败",
			"error":   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data": map[string]interface{}{
			"log_id": resp.LogId,
		},
	})
}

// QueryAuditLog 查询审计日志
func QueryAuditLog(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		OrderNo      *string `json:"order_no,omitempty"`
		Action       *string `json:"action,omitempty"`
		OperatorID   *int64  `json:"operator_id,omitempty"`
		OperatorType *string `json:"operator_type,omitempty"`
		StartTime    *string `json:"start_time,omitempty"`
		EndTime      *string `json:"end_time,omitempty"`
		OperatorIP   *string `json:"operator_ip,omitempty"`
		Page         int32   `json:"page"`
		PageSize     int32   `json:"page_size"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 设置默认分页参数
	if reqBody.Page <= 0 {
		reqBody.Page = 1
	}
	if reqBody.PageSize <= 0 {
		reqBody.PageSize = 20
	}

	req := order_audit.NewAuditLogQueryReq()
	req.OrderNo = reqBody.OrderNo
	req.Action = reqBody.Action
	req.OperatorId = reqBody.OperatorID
	req.OperatorType = reqBody.OperatorType
	req.StartTime = reqBody.StartTime
	req.EndTime = reqBody.EndTime
	req.OperatorIp = reqBody.OperatorIP
	req.Page = reqBody.Page
	req.PageSize = reqBody.PageSize

	resp, err := cli.AuditLogQuery(context.Background(), req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "查询审计日志失败",
			"error":   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":     200,
		"message":  resp.Message,
		"data":     resp.Logs,
		"total":    resp.Total,
		"page":     resp.Page,
		"pageSize": resp.PageSize,
	})
}

// DetectAnomaly 检测异常行为
func DetectAnomaly(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		UserID        *int64  `json:"user_id,omitempty"`
		UserType      *string `json:"user_type,omitempty"`
		DetectionType *string `json:"detection_type,omitempty"`
		StartTime     *string `json:"start_time,omitempty"`
		EndTime       *string `json:"end_time,omitempty"`
		Page          int32   `json:"page"`
		PageSize      int32   `json:"page_size"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 设置默认分页参数
	if reqBody.Page <= 0 {
		reqBody.Page = 1
	}
	if reqBody.PageSize <= 0 {
		reqBody.PageSize = 20
	}

	req := order_audit.NewAnomalyDetectionReq()
	req.UserId = reqBody.UserID
	req.UserType = reqBody.UserType
	req.DetectionType = reqBody.DetectionType
	req.StartTime = reqBody.StartTime
	req.EndTime = reqBody.EndTime
	req.Page = reqBody.Page
	req.PageSize = reqBody.PageSize

	resp, err := cli.AnomalyDetection(context.Background(), req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "检测异常行为失败",
			"error":   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":     200,
		"message":  resp.Message,
		"data":     resp.Anomalies,
		"total":    resp.Total,
		"page":     resp.Page,
		"pageSize": resp.PageSize,
	})
}

// HandleAnomaly 处理异常行为
func HandleAnomaly(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		AnomalyID    int64   `json:"anomaly_id"`
		HandlerID    int64   `json:"handler_id"`
		HandlerName  string  `json:"handler_name"`
		HandleAction string  `json:"handle_action"`
		HandleResult string  `json:"handle_result"`
		Remark       *string `json:"remark,omitempty"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	req := order_audit.NewAnomalyHandleReq()
	req.AnomalyId = reqBody.AnomalyID
	req.HandlerId = reqBody.HandlerID
	req.HandlerName = reqBody.HandlerName
	req.HandleAction = reqBody.HandleAction
	req.HandleResult_ = reqBody.HandleResult
	req.Remark = reqBody.Remark

	resp, err := cli.AnomalyHandle(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "处理异常失败",
			"error":   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
	})
}

// GetMonitorStats 获取监控统计
func GetMonitorStats(ctx context.Context, c *app.RequestContext) {
	timeRange := c.Query("time_range")
	if timeRange == "" {
		timeRange = "24h"
	}

	req := order_audit.NewMonitorStatsReq()
	req.TimeRange = &timeRange

	resp, err := cli.MonitorStats(context.Background(), req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取监控统计失败",
			"error":   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data": map[string]interface{}{
			"total_orders":        resp.TotalOrders,
			"today_orders":        resp.TodayOrders,
			"pending_orders":      resp.PendingOrders,
			"completed_orders":    resp.CompletedOrders,
			"cancelled_orders":    resp.CancelledOrders,
			"total_audit_logs":    resp.TotalAuditLogs,
			"today_audit_logs":    resp.TodayAuditLogs,
			"total_anomalies":     resp.TotalAnomalies,
			"pending_anomalies":   resp.PendingAnomalies,
			"high_risk_anomalies": resp.HighRiskAnomalies,
			"avg_response_time":   resp.AvgResponseTime,
			"action_stats":        resp.ActionStats,
			"anomaly_stats":       resp.AnomalyStats,
		},
	})
}

// ExportAuditLog 导出审计日志
func ExportAuditLog(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		OrderNo      *string `json:"order_no,omitempty"`
		Action       *string `json:"action,omitempty"`
		OperatorID   *int64  `json:"operator_id,omitempty"`
		StartTime    *string `json:"start_time,omitempty"`
		EndTime      *string `json:"end_time,omitempty"`
		ExportFormat string  `json:"export_format"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 默认导出格式为CSV
	if reqBody.ExportFormat == "" {
		reqBody.ExportFormat = "csv"
	}

	req := order_audit.NewAuditLogExportReq()
	req.OrderNo = reqBody.OrderNo
	req.Action = reqBody.Action
	req.OperatorId = reqBody.OperatorID
	req.StartTime = reqBody.StartTime
	req.EndTime = reqBody.EndTime
	req.ExportFormat = reqBody.ExportFormat

	resp, err := cli.AuditLogExport(context.Background(), req, callopt.WithRPCTimeout(30*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "导出审计日志失败",
			"error":   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data": map[string]interface{}{
			"download_url": resp.DownloadUrl,
			"file_name":    resp.FileName,
			"record_count": resp.RecordCount,
		},
	})
}
