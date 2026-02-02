package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

// UploadConfig 上传配置
type UploadConfig struct {
	// 本地存储配置
	LocalStoragePath string // 本地存储路径
	LocalBaseURL     string // 本地访问基础URL

	// 阿里云OSS配置
	OSSEndpoint        string // OSS端点
	OSSAccessKeyID     string // OSS访问密钥ID
	OSSAccessKeySecret string // OSS访问密钥
	OSSBucketName      string // OSS存储桶名称
	OSSRegion          string // OSS区域
}

// UploadResult 上传结果
type UploadResult struct {
	FileName string // 文件名
	FileSize int64  // 文件大小
	FileURL  string // 文件访问URL
	FileType string // 文件类型
}

// 允许的图片格式
var allowedImageTypes = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".webp": true,
}

// IsImageFile 检查是否为图片文件
func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return allowedImageTypes[ext]
}

// UploadToLocal 上传文件到本地存储
func UploadToLocal(file *multipart.FileHeader, config *UploadConfig) (*UploadResult, error) {
	// 检查文件类型
	if !IsImageFile(file.Filename) {
		return nil, fmt.Errorf("不支持的文件类型，仅支持图片格式")
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 按日期创建子目录
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(config.LocalStoragePath, dateDir)

	// 创建目录
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	// 创建目标文件
	destPath := filepath.Join(fullDir, newFileName)
	dst, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	fileSize, err := io.Copy(dst, src)
	if err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 生成访问URL
	fileURL := fmt.Sprintf("%s/%s/%s", config.LocalBaseURL, dateDir, newFileName)

	return &UploadResult{
		FileName: newFileName,
		FileSize: fileSize,
		FileURL:  fileURL,
		FileType: ext,
	}, nil
}

// UploadToOSS 上传文件到阿里云OSS
func UploadToOSS(file *multipart.FileHeader, config *UploadConfig) (*UploadResult, error) {
	// 检查文件类型
	if !IsImageFile(file.Filename) {
		return nil, fmt.Errorf("不支持的文件类型，仅支持图片格式")
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	// 创建OSS客户端
	client, err := oss.New(config.OSSEndpoint, config.OSSAccessKeyID, config.OSSAccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %v", err)
	}

	// 获取存储桶
	bucket, err := client.Bucket(config.OSSBucketName)
	if err != nil {
		return nil, fmt.Errorf("获取存储桶失败: %v", err)
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	dateDir := time.Now().Format("2006/01/02")
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	objectKey := fmt.Sprintf("images/%s/%s", dateDir, newFileName)

	// 上传文件
	err = bucket.PutObject(objectKey, src)
	if err != nil {
		return nil, fmt.Errorf("上传到OSS失败: %v", err)
	}

	// 生成访问URL
	fileURL := fmt.Sprintf("https://%s.%s/%s", config.OSSBucketName, config.OSSEndpoint, objectKey)

	return &UploadResult{
		FileName: newFileName,
		FileSize: file.Size,
		FileURL:  fileURL,
		FileType: ext,
	}, nil
}

// UploadToOSSWithEnv 使用环境变量配置上传到阿里云OSS
func UploadToOSSWithEnv(file *multipart.FileHeader, bucketName, region string) (*UploadResult, error) {
	// 检查文件类型
	if !IsImageFile(file.Filename) {
		return nil, fmt.Errorf("不支持的文件类型，仅支持图片格式")
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	// 从环境变量获取凭证
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		return nil, fmt.Errorf("获取凭证失败: %v", err)
	}

	// 创建客户端选项
	clientOptions := []oss.ClientOption{oss.SetCredentialsProvider(&provider)}
	clientOptions = append(clientOptions, oss.Region(region))
	clientOptions = append(clientOptions, oss.AuthVersion(oss.AuthV4))

	// 创建OSS客户端
	endpoint := fmt.Sprintf("https://oss-%s.aliyuncs.com", region)
	client, err := oss.New(endpoint, "", "", clientOptions...)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %v", err)
	}

	// 获取存储桶
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("获取存储桶失败: %v", err)
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	dateDir := time.Now().Format("2006/01/02")
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	objectKey := fmt.Sprintf("images/%s/%s", dateDir, newFileName)

	// 上传文件
	err = bucket.PutObject(objectKey, src)
	if err != nil {
		return nil, fmt.Errorf("上传到OSS失败: %v", err)
	}

	// 生成访问URL
	fileURL := fmt.Sprintf("https://%s.oss-%s.aliyuncs.com/%s", bucketName, region, objectKey)

	return &UploadResult{
		FileName: newFileName,
		FileSize: file.Size,
		FileURL:  fileURL,
		FileType: ext,
	}, nil
}

// DeleteFromOSS 从阿里云OSS删除文件
func DeleteFromOSS(objectKey string, config *UploadConfig) error {
	// 创建OSS客户端
	client, err := oss.New(config.OSSEndpoint, config.OSSAccessKeyID, config.OSSAccessKeySecret)
	if err != nil {
		return fmt.Errorf("创建OSS客户端失败: %v", err)
	}

	// 获取存储桶
	bucket, err := client.Bucket(config.OSSBucketName)
	if err != nil {
		return fmt.Errorf("获取存储桶失败: %v", err)
	}

	// 删除文件
	err = bucket.DeleteObject(objectKey)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	return nil
}

// DeleteFromLocal 从本地存储删除文件
func DeleteFromLocal(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}
	return nil
}
