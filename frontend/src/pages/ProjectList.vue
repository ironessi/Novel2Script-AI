<template>
  <div>
    <div class="flex justify-between items-center mb-4">
      <h2>我的项目</h2>
      <button class="btn btn-primary" @click="showCreate = true">+ 新建项目</button>
    </div>

    <!-- 新建项目表单 -->
    <div v-if="showCreate" class="card mb-4">
      <form @submit.prevent="handleCreate">
        <div class="form-group">
          <label>项目标题</label>
          <input v-model="newProject.title" required placeholder="输入小说名称" />
        </div>
        <div class="form-group">
          <label>描述</label>
          <textarea v-model="newProject.description" rows="2" placeholder="选填"></textarea>
        </div>
        <div class="form-group">
          <label>改编模式</label>
          <select v-model="newProject.adaptation_mode">
            <option value="screen_script">影视剧本</option>
            <option value="stage_play">舞台剧</option>
            <option value="short_video">短视频分镜</option>
            <option value="radio_drama">广播剧</option>
          </select>
        </div>
        <div class="flex gap-2">
          <button type="submit" class="btn btn-primary" :disabled="creating">创建</button>
          <button type="button" class="btn btn-outline" @click="showCreate = false">取消</button>
        </div>
      </form>
    </div>

    <!-- 项目列表 -->
    <div v-if="loading" class="text-center text-gray">加载中...</div>
    <div v-else-if="projects.length === 0" class="text-center text-gray">暂无项目，点击右上角创建</div>
    <div v-else class="project-grid">
      <div v-for="p in projects" :key="p.id" class="card project-card" @click="$router.push(`/projects/${p.id}`)">
        <h3>{{ p.title }}</h3>
        <p class="text-sm text-gray">{{ p.description || '暂无描述' }}</p>
        <div class="flex justify-between items-center mt-4">
          <span class="badge" :class="`badge-${p.status}`">{{ statusText(p.status) }}</span>
          <span class="text-sm text-gray">{{ p.created_at }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { projectApi } from '../api'

const projects = ref([])
const loading = ref(true)
const showCreate = ref(false)
const creating = ref(false)
const newProject = ref({ title: '', description: '', adaptation_mode: 'screen_script' })

onMounted(loadProjects)

async function loadProjects() {
  loading.value = true
  try {
    const res = await projectApi.list({ page: 1, page_size: 50 })
    projects.value = res.data.projects || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  creating.value = true
  try {
    await projectApi.create(newProject.value)
    showCreate.value = false
    newProject.value = { title: '', description: '', adaptation_mode: 'screen_script' }
    await loadProjects()
  } catch (e) {
    alert(e.message || '创建失败')
  } finally {
    creating.value = false
  }
}

function statusText(s) {
  const map = { created: '已创建', uploaded: '已上传', processing: '处理中', completed: '已完成', archived: '已归档' }
  return map[s] || s
}
</script>

<style scoped>
.project-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 16px; }
.project-card { cursor: pointer; transition: box-shadow 0.2s; }
.project-card:hover { box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
h3 { margin-bottom: 8px; }
</style>
