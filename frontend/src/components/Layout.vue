<template>
  <div class="layout">
    <!-- 顶部导航 -->
    <el-header class="header">
      <div class="header-left">
        <h2>拼约出行管理平台</h2>
      </div>
      <div class="header-right">
        <el-dropdown>
          <span class="user-info">
            <el-icon><User /></el-icon>
            admin
            <el-icon class="el-icon--right"><arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="logout">退出</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>

    <el-container>
      <!-- 侧边栏 -->
      <el-aside class="sidebar" width="200px">
        <el-menu
          :default-active="$route.path"
          router
          background-color="#f5f5f5"
          text-color="#333"
          active-text-color="#409EFF"
        >
          <el-menu-item index="/dashboard">
            <el-icon><House /></el-icon>
            <span>首页管理</span>
          </el-menu-item>
          
          <el-sub-menu index="route">
            <template #title>
              <el-icon><Location /></el-icon>
              <span>班车管理</span>
            </template>
            <el-menu-item index="/route/list">班线管理</el-menu-item>
            <el-menu-item index="/route/create">新增班线</el-menu-item>
            <el-menu-item index="/route/schedule">班次管理</el-menu-item>
          </el-sub-menu>

          <el-menu-item index="/vehicle">
            <el-icon><Van /></el-icon>
            <span>车辆管理</span>
          </el-menu-item>

          <el-menu-item index="/driver">
            <el-icon><User /></el-icon>
            <span>人员管理</span>
          </el-menu-item>

          <el-sub-menu index="finance">
            <template #title>
              <el-icon><Money /></el-icon>
              <span>数据分析</span>
            </template>
            <el-menu-item index="/finance/balance">收支对账表</el-menu-item>
            <el-menu-item index="/finance/income">收入对账表</el-menu-item>
            <el-menu-item index="/finance/route-settle">线路结算表</el-menu-item>
            <el-menu-item index="/finance/schedule-settle">班次结算表</el-menu-item>
            <el-menu-item index="/finance/station-settle">站点结算表</el-menu-item>
            <el-menu-item index="/finance/driver-settle">司机结算表</el-menu-item>
          </el-sub-menu>

          <el-menu-item index="/order">
            <el-icon><Document /></el-icon>
            <span>订单管理</span>
          </el-menu-item>

          <el-menu-item index="/finance-manage">
            <el-icon><CreditCard /></el-icon>
            <span>财务管理</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <!-- 主内容区 -->
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<script>
import { User, ArrowDown, House, Location, Van, Money, Document, CreditCard } from '@element-plus/icons-vue'

export default {
  name: 'Layout',
  components: {
    User,
    ArrowDown,
    House,
    Location,
    Van,
    Money,
    Document,
    CreditCard
  },
  methods: {
    logout() {
      this.$router.push('/login')
    }
  }
}
</script>

<style scoped>
.layout {
  height: 100vh;
}

.header {
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}

.header-left h2 {
  color: #333;
  font-size: 18px;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  color: #666;
}

.user-info .el-icon {
  margin-right: 5px;
}

.sidebar {
  background: #f5f5f5;
  border-right: 1px solid #e6e6e6;
}

.main-content {
  background: #fff;
  padding: 20px;
}
</style>