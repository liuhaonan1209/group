package model

import (
	"time"

	"gorm.io/gorm"
)

// AuditLog 操作审计日志表
type AuditLog struct {
	gorm.Model
	UserID       uint      `gorm:"not null;index;comment:操作用户ID" json:"user_id"`
	UserType     string    `gorm:"type:varchar(20);not null;comment:用户类型(passenger/driver/admin)" json:"user_type"`
	UserName     string    `gorm:"type:varchar(100);comment:用户姓名" json:"user_name"`
	Action       string    `gorm:"type:varchar(100);not null;index;comment:操作动作" json:"action"`
	Resource     string    `gorm:"type:varchar(100);not null;comment:操作资源" json:"resource"`
	ResourceID   string    `gorm:"type:varchar(100);index;comment:资源ID" json:"resource_id"`
	Method       string    `gorm:"type:varchar(20);comment:请求方法(GET/POST/PUT/DELETE)" json:"method"`
	Path         string    `gorm:"type:varchar(500);comment:请求路径" json:"path"`
	IPAddress    string    `gorm:"type:varchar(50);comment:IP地址" json:"ip_address"`
	UserAgent    string    `gorm:"type:varchar(500);comment:用户代理" json:"user_agent"`
	RequestBody  string    `gorm:"type:text;comment:请求体" json:"request_body"`
	ResponseCode int       `gorm:"comment:响应状态码" json:"response_code"`
	ErrorMessage string    `gorm:"type:text;comment:错误信息" json:"error_message"`
	Duration     int64     `gorm:"comment:执行时长(毫秒)" json:"duration"`
	OperatedAt   time.Time `gorm:"type:datetime;not null;index;comment:操作时间" json:"operated_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// ExportLog 数据导出日志表
type ExportLog struct {
	gorm.Model
	ExportRecordID uint      `gorm:"not null;index;comment:导出记录ID" json:"export_record_id"`
	UserID         uint      `gorm:"not null;index;comment:用户ID" json:"user_id"`
	UserType       string    `gorm:"type:varchar(20);not null;comment:用户类型" json:"user_type"`
	UserName       string    `gorm:"type:varchar(100);comment:用户姓名" json:"user_name"`
	ExportType     string    `gorm:"type:varchar(50);not null;comment:导出类型(orders/payments/ratings等)" json:"export_type"`
	ExportFields   string    `gorm:"type:text;comment:导出字段列表" json:"export_fields"`
	FilterCondition string   `gorm:"type:text;comment:筛选条件" json:"filter_condition"`
	RecordCount    int       `gorm:"comment:导出记录数" json:"record_count"`
	FileSize       int64     `gorm:"comment:文件大小(字节)" json:"file_size"`
	FileFormat     string    `gorm:"type:varchar(20);comment:文件格式" json:"file_format"`
	DownloadUrl    string    `gorm:"type:varchar(500);comment:下载链接" json:"download_url"`
	Status         string    `gorm:"type:varchar(20);not null;index;comment:状态(pending/processing/completed/failed)" json:"status"`
	ErrorMessage   string    `gorm:"type:text;comment:错误信息" json:"error_message"`
	StartedAt      time.Time `gorm:"type:datetime;comment:开始时间" json:"started_at"`
	CompletedAt    *time.Time `gorm:"type:datetime;comment:完成时间" json:"completed_at"`
	IPAddress      string    `gorm:"type:varchar(50);comment:IP地址" json:"ip_address"`
}

func (ExportLog) TableName() string {
	return "export_logs"
}

// SystemLog 系统异常日志表
type SystemLog struct {
	gorm.Model
	Level        string    `gorm:"type:varchar(20);not null;index;comment:日志级别(INFO/WARN/ERROR/FATAL)" json:"level"`
	Module       string    `gorm:"type:varchar(100);not null;index;comment:模块名称" json:"module"`
	Function     string    `gorm:"type:varchar(100);comment:函数名称" json:"function"`
	Message      string    `gorm:"type:text;not null;comment:日志消息" json:"message"`
	ErrorStack   string    `gorm:"type:text;comment:错误堆栈" json:"error_stack"`
	RequestID    string    `gorm:"type:varchar(100);index;comment:请求ID" json:"request_id"`
	UserID       *uint     `gorm:"index;comment:用户ID" json:"user_id"`
	IPAddress    string    `gorm:"type:varchar(50);comment:IP地址" json:"ip_address"`
	ExtraData    string    `gorm:"type:text;comment:额外数据(JSON)" json:"extra_data"`
	OccurredAt   time.Time `gorm:"type:datetime;not null;index;comment:发生时间" json:"occurred_at"`
}

func (SystemLog) TableName() string {
	return "system_logs"
}

// AccessLog 访问日志表（用于安全审计）
type AccessLog struct {
	gorm.Model
	UserID       *uint     `gorm:"index;comment:用户ID" json:"user_id"`
	UserType     string    `gorm:"type:varchar(20);comment:用户类型" json:"user_type"`
	Method       string    `gorm:"type:varchar(20);not null;comment:请求方法" json:"method"`
	Path         string    `gorm:"type:varchar(500);not null;index;comment:请求路径" json:"path"`
	StatusCode   int       `gorm:"not null;comment:响应状态码" json:"status_code"`
	IPAddress    string    `gorm:"type:varchar(50);not null;index;comment:IP地址" json:"ip_address"`
	UserAgent    string    `gorm:"type:varchar(500);comment:用户代理" json:"user_agent"`
	Referer      string    `gorm:"type:varchar(500);comment:来源页面" json:"referer"`
	Duration     int64     `gorm:"comment:响应时长(毫秒)" json:"duration"`
	RequestSize  int64     `gorm:"comment:请求大小(字节)" json:"request_size"`
	ResponseSize int64     `gorm:"comment:响应大小(字节)" json:"response_size"`
	AccessedAt   time.Time `gorm:"type:datetime;not null;index;comment:访问时间" json:"accessed_at"`
}

func (AccessLog) TableName() string {
	return "access_logs"
}
