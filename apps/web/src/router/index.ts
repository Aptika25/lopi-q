import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/public/HomeView.vue'),
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/public/LoginView.vue'),
    meta: { guestOnly: true, hideNavbarFooter: true },
  },
  {
    path: '/intern/dashboard',
    name: 'InternDashboard',
    component: () => import('@/views/intern/DashboardView.vue'),
    meta: { requiresAuth: true, hideNavbarFooter: true },
  },
  {
    path: '/intern/history',
    name: 'InternHistory',
    component: () => import('@/views/intern/HistoryView.vue'),
    meta: { requiresAuth: true, hideNavbarFooter: true },
  },
  {
    path: '/presensi',
    name: 'InternPresensi',
    component: () => import('@/views/intern/ScanView.vue'),
    meta: { requiresAuth: true, hideNavbarFooter: true },
  },
  {
    path: '/intern/scan',
    name: 'InternScan',
    component: () => import('@/views/intern/ScanView.vue'),
    meta: { requiresAuth: true, hideNavbarFooter: true },
  },
  {
    path: '/admin',
    name: 'AdminDashboard',
    component: () => import('@/views/super-admin/DashboardView.vue'),
    meta: { requiresAuth: true, requiresAdmin: true, hideNavbarFooter: true },
  },
  {
    path: '/admin/admins',
    name: 'AdminAdmins',
    component: () => import('@/views/super-admin/AdminsView.vue'),
    meta: { requiresAuth: true, requiresAdmin: true, hideNavbarFooter: true },
  },
  {
    path: '/admin/interns',
    name: 'AdminInterns',
    component: () => import('@/views/super-admin/InternsView.vue'),
    meta: { requiresAuth: true, requiresAdmin: true, hideNavbarFooter: true },
  },
  {
    path: '/admin/attendance-recap',
    name: 'AdminAttendanceRecap',
    component: () => import('@/views/super-admin/AttendanceRecapView.vue'),
    meta: { requiresAuth: true, requiresAdmin: true, hideNavbarFooter: true },
  },
  {
    path: '/admin/location',
    name: 'AdminLocation',
    component: () => import('@/views/super-admin/LocationView.vue'),
    meta: { requiresAuth: true, requiresAdmin: true, hideNavbarFooter: true },
  },
  {
    path: '/admin/activity-logs',
    name: 'AdminActivityLogs',
    component: () => import('@/views/super-admin/ActivityLogsView.vue'),
    meta: { requiresAuth: true, requiresAdmin: true, hideNavbarFooter: true },
  },
  {
    path: '/admin/security',
    name: 'AdminSecurity',
    component: () => import('@/views/super-admin/SecuritySettingsView.vue'),
    meta: { requiresAuth: true, requiresAdmin: true, hideNavbarFooter: true },
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) return savedPosition
    return { top: 0, behavior: 'instant' }
  }
})

router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()

  if (authStore.token && !authStore.user) {
    await authStore.fetchProfile()
  }

  const isLoggedIn = authStore.isAuthenticated

  if (to.matched.some((record) => record.meta.requiresAuth)) {
    if (!isLoggedIn) {
      next({ name: 'Login' })
    } else if (to.matched.some((record) => record.meta.requiresAdmin) && !authStore.isAdmin) {
      next({ name: 'Home' })
    } else {
      next()
    }
  } else if (to.matched.some((record) => record.meta.guestOnly)) {
    if (isLoggedIn) {
      if (authStore.isAdmin) {
        next({ name: 'AdminDashboard' })
      } else {
        next({ name: 'InternDashboard' })
      }
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
