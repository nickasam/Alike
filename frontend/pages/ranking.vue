<script setup lang="ts">
/**
 * 排行榜 — 三个榜单 Tab：
 *   最受共情榜 GET /api/ranking/empathy（帖子榜，条目为消息）
 *   最暖牛马榜 GET /api/ranking/warmest（用户榜，指标=被共情数）
 *   连续打卡榜 GET /api/ranking/streak（用户榜，指标=连续天数）
 * 展示：前三名领奖台 + 完整榜单前 10 名 + 我的排名（仅用户榜，在已加载列表中查找）。
 * 切换 Tab 时惰性加载对应数据并缓存。
 */
import { useEmotions } from '~/composables/useEmotions'
import { useAuthStore } from '~/stores/auth'

useHead({ title: '排行榜 · Alike' })

const api = useApi()
const auth = useAuthStore()
const { find: findEmotion } = useEmotions()

type TabKey = 'empathy' | 'warmest' | 'streak'

const tabs: { key: TabKey; label: string; path: string; unit: string; icon: string }[] = [
  { key: 'empathy', label: '最受共情榜', path: '/ranking/empathy', unit: '共情', icon: 'heart-handshake' },
  { key: 'warmest', label: '最暖牛马榜', path: '/ranking/warmest', unit: '被共情', icon: 'trophy' },
  { key: 'streak', label: '连续打卡榜', path: '/ranking/streak', unit: '连续打卡', icon: 'sparkles' },
]

