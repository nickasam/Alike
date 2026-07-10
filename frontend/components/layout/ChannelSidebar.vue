<script setup lang="ts">
/**
 * ChannelSidebar — 频道侧边栏。
 * 按分类（行业/岗位/主题/自建）分组展示从后端拉取的真实频道，
 * 保证侧边栏点击的 id 与后端频道一致（避免文字与实际频道对不上）。
 */
import { useChannelStore, type Channel, type ChannelCategory } from '~/stores/channel'

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

const api = useApi()
const channelStore = useChannelStore()

// 分类展示顺序与标题。
const categoryOrder: { key: ChannelCategory; title: string }[] = [
  { key: 'industry', title: '行业' },
  { key: 'job', title: '岗位' },
  { key: 'topic', title: '主题' },
  { key: 'custom', title: '自建' },
]

const groups = computed(() =>
  categoryOrder
    .map(({ key, title }) => ({ title, channels: channelStore.byCategory(key) }))
    .filter((g) => g.channels.length > 0),
)

/** 频道名可能已含 # 前缀，去掉后模板统一处理。 */
function displayName(name: string): string {
  return name.replace(/^#/, '')
}

onMounted(async () => {
  // 已有数据则复用，避免重复拉取。
  if (channelStore.channels.length > 0) return
  try {
    const res = await api.get<{ list: Channel[] }>('/channels?page=1&page_size=50')
    channelStore.setChannels(res?.list ?? [])
  } catch {
    // 拉取失败时侧边栏仅显示主导航
  }
})
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
        <span class="truncate text-base">{{ displayName(ch.name) }}</span>
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
