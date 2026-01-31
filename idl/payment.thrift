namespace go car.payment

// ==================== 支付接口 ====================

// 支付请求
struct PaymentReq {
    1: i64 orderId                  // 订单ID
    2: i64 userId                   // 用户ID（支付信息，如加密卡号）
    3: string paymentMethod         // 支付方式（alipay/wechat/card）
    4: optional string cardInfo     // 加密的卡号信息（可选）
}

// 支付响应
struct PaymentResp {
    1: i64 paymentId                // 支付ID
    2: string paymentNo             // 支付单号
    3: string status                // 支付状态（success/failed/pending）
    4: string message               // 返回消息
    5: bool success                 // 是否成功
}

// ==================== 支付问题处理接口 ====================

// 支付问题处理请求
struct PaymentIssueReq {
    1: i64 orderId                  // 订单ID
    2: i64 userId                   // 用户ID
    3: string issueType             // 问题类型（payment_failed/refund_issue/other）
    4: string description           // 问题描述
}

// 支付问题处理响应
struct PaymentIssueResp {
    1: i64 issueId                  // 问题ID
    2: string solution              // 解决方案
    3: string message               // 返回消息
    4: bool success                 // 是否成功
}

// ==================== 支付记录查询接口 ====================

// 支付记录查询请求
struct PaymentRecordQueryReq {
    1: optional i64 userId          // 用户ID（可选）
    2: optional i64 orderId         // 订单ID（可选）
    3: optional string status       // 支付状态（可选）
    4: optional string startDate    // 开始日期（可选）
    5: optional string endDate      // 结束日期（可选）
    6: optional i32 page            // 页码（可选，默认1）
    7: optional i32 pageSize        // 每页数量（可选，默认10）
}

// 支付记录信息
struct PaymentRecordInfo {
    1: i64 paymentId                // 支付ID
    2: string paymentNo             // 支付单号
    3: i64 orderId                  // 订单ID
    4: i64 userId                   // 用户ID
    5: string paymentMethod         // 支付方式
    6: double amount                // 支付金额
    7: string status                // 支付状态
    8: string paymentTime           // 支付时间
    9: optional string failReason   // 失败原因
}

// 支付记录查询响应
struct PaymentRecordQueryResp {
    1: list<PaymentRecordInfo> records // 支付记录列表
    2: i32 total                    // 总数量
    3: i32 page                     // 当前页码
    4: i32 pageSize                 // 每页数量
    5: string message               // 返回消息
    6: bool success                 // 是否成功
}

// ==================== 支付状态查询接口 ====================

// 支付状态查询请求
struct PaymentStatusQueryReq {
    1: i64 paymentId                // 支付ID
    2: i64 userId                   // 用户ID
}

// 支付状态查询响应
struct PaymentStatusQueryResp {
    1: optional PaymentRecordInfo paymentInfo // 支付信息
    2: string message               // 返回消息
    3: bool success                 // 是否成功
}

// ==================== 服务定义 ====================

service PaymentService {
    // 支付接口
    PaymentResp Payment(1: PaymentReq req)
    
    // 支付问题处理接口
    PaymentIssueResp PaymentIssue(1: PaymentIssueReq req)
    
    // 支付记录查询接口
    PaymentRecordQueryResp PaymentRecordQuery(1: PaymentRecordQueryReq req)
    
    // 支付状态查询接口
    PaymentStatusQueryResp PaymentStatusQuery(1: PaymentStatusQueryReq req)
}
