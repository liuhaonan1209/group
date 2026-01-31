// Package main 司机服务RPC处理器
// 实现司机相关的业务逻辑，包括注册、查询、身份验证等功能
package main

import (
	"context"
	"fmt"
	"group/handler/dao"
	"group/handler/model"
	"group/kitex_gen/car/driver"
	"group/pkg/idcard"
	"log"
	"regexp"
	"time"
)

// DriverServiceImpl 司机服务实现结构体
// 实现了IDL中定义的DriverService接口
type DriverServiceImpl struct{}

// DriverDetail 查询司机详情接口
// 根据司机ID查询司机的详细信息
// 参数:
//   - ctx: 上下文对象
//   - req: 查询请求，包含司机ID
//
// 返回:
//   - resp: 司机详情响应，包含姓名、手机号、身份证号、驾驶证号、注册时间、评分
//   - err: 错误信息
func (s *DriverServiceImpl) DriverDetail(ctx context.Context, req *driver.DriverDetailReq) (resp *driver.DriverDetailResp, err error) {

	// 查询司机信息
	var drivers model.Driver
	id, err := dao.GetDriverByID(&drivers, req.Id)
	if err != nil {
		return nil, fmt.Errorf("查询司机信息失败")
	}

	// 构造响应数据
	resp = &driver.DriverDetailResp{
		Name:         id.Name,
		Tel:          id.Tel,
		IDCard:       id.IDCard,
		License:      id.License,
		RegisterDate: id.RegisterDate.Format("2006-01-02"),
		Rating:       fmt.Sprintf("%.2f", id.Rating),
	}

	log.Printf("查询司机信息成功: %d", req.Id)
	return resp, nil
}

// DriverRegister 司机注册接口
// 处理司机注册请求，包括参数验证、身份证实名验证、数据库存储等
// 参数:
//   - ctx: 上下文对象
//   - req: 司机注册请求，包含姓名、手机号、身份证号、驾驶证号
//
// 返回:
//   - resp: 注册响应，包含注册结果、司机ID、提示信息
//   - err: 错误信息
func (s *DriverServiceImpl) DriverRegister(ctx context.Context, req *driver.DriverRegisterReq) (resp *driver.DriverRegisterResp, err error) {

	//  验证手机号格式
	if !validatePhone(req.Tel) {
		return &driver.DriverRegisterResp{
			Id:      0,
			Message: "手机号格式不正确",
			Success: false,
		}, nil
	}

	//  验证身份证号格式
	if !validateIDCard(req.IDCard) {
		return &driver.DriverRegisterResp{
			Id:      0,
			Message: "身份证号格式不正确",
			Success: false,
		}, nil
	}

	//  验证驾驶证号格式（通常与身份证号相同）
	if !validateLicense(req.License) {
		return &driver.DriverRegisterResp{
			Id:      0,
			Message: "驾驶证号格式不正确",
			Success: false,
		}, nil
	}

	//  调用身份证实名验证接口
	// 验证姓名和身份证号是否匹配
	err = idcard.Authorization(req.IDCard, req.Name)
	if err != nil {
		log.Printf("身份证实名验证失败: %v", err)
		return &driver.DriverRegisterResp{
			Id:      0,
			Message: fmt.Sprintf("身份证实名验证失败: %s", err.Error()),
			Success: false,
		}, nil
	}
	log.Printf("身份证实名验证通过: 姓名=%s", req.Name)

	//  检查手机号是否已注册
	var existingDriver model.Driver
	existingDriver, _ = dao.GetDriverByTel(&existingDriver, req.Tel)
	if existingDriver.ID != 0 {
		return &driver.DriverRegisterResp{
			Id:      0,
			Message: "该手机号已被注册",
			Success: false,
		}, nil
	}

	//  创建司机记录
	newDriver := model.Driver{
		Name:         req.Name,
		Tel:          req.Tel,
		IDCard:       req.IDCard,
		License:      req.License,
		RegisterDate: time.Now(),
		Rating:       5.00, // 默认评分5.00
	}

	err = dao.CreateDriver(&newDriver)
	if err != nil {
		log.Printf("创建司机失败: %v", err)
		return &driver.DriverRegisterResp{
			Id:      0,
			Message: "注册失败，手机号、身份证号或驾驶证号可能已被注册",
			Success: false,
		}, nil
	}

	log.Printf("司机注册成功: ID=%d", newDriver.ID)
	return &driver.DriverRegisterResp{
		Id:      int64(newDriver.ID),
		Message: "注册成功",
		Success: true,
	}, nil
}

