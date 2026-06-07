<template>
  <AppLayout>
    <div class="workflow-page">
      <h2>AI 分析工作流</h2>

      <div class="workflow-steps">
        <div
          v-for="(step, i) in store.workflowSteps"
          :key="step.key"
          class="workflow-step"
          :class="step.status"
        >
          <div class="step-indicator">
            <div class="step-number">{{ i + 1 }}</div>
            <div v-if="i < store.workflowSteps.length - 1" class="step-line" />
          </div>
          <div class="step-content">
            <div class="step-header">
              <span class="step-label">{{ step.label }}</span>
              <n-tag v-if="step.count" size="small" :bordered="false" type="info">
                {{ step.count }}{{ step.unit }}
              </n-tag>
              <n-tag :type="statusType(step.status)" size="small" :bordered="false">
                {{ statusLabel(step.status) }}
              </n-tag>
            </div>
            <p v-if="step.detail" class="step-detail">{{ step.detail }}</p>
          </div>
        </div>
      </div>

      <!-- 生成日志 -->
      <div class="log-section">
        <h3>生成日志</h3>
        <div class="log-list">
          <div v-for="(log, i) in logs" :key="i" class="log-item">
            <span class="log-time">{{ log.time }}</span>
            <span class="log-msg">{{ log.msg }}</span>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { useProjectStore } from '@/stores/project'
import AppLayout from '@/layouts/AppLayout.vue'

const store = useProjectStore()

function statusType(s: string) {
  const map: Record<string, 'success' | 'warning' | 'error' | 'info' | 'default'> = {
    completed: 'success', running: 'warning', warning: 'warning', failed: 'error', pending: 'default'
  }
  return map[s] || 'default'
}

function statusLabel(s: string) {
  const map: Record<string, string> = {
    completed: '已完成', running: '处理中', warning: '有警告', failed: '失败', pending: '等待中'
  }
  return map[s] || s
}

const logs = [
  { time: '12:00:01', msg: '开始处理 3 个章节，共 366 字' },
  { time: '12:00:03', msg: '文本清洗完成，去除 12 处乱码' },
  { time: '12:00:05', msg: '识别到 3 个章节标题' },
  { time: '12:00:15', msg: '人物抽取完成，提取 5 个人物档案' },
  { time: '12:00:25', msg: '剧情事件链构建完成，8 个事件' },
  { time: '12:00:35', msg: '场景拆分完成，6 个场景' },
  { time: '12:01:00', msg: 'YAML 剧本生成完成' },
  { time: '12:01:05', msg: 'Schema 校验发现 5 个问题' },
  { time: '12:01:10', msg: '幻觉检测完成，风险等级：低' },
  { time: '12:01:15', msg: '安全审查完成，未发现风险' }
]
</script>

<style scoped>
.workflow-page { max-width: 640px; }
.workflow-page h2 { font-size: 18px; font-weight: 600; margin-bottom: 24px; }
.workflow-steps { display: flex; flex-direction: column; gap: 0; }
.workflow-step { display: flex; gap: 16px; }
.step-indicator { display: flex; flex-direction: column; align-items: center; }
.step-number {
  width: 28px; height: 28px; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  font-size: 12px; font-weight: 600;
  background: #e5e7eb; color: #6b7280;
}
.workflow-step.completed .step-number { background: #dcfce7; color: #16a34a; }
.workflow-step.running .step-number { background: #fef3c7; color: #d97706; }
.workflow-step.warning .step-number { background: #fef3c7; color: #d97706; }
.workflow-step.failed .step-number { background: #fee2e2; color: #dc2626; }
.step-line { width: 2px; flex: 1; min-height: 24px; background: #e5e7eb; margin: 4px 0; }
.workflow-step.completed .step-line { background: #86efac; }
.step-content { flex: 1; padding-bottom: 20px; }
.step-header { display: flex; align-items: center; gap: 8px; }
.step-label { font-size: 14px; font-weight: 500; }
.step-detail { font-size: 12px; color: var(--color-text-secondary); margin-top: 4px; }
.log-section { margin-top: 32px; }
.log-section h3 { font-size: 15px; font-weight: 600; margin-bottom: 12px; }
.log-list {
  background: #f9fafb; border: 1px solid var(--color-border); border-radius: 8px;
  padding: 12px; max-height: 300px; overflow-y: auto;
}
.log-item { display: flex; gap: 12px; padding: 3px 0; font-size: 12px; font-family: 'SF Mono', monospace; }
.log-time { color: var(--color-text-secondary); flex-shrink: 0; }
.log-msg { color: var(--color-text); }
</style>
