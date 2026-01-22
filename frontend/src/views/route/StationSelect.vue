<template>
  <div class="station-select">
    <!-- 选择上车点弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      title="选择上车点"
      width="80%"
      :before-close="handleClose"
    >
      <div class="dialog-content">
        <!-- 左侧站点列表 -->
        <div class="left-panel">
          <div class="search-bar">
            <el-input
              v-model="searchKeyword"
              placeholder="请输入关键词"
              class="search-input"
            >
              <template #append>
                <el-button>搜索</el-button>
              </template>
            </el-input>
            <el-button type="success" class="add-btn">
              <el-icon><Plus /></el-icon>
              添加站点
            </el-button>
          </div>

          <div class="station-list-header">
            <h3>站点列表</h3>
          </div>

          <div class="station-list">
            <div 
              v-for="station in filteredStations" 
              :key="station.id"
              class="station-item"
              @click="selectStation(station)"
            >
              <div class="station-info">
                <div class="station-name">{{ station.name }}</div>
                <div class="station-detail">{{ station.address }}</div>
              </div>
              <div class="station-actions">
                <el-button 
                  v-if="station.canSelect"
                  type="primary" 
                  size="small"
                  @click.stop="addToSelected(station)"
                >
                  选择
                </el-button>
                <el-button 
                  v-else
                  type="info" 
                  size="small"
                  disabled
                >
                  已选
                </el-button>
              </div>
            </div>
          </div>

          <!-- 分页 -->
          <div class="pagination">
            <el-pagination
              v-model:current-page="currentPage"
              :page-size="pageSize"
              :total="totalStations"
              layout="prev, pager, next"
              @current-change="handlePageChange"
            />
          </div>
        </div>

        <!-- 右侧添加下车点 -->
        <div class="right-panel">
          <div class="add-station-header">
            <h3>添加下车点</h3>
          </div>

          <div class="station-form">
            <div class="form-item">
              <label>站点名称：</label>
              <span class="station-name-display">{{ selectedStation?.name || '请先选择上车站点' }}</span>
            </div>

            <div class="form-item">
              <label>预计发车时间：</label>
              <span class="time-display">预计时间(分钟)</span>
            </div>

            <div class="save-actions">
              <el-button type="primary" class="save-btn">保存</el-button>
              <el-button class="cancel-btn">取消</el-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 新增站点区域 -->
      <div class="new-station-section">
        <div class="section-header">
          <h3>新增站点</h3>
        </div>

        <div class="station-form-container">
          <div class="form-row">
            <div class="form-group">
              <label>站点名称：</label>
              <el-input v-model="newStation.name" placeholder="请输入站点名称" />
            </div>
            <div class="form-group">
              <label>站点编号：</label>
              <el-input v-model="newStation.code" placeholder="请输入站点编号" />
            </div>
          </div>

          <div class="location-section">
            <div class="location-header">
              <label>选定站点位置</label>
              <span class="location-tip">请在地图上选择上车站点位置，或者在搜索框中搜索</span>
            </div>

            <!-- 地图区域 -->
            <div class="map-container">
              <div class="map-placeholder">
                <!-- 这里应该集成真实的地图组件，如高德地图、百度地图等 -->
                <div class="map-content">
                  <p>地图加载中...</p>
                  <p>请在此处集成地图组件</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 触发按钮 -->
    <el-button type="primary" @click="openDialog">选择上车点</el-button>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'

