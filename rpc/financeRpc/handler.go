package main

import (
	"context"
	"time"

	"group/handler/dao"
	"group/handler/model"
	"group/kitex_gen/car/finance"
)

type FinanceServiceImpl struct{}

// ========== 收支对账 ==========
// 财务收支对账的分页查询接口
func (s *FinanceServiceImpl) GetBalanceSheet(ctx context.Context, req *finance.BalanceSheetReq) (resp *finance.BalanceSheetResp, err error) {
	resp = &finance.BalanceSheetResp{}
	//分页参数容错处理
	page, size := int(req.Page), int(req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	//调用 DAO 层，执行带条件的分页对账明细查询
	list, total, err := dao.GetBalanceSheetList(req.StartDate, req.EndDate, int(req.Status), page, size)
	if err != nil {
		resp.Success = false
		resp.Msg = "查询失败: " + err.Error()
		return resp, nil
	}
	//数据格式转换：Model 层 → API 层（财务专属格式处理）
	items := make([]*finance.BalanceSheetItem, len(list))
	for i, item := range list {
		items[i] = &finance.BalanceSheetItem{
			Id:          int64(item.ID),
			SettleDate:  item.SettleDate.Format("2006-01-02"),
			OrderIncome: item.OrderIncome,
			OrderRefund: item.OrderRefund,
			OrderSettle: item.OrderSettle,
			Status:      int32(item.Status),
			StatusText:  model.GetBalanceStatusText(item.Status),
		}
	}

	// 获取汇总数据
	totalIncome, totalRefund, totalSettle, _ := dao.GetBalanceSheetSummary(req.StartDate, req.EndDate)

	resp.Success = true
	resp.Msg = "查询成功"
	resp.List = items
	resp.Total = total
	resp.TotalIncome = totalIncome
	resp.TotalRefund = totalRefund
	resp.TotalSettle = totalSettle
	return
}

// ========== 收入对账 ==========

func (s *FinanceServiceImpl) GetIncomeSheet(ctx context.Context, req *finance.IncomeSheetReq) (resp *finance.IncomeSheetResp, err error) {
	resp = &finance.IncomeSheetResp{}

	page, size := int(req.Page), int(req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	list, total, err := dao.GetIncomeSheetList(req.StartDate, req.EndDate, int(req.Status), page, size)
	if err != nil {
		resp.Success = false
		resp.Msg = "查询失败: " + err.Error()
		return resp, nil
	}

	items := make([]*finance.IncomeSheetItem, len(list))
	for i, item := range list {
		items[i] = &finance.IncomeSheetItem{
			Id:                     int64(item.ID),
			SettleDate:             item.SettleDate.Format("2006-01-02"),
			RouteNo:                item.RouteNo,
			RouteName:              item.RouteName,
			Fleet:                  item.Fleet,
			TicketPrice:            item.TicketPrice,
			DiscountAmount:         item.DiscountAmount,
			CheckedIncome:          item.CheckedIncome,
			CheckedTickets:         int32(item.CheckedTickets),
			UncheckedIncome:        item.UncheckedIncome,
			UncheckedTickets:       int32(item.UncheckedTickets),
			CheckedRefund:          item.CheckedRefund,
			CheckedRefundTickets:   int32(item.CheckedRefundTickets),
			UncheckedRefund:        item.UncheckedRefund,
			UncheckedRefundTickets: int32(item.UncheckedRefundTickets),
			Status:                 int32(item.Status),
		}
	}

	resp.Success = true
	resp.Msg = "查询成功"
	resp.List = items
	resp.Total = total
	return
}

// ========== 线路结算 ==========
// 财务模块的线路结算分页查询接口
func (s *FinanceServiceImpl) GetRouteSettle(ctx context.Context, req *finance.RouteSettleReq) (resp *finance.RouteSettleResp, err error) {
	resp = &finance.RouteSettleResp{}
	// 分页参数容错处理
	page, size := int(req.Page), int(req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	//执行线路维度的分页筛选查询
	list, total, err := dao.GetRouteSettleList(req.StartDate, req.EndDate, req.RouteId, req.Fleet, page, size)
	if err != nil {
		resp.Success = false
		resp.Msg = "查询失败: " + err.Error()
		return resp, nil
	}
	//数据格式转换：Model 层 → API 层
	items := make([]*finance.RouteSettleItem, len(list))
	for i, item := range list {
		items[i] = &finance.RouteSettleItem{
			Id:                     int64(item.ID),
			SettleDate:             item.SettleDate.Format("2006-01-02"),
			RouteNo:                item.RouteNo,
			RouteName:              item.RouteName,
			Fleet:                  item.Fleet,
			TicketPrice:            item.TicketPrice,
			DiscountAmount:         item.DiscountAmount,
			CheckedIncome:          item.CheckedIncome,
			CheckedTickets:         int32(item.CheckedTickets),
			UncheckedIncome:        item.UncheckedIncome,
			UncheckedTickets:       int32(item.UncheckedTickets),
			CheckedRefund:          item.CheckedRefund,
			CheckedRefundTickets:   int32(item.CheckedRefundTickets),
			UncheckedRefund:        item.UncheckedRefund,
			UncheckedRefundTickets: int32(item.UncheckedRefundTickets),
		}
	}
	//赋值成功响应，返回完整结果
	resp.Success = true
	resp.Msg = "查询成功"
	resp.List = items
	resp.Total = total
	return
}

// ========== 班次结算 ==========

func (s *FinanceServiceImpl) GetScheduleSettle(ctx context.Context, req *finance.ScheduleSettleReq) (resp *finance.ScheduleSettleResp, err error) {
	resp = &finance.ScheduleSettleResp{}
	// 分页参数容错处理
	page, size := int(req.Page), int(req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	//执行线路维度的分页筛选查询
	list, total, err := dao.GetScheduleSettleList(req.StartDate, req.EndDate, req.RouteId, req.ScheduleId, page, size)
	if err != nil {
		resp.Success = false
		resp.Msg = "查询失败: " + err.Error()
		return resp, nil
	}
	//数据格式转换：Model 层 → API 层
	items := make([]*finance.ScheduleSettleItem, len(list))
	for i, item := range list {
		items[i] = &finance.ScheduleSettleItem{
			Id:                     int64(item.ID),
			SettleDate:             item.SettleDate.Format("2006-01-02"),
			RouteNo:                item.RouteNo,
			RouteName:              item.RouteName,
			DepartureTime:          item.DepartureTime,
			TicketPrice:            item.TicketPrice,
			DiscountAmount:         item.DiscountAmount,
			CheckedIncome:          item.CheckedIncome,
			CheckedTickets:         int32(item.CheckedTickets),
			UncheckedIncome:        item.UncheckedIncome,
			UncheckedTickets:       int32(item.UncheckedTickets),
			CheckedRefund:          item.CheckedRefund,
			CheckedRefundTickets:   int32(item.CheckedRefundTickets),
			UncheckedRefund:        item.UncheckedRefund,
			UncheckedRefundTickets: int32(item.UncheckedRefundTickets),
		}
	}

	resp.Success = true
	resp.Msg = "查询成功"
	resp.List = items
	resp.Total = total
	return
}

// ========== 站点结算 ==========

func (s *FinanceServiceImpl) GetStationSettle(ctx context.Context, req *finance.StationSettleReq) (resp *finance.StationSettleResp, err error) {
	resp = &finance.StationSettleResp{}
	// 分页参数容错处理
	page, size := int(req.Page), int(req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	//执行线路维度的分页筛选查询
	list, total, err := dao.GetStationSettleList(req.StartDate, req.EndDate, req.StationId, page, size)
	if err != nil {
		resp.Success = false
		resp.Msg = "查询失败: " + err.Error()
		return resp, nil
	}
	//数据格式转换：Model 层 → API 层
	items := make([]*finance.StationSettleItem, len(list))
	for i, item := range list {
		items[i] = &finance.StationSettleItem{
			Id:                     int64(item.ID),
			SettleDate:             item.SettleDate.Format("2006-01-02"),
			StationName:            item.StationName,
			TicketPrice:            item.TicketPrice,
			DiscountAmount:         item.DiscountAmount,
			CheckedIncome:          item.CheckedIncome,
			CheckedTickets:         int32(item.CheckedTickets),
			UncheckedIncome:        item.UncheckedIncome,
			UncheckedTickets:       int32(item.UncheckedTickets),
			CheckedRefund:          item.CheckedRefund,
			CheckedRefundTickets:   int32(item.CheckedRefundTickets),
			UncheckedRefund:        item.UncheckedRefund,
			UncheckedRefundTickets: int32(item.UncheckedRefundTickets),
		}
	}

	resp.Success = true
	resp.Msg = "查询成功"
	resp.List = items
	resp.Total = total
	return
}

// ========== 司机结算 ==========

func (s *FinanceServiceImpl) GetDriverSettle(ctx context.Context, req *finance.DriverSettleReq) (resp *finance.DriverSettleResp, err error) {
	resp = &finance.DriverSettleResp{}
	// 分页参数容错处理
	page, size := int(req.Page), int(req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	//执行线路维度的分页筛选查询
	list, total, err := dao.GetDriverSettleList(req.StartDate, req.EndDate, req.DriverId, req.DriverName, page, size)
	if err != nil {
		resp.Success = false
		resp.Msg = "查询失败: " + err.Error()
		return resp, nil
	}
	//数据格式转换：Model 层 → API 层
	items := make([]*finance.DriverSettleItem, len(list))
	for i, item := range list {
		items[i] = &finance.DriverSettleItem{
			Id:          int64(item.ID),
			SettleDate:  item.SettleDate.Format("2006-01-02"),
			DriverId:    int64(item.DriverId),
			DriverName:  item.DriverName,
			Phone:       item.Phone,
			TripCount:   int32(item.TripCount),
			TotalIncome: item.TotalIncome,
			Commission:  item.Commission,
			NetIncome:   item.NetIncome,
			Status:      int32(item.Status),
		}
	}

	resp.Success = true
	resp.Msg = "查询成功"
	resp.List = items
	resp.Total = total
	return
}

// ========== 交易明细 ==========

func (s *FinanceServiceImpl) GetTransactions(ctx context.Context, req *finance.TransactionReq) (resp *finance.TransactionResp, err error) {
	resp = &finance.TransactionResp{}
	// 分页参数容错处理
	page, size := int(req.Page), int(req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	//执行线路维度的分页筛选查询
	list, total, err := dao.GetTransactionList(req.StartDate, req.EndDate, req.OrderNo, int(req.PaymentMethod), int(req.TransType), page, size)
	if err != nil {
		resp.Success = false
		resp.Msg = "查询失败: " + err.Error()
		return resp, nil
	}
	//数据格式转换：Model 层 → API 层
	items := make([]*finance.TransactionItem, len(list))
	for i, item := range list {
		items[i] = &finance.TransactionItem{
			Id:                int64(item.ID),
			TransDate:         item.TransDate.Format("2006-01-02"),
			TransTime:         item.TransTime.Format("15:04:05"),
			OrderNo:           item.OrderNo,
			DriverId:          int64(item.DriverId),
			DriverName:        item.DriverName,
			PassengerId:       int64(item.PassengerId),
			PassengerName:     item.PassengerName,
			StartStation:      item.StartStation,
			EndStation:        item.EndStation,
			Amount:            item.Amount,
			PaymentMethod:     int32(item.PaymentMethod),
			PaymentMethodText: model.GetPaymentMethodText(item.PaymentMethod),
			TransType:         int32(item.TransType),
			TransTypeText:     model.GetTransTypeText(item.TransType),
			Remark:            item.Remark,
		}
	}

	resp.Success = true
	resp.Msg = "查询成功"
	resp.List = items
	resp.Total = total
	return
}

// ========== 异常交易 ==========

func (s *FinanceServiceImpl) GetAbnormalTrans(ctx context.Context, req *finance.AbnormalTransReq) (resp *finance.AbnormalTransResp, err error) {
	resp = &finance.AbnormalTransResp{}

	page, size := int(req.Page), int(req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	list, total, err := dao.GetAbnormalTransList(req.StartDate, req.EndDate, int(req.AbnormalType), int(req.Status), page, size)
	if err != nil {
		resp.Success = false
		resp.Msg = "查询失败: " + err.Error()
		return resp, nil
	}

	items := make([]*finance.AbnormalTransItem, len(list))
	for i, item := range list {
		handleTime := ""
		if item.HandleTime != nil {
			handleTime = item.HandleTime.Format("2006-01-02 15:04:05")
		}

		statusText := "待处理"
		if item.Status == model.SettleStatusDone {
			statusText = "已处理"
		}

		items[i] = &finance.AbnormalTransItem{
			Id:               int64(item.ID),
			TransDate:        item.TransDate.Format("2006-01-02"),
			OrderNo:          item.OrderNo,
			AbnormalType:     int32(item.AbnormalType),
			AbnormalTypeText: model.GetAbnormalTypeText(item.AbnormalType),
			ExpectedAmount:   item.ExpectedAmount,
			ActualAmount:     item.ActualAmount,
			DiffAmount:       item.DiffAmount,
			Status:           int32(item.Status),
			StatusText:       statusText,
			Handler:          item.Handler,
			HandleTime:       handleTime,
			HandleRemark:     item.HandleRemark,
		}
	}

	resp.Success = true
	resp.Msg = "查询成功"
	resp.List = items
	resp.Total = total
	return
}

func (s *FinanceServiceImpl) HandleAbnormal(ctx context.Context, req *finance.HandleAbnormalReq) (resp *finance.HandleAbnormalResp, err error) {
	resp = &finance.HandleAbnormalResp{}

	if req.Id <= 0 {
		resp.Success = false
		resp.Msg = "异常交易ID不能为空"
		return resp, nil
	}

	err = dao.HandleAbnormalTrans(req.Id, req.Handler, req.Remark, int(req.HandleType))
	if err != nil {
		resp.Success = false
		resp.Msg = "处理失败: " + err.Error()
		return resp, nil
	}

	resp.Success = true
	resp.Msg = "处理成功"
	return
}

// ========== 生成账单 ==========
// 财务模块的账单生成核心接口
func (s *FinanceServiceImpl) GenerateBill(ctx context.Context, req *finance.GenerateBillReq) (resp *finance.GenerateBillResp, err error) {
	resp = &finance.GenerateBillResp{}
	//按账单类型分支处理（
	var bill *model.Bill

	switch req.BillType {
	case model.BillTypeDaily:
		date, _ := time.Parse("2006-01-02", req.StartDate)
		bill, err = dao.GenerateDailyBill(date)
	case model.BillTypeWeekly:
		startDate, _ := time.Parse("2006-01-02", req.StartDate)
		endDate, _ := time.Parse("2006-01-02", req.EndDate)
		bill, err = dao.GenerateWeeklyBill(startDate, endDate)
	case model.BillTypeMonthly:
		date, _ := time.Parse("2006-01-02", req.StartDate)
		bill, err = dao.GenerateMonthlyBill(date.Year(), int(date.Month()))
	default:
		resp.Success = false
		resp.Msg = "无效的账单类型"
		return resp, nil
	}
	// 账单生成错误处理
	if err != nil {
		resp.Success = false
		resp.Msg = "生成账单失败: " + err.Error()
		return resp, nil
	}
	// 赋值成功响应，返回账单 ID
	resp.Success = true
	resp.Msg = "生成成功"
	resp.BillId = int64(bill.ID)
	return
}

// ========== 导出 ==========

func (s *FinanceServiceImpl) Export(ctx context.Context, req *finance.ExportReq) (resp *finance.ExportResp, err error) {
	resp = &finance.ExportResp{}

	// TODO: 实现导出逻辑，生成Excel文件并返回下载链接
	// 这里返回一个示例URL
	resp.Success = true
	resp.Msg = "导出成功"
	resp.FileUrl = "/exports/finance_" + time.Now().Format("20060102150405") + ".xlsx"
	return
}
