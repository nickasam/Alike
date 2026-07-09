<script setup lang="ts">
/** 个人主页（占位骨架）。阶段二/八接入用户信息与统计。 */
const route = useRoute()
useHead({ title: '个人主页 · Alike' })
</script>

<template>
  <div class="flex flex-col gap-4">
    <h1 class="text-xl font-semibold">牛马主页 #{{ route.params.id }}</h1>
    <div class="glass-card p-6">
      <p class="text-sm text-mute">个人主页占位（阶段二实现）</p>
    </div>
  </div>
</template>
