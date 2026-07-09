<script setup lang="ts">
/**
 * EmotionBoard — 频道实时情绪看板。
 *
 * - 拉取 GET /api/channels/:id/emotion-board，渲染 8 种情绪的分布进度条；
 * - channelId 未提供时从路由 /channel/:id 推断；非频道页则展示空态；
 * - 每种情绪显示数量与百分比，颜色/标签取自 useEmotions。
 */
import { useEmotions, type EmotionMeta } from '~/composables/useEmotions'

const props = defineProps<{ channelId?: number; compact?: boolean }>()

const api = useApi()
const route = useRoute()
const { emotions, find } = useEmotions()

interface Count {
  emotion: string
  count: number
}
interface Board {
  emotions: Count[]
  total: number
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

async function load() {
  const id = channelId.value
  if (!id) {
    board.value = null
    return
  }
  loading.value = true
  error.value = ''
  try {
    board.value = await api.get<Board>(`/channels/${id}/emotion-board`)
  } catch {
    error.value = '情绪看板加载失败'
    board.value = null
  } finally {
    loading.value = false
  }
}

/** 各情绪计数（含 0），保证按 useEmotions 顺序稳定渲染。 */
interface Row extends EmotionMeta {
  count: number
  percent: number
}
const rows = computed<Row[]>(() => {
  const total = board.value?.total ?? 0
  const counts = new Map(board.value?.emotions.map((e) => [e.emotion, e.count]))
  return emotions.map((meta) => {
    const count = counts.get(meta.key) ?? 0
    return {
      ...meta,
      count,
      percent: total > 0 ? Math.round((count / total) * 100) : 0,
    }
  })
})

const total = computed(() => board.value?.total ?? 0)

watch(channelId, load, { immediate: true })
</script>

<template>
  <section class="glass-card p-4" aria-label="情绪看板">
    <h3 class="mb-3 flex items-center gap-2 text-md font-semibold">
      <AppIcon name="sparkles" :size="18" />
      今日情绪看板
      <span v-if="total" class="ml-auto text-xs font-normal text-mute">{{ total }} 条心声</span>
    </h3>

    <p v-if="!channelId" class="text-sm text-mute">进入频道查看实时情绪分布。</p>
    <p v-else-if="loading && !board" class="text-sm text-mute">加载中…</p>
    <p v-else-if="error" class="text-sm text-danger">{{ error }}</p>
    <p v-else-if="total === 0" class="text-sm text-mute">还没有情绪记录，来发第一条吧。</p>

    <ul v-else class="flex flex-col gap-2">
      <li v-for="row in rows" :key="row.key" class="flex items-center gap-3">
        <span
          class="w-12 shrink-0 rounded-full px-2 py-0.5 text-center text-xs font-medium"
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
        <span class="w-16 shrink-0 text-right text-xs text-mute">
          {{ row.count }} · {{ row.percent }}%
        </span>
      </li>
    </ul>
  </section>
</template>
