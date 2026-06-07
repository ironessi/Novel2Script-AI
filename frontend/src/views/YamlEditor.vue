<template>
  <AppLayout>
    <div class="yaml-page">
      <div class="yaml-toolbar">
        <n-space>
          <n-button size="small" @click="handleFormat">格式化</n-button>
          <n-button size="small" @click="handleValidate">Schema 校验</n-button>
          <n-button size="small" type="primary" @click="handleSave">保存版本</n-button>
          <n-button size="small" @click="handleExportYaml">导出 YAML</n-button>
          <n-button size="small" @click="handleExportMd">导出 Markdown</n-button>
        </n-space>
        <n-space>
          <RiskBadge :level="validationRisk" />
          <RiskBadge :level="hallucinationRisk" />
        </n-space>
      </div>

      <div class="yaml-content">
        <div class="yaml-editor">
          <textarea
            v-model="yamlContent"
            class="editor-textarea"
            spellcheck="false"
          />
        </div>
        <div class="yaml-preview">
          <h4>结构化预览</h4>
          <div v-if="parsed" class="preview-content">
            <div class="preview-section">
              <h5>元数据</h5>
              <p>标题：{{ parsed.script?.metadata?.title }}</p>
              <p>模式：{{ parsed.script?.metadata?.adaptation_mode }}</p>
            </div>
            <div class="preview-section">
              <h5>人物 ({{ parsed.script?.characters?.length || 0 }})</h5>
              <div v-for="c in parsed.script?.characters" :key="c.id" class="preview-item">
                {{ c.name }} · {{ c.role }}
              </div>
            </div>
            <div class="preview-section">
              <h5>场景 ({{ parsed.script?.scenes?.length || 0 }})</h5>
              <div v-for="s in parsed.script?.scenes" :key="s.id" class="preview-item">
                #{{ s.order }} {{ s.title }} · {{ s.time }} · {{ s.location }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useProjectStore } from '@/stores/project'
import * as yaml from 'js-yaml'
import AppLayout from '@/layouts/AppLayout.vue'
import RiskBadge from '@/components/RiskBadge.vue'

const store = useProjectStore()
const yamlContent = ref(store.yamlScript)

const parsed = computed(() => {
  try { return yaml.load(yamlContent.value) as any } catch { return null }
})

const validationRisk = computed(() => {
  const high = store.validationIssues.filter(i => i.severity === 'high' && !i.resolved).length
  return high > 0 ? 'high' : store.unresolvedIssues.length > 0 ? 'medium' : 'low'
})

const hallucinationRisk = computed(() => 'low' as const)

function handleFormat() {
  try {
    const obj = yaml.load(yamlContent.value)
    yamlContent.value = yaml.dump(obj, { indent: 2, lineWidth: 120 })
  } catch {}
}

function handleValidate() {
  // Mock
}

function handleSave() {
  store.updateYaml(yamlContent.value)
}

function handleExportYaml() {
  const blob = new Blob([yamlContent.value], { type: 'text/yaml' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url; a.download = 'script.yaml'; a.click()
  URL.revokeObjectURL(url)
}

function handleExportMd() {
  // Mock
}
</script>

<style scoped>
.yaml-page { height: calc(100vh - var(--topbar-height) - 40px); display: flex; flex-direction: column; }
.yaml-toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.yaml-content { flex: 1; display: flex; gap: 12px; min-height: 0; }
.yaml-editor { flex: 1; }
.editor-textarea {
  width: 100%; height: 100%; padding: 16px; border: 1px solid var(--color-border);
  border-radius: 8px; font-family: 'SF Mono', 'Fira Code', monospace; font-size: 13px;
  line-height: 1.6; resize: none; outline: none; tab-size: 2; background: #1e1e2e; color: #cdd6f4;
}
.editor-textarea:focus { border-color: var(--color-primary); }
.yaml-preview {
  width: 320px; background: var(--color-surface); border: 1px solid var(--color-border);
  border-radius: 8px; padding: 16px; overflow-y: auto;
}
.yaml-preview h4 { font-size: 13px; font-weight: 600; margin-bottom: 12px; }
.preview-section { margin-bottom: 16px; }
.preview-section h5 { font-size: 12px; color: var(--color-text-secondary); margin-bottom: 6px; }
.preview-item { font-size: 12px; padding: 3px 0; }
</style>
