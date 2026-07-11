<script setup lang="ts">
/**
 * EmotionBoard — 频道「今日情绪看板」。
 *
 * - 拉取 GET /api/channels/:id/emotion-board（默认今日范围），渲染 8 种情绪分布；
 * - 订阅 WebSocket emotion_update 事件，发带情绪消息后局部实时刷新（无需重进频道）；
 * - 高亮占比最高（dominant）情绪，并给出一句话氛围概括；
 * - global 模式拉全站今日聚合 GET /api/emotion/board（首页右侧栏用）；
 *   否则 channelId 未提供时从路由 /channel/:id 推断，非频道页展示空态。
 */
import { useEmotions, type EmotionMeta } from '~/composables/useEmotions'

const props = defineProps<{ channelId?: number; compact?: boolean; global?: boolean }>()

const api = useApi()
const route = useRoute()
const ws = useWebSocket()
const { emotions } = useEmotions()

interface Count {
  emotion: string
  count: number
}
interface Board {
  scope?: string
  emotions: Count[]
  total: number
  dominant?: string
}

/** 优先使用显式 prop，否则从 /channel/:id 路由推断。 */
const channelId = computed<number | null>(() => {
  if (props.channelId && props.channelId > 0) return props.channelId
  if (route.path.startsWith('/channel/')) {
    const id = Number(route.params.id)
    if (Number.isFinite(id) && id > 0) return id
  }
  return null
})

const board = ref<Board | null>(null)
const loading = ref(false)
const error = ref('')
const justUpdated = ref(false)

async function load() {
  const id = channelId.value
  if (!props.global && !id) {
    board.value = null
    return
  }
  loading.value = true
  error.value = ''
  try {
    board.value = props.global
      ? await api.get<Board>('/emotion/board')
      : await api.get<Board>(`/channels/${id}/emotion-board`)
  } catch {
    error.value = '情绪看板加载失败'
    board.value = null
  } finally {
    loading.value = false
  }
}

/** WebSocket 实时更新：仅接受当前频道的 emotion_update，局部刷新并触发高亮动画。 */
let offEmotion: (() => void) | null = null
function subscribe() {
  offEmotion?.()
  // 全站模式不订阅频道级 emotion_update（会用单频道数据覆盖全站聚合）。
  if (props.global) return
  offEmotion = ws.on<Board>('emotion_update', (data, evtChannelId) => {
    if (evtChannelId && channelId.value && evtChannelId !== channelId.value) return
    if (!data || !Array.isArray(data.emotions)) return
    board.value = data
    pulse()
  })
}

/** 更新时短暂点亮标题，给用户"看板刚刷新"的反馈。 */
let pulseTimer: ReturnType<typeof setTimeout> | null = null
function pulse() {
  justUpdated.value = true
  if (pulseTimer) clearTimeout(pulseTimer)
  pulseTimer = setTimeout(() => (justUpdated.value = false), 1200)
}

/** 各情绪计数（含 0），按 useEmotions 顺序稳定渲染，并按计数降序排列突出主导情绪。 */
interface Row extends EmotionMeta {
  count: number
  percent: number
  dominant: boolean
}
const rows = computed<Row[]>(() => {
  const total = board.value?.total ?? 0
  const dominant = board.value?.dominant ?? ''
  const counts = new Map(board.value?.emotions.map((e) => [e.emotion, e.count]))
  return emotions
    .map((meta) => {
      const count = counts.get(meta.key) ?? 0
      return {
        ...meta,
        count,
        percent: total > 0 ? Math.round((count / total) * 100) : 0,
        dominant: meta.key === dominant && count > 0,
      }
    })
    .sort((a, b) => b.count - a.count)
})

const total = computed(() => board.value?.total ?? 0)

/** 主导情绪元数据。 */
const dominantMeta = computed<Row | null>(() => {
  const row = rows.value[0]
  return row && row.dominant ? row : null
})

