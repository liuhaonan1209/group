package model

import (
	"time"

	"gorm.io/gorm"
)

// OrderAuditLog 订单审计日志表
type OrderAuditLog struct {
	gorm.Model
	OrderNo      string    `gorm:"type:varchar(50);index;not null" json:"order_no"`        // 订单号
	Action       string    `gorm:"type:varchar(50);index;not null" json:"action"`          // 操作动作
	OperatorID   int64     `gorm:"index;not null" json:"operator_id"`                      // 操作人ID
	OperatorType string    `gorm:"type:varchar(20);index;not null" json:"operator_type"`   // 操作人类型
	OperatorName string    `gorm:"type:varchar(50);not null" json:"operator_name"`         // 操作人姓名
	OperatorIP   string    `gorm:"type:varchar(50);index" json:"operator_ip"`              // 操作人IP地址
	BeforeData   string    `gorm:"type:text" json:"before_data"`                           // 操作前数据
	AfterData    string    `gorm:"type:text" json:"after_data"`                            // 操作后数据
	Reason       string    `gorm:"type:varchar(500)" json:"reason"`                        // 操作原因
	Remark       string    `gorm:"type:text" json:"remark"`                                // 备注信息
	UserAgent    string    `gorm:"type:varchar(500)" json:"user_agent"`                    // 用户代理信息
	RequestID    string    `gorm:"type:varchar(100);index" json:"request_id"`              // 请求ID
	CreatedAt    time.Time `gorm:"index;not null" json:"created_at"`                       // 创建时间
}

// TableName 指定表名
func (OrderAuditLog) TableName() string {
	return "order_audit_log"
}

// OrderAnomaly 订单异常行为表
type OrderAnomaly struct {
	gorm.Model
	UserID          int64     `gorm:"index;not null" json:"user_id"`                        // 用户ID
	UserType        string    `gorm:"type:varchar(20);index;not null" json:"user_type"`     // 用户类型
	UserName        string    `gorm:"type:varchar(50);not null" json:"user_name"`           // 用户姓名
	AnomalyType     string    `gorm:"type:varchar(50);index;not null" json:"anomaly_type"`  // 异常类型
	Description     string    `gorm:"type:text;not null" json:"description"`                // 异常描述
	RiskLevel       string    `gorm:"type:varchar(20);index;not null" json:"risk_level"`    // 风险等级
	OccurrenceCount int       `gorm:"not null;default:1" json:"occurrence_count"`           // 发生次数
	RelatedOrders   string    `gorm:"type:text" json:"related_orders"`                      // 相关订单号
	DetectedAt      time.Time `gorm:"index;not null" json:"detected_at"`                    // 检测时间
	Status          string    `gorm:"type:varchar(20);index;not null;default:'pending'" json:"status"` // 处理状态
	HandlerID       *int64    `gorm:"index" json:"handler_id"`                              // 处理人ID
	HandlerName     string    `gorm:"type:varchar(50)" json:"handler_name"`                 // 处理人姓名
	HandleResult    string    `gorm:"type:text" json:"handle_result"`                       // 处理结果
	HandledAt       *time.Time `gorm:"index" json:"handled_at"`                             // 处理时间
}

// TableName 指定表名
func (OrderAnomaly) TableName() string {
	return "order_anomaly"
}

// AuditLogExport 审计日志导出记录表
type AuditLogExport struct {
	gorm.Model
	ExportUserID   int64     `gorm:"index;not null" json:"export_user_id"`               // 导出用户ID
	ExportUserName string    `gorm:"type:varchar(50);not null" json:"export_user_name"`  // 导出用户姓名
	ExportFormat   string    `gorm:"type:varchar(20);not null" json:"export_format"`     // 导出格式
	FilterCondition string   `gorm:"type:text" json:"filter_condition"`                  // 过滤条件
	RecordCount    int64     `gorm:"not null" json:"record_count"`                       // 记录数量
	FileName       string    `gorm:"type:varchar(200);not null" json:"file_name"`        // 文件名
	FilePath       string    `gorm:"type:varchar(500);not null" json:"file_path"`        // 文件路径
	DownloadURL    string    `gorm:"type:varchar(500)" json:"download_url"`              // 下载链接
	Status         string    `gorm:"type:varchar(20);index;not null;default:'pending'" json:"status"` // 导出状态
	ExpiresAt      time.Time `gorm:"index;not null" json:"expires_at"`                   // 过期时间
}

// TableName 指定表名
func (AuditLogExport) TableName() string {
	return "audit_log_export"
}
