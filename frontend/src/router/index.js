import { createRouter, createWebHistory } from 'vue-router'
import Layout from '../components/Layout.vue'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    component: Layout,
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue')
      }
    ]
  },
  {
    path: '/route',
    component: Layout,
    children: [
      {
        path: 'list',
        name: 'RouteList',
        component: () => import('../views/route/RouteList.vue')
      },
      {
        path: 'create',
        name: 'RouteCreate',
        component: () => import('../views/route/RouteCreate.vue')
      },
      {
        path: 'edit/:id',
        name: 'RouteEdit',
        component: () => import('../views/route/RouteCreate.vue')
      },
      {
        path: 'schedule',
        name: 'ScheduleManage',
        component: () => import('../views/route/ScheduleManage.vue')
      },
      {
        path: 'station-select',
        name: 'StationSelect',
        component: () => import('../views/route/StationSelect.vue')
      }
    ]
  },
  {
    path: '/finance',
    component: Layout,
    children: [
      {
        path: 'balance',
        name: 'BalanceSheet',
        component: () => import('../views/finance/BalanceSheet.vue')
      },
      {
        path: 'income',
        name: 'IncomeSheet',
        component: () => import('../views/finance/IncomeSheet.vue')
      },
      {
        path: 'route-settle',
        name: 'RouteSettle',
        component: () => import('../views/finance/RouteSettle.vue')
      },
      {
        path: 'schedule-settle',
        name: 'ScheduleSettle',
        component: () => import('../views/finance/ScheduleSettle.vue')
      },
      {
        path: 'station-settle',
        name: 'StationSettle',
        component: () => import('../views/finance/StationSettle.vue')
      },
      {
        path: 'driver-settle',
        name: 'DriverSettle',
        component: () => import('../views/finance/DriverSettle.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router