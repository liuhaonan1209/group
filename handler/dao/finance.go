package dao

import (
	"fmt"
	"group/global"
	"group/handler/model"
	"time"
)

// ========== 收支对账表 ==========

// GetBalanceSheetList 获取收支对账列表
func GetBalanceSheetList(startDate, endDate string, status, page, size int) (list []model.BalanceSheet, total int64, err error) {
	db := global.DB.Model(&model.BalanceSheet{})

	if startDate != "" {
		db = db.Where("settle_date >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("settle_date <= ?", endDate)
	}
	if status > 0 {
		db = db.Where("status = ?", status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Order("settle_date DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return
}

// GetBalanceSheetSummary 获取收支汇总
func GetBalanceSheetSummary(startDate, endDate string) (totalIncome, totalRefund, totalSettle float64, err error) {
	db := global.DB.Model(&model.BalanceSheet{})

	if startDate != "" {
		db = db.Where("settle_date >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("settle_date <= ?", endDate)
	}

	var result struct {
		TotalIncome float64
		TotalRefund float64
		TotalSettle float64
	}

	err = db.Select("COALESCE(SUM(order_income), 0) as total_income, COALESCE(SUM(order_refund), 0) as total_refund, COALESCE(SUM(order_settle), 0) as total_settle").Scan(&result).Error
	return result.TotalIncome, result.TotalRefund, result.TotalSettle, err
}

// ========== 收入对账表 ==========

// GetIncomeSheetList 获取收入对账列表
func GetIncomeSheetList(startDate, endDate string, status, page, size int) (list []model.IncomeSheet, total int64, err error) {
	db := global.DB.Model(&model.IncomeSheet{})

	if startDate != "" {
		db = db.Where("settle_date >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("settle_date <= ?", endDate)
	}
	if status > 0 {
		db = db.Where("status = ?", status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Order("settle_date DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return
}

// ========== 线路结算表 ==========

// GetRouteSettleList 获取线路结算列表
func GetRouteSettleList(startDate, endDate string, routeId int64, fleet string, page, size int) (list []model.RouteSettle, total int64, err error) {
	db := global.DB.Model(&model.RouteSettle{})

	if startDate != "" {
		db = db.Where("settle_date >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("settle_date <= ?", endDate)
	}
	if routeId > 0 {
		db = db.Where("route_id = ?", routeId)
	}
	if fleet != "" {
		db = db.Where("fleet = ?", fleet)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Order("settle_date DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return
}

// ========== 班次结算表 ==========

// GetScheduleSettleList 获取班次结算列表
func GetScheduleSettleList(startDate, endDate string, routeId, scheduleId int64, page, size int) (list []model.ScheduleSettle, total int64, err error) {
	db := global.DB.Model(&model.ScheduleSettle{})

	if startDate != "" {
		db = db.Where("settle_date >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("settle_date <= ?", endDate)
	}
	if routeId > 0 {
		db = db.Where("route_id = ?", routeId)
	}
	if scheduleId > 0 {
		db = db.Where("schedule_id = ?", scheduleId)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Order("settle_date DESC, departure_time ASC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return
}

// ========== 站点结算表 ==========

// GetStationSettleList 获取站点结算列表
func GetStationSettleList(startDate, endDate string, stationId int64, page, size int) (list []model.StationSettle, total int64, err error) {
	db := global.DB.Model(&model.StationSettle{})

	if startDate != "" {
		db = db.Where("settle_date >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("settle_date <= ?", endDate)
	}
	if stationId > 0 {
		db = db.Where("station_id = ?", stationId)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Order("settle_date DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return
}

// ========== 司机结算表 ==========

// GetDriverSettleList 获取司机结算列表
func GetDriverSettleList(startDate, endDate string, driverId int64, driverName string, page, size int) (list []model.DriverSettle, total int64, err error) {
	db := global.DB.Model(&model.DriverSettle{})

	if startDate != "" {
		db = db.Where("settle_date >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("settle_date <= ?", endDate)
	}
	if driverId > 0 {
		db = db.Where("driver_id = ?", driverId)
	}
	if driverName != "" {
		db = db.Where("driver_name LIKE ?", "%"+driverName+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Order("settle_date DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return
}

// ========== 交易明细 ==========

// GetTransactionList 获取交易明细列表
func GetTransactionList(startDate, endDate, orderNo string, paymentMethod, transType, page, size int) (list []model.Transaction, total int64, err error) {
	db := global.DB.Model(&model.Transaction{})

	if startDate != "" {
		db = db.Where("trans_date >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("trans_date <= ?", endDate)
	}
	if orderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+orderNo+"%")
	}
	if paymentMethod > 0 {
		db = db.Where("payment_method = ?", paymentMethod)
	}
	if transType > 0 {
		db = db.Where("trans_type = ?", transType)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Order("trans_time DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return
}

// CreateTransaction 创建交易记录
func CreateTransaction(trans *model.Transaction) error {
	return global.DB.Create(trans).Error
}

// ========== 异常交易 ==========

// GetAbnormalTransList 获取异常交易列表
func GetAbnormalTransList(startDate, endDate string, abnormalType, status, page, size int) (list []model.AbnormalTransaction, total int64, err error) {
	db := global.DB.Model(&model.AbnormalTransaction{})

	if startDate != "" {
		db = db.Where("trans_date >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("trans_date <= ?", endDate)
	}
	if abnormalType > 0 {
		db = db.Where("abnormal_type = ?", abnormalType)
	}
	if status > 0 {
		db = db.Where("status = ?", status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return
}

// HandleAbnormalTrans 处理异常交易
func HandleAbnormalTrans(id int64, handler, remark string, handleType int) error {
	now := time.Now()
	return global.DB.Model(&model.AbnormalTransaction{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        model.SettleStatusDone,
		"handler":       handler,
		"handle_time":   &now,
		"handle_remark": remark,
		"handle_type":   handleType,
	}).Error
}

// CreateAbnormalTrans 创建异常交易记录
func CreateAbnormalTrans(trans *model.AbnormalTransaction) error {
	return global.DB.Create(trans).Error
}

// DetectAbnormalTransactions 检测异常交易
func DetectAbnormalTransactions() error {
	// 检测重复扣款
	var duplicates []struct {
		OrderNo string
		Count   int
	}
	global.DB.Model(&model.Transaction{}).
		Select("order_no, COUNT(*) as count").
		Where("trans_type = ?", model.TransTypeIncome).
		Group("order_no").
		Having("COUNT(*) > 1").
		Find(&duplicates)

	for _, d := range duplicates {
		// 检查是否已记录
		var existing model.AbnormalTransaction
		err := global.DB.Where("order_no = ? AND abnormal_type = ? AND status = ?",
			d.OrderNo, model.AbnormalTypeDuplicate, model.SettleStatusPending).First(&existing).Error
		if err != nil {
			// 创建异常记录
			abnormal := &model.AbnormalTransaction{
				TransDate:    time.Now(),
				OrderNo:      d.OrderNo,
				AbnormalType: model.AbnormalTypeDuplicate,
				Status:       model.SettleStatusPending,
			}
			CreateAbnormalTrans(abnormal)
		}
	}

	return nil
}

// ========== 账单生成 ==========

// GenerateDailyBill 生成日账单
func GenerateDailyBill(date time.Time) (*model.Bill, error) {
	startDate := date.Format("2006-01-02")
	endDate := startDate

	totalIncome, totalRefund, totalSettle, err := GetBalanceSheetSummary(startDate, endDate)
	if err != nil {
		return nil, err
	}

	bill := &model.Bill{
		BillNo:      fmt.Sprintf("D%s%d", date.Format("20060102"), time.Now().UnixNano()%10000),
		BillType:    model.BillTypeDaily,
		StartDate:   date,
		EndDate:     date,
		TotalIncome: totalIncome,
		TotalRefund: totalRefund,
		TotalSettle: totalSettle,
		Status:      1,
	}

	err = global.DB.Create(bill).Error
	return bill, err
}

// GenerateWeeklyBill 生成周账单
func GenerateWeeklyBill(startDate, endDate time.Time) (*model.Bill, error) {
	totalIncome, totalRefund, totalSettle, err := GetBalanceSheetSummary(
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}

	bill := &model.Bill{
		BillNo:      fmt.Sprintf("W%s%d", startDate.Format("20060102"), time.Now().UnixNano()%10000),
		BillType:    model.BillTypeWeekly,
		StartDate:   startDate,
		EndDate:     endDate,
		TotalIncome: totalIncome,
		TotalRefund: totalRefund,
		TotalSettle: totalSettle,
		Status:      1,
	}

	err = global.DB.Create(bill).Error
	return bill, err
}

// GenerateMonthlyBill 生成月账单
func GenerateMonthlyBill(year, month int) (*model.Bill, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	totalIncome, totalRefund, totalSettle, err := GetBalanceSheetSummary(
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}

	bill := &model.Bill{
		BillNo:      fmt.Sprintf("M%d%02d%d", year, month, time.Now().UnixNano()%10000),
		BillType:    model.BillTypeMonthly,
		StartDate:   startDate,
		EndDate:     endDate,
		TotalIncome: totalIncome,
		TotalRefund: totalRefund,
		TotalSettle: totalSettle,
		Status:      1,
	}

	err = global.DB.Create(bill).Error
	return bill, err
}

// ========== 数据同步与结算 ==========

// SyncDailyBalance 同步每日收支数据
func SyncDailyBalance(date time.Time) error {
	dateStr := date.Format("2006-01-02")

	// 统计当日订单收入
	var income float64
	global.DB.Model(&model.Transaction{}).
		Where("trans_date = ? AND trans_type = ?", dateStr, model.TransTypeIncome).
		Select("COALESCE(SUM(amount), 0)").Scan(&income)

	// 统计当日退款
	var refund float64
	global.DB.Model(&model.Transaction{}).
		Where("trans_date = ? AND trans_type = ?", dateStr, model.TransTypeRefund).
		Select("COALESCE(SUM(amount), 0)").Scan(&refund)

	// 计算结算金额
	settle := income - refund

	// 检查是否有异常
	status := model.BalanceStatusNormal
	var abnormalCount int64
	global.DB.Model(&model.AbnormalTransaction{}).
		Where("trans_date = ? AND status = ?", dateStr, model.SettleStatusPending).
		Count(&abnormalCount)
	if abnormalCount > 0 {
		status = model.BalanceStatusAbnormal
	}

	// 更新或创建收支记录
	var existing model.BalanceSheet
	err := global.DB.Where("settle_date = ?", dateStr).First(&existing).Error
	if err != nil {
		// 创建新记录
		balance := &model.BalanceSheet{
			SettleDate:  date,
			OrderIncome: income,
			OrderRefund: refund,
			OrderSettle: settle,
			Status:      status,
		}
		return global.DB.Create(balance).Error
	}

	// 更新现有记录
	return global.DB.Model(&existing).Updates(map[string]interface{}{
		"order_income": income,
		"order_refund": refund,
		"order_settle": settle,
		"status":       status,
	}).Error
}
