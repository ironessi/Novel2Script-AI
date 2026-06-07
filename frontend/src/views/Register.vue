<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <h1 class="login-title">Novel2Script</h1>
        <p class="login-slogan">创建账号，开始你的改编之旅。</p>
      </div>
      <n-form ref="formRef" :model="form" :rules="rules" label-placement="left">
        <n-form-item label="用户名" path="username">
          <n-input v-model:value="form.username" placeholder="3-64位字符" />
        </n-form-item>
        <n-form-item label="邮箱" path="email">
          <n-input v-model:value="form.email" placeholder="选填" />
        </n-form-item>
        <n-form-item label="密码" path="password">
          <n-input v-model:value="form.password" type="password" placeholder="至少8位" show-password-on="click" />
        </n-form-item>
        <n-button type="primary" block :loading="loading" @click="handleRegister">
          注册
        </n-button>
      </n-form>
      <div class="login-footer">
        已有账号？<router-link to="/login">登录</router-link>
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
const form = ref({ username: '', email: '', password: '' })

const rules: FormRules = {
  username: { required: true, message: '请输入用户名', trigger: 'blur' },
  password: { required: true, min: 8, message: '密码至少8位', trigger: 'blur' }
}

async function handleRegister() {
  try { await formRef.value?.validate() } catch { return }
  loading.value = true
  error.value = ''
  try {
    await auth.register(form.value.username, form.value.email, form.value.password)
    router.push('/login')
  } catch (e: any) {
    error.value = e.message || '注册失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page { height: 100vh; display: flex; align-items: center; justify-content: center; background: var(--color-bg); }
.login-card { width: 380px; padding: 40px 32px; background: var(--color-surface); border-radius: 12px; border: 1px solid var(--color-border); }
.login-header { text-align: center; margin-bottom: 32px; }
.login-title { font-size: 24px; font-weight: 700; color: var(--color-primary); margin-bottom: 8px; }
.login-slogan { font-size: 14px; color: var(--color-text-secondary); }
.login-footer { text-align: center; margin-top: 16px; font-size: 13px; color: var(--color-text-secondary); }
.login-footer a { color: var(--color-primary); text-decoration: none; }
</style>
