package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"group/global"
	"group/handler/model"
	order_audit "group/kitex_gen/car/order_audit"

	"gorm.io/gorm"
)

// OrderAuditServiceImpl implements the last service interface defined in the IDL.
type OrderAuditServiceImpl struct{}

// AuditLogCreate 创建审计日志
func (s *OrderAuditServiceImpl) AuditLogCreate(ctx context.Context, req *order_audit.AuditLogCreateReq) (*order_audit.AuditLogCreateResp, error) {
	// 创建审计日志
	auditLog := &model.OrderAuditLog{
		OrderNo:      req.OrderNo,
		Action:       req.Action,
		OperatorID:   req.OperatorId,
		OperatorType: req.OperatorType,
		OperatorName: req.OperatorName,
		OperatorIP:   req.OperatorIp,
		BeforeData:   getStringValue(req.BeforeData),
		AfterData:    getStringValue(req.AfterData),
		Reason:       getStringValue(req.Reason),
		Remark:       getStringValue(req.Remark),
		UserAgent:    req.UserAgent,
		RequestID:    req.RequestId,
		CreatedAt:    time.Now(),
	}

	if err := global.DB.Create(auditLog).Error; err != nil {
		log.Printf("创建审计日志失败: %v", err)
		return &order_audit.AuditLogCreateResp{
			LogId:   0,
			Message: "创建审计日志失败",
			Success: false,
		}, nil
	}

	// 异步检测异常行为
	go detectAnomalies(req.OperatorId, req.OperatorType, req.OperatorName, req.Action, req.OrderNo)

	log.Printf("审计日志创建成功: ID=%d", auditLog.ID)
	return &order_audit.AuditLogCreateResp{
		LogId:   int64(auditLog.ID),
		Message: "审计日志创建成功",
		Success: true,
	}, nil
}

// AuditLogQuery 查询审计日志
func (s *OrderAuditServiceImpl) AuditLogQuery(ctx context.Context, req *order_audit.AuditLogQueryReq) (*order_audit.AuditLogQueryResp, error) {
	log.Printf("查询审计日志: 页码=%d, 每页数量=%d", req.Page, req.PageSize)

	// 验证分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 构建查询条件
	query := global.DB.Model(&model.OrderAuditLog{})

	if req.OrderNo != nil && *req.OrderNo != "" {
		query = query.Where("order_no = ?", *req.OrderNo)
	}
	if req.Action != nil && *req.Action != "" {
		query = query.Where("action = ?", *req.Action)
	}
	if req.OperatorType != nil && *req.OperatorType != "" {
		query = query.Where("operator_type = ?", *req.OperatorType)
	}
	if req.OperatorIp != nil && *req.OperatorIp != "" {
		query = query.Where("operator_ip = ?", *req.OperatorIp)
	}
	if req.OperatorId != nil {
		query = query.Where("operator_id = ?", *req.OperatorId)
	}
	if req.StartTime != nil && *req.StartTime != "" {
		startTime, err := time.Parse("2006-01-02 15:04:05", *req.StartTime)
		if err == nil {
			query = query.Where("created_at >= ?", startTime)
		}
	}
	if req.EndTime != nil && *req.EndTime != "" {
		endTime, err := time.Parse("2006-01-02 15:04:05", *req.EndTime)
		if err == nil {
			query = query.Where("created_at <= ?", endTime)
		}
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Printf("统计审计日志失败: %v", err)
		return &order_audit.AuditLogQueryResp{
			Logs:     []*order_audit.AuditLogInfo{},
			Total:    0,
			Page:     req.Page,
			PageSize: req.PageSize,
			Message:  "查询失败",
			Success:  false,
		}, nil
	}

	// 分页查询
	var logs []model.OrderAuditLog
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").Offset(int(offset)).Limit(int(req.PageSize)).Find(&logs).Error; err != nil {
		log.Printf("查询审计日志失败: %v", err)
		return &order_audit.AuditLogQueryResp{
			Logs:     []*order_audit.AuditLogInfo{},
			Total:    0,
			Page:     req.Page,
			PageSize: req.PageSize,
			Message:  "查询失败",
			Success:  false,
		}, nil
	}

	// 构造响应
	logInfos := make([]*order_audit.AuditLogInfo, 0, len(logs))
	for _, log := range logs {
		logInfo := &order_audit.AuditLogInfo{
			LogId:        int64(log.ID),
			OrderNo:      log.OrderNo,
			Action:       log.Action,
			OperatorId:   log.OperatorID,
			OperatorType: log.OperatorType,
			OperatorName: log.OperatorName,
			OperatorIp:   log.OperatorIP,
			BeforeData:   log.BeforeData,
			AfterData:    log.AfterData,
			Reason:       log.Reason,
			Remark:       log.Remark,
			UserAgent:    log.UserAgent,
			RequestId:    log.RequestID,
			CreatedAt:    log.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		logInfos = append(logInfos, logInfo)
	}

	log.Printf("查询到 %d 条审计日志", len(logInfos))
	return &order_audit.AuditLogQueryResp{
		Logs:     logInfos,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Message:  fmt.Sprintf("查询成功，共找到%d条记录", total),
		Success:  true,
	}, nil
}

