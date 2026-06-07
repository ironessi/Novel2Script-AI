<template>
  <div v-if="loading" class="text-center text-gray">加载中...</div>
  <div v-else>
    <div class="flex justify-between items-center mb-4">
      <h2>剧本编辑</h2>
      <div class="flex gap-2">
        <button class="btn btn-outline" @click="handleValidate" :disabled="validating">
          {{ validating ? '校验中...' : 'Schema 校验' }}
        </button>
        <button class="btn btn-outline" @click="handleCheckHallucination" :disabled="checking">
          幻觉检测
        </button>
        <button class="btn btn-outline" @click="handleCheckSafety" :disabled="checking">
          安全审查
        </button>
        <button class="btn btn-primary" @click="handleSave" :disabled="saving">
          {{ saving ? '保存中...' : '保存' }}
        </button>
        <a :href="exportUrl('yaml')" class="btn btn-success" download>导出 YAML</a>
      </div>
    </div>

    <!-- 风险状态 -->
    <div v-if="script" class="flex gap-4 mb-4">
      <span class="badge" :class="`badge-${script.validation_status}`">
        校验: {{ script.validation_status }}
      </span>
      <span class="badge" :class="`badge-${script.hallucination_risk}`">
        幻觉: {{ script.hallucination_risk }}
      </span>
      <span class="badge" :class="`badge-${script.safety_risk}`">
        安全: {{ script.safety_risk }}
      </span>
    </div>

    <!-- 校验问题 -->
    <div v-if="validationResult" class="card mb-4">
      <h3>校验结果</h3>
      <p class="mt-4">
        状态: <span class="badge" :class="validationResult.valid ? 'badge-valid' : 'badge-invalid'">
          {{ validationResult.valid ? '通过' : '不通过' }}
        </span>
      </p>
      <div v-if="validationResult.issues?.length" class="mt-4">
        <div v-for="(issue, i) in validationResult.issues" :key="i" class="issue-item">
          <span class="badge" :class="`badge-${issue.severity}`">{{ issue.severity }}</span>
          <span>{{ issue.message }}</span>
          <span v-if="issue.location_path" class="text-sm text-gray">({{ issue.location_path }})</span>
        </div>
      </div>
    </div>

    <!-- YAML 编辑器 -->
    <div class="card">
      <textarea
        v-model="yamlContent"
        class="yaml-editor"
        placeholder="YAML 剧本内容..."
        spellcheck="false"
      ></textarea>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { scriptApi } from '../api'

const route = useRoute()
const id = route.params.id

const script = ref(null)
const yamlContent = ref('')
const loading = ref(true)
const saving = ref(false)
const validating = ref(false)
const checking = ref(false)
const validationResult = ref(null)

onMounted(async () => {
  try {
    const res = await scriptApi.get(id)
    script.value = res.data
    yamlContent.value = res.data.yaml_content || ''
  } catch (e) {
    if (e.code !== 404) console.error(e)
  } finally {
    loading.value = false
  }
})

async function handleSave() {
  saving.value = true
  try {
    await scriptApi.update(id, { yaml_content: yamlContent.value })
    alert('保存成功')
  } catch (e) {
    alert(e.message || '保存失败')
  } finally {
    saving.value = false
  }
}

async function handleValidate() {
  validating.value = true
  validationResult.value = null
  try {
    const res = await scriptApi.validate(id)
    validationResult.value = res.data
  } catch (e) {
    alert(e.message || '校验失败')
  } finally {
    validating.value = false
  }
}

async function handleCheckHallucination() {
  checking.value = true
  try {
    const res = await scriptApi.checkHallucination(id)
    if (script.value) script.value.hallucination_risk = res.data.hallucination_risk
    alert(`幻觉风险: ${res.data.hallucination_risk}`)
  } catch (e) {
    alert(e.message || '检测失败')
  } finally {
    checking.value = false
  }
}

async function handleCheckSafety() {
  checking.value = true
  try {
    const res = await scriptApi.checkSafety(id)
    if (script.value) script.value.safety_risk = res.data.safety_risk
    alert(`安全风险: ${res.data.safety_risk}`)
  } catch (e) {
    alert(e.message || '检测失败')
  } finally {
    checking.value = false
  }
}

function exportUrl(format) {
  return scriptApi.exportUrl(id, format)
}
</script>

<style scoped>
.yaml-editor {
  width: 100%;
  min-height: 500px;
  padding: 16px;
  border: 1px solid var(--gray-200);
  border-radius: 6px;
  font-family: 'Menlo', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  resize: vertical;
  tab-size: 2;
}
.yaml-editor:focus { outline: none; border-color: var(--primary); }
.issue-item { display: flex; align-items: center; gap: 8px; padding: 8px 0; border-bottom: 1px solid var(--gray-100); }
</style>
