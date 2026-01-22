package dao

import (
	"fmt"
	"group/global"
	"group/handler/model"
	"sort"
)

// ========== 班线 CRUD ==========

// RouteCreate 创建班线
func RouteCreate[T any](data *T) error {
	return global.DB.Create(data).Error
}

// RouteList 班线列表（分页）
func RouteList(page, size int, routeNo, fieet, tage string) (list []model.Route, total int64, err error) {
	db := global.DB.Model(&model.Route{})
	
	if routeNo != "" {
		db = db.Where("route_no LIKE ?", "%"+routeNo+"%")
	}
	if fieet != "" {
		db = db.Where("fieet = ?", fieet)
	}
	if tage != "" {
		db = db.Where("tage = ?", tage)
	}
	
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	
	err = db.Offset((page - 1) * size).Limit(size).Find(&list).Error
	return
}

// RouteInfo 班线详情
func RouteInfo[T any](data *T, id int) (err error) {
	err = global.DB.Model(data).Limit(1).Where("id=?", id).Find(data).Error
	return
}

// RouteUpdate 更新班线
func RouteUpdate(id int, data *model.Route) error {
	return global.DB.Model(&model.Route{}).Where("id = ?", id).Updates(data).Error
}

// RouteDelete 删除班线
func RouteDelete(id int64) error {
	return global.DB.Where("id = ?", id).Delete(&model.Route{}).Error
}

// ========== 班线站点管理 ==========

// AddRouteStops 批量添加班线站点（带排重和排序）
func AddRouteStops(routeId int, stops []model.RouteStops) error {
	if len(stops) == 0 {
		return nil
	}

	// 检查是否有重复站点（同一班线同一站点）
	stationIds := make([]int, len(stops))
	for i, s := range stops {
		stationIds[i] = s.StationId
	}

	var existingStops []model.RouteStops
	global.DB.Where("route_id = ? AND station_id IN ?", routeId, stationIds).Find(&existingStops)
	if len(existingStops) > 0 {
		return fmt.Errorf("站点ID %d 已存在于该班线中", existingStops[0].StationId)
	}

	return global.DB.Create(&stops).Error
}

// AddSingleRouteStop 添加单个站点到班线
func AddSingleRouteStop(routeId, stationId, travelTime, stopType int) error {
	// 检查站点是否已存在
	var existing model.RouteStops
	err := global.DB.Where("route_id = ? AND station_id = ?", routeId, stationId).First(&existing).Error
	if err == nil {
		return fmt.Errorf("该站点已存在于班线中")
	}

	// 检查上下车冲突
	if err := CheckStationTypeConflict(routeId, stationId, stopType); err != nil {
		return err
	}

	// 获取当前最大排序号
	var maxOrder int
	global.DB.Model(&model.RouteStops{}).
		Where("route_id = ? AND stop_type = ?", routeId, stopType).
		Select("COALESCE(MAX(stop_order), 0)").Scan(&maxOrder)

	stop := model.RouteStops{
		RouteId:    routeId,
		StationId:  stationId,
		TravelTime: travelTime,
		StopOrder:  maxOrder + 1,
		StopType:   stopType,
		IsActive:   true,
	}

	if err := global.DB.Create(&stop).Error; err != nil {
		return err
	}

	// 重新排序
	return ReorderRouteStops(routeId, stopType)
}

// GetRouteStopsByType 获取班线的上车或下车站点（已排序）
func GetRouteStopsByType(routeId int, stopType int) ([]model.RouteStops, error) {
	var stops []model.RouteStops
	err := global.DB.Where("route_id = ? AND stop_type = ? AND is_active = ?", routeId, stopType, true).
		Order("stop_order ASC").
		Find(&stops).Error
	return stops, err
}

