<script setup lang="ts">
/**
 * 排行榜 — 三个榜单 Tab：
 *   最受共情榜 GET /api/ranking/empathy（帖子榜，条目为消息）
 *   最暖牛马榜 GET /api/ranking/warmest（用户榜，指标=被共情数）
 *   连续打卡榜 GET /api/ranking/streak（用户榜，指标=连续天数）
 * 切换 Tab 时惰性加载对应数据并缓存，Top 20 列表展示排名/头像/昵称/数值。
 */
import { useEmotions } from '~/composables/useEmotions'

useHead({ title: '排行榜 · Alike' })

const api = useApi()
const { find: findEmotion } = useEmotions()

type TabKey = 'empathy' | 'warmest' | 'streak'

const tabs: { key: TabKey; label: string; path: string; unit: string }[] = [
  { key: 'empathy', label: '最受共情榜', path: '/ranking/empathy', unit: '共情' },
  { key: 'warmest', label: '最暖牛马榜', path: '/ranking/warmest', unit: '被共情' },
  { key: 'streak', label: '连续打卡榜', path: '/ranking/streak', unit: '天' },
]

interface RankMessage {
  message_id: number
  channel_id: number
  content: string
  emotion?: string
  is_anonymous: boolean
  empathy_count: number
  author?: { id: number; nickname: string; avatar_url: string }
}
interface RankUser {
  user_id: number
  nickname: string
  avatar_url: string
  level: number
  metric: number
}
interface RankStreak {
  user_id: number
  nickname: string
  avatar_url: string
  level: number
  days: number
}

const active = ref<TabKey>('empathy')
const loading = ref(false)
const error = ref('')

const cache = reactive<{
  empathy: RankMessage[] | null
  warmest: RankUser[] | null
  streak: RankStreak[] | null
}>({ empathy: null, warmest: null, streak: null })

async function load(key: TabKey) {
  if (cache[key]) return
  const tab = tabs.find((t) => t.key === key)!
  loading.value = true
  error.value = ''
  try {
    const res = await api.get<{ list: unknown[] }>(tab.path, { limit: 20 })
    cache[key] = (res.list ?? []) as never
  } catch {
    error.value = '榜单加载失败，稍后再试'
  } finally {
    loading.value = false
  }
}

function switchTab(key: TabKey) {
  active.value = key
  load(key)
}

/** 排名奖牌色：前三名金/银/铜，其余中性。 */
function rankClass(i: number): string {
  if (i === 0) return 'bg-grad-warm text-[#1a0f00]'
  if (i === 1) return 'bg-surface-hover text-text'
  if (i === 2) return 'bg-empathy-soft text-empathy'
  return 'bg-surface text-mute'
}

function avatarChar(name: string, anon = false): string {
  if (anon) return '匿'
  return name?.charAt(0) ?? '牛'
}

onMounted(() => load('empathy'))
</script>

