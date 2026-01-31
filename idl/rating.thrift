namespace go car.rating

// ==================== 评价提交接口 ====================

// 评价提交请求
struct RatingSubmitReq {
    1: i64 orderId                  // 订单ID
    2: i64 userId                   // 用户ID
    3: double score                 // 评分（1-5分）
    4: list<string> tags            // 标签（服务好、车辆干净等）
    5: string comment               // 文字意见
    6: string raterType             // 评价人类型（passenger/driver）
}

// 评价提交响应
struct RatingSubmitResp {
    1: i64 ratingId                 // 评价ID
    2: string message               // 返回消息
    3: bool success                 // 是否成功
}

// ==================== 评价查询接口 ====================

// 评价查询请求
struct RatingQueryReq {
    1: optional string timeRange    // 评价时间范围（可选）
    2: optional i64 routeId         // 线路ID（可选）
    3: optional i64 driverId        // 司机ID（可选）
    4: optional string driverPhone  // 司机手机号（可选）
    5: optional i32 page            // 页码（可选，默认1）
    6: optional i32 pageSize        // 每页数量（可选，默认10）
}

// 评价信息
struct RatingInfo {
    1: i64 ratingId                 // 评价ID
    2: i64 orderId                  // 订单ID
    3: i64 userId                   // 用户ID
    4: string userName              // 用户姓名
    5: double score                 // 评分
    6: list<string> tags            // 标签
    7: string comment               // 文字意见
    8: string raterType             // 评价人类型
    9: string createdAt             // 评价时间
    10: optional i64 driverId       // 司机ID
    11: optional string driverName  // 司机姓名
}

// 评价查询响应
struct RatingQueryResp {
    1: list<RatingInfo> ratings     // 评价列表
    2: i32 total                    // 总数量
    3: i32 page                     // 当前页码
    4: i32 pageSize                 // 每页数量
    5: string message               // 返回消息
    6: bool success                 // 是否成功
}

// ==================== 订单历史变更接口 ====================

// 订单历史变更查询请求
struct OrderHistoryQueryReq {
    1: i64 orderId                  // 订单ID
}

// 订单历史变更记录
struct OrderHistoryRecord {
    1: i64 historyId                // 历史记录ID
    2: i64 orderId                  // 订单ID
    3: string status                // 状态
    4: string reason                // 原因
    5: string changedAt             // 变更时间
    6: optional i64 operatorId      // 操作人ID
    7: optional string operatorName // 操作人姓名
}

// 订单历史变更查询响应
struct OrderHistoryQueryResp {
    1: list<OrderHistoryRecord> records // 变更记录列表
    2: string message               // 返回消息
    3: bool success                 // 是否成功
}

// ==================== 购票记录查询接口 ====================

// 购票记录查询请求
struct TicketRecordQueryReq {
    1: i64 userId                   // 用户ID
    2: optional string timeRange    // 时间范围（可选）
    3: optional i32 page            // 页码（可选，默认1）
    4: optional i32 pageSize        // 每页数量（可选，默认10）
}

// 购票记录信息
struct TicketRecordInfo {
    1: i64 orderId                  // 订单ID
    2: string orderNo               // 订单号
    3: i64 userId                   // 用户ID
    4: string userName              // 用户姓名
    5: string startPoint            // 起点
    6: string endPoint              // 终点
    7: string departureTime         // 出发时间
    8: double price                 // 价格
    9: string status                // 订单状态
    10: string createdAt            // 购票时间
}

// 购票记录查询响应
struct TicketRecordQueryResp {
    1: list<TicketRecordInfo> records // 购票记录列表
    2: i32 total                    // 总数量
    3: i32 page                     // 当前页码
    4: i32 pageSize                 // 每页数量
    5: string message               // 返回消息
    6: bool success                 // 是否成功
}

// ==================== 服务定义 ====================

service RatingService {
    // 评价提交接口
    RatingSubmitResp RatingSubmit(1: RatingSubmitReq req)
    
    // 评价查询接口
    RatingQueryResp RatingQuery(1: RatingQueryReq req)
    
    // 订单历史变更查询接口
    OrderHistoryQueryResp OrderHistoryQuery(1: OrderHistoryQueryReq req)
    
    // 购票记录查询接口
    TicketRecordQueryResp TicketRecordQuery(1: TicketRecordQueryReq req)
}