// GetRouteStopsWithStation 获取班线站点（包含站点详情）
func GetRouteStopsWithStation(routeId int, stopType int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := global.DB.Table("route_stops rs").
		Select("rs.*, ss.name as station_name, ss.address").
		Joins("LEFT JOIN start_stations ss ON rs.station_id = ss.id").
		Where("rs.route_id = ? AND rs.stop_type = ? AND rs.is_active = ?", routeId, stopType, true).
		Order("rs.stop_order ASC").
		Find(&results).Error
	return results, err
}

// RemoveRouteStop 移除班线站点
func RemoveRouteStop(routeId, stationId int) error {
	var stop model.RouteStops
	err := global.DB.Where("route_id = ? AND station_id = ?", routeId, stationId).First(&stop).Error
	if err != nil {
		return fmt.Errorf("站点不存在")
	}

	stopType := stop.StopType
	if err := global.DB.Delete(&stop).Error; err != nil {
		return err
	}

	// 重新排序
	return ReorderRouteStops(routeId, stopType)
}

// CheckStationTypeConflict 检查站点类型冲突（同一站点不能同时是上车和下车站点）
func CheckStationTypeConflict(routeId, stationId, newStopType int) error {
	var existing model.RouteStops
	err := global.DB.Where("route_id = ? AND station_id = ?", routeId, stationId).First(&existing).Error
	if err == nil && existing.StopType != newStopType {
		return fmt.Errorf("该站点已作为%s站点存在，不能同时设为%s站点",
			getStopTypeName(existing.StopType), getStopTypeName(newStopType))
	}
	return nil
}

// ReorderRouteStops 根据行驶时间重新排序站点
func ReorderRouteStops(routeId, stopType int) error {
	var stops []model.RouteStops
	err := global.DB.Where("route_id = ? AND stop_type = ?", routeId, stopType).Find(&stops).Error
	if err != nil {
		return err
	}

	if len(stops) == 0 {
		return nil
	}

	// 按行驶时间排序
	sort.Slice(stops, func(i, j int) bool {
		return stops[i].TravelTime < stops[j].TravelTime
	})

	// 第一个站点时间设为0
	stops[0].TravelTime = 0

	// 更新排序号
	for i, stop := range stops {
		global.DB.Model(&model.RouteStops{}).
			Where("id = ?", stop.ID).
			Updates(map[string]interface{}{
				"stop_order":  i + 1,
				"travel_time": stop.TravelTime,
			})
	}

	return nil
}

// ClearRouteStops 清空班线的所有站点
func ClearRouteStops(routeId int, stopType int) error {
	return global.DB.Where("route_id = ? AND stop_type = ?", routeId, stopType).
		Delete(&model.RouteStops{}).Error
}

func getStopTypeName(stopType int) string {
	if stopType == model.StopTypePickup {
		return "上车"
	}
	return "下车"
}

// ========== 班次管理 ==========

// ScheduleCreate 创建班次
func ScheduleCreate(schedule *model.BusSchedules) error {
	return global.DB.Create(schedule).Error
}

// ScheduleList 班次列表
func ScheduleList(routeId int, page, size int) (list []model.BusSchedules, total int64, err error) {
	db := global.DB.Model(&model.BusSchedules{})
	
	if routeId > 0 {
		db = db.Where("route_id = ?", routeId)
	}
	
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	
	err = db.Offset((page - 1) * size).Limit(size).Order("departure_time ASC").Find(&list).Error
	return
}

// ScheduleInfo 班次详情
func ScheduleInfo(id int) (schedule model.BusSchedules, err error) {
	err = global.DB.Where("id = ?", id).First(&schedule).Error
	return
}

// ScheduleUpdate 更新班次
func ScheduleUpdate(id int, data *model.BusSchedules) error {
	return global.DB.Model(&model.BusSchedules{}).Where("id = ?", id).Updates(data).Error
}

// ScheduleDelete 删除班次
func ScheduleDelete(id int64) error {
	return global.DB.Where("id = ?", id).Delete(&model.BusSchedules{}).Error
}
