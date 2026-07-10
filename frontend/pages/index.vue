<script setup lang="ts">
/**
 * 首页 — 热门频道 + 今日情绪看板（全站）+ 今日牛马榜（最暖）。
 * 热门频道从 /api/channels 拉真实数据；情绪看板拉 /api/emotion/board（全站今日聚合）；
 * 牛马榜拉 /api/ranking/warmest（按给出共情数降序）。
 */
import { useChannelStore, type Channel } from '~/stores/channel'
import { useEmotions } from '~/composables/useEmotions'

useHead({ title: 'Alike · 汇聚天下牛马' })

const api = useApi()
const channelStore = useChannelStore()
const { find: findEmotion } = useEmotions()

// —— 热门频道 ——
const hotChannels = ref<Channel[]>([])
const loading = ref(true)

// —— 今日情绪看板（全站）——
interface EmotionCount { emotion: string; count: number }
interface Board { scope?: string; emotions: EmotionCount[]; total: number; dominant?: string }
const board = ref<Board | null>(null)
const boardLoading = ref(true)

// —— 今日牛马榜（最暖）——
interface RankUser { user_id: number; nickname: string; avatar_url: string; level: number; metric: number }
const warmest = ref<RankUser[]>([])
const rankLoading = ref(true)

onMounted(() => {
  loadHotChannels()
  loadBoard()
  loadWarmest()
})

async function loadHotChannels() {
  try {
    const res = await api.get<{ list: Channel[] }>('/channels?page=1&page_size=6')
    const list = res?.list ?? []
    hotChannels.value = list
    channelStore.setHotChannels(list)
  } catch {
    hotChannels.value = []
  } finally {
    loading.value = false
  }
}

async function loadBoard() {
  try {
    board.value = await api.get<Board>('/emotion/board')
  } catch {
    board.value = null
  } finally {
    boardLoading.value = false
  }
}

async function loadWarmest() {
  try {
    const res = await api.get<{ list: RankUser[] }>('/ranking/warmest?limit=10')
    warmest.value = res?.list ?? []
  } catch {
    warmest.value = []
  } finally {
    rankLoading.value = false
  }
}

