package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	upload "group/kitex_gen/car/upload"

	"github.com/google/uuid"
)

// UploadServiceImpl 上传服务实现结构体
type UploadServiceImpl struct{}

// ImageUpload 图片上传接口
func (s *UploadServiceImpl) ImageUpload(ctx context.Context, req *upload.ImageUploadReq) (resp *upload.ImageUploadResp, err error) {
	log.Printf("图片上传请求: 文件名=%s, 文件大小=%d, 文件类型=%s", req.FileName, req.FileSize, req.FileType)

	// 检查文件类型
	if !isValidImageType(req.FileType) {
		return &upload.ImageUploadResp{
			FileURL:  "",
			FileName: "",
			FileSize: 0,
			Message:  "不支持的文件类型，仅支持 jpg, jpeg, png, gif, bmp, webp",
			Success:  false,
		}, nil
	}

	// 检查文件大小（限制10MB）
	maxSize := int64(10 * 1024 * 1024)
	if req.FileSize > maxSize {
		return &upload.ImageUploadResp{
			FileURL:  "",
			FileName: "",
			FileSize: 0,
			Message:  "文件大小超过限制（最大10MB）",
			Success:  false,
		}, nil
	}

	// 确定上传类型
	uploadType := "oss" // 默认使用OSS
	if req.UploadType != nil && *req.UploadType != "" {
		uploadType = *req.UploadType
	}

	// 生成唯一文件名
	ext := req.FileType
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	dateDir := time.Now().Format("2006/01/02")
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	objectKey := fmt.Sprintf("images/%s/%s", dateDir, newFileName)

	var fileURL string

	if uploadType == "oss" {
		// 上传到阿里云OSS
		// 这里使用环境变量配置
		bucketName := "2305ahby"
		region := "cn-shanghai"

		// 模拟上传到OSS（实际应该调用OSS SDK）
		fileURL = fmt.Sprintf("https://%s.oss-%s.aliyuncs.com/%s", bucketName, region, objectKey)
		log.Printf("文件上传到OSS成功: %s", fileURL)
	} else if uploadType == "local" {
		// 上传到本地存储
		fileURL = fmt.Sprintf("http://localhost:8080/uploads/%s/%s", dateDir, newFileName)
		log.Printf("文件上传到本地成功: %s", fileURL)
	} else {
		return &upload.ImageUploadResp{
			FileURL:  "",
			FileName: "",
			FileSize: 0,
			Message:  "不支持的上传类型，仅支持 oss 或 local",
			Success:  false,
		}, nil
	}

	log.Printf("图片上传成功: 文件名=%s, URL=%s", newFileName, fileURL)
	return &upload.ImageUploadResp{
		FileURL:  fileURL,
		FileName: newFileName,
		FileSize: req.FileSize,
		Message:  "上传成功",
		Success:  true,
	}, nil
}

// ImageDelete 图片删除接口
func (s *UploadServiceImpl) ImageDelete(ctx context.Context, req *upload.ImageDeleteReq) (resp *upload.ImageDeleteResp, err error) {
	log.Printf("图片删除请求: 文件URL=%s", req.FileURL)

	// 确定删除类型
	deleteType := "oss" // 默认使用OSS
	if req.DeleteType != nil && *req.DeleteType != "" {
		deleteType = *req.DeleteType
	}

	if deleteType == "oss" {
		// 从URL提取objectKey
		// 例如: https://bucket.oss-cn-shanghai.aliyuncs.com/images/2024/01/01/xxx.jpg
		// 提取: images/2024/01/01/xxx.jpg
		parts := strings.Split(req.FileURL, ".aliyuncs.com/")
		if len(parts) != 2 {
			return &upload.ImageDeleteResp{
				Message: "无效的OSS文件URL",
				Success: false,
			}, nil
		}
		objectKey := parts[1]
		log.Printf("从OSS删除文件: %s", objectKey)
		// 这里应该调用OSS SDK删除文件
		// 暂时模拟删除成功
	} else if deleteType == "local" {
		log.Printf("从本地删除文件: %s", req.FileURL)
		// 这里应该删除本地文件
		// 暂时模拟删除成功
	} else {
		return &upload.ImageDeleteResp{
			Message: "不支持的删除类型，仅支持 oss 或 local",
			Success: false,
		}, nil
	}

	log.Printf("图片删除成功: %s", req.FileURL)
	return &upload.ImageDeleteResp{
		Message: "删除成功",
		Success: true,
	}, nil
}

// BatchImageUpload 批量图片上传接口
func (s *UploadServiceImpl) BatchImageUpload(ctx context.Context, req *upload.BatchImageUploadReq) (resp *upload.BatchImageUploadResp, err error) {
	log.Printf("批量图片上传请求: 图片数量=%d", len(req.Images))

	results := make([]*upload.BatchUploadResult_, 0, len(req.Images))
	successCount := int32(0)
	failCount := int32(0)

	// 逐个上传图片
	for _, img := range req.Images {
		uploadResp, err := s.ImageUpload(ctx, img)
		if err != nil || !uploadResp.Success {
			failCount++
			message := "上传失败"
			if uploadResp != nil {
				message = uploadResp.Message
			}
			results = append(results, &upload.BatchUploadResult_{
				FileName: img.FileName,
				FileURL:  "",
				Success:  false,
				Message:  message,
			})
		} else {
			successCount++
			results = append(results, &upload.BatchUploadResult_{
				FileName: uploadResp.FileName,
				FileURL:  uploadResp.FileURL,
				Success:  true,
				Message:  "上传成功",
			})
		}
	}

	allSuccess := failCount == 0
	message := fmt.Sprintf("批量上传完成，成功%d个，失败%d个", successCount, failCount)

	log.Printf("批量图片上传完成: 成功=%d, 失败=%d", successCount, failCount)
	return &upload.BatchImageUploadResp{
		Results:      results,
		SuccessCount: successCount,
		FailCount:    failCount,
		Message:      message,
		Success:      allSuccess,
	}, nil
}

// ImageInfo 图片信息查询接口
func (s *UploadServiceImpl) ImageInfo(ctx context.Context, req *upload.ImageInfoReq) (resp *upload.ImageInfoResp, err error) {
	log.Printf("图片信息查询请求: 文件URL=%s", req.FileURL)

	// 从URL提取文件名和类型
	parts := strings.Split(req.FileURL, "/")
	if len(parts) == 0 {
		return &upload.ImageInfoResp{
			FileName:   "",
			FileSize:   0,
			FileType:   "",
			UploadTime: "",
			Message:    "无效的文件URL",
			Success:    false,
		}, nil
	}

	fileName := parts[len(parts)-1]
	fileType := strings.ToLower(strings.TrimPrefix(strings.Split(fileName, ".")[len(strings.Split(fileName, "."))-1], "."))

	// 这里应该从数据库或OSS查询文件信息
	// 暂时返回模拟数据
	log.Printf("图片信息查询成功: 文件名=%s", fileName)
	return &upload.ImageInfoResp{
		FileName:   fileName,
		FileSize:   0, // 实际应该查询真实大小
		FileType:   fileType,
		UploadTime: time.Now().Format("2006-01-02 15:04:05"),
		Message:    "查询成功",
		Success:    true,
	}, nil
}

// isValidImageType 检查是否为有效的图片类型
func isValidImageType(fileType string) bool {
	validTypes := map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"gif":  true,
		"bmp":  true,
		"webp": true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
	}
	return validTypes[strings.ToLower(fileType)]
}
