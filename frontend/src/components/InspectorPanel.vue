<template>
  <div class="inspector">
    <div class="inspector-header">
      <span class="inspector-title">检查面板</span>
    </div>

    <!-- 任务进度 -->
    <div class="inspector-section">
      <div class="section-title">生成进度</div>
      <div v-for="step in store.workflowSteps" :key="step.key" class="step-row">
        <n-icon v-if="step.status === 'completed'" size="14" color="#22c55e"><CheckmarkCircleOutline /></n-icon>
        <n-icon v-else-if="step.status === 'running'" size="14" color="#f59e0b"><SyncOutline /></n-icon>
        <n-icon v-else-if="step.status === 'warning'" size="14" color="#f59e0b"><AlertCircleOutline /></n-icon>
        <n-icon v-else-if="step.status === 'failed'" size="14" color="#ef4444"><CloseCircleOutline /></n-icon>
        <n-icon v-else size="14" color="#d1d5db"><EllipseOutline /></n-icon>
        <span class="step-label">{{ step.label }}</span>
        <span v-if="step.count" class="step-count">{{ step.count }}{{ step.unit }}</span>
      </div>
    </div>

    <!-- 校验问题 -->
    <div class="inspector-section">
      <div class="section-title">
        校验问题
        <n-badge :value="store.unresolvedIssues.length" :max="99" type="warning" />
      </div>
      <div v-if="store.unresolvedIssues.length === 0" class="empty-hint">暂无未解决问题</div>
      <div v-for="issue in store.unresolvedIssues.slice(0, 3)" :key="issue.id" class="issue-row">
        <RiskBadge :level="issue.severity" />
        <span class="issue-msg">{{ issue.message }}</span>
      </div>
    </div>

    <!-- Source Trace -->
    <div class="inspector-section" v-if="selectedScene">
      <div class="section-title">原文溯源</div>
      <div class="trace-info">
        <div class="trace-scene">{{ selectedScene.title }}</div>
        <div v-for="(trace, i) in selectedScene.source_trace" :key="i" class="trace-item">
          第{{ trace.chapter_index }}章 第{{ trace.paragraph_start }}-{{ trace.paragraph_end }}段
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useProjectStore } from '@/stores/project'
import RiskBadge from './RiskBadge.vue'
import {
  CheckmarkCircleOutline,
  SyncOutline,
  AlertCircleOutline,
  CloseCircleOutline,
  EllipseOutline
} from '@vicons/ionicons5'

const store = useProjectStore()
const selectedScene = computed(() => store.selectedScene)
</script>

<style scoped>
.inspector {
  width: var(--inspector-width);
  background: var(--color-surface);
  border-left: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  flex-shrink: 0;
}
.inspector-header {
  padding: 14px 16px;
  border-bottom: 1px solid var(--color-border);
}
.inspector-title {
  font-size: 13px;
  font-weight: 600;
}
.inspector-section {
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-border);
}
.section-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  gap: 6px;
}
.step-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;
  font-size: 12px;
}
.step-label {
  flex: 1;
}
.step-count {
  color: var(--color-primary);
  font-weight: 500;
  font-size: 11px;
}
.issue-row {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 6px 0;
  font-size: 12px;
}
.issue-msg {
  flex: 1;
  line-height: 1.4;
}
.empty-hint {
  font-size: 12px;
  color: #9ca3af;
}
.trace-info {
  font-size: 12px;
}
.trace-scene {
  font-weight: 500;
  margin-bottom: 6px;
}
.trace-item {
  color: var(--color-text-secondary);
  padding: 2px 0;
}
</style>
