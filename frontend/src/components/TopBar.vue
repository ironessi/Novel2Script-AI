<template>
  <div class="topbar">
    <div class="topbar-left">
      <n-breadcrumb>
        <n-breadcrumb-item @click="router.push('/dashboard')">项目</n-breadcrumb-item>
        <n-breadcrumb-item v-if="project">{{ project.title }}</n-breadcrumb-item>
        <n-breadcrumb-item v-if="pageTitle">{{ pageTitle }}</n-breadcrumb-item>
      </n-breadcrumb>
    </div>
    <div class="topbar-right">
      <n-space align="center" :size="12">
        <n-tag v-if="project" :bordered="false" size="small" type="info">
          {{ modeLabel }}
        </n-tag>
        <n-dropdown :options="userOptions" @select="handleUserAction">
          <n-button text quaternary size="small">
            {{ auth.username }}
          </n-button>
        </n-dropdown>
      </n-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useProjectStore } from '@/stores/project'
import { adaptationModes } from '@/mock/data'
import type { DropdownOption } from 'naive-ui'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const store = useProjectStore()

const project = computed(() => store.currentProject)

const modeLabel = computed(() => {
  if (!project.value) return ''
  return adaptationModes[project.value.adaptation_mode]?.label || ''
})

const pageTitle = computed(() => {
  const map: Record<string, string> = {
    Upload: '上传小说',
    Workflow: '生成工作流',
    Workbench: '创作工作台',
    Characters: '人物档案',
    Plot: '剧情事件链',
    Scenes: '场景板',
    YamlEditor: 'YAML 编辑',
    Validation: '校验报告',
    Versions: '版本历史',
    Audit: '审计日志'
  }
  return map[route.name as string] || ''
})

const userOptions: DropdownOption[] = [
  { label: '退出登录', key: 'logout' }
]

function handleUserAction(key: string) {
  if (key === 'logout') {
    auth.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.topbar {
  height: var(--topbar-height);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  flex-shrink: 0;
}
</style>
