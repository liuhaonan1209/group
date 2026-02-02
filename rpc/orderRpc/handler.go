package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"group/handler/dao"
	"group/handler/model"
	order "group/kitex_gen/car/order"
)

// OrderServiceImpl 订单服务实现结构体
type OrderServiceImpl struct{}

// TripPublish 行程发布接口（乘客发布行程）
func (s *OrderServiceImpl) TripPublish(ctx context.Context, req *order.TripPublishReq) (resp *order.TripPublishResp, err error) {
	log.Printf("乘客发布行程: 乘客ID=%d, 起点=%s, 终点=%s", req.PassengerId, req.StartPoint, req.EndPoint)

	// 解析出行时间
	departureTime, err := time.Parse("2006-01-02 15:04:05", req.DepartureTime)
	if err != nil {
		return &order.TripPublishResp{
			TripId:  0,
			Message: "出行时间格式不正确，请使用格式：2006-01-02 15:04:05",
			Success: false,
		}, nil
	}

	// 判断出行时间类型（即时或预约）
	tripTimeType := "scheduled" // 默认预约
	if time.Now().Add(30 * time.Minute).After(departureTime) {
		// 如果出行时间在30分钟内，视为即时出行
		tripTimeType = "immediate"
	}

	// 处理特殊需求详情（如果有）
	specialNeedsDetail := ""
	if req.SpecialNeeds != nil && *req.SpecialNeeds != "" {
		// 这里可以解析特殊需求并转换为JSON格式
		// 暂时保持原样，后续可以通过API传入结构化数据
		specialNeedsDetail = ""
	}

	// 创建行程记录
	trip := model.Trip{
		PublisherID:        uint(req.PassengerId),
		PublisherType:      "passenger",
		StartPoint:         req.StartPoint,
		EndPoint:           req.EndPoint,
		DepartureTime:      departureTime,
		TripTimeType:       tripTimeType,
		SpecialNeeds:       getStringValue(req.SpecialNeeds),
		SpecialNeedsDetail: specialNeedsDetail,
		ContactWay:         getStringValue(req.ContactWay),
		Status:             "pending",
	}

	err = dao.CreateTrip(&trip)
	if err != nil {
		log.Printf("创建行程失败: %v", err)
		return &order.TripPublishResp{
			TripId:  0,
			Message: "发布行程失败",
			Success: false,
		}, nil
	}

	log.Printf("行程发布成功: 行程ID=%d", trip.ID)
	return &order.TripPublishResp{
		TripId:  int64(trip.ID),
		Message: "行程发布成功",
		Success: true,
	}, nil
}

// DriverTripPublish 司机行程发布接口
func (s *OrderServiceImpl) DriverTripPublish(ctx context.Context, req *order.DriverTripPublishReq) (resp *order.DriverTripPublishResp, err error) {
	log.Printf("司机发布行程: 司机ID=%d, 起点=%s, 终点=%s", req.DriverId, req.StartPoint, req.EndPoint)
	// 解析发车时间
	departureTime, err := time.Parse("2006-01-02 15:04:05", req.DepartureTime)
	if err != nil {
		return &order.DriverTripPublishResp{
			TripId:  0,
			Message: "发车时间格式不正确，请使用格式：2006-01-02 15:04:05",
			Success: false,
		}, nil
	}

	// 判断出行时间类型（即时或预约）
	tripTimeType := "scheduled" // 默认预约
	if time.Now().Add(30 * time.Minute).After(departureTime) {
		// 如果出行时间在30分钟内，视为即时出行
		tripTimeType = "immediate"
	}

	// 创建行程记录
	trip := model.Trip{
		PublisherID:   uint(req.DriverId),
		PublisherType: "driver",
		StartPoint:    req.StartPoint,
		EndPoint:      req.EndPoint,
		DepartureTime: departureTime,
		TripTimeType:  tripTimeType,
		VehicleInfo:   req.VehicleInfo,
		Status:        "pending",
	}

	err = dao.CreateTrip(&trip)
	if err != nil {
		log.Printf("创建司机行程失败: %v", err)
		return &order.DriverTripPublishResp{
			TripId:  0,
			Message: "发布行程失败",
			Success: false,
		}, nil
	}

	log.Printf("司机行程发布成功: 行程ID=%d", trip.ID)
	return &order.DriverTripPublishResp{
		TripId:  int64(trip.ID),
		Message: "行程发布成功",
		Success: true,
	}, nil
}

