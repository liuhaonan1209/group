package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"group/kitex_gen/car/route"
	"group/kitex_gen/car/route/routeservice"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

var cli routeservice.Client

func main() {
	c, err := routeservice.NewClient("car.route", client.WithHostPorts("0.0.0.0:8892"))
	if err != nil {
		log.Fatal(err)
	}
	cli = c

	hz := server.New(server.WithHostPorts("localhost:8890"))

	hz.GET("/", Index)
	// 班线管理
	hz.POST("/api/route/create", RouteCreate)
	hz.GET("/api/route/list", RouteList)
	hz.GET("/api/route/info", RouteInfo)
	hz.PUT("/api/route/update", RouteUpdate)
	hz.DELETE("/api/route/delete", RouteDelete)

	// 站点管理
	hz.POST("/api/route/station/add", AddStation)
	hz.DELETE("/api/route/station/remove", RemoveStation)

	// 班次管理
	hz.POST("/api/schedule/create", ScheduleCreate)
	hz.GET("/api/schedule/list", ScheduleList)
	hz.PUT("/api/schedule/update", ScheduleUpdate)
	hz.DELETE("/api/schedule/delete", ScheduleDelete)

	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

func Index(ctx context.Context, c *app.RequestContext) {
	b, err := os.ReadFile("route-management.html")
	if err != nil {
		c.String(500, "页面未找到")
		return
	}
	c.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	c.String(200, string(b))
}

// ========== 班线管理 ==========

// RouteCreate 创建班线
func RouteCreate(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		RouteNo        string `json:"routeNo"`
		Fieet          string `json:"fieet"`
		Tage           string `json:"tage"`
		StartStationId int64  `json:"startStationId"`
		EndStationId   int64  `json:"endStationId"`
		UpStations     []struct {
			StationId  int64 `json:"stationId"`
			TravelTime int32 `json:"travelTime"`
		} `json:"upStations"`
		DownStations []struct {
			StationId  int64 `json:"stationId"`
			TravelTime int32 `json:"travelTime"`
		} `json:"downStations"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, utils.H{"success": false, "msg": "参数错误"})
		return
	}

	req := route.NewRouteCreateReq()
	req.RouteNo = reqBody.RouteNo
	req.Fieet = reqBody.Fieet
	req.Tage = reqBody.Tage
	req.StartStationId = reqBody.StartStationId
	req.EndStationId = reqBody.EndStationId

	// 转换上车站点
	for _, s := range reqBody.UpStations {
		req.UpStations = append(req.UpStations, &route.StationInfo{
			StationId:  s.StationId,
			TravelTime: s.TravelTime,
		})
	}

	// 转换下车站点
	for _, s := range reqBody.DownStations {
		req.DownStations = append(req.DownStations, &route.StationInfo{
			StationId:  s.StationId,
			TravelTime: s.TravelTime,
		})
	}

	resp, err := cli.RouteCreate(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
		"routeId": resp.RouteId,
	})
}

// RouteList 班线列表
func RouteList(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	req := route.NewRouteListReq()
	req.RouteNo = c.Query("routeNo")
	req.Fieet = c.Query("fieet")
	req.Tage = c.Query("tage")
	req.Page = int32(page)
	req.Size = int32(size)

	resp, err := cli.RouteList(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
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

// RouteInfo 班线详情
func RouteInfo(ctx context.Context, c *app.RequestContext) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id <= 0 {
		c.JSON(400, utils.H{"success": false, "msg": "班线ID不能为空"})
		return
	}

	req := route.NewRouteInfoReq()
	req.Id = id

	resp, err := cli.RouteInfo(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success":        resp.Success,
		"msg":            resp.Msg,
		"id":             resp.Id,
		"routeNo":        resp.RouteNo,
		"fieet":          resp.Fieet,
		"tage":           resp.Tage,
		"startStationId": resp.StartStationId,
		"endStationId":   resp.EndStationId,
		"upStations":     resp.UpStations,
		"downStations":   resp.DownStations,
		"isActive":       resp.IsActive,
	})
}

// RouteUpdate 更新班线
func RouteUpdate(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		Id             int64  `json:"id"`
		RouteNo        string `json:"routeNo"`
		Fieet          string `json:"fieet"`
		Tage           string `json:"tage"`
		StartStationId int64  `json:"startStationId"`
		EndStationId   int64  `json:"endStationId"`
		IsActive       bool   `json:"isActive"`
		UpStations     []struct {
			StationId  int64 `json:"stationId"`
			TravelTime int32 `json:"travelTime"`
		} `json:"upStations"`
		DownStations []struct {
			StationId  int64 `json:"stationId"`
			TravelTime int32 `json:"travelTime"`
		} `json:"downStations"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, utils.H{"success": false, "msg": "参数错误"})
		return
	}

	req := route.NewRouteUpdateReq()
	req.Id = reqBody.Id
	req.RouteNo = reqBody.RouteNo
	req.Fieet = reqBody.Fieet
	req.Tage = reqBody.Tage
	req.StartStationId = reqBody.StartStationId
	req.EndStationId = reqBody.EndStationId
	req.IsActive = reqBody.IsActive

	for _, s := range reqBody.UpStations {
		req.UpStations = append(req.UpStations, &route.StationInfo{
			StationId:  s.StationId,
			TravelTime: s.TravelTime,
		})
	}

	for _, s := range reqBody.DownStations {
		req.DownStations = append(req.DownStations, &route.StationInfo{
			StationId:  s.StationId,
			TravelTime: s.TravelTime,
		})
	}

	resp, err := cli.RouteUpdate(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
	})
}

