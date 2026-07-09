<script setup lang="ts">
/** 日记详情（占位骨架，v1.5）。 */
const route = useRoute()
useHead({ title: '日记详情 · Alike' })
</script>

<template>
  <article class="mx-auto flex max-w-content flex-col gap-4">
    <h1 class="text-xl font-semibold">日记 #{{ route.params.id }}</h1>
    <div class="glass-card p-6">
      <p class="text-sm text-mute">日记详情占位（v1.5 实现）</p>
    </div>
  </article>
</template>
