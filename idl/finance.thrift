namespace go car.finance

// ========== 收支对账表 ==========
struct BalanceSheetReq {
    1: string startDate      // 开始日期
    2: string endDate        // 结束日期
    3: i32 status            // 账目状态：0-全部，1-正常，2-异常
    4: i32 page
    5: i32 size
}

struct BalanceSheetItem {
    1: i64 id
    2: string settleDate     // 对账日期
    3: double orderIncome    // 订单实收
    4: double orderRefund    // 订单退款
    5: double orderSettle    // 订单结算
    6: i32 status            // 账目状态：1-正常，2-异常
    7: string statusText     // 状态文本
}

struct BalanceSheetResp {
    1: bool Success
    2: string Msg
    3: list<BalanceSheetItem> list
    4: i64 total
    5: double totalIncome    // 总收入
    6: double totalRefund    // 总退款
    7: double totalSettle    // 总结算
}

// ========== 收入对账表 ==========
struct IncomeSheetReq {
    1: string startDate
    2: string endDate
    3: i32 status
    4: i32 page
    5: i32 size
}

struct IncomeSheetItem {
    1: i64 id
    2: string settleDate         // 结算日期(发车时间)
    3: string routeNo            // 线路编号
    4: string routeName          // 线路名称
    5: string fleet              // 所属车队
    6: double ticketPrice        // 票价
    7: double discountAmount     // 优惠金额
    8: double checkedIncome      // 实收金额(已检票)
    9: i32 checkedTickets        // 已检票数
    10: double uncheckedIncome   // 实收金额(未检票)
    11: i32 uncheckedTickets     // 未检票数
    12: double checkedRefund     // 退款金额(已检票)
    13: i32 checkedRefundTickets // 已检票退款数
    14: double uncheckedRefund   // 退款金额(未检票)
    15: i32 uncheckedRefundTickets // 未检票退款数
    16: i32 status
}

struct IncomeSheetResp {
    1: bool Success
    2: string Msg
    3: list<IncomeSheetItem> list
    4: i64 total
}

// ========== 线路结算表 ==========
struct RouteSettleReq {
    1: string startDate
    2: string endDate
    3: i64 routeId
    4: string fleet
    5: i32 page
    6: i32 size
}

struct RouteSettleItem {
    1: i64 id
    2: string settleDate
    3: string routeNo
    4: string routeName
    5: string fleet
    6: double ticketPrice
    7: double discountAmount
    8: double checkedIncome
    9: i32 checkedTickets
    10: double uncheckedIncome
    11: i32 uncheckedTickets
    12: double checkedRefund
    13: i32 checkedRefundTickets
    14: double uncheckedRefund
    15: i32 uncheckedRefundTickets
}

struct RouteSettleResp {
    1: bool Success
    2: string Msg
    3: list<RouteSettleItem> list
    4: i64 total
}

// ========== 班次结算表 ==========
struct ScheduleSettleReq {
    1: string startDate
    2: string endDate
    3: i64 routeId
    4: i64 scheduleId
    5: i32 page
    6: i32 size
}

struct ScheduleSettleItem {
    1: i64 id
    2: string settleDate
    3: string routeNo
    4: string routeName
    5: string departureTime      // 发车时间
    6: double ticketPrice
    7: double discountAmount
    8: double checkedIncome
    9: i32 checkedTickets
    10: double uncheckedIncome
    11: i32 uncheckedTickets
    12: double checkedRefund
    13: i32 checkedRefundTickets
    14: double uncheckedRefund
    15: i32 uncheckedRefundTickets
}

struct ScheduleSettleResp {
    1: bool Success
    2: string Msg
    3: list<ScheduleSettleItem> list
    4: i64 total
}

// ========== 站点结算表 ==========
struct StationSettleReq {
    1: string startDate
    2: string endDate
    3: i64 stationId
    4: i32 page
    5: i32 size
}

struct StationSettleItem {
    1: i64 id
    2: string settleDate
    3: string stationName        // 上车站点
    4: double ticketPrice
    5: double discountAmount
    6: double checkedIncome
    7: i32 checkedTickets
    8: double uncheckedIncome
    9: i32 uncheckedTickets
    10: double checkedRefund
    11: i32 checkedRefundTickets
    12: double uncheckedRefund
    13: i32 uncheckedRefundTickets
}

struct StationSettleResp {
    1: bool Success
    2: string Msg
    3: list<StationSettleItem> list
    4: i64 total
}

// ========== 司机结算表 ==========
struct DriverSettleReq {
    1: string startDate
    2: string endDate
    3: i64 driverId
    4: string driverName
    5: i32 page
    6: i32 size
}

