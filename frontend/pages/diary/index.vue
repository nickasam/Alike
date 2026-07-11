<script setup lang="ts">
/**
 * 日记广场 — 公开打工日记流 + 写日记。
 *
 * - GET /api/diaries 游标分页（before / limit）拉取公开日记；
 * - 卡片：标题、内容摘要（100字）、作者、心情标签、时间；
 * - 「写日记」按钮需登录，弹窗内 POST /api/diaries 创建后插入流首。
 */
import { useAuthStore } from '~/stores/auth'
import { useEmotions } from '~/composables/useEmotions'

useHead({ title: '日记广场 · Alike' })

const api = useApi()
const authStore = useAuthStore()
const { emotions } = useEmotions()

interface Author {
  id: number
  nickname: string
  avatar_url: string
}
interface Diary {
  id: number
  title?: string
  content: string
  mood?: string
  is_public: boolean
  comment_count: number
  empathy_count: number
  empathized: boolean
  author?: Author
  created_at: string
}
interface ListResp {
  list: Diary[]
  has_more: boolean
  next_cursor: number
}

const diaries = ref<Diary[]>([])
const loading = ref(false)
const hasMore = ref(true)
const cursor = ref(0)
const error = ref('')

/** 心情选项 = 情绪看板同款 8 种情绪（带图标/配色），统一体验。 */

async function loadMore() {
  if (loading.value || !hasMore.value) return
  loading.value = true
  error.value = ''
  try {
    const query: Record<string, number> = { limit: 20 }
    if (cursor.value > 0) query.before = cursor.value
    const res = await api.get<ListResp>('/diaries', query)
    diaries.value.push(...(res.list ?? []))
    hasMore.value = res.has_more
    cursor.value = res.next_cursor
  } catch {
    error.value = '日记加载失败，稍后再试'
  } finally {
    loading.value = false
  }
}

// —— 写日记弹窗 ——
const showEditor = ref(false)
const form = reactive({ title: '', content: '', mood: '', is_public: true })
const submitting = ref(false)
const formError = ref('')

function openEditor() {
  if (!authStore.isAuthenticated) {
    navigateTo({ path: '/login', query: { redirect: '/diary' } })
    return
  }
  form.title = ''
  form.content = ''
  form.mood = ''
  form.is_public = true
  formError.value = ''
  showEditor.value = true
}

async function submit() {
  if (!form.content.trim()) {
    formError.value = '写点什么再发布吧。'
    return
  }
  submitting.value = true
  formError.value = ''
  try {
    const created = await api.post<Diary>('/diaries', {
      title: form.title.trim() || undefined,
      content: form.content.trim(),
      mood: form.mood || undefined,
      is_public: form.is_public,
    })
    // 公开日记插入流首，私密则仅关闭弹窗。
    if (created.is_public) diaries.value.unshift(created)
    showEditor.value = false
  } catch {
    formError.value = '发布失败，请稍后再试。'
  } finally {
    submitting.value = false
  }
}

function excerpt(text: string, n = 100): string {
  return text.length > n ? `${text.slice(0, n)}…` : text
}

