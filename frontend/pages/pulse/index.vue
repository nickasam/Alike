<script setup lang="ts">
/**
 * 最近发生（Pulse / 世界脉搏）
 * M1：GitHub Trending Tab 上线；AI 圈大事仍占位（M2 接入）。
 * 详见 docs/plans/pulse-module-design.md 与 docs/design/07-pulse.html。
 */

useHead({ title: '最近发生 · Alike' })

interface Topic {
  id: number
  slug: string
  name: string
  emoji?: string
  description?: string
  sort_order: number
  last_fetched_at?: string
  refresh_interval_min?: number
}

interface RepoExtra {
  lang?: string
  lang_color?: string
  total_stars?: number
  forks?: number
}

interface Item {
  id: number
  topic_id: number
  source: string
  source_id: string
  title: string
  summary?: string
  url: string
  author?: string
  score: number
  extra?: RepoExtra
  captured_at: string
  published_at?: string
}

interface ItemsData {
  topic: Topic
  list: Item[]
  last_fetched_at?: string
  stale: boolean
}

const api = useApi()
const topics = ref<Topic[]>([])
const topicsLoading = ref(true)
const activeSlug = ref<string>('')

const itemsCache = ref<Record<string, ItemsData>>({})
const itemsLoading = ref(false)
const itemsError = ref('')

const currentData = computed<ItemsData | null>(() =>
  activeSlug.value ? itemsCache.value[activeSlug.value] ?? null : null,
)

/** 加载专题列表 */
async function loadTopics() {
  try {
    const res = await api.get<{ list: Topic[] }>('/pulse/topics')
    topics.value = res?.list ?? []
    if (topics.value.length > 0) {
      // 默认选中排序最靠前的（即 sort_order 最小的，接口已按此排序）
      activeSlug.value = topics.value[0].slug
      await loadItems(activeSlug.value)
    }
  } catch (e) {
    console.warn('[pulse] load topics failed:', e)
  } finally {
    topicsLoading.value = false
  }
}

/** 加载某专题的条目（切 Tab 时用；带缓存） */
async function loadItems(slug: string) {
  if (!slug) return
  if (itemsCache.value[slug]) return // 缓存命中
  itemsLoading.value = true
  itemsError.value = ''
  try {
    const data = await api.get<ItemsData>(`/pulse/topics/${slug}/items`)
    itemsCache.value[slug] = data
  } catch (e: any) {
    itemsError.value = e?.message ?? '拉取失败'
  } finally {
    itemsLoading.value = false
  }
}

/** 切 Tab */
async function switchTab(slug: string) {
  activeSlug.value = slug
  await loadItems(slug)
}

/** 计算相对时间（"3 分钟前"），用于卡片和更新时间条 */
function relativeTime(iso?: string): string {
  if (!iso) return ''
  const then = new Date(iso).getTime()
  if (isNaN(then)) return ''
  const diff = Date.now() - then
  const sec = Math.floor(diff / 1000)
  if (sec < 60) return '刚刚'
  const min = Math.floor(sec / 60)
  if (min < 60) return `${min} 分钟前`
  const hr = Math.floor(min / 60)
  if (hr < 24) return `${hr} 小时前`
  const day = Math.floor(hr / 24)
  if (day < 30) return `${day} 天前`
  return new Date(iso).toLocaleDateString('zh-CN')
}

