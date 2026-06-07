import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/projects/:id/upload',
    name: 'Upload',
    component: () => import('@/views/Upload.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/projects/:id/workflow',
    name: 'Workflow',
    component: () => import('@/views/Workflow.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/projects/:id/workbench',
    name: 'Workbench',
    component: () => import('@/views/Workbench.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/projects/:id/characters',
    name: 'Characters',
    component: () => import('@/views/Characters.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/projects/:id/plot',
    name: 'Plot',
    component: () => import('@/views/Plot.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/projects/:id/scenes',
    name: 'Scenes',
    component: () => import('@/views/Scenes.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/projects/:id/yaml',
    name: 'YamlEditor',
    component: () => import('@/views/YamlEditor.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/projects/:id/validation',
    name: 'Validation',
    component: () => import('@/views/Validation.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/projects/:id/versions',
    name: 'Versions',
    component: () => import('@/views/Versions.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/projects/:id/audit',
    name: 'Audit',
    component: () => import('@/views/Audit.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth !== false && !auth.isLoggedIn) {
    return '/login'
  }
})

export default router