function formatDate(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

/** 列表卡「我懂你」：调后端增删共情，用响应更新该卡计数与已共情态。 */
async function onEmpathy(d: Diary, payload: { action: 'add' | 'remove' }) {
  if (!authStore.isAuthenticated) {
    navigateTo({ path: '/login', query: { redirect: '/diary' } })
    return
  }
  try {
    const res = payload.action === 'add'
      ? await api.post<{ empathy_count: number; empathized: boolean }>(`/diaries/${d.id}/empathy`)
      : await api.del<{ empathy_count: number; empathized: boolean }>(`/diaries/${d.id}/empathy`)
    d.empathy_count = res.empathy_count
    d.empathized = res.empathized
  } catch {
    // 失败静默：计数保持不变
  }
}

onMounted(loadMore)
</script>

<template>
  <div class="mx-auto flex max-w-content flex-col gap-5">
    <header class="glass-card animate-rise-in flex items-center justify-between gap-4 p-6">
      <div>
        <h1 class="text-gradient flex items-center gap-2 text-2xl font-extrabold">
          <AppIcon name="book-open" :size="24" />
          打工日记广场
        </h1>
        <p class="mt-1 text-sm text-dim">记录每一天的辛苦，也许有人正读着你的故事。</p>
      </div>
      <button
        type="button"
        class="btn-primary shrink-0 px-4 py-2 text-sm font-semibold"
        @click="openEditor"
      >
        写日记
      </button>
    </header>

    <!-- 日记流 -->
    <div class="flex flex-col gap-3">
      <p v-if="error" class="glass-card p-4 text-center text-sm text-danger">{{ error }}</p>

      <NuxtLink
        v-for="d in diaries"
        :key="d.id"
        :to="`/diary/${d.id}`"
        class="glass-card block p-5"
      >
        <div class="flex items-center justify-between gap-2">
          <h2 class="truncate text-md font-semibold text-text">{{ d.title || '无题日记' }}</h2>
          <span
            v-if="d.mood"
            class="shrink-0 rounded-full bg-warm/15 px-2 py-0.5 text-xs font-medium text-warm"
          >{{ d.mood }}</span>
        </div>
        <p class="mt-1 whitespace-pre-wrap text-sm leading-relaxed text-dim">
          {{ excerpt(d.content) }}
        </p>
        <div class="mt-3 flex items-center gap-2 text-xs text-mute">
          <span>{{ d.author?.nickname ?? '匿名牛马' }}</span>
          <span>·</span>
          <span>{{ formatDate(d.created_at) }}</span>
          <span class="flex items-center gap-1">
            <AppIcon name="hash" :size="12" />{{ d.comment_count }} 评论
          </span>
          <EmpathyButton
            class="ml-auto"
            size="sm"
            :count="d.empathy_count"
            :empathized="d.empathized"
            @empathy="(p) => onEmpathy(d, p)"
            @click.prevent.stop
          />
        </div>
      </NuxtLink>

      <div
        v-if="!loading && !diaries.length && !error"
        class="glass-card p-8 text-center text-sm text-mute"
      >
        还没有人写日记，来发第一篇吧。
      </div>

      <button
        v-if="hasMore && diaries.length"
        type="button"
        class="rounded-md border border-border-strong py-2 text-sm text-dim transition hover:text-text disabled:opacity-50"
        :disabled="loading"
        @click="loadMore"
      >
        {{ loading ? '加载中…' : '加载更多' }}
      </button>
      <p v-else-if="loading" class="py-2 text-center text-sm text-mute">加载中…</p>
    </div>

    <!-- 写日记弹窗 -->
    <div
      v-if="showEditor"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4"
      @click.self="showEditor = false"
    >
      <div class="glass-card w-full max-w-content p-6">
        <h3 class="mb-4 text-lg font-semibold">写一篇打工日记</h3>
        <input
          v-model="form.title"
          type="text"
          maxlength="200"
          placeholder="标题（可选）"
          class="mb-3 w-full rounded-md border border-border bg-surface px-3 py-2 text-sm text-text placeholder:text-mute focus:border-ai-1 focus:outline-none"
        />
        <textarea
          v-model="form.content"
          rows="6"
          maxlength="10000"
          placeholder="今天发生了什么？尽情倾诉吧…"
          class="mb-3 w-full resize-none rounded-md border border-border bg-surface px-3 py-2 text-sm text-text placeholder:text-mute focus:border-ai-1 focus:outline-none"
        />
        <div class="mb-3 flex flex-wrap gap-2">
          <button
            v-for="m in emotions"
            :key="m.key"
            type="button"
            class="inline-flex items-center gap-1.5 rounded-full border px-3 py-1 text-xs font-medium transition"
            :style="form.mood === m.label
              ? { background: m.bg, color: m.color, borderColor: m.color }
              : {}"
            :class="form.mood === m.label
              ? ''
              : 'border-border text-dim hover:text-text'"
            @click="form.mood = form.mood === m.label ? '' : m.label"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.6"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="h-3.5 w-3.5 shrink-0"
              aria-hidden="true"
            >
              <path v-for="(p, i) in m.icon" :key="i" :d="p" />
            </svg>
            {{ m.label }}
          </button>
        </div>
        <label class="mb-4 flex items-center gap-2 text-sm text-dim">
          <input v-model="form.is_public" type="checkbox" class="accent-ai-1" />
          公开到日记广场（取消则仅自己可见）
        </label>
        <p v-if="formError" class="mb-3 text-sm text-danger">{{ formError }}</p>
        <div class="flex justify-end gap-2">
          <button
            type="button"
            class="rounded-md border border-border-strong px-4 py-2 text-sm text-dim transition hover:text-text"
            @click="showEditor = false"
          >取消</button>
          <button
            type="button"
            class="btn-primary px-4 py-2 text-sm font-semibold disabled:opacity-50"
            :disabled="submitting"
            @click="submit"
          >{{ submitting ? '发布中…' : '发布' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>
