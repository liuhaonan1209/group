package dao

import (
	"fmt"
	"group/global"
	"time"
)

// ==================== Trip 相关操作 ====================

// CreateTrip 创建行程
func CreateTrip[T any](data *T) error {
	return global.DB.Create(data).Error
}

// GetTripByID 根据ID查询行程
func GetTripByID[T any](data *T, id int64) (trip T, err error) {
	err = global.DB.Model(data).Where("id = ?", id).First(&trip).Error
	return
}

// QueryTrips 查询行程列表
func QueryTrips[T any](data *T, startPoint, endPoint string, departureTime time.Time, tripType string) (trips []T, err error) {
	query := global.DB.Model(data).Where("start_point = ? AND end_point = ?", startPoint, endPoint)

	// 查询当天的行程
	startOfDay := time.Date(departureTime.Year(), departureTime.Month(), departureTime.Day(), 0, 0, 0, 0, departureTime.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	query = query.Where("departure_time >= ? AND departure_time < ?", startOfDay, endOfDay)

	// 如果指定了行程类型
	if tripType != "" {
		query = query.Where("publisher_type = ?", tripType)
	}

	// 只查询待确认和已确认的行程
	query = query.Where("status IN ?", []string{"pending", "confirmed"})

	err = query.Order("departure_time ASC").Find(&trips).Error
	return
}

// UpdateTripStatus 更新行程状态
func UpdateTripStatus[T any](data *T, id int64, status string) error {
	return global.DB.Model(data).Where("id = ?", id).Update("status", status).Error
}

// GetTripsByPublisher 根据发布者查询行程
func GetTripsByPublisher[T any](data *T, publisherId int64, publisherType string) (trips []T, err error) {
	err = global.DB.Model(data).Where("publisher_id = ? AND publisher_type = ?", publisherId, publisherType).
		Order("departure_time DESC").Find(&trips).Error
	return
}

// ==================== TripShare 相关操作 ====================

// CreateTripShare 创建行程分享记录
func CreateTripShare[T any](data *T) error {
	return global.DB.Create(data).Error
}

// GetTripSharesByTripID 根据行程ID查询分享记录
func GetTripSharesByTripID[T any](data *T, tripId int64) (shares []T, err error) {
	err = global.DB.Model(data).Where("trip_id = ?", tripId).Find(&shares).Error
	return
}

// ==================== PassengerHelp 相关操作 ====================

// CreatePassengerHelp 创建乘客求助记录
func CreatePassengerHelp[T any](data *T) error {
	return global.DB.Create(data).Error
}

// GetPassengerHelpByID 根据ID查询求助记录
func GetPassengerHelpByID[T any](data *T, id int64) (help T, err error) {
	err = global.DB.Model(data).Where("id = ?", id).First(&help).Error
	return
}

// UpdatePassengerHelpStatus 更新求助状态
func UpdatePassengerHelpStatus[T any](data *T, id int64, status, contactInfo string) error {
	return global.DB.Model(data).Where("id = ?", id).
		Updates(map[string]interface{}{"status": status, "contact_info": contactInfo}).Error
}

// GetPassengerHelpsByPassenger 根据乘客ID查询求助记录
func GetPassengerHelpsByPassenger[T any](data *T, passengerId int64) (helps []T, err error) {
	err = global.DB.Model(data).Where("passenger_id = ?", passengerId).
		Order("created_at DESC").Find(&helps).Error
	return
}

// ==================== ExportRecord 相关操作 ====================

// CreateExportRecord 创建导出记录
func CreateExportRecord[T any](data *T) error {
	return global.DB.Create(data).Error
}

// GetExportRecordByID 根据ID查询导出记录
func GetExportRecordByID[T any](data *T, id int64) (record T, err error) {
	err = global.DB.Model(data).Where("id = ?", id).First(&record).Error
	return
}

// GetExportRecordsByUserID 根据用户ID查询导出记录
func GetExportRecordsByUserID[T any](data *T, userId int64, limit int) (records []T, err error) {
	query := global.DB.Model(data).Where("user_id = ?", userId).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err = query.Find(&records).Error
	return
}

// UpdateExportRecordStatus 更新导出记录状态
func UpdateExportRecordStatus[T any](data *T, id int64, status, downloadUrl string) error {
	return global.DB.Model(data).Where("id = ?", id).
		Updates(map[string]interface{}{"status": status, "download_url": downloadUrl}).Error
}

// ==================== OrderManage 订单管理相关操作 ====================

// CreateOrder 创建订单
func CreateOrder[T any](data *T) error {
	return global.DB.Create(data).Error
}

// GetOrderManageByID 根据订单ID查询订单
func GetOrderManageByID[T any](data *T, id int64) (order T, err error) {
	err = global.DB.Model(data).Where("id = ?", id).First(&order).Error
	return
}

// GetOrderManageByIDAndUserID 根据订单ID和用户ID查询订单（带权限验证）
func GetOrderManageByIDAndUserID[T any](data *T, orderId int64, userId int64) (order T, err error) {
	err = global.DB.Model(data).Where("id = ? AND (passenger_id = ? OR driver_id = ?)", orderId, userId, userId).
		First(&order).Error
	return
}

// QueryOrders 多条件查询订单列表（支持分页）
func QueryOrders[T any](data *T, orderType *string, status *string, startTime *string, endTime *string,
	routeId *string, passengerTel *string, passengerId *int64, driverId *int64, page, pageSize int32) (orders []T, total int64, err error) {
	query := global.DB.Model(data)

	// 动态构建查询条件
	if orderType != nil {
		query = query.Where("order_type = ?", *orderType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}
	if routeId != nil {
		query = query.Where("route_id = ?", *routeId)
	}
	if passengerTel != nil {
		query = query.Where("passenger_tel = ?", *passengerTel)
	}
	if passengerId != nil {
		query = query.Where("passenger_id = ?", *passengerId)
	}
	if driverId != nil {
		query = query.Where("driver_id = ?", *driverId)
	}

	// 获取总数
	err = query.Count(&total).Error
	if err != nil {
		return
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").Limit(int(pageSize)).Offset(int(offset)).Find(&orders).Error
	return
}

// UpdateOrderManage 更新订单信息
func UpdateOrderManage[T any](data *T, id int64, updates map[string]interface{}) error {
	return global.DB.Model(data).Where("id = ?", id).Updates(updates).Error
}

// ==================== OrderStatusHistory 订单状态历史相关操作 ====================

// CreateOrderStatusHistory 创建订单状态历史记录
func CreateOrderStatusHistory[T any](data *T) error {
	return global.DB.Create(data).Error
}

// GetOrderStatusHistoryByOrderID 根据订单ID查询状态历史
func GetOrderStatusHistoryByOrderID[T any](data *T, orderId int64) (histories []T, err error) {
	err = global.DB.Model(data).Where("order_id = ?", orderId).Order("created_at ASC").Find(&histories).Error
	return
}

// ==================== 关联查询辅助函数 ====================

// GetPublisherNameByID 根据发布者ID和类型查询发布者姓名
// publisherType: "passenger" 或 "driver"
func GetPublisherNameByID(publisherID uint, publisherType string) (name string, err error) {
	if publisherType == "passenger" {
		var passenger struct {
			Name string
		}
		err = global.DB.Table("passengers").Select("name").Where("id = ?", publisherID).First(&passenger).Error
		if err != nil {
			return "", err
		}
		return passenger.Name, nil
	} else if publisherType == "driver" {
		var driver struct {
			Name string
		}
		err = global.DB.Table("drivers").Select("name").Where("id = ?", publisherID).First(&driver).Error
		if err != nil {
			return "", err
		}
		return driver.Name, nil
	}
	return "", fmt.Errorf("无效的发布者类型: %s", publisherType)
}

// GetDriverInfoByID 根据司机ID查询司机信息
func GetDriverInfoByID(driverID int64) (name string, tel string, err error) {
	var driver struct {
		Name string
		Tel  string
	}
	err = global.DB.Table("drivers").Select("name, tel").Where("id = ?", driverID).First(&driver).Error
	if err != nil {
		return "", "", err
	}
	return driver.Name, driver.Tel, nil
}