<template>
  <div class="flex flex-col gap-5">
    <header class="glass-card animate-rise-in p-6">
      <h1 class="text-gradient flex items-center gap-2 text-2xl font-extrabold">
        <AppIcon name="trophy" :size="24" />
        牛马排行榜
      </h1>
      <p class="mt-1 text-sm text-dim">看看谁最懂大家的辛苦，谁在坚持打卡。</p>
    </header>

    <!-- Tab 切换 -->
    <div class="flex gap-2" role="tablist">
      <button
        v-for="t in tabs"
        :key="t.key"
        type="button"
        role="tab"
        :aria-selected="active === t.key"
        class="rounded-full px-4 py-2 text-sm font-medium transition"
        :class="active === t.key
          ? 'btn-primary'
          : 'border border-border-strong text-dim hover:text-text'"
        @click="switchTab(t.key)"
      >
        {{ t.label }}
      </button>
    </div>

    <div class="glass-card p-4">
      <p v-if="loading" class="py-8 text-center text-sm text-mute">加载中…</p>
      <p v-else-if="error" class="py-8 text-center text-sm text-danger">{{ error }}</p>

      <!-- 最受共情榜（帖子） -->
      <ul v-else-if="active === 'empathy'" class="flex flex-col gap-2">
        <li v-if="!cache.empathy?.length" class="py-8 text-center text-sm text-mute">
          还没有上榜的心声。
        </li>
        <NuxtLink
          v-for="(m, i) in cache.empathy ?? []"
          :key="m.message_id"
          :to="`/channel/${m.channel_id}`"
          class="flex items-start gap-3 rounded-md p-3 transition hover:bg-surface-hover"
        >
          <span
            class="grid h-8 w-8 shrink-0 place-items-center rounded-full text-sm font-bold"
            :class="rankClass(i)"
          >{{ i + 1 }}</span>
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2">
              <span class="text-sm font-semibold text-text">
                {{ m.is_anonymous ? '匿名牛马' : (m.author?.nickname ?? '牛马') }}
              </span>
              <span
                v-if="findEmotion(m.emotion)"
                class="rounded-full px-2 py-0.5 text-xs font-medium"
                :style="{ background: findEmotion(m.emotion)!.bg, color: findEmotion(m.emotion)!.color }"
              >{{ findEmotion(m.emotion)!.label }}</span>
            </div>
            <p class="mt-0.5 truncate text-sm text-dim">{{ m.content }}</p>
          </div>
          <span class="flex shrink-0 items-center gap-1 text-sm font-semibold text-empathy">
            <AppIcon name="heart-handshake" :size="14" />
            {{ m.empathy_count }}
          </span>
        </NuxtLink>
      </ul>

      <!-- 最暖牛马榜 / 连续打卡榜（用户） -->
      <ul v-else class="flex flex-col gap-2">
        <li
          v-if="active === 'warmest' ? !cache.warmest?.length : !cache.streak?.length"
          class="py-8 text-center text-sm text-mute"
        >
          还没有上榜的牛马。
        </li>
        <template v-if="active === 'warmest'">
          <NuxtLink
            v-for="(u, i) in cache.warmest ?? []"
            :key="u.user_id"
            :to="`/profile/${u.user_id}`"
            class="flex items-center gap-3 rounded-md p-3 transition hover:bg-surface-hover"
          >
            <span
              class="grid h-8 w-8 shrink-0 place-items-center rounded-full text-sm font-bold"
              :class="rankClass(i)"
            >{{ i + 1 }}</span>
            <div class="grid h-10 w-10 shrink-0 place-items-center rounded-md bg-grad-ai text-sm font-semibold text-white">
              {{ avatarChar(u.nickname) }}
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-semibold text-text">{{ u.nickname }}</p>
              <p class="text-xs text-mute">Lv.{{ u.level }} 牛马</p>
            </div>
            <span class="shrink-0 text-sm font-semibold text-empathy">
              {{ u.metric }} 被共情
            </span>
          </NuxtLink>
        </template>
        <template v-else>
          <NuxtLink
            v-for="(u, i) in cache.streak ?? []"
            :key="u.user_id"
            :to="`/profile/${u.user_id}`"
            class="flex items-center gap-3 rounded-md p-3 transition hover:bg-surface-hover"
          >
            <span
              class="grid h-8 w-8 shrink-0 place-items-center rounded-full text-sm font-bold"
              :class="rankClass(i)"
            >{{ i + 1 }}</span>
            <div class="grid h-10 w-10 shrink-0 place-items-center rounded-md bg-grad-ai text-sm font-semibold text-white">
              {{ avatarChar(u.nickname) }}
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-semibold text-text">{{ u.nickname }}</p>
              <p class="text-xs text-mute">Lv.{{ u.level }} 牛马</p>
            </div>
            <span class="shrink-0 text-sm font-semibold text-warm">
              连续 {{ u.days }} 天
            </span>
          </NuxtLink>
        </template>
      </ul>
    </div>
  </div>
</template>
