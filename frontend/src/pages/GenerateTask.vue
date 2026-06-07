<template>
  <div>
    <h2 class="mb-4">AI 生成剧本</h2>
    <div class="card">
      <p>此页面用于查看 AI 生成任务的进度。</p>
      <p class="text-sm text-gray mt-4">请从项目详情页发起生成任务。</p>
      <router-link :to="`/projects/${$route.params.id}`" class="btn btn-outline mt-4">返回项目</router-link>
    </div>
  </div>
</template>
