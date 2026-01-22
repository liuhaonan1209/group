package model

import "gorm.io/gorm"

// RouteStops 班线途经站点关联表
// 用于存储班线与站点的关联关系，包含行驶时间和站点类型
type RouteStops struct {
	gorm.Model
	RouteId    int  `gorm:"type:int;index:idx_route_station,unique;comment:班线ID" json:"routeId"`   // 班线ID
	StationId  int  `gorm:"type:int;index:idx_route_station,unique;comment:站点ID" json:"stationId"` // 站点ID
	TravelTime int  `gorm:"type:int;default:0;comment:行驶时间" json:"travelTime"`                     // 从起点到此站点的行驶时间（分钟）
	StopOrder  int  `gorm:"type:int;comment:停靠顺序" json:"stopOrder"`                                // 停靠顺序（根据时间自动计算）
	StopType   int  `gorm:"type:tinyint;default:1;comment:站点类型" json:"stopType"`                   // 站点类型：1-上车站点，2-下车站点
	IsActive   bool `gorm:"type:tinyint(1);default:true;comment:是否启用" json:"isActive"`             // 是否启用
}

// StopType 常量定义
const (
	StopTypePickup  = 1 // 上车站点
	StopTypeDropoff = 2 // 下车站点
)
