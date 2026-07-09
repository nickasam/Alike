<script setup lang="ts">
/**
 * 频道页（占位骨架）。阶段四接入消息流 + 情绪看板 + 线程面板。
 */
const route = useRoute()
const id = computed(() => route.params.id)
useHead({ title: `频道 · Alike` })
</script>

<template>
  <div class="flex flex-col gap-4">
    <h1 class="text-xl font-semibold">频道 #{{ id }}</h1>
    <MessageList :channel-id="Number(id)" />
    <MessageInput />
  </div>
</template>
