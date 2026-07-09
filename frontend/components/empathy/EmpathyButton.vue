<script setup lang="ts">
/**
 * EmpathyButton — 共情按钮（占位骨架，核心特色）。
 * 阶段六接入乐观更新、粒子动画、WebSocket empathy 实时同步。
 */
withDefaults(
  defineProps<{
    count?: number
    empathized?: boolean
    size?: 'sm' | 'md' | 'lg'
  }>(),
  { count: 0, empathized: false, size: 'md' },
)
defineEmits<{
  (e: 'empathy', payload: { action: 'add' | 'remove' }): void
}>()
</script>

<template>
  <button
    class="inline-flex items-center gap-2 rounded-full border border-empathy px-3 py-1 text-empathy transition hover:shadow-glow-empathy"
    :class="{ 'bg-grad-empathy text-[#052e20]': empathized }"
    :aria-pressed="empathized"
    :aria-label="`我懂你，当前 ${count} 人共情`"
  >
    <AppIcon name="heart-handshake" :size="16" />
    <span class="text-sm">{{ empathized ? '已懂你' : '我懂你' }}</span>
    <span class="text-sm font-semibold">{{ count }}</span>
  </button>
</template>
