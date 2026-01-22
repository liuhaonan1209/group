package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"group/kitex_gen/car/finance"
	"group/kitex_gen/car/finance/financeservice"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

var cli financeservice.Client

func main() {
	c, err := financeservice.NewClient("car.finance", client.WithHostPorts("0.0.0.0:8894"))
	if err != nil {
		log.Fatal(err)
	}
	cli = c

	hz := server.New(server.WithHostPorts("localhost:8895"))

	// 收支对账
	hz.GET("/api/finance/balance", GetBalanceSheet)

	// 收入对账
	hz.GET("/api/finance/income", GetIncomeSheet)

	// 线路结算
	hz.GET("/api/finance/route-settle", GetRouteSettle)

	// 班次结算
	hz.GET("/api/finance/schedule-settle", GetScheduleSettle)

	// 站点结算
	hz.GET("/api/finance/station-settle", GetStationSettle)

	// 司机结算
	hz.GET("/api/finance/driver-settle", GetDriverSettle)

	// 交易明细
	hz.GET("/api/finance/transactions", GetTransactions)

	// 异常交易
	hz.GET("/api/finance/abnormal", GetAbnormalTrans)
	hz.POST("/api/finance/abnormal/handle", HandleAbnormal)

	// 生成账单
	hz.POST("/api/finance/bill/generate", GenerateBill)

	// 导出
	hz.GET("/api/finance/export", Export)

	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

// ========== 收支对账 ==========

func GetBalanceSheet(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))

	req := finance.NewBalanceSheetReq()
	req.StartDate = c.Query("startDate")
	req.EndDate = c.Query("endDate")
	req.Status = int32(status)
	req.Page = int32(page)
	req.Size = int32(size)

	resp, err := cli.GetBalanceSheet(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success":     resp.Success,
		"msg":         resp.Msg,
		"list":        resp.List,
		"total":       resp.Total,
		"totalIncome": resp.TotalIncome,
		"totalRefund": resp.TotalRefund,
		"totalSettle": resp.TotalSettle,
	})
}

// ========== 收入对账 ==========

func GetIncomeSheet(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))

	req := finance.NewIncomeSheetReq()
	req.StartDate = c.Query("startDate")
	req.EndDate = c.Query("endDate")
	req.Status = int32(status)
	req.Page = int32(page)
	req.Size = int32(size)

	resp, err := cli.GetIncomeSheet(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
		"list":    resp.List,
		"total":   resp.Total,
	})
}

// ========== 线路结算 ==========

func GetRouteSettle(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	routeId, _ := strconv.ParseInt(c.Query("routeId"), 10, 64)

	req := finance.NewRouteSettleReq()
	req.StartDate = c.Query("startDate")
	req.EndDate = c.Query("endDate")
	req.RouteId = routeId
	req.Fleet = c.Query("fleet")
	req.Page = int32(page)
	req.Size = int32(size)

	resp, err := cli.GetRouteSettle(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
		"list":    resp.List,
		"total":   resp.Total,
	})
}

// ========== 班次结算 ==========

func GetScheduleSettle(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	routeId, _ := strconv.ParseInt(c.Query("routeId"), 10, 64)
	scheduleId, _ := strconv.ParseInt(c.Query("scheduleId"), 10, 64)

	req := finance.NewScheduleSettleReq()
	req.StartDate = c.Query("startDate")
	req.EndDate = c.Query("endDate")
	req.RouteId = routeId
	req.ScheduleId = scheduleId
	req.Page = int32(page)
	req.Size = int32(size)

	resp, err := cli.GetScheduleSettle(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
		"list":    resp.List,
		"total":   resp.Total,
	})
}

// ========== 站点结算 ==========

func GetStationSettle(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	stationId, _ := strconv.ParseInt(c.Query("stationId"), 10, 64)

	req := finance.NewStationSettleReq()
	req.StartDate = c.Query("startDate")
	req.EndDate = c.Query("endDate")
	req.StationId = stationId
	req.Page = int32(page)
	req.Size = int32(size)

	resp, err := cli.GetStationSettle(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
		"list":    resp.List,
		"total":   resp.Total,
	})
}

// ========== 司机结算 ==========

func GetDriverSettle(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	driverId, _ := strconv.ParseInt(c.Query("driverId"), 10, 64)

	req := finance.NewDriverSettleReq()
	req.StartDate = c.Query("startDate")
	req.EndDate = c.Query("endDate")
	req.DriverId = driverId
	req.DriverName = c.Query("driverName")
	req.Page = int32(page)
	req.Size = int32(size)

	resp, err := cli.GetDriverSettle(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
		"list":    resp.List,
		"total":   resp.Total,
	})
}

// ========== 交易明细 ==========

func GetTransactions(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	paymentMethod, _ := strconv.Atoi(c.DefaultQuery("paymentMethod", "0"))
	transType, _ := strconv.Atoi(c.DefaultQuery("transType", "0"))

	req := finance.NewTransactionReq()
	req.StartDate = c.Query("startDate")
	req.EndDate = c.Query("endDate")
	req.OrderNo = c.Query("orderNo")
	req.PaymentMethod = int32(paymentMethod)
	req.TransType = int32(transType)
	req.Page = int32(page)
	req.Size = int32(size)

	resp, err := cli.GetTransactions(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
		"list":    resp.List,
		"total":   resp.Total,
	})
}

// ========== 异常交易 ==========

func GetAbnormalTrans(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	abnormalType, _ := strconv.Atoi(c.DefaultQuery("abnormalType", "0"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))

	req := finance.NewAbnormalTransReq()
	req.StartDate = c.Query("startDate")
	req.EndDate = c.Query("endDate")
	req.AbnormalType = int32(abnormalType)
	req.Status = int32(status)
	req.Page = int32(page)
	req.Size = int32(size)

	resp, err := cli.GetAbnormalTrans(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
		"list":    resp.List,
		"total":   resp.Total,
	})
}

func HandleAbnormal(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		Id         int64  `json:"id"`
		Handler    string `json:"handler"`
		Remark     string `json:"remark"`
		HandleType int32  `json:"handleType"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, utils.H{"success": false, "msg": "参数错误"})
		return
	}

	req := finance.NewHandleAbnormalReq()
	req.Id = reqBody.Id
	req.Handler = reqBody.Handler
	req.Remark = reqBody.Remark
	req.HandleType = reqBody.HandleType

	resp, err := cli.HandleAbnormal(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
	})
}

// ========== 生成账单 ==========

func GenerateBill(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		BillType  int32  `json:"billType"`
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, utils.H{"success": false, "msg": "参数错误"})
		return
	}

	req := finance.NewGenerateBillReq()
	req.BillType = reqBody.BillType
	req.StartDate = reqBody.StartDate
	req.EndDate = reqBody.EndDate

	resp, err := cli.GenerateBill(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
		"billId":  resp.BillId,
	})
}

// ========== 导出 ==========

func Export(ctx context.Context, c *app.RequestContext) {
	exportType, _ := strconv.Atoi(c.Query("exportType"))

	req := finance.NewExportReq()
	req.ExportType = int32(exportType)
	req.StartDate = c.Query("startDate")
	req.EndDate = c.Query("endDate")

	resp, err := cli.Export(context.Background(), req, callopt.WithRPCTimeout(10*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
		"fileUrl": resp.FileUrl,
	})
}
