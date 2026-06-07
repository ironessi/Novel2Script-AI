<template>
  <AppLayout>
    <div class="audit-page">
      <h2>审计日志</h2>
      <n-data-table :columns="columns" :data="store.auditLogs" :bordered="false" size="small" />
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { h } from 'vue'
import { useProjectStore } from '@/stores/project'
import AppLayout from '@/layouts/AppLayout.vue'
import StatusTag from '@/components/StatusTag.vue'
import type { DataTableColumns } from 'naive-ui'

const store = useProjectStore()

const columns: DataTableColumns<any> = [
  { title: '时间', key: 'created_at', width: 160 },
  { title: '操作', key: 'action', width: 140 },
  { title: '项目', key: 'project_title' },
  { title: '用户', key: 'user', width: 100 },
  {
    title: '状态', key: 'status', width: 80,
    render: (row) => h(StatusTag, { status: row.status === '成功' ? 'completed' : 'has_risk' })
  },
  { title: 'Request ID', key: 'request_id', width: 120 }
]
</script>

<style scoped>
.audit-page h2 { font-size: 18px; font-weight: 600; margin-bottom: 20px; }
</style>
