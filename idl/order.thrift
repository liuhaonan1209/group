namespace go car.order

// ==================== 行程发布接口 ====================

// 行程发布请求（乘客发布行程）
struct TripPublishReq {
    1: i64 passengerId              // 乘客ID
    2: string startPoint            // 起点
    3: string endPoint              // 终点
    4: string departureTime         // 出行时间
    5: optional string specialNeeds // 特殊需求（可选）
    6: optional string contactWay   // 联系方式（可选）
}

// 行程发布响应
struct TripPublishResp {
    1: i64 tripId                   // 行程ID
    2: string message               // 返回消息
    3: bool success                 // 是否成功
}

// ==================== 司机行程发布接口 ====================

// 司机行程发布请求
struct DriverTripPublishReq {
    1: i64 driverId                 // 司机ID
    2: string startPoint            // 起点
    3: string endPoint              // 终点
    4: string departureTime         // 发车时间
    5: string vehicleInfo           // 车辆信息
}

// 司机行程发布响应
struct DriverTripPublishResp {
    1: i64 tripId                   // 行程ID
    2: string message               // 返回消息
    3: bool success                 // 是否成功
}

// ==================== 行程查询接口 ====================

// 行程查询请求
struct TripQueryReq {
    1: string startPoint            // 起点
    2: string endPoint              // 终点
    3: string departureTime         // 出行时间
    4: optional string tripType     // 行程类型（passenger/driver，可选）
}

// 行程信息
struct TripInfo {
    1: i64 tripId                   // 行程ID
    2: string startPoint            // 起点
    3: string endPoint              // 终点
    4: string departureTime         // 出行时间
    5: string tripType              // 行程类型（passenger/driver）
    6: string publisherName         // 发布者姓名
    7: optional string vehicleInfo  // 车辆信息（司机行程）
    8: optional string specialNeeds // 特殊需求（乘客行程）
    9: string status                // 行程状态
}

// 行程查询响应
struct TripQueryResp {
    1: list<TripInfo> trips         // 行程列表
    2: string message               // 返回消息
    3: bool success                 // 是否成功
}

// ==================== 乘客求助接口 ====================

// 乘客求助请求
struct PassengerHelpReq {
    1: i64 tripId                   // 行程ID
    2: i64 passengerId              // 乘客ID
    3: string helpType              // 求助类型（官方/客服）
}

// 乘客求助响应
struct PassengerHelpResp {
    1: string contactInfo           // 联系结果（客服/官方接入状态）
    2: string message               // 返回消息
    3: bool success                 // 是否成功
}

// ==================== 行程分享接口 ====================

// 行程分享请求
struct TripShareReq {
    1: i64 tripId                   // 行程ID
    2: i64 userId                   // 用户ID
    3: list<string> shareTargets    // 分享对象（好友/家人）
}

// 行程分享响应
struct TripShareResp {
    1: string shareLink             // 分享链接或信息
    2: string message               // 返回消息
    3: bool success                 // 是否成功
}

// ==================== 行程详情接口 ====================

// 行程详情请求
struct TripDetailReq {
    1: i64 tripId                   // 行程ID
}

// 行程详情响应
struct TripDetailResp {
    1: TripInfo tripInfo            // 行程信息
    2: string message               // 返回消息
    3: bool success                 // 是否成功
}

// ==================== 数据导出接口 ====================

// 数据导出请求
struct DataExportReq {
    1: i64 userId                   // 用户ID
    2: string userType              // 用户类型（passenger/driver）
    3: list<string> exportFields    // 导出字段列表
    4: optional string startDate    // 开始日期（可选）
    5: optional string endDate      // 结束日期（可选）
}

// 数据导出响应
struct DataExportResp {
    1: string downloadUrl           // 下载链接
    2: string fileFormat            // 文件格式（csv/excel）
    3: i64 exportId                 // 导出记录ID
    4: string message               // 返回消息
    5: bool success                 // 是否成功
}

// 导出记录查询请求
struct ExportRecordQueryReq {
    1: i64 userId                   // 用户ID
    2: optional i32 limit           // 查询数量限制（可选）
}

// 导出记录信息
struct ExportRecordInfo {
    1: i64 exportId                 // 导出记录ID
    2: i64 userId                   // 用户ID
    3: string userType              // 用户类型
    4: string exportFields          // 导出字段
    5: string status                // 导出状态（pending/completed/failed）
    6: string downloadUrl           // 下载链接
    7: string createdAt             // 创建时间
    8: string expiresAt             // 过期时间
}

// 导出记录查询响应
struct ExportRecordQueryResp {
    1: list<ExportRecordInfo> records // 导出记录列表
    2: string message               // 返回消息
    3: bool success                 // 是否成功
}

// ==================== 服务定义 ====================

service OrderService {
    // 行程发布接口（乘客发布行程）
    TripPublishResp TripPublish(1: TripPublishReq req)
    
    // 司机行程发布接口
    DriverTripPublishResp DriverTripPublish(1: DriverTripPublishReq req)
    
    // 行程查询接口
    TripQueryResp TripQuery(1: TripQueryReq req)
    
    // 乘客求助接口
    PassengerHelpResp PassengerHelp(1: PassengerHelpReq req)
    
    // 行程分享接口
    TripShareResp TripShare(1: TripShareReq req)
    
    // 行程详情接口
    TripDetailResp TripDetail(1: TripDetailReq req)
    
    // 数据导出接口
    DataExportResp DataExport(1: DataExportReq req)
    
    // 导出记录查询接口
    ExportRecordQueryResp ExportRecordQuery(1: ExportRecordQueryReq req)
}