// TripQuery 行程查询接口
func (s *OrderServiceImpl) TripQuery(ctx context.Context, req *order.TripQueryReq) (resp *order.TripQueryResp, err error) {
	log.Printf("查询行程: 起点=%s, 终点=%s, 时间=%s", req.StartPoint, req.EndPoint, req.DepartureTime)

	// 解析出行时间
	departureTime, err := time.Parse("2006-01-02", req.DepartureTime)
	if err != nil {
		return &order.TripQueryResp{
			Trips:   []*order.TripInfo{},
			Message: "出行时间格式不正确，请使用格式：2006-01-02",
			Success: false,
		}, nil
	}

	// 查询行程
	var trip model.Trip
	tripType := getStringValue(req.TripType)
	trips, err := dao.QueryTrips(&trip, req.StartPoint, req.EndPoint, departureTime, tripType)
	if err != nil {
		log.Printf("查询行程失败: %v", err)
		return &order.TripQueryResp{
			Trips:   []*order.TripInfo{},
			Message: "查询行程失败",
			Success: false,
		}, nil
	}

	// 构造响应
	tripInfos := make([]*order.TripInfo, 0, len(trips))
	for _, t := range trips {
		// 关联查询发布者姓名
		publisherName, err := dao.GetPublisherNameByID(t.PublisherID, t.PublisherType)
		if err != nil {
			log.Printf("查询发布者姓名失败: PublisherID=%d, Type=%s, Error=%v", t.PublisherID, t.PublisherType, err)
			publisherName = fmt.Sprintf("用户%d", t.PublisherID) // 查询失败时使用默认值
		}

		tripInfo := &order.TripInfo{
			TripId:        int64(t.ID),
			StartPoint:    t.StartPoint,
			EndPoint:      t.EndPoint,
			DepartureTime: t.DepartureTime.Format("2006-01-02 15:04:05"),
			TripType:      t.PublisherType,
			PublisherName: publisherName,
			Status:        t.Status,
		}

		if t.VehicleInfo != "" {
			tripInfo.VehicleInfo = &t.VehicleInfo
		}
		if t.SpecialNeeds != "" {
			tripInfo.SpecialNeeds = &t.SpecialNeeds
		}

		tripInfos = append(tripInfos, tripInfo)
	}

	log.Printf("查询到 %d 条行程", len(tripInfos))
	return &order.TripQueryResp{
		Trips:   tripInfos,
		Message: fmt.Sprintf("查询成功，共找到%d条行程", len(tripInfos)),
		Success: true,
	}, nil
}

// PassengerHelp 乘客求助接口
func (s *OrderServiceImpl) PassengerHelp(ctx context.Context, req *order.PassengerHelpReq) (resp *order.PassengerHelpResp, err error) {
	log.Printf("乘客求助: 行程ID=%d, 乘客ID=%d, 求助类型=%s", req.TripId, req.PassengerId, req.HelpType)

	// 验证行程是否存在
	var trip model.Trip
	_, err = dao.GetTripByID(&trip, req.TripId)
	if err != nil {
		return &order.PassengerHelpResp{
			ContactInfo: "",
			Message:     "行程不存在",
			Success:     false,
		}, nil
	}

	// 创建求助记录
	help := model.PassengerHelp{
		TripID:      uint(req.TripId),
		PassengerID: uint(req.PassengerId),
		HelpType:    req.HelpType,
		Status:      "pending",
	}

	// 根据求助类型提供不同的联系方式
	var contactInfo string
	if strings.Contains(req.HelpType, "官方") {
		contactInfo = "官方客服电话：400-123-4567，工作时间：9:00-21:00"
		help.ContactInfo = contactInfo
	} else if strings.Contains(req.HelpType, "客服") {
		contactInfo = "在线客服已接入，请稍候..."
		help.ContactInfo = contactInfo
	} else {
		contactInfo = "求助已提交，客服将尽快联系您"
		help.ContactInfo = contactInfo
	}

	err = dao.CreatePassengerHelp(&help)
	if err != nil {
		log.Printf("创建求助记录失败: %v", err)
		return &order.PassengerHelpResp{
			ContactInfo: "",
			Message:     "提交求助失败",
			Success:     false,
		}, nil
	}

	log.Printf("求助记录创建成功: ID=%d", help.ID)
	return &order.PassengerHelpResp{
		ContactInfo: contactInfo,
		Message:     "求助已提交",
		Success:     true,
	}, nil
}

