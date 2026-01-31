namespace go car.passenger

struct PassengerDetailReq{
    1: i64 id
}

struct PassengerDetailResp {
    1: string name
    2: string tel
    3: string IDCard
    4: string registerDate
}

// 乘客注册请求
struct PassengerRegisterReq {
    1: string name          // 姓名
    2: string tel           // 手机号
    3: string IDCard        // 身份证号
    4: optional string photo // 证件照片（可选）
}

// 乘客注册响应
struct PassengerRegisterResp {
    1: i64 id               // 乘客ID
    2: string message       // 返回消息
    3: bool success         // 是否成功
}

// 乘客身份验证请求
struct PassengerVerifyReq {
    1: string tel           // 手机号
    2: string verifyCode    // 验证码
}

// 乘客身份验证响应
struct PassengerVerifyResp {
    1: bool verified        // 是否验证通过
    2: string message       // 验证消息
    3: optional i64 passengerId // 乘客ID（验证通过时返回）
}

// 发送验证码请求
struct SendVerifyCodeReq {
    1: string tel           // 手机号
}

// 发送验证码响应
struct SendVerifyCodeResp {
    1: bool success         // 是否发送成功
    2: string message       // 返回消息
}

service PassengerService {
    PassengerDetailResp PassengerDetail(1:PassengerDetailReq req)
    PassengerRegisterResp PassengerRegister(1:PassengerRegisterReq req)
    PassengerVerifyResp PassengerVerify(1:PassengerVerifyReq req)
    SendVerifyCodeResp SendVerifyCode(1:SendVerifyCodeReq req)
}