package main

import (
	"context"
	"fmt"
	"group/handler/dao"
	"group/handler/model"
	"log"
	"time"

	"group/kitex_gen/car/order_manage"
)

// OrderManageServiceImpl implements the last service interface defined in the IDL.
type OrderManageServiceImpl struct{}

// OrderCreate 订单创建接口
func (s *OrderManageServiceImpl) OrderCreate(ctx context.Context, req *order_manage.OrderCreateReq) (*order_manage.OrderCreateResp, error) {
	log.Printf("订单创建请求: 乘客ID=%d, 行程ID=%d, 类型=%s", req.PassengerId, req.TripId, req.OrderType)

	// 1. 生成订单号
	orderNo := generateOrderNo()

	// 2. 解析出发时间
	departureTime, err := time.Parse("2006-01-02 15:04:05", req.DepartureTime)
	if err != nil {
		log.Printf("出发时间格式错误: %v", err)
		return &order_manage.OrderCreateResp{
			Success: false,
			Message: "出发时间格式错误",
		}, nil
	}

	// 3. 创建订单
	orderData := model.OrderManage{
		OrderNo:       orderNo,
		TripID:        req.TripId,
		OrderType:     req.OrderType,
		PassengerID:   req.PassengerId,
		PassengerName: req.PassengerName,
		PassengerTel:  req.PassengerTel,
		StartPoint:    req.StartPoint,
		EndPoint:      req.EndPoint,
		DepartureTime: departureTime,
		Status:        "pending",
		Price:         req.Price,
		PaymentStatus: "pending",
		PaymentMethod: req.PaymentMethod,
		SeatCount:     int(req.SeatCount),
	}

	// 可选字段
	if req.RouteId != nil {
		orderData.RouteID = *req.RouteId
	}
	if req.SpecialNeeds != nil {
		orderData.SpecialNeeds = *req.SpecialNeeds
	}

	// 4. 保存到数据库（调用DAO层）
	if err := dao.CreateOrder(&orderData); err != nil {
		log.Printf("订单创建失败: %v", err)
		return &order_manage.OrderCreateResp{
			Success: false,
			Message: "订单创建失败",
		}, nil
	}

	// 5. 记录状态历史（调用DAO层）
	history := model.OrderStatusHistory{
		OrderID: int64(orderData.ID),
		Status:  "pending",
	}
	err = dao.CreateOrderStatusHistory(&history)
	if err != nil {
		return nil, err
	}

	log.Printf("订单创建成功: ID=%d, 订单号=%s", orderData.ID, orderNo)

	return &order_manage.OrderCreateResp{
		OrderId: int64(orderData.ID),
		OrderNo: orderNo,
		Success: true,
		Message: "订单创建成功",
	}, nil
}

// OrderPaymentStatus 订单支付状态查询接口
func (s *OrderManageServiceImpl) OrderPaymentStatus(ctx context.Context, req *order_manage.OrderPaymentStatusReq) (*order_manage.OrderPaymentStatusResp, error) {
	log.Printf("查询订单支付状态: 订单ID=%d, 用户ID=%d", req.OrderId, req.UserId)

	// 1. 查询订单（调用DAO层，带权限验证）
	var orderData model.OrderManage
	order, err := dao.GetOrderManageByIDAndUserID(&orderData, req.OrderId, req.UserId)
	if err != nil {
		return &order_manage.OrderPaymentStatusResp{
			Success: false,
			Message: "订单不存在或无权限查看",
		}, nil
	}

	// 2. 返回支付状态
	paymentTime := ""
	if order.PaymentTime != nil {
		paymentTime = order.PaymentTime.Format("2006-01-02 15:04:05")
	}

	return &order_manage.OrderPaymentStatusResp{
		PaymentStatus: order.PaymentStatus,
		PaymentMethod: order.PaymentMethod,
		PaymentAmount: order.Price,
		PaymentTime:   paymentTime,
		Success:       true,
		Message:       "查询成功",
	}, nil
}

