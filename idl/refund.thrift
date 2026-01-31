namespace go car.refund

// ==================== 退票申请接口 ====================

// 退票申请请求
struct RefundApplyReq {
    1: i64 orderId                  // 订单ID
    2: i64 userId                   // 用户ID
    3: string userType              // 用户类型（passenger/driver）
    4: string reason                // 退票原因
    5: optional string contactWay   // 联系方式（可选）
}

// 退票申请响应
struct RefundApplyResp {
    1: i64 refundId                 // 退票申请ID
    2: string message               // 返回消息
    3: bool success                 // 是否成功
}

// ==================== 退票审核接口 ====================

// 退票审核请求
struct RefundAuditReq {
    1: i64 refundId                 // 退票申请ID
    2: i64 auditUserId              // 审核人ID
    3: string auditResult           // 审核结果（approved/rejected）
    4: optional string auditReason  // 审核意见（可选）
}

// 退票审核响应
struct RefundAuditResp {
    1: string message               // 返回消息
    2: bool success                 // 是否成功
    3: optional string refundAmount // 退款金额（审核通过时返回）
}

// ==================== 批量退票接口 ====================

// 批量退票请求
struct BatchRefundReq {
    1: list<i64> orderIds           // 订单ID列表
    2: i64 operatorId               // 操作人ID
    3: string refundType            // 退票类型（full/partial）
    4: string reason                // 批量退票原因
    5: optional double refundRatio  // 退款比例（部分退款时使用）
}

// 批量退票结果项
struct BatchRefundResultItem {
    1: i64 orderId                  // 订单ID
    2: bool success                 // 是否成功
    3: string message               // 处理结果消息
    4: optional double refundAmount // 退款金额
}

// 批量退票响应
struct BatchRefundResp {
    1: list<BatchRefundResultItem> results // 处理结果列表
    2: i32 successCount             // 成功数量
    3: i32 failCount                // 失败数量
    4: string message               // 总体消息
    5: bool success                 // 整体是否成功
}

// ==================== 退票状态查询接口 ====================

// 退票状态查询请求
struct RefundStatusQueryReq {
    1: i64 refundId                 // 退票申请ID
    2: i64 userId                   // 用户ID
}

// 退票状态信息
struct RefundStatusInfo {
    1: i64 refundId                 // 退票申请ID
    2: i64 orderId                  // 订单ID
    3: i64 userId                   // 用户ID
    4: string userType              // 用户类型
    5: string reason                // 退票原因
    6: string status                // 退票状态（pending/approved/rejected/completed）
    7: optional string auditReason  // 审核意见
    8: optional double refundAmount // 退款金额
    9: string applyTime             // 申请时间
    10: optional string auditTime   // 审核时间
    11: optional string completedTime // 完成时间
}

// 退票状态查询响应
struct RefundStatusQueryResp {
    1: optional RefundStatusInfo refundInfo // 退票信息
    2: string message               // 返回消息
    3: bool success                 // 是否成功
}

// ==================== 退票记录查询接口 ====================

// 退票记录查询请求
struct RefundRecordQueryReq {
    1: optional i64 userId          // 用户ID（可选）
    2: optional string userType     // 用户类型（可选）
    3: optional string status       // 状态筛选（可选）
    4: optional string startDate    // 开始日期（可选）
    5: optional string endDate      // 结束日期（可选）
    6: optional i32 page            // 页码（可选，默认1）
    7: optional i32 pageSize        // 每页数量（可选，默认10）
}

// 退票记录查询响应
struct RefundRecordQueryResp {
    1: list<RefundStatusInfo> records // 退票记录列表
    2: i32 total                    // 总数量
    3: i32 page                     // 当前页码
    4: i32 pageSize                 // 每页数量
    5: string message               // 返回消息
    6: bool success                 // 是否成功
}

// ==================== 服务定义 ====================

service RefundService {
    // 退票申请接口
    RefundApplyResp RefundApply(1: RefundApplyReq req)
    
    // 退票审核接口
    RefundAuditResp RefundAudit(1: RefundAuditReq req)
    
    // 批量退票接口
    BatchRefundResp BatchRefund(1: BatchRefundReq req)
    
    // 退票状态查询接口
    RefundStatusQueryResp RefundStatusQuery(1: RefundStatusQueryReq req)
    
    // 退票记录查询接口
    RefundRecordQueryResp RefundRecordQuery(1: RefundRecordQueryReq req)
}