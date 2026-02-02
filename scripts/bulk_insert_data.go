package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 数据库配置
const (
	DSN = "root:password@tcp(localhost:3306)/bus_management?charset=utf8mb4&parseTime=True&loc=Local"
	BATCH_SIZE = 1000 // 批量插入大小
	TOTAL_RECORDS = 1000000 // 总记录数
)

// 数据模型定义
type Route struct {
	ID             uint      `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RouteNo        string    `gorm:"type:varchar(20)"`
	Fieet          string    `gorm:"type:varchar(20)"`
	Tage           string    `gorm:"type:varchar(50)"`
	TagShowStatus  string    `gorm:"type:varchar(50)"`
	StartStationId uint      `gorm:"type:int"`
	EndStationId   uint      `gorm:"type:int"`
	PassStations   string    `gorm:"type:varchar(50)"`
	RoutePath      string    `gorm:"type:varchar(50)"`
	IsActive       string    `gorm:"type:int"`
}

type StartStation struct {
	ID            uint      `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string    `gorm:"type:varchar(255)"`
	Latitude      float64   `gorm:"type:float"`
	Longitude     float64   `gorm:"type:float"`
	TimeFromStart time.Time `gorm:"type:datetime"`
	Address       string    `gorm:"type:varchar(255)"`
	IsActive      bool      `gorm:"type:tinyint(1)"`
}

type Driver struct {
	ID           uint      `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string    `gorm:"type:varchar(100)"`
	Tel          string    `gorm:"type:varchar(20)"`
	IDCard       string    `gorm:"type:varchar(18)"`
	License      string    `gorm:"type:varchar(50)"`
	RegisterDate time.Time `gorm:"type:datetime"`
	Rating       float64   `gorm:"type:decimal(3,2)"`
}

type Passenger struct {
	ID           uint      `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string    `gorm:"type:varchar(100)"`
	Tel          string    `gorm:"type:varchar(20)"`
	IDCard       string    `gorm:"type:varchar(18)"`
	RegisterDate time.Time `gorm:"type:datetime"`
}

type BusSchedules struct {
	ID            uint      `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	RouteId       uint      `gorm:"type:int"`
	DepartureTime uint      `gorm:"type:int"`
	ArrivalTime   uint      `gorm:"type:int"`
	Capacity      uint      `gorm:"type:int"`
	IsActive      bool      `gorm:"type:bool"`
}

type RouteStops struct {
	ID         uint      `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	RouteId    int       `gorm:"type:int"`
	StationId  int       `gorm:"type:int"`
	TravelTime int       `gorm:"type:int"`
	StopOrder  int       `gorm:"type:int"`
	StopType   int       `gorm:"type:tinyint"`
	IsActive   bool      `gorm:"type:tinyint(1)"`
}

type Transaction struct {
	ID            uint      `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	TransDate     time.Time `gorm:"type:date"`
	TransTime     time.Time `gorm:"type:datetime"`
	OrderNo       string    `gorm:"type:varchar(50)"`
	DriverId      uint      `gorm:"type:int"`
	DriverName    string    `gorm:"type:varchar(50)"`
	PassengerId   uint      `gorm:"type:int"`
	PassengerName string    `gorm:"type:varchar(50)"`
	StartStation  string    `gorm:"type:varchar(100)"`
	EndStation    string    `gorm:"type:varchar(100)"`
	Amount        float64   `gorm:"type:decimal(10,2)"`
	PaymentMethod int       `gorm:"type:tinyint"`
	TransType     int       `gorm:"type:tinyint"`
	Remark        string    `gorm:"type:varchar(255)"`
}

type BalanceSheet struct {
	ID          uint      `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SettleDate  time.Time `gorm:"type:date"`
	OrderIncome float64   `gorm:"type:decimal(12,2)"`
	OrderRefund float64   `gorm:"type:decimal(12,2)"`
	OrderSettle float64   `gorm:"type:decimal(12,2)"`
	Status      int       `gorm:"type:tinyint"`
}

