// Package idcard 提供身份证实名验证功能
// 通过调用腾讯云API网关服务进行身份证号和姓名的二要素验证
package idcard

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	gourl "net/url"
	"strings"
	"time"
)

// calcAuthorization 计算API请求的授权签名
// 使用HMAC-SHA1算法对请求进行签名，确保请求的安全性
// 参数:
//   - source: 请求来源标识
//   - secretId: 云市场分配的密钥ID
//   - secretKey: 云市场分配的密钥Key
// 返回:
//   - auth: 生成的授权字符串
//   - datetime: GMT格式的当前时间
//   - err: 错误信息
func calcAuthorization(source string, secretId string, secretKey string) (auth string, datetime string, err error) {
	// 加载GMT时区
	timeLocation, _ := time.LoadLocation("Etc/GMT")
	// 格式化当前时间为GMT格式
	datetime = time.Now().In(timeLocation).Format("Mon, 02 Jan 2006 15:04:05 GMT")
	// 构造待签名字符串
	signStr := fmt.Sprintf("x-date: %s\nx-source: %s", datetime, source)

	// 使用HMAC-SHA1算法进行签名
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(signStr))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// 构造授权头信息
	auth = fmt.Sprintf("hmac id=\"%s\", algorithm=\"hmac-sha1\", headers=\"x-date x-source\", signature=\"%s\"",
		secretId, sign)

	return auth, datetime, nil
}

// urlencode 将参数map编码为URL查询字符串格式
// 参数:
//   - params: 需要编码的参数键值对
// 返回:
//   - 编码后的URL查询字符串
func urlencode(params map[string]string) string {
	var p = gourl.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	return p.Encode()
}

// Authorization 身份证实名验证主函数
// 调用腾讯云API网关服务，验证身份证号和真实姓名是否匹配
// 参数:
//   - cardNo: 身份证号码
//   - realName: 真实姓名
// 返回:
//   - error: 如果验证失败或姓名身份证不匹配，返回错误信息；验证通过返回nil
func Authorization(cardNo, realName string) error {
	// 云市场分配的密钥Id
	secretId := "OclsjMcLujFHTAPs"
	// 云市场分配的密钥Key
	secretKey := "iD8bMTxS9VJ6PJxpITba0YrHFpXZHMNr"
	// 请求来源标识
	source := "market"

	// 计算请求签名
	auth, datetime, _ := calcAuthorization(source, secretId, secretKey)

	// 设置请求方法为POST
	method := "POST"
	// 构造请求头
	headers := map[string]string{"X-Source": source, "X-Date": datetime, "Authorization": auth}

	// 查询参数（此接口不需要查询参数）
	queryParams := make(map[string]string)

	// 构造请求体参数
	bodyParams := make(map[string]string)
	bodyParams["cardNo"] = cardNo       // 身份证号
	bodyParams["realName"] = realName   // 真实姓名

	// API网关地址
	url := "https://service-18c38npd-1300755093.ap-beijing.apigateway.myqcloud.com/release/idcard/VerifyIdcardv2"
	// 如果有查询参数，拼接到URL后面
	if len(queryParams) > 0 {
		url = fmt.Sprintf("%s?%s", url, urlencode(queryParams))
	}

	// 定义需要body的HTTP方法
	bodyMethods := map[string]bool{"POST": true, "PUT": true, "PATCH": true}
	var body io.Reader = nil
	// 如果是POST/PUT/PATCH方法，设置请求体
	if bodyMethods[method] {
		body = strings.NewReader(urlencode(bodyParams))
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	}

	// 创建HTTP客户端，设置5秒超时
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 创建HTTP请求
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}

	// 设置请求头
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	// 发送HTTP请求
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// 读取响应体
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// 打印响应内容（用于调试）
	fmt.Println(string(bodyBytes))

	// 解析响应JSON
	var t IDCardVerifyResponse
	err = json.Unmarshal(bodyBytes, &t)
	if err != nil {
		return errors.New("解析失败")
	}

	// 检查验证结果
	if t.Result.Isok != true {
		return errors.New("姓名和身份证号不匹配")
	}

	return nil
}

// IDCardVerifyResponse 身份证验证API响应结构体
// 定义了腾讯云身份证验证接口返回的数据结构
type IDCardVerifyResponse struct {
	ErrorCode int    `json:"error_code"` // 错误码，0表示成功
	Reason    string `json:"reason"`     // 返回说明
	Result    struct {
		Realname string `json:"realname"` // 真实姓名
		Idcard   string `json:"idcard"`   // 身份证号
		Isok     bool   `json:"isok"`     // 是否验证通过
		IdCardInfor struct {
			Province string `json:"province"` // 省份
			City     string `json:"city"`     // 城市
			District string `json:"district"` // 区县
			Area     string `json:"area"`     // 地区
			Sex      string `json:"sex"`      // 性别
			Birthday string `json:"birthday"` // 出生日期
		} `json:"IdCardInfor"` // 身份证详细信息
	} `json:"result"` // 验证结果
}
