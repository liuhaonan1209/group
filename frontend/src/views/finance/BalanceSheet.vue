<template>
  <div class="balance-sheet">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="searchForm" inline>
        <el-form-item>
          <el-select v-model="searchForm.status" placeholder="正常">
            <el-option label="全部" value="0"></el-option>
            <el-option label="正常" value="1"></el-option>
            <el-option label="异常" value="2"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="searchForm.keyword" placeholder="时间：" />
        </el-form-item>
        <el-form-item>
          <el-date-picker
            v-model="searchForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="warning" @click="search">搜索</el-button>
        </el-form-item>
      </el-form>
      
      <div class="toolbar">
        <el-button type="primary" @click="exportData">导出</el-button>
      </div>
    </div>

    <!-- 数据表格 -->
    <div class="table-container">
      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="settleDate" label="对账日期" width="120" />
        <el-table-column prop="orderIncome" label="订单实收" width="120">
          <template #default="scope">
            {{ formatMoney(scope.row.orderIncome) }}
          </template>
        </el-table-column>
        <el-table-column prop="orderRefund" label="订单退款" width="120">
          <template #default="scope">
            {{ formatMoney(scope.row.orderRefund) }}
          </template>
        </el-table-column>
        <el-table-column prop="orderSettle" label="订单结算" width="120">
          <template #default="scope">
            {{ formatMoney(scope.row.orderSettle) }}
          </template>
        </el-table-column>
        <el-table-column prop="statusText" label="账目状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="scope">
            <el-button type="text" size="small" @click="viewDetail(scope.row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 汇总信息 -->
    <div class="summary-info">
      <div class="summary-item">
        <span class="label">总收入：</span>
        <span class="value income">{{ formatMoney(summaryData.totalIncome) }}</span>
      </div>
      <div class="summary-item">
        <span class="label">总退款：</span>
        <span class="value refund">{{ formatMoney(summaryData.totalRefund) }}</span>
      </div>
      <div class="summary-item">
        <span class="label">总结算：</span>
        <span class="value settle">{{ formatMoney(summaryData.totalSettle) }}</span>
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { financeService } from '@/api'

export default {
  name: 'BalanceSheet',
  setup() {
    const loading = ref(false)
    const searchForm = ref({
      status: '0',
      keyword: '',
      dateRange: []
    })

    const currentPage = ref(1)
    const pageSize = ref(10)
    const total = ref(0)
    const tableData = ref([])
    const summaryData = ref({
      totalIncome: 0,
      totalRefund: 0,
      totalSettle: 0
    })

    const formatMoney = (amount) => {
      return `¥${Number(amount || 0).toFixed(2)}`
    }

    const loadData = async () => {
      loading.value = true
      try {
        const params = {
          page: currentPage.value,
          size: pageSize.value,
          status: searchForm.value.status === '0' ? '' : searchForm.value.status
        }

        if (searchForm.value.dateRange && searchForm.value.dateRange.length === 2) {
          params.startDate = searchForm.value.dateRange[0]
          params.endDate = searchForm.value.dateRange[1]
        }

        const response = await financeService.getBalanceSheet(params)
        
        if (response.data.success) {
          tableData.value = response.data.list || []
          total.value = response.data.total || 0
          summaryData.value = {
            totalIncome: response.data.totalIncome || 0,
            totalRefund: response.data.totalRefund || 0,
            totalSettle: response.data.totalSettle || 0
          }
        } else {
          ElMessage.error(response.data.msg || '加载数据失败')
        }
      } catch (error) {
        console.error('加载数据失败:', error)
        ElMessage.error('加载数据失败')
        // 使用模拟数据
        tableData.value = [
          {
            id: 1,
            settleDate: '2024-01-15',
            orderIncome: 7865.9,
            orderRefund: -190.6,
            orderSettle: 7675.3,
            status: 1,
            statusText: '正常'
          },
          {
            id: 2,
            settleDate: '2024-01-14',
            orderIncome: 8200.5,
            orderRefund: -250.0,
            orderSettle: 7950.5,
            status: 1,
            statusText: '正常'
          },
          {
            id: 3,
            settleDate: '2024-01-13',
            orderIncome: 7500.0,
            orderRefund: -180.0,
            orderSettle: 7320.0,
            status: 2,
            statusText: '异常'
          }
        ]
        total.value = 3
        summaryData.value = {
          totalIncome: 23566.4,
          totalRefund: -620.6,
          totalSettle: 22945.8
        }
      } finally {
        loading.value = false
      }
    }

    const handlePageChange = (page) => {
      currentPage.value = page
      loadData()
    }

    const search = () => {
      currentPage.value = 1
      loadData()
    }

    const viewDetail = (row) => {
      ElMessage.info(`查看 ${row.settleDate} 的详细信息`)
    }

    const exportData = async () => {
      try {
        const params = {
          exportType: 1, // 收支对账
          startDate: searchForm.value.dateRange?.[0] || '',
          endDate: searchForm.value.dateRange?.[1] || ''
        }

        const response = await financeService.export(params)
        if (response.data.success) {
          ElMessage.success('导出成功')
          // 这里可以处理文件下载
          window.open(response.data.fileUrl)
        } else {
          ElMessage.error(response.data.msg || '导出失败')
        }
      } catch (error) {
        ElMessage.error('导出失败')
      }
    }

    onMounted(() => {
      loadData()
    })

    return {
      loading,
      searchForm,
      currentPage,
      pageSize,
      total,
      tableData,
      summaryData,
      formatMoney,
      handlePageChange,
      search,
      viewDetail,
      exportData
    }
  }
}
</script>

<style scoped>
.balance-sheet {
  padding: 20px;
}

.search-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 20px;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.toolbar {
  display: flex;
  gap: 10px;
}

.table-container {
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  margin-bottom: 20px;
}

.summary-info {
  display: flex;
  gap: 30px;
  padding: 20px;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  margin-bottom: 20px;
}

.summary-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.summary-item .label {
  font-weight: 500;
  color: #666;
}

.summary-item .value {
  font-size: 18px;
  font-weight: bold;
}

.summary-item .value.income {
  color: #67C23A;
}

.summary-item .value.refund {
  color: #F56C6C;
}

.summary-item .value.settle {
  color: #409EFF;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>