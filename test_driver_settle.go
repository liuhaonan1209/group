package main

import (
	"fmt"
	"group/core"
	"group/global"
	"group/handler/model"
	"log"
	"math/rand"
	"time"
)

func main() {
	fmt.Println(">>> 测试 DriverSettle 表插入...")

	// 初始化配置与数据库连接
	core.Nacos()
	core.Mysql()

	db := global.DB
	if db == nil {
		log.Fatal("数据库连接失败")
	}

	// 设置数据库连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取SQLDB失败: %v", err)
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)

	rand.Seed(time.Now().UnixNano())

	fmt.Println(">>> 开始插入测试数据...")

	// 测试插入100条数据
	batchSize := 10
	totalRows := 100

	for start := 0; start < totalRows; start += batchSize {
		end := start + batchSize
		if end > totalRows {
			end = totalRows
		}

		// 生成批量数据
		batch := make([]*model.DriverSettle, 0, batchSize)
		for i := start; i < end; i++ {
			totalIncome := float64(rand.Intn(20000) + 1000)
			commission := float64(int(totalIncome * 0.1 * 100)) / 100
			netIncome := float64(int((totalIncome - commission) * 100)) / 100
			
			phoneNum := fmt.Sprintf("138%08d", (i%99999999)+1)
			
			data := &model.DriverSettle{
				SettleDate:  time.Now().AddDate(0, 0, -(i%365)),
				DriverId:    uint((i%10000) + 1),
				DriverName:  fmt.Sprintf("测试司机_%d", (i%1000)+1),
				Phone:       phoneNum,
				TripCount:   rand.Intn(100) + 10,
				TotalIncome: totalIncome,
				Commission:  commission,
				NetIncome:   netIncome,
				Status:      (i%2) + 1,
			}
			batch = append(batch, data)
		}

		// 批量插入
		err := db.CreateInBatches(batch, batchSize).Error
		if err != nil {
			log.Printf("批次 %d-%d 插入失败: %v", start, end-1, err)
			
			// 尝试单条插入以找出问题数据
			fmt.Println("尝试单条插入以定位问题...")
			for j, item := range batch {
				err2 := db.Create(item).Error
				if err2 != nil {
					log.Printf("第 %d 条数据插入失败: %v", start+j, err2)
					fmt.Printf("问题数据: %+v\n", item)
				} else {
					fmt.Printf("第 %d 条数据插入成功\n", start+j)
				}
			}
			break
		} else {
			fmt.Printf("批次 %d-%d 插入成功\n", start, end-1)
		}
	}

	// 检查插入结果
	var count int64
	err = db.Model(&model.DriverSettle{}).Count(&count)
	if err != nil {
		log.Printf("获取数据量失败: %v", err)
	} else {
		fmt.Printf("当前 DriverSettle 表数据量: %d\n", count)
	}

	fmt.Println(">>> 测试完成")
}