package main

import (
	"context"
	"fmt"
	"group/handler/dao"
	"group/handler/model"
	"group/kitex_gen/car/route"
	"sort"
)

// RouteServiceImpl implements the last service interface defined in the IDL.
type RouteServiceImpl struct{}

// ========== 班线管理 ==========

// GetRoute 获取班线基础信息
func (s *RouteServiceImpl) GetRoute(ctx context.Context, req *route.RouteReq) (resp *route.RouteResp, err error) {
	resp = &route.RouteResp{}

	if req.RouteId <= 0 {
		resp.Success = false
		resp.Msg = "班线ID不能为空"
		return resp, nil
	}

	var routeInfo model.Route
	err = dao.RouteInfo(&routeInfo, int(req.RouteId))
	if err != nil {
		resp.Success = false
		resp.Msg = "查询失败"
		return resp, nil
	}

	resp.Success = true
	resp.Msg = "查询成功"
	return
}

// RouteCreate 创建班线
func (s *RouteServiceImpl) RouteCreate(ctx context.Context, req *route.RouteCreateReq) (resp *route.RouteCreateResp, err error) {
	resp = &route.RouteCreateResp{}

	if req.RouteNo == "" {
		resp.Success = false
		resp.Msg = "班线编号不能为空"
		return resp, nil
	}

	// 创建班线
	routeModel := &model.Route{
		RouteNo:        req.RouteNo,
		Fieet:          req.Fieet,
		Tage:           req.Tage,
		StartStationId: uint(req.StartStationId),
		EndStationId:   uint(req.EndStationId),
		IsActive:       "1",
	}

	err = dao.RouteCreate(routeModel)
	if err != nil {
		resp.Success = false
		resp.Msg = "创建失败: " + err.Error()
		return resp, nil
	}

	routeId := int(routeModel.ID)

	// 处理上车站点
	if len(req.UpStations) > 0 {
		upStops, err := processStations(req.UpStations, routeId, model.StopTypePickup)
		if err != nil {
			resp.Success = false
			resp.Msg = "上车站点处理失败: " + err.Error()
			return resp, nil
		}
		if err := dao.AddRouteStops(routeId, upStops); err != nil {
			resp.Success = false
			resp.Msg = "添加上车站点失败: " + err.Error()
			return resp, nil
		}
	}

	// 处理下车站点
	if len(req.DownStations) > 0 {
		// 检查上下车站点冲突
		if err := checkUpDownConflict(req.UpStations, req.DownStations); err != nil {
			resp.Success = false
			resp.Msg = err.Error()
			return resp, nil
		}

		downStops, err := processStations(req.DownStations, routeId, model.StopTypeDropoff)
		if err != nil {
			resp.Success = false
			resp.Msg = "下车站点处理失败: " + err.Error()
			return resp, nil
		}
		if err := dao.AddRouteStops(routeId, downStops); err != nil {
			resp.Success = false
			resp.Msg = "添加下车站点失败: " + err.Error()
			return resp, nil
		}
	}

	resp.Success = true
	resp.Msg = "创建成功"
	resp.RouteId = int64(routeId)
	return
}