// RouteDelete 删除班线
func RouteDelete(ctx context.Context, c *app.RequestContext) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id <= 0 {
		c.JSON(400, utils.H{"success": false, "msg": "班线ID不能为空"})
		return
	}

	req := route.NewRouteDeleteReq()
	req.Id = id

	resp, err := cli.RouteDelete(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
	})
}

// ========== 站点管理 ==========

// AddStation 添加站点到班线
func AddStation(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		RouteId    int64 `json:"routeId"`
		StationId  int64 `json:"stationId"`
		TravelTime int32 `json:"travelTime"`
		StopType   int32 `json:"stopType"` // 1-上车站点，2-下车站点
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, utils.H{"success": false, "msg": "参数错误"})
		return
	}

	req := route.NewAddStationReq()
	req.RouteId = reqBody.RouteId
	req.StationId = reqBody.StationId
	req.TravelTime = reqBody.TravelTime
	req.StopType = reqBody.StopType

	resp, err := cli.AddStation(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
	})
}

// RemoveStation 从班线移除站点
func RemoveStation(ctx context.Context, c *app.RequestContext) {
	routeId, _ := strconv.ParseInt(c.Query("routeId"), 10, 64)
	stationId, _ := strconv.ParseInt(c.Query("stationId"), 10, 64)

	if routeId <= 0 || stationId <= 0 {
		c.JSON(400, utils.H{"success": false, "msg": "班线ID和站点ID不能为空"})
		return
	}

	req := route.NewRemoveStationReq()
	req.RouteId = routeId
	req.StationId = stationId

	resp, err := cli.RemoveStation(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
	})
}

// ========== 班次管理 ==========

// ScheduleCreate 创建班次
func ScheduleCreate(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		RouteId       int64  `json:"routeId"`
		DepartureTime string `json:"departureTime"` // HH:mm
		ArrivalTime   string `json:"arrivalTime"`   // HH:mm
		Capacity      int32  `json:"capacity"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, utils.H{"success": false, "msg": "参数错误"})
		return
	}

	req := route.NewScheduleCreateReq()
	req.RouteId = reqBody.RouteId
	req.DepartureTime = reqBody.DepartureTime
	req.ArrivalTime = reqBody.ArrivalTime
	req.Capacity = reqBody.Capacity

	resp, err := cli.ScheduleCreate(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success":    resp.Success,
		"msg":        resp.Msg,
		"scheduleId": resp.ScheduleId,
	})
}

// ScheduleList 班次列表
func ScheduleList(ctx context.Context, c *app.RequestContext) {
	routeId, _ := strconv.ParseInt(c.Query("routeId"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	req := route.NewScheduleListReq()
	req.RouteId = routeId
	req.Page = int32(page)
	req.Size = int32(size)

	resp, err := cli.ScheduleList(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
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

// ScheduleUpdate 更新班次
func ScheduleUpdate(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		Id            int64  `json:"id"`
		DepartureTime string `json:"departureTime"`
		ArrivalTime   string `json:"arrivalTime"`
		Capacity      int32  `json:"capacity"`
		IsActive      bool   `json:"isActive"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, utils.H{"success": false, "msg": "参数错误"})
		return
	}

	req := route.NewScheduleUpdateReq()
	req.Id = reqBody.Id
	req.DepartureTime = reqBody.DepartureTime
	req.ArrivalTime = reqBody.ArrivalTime
	req.Capacity = reqBody.Capacity
	req.IsActive = reqBody.IsActive

	resp, err := cli.ScheduleUpdate(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
	})
}

// ScheduleDelete 删除班次
func ScheduleDelete(ctx context.Context, c *app.RequestContext) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	if id <= 0 {
		c.JSON(400, utils.H{"success": false, "msg": "班次ID不能为空"})
		return
	}

	req := route.NewScheduleDeleteReq()
	req.Id = id

	resp, err := cli.ScheduleDelete(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, utils.H{"success": false, "msg": "服务调用失败"})
		return
	}

	c.JSON(200, utils.H{
		"success": resp.Success,
		"msg":     resp.Msg,
	})
}
