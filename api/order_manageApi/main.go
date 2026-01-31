// Package main 订单管理API服务
// 提供HTTP RESTful API接口，作为订单管理RPC服务的网关层
// 负责接收HTTP请求，转换为RPC调用，并返回HTTP响应
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"group/kitex_gen/car/order_manage"
	"group/kitex_gen/car/order_manage/ordermanageservice"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

var (
	cli ordermanageservice.Client
)

// main 主函数
// 初始化RPC客户端和HTTP服务器，注册路由并启动服务
func main() {
	// 创建订单管理RPC客户端
	c, err := ordermanageservice.NewClient("car.order_manage", client.WithHostPorts("127.0.0.1:8999"))
	if err != nil {
		log.Fatal("初始化订单管理RPC客户端失败:", err)
	}
	cli = c

	// 创建Hertz HTTP服务器
	// 监听localhost:9992端口
	hz := server.New(server.WithHostPorts("127.0.0.1:9992"))

	// 注册HTTP路由
	hz.POST("/api/order/create", CreateOrder)              // 创建订单
	hz.GET("/api/order/payment-status", GetPaymentStatus)  // 查询支付状态
	hz.POST("/api/order/query", QueryOrders)               // 查询订单列表
	hz.GET("/api/order/detail/:id", GetOrderDetail)        // 查询订单详情
	hz.POST("/api/order/update-status", UpdateOrderStatus) // 更新订单状态

	log.Println("Order Manage API 服务启动在 localhost:9992")
	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

// CreateOrder 创建订单接口
// HTTP POST /api/order/create
// 乘客下单预定车位
// 请求体:
//   - trip_id: 行程ID
//   - order_type: 订单类型（bus/carpool）
//   - passenger_id: 乘客ID
//   - passenger_name: 乘客姓名
//   - passenger_tel: 乘客手机号
//   - start_point: 起点
//   - end_point: 终点
//   - departure_time: 出发时间
//   - price: 价格
//   - payment_method: 支付方式
//   - seat_count: 座位数量
//   - route_id: 线路ID（可选）
//   - special_needs: 特殊需求（可选）
//
// 返回:
//   - 200: 创建成功
//   - 400: 参数错误或创建失败
//   - 500: 服务器错误
func CreateOrder(ctx context.Context, c *app.RequestContext) {
	// 定义请求体结构
	var reqBody struct {
		TripID        int64   `json:"trip_id"`
		OrderType     string  `json:"order_type"`
		PassengerID   int64   `json:"passenger_id"`
		PassengerName string  `json:"passenger_name"`
		PassengerTel  string  `json:"passenger_tel"`
		StartPoint    string  `json:"start_point"`
		EndPoint      string  `json:"end_point"`
		DepartureTime string  `json:"departure_time"`
		Price         float64 `json:"price"`
		PaymentMethod string  `json:"payment_method"`
		SeatCount     int32   `json:"seat_count"`
		RouteID       string  `json:"route_id,omitempty"`
		SpecialNeeds  string  `json:"special_needs,omitempty"`
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
	req := order_manage.NewOrderCreateReq()
	req.TripId = reqBody.TripID
	req.OrderType = reqBody.OrderType
	req.PassengerId = reqBody.PassengerID
	req.PassengerName = reqBody.PassengerName
	req.PassengerTel = reqBody.PassengerTel
	req.StartPoint = reqBody.StartPoint
	req.EndPoint = reqBody.EndPoint
	req.DepartureTime = reqBody.DepartureTime
	req.Price = reqBody.Price
	req.PaymentMethod = reqBody.PaymentMethod
	req.SeatCount = reqBody.SeatCount
	
	if reqBody.RouteID != "" {
		req.RouteId = &reqBody.RouteID
	}
	if reqBody.SpecialNeeds != "" {
		req.SpecialNeeds = &reqBody.SpecialNeeds
	}

	// 调用RPC服务创建订单
	resp, err := cli.OrderCreate(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "创建订单失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查创建结果
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
			"order_id": resp.OrderId,
			"order_no": resp.OrderNo,
		},
	})
}

