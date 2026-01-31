package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"group/handler/model"

	"gorm.io/gorm"
)

// OrderAuditDAO 订单审计DAO
type OrderAuditDAO struct {
	db *gorm.DB
}

// NewOrderAuditDAO 创建订单审计DAO
func NewOrderAuditDAO(db *gorm.DB) *OrderAuditDAO {
	return &OrderAuditDAO{db: db}
}

// ==================== 审计日志记录 ====================

// CreateAuditLog 创建审计日志
func (dao *OrderAuditDAO) CreateAuditLog(ctx context.Context, log *model.OrderAuditLog) error {
	return dao.db.WithContext(ctx).Create(log).Error
}

// QueryAuditLogs 查询审计日志
func (dao *OrderAuditDAO) QueryAuditLogs(ctx context.Context, orderNo, action, operatorType, operatorIP string, operatorID *int64, startTime, endTime *time.Time, page, pageSize int) ([]*model.OrderAuditLog, int64, error) {
	var logs []*model.OrderAuditLog
	var total int64

	query := dao.db.WithContext(ctx).Model(&model.OrderAuditLog{})

	// 添加过滤条件
	if orderNo != "" {
		query = query.Where("order_no = ?", orderNo)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if operatorType != "" {
		query = query.Where("operator_type = ?", operatorType)
	}
	if operatorIP != "" {
		query = query.Where("operator_ip = ?", operatorIP)
	}
	if operatorID != nil {
		query = query.Where("operator_id = ?", *operatorID)
	}
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAuditLogByID 根据ID获取审计日志
func (dao *OrderAuditDAO) GetAuditLogByID(ctx context.Context, logID int64) (*model.OrderAuditLog, error) {
	var log model.OrderAuditLog
	if err := dao.db.WithContext(ctx).Where("id = ?", logID).First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// ==================== 异常行为检测 ====================

// CreateAnomaly 创建异常记录
func (dao *OrderAuditDAO) CreateAnomaly(ctx context.Context, anomaly *model.OrderAnomaly) error {
	return dao.db.WithContext(ctx).Create(anomaly).Error
}

// QueryAnomalies 查询异常记录
func (dao *OrderAuditDAO) QueryAnomalies(ctx context.Context, userID *int64, userType, detectionType string, startTime, endTime *time.Time, page, pageSize int) ([]*model.OrderAnomaly, int64, error) {
	var anomalies []*model.OrderAnomaly
	var total int64

	query := dao.db.WithContext(ctx).Model(&model.OrderAnomaly{})

	// 添加过滤条件
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if userType != "" {
		query = query.Where("user_type = ?", userType)
	}
	if detectionType != "" {
		query = query.Where("anomaly_type = ?", detectionType)
	}
	if startTime != nil {
		query = query.Where("detected_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("detected_at <= ?", *endTime)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("detected_at DESC").Offset(offset).Limit(pageSize).Find(&anomalies).Error; err != nil {
		return nil, 0, err
	}

	return anomalies, total, nil
}

// UpdateAnomalyStatus 更新异常状态
func (dao *OrderAuditDAO) UpdateAnomalyStatus(ctx context.Context, anomalyID, handlerID int64, handlerName, status, handleResult string) error {
	now := time.Now()
	return dao.db.WithContext(ctx).Model(&model.OrderAnomaly{}).
		Where("id = ?", anomalyID).
		Updates(map[string]interface{}{
			"status":        status,
			"handler_id":    handlerID,
			"handler_name":  handlerName,
			"handle_result": handleResult,
			"handled_at":    now,
		}).Error
}

// GetAnomalyByID 根据ID获取异常记录
func (dao *OrderAuditDAO) GetAnomalyByID(ctx context.Context, anomalyID int64) (*model.OrderAnomaly, error) {
	var anomaly model.OrderAnomaly
	if err := dao.db.WithContext(ctx).Where("id = ?", anomalyID).First(&anomaly).Error; err != nil {
		return nil, err
	}
	return &anomaly, nil
}

// ==================== 异常行为自动检测 ====================

// DetectFrequentCancellation 检测频繁取消订单
func (dao *OrderAuditDAO) DetectFrequentCancellation(ctx context.Context, userID int64, userType string, timeWindow time.Duration, threshold int) (bool, int, []string, error) {
	var count int64
	var orders []string

	startTime := time.Now().Add(-timeWindow)

	// 查询时间窗口内的取消订单数
	query := dao.db.WithContext(ctx).Model(&model.OrderAuditLog{}).
		Where("operator_id = ? AND operator_type = ? AND action = ? AND created_at >= ?",
			userID, userType, "cancel", startTime)

	if err := query.Count(&count).Error; err != nil {
		return false, 0, nil, err
	}

	// 如果超过阈值，获取相关订单号
	if int(count) >= threshold {
		var logs []model.OrderAuditLog
		if err := query.Select("order_no").Find(&logs).Error; err != nil {
			return false, 0, nil, err
		}
		for _, log := range logs {
			orders = append(orders, log.OrderNo)
		}
		return true, int(count), orders, nil
	}

	return false, int(count), nil, nil
}

// DetectAbnormalPayment 检测异常支付行为
func (dao *OrderAuditDAO) DetectAbnormalPayment(ctx context.Context, userID int64, timeWindow time.Duration, amountThreshold float64) (bool, int, []string, error) {
	var orders []model.OrderManage
	startTime := time.Now().Add(-timeWindow)

	// 查询时间窗口内的高额支付订单
	if err := dao.db.WithContext(ctx).
		Where("passenger_id = ? AND payment_status = 'paid' AND price >= ? AND payment_time >= ?",
			userID, amountThreshold, startTime).
		Find(&orders).Error; err != nil {
		return false, 0, nil, err
	}

	if len(orders) > 0 {
		var orderNos []string
		for _, order := range orders {
			orderNos = append(orderNos, order.OrderNo)
		}
		return true, len(orders), orderNos, nil
	}

	return false, 0, nil, nil
}

// DetectIPAnomaly 检测IP异常（同一用户短时间内多个不同IP）
func (dao *OrderAuditDAO) DetectIPAnomaly(ctx context.Context, userID int64, timeWindow time.Duration, ipCountThreshold int) (bool, []string, error) {
	var ips []string
	startTime := time.Now().Add(-timeWindow)

	// 查询时间窗口内的不同IP
	if err := dao.db.WithContext(ctx).Model(&model.OrderAuditLog{}).
		Where("operator_id = ? AND created_at >= ?", userID, startTime).
		Distinct("operator_ip").
		Pluck("operator_ip", &ips).Error; err != nil {
		return false, nil, err
	}

	if len(ips) >= ipCountThreshold {
		return true, ips, nil
	}

	return false, nil, nil
}

// ==================== 监控统计 ====================

// GetOrderStats 获取订单统计
func (dao *OrderAuditDAO) GetOrderStats(ctx context.Context, timeRange string) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 总订单数
	var totalOrders int64
	if err := dao.db.WithContext(ctx).Model(&model.OrderManage{}).Count(&totalOrders).Error; err != nil {
		return nil, err
	}
	stats["total_orders"] = totalOrders

	// 今日订单数
	today := time.Now().Truncate(24 * time.Hour)
	var todayOrders int64
	if err := dao.db.WithContext(ctx).Model(&model.OrderManage{}).
		Where("created_at >= ?", today).Count(&todayOrders).Error; err != nil {
		return nil, err
	}
	stats["today_orders"] = todayOrders

	// 各状态订单数
	statuses := []string{"pending", "accepted", "completed", "cancelled"}
	for _, status := range statuses {
		var count int64
		if err := dao.db.WithContext(ctx).Model(&model.OrderManage{}).
			Where("status = ?", status).Count(&count).Error; err != nil {
			return nil, err
		}
		stats[status+"_orders"] = count
	}

	return stats, nil
}

// GetAuditLogStats 获取审计日志统计
func (dao *OrderAuditDAO) GetAuditLogStats(ctx context.Context, timeRange string) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 总日志数
	var totalLogs int64
	if err := dao.db.WithContext(ctx).Model(&model.OrderAuditLog{}).Count(&totalLogs).Error; err != nil {
		return nil, err
	}
	stats["total_logs"] = totalLogs

	// 今日日志数
	today := time.Now().Truncate(24 * time.Hour)
	var todayLogs int64
	if err := dao.db.WithContext(ctx).Model(&model.OrderAuditLog{}).
		Where("created_at >= ?", today).Count(&todayLogs).Error; err != nil {
		return nil, err
	}
	stats["today_logs"] = todayLogs

	return stats, nil
}

// GetAnomalyStats 获取异常统计
func (dao *OrderAuditDAO) GetAnomalyStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 总异常数
	var totalAnomalies int64
	if err := dao.db.WithContext(ctx).Model(&model.OrderAnomaly{}).Count(&totalAnomalies).Error; err != nil {
		return nil, err
	}
	stats["total_anomalies"] = totalAnomalies

	// 待处理异常数
	var pendingAnomalies int64
	if err := dao.db.WithContext(ctx).Model(&model.OrderAnomaly{}).
		Where("status = ?", "pending").Count(&pendingAnomalies).Error; err != nil {
		return nil, err
	}
	stats["pending_anomalies"] = pendingAnomalies

	// 高风险异常数
	var highRiskAnomalies int64
	if err := dao.db.WithContext(ctx).Model(&model.OrderAnomaly{}).
		Where("risk_level = ?", "high").Count(&highRiskAnomalies).Error; err != nil {
		return nil, err
	}
	stats["high_risk_anomalies"] = highRiskAnomalies

	return stats, nil
}

// GetActionStats 获取操作统计
func (dao *OrderAuditDAO) GetActionStats(ctx context.Context, timeRange string) (map[string]int64, error) {
	var results []struct {
		Action string
		Count  int64
	}

	// 根据时间范围计算开始时间
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

	if err := dao.db.WithContext(ctx).Model(&model.OrderAuditLog{}).
		Select("action, COUNT(*) as count").
		Where("created_at >= ?", startTime).
		Group("action").
		Find(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, result := range results {
		stats[result.Action] = result.Count
	}

	return stats, nil
}

// GetAnomalyTypeStats 获取异常类型统计
func (dao *OrderAuditDAO) GetAnomalyTypeStats(ctx context.Context) (map[string]int64, error) {
	var results []struct {
		AnomalyType string
		Count       int64
	}

	if err := dao.db.WithContext(ctx).Model(&model.OrderAnomaly{}).
		Select("anomaly_type, COUNT(*) as count").
		Group("anomaly_type").
		Find(&results).Error; err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	for _, result := range results {
		stats[result.AnomalyType] = result.Count
	}

	return stats, nil
}

// ==================== 日志导出 ====================

// CreateExportRecord 创建导出记录
func (dao *OrderAuditDAO) CreateExportRecord(ctx context.Context, record *model.AuditLogExport) error {
	return dao.db.WithContext(ctx).Create(record).Error
}

// UpdateExportRecord 更新导出记录
func (dao *OrderAuditDAO) UpdateExportRecord(ctx context.Context, exportID int64, status, filePath, downloadURL string, recordCount int64) error {
	return dao.db.WithContext(ctx).Model(&model.AuditLogExport{}).
		Where("id = ?", exportID).
		Updates(map[string]interface{}{
			"status":       status,
			"file_path":    filePath,
			"download_url": downloadURL,
			"record_count": recordCount,
		}).Error
}

// GetExportRecordByID 根据ID获取导出记录
func (dao *OrderAuditDAO) GetExportRecordByID(ctx context.Context, exportID int64) (*model.AuditLogExport, error) {
	var record model.AuditLogExport
	if err := dao.db.WithContext(ctx).Where("id = ?", exportID).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// ExportAuditLogsToJSON 导出审计日志为JSON
func (dao *OrderAuditDAO) ExportAuditLogsToJSON(ctx context.Context, orderNo, action string, operatorID *int64, startTime, endTime *time.Time) ([]byte, int64, error) {
	logs, total, err := dao.QueryAuditLogs(ctx, orderNo, action, "", "", operatorID, startTime, endTime, 1, 10000)
	if err != nil {
		return nil, 0, err
	}

	data, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// ExportAuditLogsToCSV 导出审计日志为CSV
func (dao *OrderAuditDAO) ExportAuditLogsToCSV(ctx context.Context, orderNo, action string, operatorID *int64, startTime, endTime *time.Time) (string, int64, error) {
	logs, total, err := dao.QueryAuditLogs(ctx, orderNo, action, "", "", operatorID, startTime, endTime, 1, 10000)
	if err != nil {
		return "", 0, err
	}

	// 构建CSV内容
	csv := "ID,订单号,操作动作,操作人ID,操作人类型,操作人姓名,操作人IP,操作原因,创建时间\n"
	for _, log := range logs {
		csv += fmt.Sprintf("%d,%s,%s,%d,%s,%s,%s,%s,%s\n",
			log.ID, log.OrderNo, log.Action, log.OperatorID, log.OperatorType,
			log.OperatorName, log.OperatorIP, log.Reason, log.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	return csv, total, nil
}
