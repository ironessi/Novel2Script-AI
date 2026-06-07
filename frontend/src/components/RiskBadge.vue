<template>
  <n-tag :type="type" :bordered="false" size="small" round>
    {{ label }}
  </n-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  level: 'none' | 'low' | 'medium' | 'high'
}>()

const type = computed(() => {
  const map = { none: 'success', low: 'success', medium: 'warning', high: 'error' } as const
  return map[props.level] || 'default'
})

const label = computed(() => {
  const map = { none: '正常', low: '低风险', medium: '中风险', high: '高风险' }
  return map[props.level] || props.level
})
</script>
