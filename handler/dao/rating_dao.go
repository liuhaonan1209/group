// Package dao 数据访问层 - 评价管理模块
// 本文件包含所有评价、订单历史、购票记录相关的数据库操作函数
// 遵循DAO设计模式，将所有数据库操作封装在此层，供上层RPC Handler调用
package dao

import (
	"encoding/json"
	"group/global"
	"group/handler/model"
)

// ==================== Rating 评价相关操作 ====================

// CreateRating 创建评价记录
// 用于乘客或司机对订单进行评价
// 参数:
//   - data: 评价数据对象
//
// 返回:
//   - error: 创建失败时返回错误信息
func CreateRating(data *model.RatingManage) error {
	return global.DB.Create(data).Error
}

// GetRatingByID 根据评价ID查询评价记录
// 参数:
//   - id: 评价记录ID
//
// 返回:
//   - rating: 评价对象
//   - err: 查询失败或记录不存在时返回错误
func GetRatingByID(id int64) (rating model.RatingManage, err error) {
	err = global.DB.Where("id = ?", id).First(&rating).Error
	return
}

// GetRatingByOrderID 根据订单ID查询评价记录
// 用于检查某订单是否已被评价（防止重复评价）
// 参数:
//   - orderId: 订单ID
//
// 返回:
//   - rating: 评价对象
//   - err: 查询失败或记录不存在时返回错误
func GetRatingByOrderID(orderId int64) (rating model.RatingManage, err error) {
	err = global.DB.Where("order_id = ?", orderId).First(&rating).Error
	return
}

// QueryRatings 多条件查询评价记录列表（支持分页）
// 参数:
//   - timeRange: 时间范围（可选，格式：2024-01-01 或 2024-01-01~2024-12-31）
//   - routeId: 线路ID（可选，预留字段）
//   - driverId: 司机ID（可选，查询该司机收到的评价）
//   - driverPhone: 司机手机号（可选，通过手机号关联查询司机）
//   - page: 页码（从1开始）
//   - pageSize: 每页记录数
//
// 返回:
//   - ratings: 评价记录列表
//   - total: 符合条件的总记录数
//   - err: 查询失败时返回错误
func QueryRatings(timeRange *string, routeId *int64, driverId *int64, driverPhone *string, page, pageSize int32) (ratings []model.RatingManage, total int64, err error) {
	query := global.DB.Model(&model.RatingManage{})

	// 动态构建查询条件
	if timeRange != nil && *timeRange != "" {
		// 解析时间范围（格式：2024-01-01~2024-12-31）
		// 简化处理，实际应该解析时间范围的起止时间
		query = query.Where("created_at >= ?", *timeRange)
	}
	if driverId != nil {
		query = query.Where("driver_id = ?", *driverId)
	}
	if driverPhone != nil && *driverPhone != "" {
		// 关联查询司机表，通过手机号查找司机ID
		var driverIDs []int64
		err = global.DB.Table("drivers").Select("id").Where("tel LIKE ?", "%"+*driverPhone+"%").Pluck("id", &driverIDs).Error
		if err != nil {
			return nil, 0, err
		}
		if len(driverIDs) > 0 {
			query = query.Where("driver_id IN ?", driverIDs)
		} else {
			// 如果没有找到匹配的司机，返回空结果
			return []model.RatingManage{}, 0, nil
		}
	}

	// 获取符合条件的总记录数
	err = query.Count(&total).Error
	if err != nil {
		return
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&ratings).Error
	return
}

// ==================== OrderHistory 订单历史相关操作 ====================

// CreateOrderHistory 创建订单历史变更记录
// 用于记录订单状态的每次变更（如：创建、接单、完成、取消等）
// 参数:
//   - data: 订单历史数据对象
//
// 返回:
//   - error: 创建失败时返回错误信息
func CreateOrderHistory(data *model.OrderHistory) error {
	return global.DB.Create(data).Error
}

// GetOrderHistoryByOrderID 根据订单ID查询历史变更记录
// 按时间顺序返回该订单的所有状态变更记录
// 参数:
//   - orderId: 订单ID
//
// 返回:
//   - histories: 订单历史记录列表
//   - err: 查询失败时返回错误
func GetOrderHistoryByOrderID(orderId int64) (histories []model.OrderHistory, err error) {
	err = global.DB.Where("order_id = ?", orderId).Order("changed_at ASC").Find(&histories).Error
	return
}

// ==================== TicketRecord 购票记录相关操作 ====================

// QueryTicketRecords 查询购票记录列表（支持分页）
// 用于查询用户的购票历史
// 参数:
//   - userId: 用户ID（必填，查询该用户的购票记录）
//   - timeRange: 时间范围（可选，查询此时间之后的记录）
//   - page: 页码（从1开始）
//   - pageSize: 每页记录数
//
// 返回:
//   - orders: 订单记录列表
//   - total: 符合条件的总记录数
//   - err: 查询失败时返回错误
func QueryTicketRecords(userId int64, timeRange *string, page, pageSize int32) (orders []model.OrderManage, total int64, err error) {
	query := global.DB.Model(&model.OrderManage{}).Where("passenger_id = ?", userId)

	// 动态构建查询条件
	if timeRange != nil && *timeRange != "" {
		// 解析时间范围，查询指定时间之后的记录
		query = query.Where("created_at >= ?", *timeRange)
	}

	// 获取符合条件的总记录数
	err = query.Count(&total).Error
	if err != nil {
		return
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&orders).Error
	return
}

// ==================== 辅助函数 ====================

// TagsToJSON 将标签数组转换为JSON字符串
// 用于将评价标签（如：["服务好", "车辆干净"]）存储到数据库
// 参数:
//   - tags: 标签字符串数组
//
// 返回:
//   - string: JSON格式的字符串
func TagsToJSON(tags []string) string {
	if len(tags) == 0 {
		return "[]"
	}
	data, _ := json.Marshal(tags)
	return string(data)
}

// JSONToTags 将JSON字符串转换为标签数组
// 用于从数据库读取评价标签并转换为数组
// 参数:
//   - jsonStr: JSON格式的字符串
//
// 返回:
//   - []string: 标签字符串数组
func JSONToTags(jsonStr string) []string {
	var tags []string
	if jsonStr == "" || jsonStr == "[]" {
		return tags
	}
	json.Unmarshal([]byte(jsonStr), &tags)
	return tags
}
