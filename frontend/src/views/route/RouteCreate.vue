<template>
  <div class="route-create">
    <!-- 步骤导航 -->
    <div class="steps-container">
      <el-steps :active="currentStep" align-center>
        <el-step title="第一步" description="基础信息"></el-step>
        <el-step title="第二步" description="配置线路站点"></el-step>
        <el-step title="第三步" description="完成"></el-step>
      </el-steps>
    </div>

    <!-- 第一步：基础信息 -->
    <div v-if="currentStep === 0" class="step-content">
      <div class="form-container">
        <h3>第一步 基础信息</h3>
        
        <el-form :model="routeForm" :rules="rules" ref="routeFormRef" label-width="100px" class="route-form">
          <el-form-item label="班线编号：" prop="routeNo">
            <el-input 
              v-model="routeForm.routeNo" 
              placeholder="数字英文字母组合，最多10个字符"
              maxlength="10"
            />
            <div class="form-tip">数字英文字母组合，最多10个字符；建议格式：不可修改</div>
          </el-form-item>

          <el-form-item label="车队：" prop="fieet">
            <el-select v-model="routeForm.fieet" placeholder="请选择车队">
              <el-option label="车队A" value="fleetA"></el-option>
              <el-option label="车队B" value="fleetB"></el-option>
              <el-option label="车队C" value="fleetC"></el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="标签：" prop="tage">
            <el-input 
              v-model="routeForm.tage" 
              placeholder="请输入标签"
            />
            <div class="form-tip">蓝色表示新增标签显示（最多添加3个），灰色表示已有标签不显示</div>
          </el-form-item>

          <el-form-item label="起点站：" prop="startStationId">
            <el-select v-model="routeForm.startStationId" placeholder="请选择起点站" filterable>
              <el-option 
                v-for="station in stationList" 
                :key="station.id" 
                :label="station.name" 
                :value="station.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="终点站：" prop="endStationId">
            <el-select v-model="routeForm.endStationId" placeholder="请选择终点站" filterable>
              <el-option 
                v-for="station in stationList" 
                :key="station.id" 
                :label="station.name" 
                :value="station.id"
              />
            </el-select>
          </el-form-item>
        </el-form>
      </div>
    </div>

    <!-- 第二步：配置站点 -->
    <div v-if="currentStep === 1" class="step-content">
      <div class="station-config">
        <!-- 上车站点 -->
        <div class="station-section">
          <div class="section-header">
            <el-icon><Plus /></el-icon>
            <span>添加上车点</span>
          </div>
          
          <div class="station-list">
            <div 
              v-for="(station, index) in upStations" 
              :key="station.stationId"
              class="station-item"
            >
              <div class="station-info">
                <div class="station-name">{{ getStationName(station.stationId) }}</div>
                <div class="station-time">行驶时间: {{ station.travelTime }} 分钟</div>
              </div>
              <div class="station-actions">
                <el-button type="primary" size="small" @click="editUpStation(index)">编辑</el-button>
                <el-button type="danger" size="small" @click="removeUpStation(index)">删除</el-button>
              </div>
            </div>
            
            <div class="add-station-btn">
              <el-button type="dashed" @click="addUpStation" style="width: 100%;">
                <el-icon><Plus /></el-icon>
                添加上车站点
              </el-button>
            </div>
          </div>
        </div>

        <!-- 下车站点 -->
        <div class="station-section">
          <div class="section-header">
            <el-icon><Minus /></el-icon>
            <span>添加下车点</span>
          </div>
          
          <div class="station-list">
            <div 
              v-for="(station, index) in downStations" 
              :key="station.stationId"
              class="station-item"
            >
              <div class="station-info">
                <div class="station-name">{{ getStationName(station.stationId) }}</div>
                <div class="station-time">行驶时间: {{ station.travelTime }} 分钟</div>
              </div>
              <div class="station-actions">
                <el-button type="primary" size="small" @click="editDownStation(index)">编辑</el-button>
                <el-button type="danger" size="small" @click="removeDownStation(index)">删除</el-button>
              </div>
            </div>
            
            <div class="add-station-btn">
              <el-button type="dashed" @click="addDownStation" style="width: 100%;">
                <el-icon><Plus /></el-icon>
                添加下车站点
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 第三步：完成 -->
    <div v-if="currentStep === 2" class="step-content">
      <div class="success-content">
        <el-result icon="success" title="班线创建成功" sub-title="班线信息已保存，可以开始配置班次了">
          <template #extra>
            <el-button type="primary" @click="goToList">返回列表</el-button>
            <el-button @click="createAnother">再创建一个</el-button>
          </template>
        </el-result>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="action-buttons" v-if="currentStep < 2">
      <el-button v-if="currentStep > 0" @click="prevStep">上一步</el-button>
      <el-button 
        type="primary" 
        @click="nextStep"
        :loading="submitting"
      >
        {{ currentStep === 1 ? '创建班线' : '下一步' }}
      </el-button>
    </div>

    <!-- 站点选择弹窗 -->
    <el-dialog v-model="stationDialogVisible" title="选择站点" width="600px">
      <el-form :model="stationForm" label-width="100px">
        <el-form-item label="选择站点：">
          <el-select v-model="stationForm.stationId" placeholder="请选择站点" filterable style="width: 100%">
            <el-option 
              v-for="station in availableStations" 
              :key="station.id" 
              :label="station.name" 
              :value="station.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="行驶时间：">
          <el-input-number 
            v-model="stationForm.travelTime" 
            :min="0" 
            :max="999"
            placeholder="分钟"
            style="width: 100%"
          />
          <div class="form-tip">从起点到此站点的行驶时间（分钟）</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="stationDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAddStation">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Minus } from '@element-plus/icons-vue'
