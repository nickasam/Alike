<script setup lang="ts">
/**
 * ChannelSidebar — 频道侧边栏（占位骨架）。
 * 按分类（行业/岗位/主题/自建）分组展示频道，底部创建频道入口。
 * 阶段三接入真实频道数据与搜索/折叠/未读逻辑。
 */
interface NavLink {
  to: string
  label: string
  icon: string
}

const nav: NavLink[] = [
  { to: '/', label: '首页', icon: 'home' },
  { to: '/diary', label: '日记广场', icon: 'book-open' },
  { to: '/ranking', label: '排行榜', icon: 'trophy' },
]

// 占位频道数据（阶段三由 channel store 提供）
const groups: { title: string; channels: { id: number; name: string }[] }[] = [
  {
    title: '行业',
    channels: [
      { id: 1, name: '互联网大厂' },
      { id: 2, name: '制造业' },
    ],
  },
  {
    title: '岗位',
    channels: [
      { id: 3, name: '程序员' },
      { id: 4, name: '设计师' },
    ],
  },
  {
    title: '主题',
    channels: [
      { id: 5, name: '摸鱼交流' },
      { id: 6, name: '离职天堂' },
    ],
  },
]
</script>

<template>
  <nav class="flex h-full flex-col gap-5 p-4" aria-label="频道导航">
    <!-- 主导航 -->
    <ul class="flex flex-col gap-1">
      <li v-for="link in nav" :key="link.to">
        <NuxtLink
          :to="link.to"
          class="flex items-center gap-3 rounded-md px-3 py-2 text-dim transition hover:bg-surface-hover hover:text-text"
          active-class="bg-surface-hover text-text"
        >
          <AppIcon :name="link.icon" :size="18" />
          <span class="text-base">{{ link.label }}</span>
        </NuxtLink>
      </li>
    </ul>

    <!-- 频道分组 -->
    <div
      v-for="group in groups"
      :key="group.title"
      class="flex flex-col gap-1"
    >
      <p class="px-3 text-xs font-semibold uppercase tracking-wide text-mute">
        {{ group.title }}
      </p>
      <NuxtLink
        v-for="ch in group.channels"
        :key="ch.id"
        :to="`/channel/${ch.id}`"
        class="flex items-center gap-3 rounded-md px-3 py-2 text-dim transition hover:bg-surface-hover hover:text-text"
        active-class="bg-surface-hover text-text"
      >
        <AppIcon name="hash" :size="16" />
        <span class="truncate text-base">{{ ch.name }}</span>
      </NuxtLink>
    </div>

    <!-- 创建频道 -->
    <button
      class="mt-auto flex items-center justify-center gap-2 rounded-md border border-dashed border-border-strong px-3 py-2 text-dim transition hover:text-ai-1"
    >
      <AppIcon name="plus" :size="18" />
      <span class="text-base">创建频道</span>
    </button>
  </nav>
</template>