/** 格式化 HH:MM */
function formatHHMM(iso?: string): string {
  if (!iso) return '--:--'
  const d = new Date(iso)
  if (isNaN(d.getTime())) return '--:--'
  return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

/** 千分位友好显示 12.7k / 158.2k */
function humanNum(n?: number): string {
  if (!n || n <= 0) return '0'
  if (n < 1000) return String(n)
  if (n < 1_000_000) {
    const k = n / 1000
    return k >= 100 ? `${Math.round(k)}k` : `${k.toFixed(1)}k`
  }
  return `${(n / 1_000_000).toFixed(1)}M`
}

onMounted(() => {
  loadTopics()
})

/** 排名色号复用主体前三名的金/银/铜（顺序同 04-ranking 领奖台） */
function rankClass(i: number): string {
  if (i === 0) return 'text-gold'
  if (i === 1) return 'text-[#cbd5e1]'
  if (i === 2) return 'text-warm'
  return 'text-mute'
}
</script>

<template>
  <div class="mx-auto max-w-content-wide">
    <!-- Hero -->
    <section class="glass-card animate-rise-in relative overflow-hidden p-8">
      <div
        class="pointer-events-none absolute -right-16 -top-24 h-80 w-80 rounded-full"
        style="background: radial-gradient(circle, rgba(99,102,241,.22), transparent 65%);"
        aria-hidden="true"
      />
      <h1 class="text-gradient relative text-2xl font-extrabold leading-tight md:text-3xl">
        最近发生
      </h1>
      <p class="relative mt-2 text-md text-dim">
        打工人不切平台就能同步的世界脉搏。只搬运、不评论、看一眼就走。
      </p>
      <div class="relative mt-4 flex flex-wrap items-center gap-2">
        <span
          class="inline-flex items-center gap-1.5 rounded-full border border-empathy/30 bg-empathy-soft px-3 py-1 text-xs font-semibold text-empathy"
        >
          <span class="pulse-dot inline-block h-1.5 w-1.5 rounded-full bg-empathy" aria-hidden="true" />
          正在跳动
        </span>
        <span
          v-if="currentData?.last_fetched_at"
          class="inline-flex items-center gap-1.5 rounded-full border border-border bg-surface px-3 py-1 text-xs text-dim"
        >
          截至 <b class="text-text">{{ formatHHMM(currentData.last_fetched_at) }}</b> 更新 ·
          {{ relativeTime(currentData.last_fetched_at) }}
        </span>
        <span
          v-else-if="!topicsLoading && !itemsLoading"
          class="inline-flex items-center gap-1.5 rounded-full border border-border bg-surface px-3 py-1 text-xs text-mute"
        >
          尚未成功抓取
        </span>
      </div>
    </section>

    <!-- stale 警告 -->
    <div
      v-if="currentData?.stale"
      class="mt-4 flex items-start gap-2 rounded-lg border p-3 text-sm"
      style="background:rgba(234,179,8,.08);border-color:rgba(234,179,8,.35);color:#eab308;"
      role="status"
    >
      <span aria-hidden="true">⚠️</span>
      <span>数据可能不新鲜——最近 6 小时没成功抓到更新，你看到的是上一次成功的快照。</span>
    </div>

    <!-- Tabs -->
    <section class="mt-6">
      <h2 class="mb-3 flex items-center gap-2 text-lg font-semibold">
        <AppIcon name="activity" :size="20" />
        专题
        <span class="ml-auto text-xs font-normal text-mute">
          每 <b class="text-text">{{ topics.find(t => t.slug === activeSlug)?.refresh_interval_min ?? 60 }}</b> 分钟更新一次
        </span>
      </h2>
      <div class="mb-4 flex flex-wrap gap-2" role="tablist" aria-label="专题切换">
        <button
          v-for="t in topics"
          :key="t.slug"
          type="button"
          role="tab"
          :aria-selected="activeSlug === t.slug"
          :class="[
            'tab-pill inline-flex items-center gap-2 rounded-full border px-4 py-2 text-sm font-semibold transition',
            activeSlug === t.slug
              ? 'bg-grad-ai text-white shadow-glow-ai border-transparent'
              : 'bg-surface text-dim border-border hover:text-text hover:border-border-strong',
          ]"
          @click="switchTab(t.slug)"
        >
          <span aria-hidden="true">{{ t.emoji }}</span>
          <span>{{ t.name }}</span>
        </button>
      </div>
    </section>

    <!-- Feed -->
    <section class="mt-1">
      <p v-if="topicsLoading" class="glass-card p-5 text-sm text-mute">
        加载专题中…
      </p>
      <p v-else-if="itemsLoading" class="glass-card p-5 text-sm text-mute">
        加载卡片中…
      </p>
      <p v-else-if="itemsError" class="glass-card p-5 text-sm text-mute">
        暂时拿不到数据（{{ itemsError }}）
      </p>
      <p v-else-if="!currentData?.list || currentData.list.length === 0" class="glass-card p-5 text-sm text-mute">
        <template v-if="activeSlug === 'ai-news'">
          🤖 AI 圈大事将在 M2 上线。GitHub 每日 Star 榜先试试？
        </template>
        <template v-else>
          世界今天很安静，或者服务器在偷懒 🐢
        </template>
      </p>
      <ul v-else class="flex flex-col gap-3">
        <li v-for="(it, i) in currentData.list" :key="it.id">
          <!-- GitHub 卡片 -->
          <a
            v-if="it.source === 'github'"
            :href="it.url"
            target="_blank"
            rel="noopener"
            class="repo-row glass-card relative flex items-stretch gap-3.5 p-4 pl-5"
          >
            <span class="rank shrink-0 pt-0.5 text-sm font-extrabold text-mute" :class="rankClass(i)">
              {{ i + 1 }}
            </span>
            <div class="min-w-0 flex-1 flex flex-col gap-1.5">
              <h3 class="truncate text-md font-semibold leading-tight">
                <span class="text-dim font-medium">{{ it.author }}</span>
                <span class="text-mute font-normal">/</span>{{ it.source_id.split('/')[1] }}
              </h3>
              <p v-if="it.summary" class="line-clamp-2 text-sm text-dim">
                {{ it.summary }}
              </p>
              <div class="mt-0.5 flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-mute">
                <span v-if="it.extra?.lang" class="inline-flex items-center gap-1.5">
                  <span
                    class="inline-block h-2.5 w-2.5 rounded-full"
                    :style="{ background: it.extra?.lang_color || '#94a3b8' }"
                    aria-hidden="true"
                  />
                  {{ it.extra.lang }}
                </span>
                <span v-if="it.extra?.total_stars" class="inline-flex items-center gap-1">
                  ⭐ {{ humanNum(it.extra.total_stars) }}
                </span>
                <span v-if="it.extra?.forks" class="inline-flex items-center gap-1">
                  ⑂ {{ humanNum(it.extra.forks) }}
                </span>
              </div>
            </div>
            <div class="shrink-0 flex flex-col items-end justify-between gap-1.5">
              <span class="text-warm inline-flex items-center gap-1 text-xl font-extrabold leading-none tabular-nums">
                +{{ it.score }}
                <span class="text-warm text-lg" aria-hidden="true">⭐</span>
              </span>
              <span class="text-mute inline-flex items-center gap-1 text-xs font-semibold">
                今日
              </span>
            </div>
            <!-- 左侧橙彩条：视觉标记这是 GitHub 卡片 -->
            <span class="row-accent" aria-hidden="true" />
          </a>

          <!-- HN / AI 卡片：M2 上线，M1 阶段用简单展示 -->
          <a
            v-else
            :href="it.url"
            target="_blank"
            rel="noopener"
            class="glass-card block p-4"
          >
            <h3 class="text-md font-semibold">{{ it.title }}</h3>
            <p v-if="it.summary" class="mt-1 text-sm text-dim">{{ it.summary }}</p>
            <div class="mt-2 text-xs text-mute">
              {{ it.score }} points · {{ relativeTime(it.published_at) }}
            </div>
          </a>
        </li>
      </ul>
    </section>

    <!-- 克制感说明 -->
    <div
      class="mt-6 rounded-lg border p-4 text-sm text-dim"
      style="background:linear-gradient(135deg,rgba(99,102,241,.08),rgba(34,211,238,.06));border-color:rgba(99,102,241,.2);"
    >
      这里是<b class="text-ai-2">只读的世界之窗</b>——不做评论、不推送、不做红点。想聊感受？去
      <NuxtLink to="/diary" class="text-ai-2 font-semibold hover:underline">日记广场</NuxtLink> 吐槽。
    </div>
  </div>
</template>

<style scoped>
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

/* GitHub 仓库卡片左侧暖色彩条：视觉一眼区分类型 */
.repo-row {
  overflow: hidden;
  transition: all .25s var(--ease-out);
}
.repo-row:hover {
  transform: translateY(-2px);
  border-color: var(--border-strong);
}
.row-accent {
  position: absolute;
  left: 0; top: 0; bottom: 0;
  width: 3px;
  background: var(--grad-warm);
  opacity: .55;
}
</style>
