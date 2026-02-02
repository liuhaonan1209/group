// Package main 乘客API服务
// 提供HTTP RESTful API接口，作为乘客RPC服务的网关层
// 负责接收HTTP请求，转换为RPC调用，并返回HTTP响应
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"group/kitex_gen/car/passenger"
	"group/kitex_gen/car/passenger/passengerservice"
	"group/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

var (
	// cli 乘客服务RPC客户端
	// 用于调用乘客RPC服务的接口
	cli passengerservice.Client
)

// main 主函数
// 初始化RPC客户端和HTTP服务器，注册路由并启动服务
func main() {
	// 创建乘客服务RPC客户端
	// 连接到127.0.0.1:8988端口的RPC服务
	c, err := passengerservice.NewClient("car.passenger", client.WithHostPorts("127.0.0.1:8988"))
	if err != nil {
		log.Fatal(err)
	}
	cli = c

	// 创建Hertz HTTP服务器
	// 监听localhost:8891端口
	hz := server.New(server.WithHostPorts("127.0.0.1:8891"))

	// 使用 CORS 中间件
	hz.Use(middleware.CORS())

	// ==================== 注册HTTP路由 ====================
	hz.GET("/api/passenger/detail", GetPassengerDetail)   // 查询乘客详情 - 根据乘客ID获取乘客基本信息
	hz.POST("/api/passenger/register", RegisterPassenger) // 乘客注册 - 新乘客注册，包含身份证实名验证
	hz.POST("/api/passenger/send-code", SendVerifyCode)   // 发送验证码 - 向已注册手机号发送验证码
	hz.POST("/api/passenger/verify", VerifyPassenger)     // 乘客身份验证 - 验证手机号和验证码是否匹配

	log.Println("Passenger API 服务启动在 localhost:8891")
	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

// GetPassengerDetail 获取乘客详情接口
// HTTP GET /api/passenger/detail?id=1
// 根据乘客ID查询乘客的详细信息
// 参数:
//   - id: 查询参数，乘客ID
//
// 返回:
//   - 200: 查询成功，返回乘客详情
//   - 400: 无效的乘客ID
//   - 500: 查询失败
func GetPassengerDetail(ctx context.Context, c *app.RequestContext) {
	// 从查询参数中获取乘客ID
	id := c.Query("id")

	// 创建RPC请求对象
	req := passenger.NewPassengerDetailReq()

	// 将字符串ID转换为int64类型
	var passengerId int64
	if _, err := fmt.Sscanf(id, "%d", &passengerId); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "无效的乘客ID",
		})
		return
	}

	req.Id = passengerId

	// 调用RPC服务，设置3秒超时
	resp, err := cli.PassengerDetail(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "查询失败",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "查询成功",
		"data": map[string]interface{}{
			"name":          resp.Name,
			"tel":           resp.Tel,
			"id_card":       resp.IDCard,
			"register_date": resp.RegisterDate,
		},
	})
}

// RegisterPassenger 乘客注册接口
// HTTP POST /api/passenger/register
// 处理乘客注册请求，包括身份证实名验证
// 请求体:
//   - name: 姓名
//   - tel: 手机号
//   - id_card: 身份证号
//   - photo: 照片（可选）
//
// 返回:
//   - 200: 注册成功
//   - 400: 参数错误或注册失败
//   - 500: 服务器错误
func RegisterPassenger(ctx context.Context, c *app.RequestContext) {
	// 定义请求体结构
	var reqBody struct {
		Name   string `json:"name"`
		Tel    string `json:"tel"`
		IDCard string `json:"id_card"`
		Photo  string `json:"photo,omitempty"`
	}

	// 解析JSON请求体
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 构造RPC请求
	req := passenger.NewPassengerRegisterReq()
	req.Name = reqBody.Name
	req.Tel = reqBody.Tel
	req.IDCard = reqBody.IDCard
	if reqBody.Photo != "" {
		req.Photo = &reqBody.Photo
	}

	// 调用RPC服务进行注册（包含身份证实名验证）
	resp, err := cli.PassengerRegister(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "注册失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查注册结果
	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	// 返回成功响应
	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data": map[string]interface{}{
			"id": resp.Id,
		},
	})
}

// VerifyPassenger 乘客身份验证接口
// HTTP POST /api/passenger/verify
// 验证手机号和验证码是否匹配
// 请求体:
//   - tel: 手机号
//   - verify_code: 验证码
//
// 返回:
//   - 200: 验证通过
//   - 400: 参数错误或验证失败
//   - 500: 服务器错误
func VerifyPassenger(ctx context.Context, c *app.RequestContext) {
	// 定义请求体结构
	var reqBody struct {
		Tel        string `json:"tel"`
		VerifyCode string `json:"verify_code"`
	}

	// 解析JSON请求体
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 构造RPC请求
	req := passenger.NewPassengerVerifyReq()
	req.Tel = reqBody.Tel
	req.VerifyCode = reqBody.VerifyCode

	// 调用RPC服务进行身份验证
	resp, err := cli.PassengerVerify(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "验证失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查验证结果
	if !resp.Verified {
		c.JSON(400, map[string]interface{}{
			"code":     400,
			"message":  resp.Message,
			"verified": false,
		})
		return
	}

	// 返回成功响应
	c.JSON(200, map[string]interface{}{
		"code":     200,
		"message":  resp.Message,
		"verified": true,
		"data": map[string]interface{}{
			"passenger_id": resp.PassengerId,
		},
	})
}

// SendVerifyCode 发送验证码接口
// HTTP POST /api/passenger/send-code
// 向指定手机号发送验证码
// 请求体:
//   - tel: 手机号
//
// 返回:
//   - 200: 发送成功
//   - 400: 参数错误或发送失败
//   - 500: 服务器错误
func SendVerifyCode(ctx context.Context, c *app.RequestContext) {
	// 定义请求体结构
	var reqBody struct {
		Tel string `json:"tel"`
	}

	// 解析JSON请求体
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 构造RPC请求
	req := passenger.NewSendVerifyCodeReq()
	req.Tel = reqBody.Tel

	// 调用RPC服务发送验证码
	resp, err := cli.SendVerifyCode(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "发送验证码失败",
			"error":   err.Error(),
		})
		return
	}

	// 检查发送结果
	if !resp.Success {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": resp.Message,
		})
		return
	}

	// 返回成功响应
	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
	})
}