// TripShare 行程分享接口
func (s *OrderServiceImpl) TripShare(ctx context.Context, req *order.TripShareReq) (resp *order.TripShareResp, err error) {
	log.Printf("行程分享: 行程ID=%d, 用户ID=%d", req.TripId, req.UserId)

	// 验证行程是否存在
	var trip model.Trip
	tripDetail, err := dao.GetTripByID(&trip, req.TripId)
	if err != nil {
		return &order.TripShareResp{
			ShareLink: "",
			Message:   "行程不存在",
			Success:   false,
		}, nil
	}

	// 生成分享链接
	shareLink := fmt.Sprintf("https://carpool.example.com/trip/%d?share=%d", req.TripId, req.UserId)

	// 创建分享记录
	share := model.TripShare{
		TripID:       uint(req.TripId),
		UserID:       uint(req.UserId),
		ShareTargets: strings.Join(req.ShareTargets, ","),
		ShareLink:    shareLink,
		SharedAt:     time.Now(),
	}

	err = dao.CreateTripShare(&share)
	if err != nil {
		log.Printf("创建分享记录失败: %v", err)
		return &order.TripShareResp{
			ShareLink: "",
			Message:   "分享失败",
			Success:   false,
		}, nil
	}

	shareInfo := fmt.Sprintf("行程分享：%s -> %s，出发时间：%s。点击查看详情：%s",
		tripDetail.StartPoint, tripDetail.EndPoint,
		tripDetail.DepartureTime.Format("2006-01-02 15:04"),
		shareLink)

	log.Printf("行程分享成功: ID=%d", share.ID)
	return &order.TripShareResp{
		ShareLink: shareInfo,
		Message:   "分享成功",
		Success:   true,
	}, nil
}

// TripDetail 行程详情接口
func (s *OrderServiceImpl) TripDetail(ctx context.Context, req *order.TripDetailReq) (resp *order.TripDetailResp, err error) {
	log.Printf("查询行程详情: 行程ID=%d", req.TripId)
	// 查询行程详情
	var trip model.Trip
	tripDetail, err := dao.GetTripByID(&trip, req.TripId)
	if err != nil {
		log.Printf("查询行程详情失败: %v", err)
		return &order.TripDetailResp{
			TripInfo: nil,
			Message:  "行程不存在",
			Success:  false,
		}, nil
	}

	// 构造响应
	tripInfo := &order.TripInfo{
		TripId:        int64(tripDetail.ID),
		StartPoint:    tripDetail.StartPoint,
		EndPoint:      tripDetail.EndPoint,
		DepartureTime: tripDetail.DepartureTime.Format("2006-01-02 15:04:05"),
		TripType:      tripDetail.PublisherType,
		Status:        tripDetail.Status,
	}

	// 关联查询发布者姓名
	publisherName, err := dao.GetPublisherNameByID(tripDetail.PublisherID, tripDetail.PublisherType)
	if err != nil {
		log.Printf("查询发布者姓名失败: PublisherID=%d, Type=%s, Error=%v", tripDetail.PublisherID, tripDetail.PublisherType, err)
		tripInfo.PublisherName = fmt.Sprintf("用户%d", tripDetail.PublisherID) // 查询失败时使用默认值
	} else {
		tripInfo.PublisherName = publisherName
	}

	if tripDetail.VehicleInfo != "" {
		tripInfo.VehicleInfo = &tripDetail.VehicleInfo
	}
	if tripDetail.SpecialNeeds != "" {
		tripInfo.SpecialNeeds = &tripDetail.SpecialNeeds
	}

	log.Printf("查询行程详情成功: 行程ID=%d", req.TripId)
	return &order.TripDetailResp{
		TripInfo: tripInfo,
		Message:  "查询成功",
		Success:  true,
	}, nil
}

