<template>
  <AppLayout>
    <div class="dashboard">
      <div class="dashboard-header">
        <h2>我的项目</h2>
        <n-button type="primary" @click="showCreate = true">新建改编项目</n-button>
      </div>

      <!-- 新建项目表单 -->
      <n-card v-if="showCreate" title="新建项目" style="margin-bottom: 16px;">
        <n-form label-placement="left" label-width="80">
          <n-form-item label="小说名称">
            <n-input v-model:value="newProject.title" placeholder="输入小说名称" />
          </n-form-item>
          <n-form-item label="描述">
            <n-input v-model:value="newProject.description" type="textarea" rows="2" placeholder="选填" />
          </n-form-item>
          <n-form-item label="改编模式">
            <n-radio-group v-model:value="newProject.adaptation_mode">
              <n-space>
                <n-radio v-for="(m, key) in adaptationModes" :key="key" :value="key">
                  {{ m.label }}
                </n-radio>
              </n-space>
            </n-radio-group>
          </n-form-item>
          <n-space>
            <n-button type="primary" @click="handleCreate">创建</n-button>
            <n-button @click="showCreate = false">取消</n-button>
          </n-space>
        </n-form>
      </n-card>

      <!-- 项目卡片列表 -->
      <div class="project-grid">
        <div
          v-for="p in store.projects"
          :key="p.id"
          class="project-card"
          @click="openProject(p.id)"
        >
          <div class="card-header">
            <h3>{{ p.title }}</h3>
            <StatusTag :status="p.status" />
          </div>
          <p class="card-desc">{{ p.description }}</p>
          <div class="card-meta">
            <span>{{ adaptationModes[p.adaptation_mode]?.label }}</span>
            <span>{{ p.chapter_count }} 章</span>
            <span>{{ p.updated_at }}</span>
          </div>
          <div class="card-actions" @click.stop>
            <n-button size="small" @click="openProject(p.id)">进入工作台</n-button>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import { adaptationModes } from '@/mock/data'
import AppLayout from '@/layouts/AppLayout.vue'
import StatusTag from '@/components/StatusTag.vue'

const router = useRouter()
const store = useProjectStore()
const showCreate = ref(false)
const newProject = ref({ title: '', description: '', adaptation_mode: 'screen_script' })

function openProject(id: number) {
  store.setCurrentProject(id)
  router.push(`/projects/${id}/workbench`)
}

function handleCreate() {
  store.projects.push({
    id: Date.now(),
    title: newProject.value.title,
    description: newProject.value.description,
    adaptation_mode: newProject.value.adaptation_mode as any,
    status: 'created',
    chapter_count: 0,
    updated_at: '刚刚',
    created_at: '刚刚'
  })
  showCreate.value = false
  newProject.value = { title: '', description: '', adaptation_mode: 'screen_script' }
}
</script>

<style scoped>
.dashboard { max-width: 960px; }
.dashboard-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.dashboard-header h2 { font-size: 18px; font-weight: 600; }
.project-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 16px; }
.project-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 20px;
  cursor: pointer;
  transition: border-color 0.15s;
}
.project-card:hover { border-color: var(--color-primary); }
.card-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.card-header h3 { font-size: 15px; font-weight: 600; }
.card-desc { font-size: 13px; color: var(--color-text-secondary); margin-bottom: 12px; line-height: 1.5; }
.card-meta { display: flex; gap: 12px; font-size: 12px; color: var(--color-text-secondary); margin-bottom: 12px; }
.card-actions { display: flex; justify-content: flex-end; }
</style>
