import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建 axios 实例
const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    // 可以在这里添加 token
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response
  },
  error => {
    ElMessage.error(error.message || '请求失败')
    return Promise.reject(error)
  }
)

// 班线管理 API
export const routeApi = {
  // 创建班线
  create: (data) => api.post('/route/create', data),
  // 班线列表
  list: (params) => api.get('/route/list', { params }),
  // 班线详情
  info: (params) => api.get('/route/info', { params }),
  // 更新班线
  update: (data) => api.put('/route/update', data),
  // 删除班线
  delete: (params) => api.delete('/route/delete', { params }),
  // 添加站点到班线
  addStation: (data) => api.post('/route/station/add', data),
  // 移除站点
  removeStation: (params) => api.delete('/route/station/remove', { params })
}

// 班次管理 API
export const scheduleApi = {
  // 创建班次
  create: (data) => api.post('/schedule/create', data),
  // 班次列表
  list: (params) => api.get('/schedule/list', { params }),
  // 更新班次
  update: (data) => api.put('/schedule/update', data),
  // 删除班次
  delete: (params) => api.delete('/schedule/delete', { params })
}

// 财务管理 API (端口8895)
const financeApi = axios.create({
  baseURL: '/api/finance',
  timeout: 10000
})

financeApi.interceptors.response.use(
  response => response,
  error => {
    ElMessage.error(error.message || '请求失败')
    return Promise.reject(error)
  }
)

export const financeService = {
  // 收支对账
  getBalanceSheet: (params) => financeApi.get('/balance', { params }),
  // 收入对账
  getIncomeSheet: (params) => financeApi.get('/income', { params }),
  // 线路结算
  getRouteSettle: (params) => financeApi.get('/route-settle', { params }),
  // 班次结算
  getScheduleSettle: (params) => financeApi.get('/schedule-settle', { params }),
  // 站点结算
  getStationSettle: (params) => financeApi.get('/station-settle', { params }),
  // 司机结算
  getDriverSettle: (params) => financeApi.get('/driver-settle', { params }),
  // 交易明细
  getTransactions: (params) => financeApi.get('/transactions', { params }),
  // 异常交易
  getAbnormalTrans: (params) => financeApi.get('/abnormal', { params }),
  // 处理异常交易
  handleAbnormal: (data) => financeApi.post('/abnormal/handle', data),
  // 生成账单
  generateBill: (data) => financeApi.post('/bill/generate', data),
  // 导出
  export: (params) => financeApi.get('/export', { params })
}

export default api