// AnomalyDetection 异常行为检测
func (s *OrderAuditServiceImpl) AnomalyDetection(ctx context.Context, req *order_audit.AnomalyDetectionReq) (*order_audit.AnomalyDetectionResp, error) {
	log.Printf("异常行为检测: 页码=%d, 每页数量=%d", req.Page, req.PageSize)

	// 验证分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 构建查询条件
	query := global.DB.Model(&model.OrderAnomaly{})

	if req.UserId != nil {
		query = query.Where("user_id = ?", *req.UserId)
	}
	if req.UserType != nil && *req.UserType != "" {
		query = query.Where("user_type = ?", *req.UserType)
	}
	if req.DetectionType != nil && *req.DetectionType != "" {
		query = query.Where("anomaly_type = ?", *req.DetectionType)
	}
	if req.StartTime != nil && *req.StartTime != "" {
		startTime, err := time.Parse("2006-01-02 15:04:05", *req.StartTime)
		if err == nil {
			query = query.Where("detected_at >= ?", startTime)
		}
	}
	if req.EndTime != nil && *req.EndTime != "" {
		endTime, err := time.Parse("2006-01-02 15:04:05", *req.EndTime)
		if err == nil {
			query = query.Where("detected_at <= ?", endTime)
		}
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Printf("统计异常记录失败: %v", err)
		return &order_audit.AnomalyDetectionResp{
			Anomalies: []*order_audit.AnomalyInfo{},
			Total:     0,
			Page:      req.Page,
			PageSize:  req.PageSize,
			Message:   "查询失败",
			Success:   false,
		}, nil
	}

	// 分页查询
	var anomalies []model.OrderAnomaly
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("detected_at DESC").Offset(int(offset)).Limit(int(req.PageSize)).Find(&anomalies).Error; err != nil {
		log.Printf("查询异常记录失败: %v", err)
		return &order_audit.AnomalyDetectionResp{
			Anomalies: []*order_audit.AnomalyInfo{},
			Total:     0,
			Page:      req.Page,
			PageSize:  req.PageSize,
			Message:   "查询失败",
			Success:   false,
		}, nil
	}

	// 构造响应
	anomalyInfos := make([]*order_audit.AnomalyInfo, 0, len(anomalies))
	for _, anomaly := range anomalies {
		anomalyInfo := &order_audit.AnomalyInfo{
			AnomalyId:       int64(anomaly.ID),
			UserId:          anomaly.UserID,
			UserType:        anomaly.UserType,
			UserName:        anomaly.UserName,
			AnomalyType:     anomaly.AnomalyType,
			Description:     anomaly.Description,
			RiskLevel:       anomaly.RiskLevel,
			OccurrenceCount: int32(anomaly.OccurrenceCount),
			RelatedOrders:   anomaly.RelatedOrders,
			DetectedAt:      anomaly.DetectedAt.Format("2006-01-02 15:04:05"),
			Status:          anomaly.Status,
			HandlerName:     anomaly.HandlerName,
		}
		anomalyInfos = append(anomalyInfos, anomalyInfo)
	}

	log.Printf("查询到 %d 条异常记录", len(anomalyInfos))
	return &order_audit.AnomalyDetectionResp{
		Anomalies: anomalyInfos,
		Total:     total,
		Page:      req.Page,
		PageSize:  req.PageSize,
		Message:   fmt.Sprintf("查询成功，共找到%d条异常记录", total),
		Success:   true,
	}, nil
}