// RouteList 班线列表
func (s *RouteServiceImpl) RouteList(ctx context.Context, req *route.RouteListReq) (resp *route.RouteListResp, err error) {
	resp = &route.RouteListResp{}

	page := int(req.Page)
	size := int(req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	list, total, err := dao.RouteList(page, size, req.RouteNo, req.Fieet, req.Tage)
	if err != nil {
		resp.Success = false
		resp.Msg = "查询失败: " + err.Error()
		return resp, nil
	}

	items := make([]*route.RouteItem, len(list))
	for i, r := range list {
		items[i] = &route.RouteItem{
			Id:       int64(r.ID),
			RouteNo:  r.RouteNo,
			Fieet:    r.Fieet,
			Tage:     r.Tage,
			IsActive: r.IsActive == "1",
		}
	}

	resp.Success = true
	resp.Msg = "查询成功"
	resp.List = items
	resp.Total = total
	return
}

// RouteInfo 班线详情
func (s *RouteServiceImpl) RouteInfo(ctx context.Context, req *route.RouteInfoReq) (resp *route.RouteInfoResp, err error) {
	resp = &route.RouteInfoResp{}

	if req.Id <= 0 {
		resp.Success = false
		resp.Msg = "班线ID不能为空"
		return resp, nil
	}

	var routeInfo model.Route
	err = dao.RouteInfo(&routeInfo, int(req.Id))
	if err != nil {
		resp.Success = false
		resp.Msg = "查询失败"
		return resp, nil
	}

	// 获取上车站点
	upStops, _ := dao.GetRouteStopsByType(int(req.Id), model.StopTypePickup)
	upStations := make([]*route.StationInfo, len(upStops))
	for i, s := range upStops {
		upStations[i] = &route.StationInfo{
			StationId:  int64(s.StationId),
			TravelTime: int32(s.TravelTime),
			StopOrder:  int32(s.StopOrder),
		}
	}

	// 获取下车站点
	downStops, _ := dao.GetRouteStopsByType(int(req.Id), model.StopTypeDropoff)
	downStations := make([]*route.StationInfo, len(downStops))
	for i, s := range downStops {
		downStations[i] = &route.StationInfo{
			StationId:  int64(s.StationId),
			TravelTime: int32(s.TravelTime),
			StopOrder:  int32(s.StopOrder),
		}
	}

	resp.Success = true
	resp.Msg = "查询成功"
	resp.Id = int64(routeInfo.ID)
	resp.RouteNo = routeInfo.RouteNo
	resp.Fieet = routeInfo.Fieet
	resp.Tage = routeInfo.Tage
	resp.StartStationId = int64(routeInfo.StartStationId)
	resp.EndStationId = int64(routeInfo.EndStationId)
	resp.UpStations = upStations
	resp.DownStations = downStations
	resp.IsActive = routeInfo.IsActive == "1"
	return
}

// RouteUpdate 更新班线
func (s *RouteServiceImpl) RouteUpdate(ctx context.Context, req *route.RouteUpdateReq) (resp *route.RouteUpdateResp, err error) {
	resp = &route.RouteUpdateResp{}

	if req.Id <= 0 {
		resp.Success = false
		resp.Msg = "班线ID不能为空"
		return resp, nil
	}

	// 更新基础信息
	isActive := "0"
	if req.IsActive {
		isActive = "1"
	}

	routeModel := &model.Route{
		RouteNo:        req.RouteNo,
		Fieet:          req.Fieet,
		Tage:           req.Tage,
		StartStationId: uint(req.StartStationId),
		EndStationId:   uint(req.EndStationId),
		IsActive:       isActive,
	}

	if err := dao.RouteUpdate(int(req.Id), routeModel); err != nil {
		resp.Success = false
		resp.Msg = "更新失败: " + err.Error()
		return resp, nil
	}

	routeId := int(req.Id)

	// 更新上车站点
	if len(req.UpStations) > 0 {
		dao.ClearRouteStops(routeId, model.StopTypePickup)
		upStops, err := processStations(req.UpStations, routeId, model.StopTypePickup)
		if err != nil {
			resp.Success = false
			resp.Msg = "上车站点处理失败: " + err.Error()
			return resp, nil
		}
		dao.AddRouteStops(routeId, upStops)
	}

	// 更新下车站点
	if len(req.DownStations) > 0 {
		if err := checkUpDownConflict(req.UpStations, req.DownStations); err != nil {
			resp.Success = false
			resp.Msg = err.Error()
			return resp, nil
		}

		dao.ClearRouteStops(routeId, model.StopTypeDropoff)
		downStops, err := processStations(req.DownStations, routeId, model.StopTypeDropoff)
		if err != nil {
			resp.Success = false
			resp.Msg = "下车站点处理失败: " + err.Error()
			return resp, nil
		}
		dao.AddRouteStops(routeId, downStops)
	}

	resp.Success = true
	resp.Msg = "更新成功"
	return
}

// RouteDelete 删除班线
func (s *RouteServiceImpl) RouteDelete(ctx context.Context, req *route.RouteDeleteReq) (resp *route.RouteDeleteResp, err error) {
	resp = &route.RouteDeleteResp{}

	if req.Id <= 0 {
		resp.Success = false
		resp.Msg = "班线ID不能为空"
		return resp, nil
	}

	// 删除关联的站点
	dao.ClearRouteStops(int(req.Id), model.StopTypePickup)
	dao.ClearRouteStops(int(req.Id), model.StopTypeDropoff)

	if err := dao.RouteDelete(req.Id); err != nil {
		resp.Success = false
		resp.Msg = "删除失败: " + err.Error()
		return resp, nil
	}

	resp.Success = true
	resp.Msg = "删除成功"
	return
}


// ========== 站点管理 ==========

// AddStation 添加站点到班线
func (s *RouteServiceImpl) AddStation(ctx context.Context, req *route.AddStationReq) (resp *route.AddStationResp, err error) {
	resp = &route.AddStationResp{}

	if req.RouteId <= 0 || req.StationId <= 0 {
		resp.Success = false
		resp.Msg = "班线ID和站点ID不能为空"
		return resp, nil
	}

	if req.StopType != int32(model.StopTypePickup) && req.StopType != int32(model.StopTypeDropoff) {
		resp.Success = false
		resp.Msg = "站点类型无效，1-上车站点，2-下车站点"
		return resp, nil
	}

	err = dao.AddSingleRouteStop(int(req.RouteId), int(req.StationId), int(req.TravelTime), int(req.StopType))
	if err != nil {
		resp.Success = false
		resp.Msg = err.Error()
		return resp, nil
	}

	resp.Success = true
	resp.Msg = "添加成功"
	return
}

// RemoveStation 从班线移除站点
func (s *RouteServiceImpl) RemoveStation(ctx context.Context, req *route.RemoveStationReq) (resp *route.RemoveStationResp, err error) {
	resp = &route.RemoveStationResp{}

	if req.RouteId <= 0 || req.StationId <= 0 {
		resp.Success = false
		resp.Msg = "班线ID和站点ID不能为空"
		return resp, nil
	}

	err = dao.RemoveRouteStop(int(req.RouteId), int(req.StationId))
	if err != nil {
		resp.Success = false
		resp.Msg = err.Error()
		return resp, nil
	}

	resp.Success = true
	resp.Msg = "移除成功"
	return
}

// ========== 班次管理 ==========

// ScheduleCreate 创建班次
func (s *RouteServiceImpl) ScheduleCreate(ctx context.Context, req *route.ScheduleCreateReq) (resp *route.ScheduleCreateResp, err error) {
	resp = &route.ScheduleCreateResp{}

	if req.RouteId <= 0 {
		resp.Success = false
		resp.Msg = "班线ID不能为空"
		return resp, nil
	}

	schedule := &model.BusSchedules{
		RouteId:       uint(req.RouteId),
		DepartureTime: parseTimeToMinutes(req.DepartureTime),
		ArrivalTime:   parseTimeToMinutes(req.ArrivalTime),
		Capacity:      uint(req.Capacity),
		IsActive:      true,
	}

	if err := dao.ScheduleCreate(schedule); err != nil {
		resp.Success = false
		resp.Msg = "创建失败: " + err.Error()
		return resp, nil
	}

	resp.Success = true
	resp.Msg = "创建成功"
	resp.ScheduleId = int64(schedule.ID)
	return
}

// ScheduleList 班次列表
func (s *RouteServiceImpl) ScheduleList(ctx context.Context, req *route.ScheduleListReq) (resp *route.ScheduleListResp, err error) {
	resp = &route.ScheduleListResp{}

	page := int(req.Page)
	size := int(req.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	list, total, err := dao.ScheduleList(int(req.RouteId), page, size)
	if err != nil {
		resp.Success = false
		resp.Msg = "查询失败: " + err.Error()
		return resp, nil
	}

	items := make([]*route.ScheduleInfo, len(list))
	for i, s := range list {
		items[i] = &route.ScheduleInfo{
			Id:            int64(s.ID),
			RouteId:       int64(s.RouteId),
			DepartureTime: formatMinutesToTime(s.DepartureTime),
			ArrivalTime:   formatMinutesToTime(s.ArrivalTime),
			Capacity:      int32(s.Capacity),
			IsActive:      s.IsActive,
		}
	}

	resp.Success = true
	resp.Msg = "查询成功"
	resp.List = items
	resp.Total = total
	return
}

// ScheduleUpdate 更新班次
func (s *RouteServiceImpl) ScheduleUpdate(ctx context.Context, req *route.ScheduleUpdateReq) (resp *route.ScheduleUpdateResp, err error) {
	resp = &route.ScheduleUpdateResp{}

	if req.Id <= 0 {
		resp.Success = false
		resp.Msg = "班次ID不能为空"
		return resp, nil
	}

	schedule := &model.BusSchedules{
		DepartureTime: parseTimeToMinutes(req.DepartureTime),
		ArrivalTime:   parseTimeToMinutes(req.ArrivalTime),
		Capacity:      uint(req.Capacity),
		IsActive:      req.IsActive,
	}

	if err := dao.ScheduleUpdate(int(req.Id), schedule); err != nil {
		resp.Success = false
		resp.Msg = "更新失败: " + err.Error()
		return resp, nil
	}

	resp.Success = true
	resp.Msg = "更新成功"
	return
}

// ScheduleDelete 删除班次
func (s *RouteServiceImpl) ScheduleDelete(ctx context.Context, req *route.ScheduleDeleteReq) (resp *route.ScheduleDeleteResp, err error) {
	resp = &route.ScheduleDeleteResp{}

	if req.Id <= 0 {
		resp.Success = false
		resp.Msg = "班次ID不能为空"
		return resp, nil
	}

	if err := dao.ScheduleDelete(req.Id); err != nil {
		resp.Success = false
		resp.Msg = "删除失败: " + err.Error()
		return resp, nil
	}

	resp.Success = true
	resp.Msg = "删除成功"
	return
}


// ========== 辅助函数 ==========

// processStations 处理站点列表：排重、排序
func processStations(stations []*route.StationInfo, routeId int, stopType int) ([]model.RouteStops, error) {
	if len(stations) == 0 {
		return nil, nil
	}

	// 排重检查
	idMap := make(map[int64]bool)
	for _, s := range stations {
		if idMap[s.StationId] {
			return nil, fmt.Errorf("站点ID %d 重复", s.StationId)
		}
		idMap[s.StationId] = true
	}

	// 按行驶时间排序
	sorted := make([]*route.StationInfo, len(stations))
	copy(sorted, stations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].TravelTime < sorted[j].TravelTime
	})

	// 第一个站点时间设为0
	sorted[0].TravelTime = 0

	// 转换为 RouteStops
	stops := make([]model.RouteStops, len(sorted))
	for i, s := range sorted {
		stops[i] = model.RouteStops{
			RouteId:    routeId,
			StationId:  int(s.StationId),
			TravelTime: int(s.TravelTime),
			StopOrder:  i + 1,
			StopType:   stopType,
			IsActive:   true,
		}
	}

	return stops, nil
}

// checkUpDownConflict 检查上下车站点冲突
func checkUpDownConflict(upStations, downStations []*route.StationInfo) error {
	upIds := make(map[int64]bool)
	for _, s := range upStations {
		upIds[s.StationId] = true
	}

	for _, s := range downStations {
		if upIds[s.StationId] {
			return fmt.Errorf("站点ID %d 不能同时作为上车站点和下车站点", s.StationId)
		}
	}

	return nil
}

// parseTimeToMinutes 将 HH:mm 格式转换为分钟数
func parseTimeToMinutes(timeStr string) uint {
	if len(timeStr) < 5 {
		return 0
	}
	var hour, minute int
	fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	return uint(hour*60 + minute)
}

// formatMinutesToTime 将分钟数转换为 HH:mm 格式
func formatMinutesToTime(minutes uint) string {
	hour := minutes / 60
	minute := minutes % 60
	return fmt.Sprintf("%02d:%02d", hour, minute)
}
