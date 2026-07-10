<script setup lang="ts">
/**
 * ChannelSidebar — 频道侧边栏。
 * 按分类（行业/岗位/主题/自建）分组展示从后端拉取的真实频道，
 * 保证侧边栏点击的 id 与后端频道一致（避免文字与实际频道对不上）。
 *
 * 数据来源：始终拉取全量频道写入 channelStore.channels（与首页热门 hotChannels 互不覆盖），
 * 保证无论从首页还是直连频道页进入，侧边栏都显示完整频道列表。
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

/** 前端过滤关键词（纯本地 filter，不调后端）。 */
const filter = ref('')

/** 分类展示顺序与标题。 */
const categoryOrder: { key: ChannelCategory; title: string }[] = [
  { key: 'industry', title: '行业' },
  { key: 'job', title: '岗位' },
  { key: 'topic', title: '主题' },
  { key: 'custom', title: '自建' },
]

/** 折叠状态：默认全部展开。 */
const collapsed = reactive<Record<ChannelCategory, boolean>>({
  industry: false,
  job: false,
  topic: false,
  custom: false,
})

function toggle(key: ChannelCategory) {
  collapsed[key] = !collapsed[key]
}

/** 频道名可能已含 # 前缀，去掉后模板统一处理。 */
function displayName(name: string): string {
  return name.replace(/^#/, '')
}

const groups = computed(() => {
  const kw = filter.value.trim().toLowerCase()
  const match = (c: Channel) => {
    if (!kw) return true
    const hay = `${displayName(c.name)} ${c.slug ?? ''} ${c.description ?? ''}`.toLowerCase()
    return hay.includes(kw)
  }
  return categoryOrder
    .map(({ key, title }) => ({
      key,
      title,
      channels: channelStore.byCategory(key).filter(match),
    }))
    .filter((g) => g.channels.length > 0)
})

onMounted(async () => {
  // 始终拉取全量，保证侧边栏拥有完整频道列表（不受首页热门 6 条影响）。
  try {
    const res = await api.get<{ list: Channel[] }>('/channels?page=1&page_size=50')
    channelStore.setChannels(res?.list ?? [])
  } catch {
    // 拉取失败时侧边栏仅显示主导航
  }
})
</script>

<template>
  <nav class="flex h-full min-h-0 flex-col gap-5 p-4" aria-label="频道导航">
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

    <!-- 频道过滤输入 -->
    <div class="relative">
      <span class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-mute">
        <AppIcon name="search" :size="16" />
      </span>
      <input
        v-model="filter"
        type="text"
        placeholder="过滤频道…"
        aria-label="过滤频道"
        class="w-full rounded-md border border-border-strong bg-surface py-2 pl-9 pr-3 text-base text-text placeholder:text-mute transition focus:border-ai-1 focus:outline-none"
      />
    </div>

    <!-- 频道分组（可折叠 + 溢出滚动） -->
    <div class="flex min-h-0 flex-1 flex-col gap-3 overflow-y-auto">
      <div
        v-for="group in groups"
        :key="group.key"
        class="flex flex-col gap-1"
      >
        <button
          type="button"
          class="flex items-center justify-between rounded-md px-3 py-1 text-xs font-semibold uppercase tracking-wide text-mute transition hover:text-text"
          :aria-expanded="!collapsed[group.key]"
          @click="toggle(group.key)"
        >
          <span>{{ group.title }}</span>
          <AppIcon
            name="chevron-down"
            :size="14"
            :class="['transition-transform', collapsed[group.key] ? '-rotate-90' : '']"
          />
        </button>
        <template v-if="!collapsed[group.key]">
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
        </template>
      </div>
      <p
        v-if="groups.length === 0 && filter.trim()"
        class="px-3 py-1 text-xs text-mute"
      >
        没有匹配的频道
      </p>
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
