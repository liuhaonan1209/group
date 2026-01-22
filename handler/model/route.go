package model

import (
	"gorm.io/gorm"
)

// 线路
type Route struct {
	gorm.Model
	RouteNo        string         `gorm:"type:varchar(20);unique_index;comment:线路编号" json:"routeNo"`         // 线路编号
	Fieet          string         `gorm:"type:varchar(20);unique_index;comment:车队编号" json:"fieet"`           // 车队编号（字段名笔误，建议改为Fleet）
	Tage           string         `gorm:"type:varchar(50);unique_index;comment:线路标签" json:"tage"`            // 线路标签（字段名笔误，建议改为Tag）
	TagShowStatus  string         `gorm:"type:varchar(50);unique_index;comment:标签显示状态" json:"tagShowStatus"` // 标签显示状态（示例：show-显示、hide-隐藏）
	StartStationId uint           `gorm:"type:int;comment:起点站ID" json:"startStationId"`                      // 起点站ID（关联站点表主键）
	EndStationId   uint           `gorm:"type:int;comment:终点站ID" json:"endStationId"`                        // 终点站ID（关联站点表主键）
	UpStations     []StartStation `gorm:"-" json:"upStations"`                                               // 上车站点列表（仅业务逻辑/序列化使用）
	DownStations   []StartStation `gorm:"-" json:"downStations"`                                             // 下车站点列表（仅业务逻辑/序列化使用）
	PassStations   string         `gorm:"type:varchar(50);unique_index;comment:途经站点" json:"passStations"`    // 途经站点（站点ID/名称拼接串）
	RoutePath      string         `gorm:"type:varchar(50);unique_index;comment:线路路径" json:"routePath"`       // 线路路径（经纬度/站点顺序拼接串）
	IsActive       string         `gorm:"type:int;comment:线路启用状态" json:"isActive"`                           // 线路启用状态（类型建议改为bool/int，示例：1-启用、0-禁用）
}