// GetPaymentStatus 查询支付状态接口
// HTTP GET /api/order/payment-status?order_id=xxx&user_id=xxx
// 查询订单支付状态
// 参数:
//   - order_id: 订单ID
//   - user_id: 用户ID
//
// 返回:
//   - 200: 查询成功
//   - 400: 参数错误或查询失败
//   - 500: 服务器错误
func GetPaymentStatus(ctx context.Context, c *app.RequestContext) {
	// 获取查询参数
	orderIDStr := c.Query("order_id")
	userIDStr := c.Query("user_id")

	// 参数验证
	if orderIDStr == "" || userIDStr == "" {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "订单ID和用户ID为必填项",
		})
		return
	}

	// 转换参数类型
	var orderID, userID int64
	if _, err := fmt.Sscanf(orderIDStr, "%d", &orderID); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "无效的订单ID",
		})
		return
	}
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	// 构造RPC请求
	req := order_manage.NewOrderPaymentStatusReq()
	req.OrderId = orderID
	req.UserId = userID

	// 调用RPC服务查询支付状态
	resp, err := cli.OrderPaymentStatus(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "查询支付状态失败",
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

	// 返回成功响应
	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data": map[string]interface{}{
			"payment_status": resp.PaymentStatus,
			"payment_method": resp.PaymentMethod,
			"payment_amount": resp.PaymentAmount,
			"payment_time":   resp.PaymentTime,
		},
	})
}

// QueryOrders 查询订单列表接口
// HTTP POST /api/order/query
// 按条件查询班车/顺风车订单
// 请求体:
//   - order_type: 订单类型（可选）
//   - status: 订单状态（可选）
//   - start_time: 开始时间（可选）
//   - end_time: 结束时间（可选）
//   - route_id: 线路ID（可选）
//   - passenger_tel: 乘客手机号（可选）
//   - passenger_id: 乘客ID（可选）
//   - driver_id: 司机ID（可选）
//   - page: 页码（可选，默认1）
//   - page_size: 每页数量（可选，默认10）
//
// 返回:
//   - 200: 查询成功
//   - 400: 参数错误或查询失败
//   - 500: 服务器错误
func QueryOrders(ctx context.Context, c *app.RequestContext) {
	// 定义请求体结构
	var reqBody struct {
		OrderType    *string `json:"order_type,omitempty"`
		Status       *string `json:"status,omitempty"`
		StartTime    *string `json:"start_time,omitempty"`
		EndTime      *string `json:"end_time,omitempty"`
		RouteID      *string `json:"route_id,omitempty"`
		PassengerTel *string `json:"passenger_tel,omitempty"`
		PassengerID  *int64  `json:"passenger_id,omitempty"`
		DriverID     *int64  `json:"driver_id,omitempty"`
		Page         *int32  `json:"page,omitempty"`
		PageSize     *int32  `json:"page_size,omitempty"`
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
	req := order_manage.NewOrderQueryReq()
	req.OrderType = reqBody.OrderType
	req.Status = reqBody.Status
	req.StartTime = reqBody.StartTime
	req.EndTime = reqBody.EndTime
	req.RouteId = reqBody.RouteID
	req.PassengerTel = reqBody.PassengerTel
	req.PassengerId = reqBody.PassengerID
	req.DriverId = reqBody.DriverID
	req.Page = reqBody.Page
	req.PageSize = reqBody.PageSize

	// 调用RPC服务查询订单
	resp, err := cli.OrderQuery(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "查询订单失败",
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

	// 转换订单列表
	orders := make([]map[string]interface{}, 0, len(resp.Orders))
	for _, order := range resp.Orders {
		orderMap := map[string]interface{}{
			"order_id":       order.OrderId,
			"order_no":       order.OrderNo,
			"order_type":     order.OrderType,
			"passenger_id":   order.PassengerId,
			"passenger_name": order.PassengerName,
			"passenger_tel":  order.PassengerTel,
			"start_point":    order.StartPoint,
			"end_point":      order.EndPoint,
			"departure_time": order.DepartureTime,
			"status":         order.Status,
			"price":          order.Price,
			"payment_status": order.PaymentStatus,
			"payment_method": order.PaymentMethod,
			"created_at":     order.CreatedAt,
		}

		if order.DriverId != nil {
			orderMap["driver_id"] = *order.DriverId
		}
		if order.DriverName != nil {
			orderMap["driver_name"] = *order.DriverName
		}
		if order.VehicleInfo != nil {
			orderMap["vehicle_info"] = *order.VehicleInfo
		}
		if order.RouteId != nil {
			orderMap["route_id"] = *order.RouteId
		}

		orders = append(orders, orderMap)
	}

	// 返回成功响应
	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data": map[string]interface{}{
			"orders":    orders,
			"total":     resp.Total,
			"page":      resp.Page,
			"page_size": resp.PageSize,
		},
	})
}

