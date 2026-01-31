package main

import (
	"context"
	"group/handler/dao"
	"group/handler/model"
	"log"

	"gorm.io/gorm"
)

// RatingServiceImpl 评价服务实现
type RatingServiceImpl struct{}

// RatingSubmit 评价提交接口
func (s *RatingServiceImpl) RatingSubmit(ctx context.Context, req *RatingSubmitReq) (*RatingSubmitResp, error) {
	log.Printf("评价提交请求: 订单ID=%d, 用户ID=%d, 评分=%.2f", req.OrderId, req.UserId, req.Score)

	// 1. 查询订单信息
	orderData, err := dao.GetOrderByID(req.OrderId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &RatingSubmitResp{
				Success: false,
				Message: "订单不存在",
			}, nil
		}
		log.Printf("查询订单失败: %v", err)
		return &RatingSubmitResp{
			Success: false,
			Message: "查询订单失败",
		}, nil
	}

	// 2. 验证用户权限
	if orderData.PassengerID != req.UserId && (orderData.DriverID == nil || *orderData.DriverID != req.UserId) {
		return &RatingSubmitResp{
			Success: false,
			Message: "无权评价此订单",
		}, nil
	}

	// 3. 检查订单状态（只有已完成的订单才能评价）
	if orderData.Status != "completed" {
		return &RatingSubmitResp{
			Success: false,
			Message: "订单未完成，无法评价",
		}, nil
	}

	// 4. 检查是否已评价
	_, err = dao.GetRatingByOrderID(req.OrderId)
	if err == nil {
		return &RatingSubmitResp{
			Success: false,
			Message: "该订单已评价，请勿重复评价",
		}, nil
	}

	// 5. 创建评价记录
	ratingData := &model.RatingManage{
		OrderID:   req.OrderId,
		UserID:    req.UserId,
		UserName:  orderData.PassengerName,
		Score:     req.Score,
		Tags:      dao.TagsToJSON(req.Tags),
		Comment:   req.Comment,
		RaterType: req.RaterType,
	}

	// 如果是乘客评价司机
	if req.RaterType == "passenger" && orderData.DriverID != nil {
		ratingData.DriverID = orderData.DriverID
		ratingData.DriverName = orderData.DriverName
	}

	err = dao.CreateRating(ratingData)
	if err != nil {
		log.Printf("创建评价记录失败: %v", err)
		return &RatingSubmitResp{
			Success: false,
			Message: "创建评价记录失败",
		}, nil
	}

	log.Printf("评价提交成功: 评价ID=%d", ratingData.ID)

	return &RatingSubmitResp{
		RatingId: int64(ratingData.ID),
		Success:  true,
		Message:  "评价提交成功",
	}, nil
}

