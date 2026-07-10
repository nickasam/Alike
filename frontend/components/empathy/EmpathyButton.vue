<script setup lang="ts">
/**
 * EmpathyButton — 共情按钮（核心特色"抱团取暖"）。
 * 点击根据当前是否已共情，emit add/remove，由父级发起 REST 请求并做乐观更新，
 * WebSocket empathy 事件负责跨端同步计数。
 */
const props = withDefaults(
  defineProps<{
    count?: number
    empathized?: boolean
    size?: 'sm' | 'md' | 'lg'
    disabled?: boolean
  }>(),
  { count: 0, empathized: false, size: 'md', disabled: false },
)
const emit = defineEmits<{
  (e: 'empathy', payload: { action: 'add' | 'remove' }): void
}>()

function onClick() {
  if (props.disabled) return
  emit('empathy', { action: props.empathized ? 'remove' : 'add' })
}
</script>

<template>
  <button
    type="button"
    class="inline-flex items-center gap-2 rounded-full border border-empathy px-3 py-1 text-empathy transition hover:shadow-glow-empathy disabled:cursor-not-allowed disabled:opacity-60"
    :class="{ 'bg-grad-empathy text-[#052e20]': empathized }"
    :aria-pressed="empathized"
    :aria-label="`我懂你，当前 ${count} 人共情`"
    :disabled="disabled"
    @click="onClick"
  >
    <AppIcon name="heart-handshake" :size="16" />
    <span class="text-sm">{{ empathized ? '已懂你' : '我懂你' }}</span>
    <span class="text-sm font-semibold">{{ count }}</span>
  </button>
</template>
