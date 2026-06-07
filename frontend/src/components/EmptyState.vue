<template>
  <div class="empty-state">
    <div class="empty-icon">
      <n-icon size="48" :component="icon" />
    </div>
    <p class="empty-title">{{ title }}</p>
    <p class="empty-desc">{{ description }}</p>
    <slot />
  </div>
</template>

<script setup lang="ts">
import { DocumentTextOutline } from '@vicons/ionicons5'
import { type Component } from 'vue'

withDefaults(defineProps<{
  title: string
  description: string
  icon?: Component
}>(), {
  icon: DocumentTextOutline
})
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  text-align: center;
}
.empty-icon {
  color: #d1d5db;
  margin-bottom: 16px;
}
.empty-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--color-text);
  margin-bottom: 8px;
}
.empty-desc {
  font-size: 13px;
  color: var(--color-text-secondary);
  max-width: 320px;
  line-height: 1.6;
}
</style>
