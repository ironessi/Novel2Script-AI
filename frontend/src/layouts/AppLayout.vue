<template>
  <div class="app-layout">
    <SidebarNav />
    <div class="app-main">
      <TopBar />
      <div class="app-content">
        <router-view />
      </div>
    </div>
    <InspectorPanel />
  </div>
</template>

<script setup lang="ts">
import SidebarNav from '@/components/SidebarNav.vue'
import TopBar from '@/components/TopBar.vue'
import InspectorPanel from '@/components/InspectorPanel.vue'
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
</style>
