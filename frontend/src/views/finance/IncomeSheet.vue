<template>
  <div class="income-sheet">
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
        <el-table-column prop="settleDate" label="结算日期(发车时间)" width="150" />
        <el-table-column prop="routeName" label="线路名称" width="120" />
        <el-table-column prop="fleet" label="所属车队" width="100" />
        <el-table-column prop="ticketPrice" label="票价" width="80">
          <template #default="scope">
            {{ formatMoney(scope.row.ticketPrice) }}
          </template>
        </el-table-column>
        <el-table-column prop="discountAmount" label="优惠金额" width="100">
          <template #default="scope">
            {{ formatMoney(scope.row.discountAmount) }}
          </template>
        </el-table-column>
        <el-table-column label="实收金额(已检票)/票数(张)" width="180">
          <template #default="scope">
            {{ formatMoney(scope.row.checkedIncome) }}/{{ scope.row.checkedTickets }}
          </template>
        </el-table-column>
        <el-table-column label="实收金额(未检票)/票数(张)" width="180">
          <template #default="scope">
            {{ formatMoney(scope.row.uncheckedIncome) }}/{{ scope.row.uncheckedTickets }}
          </template>
        </el-table-column>
        <el-table-column label="退款金额(已检票)/票数(张)" width="180">
          <template #default="scope">
            {{ formatMoney(scope.row.checkedRefund) }}/{{ scope.row.checkedRefundTickets }}
          </template>
        </el-table-column>
        <el-table-column label="退款金额(未检票)/票数(张)" width="180">
          <template #default="scope">
            {{ formatMoney(scope.row.uncheckedRefund) }}/{{ scope.row.uncheckedRefundTickets }}
          </template>
        </el-table-column>
      </el-table>
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
  name: 'IncomeSheet',
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

    const formatMoney = (amount) => {
      return Number(amount || 0).toFixed(2)
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

        const response = await financeService.getIncomeSheet(params)
        
        if (response.data.success) {
          tableData.value = response.data.list || []
          total.value = response.data.total || 0
        } else {
          ElMessage.error(response.data.msg || '加载数据失败')
        }
      } catch (error) {
        console.error('加载数据失败:', error)
        // 使用模拟数据
        tableData.value = [
          {
            id: 1,
            settleDate: 'xxxx-xx-xx',
            routeName: '【城际高速】城际合作',
            fleet: '全部车队',
            ticketPrice: 480.00,
            discountAmount: 7.00,
            checkedIncome: 433.00,
            checkedTickets: 11,
            uncheckedIncome: 40.00,
            uncheckedTickets: 1,
            checkedRefund: 0.00,
            checkedRefundTickets: 0,
            uncheckedRefund: 20.00,
            uncheckedRefundTickets: 1
          },
          {
            id: 2,
            settleDate: 'xxxx-xx-xx',
            routeName: '【城际高速】城际合作',
            fleet: '全部车队',
            ticketPrice: 480.00,
            discountAmount: 7.00,
            checkedIncome: 433.00,
            checkedTickets: 11,
            uncheckedIncome: 40.00,
            uncheckedTickets: 1,
            checkedRefund: 0.00,
            checkedRefundTickets: 0,
            uncheckedRefund: 20.00,
            uncheckedRefundTickets: 1
          },
          {
            id: 3,
            settleDate: 'xxxx-xx-xx',
            routeName: '【城际高速】城际合作',
            fleet: '全部车队',
            ticketPrice: 480.00,
            discountAmount: 7.00,
            checkedIncome: 433.00,
            checkedTickets: 11,
            uncheckedIncome: 40.00,
            uncheckedTickets: 1,
            checkedRefund: 0.00,
            checkedRefundTickets: 0,
            uncheckedRefund: 20.00,
            uncheckedRefundTickets: 1
          }
        ]
        total.value = 3
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

    const exportData = async () => {
      try {
        const params = {
          exportType: 2, // 收入对账
          startDate: searchForm.value.dateRange?.[0] || '',
          endDate: searchForm.value.dateRange?.[1] || ''
        }

        const response = await financeService.export(params)
        if (response.data.success) {
          ElMessage.success('导出成功')
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
      formatMoney,
      handlePageChange,
      search,
      exportData
    }
  }
}
</script>

<style scoped>
.income-sheet {
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

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>