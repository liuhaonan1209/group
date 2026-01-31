namespace go car.driver

struct DriverDetailReq{
    1: i64 id
}

struct DriverDetailResp {
    1: string name
    2: string tel
    3: string IDCard
    4: string license
    5: string registerDate
    6: string rating
}

// 司机注册请求
struct DriverRegisterReq {
    1: string name          // 姓名
    2: string tel           // 手机号
    3: string IDCard        // 身份证号
    4: string license       // 驾驶证号
    5: optional string drivingLicense // 行驶证（可选）
    6: optional string vehicleInfo    // 车辆信息（可选）
    7: optional string photo          // 证件照片（可选）
}

// 司机注册响应
struct DriverRegisterResp {
    1: i64 id               // 司机ID
    2: string message       // 返回消息
    3: bool success         // 是否成功
}

// 司机身份验证请求
struct DriverVerifyReq {
    1: string tel           // 手机号
    2: string IDCard        // 身份证号
    3: string license       // 驾驶证号
    4: optional string drivingLicense // 行驶证（可选）
    5: optional string vehicleInfo    // 车辆信息（可选）
    6: optional string photo          // 证件照片（可选）
}

// 司机身份验证响应
struct DriverVerifyResp {
    1: bool verified        // 是否验证通过
    2: string message       // 验证消息
    3: optional i64 driverId // 司机ID（验证通过时返回）
}

// ==================== 司机配置管理接口 ====================

// 司机配置查询请求
struct DriverConfigQueryReq {
    1: i64 driverId         // 司机ID
}

// 司机配置信息
struct DriverConfigInfo {
    1: i64 driverId         // 司机ID
    2: bool canPublishTrip  // 是否可以发布行程（功能开关）
    3: bool autoNotify      // 是否自动通知（用户购票时触发）
    4: i32 maxPassengers    // 有效期设置：最多可配置乘客数（-1表示无限制）
    5: double startPrice    // 计费规则：车辆起步价
    6: double pricePerKm    // 计费规则：每公里单价
    7: string vehicleNumber // 车辆合规性：车牌号
    8: string vehicleType   // 车辆合规性：车型
    9: string insuranceExpiry // 车辆合规性：保险到期日期
    10: bool isCompliant    // 车辆是否合规
}

// 司机配置查询响应
struct DriverConfigQueryResp {
    1: optional DriverConfigInfo config // 配置信息
    2: string message       // 返回消息
    3: bool success         // 是否成功
}

// 车辆合规性验证请求
struct VehicleComplianceCheckReq {
    1: i64 driverId         // 司机ID
    2: string vehicleNumber // 车牌号
    3: string insuranceExpiry // 保险到期日期
}

// 车辆合规性验证响应
struct VehicleComplianceCheckResp {
    1: bool isCompliant     // 是否合规
    2: list<string> issues  // 不合规项列表
    3: string message       // 返回消息
    4: bool success         // 是否成功
}

service DriverService {
    DriverDetailResp DriverDetail(1:DriverDetailReq req)
    DriverRegisterResp DriverRegister(1:DriverRegisterReq req)
    DriverVerifyResp DriverVerify(1:DriverVerifyReq req)
    
    // 司机配置管理接口
    DriverConfigQueryResp DriverConfigQuery(1:DriverConfigQueryReq req)
    VehicleComplianceCheckResp VehicleComplianceCheck(1:VehicleComplianceCheckReq req)
}