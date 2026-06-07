<template>
  <AppLayout>
    <div class="plot-page">
      <h2>剧情事件链 <span class="count">{{ store.plotEvents.length }} 个事件</span></h2>

      <div class="timeline">
        <div v-for="e in store.plotEvents" :key="e.id" class="event-item" :class="e.importance">
          <div class="event-dot" />
          <div class="event-content">
            <div class="event-header">
              <n-tag :type="importanceType(e.importance)" size="small" :bordered="false">
                {{ importanceLabel(e.importance) }}
              </n-tag>
              <span class="event-chapter">第{{ e.chapter_index }}章</span>
            </div>
            <div class="event-parts">
              <div class="event-part">
                <span class="part-label">触发</span>
                <span>{{ e.trigger }}</span>
              </div>
              <div class="event-part">
                <span class="part-label">行动</span>
                <span>{{ e.action }}</span>
              </div>
              <div class="event-part">
                <span class="part-label">结果</span>
                <span>{{ e.result }}</span>
              </div>
            </div>
            <div v-if="e.characters_involved.length" class="event-chars">
              涉及：{{ e.characters_involved.map(id => getCharName(id)).join('、') }}
            </div>
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

const route = useRoute()
const store = useProjectStore()
const projectId = Number(route.params.id)

onMounted(async () => {
  store.setCurrentProject(projectId)
  await store.fetchPlotEvents(projectId)
})

function importanceType(i: string) {
  const map: Record<string, 'success' | 'warning' | 'error' | 'info'> = { high: 'error', medium: 'warning', low: 'info' }
  return map[i] || 'info'
}

function importanceLabel(i: string) {
  const map: Record<string, string> = { high: '关键', medium: '重要', low: '次要' }
  return map[i] || i
}

function getCharName(id: string) {
  return store.characters.find(c => c.id === id)?.name || id
}
</script>

<style scoped>
.plot-page h2 { font-size: 18px; font-weight: 600; margin-bottom: 20px; }
.count { font-size: 13px; color: var(--color-text-secondary); font-weight: 400; }
.timeline { position: relative; padding-left: 24px; }
.timeline::before {
  content: ''; position: absolute; left: 7px; top: 8px; bottom: 8px;
  width: 2px; background: var(--color-border);
}
.event-item { position: relative; margin-bottom: 16px; }
.event-dot {
  position: absolute; left: -20px; top: 6px; width: 10px; height: 10px;
  border-radius: 50%; background: #d1d5db; border: 2px solid var(--color-surface);
}
.event-item.high .dot, .event-item.high .event-dot { background: #ef4444; }
.event-item.medium .event-dot { background: #f59e0b; }
.event-content {
  background: var(--color-surface); border: 1px solid var(--color-border);
  border-radius: 8px; padding: 14px 16px;
}
.event-item.high .event-content { border-left: 3px solid #ef4444; }
.event-item.medium .event-content { border-left: 3px solid #f59e0b; }
.event-header { display: flex; align-items: center; gap: 8px; margin-bottom: 10px; }
.event-chapter { font-size: 12px; color: var(--color-text-secondary); }
.event-parts { display: flex; flex-direction: column; gap: 6px; margin-bottom: 8px; }
.event-part { font-size: 13px; display: flex; gap: 8px; }
.part-label {
  font-size: 11px; font-weight: 600; color: var(--color-primary);
  background: var(--color-primary-light); padding: 1px 6px; border-radius: 3px;
  flex-shrink: 0; height: fit-content;
}
.event-chars { font-size: 11px; color: var(--color-text-secondary); }
</style>
