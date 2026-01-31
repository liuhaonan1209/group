package alipay

//
//import (
//	"fmt"
//
//	"github.com/smartwalle/alipay/v3"
//)
//
//func ALiPay() {
//	var privateKey = "xxx" // 必须，上一步中使用 RSA签名验签工具 生成的私钥
//	var client, err = alipay.New(appId, privateKey, false)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	var p = alipay.TradeWapPay{}
//	p.NotifyURL = "http://xxx"
//	p.ReturnURL = "http://xxx"
//	p.Subject = "标题"
//	p.OutTradeNo = "传递一个唯一单号"
//	p.TotalAmount = "10.00"
//	p.ProductCode = "QUICK_WAP_WAY"
//
//	var url, err = client.TradeWapPay(p)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	// 这个 payURL 即是用于打开支付宝支付页面的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
//	var payURL = url.String()
//	fmt.Println(payURL)
//}
