<script setup lang="ts">
/**
 * EmotionPicker — 情绪选择器（占位骨架）。
 * 阶段五接入 8 种情绪单选、键盘导航、清除逻辑。
 */
defineEmits<{ (e: 'select', payload: { emotion: string | null }): void }>()
</script>

<template>
  <div class="glass-card p-3" role="listbox" aria-label="情绪选择器">
    <p class="text-sm text-mute">情绪选择器占位（阶段五实现）</p>
  </div>
</template>
