<template>
  <div class="sidebar">
    <div class="sidebar-logo" @click="router.push('/dashboard')">
      <span class="logo-text">N2S</span>
      <span class="logo-sub">Novel2Script</span>
    </div>

    <div class="sidebar-section" v-if="project">
      <div class="section-label">当前项目</div>
      <div class="project-info">
        <div class="project-title">{{ project.title }}</div>
        <StatusTag :status="project.status" />
      </div>
    </div>

    <div class="sidebar-section">
      <div class="section-label">项目列表</div>
      <div
        v-for="p in store.projects"
        :key="p.id"
        class="nav-item"
        :class="{ active: p.id === store.currentProjectId }"
        @click="store.setCurrentProject(p.id)"
      >
        <span class="nav-label">{{ p.title }}</span>
        <n-badge :value="p.chapter_count" :max="99" type="info" />
      </div>
    </div>

    <n-divider v-if="project" style="margin: 8px 0;" />

    <div class="sidebar-section" v-if="project">
      <div class="section-label">工作流</div>
      <router-link
        v-for="item in navItems"
        :key="item.path"
        :to="`/projects/${store.currentProjectId}${item.path}`"
        class="nav-item"
        :class="{ active: route.path === `/projects/${store.currentProjectId}${item.path}` }"
      >
        <n-icon size="16" :component="item.icon" />
        <span class="nav-label">{{ item.label }}</span>
        <n-badge v-if="item.badge" :value="item.badge" :type="item.badgeType || 'info'" />
      </router-link>
    </div>

    <div class="sidebar-footer">
      <n-button block @click="router.push(`/projects/${store.currentProjectId}/upload`)">
        上传小说
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import StatusTag from './StatusTag.vue'
import {
  CloudUploadOutline,
  GitNetworkOutline,
  CreateOutline,
  PeopleOutline,
  TimeOutline,
  FilmOutline,
  CodeSlashOutline,
  ShieldCheckmarkOutline,
  DocumentTextOutline,
  ListOutline
} from '@vicons/ionicons5'

const route = useRoute()
const router = useRouter()
const store = useProjectStore()

const project = computed(() => store.currentProject)

const navItems = computed(() => {
  const unresolved = store.unresolvedIssues.length
  return [
    { path: '/workflow', label: '生成工作流', icon: GitNetworkOutline },
    { path: '/workbench', label: '创作工作台', icon: CreateOutline },
    { path: '/characters', label: '人物档案', icon: PeopleOutline, badge: store.characters.length },
    { path: '/plot', label: '剧情事件链', icon: TimeOutline, badge: store.plotEvents.length },
    { path: '/scenes', label: '场景板', icon: FilmOutline, badge: store.scenes.length },
    { path: '/yaml', label: 'YAML 编辑', icon: CodeSlashOutline },
    { path: '/validation', label: '校验报告', icon: ShieldCheckmarkOutline, badge: unresolved > 0 ? unresolved : undefined, badgeType: unresolved > 0 ? 'warning' : 'info' },
    { path: '/versions', label: '版本历史', icon: DocumentTextOutline },
    { path: '/audit', label: '审计日志', icon: ListOutline }
  ]
})
</script>

<style scoped>
.sidebar {
  width: var(--sidebar-width);
  background: var(--color-surface);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  flex-shrink: 0;
}
.sidebar-logo {
  padding: 16px;
  cursor: pointer;
  display: flex;
  align-items: baseline;
  gap: 8px;
}
.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: var(--color-primary);
}
.logo-sub {
  font-size: 11px;
  color: var(--color-text-secondary);
}
.sidebar-section {
  padding: 8px 12px;
}
.section-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 4px 8px;
  margin-bottom: 4px;
}
.project-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px;
  background: var(--color-primary-light);
  border-radius: 6px;
}
.project-title {
  font-size: 13px;
  font-weight: 500;
}
.nav-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border-radius: 6px;
  font-size: 13px;
  color: var(--color-text-secondary);
  cursor: pointer;
  text-decoration: none;
  transition: all 0.15s;
}
.nav-item:hover {
  background: var(--color-bg);
  color: var(--color-text);
}
.nav-item.active {
  background: var(--color-primary-light);
  color: var(--color-primary);
  font-weight: 500;
}
.nav-label {
  flex: 1;
}
.sidebar-footer {
  margin-top: auto;
  padding: 12px;
}
</style>
