<template>
  <AppLayout>
    <div class="characters-page">
      <h2>人物档案 <span class="count">{{ store.characters.length }} 人</span></h2>

      <div class="character-grid">
        <div
          v-for="c in store.characters"
          :key="c.id"
          class="character-card"
          :class="{ active: c.id === store.selectedCharacterId }"
          @click="store.selectCharacter(c.id)"
        >
          <div class="char-header">
            <div class="char-avatar">{{ c.name[0] }}</div>
            <div>
              <div class="char-name">{{ c.name }}</div>
              <div class="char-role">{{ roleLabel(c.role) }}</div>
            </div>
            <n-tag v-if="c.confidence < 0.7" size="tiny" :bordered="false" type="warning">
              置信度 {{ (c.confidence * 100).toFixed(0) }}%
            </n-tag>
          </div>
          <p class="char-desc">{{ c.description }}</p>
          <div class="char-tags">
            <n-tag v-for="p in c.personality" :key="p" size="tiny" :bordered="false">{{ p }}</n-tag>
          </div>
          <div v-if="c.aliases.length" class="char-aliases">
            别名：{{ c.aliases.join('、') }}
          </div>
          <div v-if="c.relationships.length" class="char-relations">
            <div v-for="r in c.relationships" :key="r.target" class="relation-item">
              {{ getCharName(r.target) }}：{{ r.relation }}
            </div>
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

function roleLabel(r: string) {
  const map: Record<string, string> = { protagonist: '主角', antagonist: '反派', supporting: '配角', minor: '龙套' }
  return map[r] || r
}

function getCharName(id: string) {
  return store.characters.find(c => c.id === id)?.name || id
}
</script>

<style scoped>
.characters-page h2 { font-size: 18px; font-weight: 600; margin-bottom: 20px; }
.count { font-size: 13px; color: var(--color-text-secondary); font-weight: 400; }
.character-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 16px; }
.character-card {
  background: var(--color-surface); border: 1px solid var(--color-border); border-radius: 8px;
  padding: 16px; cursor: pointer; transition: border-color 0.15s;
}
.character-card:hover { border-color: var(--color-primary); }
.character-card.active { border-color: var(--color-primary); background: var(--color-primary-light); }
.char-header { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; }
.char-avatar {
  width: 36px; height: 36px; border-radius: 50%;
  background: var(--color-primary-light); color: var(--color-primary);
  display: flex; align-items: center; justify-content: center;
  font-size: 14px; font-weight: 600;
}
.char-name { font-size: 14px; font-weight: 600; }
.char-role { font-size: 11px; color: var(--color-text-secondary); }
.char-desc { font-size: 12px; color: var(--color-text-secondary); line-height: 1.5; margin-bottom: 10px; }
.char-tags { display: flex; flex-wrap: wrap; gap: 4px; margin-bottom: 8px; }
.char-aliases { font-size: 11px; color: var(--color-text-secondary); margin-bottom: 6px; }
.char-relations { font-size: 11px; color: var(--color-text-secondary); }
.relation-item { padding: 2px 0; }
</style>
