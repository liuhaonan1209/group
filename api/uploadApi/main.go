package main

import (
	"context"
	"io"
	"log"
	"time"

	"group/kitex_gen/car/upload"
	"group/kitex_gen/car/upload/uploadservice"
	"group/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

var (
	cli uploadservice.Client
)

func main() {
	// 创建 RPC 客户端
	c, err := uploadservice.NewClient("car.upload", client.WithHostPorts("127.0.0.1:8912"))
	if err != nil {
		log.Fatal(err)
	}
	cli = c

	// 创建 HTTP 服务器
	hz := server.New(server.WithHostPorts("127.0.0.1:8996"))

	// 使用 CORS 中间件
	hz.Use(middleware.CORS())

	// ==================== 注册路由 ====================
	// 图片上传路由
	hz.POST("/api/upload/image", UploadImage)        // 单张图片上传
	hz.POST("/api/upload/images", BatchUploadImages) // 批量图片上传
	hz.DELETE("/api/upload/image", DeleteImage)      // 删除图片
	hz.GET("/api/upload/image/info", GetImageInfo)   // 获取图片信息

	log.Println("Upload API 服务启动在 localhost:8996")
	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

// UploadImage 单张图片上传
func UploadImage(ctx context.Context, c *app.RequestContext) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "未找到上传文件",
			"error":   err.Error(),
		})
		return
	}

	// 获取上传类型（可选，默认oss）
	uploadType := c.PostForm("upload_type")
	if uploadType == "" {
		uploadType = "oss"
	}

	// 打开文件读取内容
	src, err := file.Open()
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "打开文件失败",
			"error":   err.Error(),
		})
		return
	}
	defer src.Close()

	// 读取文件内容
	fileData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "读取文件失败",
			"error":   err.Error(),
		})
		return
	}

	// 获取文件扩展名
	fileType := ""
	if len(file.Filename) > 0 {
		for i := len(file.Filename) - 1; i >= 0; i-- {
			if file.Filename[i] == '.' {
				fileType = file.Filename[i:]
				break
			}
		}
	}

	// 构造RPC请求
	req := upload.NewImageUploadReq()
	req.FileName = file.Filename
	req.FileData = fileData
	req.FileSize = file.Size
	req.FileType = fileType
	req.UploadType = &uploadType

	// 调用RPC服务
	resp, err := cli.ImageUpload(context.Background(), req, callopt.WithRPCTimeout(10*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "上传失败",
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
			"file_url":  resp.FileURL,
			"file_name": resp.FileName,
			"file_size": resp.FileSize,
		},
	})
}

// BatchUploadImages 批量图片上传
func BatchUploadImages(ctx context.Context, c *app.RequestContext) {
	// 获取上传类型（可选，默认oss）
	uploadType := c.PostForm("upload_type")
	if uploadType == "" {
		uploadType = "oss"
	}

	// 获取多个文件
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "获取上传文件失败",
			"error":   err.Error(),
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "未找到上传文件",
		})
		return
	}

	// 限制批量上传数量
	if len(files) > 10 {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "批量上传最多支持10张图片",
		})
		return
	}

	// 构造批量上传请求
	images := make([]*upload.ImageUploadReq, 0, len(files))

	for _, file := range files {
		// 打开文件读取内容
		src, err := file.Open()
		if err != nil {
			continue
		}

		// 读取文件内容
		fileData, err := io.ReadAll(src)
		src.Close()
		if err != nil {
			continue
		}

		// 获取文件扩展名
		fileType := ""
		if len(file.Filename) > 0 {
			for i := len(file.Filename) - 1; i >= 0; i-- {
				if file.Filename[i] == '.' {
					fileType = file.Filename[i:]
					break
				}
			}
		}

		imgReq := upload.NewImageUploadReq()
		imgReq.FileName = file.Filename
		imgReq.FileData = fileData
		imgReq.FileSize = file.Size
		imgReq.FileType = fileType
		imgReq.UploadType = &uploadType

		images = append(images, imgReq)
	}

	if len(images) == 0 {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "没有有效的上传文件",
		})
		return
	}

	// 构造RPC请求
	req := upload.NewBatchImageUploadReq()
	req.Images = images

	// 调用RPC服务
	resp, err := cli.BatchImageUpload(context.Background(), req, callopt.WithRPCTimeout(30*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "批量上传失败",
			"error":   err.Error(),
		})
		return
	}

	// 转换结果
	results := make([]map[string]interface{}, 0, len(resp.Results))
	for _, result := range resp.Results {
		resultMap := map[string]interface{}{
			"file_name": result.FileName,
			"file_url":  result.FileURL,
			"success":   result.Success,
			"message":   result.Message,
		}
		results = append(results, resultMap)
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": resp.Message,
		"data": map[string]interface{}{
			"results":       results,
			"success_count": resp.SuccessCount,
			"fail_count":    resp.FailCount,
		},
	})
}

// DeleteImage 删除图片
func DeleteImage(ctx context.Context, c *app.RequestContext) {
	fileURL := c.Query("file_url")
	if fileURL == "" {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "文件URL为必填参数",
		})
		return
	}

	deleteType := c.Query("delete_type")
	if deleteType == "" {
		deleteType = "oss"
	}

	// 构造RPC请求
	req := upload.NewImageDeleteReq()
	req.FileURL = fileURL
	req.DeleteType = &deleteType

	// 调用RPC服务
	resp, err := cli.ImageDelete(context.Background(), req, callopt.WithRPCTimeout(5*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "删除失败",
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
	})
}

// GetImageInfo 获取图片信息
func GetImageInfo(ctx context.Context, c *app.RequestContext) {
	fileURL := c.Query("file_url")
	if fileURL == "" {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "文件URL为必填参数",
		})
		return
	}

	// 构造RPC请求
	req := upload.NewImageInfoReq()
	req.FileURL = fileURL

	// 调用RPC服务
	resp, err := cli.ImageInfo(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "查询失败",
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
			"file_name":   resp.FileName,
			"file_size":   resp.FileSize,
			"file_type":   resp.FileType,
			"upload_time": resp.UploadTime,
		},
	})
}
