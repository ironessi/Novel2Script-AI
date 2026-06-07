<template>
  <div>
    <h2 class="mb-4">上传小说</h2>
    <div class="card">
      <div class="upload-area" @dragover.prevent @drop.prevent="onDrop" @click="$refs.input.click()">
        <input ref="input" type="file" accept=".txt,.md,.docx" style="display:none" @change="onFileChange" />
        <p v-if="!file">点击或拖拽文件到此处<br><span class="text-sm text-gray">支持 .txt / .md / .docx，最大 20MB</span></p>
        <p v-else>{{ file.name }} ({{ (file.size / 1024).toFixed(1) }} KB)</p>
      </div>

      <p v-if="error" class="error mt-4">{{ error }}</p>

      <button class="btn btn-primary mt-4" @click="handleUpload" :disabled="!file || uploading">
        {{ uploading ? `上传中... ${progress}%` : '上传' }}
      </button>
    </div>

    <!-- 上传结果 -->
    <div v-if="result" class="card mt-4">
      <h3>上传成功</h3>
      <p class="mt-4">识别到 <strong>{{ result.chapter_count }}</strong> 个章节：</p>
      <table>
        <thead><tr><th>序号</th><th>标题</th></tr></thead>
        <tbody>
          <tr v-for="ch in result.chapters" :key="ch.chapter_index">
            <td>{{ ch.chapter_index }}</td>
            <td>{{ ch.chapter_title }}</td>
          </tr>
        </tbody>
      </table>
      <router-link :to="`/projects/${$route.params.id}`" class="btn btn-outline mt-4">返回项目</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { uploadApi } from '../api'

const route = useRoute()
const file = ref(null)
const uploading = ref(false)
const progress = ref(0)
const error = ref('')
const result = ref(null)

function onFileChange(e) { file.value = e.target.files[0] }
function onDrop(e) { file.value = e.dataTransfer.files[0] }

async function handleUpload() {
  if (!file.value) return
  uploading.value = true
  error.value = ''
  progress.value = 0

  const formData = new FormData()
  formData.append('file', file.value)

  try {
    const res = await uploadApi.upload(route.params.id, formData)
    result.value = res.data
  } catch (e) {
    error.value = e.message || '上传失败'
  } finally {
    uploading.value = false
  }
}
</script>

<style scoped>
.upload-area {
  border: 2px dashed var(--gray-300);
  border-radius: 8px;
  padding: 48px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s;
}
.upload-area:hover { border-color: var(--primary); }
.error { color: var(--danger); font-size: 13px; }
</style>
