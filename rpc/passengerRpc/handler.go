// Package main 乘客服务RPC处理器
// 实现乘客相关的业务逻辑，包括注册、查询、身份验证等功能
package main

import (
	"context"
	"fmt"
	"group/handler/dao"
	"group/handler/model"
	passenger "group/kitex_gen/car/passenger"
	"group/pkg/idcard"
	"group/pkg/verifycode"
	"log"
	"regexp"
	"time"
)

// PassengerServiceImpl 乘客服务实现结构体
// 实现了IDL中定义的PassengerService接口
type PassengerServiceImpl struct{}

// PassengerRegister 乘客注册接口
// 处理乘客注册请求，包括参数验证、身份证实名验证、数据库存储等
// 参数:
//   - ctx: 上下文对象
//   - req: 乘客注册请求，包含姓名、手机号、身份证号
//
// 返回:
//   - resp: 注册响应，包含注册结果、乘客ID、提示信息
//   - err: 错误信息
func (s *PassengerServiceImpl) PassengerRegister(ctx context.Context, req *passenger.PassengerRegisterReq) (resp *passenger.PassengerRegisterResp, err error) {
	//  验证手机号格式
	if !validatePhone(req.Tel) {
		return &passenger.PassengerRegisterResp{
			Id:      0,
			Message: "手机号格式不正确",
			Success: false,
		}, nil
	}

	// 3. 验证身份证号格式
	if !validateIDCard(req.IDCard) {
		return &passenger.PassengerRegisterResp{
			Id:      0,
			Message: "身份证号格式不正确",
			Success: false,
		}, nil
	}

	// 4. 调用身份证实名验证接口
	// 验证姓名和身份证号是否匹配
	err = idcard.Authorization(req.IDCard, req.Name)
	if err != nil {
		log.Printf("身份证实名验证失败: %v", err)
		return &passenger.PassengerRegisterResp{
			Id:      0,
			Message: fmt.Sprintf("身份证实名验证失败: %s", err.Error()),
			Success: false,
		}, nil
	}
	log.Printf("身份证实名验证通过: 姓名=%s", req.Name)

	// 5. 检查手机号是否已注册
	var existingPassenger model.Passenger
	existingPassenger, _ = dao.GetPassengerByTel(&existingPassenger, req.Tel)
	if existingPassenger.ID != 0 {
		return &passenger.PassengerRegisterResp{
			Id:      0,
			Message: "该手机号已被注册",
			Success: false,
		}, nil
	}

	// 6. 创建乘客记录
	newPassenger := model.Passenger{
		Name:         req.Name,
		Tel:          req.Tel,
		IDCard:       req.IDCard,
		RegisterDate: time.Now(),
	}

	err = dao.CreatePassenger(&newPassenger)
	if err != nil {
		log.Printf("创建乘客失败: %v", err)
		return &passenger.PassengerRegisterResp{
			Id:      0,
			Message: "注册失败，手机号或身份证号可能已被注册",
			Success: false,
		}, nil
	}

	log.Printf("乘客注册成功: ID=%d", newPassenger.ID)
	return &passenger.PassengerRegisterResp{
		Id:      int64(newPassenger.ID),
		Message: "注册成功",
		Success: true,
	}, nil
}

// PassengerDetail 查询乘客详情接口
// 根据乘客ID查询乘客的详细信息
// 参数:
//   - ctx: 上下文对象
//   - req: 查询请求，包含乘客ID
//
// 返回:
//   - resp: 乘客详情响应，包含姓名、手机号、身份证号、注册时间
//   - err: 错误信息
func (s *PassengerServiceImpl) PassengerDetail(ctx context.Context, req *passenger.PassengerDetailReq) (resp *passenger.PassengerDetailResp, err error) {
	log.Printf("查询乘客ID: %d", req.Id)

	// 根据ID查询乘客信息
	var p model.Passenger
	id, err := dao.GetPassengerByID(&p, req.Id)
	if err != nil {
		return nil, fmt.Errorf("查询乘客失败")
	}

	// 构造响应数据
	resp = &passenger.PassengerDetailResp{
		Name:         id.Name,
		Tel:          id.Tel,
		IDCard:       id.IDCard,
		RegisterDate: id.RegisterDate.Format("2006-01-02"),
	}

	log.Printf("查询乘客信息成功: %d", req.Id)
	return resp, nil
}