export default {
  name: 'StationSelect',
  components: {
    Plus
  },
  setup() {
    const dialogVisible = ref(false)
    const searchKeyword = ref('')
    const currentPage = ref(1)
    const pageSize = ref(10)
    const selectedStation = ref(null)
    const newStation = ref({
      name: '',
      code: ''
    })

    // 模拟站点数据
    const stations = ref([
      {
        id: 1,
        name: '昆明-上水机场出发层',
        address: '昆明长水国际机场出发层',
        canSelect: false
      },
      {
        id: 2,
        name: '昆明-上水机场出发层(3层)第7通外车道',
        address: '昆明长水国际机场出发层3层第7通外车道',
        canSelect: true
      },
      {
        id: 3,
        name: '昆明-上水机场出发层(3层)第7通外车道',
        address: '昆明长水国际机场出发层3层第7通外车道',
        canSelect: true
      },
      {
        id: 4,
        name: '昆明-上水机场出发层(3层)第7通外车道',
        address: '昆明长水国际机场出发层3层第7通外车道',
        canSelect: true
      },
      {
        id: 5,
        name: '昆明-上水机场出发层(3层)第7通外车道',
        address: '昆明长水国际机场出发层3层第7通外车道',
        canSelect: true
      }
    ])

    const totalStations = computed(() => stations.value.length)

    const filteredStations = computed(() => {
      let filtered = stations.value
      if (searchKeyword.value) {
        filtered = filtered.filter(station => 
          station.name.includes(searchKeyword.value) || 
          station.address.includes(searchKeyword.value)
        )
      }
      
      const start = (currentPage.value - 1) * pageSize.value
      const end = start + pageSize.value
      return filtered.slice(start, end)
    })

    const openDialog = () => {
      dialogVisible.value = true
    }

    const handleClose = () => {
      dialogVisible.value = false
    }

    const selectStation = (station) => {
      selectedStation.value = station
    }

    const addToSelected = (station) => {
      station.canSelect = false
      selectedStation.value = station
    }

    const handlePageChange = (page) => {
      currentPage.value = page
    }

    return {
      dialogVisible,
      searchKeyword,
      currentPage,
      pageSize,
      selectedStation,
      newStation,
      stations,
      totalStations,
      filteredStations,
      openDialog,
      handleClose,
      selectStation,
      addToSelected,
      handlePageChange
    }
  }
}
</script>

<style scoped>
.station-select {
  padding: 20px;
}

.dialog-content {
  display: flex;
  gap: 20px;
  min-height: 400px;
}

.left-panel {
  flex: 1;
  border-right: 1px solid #e6e6e6;
  padding-right: 20px;
}

.right-panel {
  width: 300px;
}

.search-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.search-input {
  flex: 1;
}

.add-btn {
  white-space: nowrap;
}

.station-list-header {
  margin-bottom: 15px;
}

.station-list-header h3 {
  font-size: 16px;
  color: #333;
}

.station-list {
  max-height: 300px;
  overflow-y: auto;
}

.station-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border: 1px solid #e6e6e6;
  border-radius: 4px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.station-item:hover {
  background-color: #f5f7fa;
  border-color: #409EFF;
}

.station-info {
  flex: 1;
}

.station-name {
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.station-detail {
  font-size: 12px;
  color: #666;
}

.station-actions {
  margin-left: 10px;
}

.pagination {
  margin-top: 20px;
  text-align: center;
}

.add-station-header {
  margin-bottom: 20px;
}

.add-station-header h3 {
  font-size: 16px;
  color: #333;
}

.station-form {
  background: #f9f9f9;
  padding: 20px;
  border-radius: 4px;
}

.form-item {
  margin-bottom: 15px;
}

.form-item label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #333;
}

.station-name-display {
  color: #666;
  font-size: 14px;
}

.time-display {
  color: #666;
  font-size: 14px;
}

.save-actions {
  margin-top: 20px;
}

.save-btn {
  width: 100%;
  margin-bottom: 10px;
}

.cancel-btn {
  width: 100%;
}

.new-station-section {
  margin-top: 30px;
  border-top: 1px solid #e6e6e6;
  padding-top: 20px;
}

.section-header {
  margin-bottom: 20px;
}

.section-header h3 {
  font-size: 16px;
  color: #333;
}

.station-form-container {
  background: #f9f9f9;
  padding: 20px;
  border-radius: 4px;
}

.form-row {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
}

.form-group {
  flex: 1;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #333;
}

.location-section {
  margin-top: 20px;
}

.location-header {
  margin-bottom: 15px;
}

.location-header label {
  display: block;
  font-weight: 500;
  color: #333;
  margin-bottom: 5px;
}

.location-tip {
  font-size: 12px;
  color: #666;
}

.map-container {
  height: 300px;
  border: 1px solid #e6e6e6;
  border-radius: 4px;
  overflow: hidden;
}

.map-placeholder {
  width: 100%;
  height: 100%;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.map-content {
  text-align: center;
  color: #666;
}

.map-content p {
  margin: 5px 0;
}
</style>