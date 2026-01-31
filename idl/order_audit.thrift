namespace go car.order_audit

// ==================== 审计日志记录接口 ====================

// 审计日志记录请求
struct AuditLogCreateReq {
    1: string orderNo               // 订单号
    2: string action                // 操作动作（create/update/cancel/pay/refund等）
    3: i64 operatorId               // 操作人ID
    4: string operatorType          // 操作人类型（passenger/driver/admin/system）
    5: string operatorName          // 操作人姓名
    6: string operatorIp            // 操作人IP地址
    7: optional string beforeData   // 操作前数据（JSON格式）
    8: optional string afterData    // 操作后数据（JSON格式）
    9: optional string reason       // 操作原因
    10: optional string remark      // 备注信息
    11: string userAgent            // 用户代理信息
    12: string requestId            // 请求ID（用于追踪）
}

// 审计日志记录响应
struct AuditLogCreateResp {
    1: i64 logId                    // 日志ID
    2: string message               // 返回消息
    3: bool success                 // 是否成功
}

// ==================== 审计日志查询接口 ====================

// 审计日志查询请求
struct AuditLogQueryReq {
    1: optional string orderNo      // 订单号（可选）
    2: optional string action       // 操作动作（可选）
    3: optional i64 operatorId      // 操作人ID（可选）
    4: optional string operatorType // 操作人类型（可选）
    5: optional string startTime    // 开始时间（可选）
    6: optional string endTime      // 结束时间（可选）
    7: optional string operatorIp   // 操作人IP（可选）
    8: i32 page                     // 页码（从1开始）
    9: i32 pageSize                 // 每页数量
}

// 审计日志信息
struct AuditLogInfo {
    1: i64 logId                    // 日志ID
    2: string orderNo               // 订单号
    3: string action                // 操作动作
    4: i64 operatorId               // 操作人ID
    5: string operatorType          // 操作人类型
    6: string operatorName          // 操作人姓名
    7: string operatorIp            // 操作人IP地址
    8: string beforeData            // 操作前数据
    9: string afterData             // 操作后数据
    10: string reason               // 操作原因
    11: string remark               // 备注信息
    12: string userAgent            // 用户代理信息
    13: string requestId            // 请求ID
    14: string createdAt            // 创建时间
}

// 审计日志查询响应
struct AuditLogQueryResp {
    1: list<AuditLogInfo> logs      // 日志列表
    2: i64 total                    // 总数量
    3: i32 page                     // 当前页码
    4: i32 pageSize                 // 每页数量
    5: string message               // 返回消息
    6: bool success                 // 是否成功
}

// ==================== 异常行为检测接口 ====================

// 异常行为检测请求
struct AnomalyDetectionReq {
    1: optional i64 userId          // 用户ID（可选）
    2: optional string userType     // 用户类型（可选）
    3: optional string detectionType // 检测类型（频繁取消/异常支付/IP异常等）
    4: optional string startTime    // 开始时间（可选）
    5: optional string endTime      // 结束时间（可选）
    6: i32 page                     // 页码
    7: i32 pageSize                 // 每页数量
}

// 异常行为信息
struct AnomalyInfo {
    1: i64 anomalyId                // 异常ID
    2: i64 userId                   // 用户ID
    3: string userType              // 用户类型
    4: string userName              // 用户姓名
    5: string anomalyType           // 异常类型
    6: string description           // 异常描述
    7: string riskLevel             // 风险等级（low/medium/high/critical）
    8: i32 occurrenceCount          // 发生次数
    9: string relatedOrders         // 相关订单号（逗号分隔）
    10: string detectedAt           // 检测时间
    11: string status               // 处理状态（pending/reviewing/resolved/ignored）
    12: string handlerName          // 处理人姓名
    13: string handleResult         // 处理结果
}

// 异常行为检测响应
struct AnomalyDetectionResp {
    1: list<AnomalyInfo> anomalies  // 异常列表
    2: i64 total                    // 总数量
    3: i32 page                     // 当前页码
    4: i32 pageSize                 // 每页数量
    5: string message               // 返回消息
    6: bool success                 // 是否成功
}

// ==================== 实时监控统计接口 ====================

// 实时监控统计请求
struct MonitorStatsReq {
    1: optional string timeRange    // 时间范围（1h/24h/7d/30d）
    2: optional string statsType    // 统计类型（order/payment/anomaly/all）
}

// 实时监控统计响应
struct MonitorStatsResp {
    1: i64 totalOrders              // 总订单数
    2: i64 todayOrders              // 今日订单数
    3: i64 pendingOrders            // 待处理订单数
    4: i64 completedOrders          // 已完成订单数
    5: i64 cancelledOrders          // 已取消订单数
    6: i64 totalAuditLogs           // 总审计日志数
    7: i64 todayAuditLogs           // 今日审计日志数
    8: i64 totalAnomalies           // 总异常数
    9: i64 pendingAnomalies         // 待处理异常数
    10: i64 highRiskAnomalies       // 高风险异常数
    11: double avgResponseTime      // 平均响应时间（毫秒）
    12: map<string, i64> actionStats // 操作统计（action -> count）
    13: map<string, i64> anomalyStats // 异常统计（type -> count）
    14: string message              // 返回消息
    15: bool success                // 是否成功
}

// ==================== 异常行为处理接口 ====================

// 异常行为处理请求
struct AnomalyHandleReq {
    1: i64 anomalyId                // 异常ID
    2: i64 handlerId                // 处理人ID
    3: string handlerName           // 处理人姓名
    4: string handleAction          // 处理动作（resolve/ignore/escalate）
    5: string handleResult          // 处理结果
    6: optional string remark       // 备注
}

// 异常行为处理响应
struct AnomalyHandleResp {
    1: string message               // 返回消息
    2: bool success                 // 是否成功
}

// ==================== 日志导出接口 ====================

// 日志导出请求
struct AuditLogExportReq {
    1: optional string orderNo      // 订单号（可选）
    2: optional string action       // 操作动作（可选）
    3: optional i64 operatorId      // 操作人ID（可选）
    4: optional string startTime    // 开始时间（可选）
    5: optional string endTime      // 结束时间（可选）
    6: string exportFormat          // 导出格式（csv/excel/json）
}

// 日志导出响应
struct AuditLogExportResp {
    1: string downloadUrl           // 下载链接
    2: string fileName              // 文件名
    3: i64 recordCount              // 记录数量
    4: string message               // 返回消息
    5: bool success                 // 是否成功
}

// ==================== 服务定义 ====================

service OrderAuditService {
    // 审计日志记录接口
    AuditLogCreateResp AuditLogCreate(1: AuditLogCreateReq req)
    
    // 审计日志查询接口
    AuditLogQueryResp AuditLogQuery(1: AuditLogQueryReq req)
    
    // 异常行为检测接口
    AnomalyDetectionResp AnomalyDetection(1: AnomalyDetectionReq req)
    
    // 实时监控统计接口
    MonitorStatsResp MonitorStats(1: MonitorStatsReq req)
    
    // 异常行为处理接口
    AnomalyHandleResp AnomalyHandle(1: AnomalyHandleReq req)
    
    // 日志导出接口
    AuditLogExportResp AuditLogExport(1: AuditLogExportReq req)
}
