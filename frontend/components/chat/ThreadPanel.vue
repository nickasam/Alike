<script setup lang="ts">
/**
 * ThreadPanel — 线程回复面板（占位骨架）。
 * 阶段四接入线程回复流；断点适配：桌面内联 / 平板覆盖 / 移动全屏。
 */
defineProps<{ visible?: boolean }>()
defineEmits<{ (e: 'close'): void }>()
</script>

<template>
  <aside v-if="visible" class="glass-card p-4" role="dialog" aria-modal="true">
    <p class="text-sm text-mute">线程面板占位（阶段四实现）</p>
  </aside>
</template>
