// Package verifycode 验证码管理包
// 提供验证码生成、存储、验证功能
// 优先使用Redis存储，如果Redis不可用则使用内存存储
package verifycode

import (
	"context"
	"fmt"
	"group/global"
	"math/rand"
	"sync"
	"time"
)

// VerifyCode 验证码结构
type VerifyCode struct {
	Code      string    // 验证码
	ExpireAt  time.Time // 过期时间
	CreatedAt time.Time // 创建时间
}

// 验证码内存存储（Redis不可用时的备用方案）
var (
	codeStore = make(map[string]*VerifyCode) // key: 手机号, value: 验证码信息
	mu        sync.RWMutex                   // 读写锁
)

const (
	// 验证码有效期（5分钟）
	CodeExpiration = 5 * time.Minute
	// Redis key前缀
	RedisKeyPrefix = "verify_code:"
)

// GenerateCode 生成6位数字验证码
func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000 // 生成100000-999999之间的随机数
	return fmt.Sprintf("%06d", code)
}

// SendCode 发送验证码
// 优先使用Redis存储，如果Redis不可用则使用内存存储
// 实际项目中应该调用短信服务商API
func SendCode(tel string) (string, error) {
	// 生成验证码
	code := GenerateCode()

	// 尝试使用Redis存储
	if global.RedisClient != nil {
		err := saveToRedis(tel, code)
		if err != nil {
			fmt.Printf("⚠ Redis存储失败，使用内存存储: %v\n", err)
			saveToMemory(tel, code)
		} else {
			fmt.Printf("✓ 验证码已存储到Redis\n")
		}
	} else {
		// Redis不可用，使用内存存储
		saveToMemory(tel, code)
		fmt.Printf("⚠ Redis未连接，使用内存存储\n")
	}

	// 模拟发送短信（实际项目中调用短信API）
	fmt.Printf("【模拟短信】发送验证码到 %s: %s (5分钟内有效)\n", tel, code)

	return code, nil
}

// saveToRedis 保存验证码到Redis
func saveToRedis(tel, code string) error {
	ctx := context.Background()
	key := RedisKeyPrefix + tel
	
	// 存储验证码，设置5分钟过期
	err := global.RedisClient.Set(ctx, key, code, CodeExpiration).Err()
	if err != nil {
		return err
	}
	
	return nil
}

// saveToMemory 保存验证码到内存
func saveToMemory(tel, code string) {
	mu.Lock()
	defer mu.Unlock()

	codeStore[tel] = &VerifyCode{
		Code:      code,
		ExpireAt:  time.Now().Add(CodeExpiration),
		CreatedAt: time.Now(),
	}
}

// Verify 验证验证码
// 优先从Redis验证，如果Redis不可用则从内存验证
func Verify(tel, code string) bool {
	// 尝试从Redis验证
	if global.RedisClient != nil {
		valid, err := verifyFromRedis(tel, code)
		if err != nil {
			fmt.Printf("⚠ Redis验证失败，尝试内存验证: %v\n", err)
			return verifyFromMemory(tel, code)
		}
		return valid
	}

	// Redis不可用，从内存验证
	return verifyFromMemory(tel, code)
}

// verifyFromRedis 从Redis验证验证码
func verifyFromRedis(tel, code string) (bool, error) {
	ctx := context.Background()
	key := RedisKeyPrefix + tel

	// 获取存储的验证码
	storedCode, err := global.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// 验证码匹配
	if storedCode == code {
		// 验证成功后删除验证码（一次性使用）
		global.RedisClient.Del(ctx, key)
		fmt.Printf("✓ Redis验证码验证成功: %s\n", tel)
		return true, nil
	}

	return false, nil
}

// verifyFromMemory 从内存验证验证码
func verifyFromMemory(tel, code string) bool {
	mu.RLock()
	defer mu.RUnlock()

	// 查找验证码
	storedCode, exists := codeStore[tel]
	if !exists {
		return false
	}

	// 检查是否过期
	if time.Now().After(storedCode.ExpireAt) {
		// 删除过期验证码
		delete(codeStore, tel)
		return false
	}

	// 验证码匹配
	if storedCode.Code == code {
		// 验证成功后删除验证码（一次性使用）
		delete(codeStore, tel)
		fmt.Printf("✓ 内存验证码验证成功: %s\n", tel)
		return true
	}

	return false
}

// CleanExpiredCodes 清理过期验证码（仅用于内存存储）
// Redis会自动过期，不需要手动清理
func CleanExpiredCodes() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			mu.Lock()
			now := time.Now()
			count := 0
			for tel, code := range codeStore {
				if now.After(code.ExpireAt) {
					delete(codeStore, tel)
					count++
				}
			}
			if count > 0 {
				fmt.Printf("清理过期验证码: %d 条\n", count)
			}
			mu.Unlock()
		}
	}()
}
