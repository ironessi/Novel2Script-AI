<template>
  <div>
    <h2 class="mb-4">审计日志</h2>
    <div class="card">
      <div v-if="loading" class="text-center text-gray">加载中...</div>
      <table v-else-if="logs.length > 0">
        <thead>
          <tr><th>时间</th><th>操作</th><th>资源</th><th>IP</th><th>请求ID</th></tr>
        </thead>
        <tbody>
          <tr v-for="log in logs" :key="log.id">
            <td class="text-sm">{{ log.created_at }}</td>
            <td><span class="badge badge-pending">{{ log.action }}</span></td>
            <td class="text-sm">{{ log.resource_type }} #{{ log.resource_id }}</td>
            <td class="text-sm text-gray">{{ log.ip_address }}</td>
            <td class="text-sm text-gray">{{ log.request_id?.slice(0, 8) }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else class="text-center text-gray">暂无审计日志</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { auditApi } from '../api'

const route = useRoute()
const logs = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await auditApi.list(route.params.id, { page: 1, page_size: 50 })
    logs.value = res.data.logs || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})
</script>
