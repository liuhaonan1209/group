// Package main 评价管理API服务
// 提供HTTP RESTful API接口，作为评价管理RPC服务的网关层
package main

import (
	"context"
	"log"
	"strconv"

	"group/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// main 主函数
func main() {
	// 创建Hertz HTTP服务器
	hz := server.New(server.WithHostPorts("127.0.0.1:9995"))

	// 使用 CORS 中间件
	hz.Use(middleware.CORS())

	// 注册HTTP路由
	hz.POST("/api/rating/submit", RatingSubmit)               // 评价提交
	hz.POST("/api/rating/query", RatingQuery)                 // 评价查询
	hz.GET("/api/rating/history/:orderId", OrderHistoryQuery) // 订单历史变更查询
	hz.POST("/api/rating/tickets", TicketRecordQuery)         // 购票记录查询

	log.Println("Rating API 服务启动在 localhost:9995")
	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

// RatingSubmit 评价提交接口
func RatingSubmit(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		OrderID   int64    `json:"order_id"`
		UserID    int64    `json:"user_id"`
		Score     float64  `json:"score"`
		Tags      []string `json:"tags"`
		Comment   string   `json:"comment"`
		RaterType string   `json:"rater_type"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 模拟RPC调用
	resp := map[string]interface{}{
		"rating_id": 1,
		"success":   true,
		"message":   "评价提交成功",
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "评价提交成功",
		"data":    resp,
	})
}

// RatingQuery 评价查询接口
func RatingQuery(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		TimeRange   *string `json:"time_range,omitempty"`
		RouteID     *int64  `json:"route_id,omitempty"`
		DriverID    *int64  `json:"driver_id,omitempty"`
		DriverPhone *string `json:"driver_phone,omitempty"`
		Page        *int32  `json:"page,omitempty"`
		PageSize    *int32  `json:"page_size,omitempty"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 模拟RPC调用
	resp := map[string]interface{}{
		"ratings":   []interface{}{},
		"total":     0,
		"page":      1,
		"page_size": 10,
		"success":   true,
		"message":   "查询成功",
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "查询成功",
		"data":    resp,
	})
}

// OrderHistoryQuery 订单历史变更查询接口
func OrderHistoryQuery(ctx context.Context, c *app.RequestContext) {
	orderIDStr := c.Param("orderId")

	_, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "无效的订单ID",
		})
		return
	}

	// 模拟RPC调用
	resp := map[string]interface{}{
		"records": []interface{}{},
		"success": true,
		"message": "查询成功",
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "查询成功",
		"data":    resp,
	})
}

// TicketRecordQuery 购票记录查询接口
func TicketRecordQuery(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		UserID    int64   `json:"user_id"`
		TimeRange *string `json:"time_range,omitempty"`
		Page      *int32  `json:"page,omitempty"`
		PageSize  *int32  `json:"page_size,omitempty"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 模拟RPC调用
	resp := map[string]interface{}{
		"records":   []interface{}{},
		"total":     0,
		"page":      1,
		"page_size": 10,
		"success":   true,
		"message":   "查询成功",
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "查询成功",
		"data":    resp,
	})
}