type DriverSettle struct {
	ID          uint      `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SettleDate  time.Time `gorm:"type:date"`
	DriverId    uint      `gorm:"type:int"`
	DriverName  string    `gorm:"type:varchar(50)"`
	Phone       string    `gorm:"type:varchar(20)"`
	TripCount   int       `gorm:"type:int"`
	TotalIncome float64   `gorm:"type:decimal(12,2)"`
	Commission  float64   `gorm:"type:decimal(12,2)"`
	NetIncome   float64   `gorm:"type:decimal(12,2)"`
	Status      int       `gorm:"type:tinyint"`
}

// 随机数据生成器
var (
	cities = []string{"北京", "上海", "广州", "深圳", "杭州", "南京", "武汉", "成都", "重庆", "西安"}
	fleets = []string{"车队A", "车队B", "车队C", "车队D", "车队E"}
	tags   = []string{"快线", "普通", "夜班", "高峰", "节假日"}
	names  = []string{"张三", "李四", "王五", "赵六", "钱七", "孙八", "周九", "吴十", "郑十一", "王十二"}
	surnames = []string{"张", "李", "王", "赵", "钱", "孙", "周", "吴", "郑", "冯", "陈", "褚", "卫", "蒋"}
	givenNames = []string{"伟", "芳", "娜", "敏", "静", "丽", "强", "磊", "军", "洋", "勇", "艳", "杰", "娟", "涛", "明", "超", "秀英"}
)

func main() {
	// 连接数据库
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 关闭日志以提高性能
	})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 设置连接池
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移表结构
	fmt.Println("开始创建表结构...")
	err = db.AutoMigrate(
		&Route{},
		&StartStation{},
		&Driver{},
		&Passenger{},
		&BusSchedules{},
		&RouteStops{},
		&Transaction{},
		&BalanceSheet{},
		&DriverSettle{},
	)
	if err != nil {
		log.Fatal("创建表结构失败:", err)
	}

	fmt.Println("表结构创建完成，开始插入数据...")

	// 插入各表数据
	insertStations(db)
	insertDrivers(db)
	insertPassengers(db)
	insertRoutes(db)
	insertBusSchedules(db)
	insertRouteStops(db)
	insertTransactions(db)
	insertBalanceSheets(db)
	insertDriverSettles(db)

	fmt.Println("所有数据插入完成！")
}

// 插入站点数据
func insertStations(db *gorm.DB) {
	fmt.Println("开始插入站点数据...")
	start := time.Now()

	for i := 0; i < TOTAL_RECORDS/BATCH_SIZE; i++ {
		var stations []StartStation
		for j := 0; j < BATCH_SIZE; j++ {
			id := i*BATCH_SIZE + j + 1
			station := StartStation{
				Name:          fmt.Sprintf("%s站%d", cities[rand.Intn(len(cities))], id),
				Latitude:      39.9 + rand.Float64()*10, // 北京附近纬度
				Longitude:     116.4 + rand.Float64()*10, // 北京附近经度
				TimeFromStart: time.Now().Add(-time.Duration(rand.Intn(365*24)) * time.Hour),
				Address:       fmt.Sprintf("%s市%s区%s路%d号", cities[rand.Intn(len(cities))], cities[rand.Intn(len(cities))], cities[rand.Intn(len(cities))], rand.Intn(999)+1),
				IsActive:      rand.Intn(10) > 1, // 90%概率启用
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			stations = append(stations, station)
		}

		if err := db.CreateInBatches(stations, BATCH_SIZE).Error; err != nil {
			log.Printf("插入站点数据失败 (批次 %d): %v", i+1, err)
		}

		if (i+1)%100 == 0 {
			fmt.Printf("已插入站点数据: %d/%d 批次\n", i+1, TOTAL_RECORDS/BATCH_SIZE)
		}
	}

	fmt.Printf("站点数据插入完成，耗时: %v\n", time.Since(start))
}

// 插入司机数据
func insertDrivers(db *gorm.DB) {
	fmt.Println("开始插入司机数据...")
	start := time.Now()

	for i := 0; i < TOTAL_RECORDS/BATCH_SIZE; i++ {
		var drivers []Driver
		for j := 0; j < BATCH_SIZE; j++ {
			id := i*BATCH_SIZE + j + 1
			driver := Driver{
				Name:         generateRandomName(),
				Tel:          fmt.Sprintf("1%d%08d", rand.Intn(9)+3, rand.Intn(99999999)+1),
				IDCard:       generateIDCard(),
				License:      fmt.Sprintf("C1%d%08d", rand.Intn(9)+1, rand.Intn(99999999)+1),
				RegisterDate: time.Now().Add(-time.Duration(rand.Intn(365*3)) * time.Hour * 24),
				Rating:       4.0 + rand.Float64(),
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			drivers = append(drivers, driver)
		}

		if err := db.CreateInBatches(drivers, BATCH_SIZE).Error; err != nil {
			log.Printf("插入司机数据失败 (批次 %d): %v", i+1, err)
		}

		if (i+1)%100 == 0 {
			fmt.Printf("已插入司机数据: %d/%d 批次\n", i+1, TOTAL_RECORDS/BATCH_SIZE)
		}
	}

	fmt.Printf("司机数据插入完成，耗时: %v\n", time.Since(start))
}

// 插入乘客数据
func insertPassengers(db *gorm.DB) {
	fmt.Println("开始插入乘客数据...")
	start := time.Now()

	for i := 0; i < TOTAL_RECORDS/BATCH_SIZE; i++ {
		var passengers []Passenger
		for j := 0; j < BATCH_SIZE; j++ {
			id := i*BATCH_SIZE + j + 1
			passenger := Passenger{
				Name:         generateRandomName(),
				Tel:          fmt.Sprintf("1%d%08d", rand.Intn(9)+3, rand.Intn(99999999)+1),
				IDCard:       generateIDCard(),
				RegisterDate: time.Now().Add(-time.Duration(rand.Intn(365*2)) * time.Hour * 24),
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			passengers = append(passengers, passenger)
		}

		if err := db.CreateInBatches(passengers, BATCH_SIZE).Error; err != nil {
			log.Printf("插入乘客数据失败 (批次 %d): %v", i+1, err)
		}

		if (i+1)%100 == 0 {
			fmt.Printf("已插入乘客数据: %d/%d 批次\n", i+1, TOTAL_RECORDS/BATCH_SIZE)
		}
	}

	fmt.Printf("乘客数据插入完成，耗时: %v\n", time.Since(start))
}

// 插入路线数据
func insertRoutes(db *gorm.DB) {
	fmt.Println("开始插入路线数据...")
	start := time.Now()

	for i := 0; i < TOTAL_RECORDS/BATCH_SIZE; i++ {
		var routes []Route
		for j := 0; j < BATCH_SIZE; j++ {
			id := i*BATCH_SIZE + j + 1
			route := Route{
				RouteNo:        fmt.Sprintf("R%06d", id),
				Fieet:          fleets[rand.Intn(len(fleets))],
				Tage:           tags[rand.Intn(len(tags))],
				TagShowStatus:  []string{"show", "hide"}[rand.Intn(2)],
				StartStationId: uint(rand.Intn(1000) + 1),
				EndStationId:   uint(rand.Intn(1000) + 1),
				PassStations:   fmt.Sprintf("%d,%d,%d", rand.Intn(1000)+1, rand.Intn(1000)+1, rand.Intn(1000)+1),
				RoutePath:      fmt.Sprintf("path_%d", id),
				IsActive:       strconv.Itoa(rand.Intn(2)),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			routes = append(routes, route)
		}

		if err := db.CreateInBatches(routes, BATCH_SIZE).Error; err != nil {
			log.Printf("插入路线数据失败 (批次 %d): %v", i+1, err)
		}

		if (i+1)%100 == 0 {
			fmt.Printf("已插入路线数据: %d/%d 批次\n", i+1, TOTAL_RECORDS/BATCH_SIZE)
		}
	}

	fmt.Printf("路线数据插入完成，耗时: %v\n", time.Since(start))
}

// 插入班次数据
func insertBusSchedules(db *gorm.DB) {
	fmt.Println("开始插入班次数据...")
	start := time.Now()

	for i := 0; i < TOTAL_RECORDS/BATCH_SIZE; i++ {
		var schedules []BusSchedules
		for j := 0; j < BATCH_SIZE; j++ {
			id := i*BATCH_SIZE + j + 1
			departureTime := uint(rand.Intn(24*60)) // 分钟数
			schedule := BusSchedules{
				RouteId:       uint(rand.Intn(100000) + 1),
				DepartureTime: departureTime,
				ArrivalTime:   departureTime + uint(rand.Intn(180)+30), // 30-210分钟后到达
				Capacity:      uint(rand.Intn(50) + 20), // 20-70座
				IsActive:      rand.Intn(10) > 1, // 90%概率启用
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			schedules = append(schedules, schedule)
		}

		if err := db.CreateInBatches(schedules, BATCH_SIZE).Error; err != nil {
			log.Printf("插入班次数据失败 (批次 %d): %v", i+1, err)
		}

		if (i+1)%100 == 0 {
			fmt.Printf("已插入班次数据: %d/%d 批次\n", i+1, TOTAL_RECORDS/BATCH_SIZE)
		}
	}

	fmt.Printf("班次数据插入完成，耗时: %v\n", time.Since(start))
}

// 插入路线站点关联数据
func insertRouteStops(db *gorm.DB) {
	fmt.Println("开始插入路线站点关联数据...")
	start := time.Now()

	for i := 0; i < TOTAL_RECORDS/BATCH_SIZE; i++ {
		var routeStops []RouteStops
		for j := 0; j < BATCH_SIZE; j++ {
			id := i*BATCH_SIZE + j + 1
			routeStop := RouteStops{
				RouteId:    rand.Intn(100000) + 1,
				StationId:  rand.Intn(100000) + 1,
				TravelTime: rand.Intn(120), // 0-120分钟
				StopOrder:  rand.Intn(20) + 1, // 1-20站
				StopType:   rand.Intn(2) + 1, // 1-上车，2-下车
				IsActive:   rand.Intn(10) > 1, // 90%概率启用
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			routeStops = append(routeStops, routeStop)
		}

		if err := db.CreateInBatches(routeStops, BATCH_SIZE).Error; err != nil {
			log.Printf("插入路线站点关联数据失败 (批次 %d): %v", i+1, err)
		}

		if (i+1)%100 == 0 {
			fmt.Printf("已插入路线站点关联数据: %d/%d 批次\n", i+1, TOTAL_RECORDS/BATCH_SIZE)
		}
	}

	fmt.Printf("路线站点关联数据插入完成，耗时: %v\n", time.Since(start))
}

// 插入交易数据
func insertTransactions(db *gorm.DB) {
	fmt.Println("开始插入交易数据...")
	start := time.Now()

	for i := 0; i < TOTAL_RECORDS/BATCH_SIZE; i++ {
		var transactions []Transaction
		for j := 0; j < BATCH_SIZE; j++ {
			id := i*BATCH_SIZE + j + 1
			transTime := time.Now().Add(-time.Duration(rand.Intn(365)) * time.Hour * 24)
			transaction := Transaction{
				TransDate:     transTime,
				TransTime:     transTime,
				OrderNo:       fmt.Sprintf("ORD%013d", id),
				DriverId:      uint(rand.Intn(100000) + 1),
				DriverName:    generateRandomName(),
				PassengerId:   uint(rand.Intn(100000) + 1),
				PassengerName: generateRandomName(),
				StartStation:  fmt.Sprintf("%s站", cities[rand.Intn(len(cities))]),
				EndStation:    fmt.Sprintf("%s站", cities[rand.Intn(len(cities))]),
				Amount:        float64(rand.Intn(10000)+500) / 100, // 5-105元
				PaymentMethod: rand.Intn(3) + 1, // 1-3
				TransType:     rand.Intn(3) + 1, // 1-3
				Remark:        fmt.Sprintf("交易备注%d", id),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			transactions = append(transactions, transaction)
		}

		if err := db.CreateInBatches(transactions, BATCH_SIZE).Error; err != nil {
			log.Printf("插入交易数据失败 (批次 %d): %v", i+1, err)
		}

		if (i+1)%100 == 0 {
			fmt.Printf("已插入交易数据: %d/%d 批次\n", i+1, TOTAL_RECORDS/BATCH_SIZE)
		}
	}

	fmt.Printf("交易数据插入完成，耗时: %v\n", time.Since(start))
}

// 插入收支对账数据
func insertBalanceSheets(db *gorm.DB) {
	fmt.Println("开始插入收支对账数据...")
	start := time.Now()

	for i := 0; i < TOTAL_RECORDS/BATCH_SIZE; i++ {
		var balanceSheets []BalanceSheet
		for j := 0; j < BATCH_SIZE; j++ {
			id := i*BATCH_SIZE + j + 1
			settleDate := time.Now().Add(-time.Duration(rand.Intn(365)) * time.Hour * 24)
			orderIncome := float64(rand.Intn(100000)+10000) / 100
			orderRefund := float64(rand.Intn(5000)) / 100
			balanceSheet := BalanceSheet{
				SettleDate:  settleDate,
				OrderIncome: orderIncome,
				OrderRefund: orderRefund,
				OrderSettle: orderIncome - orderRefund,
				Status:      rand.Intn(2) + 1, // 1-正常，2-异常
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			balanceSheets = append(balanceSheets, balanceSheet)
		}

		if err := db.CreateInBatches(balanceSheets, BATCH_SIZE).Error; err != nil {
			log.Printf("插入收支对账数据失败 (批次 %d): %v", i+1, err)
		}

		if (i+1)%100 == 0 {
			fmt.Printf("已插入收支对账数据: %d/%d 批次\n", i+1, TOTAL_RECORDS/BATCH_SIZE)
		}
	}

	fmt.Printf("收支对账数据插入完成，耗时: %v\n", time.Since(start))
}

// 插入司机结算数据
func insertDriverSettles(db *gorm.DB) {
	fmt.Println("开始插入司机结算数据...")
	start := time.Now()

	for i := 0; i < TOTAL_RECORDS/BATCH_SIZE; i++ {
		var driverSettles []DriverSettle
		for j := 0; j < BATCH_SIZE; j++ {
			id := i*BATCH_SIZE + j + 1
			settleDate := time.Now().Add(-time.Duration(rand.Intn(365)) * time.Hour * 24)
			totalIncome := float64(rand.Intn(50000)+10000) / 100
			commission := totalIncome * 0.1 // 10%佣金
			driverSettle := DriverSettle{
				SettleDate:  settleDate,
				DriverId:    uint(rand.Intn(100000) + 1),
				DriverName:  generateRandomName(),
				Phone:       fmt.Sprintf("1%d%08d", rand.Intn(9)+3, rand.Intn(99999999)+1),
				TripCount:   rand.Intn(100) + 10,
				TotalIncome: totalIncome,
				Commission:  commission,
				NetIncome:   totalIncome - commission,
				Status:      rand.Intn(2) + 1, // 1-待结算，2-已结算
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			driverSettles = append(driverSettles, driverSettle)
		}

		if err := db.CreateInBatches(driverSettles, BATCH_SIZE).Error; err != nil {
			log.Printf("插入司机结算数据失败 (批次 %d): %v", i+1, err)
		}

		if (i+1)%100 == 0 {
			fmt.Printf("已插入司机结算数据: %d/%d 批次\n", i+1, TOTAL_RECORDS/BATCH_SIZE)
		}
	}

	fmt.Printf("司机结算数据插入完成，耗时: %v\n", time.Since(start))
}

// 生成随机姓名
func generateRandomName() string {
	surname := surnames[rand.Intn(len(surnames))]
	givenName := givenNames[rand.Intn(len(givenNames))]
	if rand.Intn(2) == 0 {
		// 两个字的名字
		return surname + givenName
	}
	// 三个字的名字
	givenName2 := givenNames[rand.Intn(len(givenNames))]
	return surname + givenName + givenName2
}

// 生成随机身份证号
func generateIDCard() string {
	// 简化版身份证号生成（前6位地区码 + 8位生日 + 4位随机码）
	areaCode := []string{"110101", "310101", "440101", "440301", "330101", "320101", "420101", "510101", "500101", "610101"}
	area := areaCode[rand.Intn(len(areaCode))]
	
	// 生成生日（1970-2000年）
	year := rand.Intn(30) + 1970
	month := rand.Intn(12) + 1
	day := rand.Intn(28) + 1
	birthday := fmt.Sprintf("%04d%02d%02d", year, month, day)
	
	// 4位随机码
	randomCode := fmt.Sprintf("%04d", rand.Intn(10000))
	
	return area + birthday + randomCode
}