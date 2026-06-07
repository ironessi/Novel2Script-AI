<template>
  <div class="app-layout">
    <SidebarNav />
    <div class="app-main">
      <TopBar />
      <div class="app-content">
        <slot />
      </div>
    </div>
    <InspectorPanel />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import SidebarNav from '@/components/SidebarNav.vue'
import TopBar from '@/components/TopBar.vue'
import InspectorPanel from '@/components/InspectorPanel.vue'

const route = useRoute()
const store = useProjectStore()

onMounted(async () => {
  if (store.projects.length === 0) {
    await store.fetchProjects()
  }
  const projectId = Number(route.params.id)
  if (projectId) {
    store.setCurrentProject(projectId)
  }
})
</script>

<style scoped>
.app-layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}
.app-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}
.app-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px 24px;
  background: var(--color-bg);
}

@media (max-width: 1100px) {
  .app-layout :deep(.inspector) {
    display: none;
  }
}

@media (max-width: 720px) {
  .app-layout :deep(.sidebar) {
    display: none;
  }

  .app-content {
    padding: 16px;
  }
}
</style>
