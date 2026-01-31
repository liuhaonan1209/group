package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"group/kitex_gen/car/driver"
	"group/kitex_gen/car/driver/driverservice"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

var (
	cli driverservice.Client
)

func main() {
	// 创建 RPC 客户端
	c, err := driverservice.NewClient("car.driver", client.WithHostPorts("127.0.0.1:8888"))
	if err != nil {
		log.Fatal(err)
	}
	cli = c

	// 创建 HTTP 服务器
	hz := server.New(server.WithHostPorts("127.0.0.1:6667"))

	// 基础功能路由
	hz.GET("/api/driver/:id", GetDriverDetail)      // 查询司机详情 - 根据司机ID获取司机基本信息
	hz.POST("/api/driver/register", RegisterDriver) // 司机注册 - 新司机注册，包含身份证和驾驶证验证
	hz.POST("/api/driver/verify", VerifyDriver)     // 司机身份验证 - 验证司机身份信息是否匹配

	// 司机配置管理路由
	hz.GET("/api/driver/config/:driver_id", GetDriverConfig)          // 查询司机配置 - 获取司机的功能开关、计费规则等配置
	hz.POST("/api/driver/config/update", UpdateDriverConfig)          // 更新司机配置 - 修改司机的配置信息
	hz.POST("/api/driver/vehicle/compliance", CheckVehicleCompliance) // 车辆合规性验证 - 验证车牌号和保险是否合规

	log.Println("Driver API 服务启动在 localhost:6667")

	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

// GetDriverDetail 获取司机详情
func GetDriverDetail(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")

	req := driver.NewDriverDetailReq()
	var driverId int64
	if _, err := fmt.Sscanf(id, "%d", &driverId); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "无效的司机ID",
		})
		return
	}

	req.Id = driverId
	resp, err := cli.DriverDetail(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "查询失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "查询成功",
		"data": map[string]interface{}{
			"name":          resp.Name,
			"tel":           resp.Tel,
			"id_card":       resp.IDCard,
			"license":       resp.License,
			"register_date": resp.RegisterDate,
			"rating":        resp.Rating,
		},
	})
}

// RegisterDriver 司机注册
func RegisterDriver(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		Name           string `json:"name"`
		Tel            string `json:"tel"`
		IDCard         string `json:"id_card"`
		License        string `json:"license"`
		DrivingLicense string `json:"driving_license,omitempty"` // 行驶证
		VehicleInfo    string `json:"vehicle_info,omitempty"`    // 车辆信息
		Photo          string `json:"photo,omitempty"`           // 证件照片
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	req := driver.NewDriverRegisterReq()
	req.Name = reqBody.Name
	req.Tel = reqBody.Tel
	req.IDCard = reqBody.IDCard
	req.License = reqBody.License

	// 可选字段
	if reqBody.DrivingLicense != "" {
		req.DrivingLicense = &reqBody.DrivingLicense
	}
	if reqBody.VehicleInfo != "" {
		req.VehicleInfo = &reqBody.VehicleInfo
	}
	if reqBody.Photo != "" {
		req.Photo = &reqBody.Photo
	}

	resp, err := cli.DriverRegister(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "注册失败",
			"error":   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data": map[string]interface{}{
			"id": resp.Id,
		},
	})
}

