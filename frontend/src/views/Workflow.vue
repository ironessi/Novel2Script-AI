<template>
  <AppLayout>
    <div class="workflow-page">
      <h2>AI 分析工作流</h2>

      <!-- 触发按钮 -->
      <div v-if="!taskId" class="trigger-section">
        <n-button type="primary" size="large" :loading="starting" @click="handleStart">
          开始 AI 分析
        </n-button>
        <p class="trigger-hint">将对已上传的章节进行完整分析：人物抽取 → 剧情事件 → 场景拆分 → 剧本生成</p>
      </div>

      <!-- 工作流步骤 -->
      <div v-if="taskId" class="workflow-steps">
        <div
          v-for="(step, i) in steps"
          :key="step.key"
          class="workflow-step"
          :class="step.status"
        >
          <div class="step-indicator">
            <div class="step-number">{{ i + 1 }}</div>
            <div v-if="i < steps.length - 1" class="step-line" />
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

      <!-- 任务信息 -->
      <div v-if="taskId" class="task-info">
        <p>任务 ID: {{ taskId }} | 状态: {{ taskStatus }} | 进度: {{ taskProgress }}%</p>
        <p v-if="taskStep">当前步骤: {{ taskStep }}</p>
        <p v-if="taskError" class="error-text">{{ taskError }}</p>
      </div>

      <!-- 完成后跳转 -->
      <div v-if="taskStatus === 'completed'" class="complete-section">
        <n-alert type="success" title="分析完成">
          AI 分析已完成，可以查看生成的剧本。
        </n-alert>
        <n-button type="primary" style="margin-top: 16px;" @click="router.push(`/projects/${projectId}/workbench`)">
          进入创作工作台
        </n-button>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import AppLayout from '@/layouts/AppLayout.vue'

const route = useRoute()
const router = useRouter()
const store = useProjectStore()
const projectId = Number(route.params.id)

const starting = ref(false)
const taskId = ref<number | null>(null)
const taskStatus = ref('')
const taskProgress = ref(0)
const taskStep = ref('')
const taskError = ref('')
let pollTimer: ReturnType<typeof setInterval> | null = null

const steps = ref([
  { key: 'pending', label: '等待中', status: 'pending' as string, count: 0, unit: '', detail: '' },
  { key: 'running', label: 'AI 处理中', status: 'pending' as string, count: 0, unit: '', detail: '' },
  { key: 'completed', label: '完成', status: 'pending' as string, count: 0, unit: '', detail: '' }
])

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

async function handleStart() {
  starting.value = true
  try {
    const id = await store.triggerGenerate(projectId)
    taskId.value = id
    taskStatus.value = 'pending'
    startPolling()
  } catch (e: any) {
    alert(e.message || '创建任务失败')
  } finally {
    starting.value = false
  }
}

function startPolling() {
  pollTimer = setInterval(async () => {
    if (!taskId.value) return
    try {
      const status = await store.checkTaskStatus(taskId.value)
      taskStatus.value = status.status
      taskProgress.value = status.progress
      taskStep.value = status.current_step || ''
      taskError.value = status.error_message || ''

      // 更新步骤状态
      if (status.status === 'running') {
        steps.value[0].status = 'completed'
        steps.value[1].status = 'running'
      } else if (status.status === 'completed') {
        steps.value[0].status = 'completed'
        steps.value[1].status = 'completed'
        steps.value[2].status = 'completed'
        stopPolling()
      } else if (status.status === 'failed') {
        steps.value[1].status = 'failed'
        stopPolling()
      }
    } catch (e) {
      console.error('轮询失败', e)
    }
  }, 3000)
}

function stopPolling() {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

onMounted(() => {
  store.setCurrentProject(projectId)
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.workflow-page { max-width: 640px; }
.workflow-page h2 { font-size: 18px; font-weight: 600; margin-bottom: 24px; }
.trigger-section { text-align: center; padding: 48px 0; }
.trigger-hint { font-size: 13px; color: var(--color-text-secondary); margin-top: 12px; }
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
.workflow-step.failed .step-number { background: #fee2e2; color: #dc2626; }
.step-line { width: 2px; flex: 1; min-height: 24px; background: #e5e7eb; margin: 4px 0; }
.workflow-step.completed .step-line { background: #86efac; }
.step-content { flex: 1; padding-bottom: 20px; }
.step-header { display: flex; align-items: center; gap: 8px; }
.step-label { font-size: 14px; font-weight: 500; }
.step-detail { font-size: 12px; color: var(--color-text-secondary); margin-top: 4px; }
.task-info { margin-top: 24px; padding: 16px; background: #f9fafb; border-radius: 8px; font-size: 13px; }
.error-text { color: #ef4444; }
.complete-section { margin-top: 24px; }
</style>