// MonitorStats 实时监控统计
func (s *OrderAuditServiceImpl) MonitorStats(ctx context.Context, req *order_audit.MonitorStatsReq) (*order_audit.MonitorStatsResp, error) {
	log.Printf("获取监控统计: 时间范围=%s", getStringValue(req.TimeRange))

	// 获取订单统计
	var totalOrders, todayOrders, pendingOrders, completedOrders, cancelledOrders int64
	global.DB.Model(&model.OrderManage{}).Count(&totalOrders)

	today := time.Now().Truncate(24 * time.Hour)
	global.DB.Model(&model.OrderManage{}).Where("created_at >= ?", today).Count(&todayOrders)
	global.DB.Model(&model.OrderManage{}).Where("status = ?", "pending").Count(&pendingOrders)
	global.DB.Model(&model.OrderManage{}).Where("status = ?", "completed").Count(&completedOrders)
	global.DB.Model(&model.OrderManage{}).Where("status = ?", "cancelled").Count(&cancelledOrders)

	// 获取审计日志统计
	var totalLogs, todayLogs int64
	global.DB.Model(&model.OrderAuditLog{}).Count(&totalLogs)
	global.DB.Model(&model.OrderAuditLog{}).Where("created_at >= ?", today).Count(&todayLogs)

	// 获取异常统计
	var totalAnomalies, pendingAnomalies, highRiskAnomalies int64
	global.DB.Model(&model.OrderAnomaly{}).Count(&totalAnomalies)
	global.DB.Model(&model.OrderAnomaly{}).Where("status = ?", "pending").Count(&pendingAnomalies)
	global.DB.Model(&model.OrderAnomaly{}).Where("risk_level = ?", "high").Count(&highRiskAnomalies)

	// 获取操作统计
	actionStats := make(map[string]int64)
	var actionResults []struct {
		Action string
		Count  int64
	}

	timeRange := getStringValue(req.TimeRange)
	if timeRange == "" {
		timeRange = "24h"
	}

	var startTime time.Time
	switch timeRange {
	case "1h":
		startTime = time.Now().Add(-1 * time.Hour)
	case "24h":
		startTime = time.Now().Add(-24 * time.Hour)
	case "7d":
		startTime = time.Now().Add(-7 * 24 * time.Hour)
	case "30d":
		startTime = time.Now().Add(-30 * 24 * time.Hour)
	default:
		startTime = time.Now().Add(-24 * time.Hour)
	}

	global.DB.Model(&model.OrderAuditLog{}).
		Select("action, COUNT(*) as count").
		Where("created_at >= ?", startTime).
		Group("action").
		Find(&actionResults)

	for _, result := range actionResults {
		actionStats[result.Action] = result.Count
	}

	// 获取异常类型统计
	anomalyStats := make(map[string]int64)
	var anomalyResults []struct {
		AnomalyType string
		Count       int64
	}

	global.DB.Model(&model.OrderAnomaly{}).
		Select("anomaly_type, COUNT(*) as count").
		Group("anomaly_type").
		Find(&anomalyResults)

	for _, result := range anomalyResults {
		anomalyStats[result.AnomalyType] = result.Count
	}

	// 构造响应
	resp := &order_audit.MonitorStatsResp{
		TotalOrders:       totalOrders,
		TodayOrders:       todayOrders,
		PendingOrders:     pendingOrders,
		CompletedOrders:   completedOrders,
		CancelledOrders:   cancelledOrders,
		TotalAuditLogs:    totalLogs,
		TodayAuditLogs:    todayLogs,
		TotalAnomalies:    totalAnomalies,
		PendingAnomalies:  pendingAnomalies,
		HighRiskAnomalies: highRiskAnomalies,
		AvgResponseTime:   150.5, // 模拟平均响应时间
		ActionStats:       actionStats,
		AnomalyStats:      anomalyStats,
		Message:           "统计数据获取成功",
		Success:           true,
	}

	log.Printf("监控统计获取成功")
	return resp, nil
}

