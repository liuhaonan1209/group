# 拼约出行管理平台架构文档

## 系统架构概览

```mermaid
graph TB
    subgraph "前端层 Frontend Layer"
        FE[Vue 3 + Element Plus<br/>端口: 3000]
        FE --> |HTTP/HTTPS| NGINX[Nginx 反向代理]
    end

    subgraph "API 网关层 API Gateway Layer"
        NGINX --> |路由分发| API1[班线管理 API<br/>端口: 8893]
        NGINX --> |路由分发| API2[财务管理 API<br/>端口: 8895]
        NGINX --> |路由分发| API3[站点管理 API<br/>端口: 8891]
        NGINX --> |路由分发| API4[司机管理 API<br/>端口: 8889]
        NGINX --> |路由分发| API5[乘客管理 API<br/>端口: 8890]
    end

    subgraph "业务服务层 Business Service Layer"
        API1 --> |Kitex RPC| RPC1[班线 RPC 服务<br/>端口: 8892]
        API2 --> |Kitex RPC| RPC2[财务 RPC 服务<br/>端口: 8894]
        API3 --> |Kitex RPC| RPC3[站点 RPC 服务<br/>端口: 8888]
        API4 --> |Kitex RPC| RPC4[司机 RPC 服务<br/>端口: 8888]
        API5 --> |Kitex RPC| RPC5[乘客 RPC 服务<br/>端口: 8888]
        RPC6[订单 RPC 服务<br/>端口: 8888]
    end

    subgraph "数据访问层 Data Access Layer"
        RPC1 --> DAO1[Route DAO]
        RPC2 --> DAO2[Finance DAO]
        RPC3 --> DAO3[Station DAO]
        RPC4 --> DAO4[Driver DAO]
        RPC5 --> DAO5[Passenger DAO]
        RPC6 --> DAO6[Order DAO]
    end

    subgraph "数据存储层 Data Storage Layer"
        DAO1 --> DB[(MySQL 数据库)]
        DAO2 --> DB
        DAO3 --> DB
        DAO4 --> DB
        DAO5 --> DB
        DAO6 --> DB
    end

    subgraph "配置与注册中心"
        NACOS[Nacos 配置中心]
        RPC1 -.-> NACOS
        RPC2 -.-> NACOS
        RPC3 -.-> NACOS
        RPC4 -.-> NACOS
        RPC5 -.-> NACOS
        RPC6 -.-> NACOS
    end

    style FE fill:#e1f5fe
    style API1 fill:#f3e5f5
    style API2 fill:#f3e5f5
    style API3 fill:#f3e5f5
    style API4 fill:#f3e5f5
    style API5 fill:#f3e5f5
    style RPC1 fill:#e8f5e8
    style RPC2 fill:#e8f5e8
    style RPC3 fill:#e8f5e8
    style RPC4 fill:#e8f5e8
    style RPC5 fill:#e8f5e8
    style RPC6 fill:#e8f5e8
    style DB fill:#fff3e0
    style NACOS fill:#fce4ec
```

## 技术栈详情

### 前端技术栈
- **框架**: Vue 3 + Composition API
- **UI 组件库**: Element Plus
- **路由**: Vue Router 4
- **状态管理**: Pinia
- **HTTP 客户端**: Axios
- **构建工具**: Vite
- **开发语言**: JavaScript/TypeScript

### 后端技术栈
- **开发语言**: Go 1.19+
- **RPC 框架**: Kitex (字节跳动)
- **HTTP 框架**: Hertz (字节跳动)
- **IDL**: Thrift
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **配置中心**: Nacos
- **服务注册与发现**: Nacos

## 模块架构详解

### 1. 班线管理模块
```mermaid
graph LR
    subgraph "班线管理"
        A[前端页面] --> B[班线 API<br/>8893]
        B --> C[班线 RPC<br/>8892]
        C --> D[Route DAO]
        D --> E[MySQL]
        
        C --> F[RouteStops DAO]
        F --> E
        
        C --> G[BusSchedules DAO]
        G --> E
    end
```

**核心功能**:
- 班线 CRUD 操作
- 上下车站点管理（排重、排序）
- 班次管理
- 站点时间自动排序（起点为0）

### 2. 财务管理模块
```mermaid
graph LR
    subgraph "财务管理"
        A[前端页面] --> B[财务 API<br/>8895]
        B --> C[财务 RPC<br/>8894]
        C --> D[Finance DAO]
        D --> E[MySQL]
        
        C --> F[异常检测]
        F --> G[自动对账]
        G --> H[账单生成]
    end
```

**核心功能**:
- 收支对账表
- 收入对账表
- 多维度结算（线路、班次、站点、司机）
- 异常交易检测与处理
- 自动账单生成

### 3. 站点管理模块
```mermaid
graph LR
    subgraph "站点管理"
        A[站点选择] --> B[站点 API]
        B --> C[站点 RPC<br/>8888]
        C --> D[Station DAO]
        D --> E[MySQL]
        
        C --> F[地图集成]
        F --> G[位置服务]
    end
```

## 数据模型关系

