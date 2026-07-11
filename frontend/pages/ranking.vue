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

/** 头像首字（昵称首字符），匿名固定「匿」。 */
function avatarChar(name: string, anon = false): string {
  if (anon) return '匿'
  return name?.charAt(0) ?? '牛'
}

/** 牛马等级名（参照原型：牛马之王 / 老牛马 / 小牛马）。 */
function levelLabel(level: number): string {
  if (level >= 5) return '牛马之王'
  if (level >= 3) return '老牛马'
  return '小牛马'
}
/** 等级徽章配色类。 */
function levelClass(level: number): string {
  if (level >= 5) return 'bg-grad-warm text-[#3d2c00]'
  if (level >= 3) return 'bg-warm/15 text-warm'
  return 'bg-empathy-soft text-empathy'
}
/** 每个 Tab 的指标图标名（AppIcon）。 */
function tabIcon(key: TabKey): string {
  if (key === 'empathy') return 'heart-handshake'
  if (key === 'warmest') return 'trophy'
  return 'sparkles'
}

onMounted(() => load('empathy'))
</script>

<template>
  <div class="mx-auto flex max-w-content flex-col gap-6">
    <!-- 页面标题：居中 + 皇冠 + 渐变字 -->
    <header class="animate-rise-in text-center">
      <h1 class="inline-flex items-center gap-3 text-2xl font-extrabold md:text-3xl">
        <AppIcon name="trophy" :size="30" class="text-gold" />
        牛马<span class="text-gradient">排行榜</span>
      </h1>
      <p class="mt-1.5 text-sm text-dim">最努力扛住生活的人，值得被看见。</p>
    </header>

    <!-- Tab 切换：玻璃胶囊 + 图标 + 选中渐变发光 -->
    <div class="flex flex-wrap justify-center gap-2.5" role="tablist">
      <button
        v-for="t in tabs"
        :key="t.key"
        type="button"
        role="tab"
        :aria-selected="active === t.key"
        class="inline-flex items-center gap-2 rounded-full px-5 py-2.5 text-sm font-bold transition duration-std ease-out"
        :class="active === t.key
          ? 'bg-grad-ai text-white shadow-glow-ai'
          : 'glass-card !rounded-full text-dim hover:-translate-y-0.5 hover:text-text'"
        @click="switchTab(t.key)"
      >
        <AppIcon :name="tabIcon(t.key)" :size="16" />
        {{ t.label }}
      </button>
    </div>

    <div class="glass-card overflow-hidden">
      <p v-if="loading" class="py-10 text-center text-sm text-mute">加载中…</p>
      <p v-else-if="error" class="py-10 text-center text-sm text-danger">{{ error }}</p>

      <!-- 最受共情榜（帖子） -->
      <ul v-else-if="active === 'empathy'" class="flex flex-col">
        <li v-if="!cache.empathy?.length" class="py-10 text-center text-sm text-mute">
          还没有上榜的心声。
        </li>
        <NuxtLink
          v-for="(m, i) in cache.empathy ?? []"
          :key="m.message_id"
          :to="`/channel/${m.channel_id}`"
          class="flex items-center gap-3.5 border-b border-border px-5 py-3.5 transition last:border-b-0 hover:bg-surface-hover"
        >
          <span
            class="w-8 shrink-0 text-center text-md font-extrabold"
            :class="i < 3 ? 'text-warm' : 'text-mute'"
          >{{ i + 1 }}</span>
          <div class="min-w-0 flex-1">
            <div class="flex flex-wrap items-center gap-2">
              <span class="text-sm font-bold text-text">
                {{ m.is_anonymous ? '匿名牛马' : (m.author?.nickname ?? '牛马') }}
              </span>
              <span
                v-if="findEmotion(m.emotion)"
                class="rounded-full px-2 py-0.5 text-xs font-medium"
                :style="{ background: findEmotion(m.emotion)!.bg, color: findEmotion(m.emotion)!.color }"
              >{{ findEmotion(m.emotion)!.label }}</span>
            </div>
            <p class="mt-0.5 truncate text-xs text-mute">{{ m.content }}</p>
          </div>
          <span class="flex shrink-0 flex-col items-end">
            <span class="text-lg font-extrabold text-empathy">{{ m.empathy_count }}</span>
            <span class="flex items-center gap-1 text-xs text-mute">
              <AppIcon name="heart-handshake" :size="12" />共情
            </span>
          </span>
        </NuxtLink>
      </ul>

      <!-- 最暖牛马榜 / 连续打卡榜（用户） -->
      <ul v-else class="flex flex-col">
        <li
          v-if="active === 'warmest' ? !cache.warmest?.length : !cache.streak?.length"
          class="py-10 text-center text-sm text-mute"
        >
          还没有上榜的牛马。
        </li>
        <template v-if="active === 'warmest'">
          <NuxtLink
            v-for="(u, i) in cache.warmest ?? []"
            :key="u.user_id"
            :to="`/profile/${u.user_id}`"
            class="flex items-center gap-3.5 border-b border-border px-5 py-3.5 transition last:border-b-0 hover:bg-surface-hover"
          >
            <span
              class="w-8 shrink-0 text-center text-md font-extrabold"
              :class="i < 3 ? 'text-warm' : 'text-mute'"
            >{{ i + 1 }}</span>
            <div class="grid h-11 w-11 shrink-0 place-items-center overflow-hidden rounded-md bg-grad-ai text-md font-bold text-white">
              <img
                v-if="u.avatar_url"
                :src="u.avatar_url"
                :alt="u.nickname"
                class="h-full w-full object-cover"
              />
              <template v-else>{{ avatarChar(u.nickname) }}</template>
            </div>
            <div class="min-w-0 flex-1">
              <div class="flex flex-wrap items-center gap-2">
                <span class="truncate text-sm font-bold text-text">{{ u.nickname }}</span>
                <span class="rounded-full px-2 py-0.5 text-xs font-bold" :class="levelClass(u.level)">
                  {{ levelLabel(u.level) }}
                </span>
              </div>
              <p class="mt-0.5 text-xs text-mute">Lv.{{ u.level }} 牛马</p>
            </div>
            <span class="flex shrink-0 flex-col items-end">
              <span class="text-lg font-extrabold text-empathy">{{ u.metric }}</span>
              <span class="flex items-center gap-1 text-xs text-mute">
                <AppIcon name="heart-handshake" :size="12" />被共情
              </span>
            </span>
          </NuxtLink>
        </template>
        <template v-else>
          <NuxtLink
            v-for="(u, i) in cache.streak ?? []"
            :key="u.user_id"
            :to="`/profile/${u.user_id}`"
            class="flex items-center gap-3.5 border-b border-border px-5 py-3.5 transition last:border-b-0 hover:bg-surface-hover"
          >
            <span
              class="w-8 shrink-0 text-center text-md font-extrabold"
              :class="i < 3 ? 'text-warm' : 'text-mute'"
            >{{ i + 1 }}</span>
            <div class="grid h-11 w-11 shrink-0 place-items-center overflow-hidden rounded-md bg-grad-ai text-md font-bold text-white">
              <img
                v-if="u.avatar_url"
                :src="u.avatar_url"
                :alt="u.nickname"
                class="h-full w-full object-cover"
              />
              <template v-else>{{ avatarChar(u.nickname) }}</template>
            </div>
            <div class="min-w-0 flex-1">
              <div class="flex flex-wrap items-center gap-2">
                <span class="truncate text-sm font-bold text-text">{{ u.nickname }}</span>
                <span class="rounded-full px-2 py-0.5 text-xs font-bold" :class="levelClass(u.level)">
                  {{ levelLabel(u.level) }}
                </span>
              </div>
              <p class="mt-0.5 text-xs text-mute">Lv.{{ u.level }} 牛马</p>
            </div>
            <span class="flex shrink-0 flex-col items-end">
              <span class="text-lg font-extrabold text-warm">{{ u.days }}</span>
              <span class="text-xs text-mute">连续打卡</span>
            </span>
          </NuxtLink>
        </template>
      </ul>
    </div>
  </div>
</template>
