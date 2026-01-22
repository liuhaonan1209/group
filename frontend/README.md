# 拼约出行管理平台前端

基于 Vue 3 + Element Plus 的班车管理系统前端界面。

## 功能模块

### 1. 班车管理
- **站点选择** - 支持地图选点、站点搜索、上下车站点配置
- **线路创建** - 分步骤创建班线，配置基础信息和站点
- **班线列表** - 班线数据展示和管理

### 2. 收支对账系统
- **收支对账表** - 日常收支数据统计和异常监控
- **收入对账表** - 详细的收入明细和票务统计
- **线路结算表** - 按线路维度的收支结算
- **班次结算表** - 按班次维度的收支结算
- **站点结算表** - 按站点维度的收支结算
- **司机结算表** - 司机收入和佣金结算

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **Element Plus** - Vue 3 组件库
- **Vue Router** - 路由管理
- **Pinia** - 状态管理
- **Axios** - HTTP 客户端
- **Vite** - 构建工具

## 项目结构

```
frontend/
├── src/
│   ├── components/          # 公共组件
│   │   └── Layout.vue      # 布局组件
│   ├── views/              # 页面组件
│   │   ├── route/          # 班车管理
│   │   │   ├── RouteList.vue
│   │   │   ├── RouteCreate.vue
│   │   │   └── StationSelect.vue
│   │   └── finance/        # 收支对账
│   │       ├── BalanceSheet.vue
│   │       ├── IncomeSheet.vue
│   │       ├── RouteSettle.vue
│   │       ├── ScheduleSettle.vue
│   │       ├── StationSettle.vue
│   │       └── DriverSettle.vue
│   ├── router/             # 路由配置
│   ├── api/                # API 接口
│   └── main.js            # 入口文件
├── package.json
├── vite.config.js
└── index.html
```

## 安装和运行

### 1. 安装依赖
```bash
cd frontend
npm install
```

### 2. 启动开发服务器
```bash
npm run dev
```

### 3. 构建生产版本
```bash
npm run build
```

## 界面特性

### 1. 站点选择界面
- 弹窗式站点选择器
- 支持关键词搜索
- 地图集成预留接口
- 站点信息表单填写
- 分页展示站点列表

### 2. 线路创建流程
- 三步骤创建向导
- 基础信息配置
- 上下车站点配置
- 实时预览和验证

### 3. 收支对账界面
- 统一的搜索和筛选
- 数据表格展示
- 汇总信息显示
- 导出功能支持
- 分页导航

### 4. 响应式设计
- 适配不同屏幕尺寸
- 移动端友好
- 现代化 UI 设计

## API 接口

前端通过 Axios 调用后端 API，支持：

- 班线管理 API (`/api/route/*`)
- 收支对账 API (`/api/finance/*`)
- 自动错误处理和消息提示
- 请求/响应拦截器

## 开发说明

### 1. 添加新页面
1. 在 `src/views/` 下创建 Vue 组件
2. 在 `src/router/index.js` 中添加路由
3. 在 `Layout.vue` 中添加菜单项

### 2. 调用 API
```javascript
import api from '@/api'

// GET 请求
const response = await api.get('/route/list', { params })

// POST 请求
const response = await api.post('/route/create', data)
```

### 3. 样式规范
- 使用 Element Plus 主题色
- 统一的间距和圆角
- 响应式布局
- 阴影和过渡效果

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88

## 部署

构建后的文件在 `dist/` 目录，可部署到任何静态文件服务器。

推荐使用 Nginx 配置：

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://backend-server:8893;
    }
}
```