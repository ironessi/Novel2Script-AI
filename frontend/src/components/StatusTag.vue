<template>
  <n-tag :type="type" :bordered="false" size="small">
    {{ label }}
  </n-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: string
}>()

const type = computed(() => {
  const map: Record<string, 'success' | 'warning' | 'error' | 'info' | 'default'> = {
    created: 'default',
    uploaded: 'info',
    processing: 'warning',
    completed: 'success',
    has_risk: 'warning',
    valid: 'success',
    invalid: 'error',
    pending: 'default',
    running: 'warning',
    failed: 'error',
    cancelled: 'default'
  }
  return map[props.status] || 'default'
})

const label = computed(() => {
  const map: Record<string, string> = {
    created: '已创建',
    uploaded: '已上传',
    processing: '处理中',
    completed: '已完成',
    has_risk: '有风险',
    valid: '通过',
    invalid: '不通过',
    pending: '等待中',
    running: '运行中',
    failed: '失败',
    cancelled: '已取消',
    screen_script: '影视剧本',
    stage_play: '舞台剧',
    short_video: '短视频分镜',
    radio_drama: '广播剧'
  }
  return map[props.status] || props.status
})
</script>