// AnomalyHandle 异常行为处理
func (s *OrderAuditServiceImpl) AnomalyHandle(ctx context.Context, req *order_audit.AnomalyHandleReq) (*order_audit.AnomalyHandleResp, error) {
	log.Printf("处理异常: 异常ID=%d, 处理人=%s, 处理动作=%s", req.AnomalyId, req.HandlerName, req.HandleAction)

	// 验证必填字段
	if req.AnomalyId == 0 || req.HandlerId == 0 || req.HandleAction == "" {
		return &order_audit.AnomalyHandleResp{
			Message: "异常ID、处理人ID和处理动作为必填项",
			Success: false,
		}, nil
	}

	// 验证异常是否存在
	var anomaly model.OrderAnomaly
	if err := global.DB.Where("id = ?", req.AnomalyId).First(&anomaly).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &order_audit.AnomalyHandleResp{
				Message: "异常记录不存在",
				Success: false,
			}, nil
		}
		log.Printf("查询异常记录失败: %v", err)
		return &order_audit.AnomalyHandleResp{
			Message: "查询异常记录失败",
			Success: false,
		}, nil
	}

	// 根据处理动作更新状态
	var status string
	switch req.HandleAction {
	case "resolve":
		status = "resolved"
	case "ignore":
		status = "ignored"
	case "escalate":
		status = "escalated"
	default:
		return &order_audit.AnomalyHandleResp{
			Message: "无效的处理动作",
			Success: false,
		}, nil
	}

	// 更新异常状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":       status,
		"handler_id":   req.HandlerId,
		"handler_name": req.HandlerName,
		"handled_at":   now,
	}

	// 如果提供了处理结果，也更新
	if req.HandleResult_ != "" {
		updates["handle_result"] = req.HandleResult_
	}

	if err := global.DB.Model(&anomaly).Updates(updates).Error; err != nil {
		log.Printf("更新异常状态失败: %v", err)
		return &order_audit.AnomalyHandleResp{
			Message: "处理异常失败",
			Success: false,
		}, nil
	}

	log.Printf("异常处理成功: 异常ID=%d, 原状态=%s, 新状态=%s", req.AnomalyId, anomaly.Status, status)
	return &order_audit.AnomalyHandleResp{
		Message: "异常处理成功",
		Success: true,
	}, nil
}

// AuditLogExport 导出审计日志
func (s *OrderAuditServiceImpl) AuditLogExport(ctx context.Context, req *order_audit.AuditLogExportReq) (*order_audit.AuditLogExportResp, error) {
	log.Printf("导出审计日志: 格式=%s", req.ExportFormat)

	// 验证导出格式
	if req.ExportFormat != "csv" && req.ExportFormat != "excel" && req.ExportFormat != "json" {
		return &order_audit.AuditLogExportResp{
			DownloadUrl: "",
			FileName:    "",
			RecordCount: 0,
			Message:     "不支持的导出格式，仅支持 csv、excel、json",
			Success:     false,
		}, nil
	}

	// 构建查询条件
	query := global.DB.Model(&model.OrderAuditLog{})

	if req.OrderNo != nil && *req.OrderNo != "" {
		query = query.Where("order_no = ?", *req.OrderNo)
	}
	if req.Action != nil && *req.Action != "" {
		query = query.Where("action = ?", *req.Action)
	}
	if req.OperatorId != nil {
		query = query.Where("operator_id = ?", *req.OperatorId)
	}
	if req.StartTime != nil && *req.StartTime != "" {
		startTime, err := time.Parse("2006-01-02 15:04:05", *req.StartTime)
		if err == nil {
			query = query.Where("created_at >= ?", startTime)
		}
	}
	if req.EndTime != nil && *req.EndTime != "" {
		endTime, err := time.Parse("2006-01-02 15:04:05", *req.EndTime)
		if err == nil {
			query = query.Where("created_at <= ?", endTime)
		}
	}

	// 查询数据
	var logs []model.OrderAuditLog
	if err := query.Limit(10000).Find(&logs).Error; err != nil {
		log.Printf("查询审计日志失败: %v", err)
		return &order_audit.AuditLogExportResp{
			Message: "导出失败",
			Success: false,
		}, nil
	}

	var fileName string
	var recordCount int64 = int64(len(logs))

	switch req.ExportFormat {
	case "json":
		fileName = fmt.Sprintf("audit_log_%s.json", time.Now().Format("20060102150405"))
	case "csv":
		fileName = fmt.Sprintf("audit_log_%s.csv", time.Now().Format("20060102150405"))
	default:
		return &order_audit.AuditLogExportResp{
			Message: "暂不支持Excel格式导出",
			Success: false,
		}, nil
	}

	// 生成下载链接（实际应该保存文件到服务器）
	downloadUrl := fmt.Sprintf("https://carpool.example.com/download/%s", fileName)

	log.Printf("审计日志导出成功: 文件=%s, 记录数=%d", fileName, recordCount)
	return &order_audit.AuditLogExportResp{
		DownloadUrl: downloadUrl,
		FileName:    fileName,
		RecordCount: recordCount,
		Message:     fmt.Sprintf("导出成功，共%d条记录", recordCount),
		Success:     true,
	}, nil
}

