import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/login', name: 'Login', component: () => import('../pages/Login.vue') },
  { path: '/register', name: 'Register', component: () => import('../pages/Register.vue') },
  {
    path: '/',
    component: () => import('../layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', name: 'Projects', component: () => import('../pages/ProjectList.vue') },
      { path: 'projects/:id', name: 'ProjectDetail', component: () => import('../pages/ProjectDetail.vue') },
      { path: 'projects/:id/upload', name: 'Upload', component: () => import('../pages/UploadNovel.vue') },
      { path: 'projects/:id/generate', name: 'Generate', component: () => import('../pages/GenerateTask.vue') },
      { path: 'projects/:id/script', name: 'Script', component: () => import('../pages/ScriptEditor.vue') },
      { path: 'projects/:id/audit', name: 'Audit', component: () => import('../pages/AuditLog.vue') },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
