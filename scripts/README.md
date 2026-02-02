# 批量数据插入脚本

本目录包含用于班车管理系统的批量数据插入脚本，可以快速生成大量测试数据。

## 📋 脚本说明

### 1. Python版本 (推荐)
- **文件**: `bulk_insert_sql.py`
- **特点**: 高性能、多线程、错误处理完善
- **依赖**: `pymysql`

### 2. Go版本
- **文件**: `bulk_insert_data.go`
- **特点**: 使用GORM，类型安全
- **依赖**: `gorm.io/gorm`, `gorm.io/driver/mysql`

### 3. 执行脚本
- **文件**: `run_bulk_insert.sh`
- **功能**: 自动检查环境并执行插入

## 🚀 快速开始

### 方法一：使用Python脚本 (推荐)

```bash
# 1. 安装依赖
pip3 install pymysql

# 2. 修改数据库配置 (如需要)
# 编辑 bulk_insert_sql.py 中的 DB_CONFIG

# 3. 执行插入
python3 bulk_insert_sql.py
```

### 方法二：使用执行脚本

```bash
# 给脚本执行权限
chmod +x run_bulk_insert.sh

# 执行脚本
./run_bulk_insert.sh
```

### 方法三：使用Go脚本

```bash
# 1. 初始化Go模块 (如果还没有)
go mod init bulk_insert

# 2. 安装依赖
go get gorm.io/gorm
go get gorm.io/driver/mysql

# 3. 执行插入
go run bulk_insert_data.go
```

## ⚙️ 配置参数

### Python脚本配置
```python
# 数据库配置
DB_CONFIG = {
    'host': 'localhost',
    'port': 3306,
    'user': 'root',
    'password': 'password',  # 修改为你的密码
    'database': 'bus_management',
    'charset': 'utf8mb4'
}

# 性能配置
TOTAL_RECORDS = 1000000  # 每个表插入100万条数据
BATCH_SIZE = 5000        # 每批插入5000条
THREAD_COUNT = 4         # 并发线程数
```

### Go脚本配置
```go
const (
    DSN = "root:password@tcp(localhost:3306)/bus_management?charset=utf8mb4&parseTime=True&loc=Local"
    BATCH_SIZE = 1000     // 批量插入大小
    TOTAL_RECORDS = 1000000 // 总记录数
)
```

## 📊 数据表说明

脚本将为以下9个表各插入100万条数据：

| 表名 | 说明 | 记录数 |
|------|------|--------|
| `start_stations` | 站点表 | 1,000,000 |
| `drivers` | 司机表 | 1,000,000 |
| `passengers` | 乘客表 | 1,000,000 |
| `routes` | 路线表 | 1,000,000 |
| `bus_schedules` | 班次表 | 1,000,000 |
| `route_stops` | 路线站点关联表 | 1,000,000 |
| `transactions` | 交易明细表 | 1,000,000 |
| `balance_sheets` | 收支对账表 | 1,000,000 |
| `driver_settles` | 司机结算表 | 1,000,000 |

**总计**: 9,000,000 条记录

## 🔧 性能优化

### Python脚本优化特性
- ✅ 多线程并发插入 (默认4线程)
- ✅ 批量插入 (默认5000条/批)
- ✅ 连接池管理
- ✅ 错误处理和重试
- ✅ 进度显示
- ✅ 性能统计

### 预期性能
- **插入速度**: 约 50,000-100,000 记录/秒
- **总耗时**: 约 2-5 分钟 (取决于硬件配置)
- **内存占用**: < 500MB
- **磁盘空间**: 约 2-5GB (取决于数据类型)

## 📝 数据特征

### 生成的数据具有以下特征：
- **真实性**: 使用中国城市名、真实姓名格式
- **多样性**: 随机生成各种组合数据
- **关联性**: 表间数据具有合理的关联关系
- **时间分布**: 数据时间跨度1-3年
- **地理分布**: 覆盖主要城市和地区

### 示例数据：
```
站点: 北京站1, 上海站2, 广州站3...
司机: 张伟, 李娜, 王强...
路线: R000001, R000002, R000003...
交易: ORD0000000000001, ORD0000000000002...
```

## ⚠️ 注意事项

### 执行前准备
1. **数据库准备**: 确保MySQL服务已启动
2. **权限检查**: 确保数据库用户有创建表和插入数据的权限
3. **空间检查**: 确保有足够的磁盘空间 (建议10GB+)
4. **备份数据**: 如果数据库中有重要数据，请先备份

### 执行过程中
- 脚本会自动创建所需的数据表
- 如果表已存在，会继续插入数据 (不会删除现有数据)
- 可以随时 Ctrl+C 中断执行
- 支持断点续传 (重新执行会继续插入)

### 执行后
- 可以通过SQL查询验证数据
- 建议为大表创建适当的索引
- 可以根据需要调整数据或删除测试数据

## 🔍 验证数据

执行完成后，可以使用以下SQL验证数据：

```sql
-- 检查各表记录数
SELECT 'start_stations' as table_name, COUNT(*) as count FROM start_stations
UNION ALL
SELECT 'drivers', COUNT(*) FROM drivers
UNION ALL
SELECT 'passengers', COUNT(*) FROM passengers
UNION ALL
SELECT 'routes', COUNT(*) FROM routes
UNION ALL
SELECT 'bus_schedules', COUNT(*) FROM bus_schedules
UNION ALL
SELECT 'route_stops', COUNT(*) FROM route_stops
UNION ALL
SELECT 'transactions', COUNT(*) FROM transactions
UNION ALL
SELECT 'balance_sheets', COUNT(*) FROM balance_sheets
UNION ALL
SELECT 'driver_settles', COUNT(*) FROM driver_settles;

-- 查看示例数据
SELECT * FROM start_stations LIMIT 5;
SELECT * FROM drivers LIMIT 5;
SELECT * FROM transactions LIMIT 5;
```

## 🛠️ 故障排除

### 常见问题

1. **连接失败**
   ```
   错误: (2003, "Can't connect to MySQL server")
   解决: 检查MySQL服务是否启动，配置是否正确
   ```

2. **权限不足**
   ```
   错误: (1045, "Access denied for user")
   解决: 检查用户名密码，确保有足够权限
   ```

3. **内存不足**
   ```
   解决: 减少 BATCH_SIZE 或 THREAD_COUNT
   ```

4. **磁盘空间不足**
   ```
   解决: 清理磁盘空间或减少 TOTAL_RECORDS
   ```

### 性能调优

- **提高插入速度**: 增加 `THREAD_COUNT` 和 `BATCH_SIZE`
- **降低内存占用**: 减少 `BATCH_SIZE`
- **减少数据量**: 调整 `TOTAL_RECORDS`

## 📞 技术支持

如果遇到问题，请检查：
1. 数据库连接配置
2. Python/Go环境和依赖
3. 系统资源 (内存、磁盘)
4. 数据库权限设置