<script setup lang="ts">
/**
 * MessageInput — 消息输入框（占位骨架）。
 * 阶段四接入表情 / 情绪标签 / 匿名开关 / 发送态 / typing 广播。
 */
</script>

<template>
  <div class="glass-card flex items-center gap-2 p-3">
    <input
      type="text"
      aria-label="发送消息"
      placeholder="说点什么，总有人懂你的辛苦..."
      class="h-10 flex-1 rounded-md border border-border bg-surface px-3 text-base text-text outline-none placeholder:text-mute focus:border-ai-1"
    />
    <button class="btn-primary px-4 py-2 text-base">发送</button>
  </div>
</template>
