<script setup lang="ts">
/**
 * EmotionBoard — 情绪看板（占位骨架）。
 * 阶段五接入实时情绪热力图、时间范围 Tab、WebSocket emotion_update。
 */
defineProps<{ compact?: boolean }>()
</script>

<template>
  <section class="glass-card p-4" aria-label="情绪看板">
    <h3 class="mb-2 flex items-center gap-2 text-md font-semibold">
      <AppIcon name="sparkles" :size="18" />
      今日情绪看板
    </h3>
    <p class="text-sm text-mute">情绪看板占位（阶段五实现）</p>
  </section>
</template>
