package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"group/kitex_gen/car/order"
	"group/kitex_gen/car/order/orderservice"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

var (
	cli orderservice.Client
)

func main() {
	// 创建 RPC 客户端
	c, err := orderservice.NewClient("car.order", client.WithHostPorts("127.0.0.1:8892"))
	if err != nil {
		log.Fatal(err)
	}
	cli = c

	// 创建 HTTP 服务器
	hz := server.New(server.WithHostPorts("127.0.0.1:8993"))

	// ==================== 注册路由 ====================
	// 行程管理路由
	hz.POST("/api/trip/publish", PublishTrip)              // 乘客发布行程 - 乘客发布拼车需求
	hz.POST("/api/trip/driver/publish", PublishDriverTrip) // 司机发布行程 - 司机发布可拼车行程
	hz.GET("/api/trip/query", QueryTrips)                  // 查询行程列表 - 根据起点终点时间查询可用行程
	hz.GET("/api/trip/:id", GetTripDetail)                 // 获取行程详情 - 查看单个行程的详细信息
	
	// 辅助功能路由
	hz.POST("/api/trip/help", PassengerHelp)             // 乘客求助 - 乘客在行程中请求帮助
	hz.POST("/api/trip/share", ShareTrip)                // 分享行程 - 将行程分享给好友
	
	// 数据导出路由
	hz.POST("/api/data/export", ExportData)              // 数据导出 - 导出用户的订单、支付、评价等数据
	hz.GET("/api/data/export/records", GetExportRecords) // 查询导出记录 - 查看历史导出任务状态

	log.Println("Order API 服务启动在 localhost:8993")
	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

// PublishTrip 乘客发布行程
func PublishTrip(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		PassengerID   int64  `json:"passenger_id"`
		StartPoint    string `json:"start_point"`
		EndPoint      string `json:"end_point"`
		DepartureTime string `json:"departure_time"`
		SpecialNeeds  string `json:"special_needs,omitempty"`
		ContactWay    string `json:"contact_way,omitempty"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	req := order.NewTripPublishReq()
	req.PassengerId = reqBody.PassengerID
	req.StartPoint = reqBody.StartPoint
	req.EndPoint = reqBody.EndPoint
	req.DepartureTime = reqBody.DepartureTime
	if reqBody.SpecialNeeds != "" {
		req.SpecialNeeds = &reqBody.SpecialNeeds
	}
	if reqBody.ContactWay != "" {
		req.ContactWay = &reqBody.ContactWay
	}

	resp, err := cli.TripPublish(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "发布行程失败",
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
			"trip_id": resp.TripId,
		},
	})
}

// PublishDriverTrip 司机发布行程
func PublishDriverTrip(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		DriverID      int64  `json:"driver_id"`
		StartPoint    string `json:"start_point"`
		EndPoint      string `json:"end_point"`
		DepartureTime string `json:"departure_time"`
		VehicleInfo   string `json:"vehicle_info"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	req := order.NewDriverTripPublishReq()
	req.DriverId = reqBody.DriverID
	req.StartPoint = reqBody.StartPoint
	req.EndPoint = reqBody.EndPoint
	req.DepartureTime = reqBody.DepartureTime
	req.VehicleInfo = reqBody.VehicleInfo

	resp, err := cli.DriverTripPublish(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "发布行程失败",
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
			"trip_id": resp.TripId,
		},
	})
}

// QueryTrips 查询行程
func QueryTrips(ctx context.Context, c *app.RequestContext) {
	startPoint := c.Query("start_point")
	endPoint := c.Query("end_point")
	departureTime := c.Query("departure_time")
	tripType := c.Query("trip_type") // 可选：passenger/driver

	if startPoint == "" || endPoint == "" || departureTime == "" {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "起点、终点和出行时间为必填参数",
		})
		return
	}

	req := order.NewTripQueryReq()
	req.StartPoint = startPoint
	req.EndPoint = endPoint
	req.DepartureTime = departureTime
	if tripType != "" {
		req.TripType = &tripType
	}

	resp, err := cli.TripQuery(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "查询行程失败",
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

	// 转换行程列表
	trips := make([]map[string]interface{}, 0, len(resp.Trips))
	for _, trip := range resp.Trips {
		tripMap := map[string]interface{}{
			"trip_id":        trip.TripId,
			"start_point":    trip.StartPoint,
			"end_point":      trip.EndPoint,
			"departure_time": trip.DepartureTime,
			"trip_type":      trip.TripType,
			"publisher_name": trip.PublisherName,
			"status":         trip.Status,
		}
		if trip.VehicleInfo != nil {
			tripMap["vehicle_info"] = *trip.VehicleInfo
		}
		if trip.SpecialNeeds != nil {
			tripMap["special_needs"] = *trip.SpecialNeeds
		}
		trips = append(trips, tripMap)
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data": map[string]interface{}{
			"trips": trips,
			"total": len(trips),
		},
	})
}

