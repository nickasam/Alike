<script setup lang="ts">
/**
 * EmpathyButton — 共情按钮（核心特色"抱团取暖"）。
 * 两态：
 *   - 未共情：柔和描边胶囊 + 空心手心，hover 微亮；
 *   - 已共情：empathy 柔色底 + 实色边框文字 + 轻微发光，图标高亮，点击有回弹动效。
 * 匿名/隐私不受影响。点击根据当前态 emit add/remove，由父级发起请求。
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
    class="empathy-btn"
    :class="[
      empathized ? 'empathy-btn--on' : 'empathy-btn--off',
      size === 'sm' ? 'empathy-btn--sm' : '',
    ]"
    :aria-pressed="empathized"
    :aria-label="`我懂你，当前 ${count} 人共情`"
    :disabled="disabled"
    @click="onClick"
  >
    <AppIcon
      name="heart-handshake"
      :size="size === 'sm' ? 14 : 16"
      class="empathy-btn__icon"
    />
    <span>{{ empathized ? '已懂你' : '我懂你' }}</span>
    <span v-if="count > 0" class="empathy-btn__count">{{ count }}</span>
  </button>
</template>

<style scoped>
.empathy-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  border-radius: 9999px;
  padding: 0.3rem 0.85rem;
  font-size: 0.875rem;
  font-weight: 500;
  line-height: 1;
  border: 1px solid transparent;
  cursor: pointer;
  transition:
    background 0.2s ease,
    border-color 0.2s ease,
    color 0.2s ease,
    box-shadow 0.2s ease,
    transform 0.12s ease;
}
.empathy-btn--sm {
  padding: 0.25rem 0.7rem;
  font-size: 0.8125rem;
}
.empathy-btn:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}
.empathy-btn:active {
  transform: scale(0.94);
}

/* 未共情：柔和描边，hover 微微点亮为 empathy 色 */
.empathy-btn--off {
  background: transparent;
  border-color: var(--border-strong, rgba(148, 163, 184, 0.35));
  color: var(--text-dim, #94a3b8);
}
.empathy-btn--off:hover:not(:disabled) {
  border-color: var(--empathy);
  color: var(--empathy);
  background: var(--empathy-soft);
}

/* 已共情：empathy 柔色底 + 实色边框/文字 + 轻发光，清爽高对比 */
.empathy-btn--on {
  background: var(--empathy-soft);
  border-color: var(--empathy);
  color: var(--empathy);
  box-shadow: var(--glow-empathy);
  animation: empathy-pop 0.28s ease;
}
.empathy-btn--on:hover:not(:disabled) {
  background: color-mix(in srgb, var(--empathy) 18%, transparent);
}
.empathy-btn--on .empathy-btn__icon {
  transform: scale(1.08);
}

/* 计数徽标 */
.empathy-btn__count {
  font-weight: 700;
  font-variant-numeric: tabular-nums;
}
.empathy-btn__icon {
  transition: transform 0.2s ease;
}

@keyframes empathy-pop {
  0% {
    transform: scale(1);
  }
  45% {
    transform: scale(1.12);
  }
  100% {
    transform: scale(1);
  }
}
</style>