interface RankMessage {
  message_id: number
  channel_id: number
  channel_name: string
  content: string
  emotion?: string
  is_anonymous: boolean
  empathy_count: number
  author?: { id: number; nickname: string; avatar_url: string; level: number }
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

/** 归一化后的榜单行，供领奖台 / 列表 / 我的排名统一渲染。 */
interface RankRow {
  key: string | number
  name: string
  sub: string
  emotion?: string
  avatarUrl: string
  avatarChar: string
  level: number | null
  metric: number
  unit: string
  link: string
  /** 关联用户 id（用户榜有；帖子榜匿名则无），用于「我的排名」匹配 */
  userId: number | null
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
  if (key === 'warmest' || key === 'streak') loadMyRank(key)
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

/** 当前 Tab 单位文案。 */
const activeUnit = computed(() => tabs.find((t) => t.key === active.value)!.unit)

/** 将各 Tab 的缓存归一化为统一行结构。 */
const rows = computed<RankRow[]>(() => {
  if (active.value === 'empathy') {
    return (cache.empathy ?? []).map((m) => ({
      key: m.message_id,
      name: m.is_anonymous ? '匿名牛马' : (m.author?.nickname ?? '牛马'),
      sub: m.channel_name ? `#${m.channel_name} · ${m.content}` : m.content,
      emotion: m.emotion,
      avatarUrl: m.is_anonymous ? '' : (m.author?.avatar_url ?? ''),
      avatarChar: avatarChar(m.author?.nickname ?? '', m.is_anonymous),
      level: m.is_anonymous ? null : (m.author?.level ?? null),
      metric: m.empathy_count,
      unit: '共情',
      link: `/channel/${m.channel_id}`,
      userId: m.is_anonymous ? null : (m.author?.id ?? null),
    }))
  }
  if (active.value === 'warmest') {
    return (cache.warmest ?? []).map((u) => ({
      key: u.user_id,
      name: u.nickname,
      sub: `Lv.${u.level} 牛马`,
      avatarUrl: u.avatar_url,
      avatarChar: avatarChar(u.nickname),
      level: u.level,
      metric: u.metric,
      unit: '被共情',
      link: `/profile/${u.user_id}`,
      userId: u.user_id,
    }))
  }
  return (cache.streak ?? []).map((u) => ({
    key: u.user_id,
    name: u.nickname,
    sub: `Lv.${u.level} 牛马`,
    avatarUrl: u.avatar_url,
    avatarChar: avatarChar(u.nickname),
    level: u.level,
    metric: u.days,
    unit: '连续打卡',
    link: `/profile/${u.user_id}`,
    userId: u.user_id,
  }))
})

/** 前三名（领奖台，按 亚/冠/季 顺序在模板里排布）。 */
const podium = computed(() => rows.value.slice(0, 3))
/** 第 4-10 名（完整榜单）。 */
const restRows = computed(() => rows.value.slice(3, 10))

/** 我的排名：仅用户榜、已登录时，向后端拉取精确名次（不受前 50 名限制）。
 *  rank>0 上榜；rank=0 未上榜（指标为 0）。帖子榜/未登录为 null。 */
interface MyRankResp { rank: number; metric: number }
const myRankData = reactive<{ warmest: MyRankResp | null; streak: MyRankResp | null }>({
  warmest: null,
  streak: null,
})

async function loadMyRank(key: 'warmest' | 'streak') {
  if (!auth.user || myRankData[key]) return
  try {
    myRankData[key] = await api.get<MyRankResp>(`/ranking/${key}/me`)
  } catch {
    // 拉取失败静默：不显示我的排名卡片
  }
}

const myRank = computed<{ rank: number; metric: number } | null>(() => {
  if (active.value === 'empathy' || !auth.user) return null
  return myRankData[active.value]
})

/** 领奖台名次样式（金/银/铜）。数组顺序对应 podium 的 [冠,亚,季]。 */
const podiumStyle = [
  // 冠军：金 order-2 居中拔高
  { order: 'order-2 md:scale-105', medal: 'text-[#3d2c00]', medalBg: 'linear-gradient(135deg,#fcd34d,#f59e0b)', val: 'text-gold',
    cardBg: 'linear-gradient(180deg,rgba(251,191,36,.16),rgba(26,34,54,.6))', border: 'rgba(251,191,36,.4)' },
  // 亚军：银 order-1 居左
  { order: 'order-1', medal: 'text-[#2a2f3a]', medalBg: 'linear-gradient(135deg,#e2e8f0,#94a3b8)', val: 'text-[#cbd5e1]',
    cardBg: 'linear-gradient(180deg,rgba(203,213,225,.16),rgba(26,34,54,.6))', border: 'rgba(203,213,225,.35)' },
  // 季军：铜 order-3 居右
  { order: 'order-3', medal: 'text-[#3a2410]', medalBg: 'linear-gradient(135deg,#f0a875,#c2703c)', val: 'text-[#e0955f]',
    cardBg: 'linear-gradient(180deg,rgba(224,149,95,.16),rgba(26,34,54,.6))', border: 'rgba(224,149,95,.35)' },
]
function podStyleOf(i: number) {
  return podiumStyle[i] ?? podiumStyle[2]
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
        <AppIcon :name="t.icon" :size="16" />
        {{ t.label }}
      </button>
    </div>

    <p v-if="loading" class="glass-card py-10 text-center text-sm text-mute">加载中…</p>
    <p v-else-if="error" class="glass-card py-10 text-center text-sm text-danger">{{ error }}</p>
    <p v-else-if="!rows.length" class="glass-card py-10 text-center text-sm text-mute">
      还没有上榜的{{ active === 'empathy' ? '心声' : '牛马' }}。
    </p>

    <template v-else>
      <!-- 领奖台：前三名（顺序 亚-冠-季，冠军居中拔高） -->
      <div class="flex items-end justify-center gap-3 md:gap-4">
        <NuxtLink
          v-for="(r, i) in podium"
          :key="r.key"
          :to="r.link"
          class="glass-card relative flex flex-1 flex-col items-center px-3 pb-5 pt-6 text-center"
          :class="podStyleOf(i).order"
          :style="{ maxWidth: '200px', background: podStyleOf(i).cardBg, borderColor: podStyleOf(i).border }"
        >
          <!-- 冠军皇冠 -->
          <AppIcon
            v-if="i === 0"
            name="trophy"
            :size="26"
            class="absolute -top-3 left-1/2 -translate-x-1/2 text-gold drop-shadow-[0_2px_8px_rgba(251,191,36,.5)]"
          />
          <div class="relative">
            <div class="grid h-16 w-16 place-items-center overflow-hidden rounded-xl bg-grad-ai text-xl font-extrabold text-white shadow-md">
              <img v-if="r.avatarUrl" :src="r.avatarUrl" :alt="r.name" class="h-full w-full object-cover" />
              <template v-else>{{ r.avatarChar }}</template>
            </div>
            <span
              class="absolute -bottom-1.5 -right-1.5 grid h-7 w-7 place-items-center rounded-full text-xs font-extrabold ring-2 ring-surface-solid"
              :class="podStyleOf(i).medal"
              :style="{ background: podStyleOf(i).medalBg }"
            >{{ i + 1 }}</span>
          </div>
          <p class="mt-3 truncate text-sm font-extrabold text-text">{{ r.name }}</p>
          <span
            v-if="r.level !== null"
            class="mt-1.5 rounded-full px-2.5 py-0.5 text-xs font-bold"
            :class="levelClass(r.level)"
          >{{ levelLabel(r.level) }}</span>
          <p class="mt-2.5 text-2xl font-extrabold leading-none" :class="podStyleOf(i).val">{{ r.metric }}</p>
          <p class="mt-1 text-xs text-mute">{{ r.unit }}</p>
        </NuxtLink>
      </div>

      <!-- 完整榜单：第 4-10 名 -->
      <div class="glass-card overflow-hidden">
        <div class="flex items-center gap-2 border-b border-border px-5 py-3.5 text-sm font-bold">
          <AppIcon name="trophy" :size="16" class="text-ai-2" />
          完整榜单 · {{ tabs.find((t) => t.key === active)!.label }}
        </div>
        <p v-if="!restRows.length" class="py-8 text-center text-xs text-mute">
          仅有前三名上榜。
        </p>
        <NuxtLink
          v-for="(r, i) in restRows"
          :key="r.key"
          :to="r.link"
          class="flex items-center gap-3.5 border-b border-border px-5 py-3.5 transition last:border-b-0 hover:bg-surface-hover"
        >
          <span class="w-8 shrink-0 text-center text-md font-extrabold text-mute">{{ i + 4 }}</span>
          <div class="grid h-11 w-11 shrink-0 place-items-center overflow-hidden rounded-md bg-grad-ai text-md font-bold text-white">
            <img v-if="r.avatarUrl" :src="r.avatarUrl" :alt="r.name" class="h-full w-full object-cover" />
            <template v-else>{{ r.avatarChar }}</template>
          </div>
          <div class="min-w-0 flex-1">
            <div class="flex flex-wrap items-center gap-2">
              <span class="truncate text-sm font-bold text-text">{{ r.name }}</span>
              <span
                v-if="r.level !== null"
                class="rounded-full px-2 py-0.5 text-xs font-bold"
                :class="levelClass(r.level)"
              >{{ levelLabel(r.level) }}</span>
              <span
                v-if="findEmotion(r.emotion)"
                class="rounded-full px-2 py-0.5 text-xs font-medium"
                :style="{ background: findEmotion(r.emotion)!.bg, color: findEmotion(r.emotion)!.color }"
              >{{ findEmotion(r.emotion)!.label }}</span>
            </div>
            <p class="mt-0.5 truncate text-xs text-mute">{{ r.sub }}</p>
          </div>
          <span class="flex shrink-0 flex-col items-end">
            <span class="text-lg font-extrabold" :class="active === 'streak' ? 'text-warm' : 'text-empathy'">{{ r.metric }}</span>
            <span class="text-xs text-mute">{{ r.unit }}</span>
          </span>
        </NuxtLink>
      </div>

      <!-- 我的排名（仅用户榜、已登录时展示） -->
      <div
        v-if="myRank"
        class="flex items-center gap-3.5 rounded-lg border border-dashed border-ai-1 px-5 py-4"
        style="background:linear-gradient(135deg,rgba(99,102,241,.12),var(--glass-bg))"
      >
        <span class="w-10 shrink-0 text-center text-md font-extrabold text-ai-1">
          {{ myRank.rank > 0 ? `#${myRank.rank}` : '—' }}
        </span>
        <div class="grid h-11 w-11 shrink-0 place-items-center overflow-hidden rounded-md bg-grad-warm text-md font-bold text-white">
          <img v-if="auth.user?.avatar_url" :src="auth.user.avatar_url" :alt="auth.user?.nickname" class="h-full w-full object-cover" />
          <template v-else>{{ avatarChar(auth.user?.nickname ?? '') }}</template>
        </div>
        <div class="min-w-0 flex-1">
          <div class="flex flex-wrap items-center gap-2">
            <span class="truncate text-sm font-bold text-text">{{ auth.user?.nickname }}</span>
            <span class="rounded-full bg-grad-ai px-2 py-0.5 text-xs font-extrabold text-white">就是你</span>
          </div>
          <p class="mt-0.5 text-xs text-mute">
            {{ myRank.rank > 0 ? '继续加油，稳住名次！' : '还没上榜，多互动就能冲榜啦' }}
          </p>
        </div>
        <span v-if="myRank.rank > 0" class="flex shrink-0 flex-col items-end">
          <span class="text-lg font-extrabold" :class="active === 'streak' ? 'text-warm' : 'text-empathy'">{{ myRank.metric }}</span>
          <span class="text-xs text-mute">{{ activeUnit }}</span>
        </span>
      </div>
    </template>
  </div>
</template>