import { routeApi } from '@/api'
import { useRouter } from 'vue-router'

export default {
  name: 'RouteCreate',
  components: {
    Plus,
    Minus
  },
  setup() {
    const router = useRouter()
    const currentStep = ref(0)
    const submitting = ref(false)
    const routeFormRef = ref(null)
    const stationDialogVisible = ref(false)
    const currentStationType = ref('') // 'up' or 'down'
    const editingIndex = ref(-1)

    const routeForm = ref({
      routeNo: '',
      fieet: '',
      tage: '',
      startStationId: null,
      endStationId: null
    })

    const rules = {
      routeNo: [
        { required: true, message: '请输入班线编号', trigger: 'blur' },
        { max: 10, message: '班线编号不能超过10个字符', trigger: 'blur' }
      ],
      fieet: [
        { required: true, message: '请选择车队', trigger: 'change' }
      ],
      startStationId: [
        { required: true, message: '请选择起点站', trigger: 'change' }
      ],
      endStationId: [
        { required: true, message: '请选择终点站', trigger: 'change' }
      ]
    }

    const upStations = ref([])
    const downStations = ref([])
    const stationList = ref([])

    const stationForm = ref({
      stationId: null,
      travelTime: 0
    })

    // 可选择的站点（排除已选择的）
    const availableStations = computed(() => {
      const selectedIds = [
        ...upStations.value.map(s => s.stationId),
        ...downStations.value.map(s => s.stationId)
      ]
      return stationList.value.filter(station => !selectedIds.includes(station.id))
    })

    const getStationName = (stationId) => {
      const station = stationList.value.find(s => s.id === stationId)
      return station ? station.name : '未知站点'
    }

    const loadStations = async () => {
      // 这里应该调用获取站点列表的API
      // 暂时使用模拟数据
      stationList.value = [
        { id: 1, name: '昆明长水机场' },
        { id: 2, name: '昆明火车站' },
        { id: 3, name: '昆明汽车站' },
        { id: 4, name: '西双版纳客运站' },
        { id: 5, name: '大理古城' }
      ]
    }

    const nextStep = async () => {
      if (currentStep.value === 0) {
        // 验证第一步表单
        const valid = await routeFormRef.value.validate().catch(() => false)
        if (!valid) return
        
        currentStep.value++
      } else if (currentStep.value === 1) {
        // 创建班线
        await createRoute()
      }
    }

    const prevStep = () => {
      if (currentStep.value > 0) {
        currentStep.value--
      }
    }

    const createRoute = async () => {
      if (upStations.value.length === 0) {
        ElMessage.warning('请至少添加一个上车站点')
        return
      }

      submitting.value = true
      try {
        const data = {
          ...routeForm.value,
          upStations: upStations.value,
          downStations: downStations.value
        }

        const response = await routeApi.create(data)
        
        if (response.data.success) {
          ElMessage.success('班线创建成功')
          currentStep.value = 2
        } else {
          ElMessage.error(response.data.msg || '创建失败')
        }
      } catch (error) {
        ElMessage.error('创建失败')
      } finally {
        submitting.value = false
      }
    }

    const addUpStation = () => {
      currentStationType.value = 'up'
      editingIndex.value = -1
      stationForm.value = { stationId: null, travelTime: 0 }
      stationDialogVisible.value = true
    }

    const addDownStation = () => {
      currentStationType.value = 'down'
      editingIndex.value = -1
      stationForm.value = { stationId: null, travelTime: 0 }
      stationDialogVisible.value = true
    }

    const editUpStation = (index) => {
      currentStationType.value = 'up'
      editingIndex.value = index
      stationForm.value = { ...upStations.value[index] }
      stationDialogVisible.value = true
    }

    const editDownStation = (index) => {
      currentStationType.value = 'down'
      editingIndex.value = index
      stationForm.value = { ...downStations.value[index] }
      stationDialogVisible.value = true
    }

    const confirmAddStation = () => {
      if (!stationForm.value.stationId) {
        ElMessage.warning('请选择站点')
        return
      }

      const stationData = {
        stationId: stationForm.value.stationId,
        travelTime: stationForm.value.travelTime || 0
      }

      if (currentStationType.value === 'up') {
        if (editingIndex.value >= 0) {
          upStations.value[editingIndex.value] = stationData
        } else {
          upStations.value.push(stationData)
        }
        // 按时间排序
        upStations.value.sort((a, b) => a.travelTime - b.travelTime)
        // 第一个站点时间设为0
        if (upStations.value.length > 0) {
          upStations.value[0].travelTime = 0
        }
      } else {
        if (editingIndex.value >= 0) {
          downStations.value[editingIndex.value] = stationData
        } else {
          downStations.value.push(stationData)
        }
        // 按时间排序
        downStations.value.sort((a, b) => a.travelTime - b.travelTime)
      }

      stationDialogVisible.value = false
    }

    const removeUpStation = (index) => {
      upStations.value.splice(index, 1)
    }

    const removeDownStation = (index) => {
      downStations.value.splice(index, 1)
    }

    const goToList = () => {
      router.push('/route/list')
    }

    const createAnother = () => {
      // 重置表单
      currentStep.value = 0
      routeForm.value = {
        routeNo: '',
        fieet: '',
        tage: '',
        startStationId: null,
        endStationId: null
      }
      upStations.value = []
      downStations.value = []
    }

    onMounted(() => {
      loadStations()
    })

    return {
      currentStep,
      submitting,
      routeFormRef,
      stationDialogVisible,
      currentStationType,
      editingIndex,
      routeForm,
      rules,
      upStations,
      downStations,
      stationList,
      stationForm,
      availableStations,
      getStationName,
      nextStep,
      prevStep,
      addUpStation,
      addDownStation,
      editUpStation,
      editDownStation,
      confirmAddStation,
      removeUpStation,
      removeDownStation,
      goToList,
      createAnother
    }
  }
}
</script>

