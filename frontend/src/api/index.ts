import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 300000
})

// 请求拦截：自动添加 Token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截：处理 401
api.interceptors.response.use(
  (res) => res.data,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(err.response?.data || err)
  }
)

// ========== Auth ==========
export const authApi = {
  register: (data: { username: string; email?: string; password: string }) =>
    api.post('/auth/register', data),
  login: (data: { username: string; password: string }) =>
    api.post('/auth/login', data),
  me: () => api.get('/auth/me')
}

// ========== Projects ==========
export const projectApi = {
  list: (params?: { page?: number; page_size?: number }) =>
    api.get('/projects', { params }),
  create: (data: { title: string; description?: string; adaptation_mode?: string }) =>
    api.post('/projects', data),
  get: (id: number) => api.get(`/projects/${id}`),
  update: (id: number, data: Record<string, any>) =>
    api.put(`/projects/${id}`, data),
  delete: (id: number) => api.delete(`/projects/${id}`)
}

// ========== Upload ==========
export const uploadApi = {
  upload: (projectId: number, formData: FormData) =>
    api.post(`/projects/${projectId}/upload`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      timeout: 600000
    })
}

// ========== Chapters ==========
export const chapterApi = {
  list: (projectId: number) => api.get(`/projects/${projectId}/chapters`),
  get: (projectId: number, chapterId: number) =>
    api.get(`/projects/${projectId}/chapters/${chapterId}`)
}

// ========== Tasks ==========
export const taskApi = {
  create: (projectId: number, data: { task_type: string }) =>
    api.post(`/projects/${projectId}/generate`, data),
  status: (taskId: number) => api.get(`/tasks/${taskId}/status`)
}

// ========== Script ==========
export const scriptApi = {
  get: (projectId: number) => api.get(`/projects/${projectId}/script`),
  update: (projectId: number, data: { yaml_content: string }) =>
    api.put(`/projects/${projectId}/script`, data),
  validate: (projectId: number) =>
    api.post(`/projects/${projectId}/validate`),
  checkHallucination: (projectId: number) =>
    api.post(`/projects/${projectId}/check-hallucination`),
  checkSafety: (projectId: number) =>
    api.post(`/projects/${projectId}/check-safety`),
  exportUrl: (projectId: number, format = 'yaml') =>
    `/api/projects/${projectId}/export?format=${format}`
}

// ========== Audit ==========
export const auditApi = {
  list: (projectId: number, params?: { page?: number; page_size?: number }) =>
    api.get(`/projects/${projectId}/audit`, { params })
}
