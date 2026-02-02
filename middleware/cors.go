package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// CORS 跨域中间件
func CORS() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 设置允许的源
		c.Header("Access-Control-Allow-Origin", "*")
		
		// 设置允许的请求方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		
		// 设置允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")
		
		// 设置允许携带凭证
		c.Header("Access-Control-Allow-Credentials", "true")
		
		// 设置预检请求的缓存时间（秒）
		c.Header("Access-Control-Max-Age", "86400")
		
		// 设置允许暴露的响应头
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		// 处理 OPTIONS 预检请求
		if string(c.Method()) == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// 继续处理请求
		c.Next(ctx)
	}
}