<style scoped>
.route-create {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.steps-container {
  margin-bottom: 40px;
}

.step-content {
  min-height: 400px;
  margin-bottom: 30px;
}

.form-container h3 {
  margin-bottom: 30px;
  color: #333;
  font-size: 18px;
}

.route-form {
  max-width: 600px;
}

.form-tip {
  font-size: 12px;
  color: #666;
  margin-top: 5px;
}

.station-config {
  display: flex;
  gap: 30px;
}

.station-section {
  flex: 1;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: linear-gradient(90deg, #409EFF 0%, #66B3FF 100%);
  color: white;
  border-radius: 4px;
  margin-bottom: 15px;
  font-weight: 500;
}

.station-section:nth-child(2) .section-header {
  background: linear-gradient(90deg, #FF8C00 0%, #FFA500 100%);
}

.station-list {
  border: 1px solid #e6e6e6;
  border-radius: 4px;
  min-height: 200px;
  padding: 15px;
}

.station-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border: 1px solid #f0f0f0;
  border-radius: 4px;
  margin-bottom: 10px;
  background: #fafafa;
}

.station-info {
  flex: 1;
}

.station-name {
  font-weight: 500;
  color: #333;
  margin-bottom: 5px;
}

.station-time {
  font-size: 12px;
  color: #666;
}

.station-actions {
  display: flex;
  gap: 8px;
}

.add-station-btn {
  margin-top: 10px;
}

.success-content {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 40px;
}

.action-buttons .el-button {
  min-width: 100px;
}
</style>