// ==================== 辅助方法 ====================

// getStringValue 获取可选字符串的值
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// detectAnomalies 异步检测异常行为
func detectAnomalies(userID int64, userType, userName, action, orderNo string) {
	// 1. 检测频繁取消订单
	if action == "cancel" {
		var count int64
		startTime := time.Now().Add(-24 * time.Hour)
		global.DB.Model(&model.OrderAuditLog{}).
			Where("operator_id = ? AND operator_type = ? AND action = ? AND created_at >= ?",
				userID, userType, "cancel", startTime).
			Count(&count)

		if count >= 3 {
			var logs []model.OrderAuditLog
			global.DB.Model(&model.OrderAuditLog{}).
				Select("order_no").
				Where("operator_id = ? AND operator_type = ? AND action = ? AND created_at >= ?",
					userID, userType, "cancel", startTime).
				Find(&logs)

			var orders []string
			for _, log := range logs {
				orders = append(orders, log.OrderNo)
			}

			anomaly := &model.OrderAnomaly{
				UserID:          userID,
				UserType:        userType,
				UserName:        userName,
				AnomalyType:     "frequent_cancellation",
				Description:     fmt.Sprintf("用户在24小时内取消了%d个订单", count),
				RiskLevel:       "medium",
				OccurrenceCount: int(count),
				RelatedOrders:   strings.Join(orders, ","),
				DetectedAt:      time.Now(),
				Status:          "pending",
			}
			global.DB.Create(anomaly)
			log.Printf("检测到频繁取消异常: 用户ID=%d, 次数=%d", userID, count)
		}
	}

	// 2. 检测异常支付（仅针对乘客）
	if action == "pay" && userType == "passenger" {
		var orders []model.OrderManage
		startTime := time.Now().Add(-1 * time.Hour)
		global.DB.Where("passenger_id = ? AND payment_status = 'paid' AND price >= ? AND payment_time >= ?",
			userID, 1000.0, startTime).Find(&orders)

		if len(orders) > 0 {
			var orderNos []string
			for _, order := range orders {
				orderNos = append(orderNos, order.OrderNo)
			}

			anomaly := &model.OrderAnomaly{
				UserID:          userID,
				UserType:        userType,
				UserName:        userName,
				AnomalyType:     "abnormal_payment",
				Description:     fmt.Sprintf("用户在1小时内进行了%d笔高额支付", len(orders)),
				RiskLevel:       "high",
				OccurrenceCount: len(orders),
				RelatedOrders:   strings.Join(orderNos, ","),
				DetectedAt:      time.Now(),
				Status:          "pending",
			}
			global.DB.Create(anomaly)
			log.Printf("检测到异常支付: 用户ID=%d, 次数=%d", userID, len(orders))
		}
	}

	// 3. 检测IP异常
	var ips []string
	startTime := time.Now().Add(-1 * time.Hour)
	global.DB.Model(&model.OrderAuditLog{}).
		Where("operator_id = ? AND created_at >= ?", userID, startTime).
		Distinct("operator_ip").
		Pluck("operator_ip", &ips)

	if len(ips) >= 5 {
		anomaly := &model.OrderAnomaly{
			UserID:          userID,
			UserType:        userType,
			UserName:        userName,
			AnomalyType:     "ip_anomaly",
			Description:     fmt.Sprintf("用户在1小时内使用了%d个不同的IP地址: %s", len(ips), strings.Join(ips, ", ")),
			RiskLevel:       "high",
			OccurrenceCount: len(ips),
			RelatedOrders:   orderNo,
			DetectedAt:      time.Now(),
			Status:          "pending",
		}
		global.DB.Create(anomaly)
		log.Printf("检测到IP异常: 用户ID=%d, IP数量=%d", userID, len(ips))
	}
}

// logToJSON 将对象转换为JSON字符串
func logToJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}