// OrderQuery 订单查询接口
func (s *OrderManageServiceImpl) OrderQuery(ctx context.Context, req *order_manage.OrderQueryReq) (*order_manage.OrderQueryResp, error) {
	log.Printf("订单查询请求: 类型=%v, 状态=%v", req.OrderType, req.Status)

	// 1. 分页参数
	page := int32(1)
	pageSize := int32(10)
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.PageSize != nil && *req.PageSize > 0 {
		pageSize = *req.PageSize
	}

	// 2. 调用DAO层查询订单列表
	var orderData model.OrderManage
	orders, total, err := dao.QueryOrders(&orderData, req.OrderType, req.Status, req.StartTime, req.EndTime,
		req.RouteId, req.PassengerTel, req.PassengerId, req.DriverId, page, pageSize)
	if err != nil {
		log.Printf("查询订单列表失败: %v", err)
		return &order_manage.OrderQueryResp{
			Success: false,
			Message: "查询失败",
		}, nil
	}

	// 3. 转换为响应格式
	var orderList []*order_manage.OrderInfo
	for _, o := range orders {
		orderInfo := &order_manage.OrderInfo{
			OrderId:       int64(o.ID),
			OrderNo:       o.OrderNo,
			OrderType:     o.OrderType,
			PassengerId:   o.PassengerID,
			PassengerName: o.PassengerName,
			PassengerTel:  o.PassengerTel,
			StartPoint:    o.StartPoint,
			EndPoint:      o.EndPoint,
			DepartureTime: o.DepartureTime.Format("2006-01-02 15:04:05"),
			Status:        o.Status,
			Price:         o.Price,
			PaymentStatus: o.PaymentStatus,
			PaymentMethod: o.PaymentMethod,
			CreatedAt:     o.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if o.DriverID != nil {
			orderInfo.DriverId = o.DriverID
			orderInfo.DriverName = &o.DriverName
		}
		if o.VehicleInfo != "" {
			orderInfo.VehicleInfo = &o.VehicleInfo
		}
		if o.RouteID != "" {
			orderInfo.RouteId = &o.RouteID
		}

		orderList = append(orderList, orderInfo)
	}

	log.Printf("订单查询成功: 总数=%d, 当前页=%d", total, page)

	return &order_manage.OrderQueryResp{
		Orders:   orderList,
		Total:    int32(total),
		Page:     page,
		PageSize: pageSize,
		Success:  true,
		Message:  "查询成功",
	}, nil
}

// OrderDetail 订单详情接口
func (s *OrderManageServiceImpl) OrderDetail(ctx context.Context, req *order_manage.OrderDetailReq) (*order_manage.OrderDetailResp, error) {
	log.Printf("查询订单详情: 订单ID=%d, 用户ID=%d", req.OrderId, req.UserId)

	// 1. 查询订单（调用DAO层，带权限验证）
	var orderData model.OrderManage
	order, err := dao.GetOrderManageByIDAndUserID(&orderData, req.OrderId, req.UserId)
	if err != nil {
		return &order_manage.OrderDetailResp{
			Success: false,
			Message: "订单不存在或无权限查看",
		}, nil
	}

	// 2. 查询状态历史（调用DAO层）
	var historyData model.OrderStatusHistory
	histories, err := dao.GetOrderStatusHistoryByOrderID(&historyData, req.OrderId)
	if err != nil {
		log.Printf("查询订单状态历史失败: %v", err)
		// 不影响主流程，继续返回订单详情
	}

	var statusHistory []*order_manage.OrderStatusHistory
	for _, h := range histories {
		history := &order_manage.OrderStatusHistory{
			Status:    h.Status,
			ChangedAt: h.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if h.Reason != "" {
			history.Reason = &h.Reason
		}
		if h.OperatorID != nil {
			history.OperatorId = h.OperatorID
			history.OperatorName = &h.OperatorName
		}
		statusHistory = append(statusHistory, history)
	}

	// 3. 构建订单详情
	detail := &order_manage.OrderDetail{
		OrderId:       int64(order.ID),
		OrderNo:       order.OrderNo,
		OrderType:     order.OrderType,
		PassengerId:   order.PassengerID,
		PassengerName: order.PassengerName,
		PassengerTel:  order.PassengerTel,
		StartPoint:    order.StartPoint,
		EndPoint:      order.EndPoint,
		DepartureTime: order.DepartureTime.Format("2006-01-02 15:04:05"),
		Status:        order.Status,
		Price:         order.Price,
		PaymentStatus: order.PaymentStatus,
		PaymentMethod: order.PaymentMethod,
		CreatedAt:     order.CreatedAt.Format("2006-01-02 15:04:05"),
		StatusHistory: statusHistory,
	}

	// 可选字段
	if order.AcceptedAt != nil {
		acceptedAt := order.AcceptedAt.Format("2006-01-02 15:04:05")
		detail.AcceptedAt = &acceptedAt
	}
	if order.CompletedAt != nil {
		completedAt := order.CompletedAt.Format("2006-01-02 15:04:05")
		detail.CompletedAt = &completedAt
	}
	if order.CancelledAt != nil {
		cancelledAt := order.CancelledAt.Format("2006-01-02 15:04:05")
		detail.CancelledAt = &cancelledAt
	}
	if order.DriverID != nil {
		detail.DriverId = order.DriverID
		detail.DriverName = &order.DriverName
		detail.DriverTel = &order.DriverTel
	}
	if order.VehicleInfo != "" {
		detail.VehicleInfo = &order.VehicleInfo
	}
	if order.RouteID != "" {
		detail.RouteId = &order.RouteID
	}
	if order.SpecialNeeds != "" {
		detail.SpecialNeeds = &order.SpecialNeeds
	}
	seatCount := int32(order.SeatCount)
	detail.SeatCount = &seatCount

	log.Printf("订单详情查询成功: 订单号=%s", order.OrderNo)

	return &order_manage.OrderDetailResp{
		OrderDetail: detail,
		Success:     true,
		Message:     "查询成功",
	}, nil
}

// OrderStatusUpdate 订单状态更新接口
func (s *OrderManageServiceImpl) OrderStatusUpdate(ctx context.Context, req *order_manage.OrderStatusUpdateReq) (*order_manage.OrderStatusUpdateResp, error) {
	log.Printf("订单状态更新请求: 订单ID=%d, 操作人ID=%d, 新状态=%s", req.OrderId, req.OperatorId, req.NewStatus_)

	// 1. 查询订单（调用DAO层）
	var orderData model.OrderManage
	order, err := dao.GetOrderManageByID(&orderData, req.OrderId)
	if err != nil {
		return &order_manage.OrderStatusUpdateResp{
			Success: false,
			Message: "订单不存在",
		}, nil
	}

	// 2. 验证状态转换
	if !isValidStatusTransition(order.Status, req.NewStatus_) {
		return &order_manage.OrderStatusUpdateResp{
			Success: false,
			Message: fmt.Sprintf("不能从 %s 状态转换到 %s 状态", order.Status, req.NewStatus_),
		}, nil
	}

	// 3. 取消订单时必须提供原因
	if req.NewStatus_ == "cancelled" && (req.Reason == nil || *req.Reason == "") {
		return &order_manage.OrderStatusUpdateResp{
			Success: false,
			Message: "取消订单必须提供原因",
		}, nil
	}

	// 4. 更新订单状态（调用DAO层）
	now := time.Now()
	updates := map[string]interface{}{
		"status":     req.NewStatus_,
		"updated_at": now,
	}

	// 根据状态更新对应的时间字段
	switch req.NewStatus_ {
	case "accepted":
		updates["accepted_at"] = now
		// 如果是司机接单，更新司机信息
		if req.OperatorType == "driver" {
			updates["driver_id"] = req.OperatorId
			// 关联查询司机信息
			driverName, driverTel, err := dao.GetDriverInfoByID(req.OperatorId)
			if err != nil {
				log.Printf("查询司机信息失败: DriverID=%d, Error=%v", req.OperatorId, err)
				updates["driver_name"] = fmt.Sprintf("司机%d", req.OperatorId) // 查询失败时使用默认值
				updates["driver_tel"] = ""
			} else {
				updates["driver_name"] = driverName
				updates["driver_tel"] = driverTel
			}
		}
	case "completed":
		updates["completed_at"] = now
	case "cancelled":
		updates["cancelled_at"] = now
	}

	if err := dao.UpdateOrderManage(&orderData, req.OrderId, updates); err != nil {
		log.Printf("更新订单状态失败: %v", err)
		return &order_manage.OrderStatusUpdateResp{
			Success: false,
			Message: "更新订单状态失败",
		}, nil
	}

	// 5. 记录状态历史（调用DAO层）
	reason := ""
	if req.Reason != nil {
		reason = *req.Reason
	}
	history := model.OrderStatusHistory{
		OrderID:      req.OrderId,
		Status:       req.NewStatus_,
		Reason:       reason,
		OperatorID:   &req.OperatorId,
		OperatorType: req.OperatorType,
	}
	err = dao.CreateOrderStatusHistory(&history)
	if err != nil {
		return nil, err
	}

	log.Printf("订单状态更新成功: 订单ID=%d, 新状态=%s", req.OrderId, req.NewStatus_)

	return &order_manage.OrderStatusUpdateResp{
		Success: true,
		Message: "订单状态更新成功",
	}, nil
}

// generateOrderNo 生成订单号
func generateOrderNo() string {
	return fmt.Sprintf("ORD%d", time.Now().UnixNano()/1000)
}

// isValidStatusTransition 验证状态转换是否合法
func isValidStatusTransition(oldStatus, newStatus string) bool {
	// 定义合法的状态转换
	validTransitions := map[string][]string{
		"pending":     {"accepted", "cancelled"},
		"accepted":    {"in_progress", "cancelled"},
		"in_progress": {"completed", "cancelled"},
		"completed":   {},
		"cancelled":   {},
	}

	allowedStatuses, exists := validTransitions[oldStatus]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}

	return false
}