// PassengerVerify 乘客身份验证接口
// 验证手机号和验证码是否匹配
// 参数:
//   - ctx: 上下文对象
//   - req: 验证请求，包含手机号和验证码
//
// 返回:
//   - resp: 验证响应，包含验证结果、乘客ID、提示信息
//   - err: 错误信息
func (s *PassengerServiceImpl) PassengerVerify(ctx context.Context, req *passenger.PassengerVerifyReq) (resp *passenger.PassengerVerifyResp, err error) {
	log.Printf("乘客身份验证请求: 手机号=%s", req.Tel)
	// . 验证手机号格式
	if !validatePhone(req.Tel) {
		return &passenger.PassengerVerifyResp{
			Verified: false,
			Message:  "手机号格式不正确",
		}, nil
	}

	// 3. 验证验证码
	if !verifycode.Verify(req.Tel, req.VerifyCode) {
		log.Printf("验证码验证失败: 手机号=%s", req.Tel)
		return &passenger.PassengerVerifyResp{
			Verified: false,
			Message:  "验证码错误或已过期",
		}, nil
	}

	// 4. 根据手机号查询乘客信息
	var p model.Passenger
	passengerDetail, err := dao.GetPassengerByTel(&p, req.Tel)
	if err != nil {
		log.Printf("乘客身份验证失败: %v", err)
		return &passenger.PassengerVerifyResp{
			Verified: false,
			Message:  "身份验证失败，未找到匹配的乘客信息",
		}, nil
	}

	log.Printf("乘客身份验证成功: ID=%d", passengerDetail.ID)
	passengerId := int64(passengerDetail.ID)

	return &passenger.PassengerVerifyResp{
		Verified:    true,
		Message:     "身份验证通过",
		PassengerId: &passengerId,
	}, nil
}

// SendVerifyCode 发送验证码接口
// 向指定手机号发送验证码
// 参数:
//   - ctx: 上下文对象
//   - req: 发送验证码请求，包含手机号
//
// 返回:
//   - resp: 发送响应，包含发送结果、提示信息
//   - err: 错误信息
func (s *PassengerServiceImpl) SendVerifyCode(ctx context.Context, req *passenger.SendVerifyCodeReq) (resp *passenger.SendVerifyCodeResp, err error) {
	log.Printf("发送验证码请求: 手机号=%s", req.Tel)

	// . 验证手机号格式
	if !validatePhone(req.Tel) {
		return &passenger.SendVerifyCodeResp{
			Success: false,
			Message: "手机号格式不正确",
		}, nil
	}

	// 3. 检查手机号是否已注册
	var p model.Passenger
	_, err = dao.GetPassengerByTel(&p, req.Tel)
	if err != nil {
		return &passenger.SendVerifyCodeResp{
			Success: false,
			Message: "该手机号未注册",
		}, nil
	}

	// 4. 发送验证码
	code, err := verifycode.SendCode(req.Tel)
	if err != nil {
		log.Printf("发送验证码失败: %v", err)
		return &passenger.SendVerifyCodeResp{
			Success: false,
			Message: "发送验证码失败",
		}, nil
	}

	log.Printf("验证码发送成功: 手机号=%s, 验证码=%s", req.Tel, code)
	return &passenger.SendVerifyCodeResp{
		Success: true,
		Message: fmt.Sprintf("验证码已发送，5分钟内有效（测试环境验证码：%s）", code),
	}, nil
}

// validatePhone 验证手机号格式
// 使用正则表达式验证中国大陆手机号格式
// 规则：1开头，第二位是3-9，共11位数字
// 参数:
//   - phone: 待验证的手机号
//
// 返回:
//   - bool: 格式正确返回true，否则返回false
func validatePhone(phone string) bool {
	// 中国大陆手机号正则：1开头，第二位是3-9，共11位
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// validateIDCard 验证身份证号格式
// 使用正则表达式验证18位身份证号格式
// 规则：
//   - 前6位为地区码
//   - 第7-14位为出生日期（年月日）
//   - 第15-17位为顺序码
//   - 第18位为校验码（数字或X）
//
// 参数:
//   - idCard: 待验证的身份证号
//
// 返回:
//   - bool: 格式正确返回true，否则返回false
func validateIDCard(idCard string) bool {
	// 18位身份证号正则
	pattern := `^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`
	matched, _ := regexp.MatchString(pattern, idCard)
	return matched
}
