<template>
  <AppLayout>
    <div class="upload-page">
      <h2>上传小说稿件</h2>

      <div class="upload-area" @dragover.prevent @drop.prevent="onDrop" @click="fileInput?.click()">
        <input ref="fileInput" type="file" accept=".txt,.md,.docx" hidden @change="onFileChange" />
        <div v-if="!file" class="upload-placeholder">
          <n-icon size="48" color="#d1d5db"><CloudUploadOutline /></n-icon>
          <p class="upload-text">点击或拖拽文件到此处</p>
          <p class="upload-hint">支持 .txt / .md / .docx，最大 20MB</p>
        </div>
        <div v-else class="upload-file">
          <n-icon size="24" color="#4f6ef7"><DocumentTextOutline /></n-icon>
          <span>{{ file.name }}</span>
          <span class="file-size">{{ (file.size / 1024).toFixed(1) }} KB</span>
        </div>
      </div>

      <!-- 改编模式选择 -->
      <div class="mode-section">
        <h3>选择改编模式</h3>
        <div class="mode-grid">
          <div
            v-for="(m, key) in adaptationModes"
            :key="key"
            class="mode-card"
            :class="{ selected: selectedMode === key }"
            @click="selectedMode = key"
          >
            <div class="mode-label">{{ m.label }}</div>
            <div class="mode-desc">{{ m.desc }}</div>
          </div>
        </div>
      </div>

      <!-- 上传按钮 -->
      <div v-if="file" class="upload-actions">
        <n-button type="primary" :loading="uploading" @click="handleUpload">
          {{ uploading ? '上传中...' : '上传并识别章节' }}
        </n-button>
      </div>

      <!-- 章节预览 -->
      <div v-if="chapters.length > 0" class="preview-section">
        <h3>章节识别预览</h3>
        <n-data-table :columns="columns" :data="chapters" :bordered="false" size="small" />
        <div class="preview-actions">
          <n-button type="primary" @click="handleConfirm">确认，开始 AI 分析</n-button>
        </div>
      </div>

      <n-alert v-if="error" type="error" style="margin-top: 16px;" closable @close="error = ''">
        {{ error }}
      </n-alert>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import { uploadApi } from '@/api'
import { adaptationModes } from '@/mock/data'
import AppLayout from '@/layouts/AppLayout.vue'
import { CloudUploadOutline, DocumentTextOutline } from '@vicons/ionicons5'
import type { DataTableColumns } from 'naive-ui'

const route = useRoute()
const router = useRouter()
const store = useProjectStore()
const projectId = Number(route.params.id)

const fileInput = ref<HTMLInputElement | null>(null)
const file = ref<File | null>(null)
const selectedMode = ref('screen_script')
const uploading = ref(false)
const error = ref('')
const chapters = ref<any[]>([])

const columns: DataTableColumns<any> = [
  { title: '序号', key: 'chapter_index', width: 60 },
  { title: '章节标题', key: 'chapter_title' },
  { title: '哈希', key: 'content_hash', width: 200, render: (row) => row.content_hash?.slice(0, 16) + '...' }
]

onMounted(() => {
  store.setCurrentProject(projectId)
})

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  if (input.files?.[0]) file.value = input.files[0]
}

function onDrop(e: DragEvent) {
  if (e.dataTransfer?.files?.[0]) file.value = e.dataTransfer.files[0]
}

async function handleUpload() {
  if (!file.value) return
  uploading.value = true
  error.value = ''

  const formData = new FormData()
  formData.append('file', file.value)

  try {
    const res: any = await uploadApi.upload(projectId, formData)
    if (res.code === 0) {
      chapters.value = res.data.chapters || []
    } else {
      error.value = res.message || '上传失败'
    }
  } catch (e: any) {
    error.value = e.message || '上传失败'
  } finally {
    uploading.value = false
  }
}

function handleConfirm() {
  router.push(`/projects/${projectId}/workflow`)
}
</script>

<style scoped>
.upload-page { max-width: 720px; }
.upload-page h2 { font-size: 18px; font-weight: 600; margin-bottom: 20px; }
.upload-area {
  border: 2px dashed var(--color-border);
  border-radius: 8px;
  padding: 48px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.15s;
  background: var(--color-surface);
}
.upload-area:hover { border-color: var(--color-primary); }
.upload-text { font-size: 15px; margin-top: 12px; }
.upload-hint { font-size: 12px; color: var(--color-text-secondary); margin-top: 4px; }
.upload-file { display: flex; align-items: center; gap: 8px; justify-content: center; }
.file-size { color: var(--color-text-secondary); font-size: 12px; }
.upload-actions { margin-top: 16px; display: flex; justify-content: center; }
.mode-section { margin-top: 24px; }
.mode-section h3 { font-size: 15px; font-weight: 600; margin-bottom: 12px; }
.mode-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 12px; }
.mode-card {
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s;
  background: var(--color-surface);
}
.mode-card:hover { border-color: var(--color-primary); }
.mode-card.selected { border-color: var(--color-primary); background: var(--color-primary-light); }
.mode-label { font-size: 14px; font-weight: 500; margin-bottom: 4px; }
.mode-desc { font-size: 12px; color: var(--color-text-secondary); line-height: 1.5; }
.preview-section { margin-top: 24px; }
.preview-section h3 { font-size: 15px; font-weight: 600; margin-bottom: 12px; }
.preview-actions { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
