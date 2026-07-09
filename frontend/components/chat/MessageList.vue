<script setup lang="ts">
/**
 * MessageList — 消息流（占位骨架）。阶段四接入分页加载与实时消息。
 */
defineProps<{ channelId?: number }>()
</script>

<template>
  <div class="flex flex-col gap-3" aria-label="消息列表">
    <p class="text-sm text-mute">消息列表占位（阶段四实现）</p>
  </div>
</template>