// getStringValue 获取可选字符串的值
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// DataExport 数据导出接口
func (s *OrderServiceImpl) DataExport(ctx context.Context, req *order.DataExportReq) (resp *order.DataExportResp, err error) {
	log.Printf("数据导出请求: 用户ID=%d, 用户类型=%s", req.UserId, req.UserType)
	// 验证用户类型
	if req.UserType != "passenger" && req.UserType != "driver" {
		return &order.DataExportResp{
			DownloadUrl: "",
			FileFormat:  "",
			ExportId:    0,
			Message:     "用户类型必须是 passenger 或 driver",
			Success:     false,
		}, nil
	}

	// 创建导出记录
	exportRecord := model.ExportRecord{
		UserID:       uint(req.UserId),
		UserType:     req.UserType,
		ExportFields: strings.Join(req.ExportFields, ","),
		Status:       "pending",
		FileFormat:   "csv",
		ExpiresAt:    time.Now().AddDate(0, 0, 7), // 7天后过期
	}

	err = dao.CreateExportRecord(&exportRecord)
	if err != nil {
		log.Printf("创建导出记录失败: %v", err)
		return &order.DataExportResp{
			DownloadUrl: "",
			FileFormat:  "",
			ExportId:    0,
			Message:     "创建导出任务失败",
			Success:     false,
		}, nil
	}

	// 生成下载链接（实际应该是异步处理导出任务）
	downloadUrl := fmt.Sprintf("https://carpool.example.com/download/export_%d.csv", exportRecord.ID)

	// 更新导出记录状态
	err = dao.UpdateExportRecordStatus(&exportRecord, int64(exportRecord.ID), "completed", downloadUrl)
	if err != nil {
		return nil, err
	}

	log.Printf("数据导出任务创建成功: 导出ID=%d", exportRecord.ID)
	return &order.DataExportResp{
		DownloadUrl: downloadUrl,
		FileFormat:  "csv",
		ExportId:    int64(exportRecord.ID),
		Message:     "导出任务已创建，数据正在生成中",
		Success:     true,
	}, nil
}

// ExportRecordQuery 导出记录查询接口
func (s *OrderServiceImpl) ExportRecordQuery(ctx context.Context, req *order.ExportRecordQueryReq) (resp *order.ExportRecordQueryResp, err error) {
	log.Printf("查询导出记录: 用户ID=%d", req.UserId)
	// 查询导出记录
	var record model.ExportRecord
	limit := 10 // 默认查询10条
	if req.Limit != nil && *req.Limit > 0 {
		limit = int(*req.Limit)
	}

	records, err := dao.GetExportRecordsByUserID(&record, req.UserId, limit)
	if err != nil {
		log.Printf("查询导出记录失败: %v", err)
		return &order.ExportRecordQueryResp{
			Records: []*order.ExportRecordInfo{},
			Message: "查询导出记录失败",
			Success: false,
		}, nil
	}

	// 构造响应
	recordInfos := make([]*order.ExportRecordInfo, 0, len(records))
	for _, r := range records {
		recordInfo := &order.ExportRecordInfo{
			ExportId:     int64(r.ID),
			UserId:       int64(r.UserID),
			UserType:     r.UserType,
			ExportFields: r.ExportFields,
			Status:       r.Status,
			DownloadUrl:  r.DownloadUrl,
			CreatedAt:    r.CreatedAt.Format("2006-01-02 15:04:05"),
			ExpiresAt:    r.ExpiresAt.Format("2006-01-02 15:04:05"),
		}
		recordInfos = append(recordInfos, recordInfo)
	}

	log.Printf("查询到 %d 条导出记录", len(recordInfos))
	return &order.ExportRecordQueryResp{
		Records: recordInfos,
		Message: fmt.Sprintf("查询成功，共找到%d条记录", len(recordInfos)),
		Success: true,
	}, nil
}
