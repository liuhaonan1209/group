<template>
  <div class="schedule-manage">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="searchForm" inline>
        <el-form-item>
          <el-select v-model="searchForm.routeId" placeholder="选择班线" filterable clearable>
            <el-option 
              v-for="route in routeList" 
              :key="route.id" 
              :label="route.routeNo" 
              :value="route.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="warning" @click="search">搜索</el-button>
        </el-form-item>
      </el-form>
      
      <div class="toolbar">
        <el-button type="primary" @click="createSchedule">新增班次</el-button>
      </div>
    </div>

    <!-- 数据表格 -->
    <div class="table-container">
      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="routeNo" label="班线编号" width="120" />
        <el-table-column prop="departureTime" label="发车时间" width="100" />
        <el-table-column prop="arrivalTime" label="预计到达时间" width="120" />
        <el-table-column prop="capacity" label="座位容量" width="100" />
        <el-table-column prop="isActive" label="状态" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.isActive ? 'success' : 'danger'">
              {{ scope.row.isActive ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button type="text" size="small" @click="editSchedule(scope.row)">编辑</el-button>
            <el-button type="text" size="small" @click="deleteSchedule(scope.row)" style="color: #f56c6c">删除</el-button>
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

    <!-- 班次编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑班次' : '新增班次'" width="500px">
      <el-form :model="scheduleForm" :rules="rules" ref="scheduleFormRef" label-width="100px">
        <el-form-item label="班线：" prop="routeId">
          <el-select v-model="scheduleForm.routeId" placeholder="请选择班线" style="width: 100%">
            <el-option 
              v-for="route in routeList" 
              :key="route.id" 
              :label="route.routeNo" 
              :value="route.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="发车时间：" prop="departureTime">
          <el-time-picker
            v-model="scheduleForm.departureTime"
            format="HH:mm"
            value-format="HH:mm"
            placeholder="选择发车时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="到达时间：" prop="arrivalTime">
          <el-time-picker
            v-model="scheduleForm.arrivalTime"
            format="HH:mm"
            value-format="HH:mm"
            placeholder="选择到达时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="座位容量：" prop="capacity">
          <el-input-number 
            v-model="scheduleForm.capacity" 
            :min="1" 
            :max="100"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="状态：" v-if="isEdit">
          <el-switch v-model="scheduleForm.isActive" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitSchedule" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { scheduleApi, routeApi } from '@/api'

export default {
  name: 'ScheduleManage',
  setup() {
    const loading = ref(false)
    const submitting = ref(false)
    const dialogVisible = ref(false)
    const isEdit = ref(false)
    const scheduleFormRef = ref(null)
    
    const searchForm = ref({
      routeId: null
    })

    const currentPage = ref(1)
    const pageSize = ref(10)
    const total = ref(0)
    const tableData = ref([])
    const routeList = ref([])

    const scheduleForm = ref({
      id: null,
      routeId: null,
      departureTime: '',
      arrivalTime: '',
      capacity: 30,
      isActive: true
    })

    const rules = {
      routeId: [
        { required: true, message: '请选择班线', trigger: 'change' }
      ],
      departureTime: [
        { required: true, message: '请选择发车时间', trigger: 'change' }
      ],
      arrivalTime: [
        { required: true, message: '请选择到达时间', trigger: 'change' }
      ],
      capacity: [
        { required: true, message: '请输入座位容量', trigger: 'blur' }
      ]
    }

    const loadRoutes = async () => {
      try {
        const response = await routeApi.list({ page: 1, size: 100 })
        if (response.data.success) {
          routeList.value = response.data.list || []
        }
      } catch (error) {
        console.error('加载班线失败:', error)
      }
    }

    const loadData = async () => {
      loading.value = true
      try {
        const params = {
          page: currentPage.value,
          size: pageSize.value,
          routeId: searchForm.value.routeId || ''
        }

        const response = await scheduleApi.list(params)
        
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

    const createSchedule = () => {
      isEdit.value = false
      scheduleForm.value = {
        id: null,
        routeId: null,
        departureTime: '',
        arrivalTime: '',
        capacity: 30,
        isActive: true
      }
      dialogVisible.value = true
    }

    const editSchedule = (row) => {
      isEdit.value = true
      scheduleForm.value = {
        id: row.id,
        routeId: row.routeId,
        departureTime: row.departureTime,
        arrivalTime: row.arrivalTime,
        capacity: row.capacity,
        isActive: row.isActive
      }
      dialogVisible.value = true
    }

    const submitSchedule = async () => {
      const valid = await scheduleFormRef.value.validate().catch(() => false)
      if (!valid) return

      submitting.value = true
      try {
        let response
        if (isEdit.value) {
          response = await scheduleApi.update(scheduleForm.value)
        } else {
          response = await scheduleApi.create(scheduleForm.value)
        }

        if (response.data.success) {
          ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
          dialogVisible.value = false
          loadData()
        } else {
          ElMessage.error(response.data.msg || '操作失败')
        }
      } catch (error) {
        ElMessage.error('操作失败')
      } finally {
        submitting.value = false
      }
    }

    const deleteSchedule = async (row) => {
      try {
        await ElMessageBox.confirm('确定要删除这个班次吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })

        const response = await scheduleApi.delete({ id: row.id })
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
      loadRoutes()
      loadData()
    })

    return {
      loading,
      submitting,
      dialogVisible,
      isEdit,
      scheduleFormRef,
      searchForm,
      currentPage,
      pageSize,
      total,
      tableData,
      routeList,
      scheduleForm,
      rules,
      search,
      handlePageChange,
      createSchedule,
      editSchedule,
      submitSchedule,
      deleteSchedule
    }
  }
}
</script>

<style scoped>
.schedule-manage {
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
</style>