struct DriverSettleItem {
    1: i64 id
    2: string settleDate
    3: i64 driverId
    4: string driverName
    5: string phone
    6: i32 tripCount             // 行程数
    7: double totalIncome        // 总收入
    8: double commission         // 佣金
    9: double netIncome          // 净收入
    10: i32 status               // 结算状态：1-待结算，2-已结算
}

struct DriverSettleResp {
    1: bool Success
    2: string Msg
    3: list<DriverSettleItem> list
    4: i64 total
}

// ========== 交易明细 ==========
struct TransactionReq {
    1: string startDate
    2: string endDate
    3: string orderNo            // 订单号
    4: i32 paymentMethod         // 支付方式：0-全部，1-微信，2-支付宝，3-银行卡
    5: i32 transType             // 交易类型：0-全部，1-收入，2-支出，3-退款
    6: i32 page
    7: i32 size
}

struct TransactionItem {
    1: i64 id
    2: string transDate          // 交易日期
    3: string transTime          // 交易时间
    4: string orderNo            // 订单号
    5: i64 driverId
    6: string driverName
    7: i64 passengerId
    8: string passengerName
    9: string startStation       // 起点
    10: string endStation        // 终点
    11: double amount            // 金额
    12: i32 paymentMethod        // 支付方式
    13: string paymentMethodText
    14: i32 transType            // 交易类型
    15: string transTypeText
    16: string remark            // 备注
}

struct TransactionResp {
    1: bool Success
    2: string Msg
    3: list<TransactionItem> list
    4: i64 total
}

// ========== 异常交易 ==========
struct AbnormalTransReq {
    1: string startDate
    2: string endDate
    3: i32 abnormalType          // 异常类型：0-全部，1-重复扣款，2-金额不符，3-其他
    4: i32 status                // 处理状态：0-全部，1-待处理，2-已处理
    5: i32 page
    6: i32 size
}

struct AbnormalTransItem {
    1: i64 id
    2: string transDate
    3: string orderNo
    4: i32 abnormalType
    5: string abnormalTypeText
    6: double expectedAmount     // 应收金额
    7: double actualAmount       // 实收金额
    8: double diffAmount         // 差异金额
    9: i32 status
    10: string statusText
    11: string handler           // 处理人
    12: string handleTime        // 处理时间
    13: string handleRemark      // 处理备注
}

struct AbnormalTransResp {
    1: bool Success
    2: string Msg
    3: list<AbnormalTransItem> list
    4: i64 total
}

// ========== 处理异常交易 ==========
struct HandleAbnormalReq {
    1: i64 id
    2: string handler
    3: string remark
    4: i32 handleType            // 处理方式：1-退款，2-补收，3-忽略
}

struct HandleAbnormalResp {
    1: bool Success
    2: string Msg
}

// ========== 生成账单 ==========
struct GenerateBillReq {
    1: i32 billType              // 账单类型：1-日账单，2-周账单，3-月账单
    2: string startDate
    3: string endDate
}

struct GenerateBillResp {
    1: bool Success
    2: string Msg
    3: i64 billId
}

// ========== 导出 ==========
struct ExportReq {
    1: i32 exportType            // 导出类型：1-收支对账，2-收入对账，3-线路结算，4-班次结算，5-站点结算，6-司机结算
    2: string startDate
    3: string endDate
}

struct ExportResp {
    1: bool Success
    2: string Msg
    3: string fileUrl
}

service FinanceService {
    // 收支对账
    BalanceSheetResp GetBalanceSheet(1: BalanceSheetReq req)
    
    // 收入对账
    IncomeSheetResp GetIncomeSheet(1: IncomeSheetReq req)
    
    // 线路结算
    RouteSettleResp GetRouteSettle(1: RouteSettleReq req)
    
    // 班次结算
    ScheduleSettleResp GetScheduleSettle(1: ScheduleSettleReq req)
    
    // 站点结算
    StationSettleResp GetStationSettle(1: StationSettleReq req)
    
    // 司机结算
    DriverSettleResp GetDriverSettle(1: DriverSettleReq req)
    
    // 交易明细
    TransactionResp GetTransactions(1: TransactionReq req)
    
    // 异常交易
    AbnormalTransResp GetAbnormalTrans(1: AbnormalTransReq req)
    HandleAbnormalResp HandleAbnormal(1: HandleAbnormalReq req)
    
    // 生成账单
    GenerateBillResp GenerateBill(1: GenerateBillReq req)
    
    // 导出
    ExportResp Export(1: ExportReq req)
}
