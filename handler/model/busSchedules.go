package model

import "gorm.io/gorm"

// 班次表
type BusSchedules struct {
	gorm.Model
	RouteId       uint `gorm:"type:int" json:"routeId"`       //班线id
	DepartureTime uint `gorm:"type:int" json:"departureTime"` //发车时间
	ArrivalTime   uint `gorm:"type:int" json:"arrivalTime"`   //预计到达时间
	Capacity      uint `gorm:"type:int" json:"capacity"`      //座位容量
	IsActive      bool `gorm:"type:bool" json:"isActive"`     //是否启用
}
