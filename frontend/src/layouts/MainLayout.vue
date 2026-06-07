<template>
  <div class="layout">
    <header class="header">
      <div class="header-left">
        <router-link to="/" class="logo">Novel2Script-AI</router-link>
      </div>
      <div class="header-right">
        <span class="username">{{ user?.username }}</span>
        <button class="btn btn-outline" @click="logout">退出</button>
      </div>
    </header>
    <main class="main">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const user = computed(() => {
  try { return JSON.parse(localStorage.getItem('user')) } catch { return null }
})

function logout() {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}
</script>

<style scoped>
.layout { min-height: 100vh; }
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 24px;
  background: white;
  border-bottom: 1px solid var(--gray-200);
  position: sticky;
  top: 0;
  z-index: 10;
}
.logo { font-size: 18px; font-weight: 700; color: var(--primary); }
.logo:hover { text-decoration: none; }
.header-right { display: flex; align-items: center; gap: 12px; }
.username { font-size: 14px; color: var(--gray-500); }
.main { max-width: 1200px; margin: 0 auto; padding: 24px; }
</style>