/** 一句话氛围概括，依据主导情绪与占比生成。 */
const moodSummary = computed<string>(() => {
  const d = dominantMeta.value
  if (!d || total.value === 0) return ''
  const phrases: Record<string, string> = {
    tired: '大家今天都挺累的，抱抱',
    angry: '火气有点大，深呼吸一下',
    wronged: '好多委屈没处说，有人懂你',
    break: '有人快撑不住了，别一个人扛',
    numb: '麻了，但我们还在',
    quit: '想润的人不少，先喝口水',
    anxious: '焦虑在蔓延，慢慢来',
    cheer: '今天氛围不错，继续加油',
  }
  const base = phrases[d.key] ?? ''
  return d.percent >= 50 ? `${base}（${d.label}占了一多半）` : base
})

watch(channelId, () => {
  subscribe()
  load()
}, { immediate: true })

onBeforeUnmount(() => {
  offEmotion?.()
  offEmotion = null
  if (pulseTimer) clearTimeout(pulseTimer)
})
</script>

<template>
  <section class="glass-card p-4" aria-label="今日情绪看板">
    <h3 class="mb-1 flex items-center gap-2 text-md font-semibold">
      <AppIcon
        name="sparkles"
        :size="18"
        :class="justUpdated ? 'text-ai-1 transition-colors' : 'transition-colors'"
      />
      今日情绪看板
      <span
        v-if="justUpdated"
        class="rounded-full bg-ai-1/15 px-1.5 py-0.5 text-[10px] font-normal text-ai-1"
        >刚刚更新</span
      >
      <span v-if="total" class="ml-auto text-xs font-normal text-mute">{{ total }} 条心声</span>
    </h3>

    <!-- 氛围概括 -->
    <p
      v-if="moodSummary"
      class="mb-3 flex items-center gap-1.5 text-xs text-dim"
    >
      <span
        v-if="dominantMeta"
        class="inline-flex items-center gap-1 rounded-full px-1.5 py-0.5 font-medium"
        :style="{ background: dominantMeta.bg, color: dominantMeta.color }"
        ><svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.6"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="h-3 w-3 shrink-0"
          aria-hidden="true"
        ><path v-for="(p, i) in dominantMeta.icon" :key="i" :d="p" /></svg>{{ dominantMeta.label }}</span
      >
      {{ moodSummary }}
    </p>

    <p v-if="!global && !channelId" class="text-sm text-mute">进入频道查看今日情绪分布。</p>
    <p v-else-if="loading && !board" class="text-sm text-mute">加载中…</p>
    <p v-else-if="error" class="text-sm text-danger">{{ error }}</p>
    <p v-else-if="total === 0" class="text-sm text-mute">今天还没有情绪记录，来发第一条吧。</p>

    <ul v-else class="flex flex-col gap-2">
      <li
        v-for="row in rows"
        :key="row.key"
        class="flex items-center gap-3"
        :class="{ 'opacity-45': row.count === 0 }"
      >
        <span
          class="inline-flex w-16 shrink-0 items-center justify-center gap-1 rounded-full px-2 py-0.5 text-center text-xs font-medium ring-offset-0 transition"
          :class="{ 'ring-2': row.dominant }"
          :style="{
            background: row.bg,
            color: row.color,
            boxShadow: row.dominant ? `0 0 0 2px ${row.color}` : 'none',
          }"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.6"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="h-3 w-3 shrink-0"
            aria-hidden="true"
          >
            <path v-for="(p, i) in row.icon" :key="i" :d="p" />
          </svg>
          {{ row.label }}
        </span>
        <div class="h-2 flex-1 overflow-hidden rounded-full bg-surface-hover">
          <div
            class="h-full rounded-full transition-all duration-500 ease-out"
            :style="{ width: `${row.percent}%`, background: row.color }"
          />
        </div>
        <span class="w-16 shrink-0 text-right text-xs text-mute">
          {{ row.count }} · {{ row.percent }}%
        </span>
      </li>
    </ul>
  </section>
</template>
