<template>
  <AppLayout>
    <div class="audit-page">
      <h2>审计日志</h2>
      <n-data-table :columns="columns" :data="store.auditLogs" :bordered="false" size="small" />
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import AppLayout from '@/layouts/AppLayout.vue'
import type { DataTableColumns } from 'naive-ui'

const route = useRoute()
const store = useProjectStore()
const projectId = Number(route.params.id)

const columns: DataTableColumns<any> = [
  { title: '时间', key: 'created_at', width: 160 },
  { title: '操作', key: 'action', width: 140 },
  { title: '资源类型', key: 'resource_type' },
  { title: '用户 ID', key: 'user_id', width: 100 },
  { title: 'IP 地址', key: 'ip_address', width: 130 },
  { title: 'Request ID', key: 'request_id', width: 120 }
]

onMounted(async () => {
  store.setCurrentProject(projectId)
  await store.fetchAuditLogs(projectId)
})
</script>

<style scoped>
.audit-page h2 { font-size: 18px; font-weight: 600; margin-bottom: 20px; }
</style>
