<script setup lang="ts">
/**
 * AppIcon — 统一 SVG 线性图标（Lucide 风格，stroke-width 1.5）。
 * 禁用 emoji 作为 UI 元素。按 name 渲染对应 path。
 */
const props = withDefaults(
  defineProps<{ name: string; size?: number }>(),
  { size: 20 },
)

// 核心图标 path 集合（d 属性），viewBox 统一 24x24。
const paths: Record<string, string[]> = {
  search: ['M11 11m-8 0a8 8 0 1 0 16 0a8 8 0 1 0 -16 0', 'M21 21l-4.35-4.35'],
  bell: [
    'M18 8a6 6 0 0 0-12 0c0 7-3 9-3 9h18s-3-2-3-9',
    'M13.73 21a2 2 0 0 1-3.46 0',
  ],
  home: ['M3 9.5 12 3l9 6.5', 'M5 10v10h14V10'],
  hash: ['M4 9h16', 'M4 15h16', 'M10 3 8 21', 'M16 3l-2 18'],
  'book-open': [
    'M12 7v14',
    'M3 5h6a2 2 0 0 1 2 2v14a2.5 2.5 0 0 0-2-1H3z',
    'M21 5h-6a2 2 0 0 0-2 2v14a2.5 2.5 0 0 1 2-1h6z',
  ],
  trophy: [
    'M6 4h12v4a6 6 0 0 1-12 0z',
    'M6 6H4a2 2 0 0 0 0 4h2',
    'M18 6h2a2 2 0 0 1 0 4h-2',
    'M12 14v4',
    'M8 21h8',
    'M10 18h4',
  ],
  user: ['M12 12m-4 0a4 4 0 1 0 8 0a4 4 0 1 0 -8 0', 'M4 21c0-4 4-6 8-6s8 2 8 6'],
  'heart-handshake': [
    'M12 20 4 12a4 4 0 0 1 6-6l2 2 2-2a4 4 0 0 1 6 6z',
  ],
  sparkles: [
    'M12 3l1.8 4.2L18 9l-4.2 1.8L12 15l-1.8-4.2L6 9l4.2-1.8z',
    'M19 15l.9 2.1L22 18l-2.1.9L19 21l-.9-2.1L16 18l2.1-.9z',
  ],
  sun: [
    'M12 12m-4 0a4 4 0 1 0 8 0a4 4 0 1 0 -8 0',
    'M12 2v2 M12 20v2 M4.9 4.9l1.4 1.4 M17.7 17.7l1.4 1.4 M2 12h2 M20 12h2 M4.9 19.1l1.4-1.4 M17.7 6.3l1.4-1.4',
  ],
  moon: ['M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8z'],
  menu: ['M4 6h16', 'M4 12h16', 'M4 18h16'],
  plus: ['M12 5v14', 'M5 12h14'],
  'log-out': [
    'M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4',
    'M16 17l5-5-5-5',
    'M21 12H9',
  ],
}

const d = computed(() => paths[props.name] ?? [])
</script>

<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.5"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
    class="shrink-0"
  >
    <path v-for="(p, i) in d" :key="i" :d="p" />
  </svg>
</template>
