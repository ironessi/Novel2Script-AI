<template>
  <AppLayout>
    <div class="validation-page">
      <h2>校验报告</h2>

      <!-- 概览 -->
      <div class="summary-row">
        <div class="summary-card">
          <div class="summary-num">{{ store.validationIssues.length }}</div>
          <div class="summary-label">总问题数</div>
        </div>
        <div class="summary-card">
          <div class="summary-num high">{{ highCount }}</div>
          <div class="summary-label">高风险</div>
        </div>
        <div class="summary-card">
          <div class="summary-num medium">{{ mediumCount }}</div>
          <div class="summary-label">中风险</div>
        </div>
        <div class="summary-card">
          <div class="summary-num low">{{ lowCount }}</div>
          <div class="summary-label">低风险</div>
        </div>
      </div>

      <!-- 问题列表 -->
      <n-data-table :columns="columns" :data="store.validationIssues" :bordered="false" size="small" />
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, h, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import AppLayout from '@/layouts/AppLayout.vue'
import RiskBadge from '@/components/RiskBadge.vue'
import type { DataTableColumns } from 'naive-ui'
import { NButton } from 'naive-ui'

const route = useRoute()
const router = useRouter()
const store = useProjectStore()
const projectId = Number(route.params.id)

onMounted(async () => {
  store.setCurrentProject(projectId)
  await store.fetchValidationIssues(projectId)
})

const highCount = computed(() => store.validationIssues.filter(i => i.severity === 'high').length)
const mediumCount = computed(() => store.validationIssues.filter(i => i.severity === 'medium').length)
const lowCount = computed(() => store.validationIssues.filter(i => i.severity === 'low').length)

const columns: DataTableColumns<any> = [
  { title: '类型', key: 'issue_type', width: 120 },
  {
    title: '严重程度', key: 'severity', width: 80,
    render: (row) => h(RiskBadge, { level: row.severity })
  },
  { title: '说明', key: 'message' },
  { title: '位置', key: 'location_path', width: 180 },
  { title: '建议', key: 'suggestion' },
  {
    title: '状态', key: 'resolved', width: 80,
    render: (row) => row.resolved ? '已解决' : '待处理'
  },
  {
    title: '操作', key: 'actions', width: 100,
    render: (row) => h(NButton, { size: 'tiny', onClick: () => router.push(`/projects/${store.currentProjectId}/yaml`) }, () => '去修复')
  }
]
</script>

<style scoped>
.validation-page h2 { font-size: 18px; font-weight: 600; margin-bottom: 20px; }
.summary-row { display: flex; gap: 16px; margin-bottom: 24px; }
.summary-card {
  flex: 1; background: var(--color-surface); border: 1px solid var(--color-border);
  border-radius: 8px; padding: 16px; text-align: center;
}
.summary-num { font-size: 24px; font-weight: 700; }
.summary-num.high { color: #ef4444; }
.summary-num.medium { color: #f59e0b; }
.summary-num.low { color: #6b7280; }
.summary-label { font-size: 12px; color: var(--color-text-secondary); margin-top: 4px; }
</style>