// DriverVerify 司机身份验证接口
// 验证手机号、身份证号和驾驶证号是否匹配系统中已注册的司机信息
// 参数:
//   - ctx: 上下文对象
//   - req: 验证请求，包含手机号、身份证号、驾驶证号
//
// 返回:
//   - resp: 验证响应，包含验证结果、司机ID、提示信息
//   - err: 错误信息
func (s *DriverServiceImpl) DriverVerify(ctx context.Context, req *driver.DriverVerifyReq) (resp *driver.DriverVerifyResp, err error) {

	// . 验证手机号格式
	if !validatePhone(req.Tel) {
		return &driver.DriverVerifyResp{
			Verified: false,
			Message:  "手机号格式不正确",
		}, nil
	}

	//  验证身份证号格式
	if !validateIDCard(req.IDCard) {
		return &driver.DriverVerifyResp{
			Verified: false,
			Message:  "身份证号格式不正确",
		}, nil
	}

	// 验证驾驶证号格式
	if !validateLicense(req.License) {
		return &driver.DriverVerifyResp{
			Verified: false,
			Message:  "驾驶证号格式不正确",
		}, nil
	}

	// 根据手机号、身份证和驾驶证查询司机信息
	var d model.Driver
	driverDetail, err := dao.GetDriverByTelAndIDCardAndLicense(&d, req.Tel, req.IDCard, req.License)
	if err != nil {
		log.Printf("司机身份验证失败: %v", err)
		return &driver.DriverVerifyResp{
			Verified: false,
			Message:  "身份验证失败，未找到匹配的司机信息",
		}, nil
	}

	log.Printf("司机身份验证成功: ID=%d", driverDetail.ID)
	driverId := int64(driverDetail.ID)

	return &driver.DriverVerifyResp{
		Verified: true,
		Message:  "身份验证通过",
		DriverId: &driverId,
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

// validateLicense 验证驾驶证号格式
// 驾驶证号通常是18位，与身份证号格式相同
// 参数:
//   - license: 待验证的驾驶证号
//
// 返回:
//   - bool: 格式正确返回true，否则返回false
func validateLicense(license string) bool {
	// 驾驶证号通常是18位，与身份证号格式相同
	pattern := `^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`
	matched, _ := regexp.MatchString(pattern, license)
	return matched
}

// ==================== 司机配置管理接口实现 ====================

// DriverConfigQuery 司机配置查询接口
// 根据司机ID查询司机的配置信息，包括功能开关、有效期设置、计费规则、车辆合规性等
// 参数:
//   - ctx: 上下文对象
//   - req: 查询请求，包含司机ID
//
// 返回:
//   - resp: 配置查询响应，包含配置信息
//   - err: 错误信息
func (s *DriverServiceImpl) DriverConfigQuery(ctx context.Context, req *driver.DriverConfigQueryReq) (resp *driver.DriverConfigQueryResp, err error) {

	// . 验证司机是否存在
	var driverModel model.Driver
	driverInfo, err := dao.GetDriverByID(&driverModel, req.DriverId)
	if err != nil {
		log.Printf("查询司机信息失败: %v", err)
		return &driver.DriverConfigQueryResp{
			Config:  nil,
			Message: "司机不存在",
			Success: false,
		}, nil
	}

	// . 查询司机配置信息（调用DAO函数）
	var configModel model.DriverConfig
	configInfo, err := dao.GetDriverConfigByDriverID(&configModel, req.DriverId)
	if err != nil {
		// 如果配置不存在，创建默认配置
		log.Printf("司机配置不存在，创建默认配置: 司机ID=%d", req.DriverId)

		defaultConfig := model.DriverConfig{
			DriverID:        req.DriverId,
			CanPublishTrip:  true,  // 默认可以发布行程
			AutoNotify:      true,  // 默认开启自动通知
			MaxPassengers:   -1,    // 默认无限制
			StartPrice:      0.0,   // 默认起步价为0
			PricePerKm:      0.0,   // 默认每公里单价为0
			VehicleNumber:   "",    // 默认无车牌号
			VehicleType:     "",    // 默认无车型
			InsuranceExpiry: "",    // 默认无保险到期日期
			IsCompliant:     false, // 默认不合规
		}

		// 创建默认配置
		err = dao.CreateDriverConfig(&defaultConfig)
		if err != nil {
			log.Printf("创建默认配置失败: %v", err)
			return &driver.DriverConfigQueryResp{
				Config:  nil,
				Message: "查询配置失败",
				Success: false,
			}, nil
		}

		configInfo = defaultConfig
	}

	// 4. 构造响应数据
	configResp := &driver.DriverConfigInfo{
		DriverId:        configInfo.DriverID,
		CanPublishTrip:  configInfo.CanPublishTrip,
		AutoNotify:      configInfo.AutoNotify,
		MaxPassengers:   int32(configInfo.MaxPassengers),
		StartPrice:      configInfo.StartPrice,
		PricePerKm:      configInfo.PricePerKm,
		VehicleNumber:   configInfo.VehicleNumber,
		VehicleType:     configInfo.VehicleType,
		InsuranceExpiry: configInfo.InsuranceExpiry,
		IsCompliant:     configInfo.IsCompliant,
	}

	log.Printf("司机配置查询成功: 司机ID=%d, 司机姓名=%s", req.DriverId, driverInfo.Name)
	return &driver.DriverConfigQueryResp{
		Config:  configResp,
		Message: "查询成功",
		Success: true,
	}, nil
}

// VehicleComplianceCheck 车辆合规性验证接口
// 验证司机的车辆是否合规，包括车牌号、保险到期日期等
// 参数:
//   - ctx: 上下文对象
//   - req: 验证请求，包含司机ID、车牌号、保险到期日期
//
// 返回:
//   - resp: 验证响应，包含是否合规、不合规项列表
//   - err: 错误信息
func (s *DriverServiceImpl) VehicleComplianceCheck(ctx context.Context, req *driver.VehicleComplianceCheckReq) (resp *driver.VehicleComplianceCheckResp, err error) {

	// . 验证司机是否存在
	var driverModel model.Driver
	_, err = dao.GetDriverByID(&driverModel, req.DriverId)
	if err != nil {
		log.Printf("查询司机信息失败: %v", err)
		return &driver.VehicleComplianceCheckResp{
			IsCompliant: false,
			Issues:      []string{"司机不存在"},
			Message:     "司机不存在",
			Success:     false,
		}, nil
	}

	//  验证车牌号格式
	issues := []string{}
	if req.VehicleNumber == "" {
		issues = append(issues, "车牌号不能为空")
	} else if !validateVehicleNumber(req.VehicleNumber) {
		issues = append(issues, "车牌号格式不正确")
	}

	//  验证保险到期日期
	if req.InsuranceExpiry == "" {
		issues = append(issues, "保险到期日期不能为空")
	} else {
		// 解析保险到期日期
		expiryDate, err := time.Parse("2006-01-02", req.InsuranceExpiry)
		if err != nil {
			issues = append(issues, "保险到期日期格式不正确，应为：YYYY-MM-DD")
		} else {
			// 检查保险是否已过期
			if expiryDate.Before(time.Now()) {
				issues = append(issues, "保险已过期")
			}
		}
	}

	// . 判断是否合规
	isCompliant := len(issues) == 0

	// . 更新司机配置中的车辆合规性信息（调用DAO函数）
	var configModel model.DriverConfig
	updates := map[string]interface{}{
		"vehicle_number":   req.VehicleNumber,
		"insurance_expiry": req.InsuranceExpiry,
		"is_compliant":     isCompliant,
	}

	err = dao.UpdateDriverConfigFields(&configModel, req.DriverId, updates)
	if err != nil {
		log.Printf("更新车辆合规性信息失败: %v", err)
		// 不影响验证结果，继续返回
	}

	// 7. 构造响应
	message := "车辆合规"
	if !isCompliant {
		message = "车辆不合规"
	}

	log.Printf("车辆合规性验证完成: 司机ID=%d, 是否合规=%v", req.DriverId, isCompliant)
	return &driver.VehicleComplianceCheckResp{
		IsCompliant: isCompliant,
		Issues:      issues,
		Message:     message,
		Success:     true,
	}, nil
}

// validateVehicleNumber 验证车牌号格式
// 支持普通车牌和新能源车牌格式
// 规则：
//   - 普通车牌：省份简称 + 字母 + 5位数字或字母（如：京A12345）
//   - 新能源车牌：省份简称 + 字母 + 6位数字或字母（如：京AD12345）
//
// 参数:
//   - vehicleNumber: 待验证的车牌号
//
// 返回:
//   - bool: 格式正确返回true，否则返回false
func validateVehicleNumber(vehicleNumber string) bool {
	// 普通车牌：省份简称 + 字母 + 5位数字或字母
	pattern1 := `^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领][A-Z][A-HJ-NP-Z0-9]{5}$`
	// 新能源车牌：省份简称 + 字母 + 6位数字或字母
	pattern2 := `^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领][A-Z][A-HJ-NP-Z0-9]{6}$`

	matched1, _ := regexp.MatchString(pattern1, vehicleNumber)
	matched2, _ := regexp.MatchString(pattern2, vehicleNumber)

	return matched1 || matched2
}
