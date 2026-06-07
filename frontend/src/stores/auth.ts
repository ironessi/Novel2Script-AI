import { defineStore } from 'pinia'
import { authApi } from '@/api'

interface User {
  id: number
  username: string
  email: string
  role: string
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    user: JSON.parse(localStorage.getItem('user') || 'null') as User | null
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    username: (state) => state.user?.username || ''
  },
  actions: {
    async login(username: string, password: string) {
      const res: any = await authApi.login({ username, password })
      if (res.code === 0) {
        this.token = res.data.token
        this.user = res.data.user
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('user', JSON.stringify(res.data.user))
        return true
      }
      throw new Error(res.message || '登录失败')
    },
    async register(username: string, email: string, password: string) {
      const res: any = await authApi.register({ username, email, password })
      if (res.code === 0) return true
      throw new Error(res.message || '注册失败')
    },
    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  }
})
