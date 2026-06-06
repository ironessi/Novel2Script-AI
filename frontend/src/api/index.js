import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 300000,
})

// 请求拦截器：自动添加 Token
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器：处理 401
api.interceptors.response.use(
  res => res.data,
  err => {
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
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
  me: () => api.get('/auth/me'),
}

// ========== Projects ==========
export const projectApi = {
  list: (params) => api.get('/projects', { params }),
  create: (data) => api.post('/projects', data),
  get: (id) => api.get(`/projects/${id}`),
  update: (id, data) => api.put(`/projects/${id}`, data),
  delete: (id) => api.delete(`/projects/${id}`),
}

// ========== Upload ==========
export const uploadApi = {
  upload: (projectId, formData) => api.post(`/projects/${projectId}/upload`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    timeout: 600000,
  }),
}

// ========== Chapters ==========
export const chapterApi = {
  list: (projectId) => api.get(`/projects/${projectId}/chapters`),
  get: (projectId, chapterId) => api.get(`/projects/${projectId}/chapters/${chapterId}`),
}

// ========== Tasks ==========
export const taskApi = {
  create: (projectId, data) => api.post(`/projects/${projectId}/generate`, data),
  status: (taskId) => api.get(`/tasks/${taskId}/status`),
}

// ========== Script ==========
export const scriptApi = {
  get: (projectId) => api.get(`/projects/${projectId}/script`),
  update: (projectId, data) => api.put(`/projects/${projectId}/script`, data),
  validate: (projectId) => api.post(`/projects/${projectId}/validate`),
  checkHallucination: (projectId) => api.post(`/projects/${projectId}/check-hallucination`),
  checkSafety: (projectId) => api.post(`/projects/${projectId}/check-safety`),
  exportUrl: (projectId, format = 'yaml') => `/api/projects/${projectId}/export?format=${format}`,
}

// ========== Audit ==========
export const auditApi = {
  list: (projectId, params) => api.get(`/projects/${projectId}/audit`, { params }),
}