```mermaid
erDiagram
    Route ||--o{ RouteStops : contains
    Route ||--o{ BusSchedules : has
    StartStation ||--o{ RouteStops : belongs_to
    Route }o--|| StartStation : start_station
    Route }o--|| StartStation : end_station
    
    Route {
        uint ID
        string RouteNo
        string Fieet
        string Tage
        uint StartStationId
        uint EndStationId
        bool IsActive
    }
    
    RouteStops {
        uint ID
        int RouteId
        int StationId
        int TravelTime
        int StopOrder
        int StopType
        bool IsActive
    }
    
    StartStation {
        uint ID
        string Name
        float64 Latitude
        float64 Longitude
        string Address
        bool IsActive
    }
    
    BusSchedules {
        uint ID
        uint RouteId
        uint DepartureTime
        uint ArrivalTime
        uint Capacity
        bool IsActive
    }
    
    BalanceSheet {
        uint ID
        date SettleDate
        decimal OrderIncome
        decimal OrderRefund
        decimal OrderSettle
        int Status
    }
    
    Transaction {
        uint ID
        date TransDate
        datetime TransTime
        string OrderNo
        uint DriverId
        uint PassengerId
        decimal Amount
        int PaymentMethod
        int TransType
    }
```

## 服务通信流程

### 班线创建流程
```mermaid
sequenceDiagram
    participant F as 前端
    participant A as 班线API
    participant R as 班线RPC
    participant D as DAO层
    participant DB as MySQL

    F->>A: POST /api/route/create
    A->>R: RouteCreate(req)
    R->>D: ValidateBusLineStops()
    D->>D: 站点排重检查
    D->>D: 时间排序处理
    R->>D: RouteCreate()
    D->>DB: INSERT INTO routes
    R->>D: AddRouteStops()
    D->>DB: INSERT INTO route_stops
    DB-->>D: 返回结果
    D-->>R: 返回结果
    R-->>A: RouteCreateResp
    A-->>F: JSON Response
```

### 财务对账流程
```mermaid
sequenceDiagram
    participant F as 前端
    participant A as 财务API
    participant R as 财务RPC
    participant D as DAO层
    participant DB as MySQL

    F->>A: GET /api/finance/balance
    A->>R: GetBalanceSheet(req)
    R->>D: GetBalanceSheetList()
    D->>DB: SELECT FROM balance_sheets
    R->>D: GetBalanceSheetSummary()
    D->>DB: SELECT SUM() FROM balance_sheets
    DB-->>D: 返回数据
    D-->>R: 返回结果
    R-->>A: BalanceSheetResp
    A-->>F: JSON Response
```

## 部署架构

```mermaid
graph TB
    subgraph "负载均衡层"
        LB[负载均衡器<br/>Nginx/HAProxy]
    end
    
    subgraph "Web 服务器集群"
        WEB1[Web Server 1<br/>前端静态资源]
        WEB2[Web Server 2<br/>前端静态资源]
    end
    
    subgraph "API 网关集群"
        GW1[API Gateway 1<br/>路由分发]
        GW2[API Gateway 2<br/>路由分发]
    end
    
    subgraph "微服务集群"
        MS1[班线服务集群]
        MS2[财务服务集群]
        MS3[站点服务集群]
        MS4[司机服务集群]
        MS5[乘客服务集群]
    end
    
    subgraph "数据层"
        DB1[(MySQL 主库)]
        DB2[(MySQL 从库)]
        REDIS[(Redis 缓存)]
    end
    
    subgraph "基础设施"
        NACOS[Nacos 注册中心]
        LOG[日志系统<br/>ELK Stack]
        MONITOR[监控系统<br/>Prometheus]
    end
    
    LB --> WEB1
    LB --> WEB2
    LB --> GW1
    LB --> GW2
    
    GW1 --> MS1
    GW1 --> MS2
    GW2 --> MS3
    GW2 --> MS4
    GW2 --> MS5
    
    MS1 --> DB1
    MS2 --> DB1
    MS3 --> DB2
    MS4 --> DB2
    MS5 --> DB2
    
    MS1 --> REDIS
    MS2 --> REDIS
    MS3 --> REDIS
    
    MS1 -.-> NACOS
    MS2 -.-> NACOS
    MS3 -.-> NACOS
    MS4 -.-> NACOS
    MS5 -.-> NACOS
```

## 关键特性

### 1. 微服务架构
- **服务拆分**: 按业务领域拆分为独立的微服务
- **服务通信**: 使用 Kitex RPC 框架进行高性能通信
- **服务注册**: 基于 Nacos 的服务注册与发现

### 2. 数据一致性
- **事务管理**: 使用 GORM 事务确保数据一致性
- **分布式事务**: 关键业务流程采用分布式事务
- **数据同步**: 实时同步交易数据到财务系统

### 3. 高可用设计
- **负载均衡**: 多实例部署，支持水平扩展
- **故障转移**: 服务实例故障自动切换
- **数据备份**: 主从复制，定期备份

### 4. 安全性
- **API 鉴权**: JWT Token 认证
- **数据加密**: 敏感数据加密存储
- **访问控制**: 基于角色的权限控制

### 5. 监控与运维
- **服务监控**: Prometheus + Grafana
- **日志收集**: ELK Stack
- **链路追踪**: Jaeger 分布式追踪
- **健康检查**: 服务健康状态监控

## 性能优化

### 1. 缓存策略
- **Redis 缓存**: 热点数据缓存
- **本地缓存**: 配置信息本地缓存
- **CDN**: 静态资源 CDN 加速

### 2. 数据库优化
- **读写分离**: 主库写入，从库读取
- **索引优化**: 关键字段建立索引
- **分库分表**: 大表水平拆分

### 3. 接口优化
- **批量操作**: 减少数据库交互次数
- **异步处理**: 耗时操作异步执行
- **连接池**: 数据库连接池管理

这个架构设计确保了系统的可扩展性、可维护性和高性能，能够支撑大规模的班车管理业务需求。