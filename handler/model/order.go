package model

import (
	"time"

	"gorm.io/gorm"
)

// Trip 行程表模型
type Trip struct {
	gorm.Model
	PublisherID   uint      `gorm:"not null;comment:发布者ID" json:"publisher_id"`
	PublisherType string    `gorm:"type:varchar(20);not null;comment:发布者类型(passenger/driver)" json:"publisher_type"`
	StartPoint    string    `gorm:"type:varchar(200);not null;comment:起点" json:"start_point"`
	EndPoint      string    `gorm:"type:varchar(200);not null;comment:终点" json:"end_point"`
	DepartureTime time.Time `gorm:"type:datetime;not null;comment:出行时间" json:"departure_time"`
	VehicleInfo   string    `gorm:"type:varchar(200);comment:车辆信息" json:"vehicle_info"`
	SpecialNeeds  string    `gorm:"type:text;comment:特殊需求" json:"special_needs"`
	ContactWay    string    `gorm:"type:varchar(100);comment:联系方式" json:"contact_way"`
	Status        string    `gorm:"type:varchar(20);default:'pending';comment:行程状态(pending/confirmed/completed/cancelled)" json:"status"`
}

func (Trip) TableName() string {
	return "trips"
}

// TripShare 行程分享记录表
type TripShare struct {
	gorm.Model
	TripID       uint      `gorm:"not null;comment:行程ID" json:"trip_id"`
	UserID       uint      `gorm:"not null;comment:分享用户ID" json:"user_id"`
	ShareTargets string    `gorm:"type:text;comment:分享对象" json:"share_targets"`
	ShareLink    string    `gorm:"type:varchar(500);comment:分享链接" json:"share_link"`
	SharedAt     time.Time `gorm:"type:datetime;comment:分享时间" json:"shared_at"`
}

func (TripShare) TableName() string {
	return "trip_shares"
}

// PassengerHelp 乘客求助记录表
type PassengerHelp struct {
	gorm.Model
	TripID      uint   `gorm:"not null;comment:行程ID" json:"trip_id"`
	PassengerID uint   `gorm:"not null;comment:乘客ID" json:"passenger_id"`
	HelpType    string `gorm:"type:varchar(50);not null;comment:求助类型" json:"help_type"`
	Status      string `gorm:"type:varchar(20);default:'pending';comment:处理状态" json:"status"`
	ContactInfo string `gorm:"type:text;comment:联系信息" json:"contact_info"`
	// 移除额外的CreatedAt字段，因为gorm.Model已经包含了CreatedAt
}

func (PassengerHelp) TableName() string {
	return "passenger_helps"
}

// ExportRecord 数据导出记录表
type ExportRecord struct {
	gorm.Model
	UserID       uint      `gorm:"not null;comment:用户ID" json:"user_id"`
	UserType     string    `gorm:"type:varchar(20);not null;comment:用户类型(passenger/driver)" json:"user_type"`
	ExportFields string    `gorm:"type:text;comment:导出字段" json:"export_fields"`
	Status       string    `gorm:"type:varchar(20);default:'pending';comment:导出状态(pending/completed/failed)" json:"status"`
	DownloadUrl  string    `gorm:"type:varchar(500);comment:下载链接" json:"download_url"`
	FileFormat   string    `gorm:"type:varchar(20);default:'csv';comment:文件格式" json:"file_format"`
	ExpiresAt    time.Time `gorm:"type:datetime;comment:过期时间" json:"expires_at"`
}

func (ExportRecord) TableName() string {
	return "export_records"
}
