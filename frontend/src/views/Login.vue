<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <h1 class="login-title">Novel2Script</h1>
        <p class="login-slogan">把小说拆成场景，把故事变成剧本。</p>
      </div>
      <n-form ref="formRef" :model="form" :rules="rules" label-placement="left">
        <n-form-item label="用户名" path="username">
          <n-input v-model:value="form.username" placeholder="请输入用户名" @keyup.enter="handleLogin" />
        </n-form-item>
        <n-form-item label="密码" path="password">
          <n-input v-model:value="form.password" type="password" placeholder="请输入密码" show-password-on="click" @keyup.enter="handleLogin" />
        </n-form-item>
        <n-button type="primary" block :loading="loading" @click="handleLogin">
          登录
        </n-button>
      </n-form>
      <div class="login-footer">
        还没有账号？<router-link to="/register">注册</router-link>
      </div>
      <n-alert v-if="error" type="error" style="margin-top: 12px;" closable @close="error = ''">
        {{ error }}
      </n-alert>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import type { FormInst, FormRules } from 'naive-ui'

const router = useRouter()
const auth = useAuthStore()
const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const error = ref('')

const form = ref({ username: '', password: '' })
const rules: FormRules = {
  username: { required: true, message: '请输入用户名', trigger: 'blur' },
  password: { required: true, message: '请输入密码', trigger: 'blur' }
}

async function handleLogin() {
  try {
    await formRef.value?.validate()
  } catch { return }

  loading.value = true
  error.value = ''
  try {
    await auth.login(form.value.username, form.value.password)
    router.push('/dashboard')
  } catch (e: any) {
    error.value = e.message || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg);
}
.login-card {
  width: 380px;
  padding: 40px 32px;
  background: var(--color-surface);
  border-radius: 12px;
  border: 1px solid var(--color-border);
}
.login-header {
  text-align: center;
  margin-bottom: 32px;
}
.login-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-primary);
  margin-bottom: 8px;
}
.login-slogan {
  font-size: 14px;
  color: var(--color-text-secondary);
}
.login-footer {
  text-align: center;
  margin-top: 16px;
  font-size: 13px;
  color: var(--color-text-secondary);
}
.login-footer a {
  color: var(--color-primary);
  text-decoration: none;
}
</style>
