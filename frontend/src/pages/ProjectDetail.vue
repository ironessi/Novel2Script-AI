<template>
  <div v-if="loading" class="text-center text-gray">加载中...</div>
  <div v-else-if="!project" class="text-center text-gray">项目不存在</div>
  <div v-else>
    <!-- 项目信息 -->
    <div class="card mb-4">
      <div class="flex justify-between items-center">
        <div>
          <h2>{{ project.title }}</h2>
          <p class="text-sm text-gray mt-4">{{ project.description || '暂无描述' }}</p>
        </div>
        <span class="badge" :class="`badge-${project.status}`">{{ statusText(project.status) }}</span>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="flex gap-2 mb-4">
      <router-link :to="`/projects/${id}/upload`" class="btn btn-primary">上传小说</router-link>
      <button class="btn btn-success" @click="startGenerate" :disabled="generating || chapters.length === 0">
        {{ generating ? '生成中...' : 'AI 生成剧本' }}
      </button>
      <router-link v-if="hasScript" :to="`/projects/${id}/script`" class="btn btn-outline">查看剧本</router-link>
      <router-link :to="`/projects/${id}/audit`" class="btn btn-outline">审计日志</router-link>
    </div>

    <!-- 任务状态 -->
    <div v-if="task" class="card mb-4">
      <h3>任务状态</h3>
      <div class="mt-4">
        <div class="flex justify-between mb-4">
          <span class="badge" :class="`badge-${task.status}`">{{ task.status }}</span>
          <span class="text-sm">{{ task.progress }}%</span>
        </div>
        <div class="progress-bar">
          <div class="progress-bar-fill" :style="{ width: task.progress + '%' }"></div>
        </div>
        <p v-if="task.current_step" class="text-sm text-gray mt-4">{{ task.current_step }}</p>
        <p v-if="task.error_message" class="text-sm" style="color:var(--danger)">{{ task.error_message }}</p>
      </div>
    </div>

    <!-- 章节列表 -->
    <div class="card">
      <h3>章节列表 ({{ chapters.length }})</h3>
      <table v-if="chapters.length > 0">
        <thead><tr><th>序号</th><th>标题</th><th>哈希</th><th>时间</th></tr></thead>
        <tbody>
          <tr v-for="ch in chapters" :key="ch.id">
            <td>{{ ch.chapter_index }}</td>
            <td>{{ ch.chapter_title }}</td>
            <td class="text-sm text-gray">{{ ch.content_hash?.slice(0, 16) }}...</td>
            <td class="text-sm text-gray">{{ ch.created_at }}</td>
          </tr>
        </tbody>
      </table>
      <p v-else class="text-center text-gray mt-4">暂无章节，请先上传小说</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { projectApi, chapterApi, taskApi, scriptApi } from '../api'

const route = useRoute()
const id = route.params.id

const project = ref(null)
const chapters = ref([])
const task = ref(null)
const loading = ref(true)
const generating = ref(false)
const polling = ref(null)

const hasScript = computed(() => project.value?.status === 'completed')

onMounted(async () => {
  try {
    const [projRes, chRes] = await Promise.all([
      projectApi.get(id),
      chapterApi.list(id),
    ])
    project.value = projRes.data
    chapters.value = chRes.data.chapters || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  if (polling.value) clearInterval(polling.value)
})

async function startGenerate() {
  generating.value = true
  try {
    const res = await taskApi.create(id, { task_type: 'full_generate' })
    task.value = { ...res.data, progress: 0, current_step: '等待中' }
    // 轮询任务状态
    polling.value = setInterval(async () => {
      try {
        const statusRes = await taskApi.status(res.data.task_id)
        task.value = statusRes.data
        if (['completed', 'failed', 'cancelled'].includes(statusRes.data.status)) {
          clearInterval(polling.value)
          generating.value = false
          if (statusRes.data.status === 'completed') {
            project.value.status = 'completed'
          }
        }
      } catch (e) {
        console.error(e)
      }
    }, 3000)
  } catch (e) {
    alert(e.message || '创建任务失败')
    generating.value = false
  }
}

function statusText(s) {
  const map = { created: '已创建', uploaded: '已上传', processing: '处理中', completed: '已完成', archived: '已归档' }
  return map[s] || s
}
</script>
