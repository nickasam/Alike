<script setup lang="ts">
/**
 * ChannelSidebar — 频道侧边栏。
 * 只展示「主题频道」（category=topic，如吐槽大会/打工日记/薪资揭秘…约 6 个），
 * 保持列表精简；其余分类（行业/岗位/自建）不在侧边栏铺开。
 *
 * 数据来源：始终拉取全量频道写入 channelStore.channels（与首页热门 hotChannels 互不覆盖），
 * 保证无论从首页还是直连频道页进入，侧边栏都显示完整主题频道列表。
 */
import { useChannelStore, type Channel } from '~/stores/channel'

interface NavLink {
  to: string
  label: string
  icon: string
  pulse?: boolean // 呼吸绿点：暗示"正在跳动"
}

const nav: NavLink[] = [
  { to: '/', label: '首页', icon: 'home' },
  { to: '/diary', label: '日记广场', icon: 'book-open' },
  { to: '/pulse', label: '最近发生', icon: 'activity', pulse: true },
  { to: '/ranking', label: '排行榜', icon: 'trophy' },
]

const api = useApi()
const channelStore = useChannelStore()

/** 前端过滤关键词（纯本地 filter，不调后端）。 */
const filter = ref('')

/** 频道名可能已含 # 前缀，去掉后模板统一处理。 */
function displayName(name: string): string {
  return name.replace(/^#/, '')
}

/** 只取主题频道，并按关键词本地过滤。 */
const topicChannels = computed(() => {
  const kw = filter.value.trim().toLowerCase()
  return channelStore.byCategory('topic').filter((c: Channel) => {
    if (!kw) return true
    const hay = `${displayName(c.name)} ${c.slug ?? ''} ${c.description ?? ''}`.toLowerCase()
    return hay.includes(kw)
  })
})

onMounted(async () => {
  // 始终拉取全量，保证侧边栏拥有完整主题频道列表（不受首页热门 6 条影响）。
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
          <span
            v-if="link.pulse"
            class="ml-auto h-1.5 w-1.5 rounded-full bg-empathy pulse-dot"
            aria-hidden="true"
          />
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

    <!-- 主题频道（溢出滚动） -->
    <div class="flex min-h-0 flex-1 flex-col gap-1 overflow-y-auto">
      <div class="px-3 py-1 text-xs font-semibold uppercase tracking-wide text-mute">
        主题频道
      </div>
      <NuxtLink
        v-for="ch in topicChannels"
        :key="ch.id"
        :to="`/channel/${ch.id}`"
        class="flex items-center gap-3 rounded-md px-3 py-2 text-dim transition hover:bg-surface-hover hover:text-text"
        active-class="bg-surface-hover text-text"
      >
        <AppIcon name="hash" :size="16" />
        <span class="truncate text-base">{{ displayName(ch.name) }}</span>
      </NuxtLink>
      <p
        v-if="topicChannels.length === 0"
        class="px-3 py-1 text-xs text-mute"
      >
        {{ filter.trim() ? '没有匹配的频道' : '暂无主题频道' }}
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

<style scoped>
/*「最近发生」入口的呼吸绿点：暗示这里正在跳动。 */
.pulse-dot {
  box-shadow: 0 0 8px var(--empathy);
  animation: pulseDot 1.6s ease-in-out infinite;
}
@keyframes pulseDot {
  0%, 100% { opacity: .4; transform: scale(.9); }
  50%      { opacity: 1;  transform: scale(1.15); }
}
@media (prefers-reduced-motion: reduce) {
  .pulse-dot { animation: none; opacity: .9; }
}
</style>
