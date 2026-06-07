<template>
  <AppLayout>
    <div class="workbench">
      <n-tabs v-model:value="activeTab" type="line" animated>
        <n-tab-pane name="script" tab="YAML 剧本">
          <div class="script-preview">
            <pre>{{ store.yamlScript }}</pre>
          </div>
        </n-tab-pane>
        <n-tab-pane name="characters" tab="人物档案">
          <div class="mini-grid">
            <div v-for="c in store.characters" :key="c.id" class="mini-card">
              <div class="mini-name">{{ c.name }}</div>
              <div class="mini-desc">{{ c.description }}</div>
            </div>
          </div>
        </n-tab-pane>
        <n-tab-pane name="scenes" tab="场景板">
          <div class="mini-grid">
            <div
              v-for="s in store.scenes"
              :key="s.id"
              class="mini-card"
              @click="store.selectScene(s.id)"
            >
              <div class="mini-name">#{{ s.order }} {{ s.title }}</div>
              <div class="mini-desc">{{ s.time }} · {{ s.location }}</div>
            </div>
          </div>
        </n-tab-pane>
        <n-tab-pane name="plot" tab="剧情事件链">
          <div v-for="e in store.plotEvents" :key="e.id" class="plot-mini">
            <n-tag :type="importanceType(e.importance)" size="tiny" :bordered="false">
              {{ importanceLabel(e.importance) }}
            </n-tag>
            <span>{{ e.trigger }} → {{ e.result }}</span>
          </div>
        </n-tab-pane>
        <n-tab-pane name="original" tab="原文视图">
          <div v-for="ch in store.chapters" :key="ch.id" class="original-chapter">
            <h3>{{ ch.chapter_title }}</h3>
            <p v-for="(para, i) in ch.content.split('\n\n')" :key="i">{{ para }}</p>
          </div>
        </n-tab-pane>
      </n-tabs>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useProjectStore } from '@/stores/project'
import AppLayout from '@/layouts/AppLayout.vue'

const store = useProjectStore()
const activeTab = ref('script')

function importanceType(i: string) {
  const map: Record<string, 'success' | 'warning' | 'error' | 'info'> = { high: 'error', medium: 'warning', low: 'info' }
  return map[i] || 'info'
}

function importanceLabel(i: string) {
  const map: Record<string, string> = { high: '关键', medium: '重要', low: '次要' }
  return map[i] || i
}
</script>

<style scoped>
.workbench { height: calc(100vh - var(--topbar-height) - 40px); }
.script-preview {
  background: #1e1e2e; color: #cdd6f4; padding: 20px; border-radius: 8px;
  font-family: 'SF Mono', 'Fira Code', monospace; font-size: 13px; line-height: 1.6;
  overflow: auto; max-height: calc(100vh - 160px); white-space: pre-wrap;
}
.mini-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 12px; }
.mini-card {
  background: var(--color-surface); border: 1px solid var(--color-border);
  border-radius: 6px; padding: 12px; cursor: pointer; transition: border-color 0.15s;
}
.mini-card:hover { border-color: var(--color-primary); }
.mini-name { font-size: 13px; font-weight: 500; margin-bottom: 4px; }
.mini-desc { font-size: 12px; color: var(--color-text-secondary); }
.plot-mini { display: flex; align-items: center; gap: 8px; padding: 6px 0; font-size: 13px; }
.original-chapter { margin-bottom: 24px; }
.original-chapter h3 { font-size: 15px; font-weight: 600; margin-bottom: 10px; }
.original-chapter p { font-size: 14px; line-height: 1.8; margin-bottom: 12px; color: var(--color-text); }
</style>
