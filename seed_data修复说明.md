# seed_data.go 问题修复说明

## 🔍 发现的问题

### 1. 随机数初始化问题
**问题**: `rand.NewSource(time.Now().UnixNano())` 语法错误
**修复**: 改为 `rand.Seed(time.Now().UnixNano())`

### 2. 数据类型问题
**问题**: 返回结构体值而不是指针，导致GORM无法正确处理
**修复**: 所有Generator函数返回指针类型 `&model.StructName{}`

### 3. 时间字段问题
**问题**: BusSchedules表的时间字段使用时间戳，但模型定义为uint分钟数
**修复**: 使用分钟数 (0-1440) 而不是时间戳

### 4. 缺失字段问题
**问题**: Transaction表缺少StartStation、EndStation、Remark等字段
**修复**: 添加所有必需字段

### 5. 缺失表问题
**问题**: 缺少重要的数据表
**修复**: 添加以下表：
- StartStation (站点表)
- RouteStops (路线站点关联表)
- RouteSettle (线路结算表)
- ScheduleSettle (班次结算表)
- StationSettle (站点结算表)
- DriverSettle (司机结算表)

### 6. 数据逻辑问题
**问题**: 收支数据不合理，可能出现负数
**修复**: 确保退款不超过收入，结算金额合理

## 📋 修复后的表清单

现在脚本支持以下9个表，每表100万条数据：

1. **StartStation** - 站点表
2. **Driver** - 司机表
3. **Passenger** - 乘客表
4. **Route** - 线路表
5. **BusSchedules** - 班次表
6. **Transaction** - 交易明细表
7. **RouteStops** - 路线站点关联表
8. **BalanceSheet** - 收支对账表
9. **RouteSettle** - 线路结算表
10. **ScheduleSettle** - 班次结算表
11. **StationSettle** - 站点结算表
12. **DriverSettle** - 司机结算表

## 🚀 使用方法

### Windows系统
```cmd
# 双击运行
run_seed.bat

# 或命令行运行
go run seed_data.go

# 自定义参数
go run seed_data.go -n 10000 -b 1000
```

### Linux/Mac系统
```bash
# 使用脚本
chmod +x run_seed.sh
./run_seed.sh

# 或直接运行
go run seed_data.go

# 自定义参数
go run seed_data.go -n 10000 -b 1000
```

## ⚙️ 参数说明

- `-n`: 每个表插入的数据总量 (默认: 1000000)
- `-b`: 每批插入的数据量 (默认: 2000)

## 🔧 性能优化

修复后的脚本包含以下优化：

1. **连接池优化**: 设置合理的数据库连接池参数
2. **批量插入**: 使用GORM的CreateInBatches方法
3. **进度显示**: 每10批显示一次进度
4. **错误处理**: 完善的错误处理和日志记录
5. **内存优化**: 避免内存泄漏，及时释放资源

## 📊 预期性能

- **插入速度**: 约10,000-50,000条/秒 (取决于硬件)
- **总耗时**: 约20-100分钟 (1200万条数据)
- **内存占用**: < 1GB
- **磁盘空间**: 约5-10GB

## ⚠️ 注意事项

1. **数据库配置**: 确保Nacos配置正确
2. **数据库权限**: 确保有创建表和插入数据的权限
3. **磁盘空间**: 确保有足够的磁盘空间
4. **网络连接**: 确保数据库连接稳定

## 🐛 常见问题

### 问题1: "unsupported data type"
**原因**: 返回结构体值而不是指针
**解决**: 已修复，现在返回指针类型

### 问题2: 随机数初始化失败
**原因**: 语法错误
**解决**: 已修复随机数初始化

### 问题3: 字段类型不匹配
**原因**: 时间字段类型错误
**解决**: 已修复所有字段类型

### 问题4: 数据库连接失败
**原因**: Nacos配置问题
**解决**: 检查core/nacos.go和core/mysql.go配置

## 📝 验证数据

插入完成后，可以使用以下SQL验证：

```sql
-- 检查各表记录数
SELECT 
    table_name,
    table_rows
FROM information_schema.tables 
WHERE table_schema = 'your_database_name'
ORDER BY table_rows DESC;

-- 查看示例数据
SELECT * FROM start_stations LIMIT 5;
SELECT * FROM drivers LIMIT 5;
SELECT * FROM transactions LIMIT 5;
```

## 🎯 总结

修复后的seed_data.go脚本现在可以：
- ✅ 正确插入12个表的数据
- ✅ 每个表100万条记录
- ✅ 高性能批量插入
- ✅ 完善的错误处理
- ✅ 实时进度显示
- ✅ 合理的数据关联关系

脚本已经可以正常运行，直接执行即可开始数据插入。