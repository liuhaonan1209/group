package dao

import (
	"fmt"
	"group/global"
	"group/handler/model"
	"sort"
	"time"
)

// 添加站点
func StartStationCreate[T any](data *T) error {

	return global.DB.Create(data).Error
}

//修改站点

func StartStationUpdate[T any](data *T, id int) error {
	return global.DB.Model(data).Update("id=?", id).Error
}

// 站点展示
func StartStationList[T any](data *[]T, page, size int) (list []T, err error) {
	err = global.DB.Offset((page - 1) * size).Limit(size).Find(data).Error
	return
}

// 站点详情
func StartStationInfo[T any](data *[]T, id int) (list []T, err error) {
	err = global.DB.Model(data).Where("id=?", id).Limit(1).Find(data).Error
	return
}

// 站点排重
func RemoveDuplicateStops(stops []model.StartStation, stopsType string) ([]model.StartStation, error) {
	if len(stops) == 0 {
		return []model.StartStation{}, nil
	}

	uniqueStops := make([]model.StartStation, 0, len(stops))
	idMap := make(map[int64]bool)    //按照id去重
	nameMap := make(map[string]bool) //按照名称去重

	for _, stop := range stops {
		//判断站点是否启用
		if !stop.IsActive {
			return nil, fmt.Errorf("%s站点[%s]已禁用，无法加入", stopsType, stop.Name)
		}
		if stop.ID == 0 && stop.Name == "" {
			return nil, fmt.Errorf("站点和名称不可以为空")
		}
		if stop.ID != 0 {
			if idMap[int64(stop.ID)] {
				return nil, fmt.Errorf("%s站点重复（ID）：%d-%s", stopsType, stop.ID, stop.Name)
			}
			idMap[int64(stop.ID)] = true
		} else {
			if nameMap[stop.Name] {
				return nil, fmt.Errorf("%s站点重复（名称）：%s", stopsType, stop.Name)
			}
			nameMap[stop.Name] = true
		}
		uniqueStops = append(uniqueStops, stop)
	}
	return uniqueStops, nil
}

// 校验班车线路上下车站点（批量排重）
func ValidateBusLineStops(line *model.Route) error {
	// 1. 先对上车站点内部排重
	upstops, err := RemoveDuplicateStops(line.UpStations, "上车")
	if err != nil {
		return err
	}

	// 2. 对下车站点内部排重
	downstops, err := RemoveDuplicateStops(line.DownStations, "下车")
	if err != nil {
		return err
	}

	// 3. 检查上下车站点之间是否有重复（同一站点不能同时出现在上车和下车位置）
	if err := CheckUpDownStationConflict(upstops, downstops); err != nil {
		return err
	}

	// 4. 按时间排序上车站点
	sortedUpStops, err := SortStopsByTime(upstops)
	if err != nil {
		return err
	}
	line.UpStations = sortedUpStops

	// 5. 按时间排序下车站点
	sortedDownStops, err := SortStopsByTime(downstops)
	if err != nil {
		return err
	}
	line.DownStations = sortedDownStops

	return nil
}

// CheckUpDownStationConflict 检查上下车站点之间是否有重复
// 同一个站点不能同时出现在上车站点和下车站点的位置
func CheckUpDownStationConflict(upStops, downStops []model.StartStation) error {
	// 构建上车站点的ID和名称映射
	upIdMap := make(map[uint]string)
	upNameMap := make(map[string]bool)

	for _, stop := range upStops {
		if stop.ID != 0 {
			upIdMap[stop.ID] = stop.Name
		}
		if stop.Name != "" {
			upNameMap[stop.Name] = true
		}
	}

	// 检查下车站点是否与上车站点冲突
	for _, stop := range downStops {
		if stop.ID != 0 {
			if name, exists := upIdMap[stop.ID]; exists {
				return fmt.Errorf("站点[%s]不能同时作为上车站点和下车站点", name)
			}
		}
		if stop.Name != "" {
			if upNameMap[stop.Name] {
				return fmt.Errorf("站点[%s]不能同时作为上车站点和下车站点", stop.Name)
			}
		}
	}

	return nil
}

// SortStopsByTime 按照从起点的行驶时间对站点进行排序
// 排序规则：时间由小到大
// 第一个站点的时间自动设为0（起点）
// 举例：A到B 20分钟, A到C 40分钟, A到D 30分钟 => 排序为 A, B, D, C
func SortStopsByTime(stops []model.StartStation) ([]model.StartStation, error) {
	if len(stops) == 0 {
		return []model.StartStation{}, nil
	}

	// 复制切片避免修改原数据
	sortedStops := make([]model.StartStation, len(stops))
	copy(sortedStops, stops)

	// 按照 TimeFromStart 从小到大排序
	sort.Slice(sortedStops, func(i, j int) bool {
		return sortedStops[i].TimeFromStart.Before(sortedStops[j].TimeFromStart)
	})

	// 第一个站点（起点）时间设为0
	zeroTime := time.Time{}
	sortedStops[0].TimeFromStart = zeroTime

	return sortedStops, nil
}
