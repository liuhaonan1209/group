<template>
  <div class="route-list">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="searchForm" inline>
        <el-form-item>
          <el-select v-model="searchForm.fieet" placeholder="车队">
            <el-option label="全部" value=""></el-option>
            <el-option label="车队A" value="fleetA"></el-option>
            <el-option label="车队B" value="fleetB"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-input v-model="searchForm.routeNo" placeholder="班线编号" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="searchForm.tage" placeholder="标签" />
        </el-form-item>
        <el-form-item>
          <el-button type="warning" @click="search">搜索</el-button>
        </el-form-item>
      </el-form>
      
      <div class="toolbar">
        <el-button type="primary" @click="createRoute">新增班线</el-button>
      </div>
    </div>

    <!-- 数据表格 -->
    <div class="table-container">
      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="routeNo" label="班线编号" width="120" />
        <el-table-column prop="fieet" label="车队" width="100" />
        <el-table-column prop="tage" label="标签" width="100" />
        <el-table-column prop="startStationName" label="起点站" width="150" />
        <el-table-column prop="endStationName" label="终点站" width="150" />
        <el-table-column prop="isActive" label="状态" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.isActive ? 'success' : 'danger'">
              {{ scope.row.isActive ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button type="text" size="small" @click="viewDetail(scope.row)">详情</el-button>
            <el-button type="text" size="small" @click="editRoute(scope.row)">编辑</el-button>
            <el-button type="text" size="small" @click="deleteRoute(scope.row)" style="color: #f56c6c">删除</el-button>
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

    <!-- 班线详情弹窗 -->
    <el-dialog v-model="detailVisible" title="班线详情" width="80%">
      <div v-if="currentRoute">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="班线编号">{{ currentRoute.routeNo }}</el-descriptions-item>
          <el-descriptions-item label="车队">{{ currentRoute.fieet }}</el-descriptions-item>
          <el-descriptions-item label="标签">{{ currentRoute.tage }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="currentRoute.isActive ? 'success' : 'danger'">
              {{ currentRoute.isActive ? '启用' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <div class="station-info">
          <h4>上车站点</h4>
          <el-table :data="currentRoute.upStations" style="width: 100%; margin-bottom: 20px;">
            <el-table-column prop="name" label="站点名称" />
            <el-table-column prop="travelTime" label="行驶时间(分钟)" />
            <el-table-column prop="stopOrder" label="停靠顺序" />
          </el-table>

          <h4>下车站点</h4>
          <el-table :data="currentRoute.downStations" style="width: 100%;">
            <el-table-column prop="name" label="站点名称" />
            <el-table-column prop="travelTime" label="行驶时间(分钟)" />
            <el-table-column prop="stopOrder" label="停靠顺序" />
          </el-table>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { routeApi } from '@/api'
import { useRouter } from 'vue-router'

export default {
  name: 'RouteList',
  setup() {
    const router = useRouter()
    const loading = ref(false)
    const detailVisible = ref(false)
    const currentRoute = ref(null)
    
    const searchForm = ref({
      routeNo: '',
      fieet: '',
      tage: ''
    })

    const currentPage = ref(1)
    const pageSize = ref(10)
    const total = ref(0)
    const tableData = ref([])

    const loadData = async () => {
      loading.value = true
      try {
        const params = {
          page: currentPage.value,
          size: pageSize.value,
          ...searchForm.value
        }

        const response = await routeApi.list(params)
        
        if (response.data.success) {
          tableData.value = response.data.list || []
          total.value = response.data.total || 0
        } else {
          ElMessage.error(response.data.msg || '加载数据失败')
        }
      } catch (error) {
        console.error('加载数据失败:', error)
        ElMessage.error('加载数据失败')
      } finally {
        loading.value = false
      }
    }

    const search = () => {
      currentPage.value = 1
      loadData()
    }

    const handlePageChange = (page) => {
      currentPage.value = page
      loadData()
    }

    const createRoute = () => {
      router.push('/route/create')
    }

    const viewDetail = async (row) => {
      try {
        const response = await routeApi.info({ id: row.id })
        if (response.data.success) {
          currentRoute.value = response.data
          detailVisible.value = true
        } else {
          ElMessage.error(response.data.msg || '获取详情失败')
        }
      } catch (error) {
        ElMessage.error('获取详情失败')
      }
    }

    const editRoute = (row) => {
      router.push(`/route/edit/${row.id}`)
    }

    const deleteRoute = async (row) => {
      try {
        await ElMessageBox.confirm('确定要删除这条班线吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })

        const response = await routeApi.delete({ id: row.id })
        if (response.data.success) {
          ElMessage.success('删除成功')
          loadData()
        } else {
          ElMessage.error(response.data.msg || '删除失败')
        }
      } catch (error) {
        if (error !== 'cancel') {
          ElMessage.error('删除失败')
        }
      }
    }

    onMounted(() => {
      loadData()
    })

    return {
      loading,
      detailVisible,
      currentRoute,
      searchForm,
      currentPage,
      pageSize,
      total,
      tableData,
      search,
      handlePageChange,
      createRoute,
      viewDetail,
      editRoute,
      deleteRoute
    }
  }
}
</script>

<style scoped>
.route-list {
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
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.station-info {
  margin-top: 20px;
}

.station-info h4 {
  margin: 20px 0 10px 0;
  color: #333;
}
</style>