/** 频道名可能已含 # 前缀，去掉后由模板统一加 #，避免 ## 重复。 */
function displayName(name: string): string {
  return name.replace(/^#/, '')
}

/** 「进入频道吐槽」目标：优先热门首个频道，其次侧边栏全量首个频道。 */
const firstChannelId = computed<number | null>(
  () => hotChannels.value[0]?.id ?? channelStore.channels[0]?.id ?? null,
)

function enterChannel() {
  if (firstChannelId.value != null) {
    navigateTo(`/channel/${firstChannelId.value}`)
  } else if (import.meta.client) {
    document.getElementById('hot-channels')?.scrollIntoView({ behavior: 'smooth' })
  }
}

/** 情绪看板行数据：8 种情绪 + 百分比 + 是否主导。 */
const emotionRows = computed(() => {
  const total = board.value?.total ?? 0
  const dominant = board.value?.dominant ?? ''
  const counts = new Map((board.value?.emotions ?? []).map((e) => [e.emotion, e.count]))
  const order = ['tired', 'angry', 'wronged', 'break', 'numb', 'quit', 'anxious', 'cheer']
  return order.map((key) => {
    const meta = findEmotion(key)
    const count = counts.get(key) ?? 0
    return {
      key,
      label: meta?.label ?? key,
      bg: meta?.bg ?? 'rgba(148,163,184,.14)',
      color: meta?.color ?? '#94a3b8',
      count,
      percent: total > 0 ? Math.round((count / total) * 100) : 0,
      isDominant: key === dominant && count > 0,
    }
  })
})
const boardTotal = computed(() => board.value?.total ?? 0)
const dominantLabel = computed(() => findEmotion(board.value?.dominant)?.label ?? '')
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Hero -->
    <section class="glass-card animate-rise-in p-8">
      <h1 class="text-gradient text-2xl font-extrabold leading-tight md:text-3xl">
        汇聚天下牛马
      </h1>
      <p class="mt-2 text-md text-dim">
        总有人懂你的辛苦，你不是一个人在扛。
      </p>
      <div class="mt-5 flex flex-wrap gap-3">
        <button
          class="btn-primary px-5 py-2 text-base font-semibold"
          @click="enterChannel"
        >
          进入频道吐槽
        </button>
        <NuxtLink
          to="/ranking"
          class="rounded-md border border-border-strong px-5 py-2 text-base text-dim transition hover:text-text"
        >
          今日牛马榜
        </NuxtLink>
      </div>
    </section>

    <!-- 今日情绪看板（全站） + 今日牛马榜（最暖），桌面并排 -->
    <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <!-- 今日情绪看板 -->
      <section class="glass-card p-5">
        <h2 class="mb-3 flex items-center gap-2 text-lg font-semibold">
          <AppIcon name="sparkles" :size="20" />
          今日情绪看板
          <span v-if="boardTotal" class="ml-auto text-xs font-normal text-mute">
            {{ boardTotal }} 条心声
          </span>
        </h2>
        <p v-if="boardLoading" class="text-sm text-mute">加载中…</p>
        <p v-else-if="boardTotal === 0" class="text-sm text-mute">
          今天还没有人吐露情绪，来发第一条吧。
        </p>
        <template v-else>
          <p v-if="dominantLabel" class="mb-3 text-sm text-dim">
            今天大家最多的情绪是
            <span class="font-semibold text-text">「{{ dominantLabel }}」</span>，抱抱。
          </p>
          <ul class="flex flex-col gap-2">
            <li v-for="row in emotionRows" :key="row.key" class="flex items-center gap-3">
              <span
                class="w-12 shrink-0 rounded-full px-2 py-0.5 text-center text-xs font-medium"
                :class="{ 'ring-2 ring-offset-0': row.isDominant }"
                :style="{ background: row.bg, color: row.color }"
              >
                {{ row.label }}
              </span>
              <div class="h-2 flex-1 overflow-hidden rounded-full bg-surface-hover">
                <div
                  class="h-full rounded-full transition-all duration-500 ease-out"
                  :style="{ width: `${row.percent}%`, background: row.color }"
                />
              </div>
              <span class="w-14 shrink-0 text-right text-xs text-mute">
                {{ row.count }} · {{ row.percent }}%
              </span>
            </li>
          </ul>
        </template>
      </section>

      <!-- 今日牛马榜（最暖） -->
      <section class="glass-card p-5">
        <h2 class="mb-3 flex items-center gap-2 text-lg font-semibold">
          <AppIcon name="trophy" :size="20" />
          今日牛马榜
          <span class="ml-auto text-xs font-normal text-mute">最暖牛马</span>
        </h2>
        <p v-if="rankLoading" class="text-sm text-mute">加载中…</p>
        <p v-else-if="warmest.length === 0" class="text-sm text-mute">
          还没有人给出共情，第一个抱抱别人的就是你。
        </p>
        <ol v-else class="flex flex-col gap-2">
          <li v-for="(u, i) in warmest" :key="u.user_id">
            <NuxtLink
              :to="`/profile/${u.user_id}`"
              class="flex items-center gap-3 rounded-md px-2 py-1.5 transition hover:bg-surface-hover"
            >
              <span
                class="w-5 shrink-0 text-center text-sm font-bold"
                :class="i < 3 ? 'text-warm' : 'text-mute'"
              >
                {{ i + 1 }}
              </span>
              <img
                v-if="u.avatar_url"
                :src="u.avatar_url"
                :alt="u.nickname"
                class="h-8 w-8 rounded-md object-cover"
              />
              <span
                v-else
                class="grid h-8 w-8 shrink-0 place-items-center rounded-md bg-grad-ai text-xs font-semibold text-white"
              >
                {{ u.nickname?.[0] ?? '牛' }}
              </span>
              <span class="min-w-0 flex-1 truncate text-sm font-medium text-text">
                {{ u.nickname }}
              </span>
              <span class="shrink-0 text-xs text-empathy">给出 {{ u.metric }} 次共情</span>
            </NuxtLink>
          </li>
        </ol>
      </section>
    </div>

    <!-- 热门频道 -->
    <section id="hot-channels">
      <h2 class="mb-3 flex items-center gap-2 text-lg font-semibold">
        <AppIcon name="hash" :size="20" />
        热门频道
      </h2>
      <p v-if="loading" class="glass-card p-5 text-sm text-mute">加载中…</p>
      <p v-else-if="hotChannels.length === 0" class="glass-card p-5 text-sm text-mute">
        暂无频道
      </p>
      <div v-else class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <NuxtLink
          v-for="ch in hotChannels"
          :key="ch.id"
          :to="`/channel/${ch.id}`"
          class="glass-card block p-5"
        >
          <div class="flex items-center justify-between">
            <h3 class="text-md font-semibold">#{{ displayName(ch.name) }}</h3>
            <span class="text-xs text-mute">{{ ch.member_count }} 位牛马</span>
          </div>
          <p v-if="ch.description" class="mt-1 text-sm text-dim">{{ ch.description }}</p>
        </NuxtLink>
      </div>
    </section>
  </div>
</template>