// RatingQuery 评价查询接口
func (s *RatingServiceImpl) RatingQuery(ctx context.Context, req *RatingQueryReq) (*RatingQueryResp, error) {
	log.Printf("评价查询请求: 时间范围=%v, 线路ID=%v, 司机ID=%v",
		req.TimeRange, req.RouteId, req.DriverId)

	// 分页参数
	page := int32(1)
	pageSize := int32(10)
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.PageSize != nil && *req.PageSize > 0 {
		pageSize = *req.PageSize
	}

	// 查询评价记录
	ratingList, total, err := dao.QueryRatings(
		req.TimeRange,
		req.RouteId,
		req.DriverId,
		req.DriverPhone,
		page,
		pageSize,
	)
	if err != nil {
		log.Printf("查询评价记录失败: %v", err)
		return &RatingQueryResp{
			Success: false,
			Message: "查询失败",
		}, nil
	}

	// 构建返回数据
	var ratings []*RatingInfo
	for _, rating := range ratingList {
		info := &RatingInfo{
			RatingId:  int64(rating.ID),
			OrderId:   rating.OrderID,
			UserId:    rating.UserID,
			UserName:  rating.UserName,
			Score:     rating.Score,
			Tags:      dao.JSONToTags(rating.Tags),
			Comment:   rating.Comment,
			RaterType: rating.RaterType,
			CreatedAt: rating.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if rating.DriverID != nil {
			info.DriverId = rating.DriverID
			info.DriverName = &rating.DriverName
		}

		ratings = append(ratings, info)
	}

	return &RatingQueryResp{
		Ratings:  ratings,
		Total:    int32(total),
		Page:     page,
		PageSize: pageSize,
		Success:  true,
		Message:  "查询成功",
	}, nil
}

// OrderHistoryQuery 订单历史变更查询接口
func (s *RatingServiceImpl) OrderHistoryQuery(ctx context.Context, req *OrderHistoryQueryReq) (*OrderHistoryQueryResp, error) {
	log.Printf("订单历史变更查询请求: 订单ID=%d", req.OrderId)

	// 查询订单历史记录
	histories, err := dao.GetOrderHistoryByOrderID(req.OrderId)
	if err != nil {
		log.Printf("查询订单历史记录失败: %v", err)
		return &OrderHistoryQueryResp{
			Success: false,
			Message: "查询失败",
		}, nil
	}

	// 构建返回数据
	var records []*OrderHistoryRecord
	for _, history := range histories {
		record := &OrderHistoryRecord{
			HistoryId: int64(history.ID),
			OrderId:   history.OrderID,
			Status:    history.Status,
			Reason:    history.Reason,
			ChangedAt: history.ChangedAt.Format("2006-01-02 15:04:05"),
		}

		if history.OperatorID != nil {
			record.OperatorId = history.OperatorID
			record.OperatorName = &history.OperatorName
		}

		records = append(records, record)
	}

	return &OrderHistoryQueryResp{
		Records: records,
		Success: true,
		Message: "查询成功",
	}, nil
}

// TicketRecordQuery 购票记录查询接口
func (s *RatingServiceImpl) TicketRecordQuery(ctx context.Context, req *TicketRecordQueryReq) (*TicketRecordQueryResp, error) {
	log.Printf("购票记录查询请求: 用户ID=%d, 时间范围=%v", req.UserId, req.TimeRange)

	// 分页参数
	page := int32(1)
	pageSize := int32(10)
	if req.Page != nil && *req.Page > 0 {
		page = *req.Page
	}
	if req.PageSize != nil && *req.PageSize > 0 {
		pageSize = *req.PageSize
	}

	// 查询购票记录
	orderList, total, err := dao.QueryTicketRecords(
		req.UserId,
		req.TimeRange,
		page,
		pageSize,
	)
	if err != nil {
		log.Printf("查询购票记录失败: %v", err)
		return &TicketRecordQueryResp{
			Success: false,
			Message: "查询失败",
		}, nil
	}

	// 构建返回数据
	var records []*TicketRecordInfo
	for _, order := range orderList {
		record := &TicketRecordInfo{
			OrderId:       int64(order.ID),
			OrderNo:       order.OrderNo,
			UserId:        order.PassengerID,
			UserName:      order.PassengerName,
			StartPoint:    order.StartPoint,
			EndPoint:      order.EndPoint,
			DepartureTime: order.DepartureTime.Format("2006-01-02 15:04:05"),
			Price:         order.Price,
			Status:        order.Status,
			CreatedAt:     order.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		records = append(records, record)
	}

	return &TicketRecordQueryResp{
		Records:  records,
		Total:    int32(total),
		Page:     page,
		PageSize: pageSize,
		Success:  true,
		Message:  "查询成功",
	}, nil
}

// 临时定义结构体（实际应该从生成的kitex代码中导入）
type RatingSubmitReq struct {
	OrderId   int64
	UserId    int64
	Score     float64
	Tags      []string
	Comment   string
	RaterType string
}

type RatingSubmitResp struct {
	RatingId int64
	Message  string
	Success  bool
}

type RatingQueryReq struct {
	TimeRange   *string
	RouteId     *int64
	DriverId    *int64
	DriverPhone *string
	Page        *int32
	PageSize    *int32
}

type RatingInfo struct {
	RatingId   int64
	OrderId    int64
	UserId     int64
	UserName   string
	Score      float64
	Tags       []string
	Comment    string
	RaterType  string
	CreatedAt  string
	DriverId   *int64
	DriverName *string
}

type RatingQueryResp struct {
	Ratings  []*RatingInfo
	Total    int32
	Page     int32
	PageSize int32
	Message  string
	Success  bool
}

type OrderHistoryQueryReq struct {
	OrderId int64
}

type OrderHistoryRecord struct {
	HistoryId    int64
	OrderId      int64
	Status       string
	Reason       string
	ChangedAt    string
	OperatorId   *int64
	OperatorName *string
}

type OrderHistoryQueryResp struct {
	Records []*OrderHistoryRecord
	Message string
	Success bool
}

type TicketRecordQueryReq struct {
	UserId    int64
	TimeRange *string
	Page      *int32
	PageSize  *int32
}

type TicketRecordInfo struct {
	OrderId       int64
	OrderNo       string
	UserId        int64
	UserName      string
	StartPoint    string
	EndPoint      string
	DepartureTime string
	Price         float64
	Status        string
	CreatedAt     string
}

type TicketRecordQueryResp struct {
	Records  []*TicketRecordInfo
	Total    int32
	Page     int32
	PageSize int32
	Message  string
	Success  bool
}
