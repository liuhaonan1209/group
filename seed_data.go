package main

import (
	"flag"
	"fmt"
	"group/core"
	"group/global"
	"group/handler/model"
	"log"
	"math/rand"
	"time"
)

var (
	totalRows int
	batchSize int
)

func init() {
	flag.IntVar(&totalRows, "n", 1000000, "每个表插入的数据总量")
	flag.IntVar(&batchSize, "b", 2000, "每批插入的数据量")
	// 修复随机数初始化
	rand.Seed(time.Now().UnixNano())
}

func main() {
	flag.Parse()

	fmt.Println(">>> 脚本启动中...")

	// 1. 初始化配置与数据库连接
	fmt.Println(">>> 正在从 Nacos 读取配置...")
	core.Nacos()

	fmt.Println(">>> 正在连接 MySQL 并执行自动迁移...")
	core.Mysql()

	db := global.DB
	if db == nil {
		log.Fatal("！！！错误：数据库连接对象 global.DB 为空")
	}

	// 设置数据库连接池优化大批量插入
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("！！！获取SQLDB失败: %v", err)
	}
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	fmt.Printf(">>> 开始为每个表插入 %d 条数据，每批 %d 条...\n", totalRows, batchSize)

	seeders := []struct {
		Name      string
		Generator func(index int) interface{}
		Model     interface{}
	}{
		{
			Name:  "StartStation (站点表)",
			Model: &model.StartStation{},
			Generator: func(i int) interface{} {
				return &model.StartStation{
					Name:          fmt.Sprintf("站点_%d", i),
					Latitude:      39.9 + rand.Float64()*10 - 5,  // 北京附近纬度 ±5度
					Longitude:     116.4 + rand.Float64()*10 - 5, // 北京附近经度 ±5度
					TimeFromStart: time.Now().AddDate(0, 0, -rand.Intn(365)),
					Address:       fmt.Sprintf("北京市朝阳区第%d大街%d号", i%100, i%1000),
					IsActive:      rand.Intn(10) > 1, // 90%概率启用
				}
			},
		},
		{
			Name:  "Driver (司机表)",
			Model: &model.Driver{},
			Generator: func(i int) interface{} {
				// 修复评分范围：保证评分在4.0~5.0之间
				rating := 4.0 + rand.Float64()
				if rating > 5.0 {
					rating = 5.0
				}
				// 修复：返回指针类型，匹配GORM要求
				return &model.Driver{
					Name:         fmt.Sprintf("司机_%d", i),
					Tel:          fmt.Sprintf("138%08d", i%100000000),
					IDCard:       fmt.Sprintf("53010119900101%04d", i%10000),
					License:      fmt.Sprintf("C1%06d", i),
					RegisterDate: time.Now().AddDate(0, 0, -rand.Intn(365)),
					Rating:       rating,
				}
			},
		},
		{
			Name:  "Passenger (乘客表)",
			Model: &model.Passenger{},
			Generator: func(i int) interface{} {
				return &model.Passenger{
					Name:         fmt.Sprintf("乘客_%d", i),
					Tel:          fmt.Sprintf("139%08d", i%100000000),
					IDCard:       fmt.Sprintf("53010120000101%04d", i%10000),
					RegisterDate: time.Now().AddDate(0, 0, -rand.Intn(365)),
				}
			},
		},
		{
			Name:  "Route (线路表)",
			Model: &model.Route{},
			Generator: func(i int) interface{} {
				return &model.Route{
					RouteNo:        fmt.Sprintf("R%06d-%d", i, time.Now().UnixNano()%1000),
					Fieet:          fmt.Sprintf("F%03d", i%100),
					Tage:           "高峰线",
					TagShowStatus:  "show",
					StartStationId: uint(rand.Intn(1000) + 1),
					EndStationId:   uint(rand.Intn(1000) + 1),
					PassStations:   fmt.Sprintf("站点_%d,站点_%d", rand.Intn(100), rand.Intn(100)),
					RoutePath:      fmt.Sprintf("路径_%d", i),
					IsActive:       "1",
				}
			},
		},
		{
			Name:  "BusSchedules (班次表)",
			Model: &model.BusSchedules{},
			Generator: func(i int) interface{} {
				// 修复时间字段：使用分钟数而不是时间戳
				departureMinutes := uint(rand.Intn(1440))                    // 0-1440分钟 (24小时)
				arrivalMinutes := departureMinutes + uint(rand.Intn(180)+30) // 30-210分钟后到达
				return &model.BusSchedules{
					RouteId:       uint(rand.Intn(10000) + 1),
					DepartureTime: departureMinutes,
					ArrivalTime:   arrivalMinutes,
					Capacity:      uint(rand.Intn(50) + 10),
					IsActive:      true,
				}
			},
		},
		{
			Name:  "BalanceSheet (收支对账表)",
			Model: &model.BalanceSheet{},
			Generator: func(i int) interface{} {
				// 修复：确保decimal精度正确
				income := float64(rand.Intn(8000) + 1000)         // 1000-9000，避免过大
				refund := float64(rand.Intn(int(income / 20)))    // 退款不超过收入的5%
				settle := float64(int((income-refund)*100)) / 100 // 保留2位小数

				return &model.BalanceSheet{
					SettleDate:  time.Now().AddDate(0, 0, -(i % 365)),
					OrderIncome: float64(int(income*100)) / 100, // 保留2位小数
					OrderRefund: float64(int(refund*100)) / 100, // 保留2位小数
					OrderSettle: settle,
					Status:      (i % 2) + 1, // 交替状态
				}
			},
		},
		{
			Name:  "Transaction (交易明细表)",
			Model: &model.Transaction{},
			Generator: func(i int) interface{} {
				transDate := time.Now().AddDate(0, 0, -(i % 365))
				// 修复：确保订单号唯一性和金额精度
				return &model.Transaction{
					TransDate:     transDate,
					TransTime:     transDate.Add(time.Duration(rand.Intn(24)) * time.Hour),
					OrderNo:       fmt.Sprintf("ORD%012d", i+1), // 确保唯一性
					DriverId:      uint((i % 10000) + 1),
					DriverName:    fmt.Sprintf("司机_%d", (i%1000)+1),
					PassengerId:   uint((i % 10000) + 1),
					PassengerName: fmt.Sprintf("乘客_%d", (i%1000)+1),
					StartStation:  fmt.Sprintf("起点站_%d", (i%100)+1),
					EndStation:    fmt.Sprintf("终点站_%d", ((i+50)%100)+1),
					Amount:        float64(int((float64(rand.Intn(150)+10) * 100))) / 100, // 10-160元，保留2位小数
					PaymentMethod: (i % 3) + 1,                                            // 循环1,2,3
					TransType:     1,
					Remark:        fmt.Sprintf("交易备注_%d", i+1),
				}
			},
		},
		{
			Name:  "RouteStops (路线站点关联表)",
			Model: &model.RouteStops{},
			Generator: func(i int) interface{} {
				return &model.RouteStops{
					RouteId:    rand.Intn(10000) + 1,
					StationId:  rand.Intn(10000) + 1,
					TravelTime: rand.Intn(120),    // 0-120分钟
					StopOrder:  rand.Intn(20) + 1, // 1-20站
					StopType:   rand.Intn(2) + 1,  // 1-上车，2-下车
					IsActive:   rand.Intn(10) > 1, // 90%概率启用
				}
			},
		},
		{
			Name:  "RouteSettle (线路结算表)",
			Model: &model.RouteSettle{},
			Generator: func(i int) interface{} {
				// 修复：确保decimal精度和数据合理性
				income := float64(rand.Intn(30000) + 5000)     // 5000-35000
				refund := float64(rand.Intn(int(income / 50))) // 退款不超过收入的2%
				ticketPrice := float64(rand.Intn(40) + 10)     // 10-50元

				return &model.RouteSettle{
					SettleDate:             time.Now().AddDate(0, 0, -(i % 365)),
					RouteId:                uint((i % 10000) + 1),
					RouteNo:                fmt.Sprintf("R%06d", (i%100000)+1),
					RouteName:              fmt.Sprintf("线路_%d", (i%1000)+1),
					Fleet:                  fmt.Sprintf("车队_%d", (i%10)+1),
					TicketPrice:            float64(int(ticketPrice*100)) / 100,
					DiscountAmount:         float64(rand.Intn(8) + 2), // 2-10元折扣
					CheckedIncome:          float64(int(income*100)) / 100,
					CheckedTickets:         rand.Intn(800) + 200, // 200-1000张
					UncheckedIncome:        float64(int((float64(rand.Intn(3000)) * 100))) / 100,
					UncheckedTickets:       rand.Intn(80) + 20,
					CheckedRefund:          float64(int(refund*100)) / 100,
					CheckedRefundTickets:   rand.Intn(40) + 10,
					UncheckedRefund:        float64(int((float64(rand.Intn(800)) * 100))) / 100,
					UncheckedRefundTickets: rand.Intn(15) + 5,
				}
			},
		},
		{
			Name:  "ScheduleSettle (班次结算表)",
			Model: &model.ScheduleSettle{},
			Generator: func(i int) interface{} {
				income := float64(rand.Intn(10000) + 500)
				refund := float64(rand.Intn(int(income / 20)))
				return &model.ScheduleSettle{
					SettleDate:             time.Now().AddDate(0, 0, -i%365),
					RouteId:                uint(rand.Intn(10000) + 1),
					ScheduleId:             uint(rand.Intn(10000) + 1),
					RouteNo:                fmt.Sprintf("R%06d", rand.Intn(100000)),
					RouteName:              fmt.Sprintf("线路_%d", rand.Intn(1000)),
					DepartureTime:          fmt.Sprintf("%02d:%02d", rand.Intn(24), rand.Intn(60)),
					TicketPrice:            float64(rand.Intn(50) + 10),
					DiscountAmount:         float64(rand.Intn(10)),
					CheckedIncome:          income,
					CheckedTickets:         rand.Intn(100) + 10,
					UncheckedIncome:        float64(rand.Intn(1000)),
					UncheckedTickets:       rand.Intn(20),
					CheckedRefund:          refund,
					CheckedRefundTickets:   rand.Intn(10),
					UncheckedRefund:        float64(rand.Intn(500)),
					UncheckedRefundTickets: rand.Intn(5),
				}
			},
		},
		{
			Name:  "StationSettle (站点结算表)",
			Model: &model.StationSettle{},
			Generator: func(i int) interface{} {
				income := float64(rand.Intn(20000) + 1000)
				refund := float64(rand.Intn(int(income / 20)))
				return &model.StationSettle{
					SettleDate:             time.Now().AddDate(0, 0, -i%365),
					StationId:              uint(rand.Intn(10000) + 1),
					StationName:            fmt.Sprintf("站点_%d", rand.Intn(1000)),
					TicketPrice:            float64(rand.Intn(50) + 10),
					DiscountAmount:         float64(rand.Intn(10)),
					CheckedIncome:          income,
					CheckedTickets:         rand.Intn(500) + 50,
					UncheckedIncome:        float64(rand.Intn(2000)),
					UncheckedTickets:       rand.Intn(50),
					CheckedRefund:          refund,
					CheckedRefundTickets:   rand.Intn(25),
					UncheckedRefund:        float64(rand.Intn(1000)),
					UncheckedRefundTickets: rand.Intn(10),
				}
			},
		},
		{
			Name:  "DriverSettle (司机结算表)",
			Model: &model.DriverSettle{},
			Generator: func(i int) interface{} {
				// 修复：确保数据精度和范围合理
				totalIncome := float64(rand.Intn(20000) + 1000)               // 1000-21000，避免过大数值
				commission := float64(int(totalIncome*0.1*100)) / 100         // 保留2位小数
				netIncome := float64(int((totalIncome-commission)*100)) / 100 // 保留2位小数

				// 修复：确保手机号格式正确且不重复
				phoneNum := fmt.Sprintf("138%08d", (i%99999999)+1) // 确保唯一性

				return &model.DriverSettle{
					SettleDate:  time.Now().AddDate(0, 0, -(i % 365)),
					DriverId:    uint((i % 10000) + 1), // 确保在合理范围内
					DriverName:  fmt.Sprintf("司机_%d", (i%1000)+1),
					Phone:       phoneNum,
					TripCount:   rand.Intn(100) + 10, // 10-110次，合理范围
					TotalIncome: totalIncome,
					Commission:  commission,
					NetIncome:   netIncome,
					Status:      (i % 2) + 1, // 交替1和2，避免随机重复
				}
			},
		},
	}

	// 遍历每个表执行插入
	for _, seeder := range seeders {
		fmt.Printf(">>> 正在处理表: %s\n", seeder.Name)
		startTime := time.Now()
		insertErr := false

		for start := 0; start < totalRows && !insertErr; start += batchSize {
			end := start + batchSize
			if end > totalRows {
				end = totalRows
			}
			batchCount := end - start

			// 组装批量数据：使用指针切片
			batch := make([]interface{}, 0, batchCount)
			for i := start; i < end; i++ {
				batch = append(batch, seeder.Generator(i))
			}

			// 修复：使用CreateInBatches正确的方式
			err := db.CreateInBatches(batch, batchSize).Error
			if err != nil {
				log.Printf("！！！插入表 %s 出错: %v\n", seeder.Name, err)
				insertErr = true
				break
			}

			// 每10批打印进度，减少IO开销
			if batchNo := start / batchSize; batchNo > 0 && batchNo%10 == 0 {
				elapsed := time.Since(startTime)
				progress := float64(end) / float64(totalRows) * 100
				estRemain := elapsed * time.Duration((100-progress)/progress)
				fmt.Printf("    进度: %d/%d (%.2f%%) | 已耗时: %v | 预估剩余: %v\n",
					end, totalRows, progress, elapsed, estRemain)
			}
		}

		if insertErr {
			log.Printf("！！！表 %s 插入中断\n", seeder.Name)
		} else {
			elapsed := time.Since(startTime)
			avgSpeed := float64(totalRows) / elapsed.Seconds()
			fmt.Printf(">>> 表 %s 插入完成! 总耗时: %v，平均每秒插入: %.0f 条\n",
				seeder.Name, elapsed, avgSpeed)
		}
	}

	fmt.Println(">>> 所有数据插入任务执行完毕!")
}