// GetTripDetail 获取行程详情
func GetTripDetail(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")

	var tripId int64
	if _, err := fmt.Sscanf(id, "%d", &tripId); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "无效的行程ID",
		})
		return
	}

	req := order.NewTripDetailReq()
	req.TripId = tripId

	resp, err := cli.TripDetail(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "查询行程详情失败",
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

	trip := resp.TripInfo
	tripData := map[string]interface{}{
		"trip_id":        trip.TripId,
		"start_point":    trip.StartPoint,
		"end_point":      trip.EndPoint,
		"departure_time": trip.DepartureTime,
		"trip_type":      trip.TripType,
		"publisher_name": trip.PublisherName,
		"status":         trip.Status,
	}
	if trip.VehicleInfo != nil {
		tripData["vehicle_info"] = *trip.VehicleInfo
	}
	if trip.SpecialNeeds != nil {
		tripData["special_needs"] = *trip.SpecialNeeds
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data":    tripData,
	})
}

// PassengerHelp 乘客求助
func PassengerHelp(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		TripID      int64  `json:"trip_id"`
		PassengerID int64  `json:"passenger_id"`
		HelpType    string `json:"help_type"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	req := order.NewPassengerHelpReq()
	req.TripId = reqBody.TripID
	req.PassengerId = reqBody.PassengerID
	req.HelpType = reqBody.HelpType

	resp, err := cli.PassengerHelp(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "提交求助失败",
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
			"contact_info": resp.ContactInfo,
		},
	})
}

// ShareTrip 分享行程
func ShareTrip(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		TripID       int64    `json:"trip_id"`
		UserID       int64    `json:"user_id"`
		ShareTargets []string `json:"share_targets"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	req := order.NewTripShareReq()
	req.TripId = reqBody.TripID
	req.UserId = reqBody.UserID
	req.ShareTargets = reqBody.ShareTargets

	resp, err := cli.TripShare(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "分享行程失败",
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
			"share_info": resp.ShareLink,
		},
	})
}

// ExportData 数据导出
func ExportData(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		UserID       int64    `json:"user_id"`
		UserType     string   `json:"user_type"`
		ExportFields []string `json:"export_fields"`
		StartDate    string   `json:"start_date,omitempty"`
		EndDate      string   `json:"end_date,omitempty"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	req := order.NewDataExportReq()
	req.UserId = reqBody.UserID
	req.UserType = reqBody.UserType
	req.ExportFields = reqBody.ExportFields
	if reqBody.StartDate != "" {
		req.StartDate = &reqBody.StartDate
	}
	if reqBody.EndDate != "" {
		req.EndDate = &reqBody.EndDate
	}

	resp, err := cli.DataExport(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "导出失败",
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
			"export_id":    resp.ExportId,
			"download_url": resp.DownloadUrl,
			"file_format":  resp.FileFormat,
		},
	})
}

// GetExportRecords 查询导出记录
func GetExportRecords(ctx context.Context, c *app.RequestContext) {
	userIdStr := c.Query("user_id")
	limitStr := c.Query("limit")

	if userIdStr == "" {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "用户ID为必填参数",
		})
		return
	}

	var userId int64
	if _, err := fmt.Sscanf(userIdStr, "%d", &userId); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	req := order.NewExportRecordQueryReq()
	req.UserId = userId

	if limitStr != "" {
		var limit int32
		if _, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil {
			req.Limit = &limit
		}
	}

	resp, err := cli.ExportRecordQuery(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "查询导出记录失败",
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

	// 转换导出记录列表
	records := make([]map[string]interface{}, 0, len(resp.Records))
	for _, record := range resp.Records {
		recordMap := map[string]interface{}{
			"export_id":     record.ExportId,
			"user_id":       record.UserId,
			"user_type":     record.UserType,
			"export_fields": record.ExportFields,
			"status":        record.Status,
			"download_url":  record.DownloadUrl,
			"created_at":    record.CreatedAt,
			"expires_at":    record.ExpiresAt,
		}
		records = append(records, recordMap)
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data": map[string]interface{}{
			"records": records,
			"total":   len(records),
		},
	})
}