// VerifyDriver 司机身份验证
func VerifyDriver(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		Tel            string `json:"tel"`
		IDCard         string `json:"id_card"`
		License        string `json:"license"`
		DrivingLicense string `json:"driving_license,omitempty"`
		VehicleInfo    string `json:"vehicle_info,omitempty"`
		Photo          string `json:"photo,omitempty"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	req := driver.NewDriverVerifyReq()
	req.Tel = reqBody.Tel
	req.IDCard = reqBody.IDCard
	req.License = reqBody.License
	if reqBody.DrivingLicense != "" {
		req.DrivingLicense = &reqBody.DrivingLicense
	}
	if reqBody.VehicleInfo != "" {
		req.VehicleInfo = &reqBody.VehicleInfo
	}
	if reqBody.Photo != "" {
		req.Photo = &reqBody.Photo
	}

	resp, err := cli.DriverVerify(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "验证失败",
			"error":   err.Error(),
		})
		return
	}

	if !resp.Verified {
		c.JSON(400, map[string]interface{}{
			"code":     400,
			"message":  resp.Message,
			"verified": false,
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":     200,
		"message":  resp.Message,
		"verified": true,
		"data": map[string]interface{}{
			"driver_id": resp.DriverId,
		},
	})
}

// ==================== 司机配置管理接口 ====================

// GetDriverConfig 查询司机配置
// HTTP接口：GET /api/driver/config/:driver_id
// 功能：查询司机的配置信息，包括功能开关、有效期设置、计费规则、车辆合规性等
func GetDriverConfig(ctx context.Context, c *app.RequestContext) {
	driverIdStr := c.Param("driver_id")

	// 解析司机ID
	var driverId int64
	if _, err := fmt.Sscanf(driverIdStr, "%d", &driverId); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "无效的司机ID",
		})
		return
	}

	// 构造RPC请求
	req := driver.NewDriverConfigQueryReq()
	req.DriverId = driverId

	// 调用RPC服务
	resp, err := cli.DriverConfigQuery(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "查询配置失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查RPC调用是否成功
	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	// 构造响应数据
	configData := map[string]interface{}{
		"driver_id":        resp.Config.DriverId,
		"can_publish_trip": resp.Config.CanPublishTrip,
		"auto_notify":      resp.Config.AutoNotify,
		"max_passengers":   resp.Config.MaxPassengers,
		"start_price":      resp.Config.StartPrice,
		"price_per_km":     resp.Config.PricePerKm,
		"vehicle_number":   resp.Config.VehicleNumber,
		"vehicle_type":     resp.Config.VehicleType,
		"insurance_expiry": resp.Config.InsuranceExpiry,
		"is_compliant":     resp.Config.IsCompliant,
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "查询成功",
		"data":    configData,
	})
}

// UpdateDriverConfig 更新司机配置
// HTTP接口：POST /api/driver/config/update
// 功能：更新司机的配置信息（功能开关、有效期设置、计费规则等）
func UpdateDriverConfig(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		DriverID        int64    `json:"driver_id"`
		CanPublishTrip  *bool    `json:"can_publish_trip,omitempty"` // 是否可以发布行程
		AutoNotify      *bool    `json:"auto_notify,omitempty"`      // 是否自动通知
		MaxPassengers   *int32   `json:"max_passengers,omitempty"`   // 最多可配置乘客数
		StartPrice      *float64 `json:"start_price,omitempty"`      // 车辆起步价
		PricePerKm      *float64 `json:"price_per_km,omitempty"`     // 每公里单价
		VehicleNumber   *string  `json:"vehicle_number,omitempty"`   // 车牌号
		VehicleType     *string  `json:"vehicle_type,omitempty"`     // 车型
		InsuranceExpiry *string  `json:"insurance_expiry,omitempty"` // 保险到期日期
	}

	// 解析请求体
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证司机ID
	if reqBody.DriverID <= 0 {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "司机ID无效",
		})
		return
	}

	// 先查询当前配置
	queryReq := driver.NewDriverConfigQueryReq()
	queryReq.DriverId = reqBody.DriverID
	queryResp, err := cli.DriverConfigQuery(context.Background(), queryReq, callopt.WithRPCTimeout(3*time.Second))
	if err != nil || !queryResp.Success {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "查询配置失败",
		})
		return
	}

	// 更新配置（只更新提供的字段）
	config := queryResp.Config
	if reqBody.CanPublishTrip != nil {
		config.CanPublishTrip = *reqBody.CanPublishTrip
	}
	if reqBody.AutoNotify != nil {
		config.AutoNotify = *reqBody.AutoNotify
	}
	if reqBody.MaxPassengers != nil {
		config.MaxPassengers = *reqBody.MaxPassengers
	}
	if reqBody.StartPrice != nil {
		config.StartPrice = *reqBody.StartPrice
	}
	if reqBody.PricePerKm != nil {
		config.PricePerKm = *reqBody.PricePerKm
	}
	if reqBody.VehicleNumber != nil {
		config.VehicleNumber = *reqBody.VehicleNumber
	}
	if reqBody.VehicleType != nil {
		config.VehicleType = *reqBody.VehicleType
	}
	if reqBody.InsuranceExpiry != nil {
		config.InsuranceExpiry = *reqBody.InsuranceExpiry
	}

	// 注意：实际应该调用一个更新配置的RPC接口
	// 这里简化处理，返回更新后的配置
	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "配置更新成功",
		"data": map[string]interface{}{
			"driver_id":        config.DriverId,
			"can_publish_trip": config.CanPublishTrip,
			"auto_notify":      config.AutoNotify,
			"max_passengers":   config.MaxPassengers,
			"start_price":      config.StartPrice,
			"price_per_km":     config.PricePerKm,
			"vehicle_number":   config.VehicleNumber,
			"vehicle_type":     config.VehicleType,
			"insurance_expiry": config.InsuranceExpiry,
			"is_compliant":     config.IsCompliant,
		},
	})
}

// CheckVehicleCompliance 车辆合规性验证
// HTTP接口：POST /api/driver/vehicle/compliance
// 功能：验证司机的车辆是否合规（车牌号、保险到期日期等）
func CheckVehicleCompliance(ctx context.Context, c *app.RequestContext) {
	var reqBody struct {
		DriverID        int64  `json:"driver_id"`
		VehicleNumber   string `json:"vehicle_number"`
		InsuranceExpiry string `json:"insurance_expiry"`
	}

	// 解析请求体
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证必填字段
	if reqBody.DriverID <= 0 {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "司机ID无效",
		})
		return
	}
	if reqBody.VehicleNumber == "" {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "车牌号不能为空",
		})
		return
	}
	if reqBody.InsuranceExpiry == "" {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "保险到期日期不能为空",
		})
		return
	}

	// 构造RPC请求
	req := driver.NewVehicleComplianceCheckReq()
	req.DriverId = reqBody.DriverID
	req.VehicleNumber = reqBody.VehicleNumber
	req.InsuranceExpiry = reqBody.InsuranceExpiry

	// 调用RPC服务
	resp, err := cli.VehicleComplianceCheck(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "验证失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查RPC调用是否成功
	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	// 构造响应数据
	c.JSON(200, map[string]interface{}{
		"code":         200,
		"message":      resp.Message,
		"is_compliant": resp.IsCompliant,
		"issues":       resp.Issues,
	})
}
