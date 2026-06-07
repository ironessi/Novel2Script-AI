<template>
  <AppLayout>
    <div class="scenes-page">
      <h2>场景板 <span class="count">{{ store.scenes.length }} 个场景</span></h2>

      <div class="scene-board">
        <div
          v-for="s in store.scenes"
          :key="s.id"
          class="scene-card"
          :class="{ active: s.id === store.selectedSceneId }"
          @click="store.selectScene(s.id)"
        >
          <div class="scene-header">
            <span class="scene-order">#{{ s.order }}</span>
            <span class="scene-title">{{ s.title }}</span>
            <RiskBadge v-if="s.risk_level !== 'none'" :level="s.risk_level" />
          </div>
          <div class="scene-meta">
            <span>{{ s.time }}</span>
            <span>{{ s.location }}</span>
          </div>
          <p class="scene-summary">{{ s.summary }}</p>
          <div class="scene-chars">
            <n-tag v-for="cid in s.characters" :key="cid" size="tiny" :bordered="false">
              {{ getCharName(cid) }}
            </n-tag>
          </div>
          <div v-if="s.risk_type" class="scene-risk">
            <n-icon size="12" color="#f59e0b"><AlertCircleOutline /></n-icon>
            {{ s.risk_type }}
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import AppLayout from '@/layouts/AppLayout.vue'
import RiskBadge from '@/components/RiskBadge.vue'
import { AlertCircleOutline } from '@vicons/ionicons5'

const route = useRoute()
const store = useProjectStore()
const projectId = Number(route.params.id)

onMounted(async () => {
  store.setCurrentProject(projectId)
  await store.fetchScript(projectId)
})

function getCharName(id: string) {
  return store.characters.find(c => c.id === id)?.name || id
}
</script>

<style scoped>
.scenes-page h2 { font-size: 18px; font-weight: 600; margin-bottom: 20px; }
.count { font-size: 13px; color: var(--color-text-secondary); font-weight: 400; }
.scene-board { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 14px; }
.scene-card {
  background: var(--color-surface); border: 1px solid var(--color-border);
  border-radius: 8px; padding: 14px; cursor: pointer; transition: border-color 0.15s;
}
.scene-card:hover { border-color: var(--color-primary); }
.scene-card.active { border-color: var(--color-primary); background: var(--color-primary-light); }
.scene-header { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.scene-order { font-size: 12px; color: var(--color-text-secondary); font-weight: 600; }
.scene-title { font-size: 14px; font-weight: 500; flex: 1; }
.scene-meta { display: flex; gap: 12px; font-size: 11px; color: var(--color-text-secondary); margin-bottom: 8px; }
.scene-summary { font-size: 12px; color: var(--color-text-secondary); line-height: 1.5; margin-bottom: 8px; }
.scene-chars { display: flex; flex-wrap: wrap; gap: 4px; margin-bottom: 6px; }
.scene-risk { font-size: 11px; color: #d97706; display: flex; align-items: center; gap: 4px; }
</style>
