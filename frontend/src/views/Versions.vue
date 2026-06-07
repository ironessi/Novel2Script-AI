<template>
  <AppLayout>
    <div class="versions-page">
      <h2>版本历史</h2>

      <div class="version-list">
        <div v-for="v in store.versions" :key="v.id" class="version-item">
          <div class="version-header">
            <span class="version-no">v{{ v.version_no }}</span>
            <n-tag size="tiny" :bordered="false" :type="createdByType(v.created_by)">
              {{ createdByLabel(v.created_by) }}
            </n-tag>
            <span class="version-time">{{ v.created_at }}</span>
          </div>
          <div class="version-status">
            <span>校验：<StatusTag :status="v.validation_status" /></span>
            <span>幻觉：<RiskBadge :level="v.hallucination_risk" /></span>
            <span>安全：<RiskBadge :level="v.safety_risk" /></span>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { useProjectStore } from '@/stores/project'
import AppLayout from '@/layouts/AppLayout.vue'
import StatusTag from '@/components/StatusTag.vue'
import RiskBadge from '@/components/RiskBadge.vue'

const store = useProjectStore()

function createdByType(c: string) {
  const map: Record<string, 'success' | 'warning' | 'info'> = { ai: 'info', user: 'success', system_repair: 'warning' }
  return map[c] || 'info'
}

function createdByLabel(c: string) {
  const map: Record<string, string> = { ai: 'AI 生成', user: '用户编辑', system_repair: '系统修复' }
  return map[c] || c
}
</script>

<style scoped>
.versions-page h2 { font-size: 18px; font-weight: 600; margin-bottom: 20px; }
.version-list { display: flex; flex-direction: column; gap: 12px; }
.version-item {
  background: var(--color-surface); border: 1px solid var(--color-border);
  border-radius: 8px; padding: 16px;
}
.version-header { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; }
.version-no { font-size: 16px; font-weight: 700; color: var(--color-primary); }
.version-time { font-size: 12px; color: var(--color-text-secondary); }
.version-status { display: flex; gap: 16px; font-size: 13px; align-items: center; }
</style>