// GetOrderDetail 查询订单详情接口
// HTTP GET /api/order/detail/:id?user_id=xxx
// 查询单个订单完整信息
// 参数:
//   - id: 路径参数，订单ID
//   - user_id: 查询参数，用户ID
//
// 返回:
//   - 200: 查询成功
//   - 400: 参数错误或查询失败
//   - 500: 服务器错误
func GetOrderDetail(ctx context.Context, c *app.RequestContext) {
	// 获取路径参数
	orderIDStr := c.Param("id")
	userIDStr := c.Query("user_id")

	// 参数验证
	if orderIDStr == "" || userIDStr == "" {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "订单ID和用户ID为必填项",
		})
		return
	}

	// 转换参数类型
	var orderID, userID int64
	if _, err := fmt.Sscanf(orderIDStr, "%d", &orderID); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "无效的订单ID",
		})
		return
	}
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	// 构造RPC请求
	req := order_manage.NewOrderDetailReq()
	req.OrderId = orderID
	req.UserId = userID

	// 调用RPC服务查询订单详情
	resp, err := cli.OrderDetail(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "查询订单详情失败",
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

	// 构建订单详情响应
	detail := resp.OrderDetail
	orderDetail := map[string]interface{}{
		"order_id":       detail.OrderId,
		"order_no":       detail.OrderNo,
		"order_type":     detail.OrderType,
		"passenger_id":   detail.PassengerId,
		"passenger_name": detail.PassengerName,
		"passenger_tel":  detail.PassengerTel,
		"start_point":    detail.StartPoint,
		"end_point":      detail.EndPoint,
		"departure_time": detail.DepartureTime,
		"status":         detail.Status,
		"price":          detail.Price,
		"payment_status": detail.PaymentStatus,
		"payment_method": detail.PaymentMethod,
		"created_at":     detail.CreatedAt,
	}

	// 可选字段
	if detail.DriverId != nil {
		orderDetail["driver_id"] = *detail.DriverId
	}
	if detail.DriverName != nil {
		orderDetail["driver_name"] = *detail.DriverName
	}
	if detail.DriverTel != nil {
		orderDetail["driver_tel"] = *detail.DriverTel
	}
	if detail.VehicleInfo != nil {
		orderDetail["vehicle_info"] = *detail.VehicleInfo
	}
	if detail.RouteId != nil {
		orderDetail["route_id"] = *detail.RouteId
	}
	if detail.SpecialNeeds != nil {
		orderDetail["special_needs"] = *detail.SpecialNeeds
	}
	if detail.SeatCount != nil {
		orderDetail["seat_count"] = *detail.SeatCount
	}
	if detail.AcceptedAt != nil {
		orderDetail["accepted_at"] = *detail.AcceptedAt
	}
	if detail.CompletedAt != nil {
		orderDetail["completed_at"] = *detail.CompletedAt
	}
	if detail.CancelledAt != nil {
		orderDetail["cancelled_at"] = *detail.CancelledAt
	}

	// 状态历史
	statusHistory := make([]map[string]interface{}, 0, len(detail.StatusHistory))
	for _, history := range detail.StatusHistory {
		historyMap := map[string]interface{}{
			"status":     history.Status,
			"changed_at": history.ChangedAt,
		}
		if history.Reason != nil {
			historyMap["reason"] = *history.Reason
		}
		if history.OperatorId != nil {
			historyMap["operator_id"] = *history.OperatorId
		}
		if history.OperatorName != nil {
			historyMap["operator_name"] = *history.OperatorName
		}
		statusHistory = append(statusHistory, historyMap)
	}
	orderDetail["status_history"] = statusHistory

	// 返回成功响应
	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data":    orderDetail,
	})
}

// UpdateOrderStatus 更新订单状态接口
// HTTP POST /api/order/update-status
// 更新订单状态（接单/完成/取消）
// 请求体:
//   - order_id: 订单ID
//   - operator_id: 操作人ID
//   - operator_type: 操作人类型（passenger/driver）
//   - new_status: 新状态（accepted/in_progress/completed/cancelled）
//   - reason: 变更原因（取消时必填）
//
// 返回:
//   - 200: 更新成功
//   - 400: 参数错误或更新失败
//   - 500: 服务器错误
func UpdateOrderStatus(ctx context.Context, c *app.RequestContext) {
	// 定义请求体结构
	var reqBody struct {
		OrderID      int64  `json:"order_id"`
		OperatorID   int64  `json:"operator_id"`
		OperatorType string `json:"operator_type"`
		NewStatus    string `json:"new_status"`
		Reason       string `json:"reason,omitempty"`
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
	req := order_manage.NewOrderStatusUpdateReq()
	req.OrderId = reqBody.OrderID
	req.OperatorId = reqBody.OperatorID
	req.OperatorType = reqBody.OperatorType
	req.NewStatus_ = reqBody.NewStatus

	if reqBody.Reason != "" {
		req.Reason = &reqBody.Reason
	}

	// 调用RPC服务更新订单状态
	resp, err := cli.OrderStatusUpdate(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "更新订单状态失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查更新结果
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
	})
}
