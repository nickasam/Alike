<script setup lang="ts">
/**
 * 日记详情 — 正文 + 共鸣 + 评论。
 *
 * - GET /api/diaries/:id 日记正文；
 * - GET /api/diaries/:id/comments 评论列表（分页）；
 * - POST /api/diaries/:id/comments 发表评论（需登录，支持匿名）；
 * - 共情按钮（EmpathyButton）为本地共鸣手势——日记暂无独立共情端点。
 */
import { useAuthStore } from '~/stores/auth'

const route = useRoute()
const api = useApi()
const authStore = useAuthStore()

const diaryId = computed(() => Number(route.params.id))

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
  author?: Author
  created_at: string
}
interface Comment {
  id: number
  diary_id: number
  content: string
  is_anonymous: boolean
  is_deleted: boolean
  author?: Author
  created_at: string
}
interface CommentPage {
  list: Comment[]
  total: number
  page: number
  page_size: number
}

const diary = ref<Diary | null>(null)
const comments = ref<Comment[]>([])
const loading = ref(true)
const error = ref('')

// 评论分页
const page = ref(1)
const total = ref(0)
const commentsLoading = ref(false)

// 本地共鸣手势（日记无独立共情端点）
const empathized = ref(false)
const empathyCount = ref(0)

useHead(() => ({
  title: diary.value ? `${diary.value.title || '无题日记'} · Alike` : '日记详情 · Alike',
}))

async function loadDiary() {
  loading.value = true
  error.value = ''
  try {
    diary.value = await api.get<Diary>(`/diaries/${diaryId.value}`)
  } catch {
    error.value = '日记不存在或加载失败'
  } finally {
    loading.value = false
  }
}

async function loadComments(reset = false) {
  if (commentsLoading.value) return
  if (reset) {
    page.value = 1
    comments.value = []
  }
  commentsLoading.value = true
  try {
    const res = await api.get<CommentPage>(`/diaries/${diaryId.value}/comments`, {
      page: page.value,
      page_size: 20,
    })
    comments.value.push(...(res.list ?? []))
    total.value = res.total
  } catch {
    // 评论加载失败不阻塞正文
  } finally {
    commentsLoading.value = false
  }
}

const hasMoreComments = computed(() => comments.value.length < total.value)

function loadMoreComments() {
  if (!hasMoreComments.value) return
  page.value += 1
  loadComments()
}

// —— 发表评论 ——
const commentText = ref('')
const commentAnon = ref(false)
const posting = ref(false)
const commentError = ref('')

async function postComment() {
  if (!authStore.isAuthenticated) {
    navigateTo({ path: '/login', query: { redirect: `/diary/${diaryId.value}` } })
    return
  }
  if (!commentText.value.trim()) {
    commentError.value = '评论不能为空。'
    return
  }
  posting.value = true
  commentError.value = ''
  try {
    const created = await api.post<Comment>(`/diaries/${diaryId.value}/comments`, {
      content: commentText.value.trim(),
      is_anonymous: commentAnon.value,
    })
    comments.value.unshift(created)
    total.value += 1
    if (diary.value) diary.value.comment_count += 1
    commentText.value = ''
    commentAnon.value = false
  } catch {
    commentError.value = '发表失败，请稍后再试。'
  } finally {
    posting.value = false
  }
}

function onEmpathy(payload: { action: 'add' | 'remove' }) {
  if (payload.action === 'add') {
    empathized.value = true
    empathyCount.value += 1
  } else {
    empathized.value = false
    empathyCount.value = Math.max(0, empathyCount.value - 1)
  }
}

function displayName(c: Comment): string {
  if (c.is_anonymous) return '匿名牛马'
  return c.author?.nickname ?? '牛马'
}

function formatDate(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  const date = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
  const time = `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
  return `${date} ${time}`
}

watch(diaryId, () => {
  loadDiary()
  loadComments(true)
}, { immediate: true })
</script>

<template>
  <article class="mx-auto flex max-w-content flex-col gap-5">
    <p v-if="loading" class="glass-card p-8 text-center text-sm text-mute">加载中…</p>
    <p v-else-if="error" class="glass-card p-8 text-center text-sm text-danger">{{ error }}</p>

    <template v-else-if="diary">
      <!-- 正文 -->
      <section class="glass-card animate-rise-in p-6">
        <div class="flex items-start justify-between gap-2">
          <h1 class="text-xl font-bold text-text">{{ diary.title || '无题日记' }}</h1>
          <span
            v-if="diary.mood"
            class="shrink-0 rounded-full bg-warm/15 px-2 py-0.5 text-xs font-medium text-warm"
          >{{ diary.mood }}</span>
        </div>
        <div class="mt-2 flex items-center gap-2 text-xs text-mute">
          <NuxtLink
            v-if="diary.author"
            :to="`/profile/${diary.author.id}`"
            class="transition hover:text-text"
          >{{ diary.author.nickname }}</NuxtLink>
          <span v-else>匿名牛马</span>
          <span>·</span>
          <span>{{ formatDate(diary.created_at) }}</span>
        </div>
        <p class="mt-4 whitespace-pre-wrap break-words text-base leading-relaxed text-dim">
          {{ diary.content }}
        </p>
        <div class="mt-5">
          <EmpathyButton
            :count="empathyCount"
            :empathized="empathized"
            @empathy="onEmpathy"
          />
        </div>
      </section>

      <!-- 评论区 -->
      <section>
        <h2 class="mb-3 flex items-center gap-2 text-lg font-semibold">
          <AppIcon name="hash" :size="20" />
          评论 {{ total }}
        </h2>

        <!-- 发表评论 -->
        <div class="glass-card mb-4 p-4">
          <textarea
            v-model="commentText"
            rows="3"
            maxlength="2000"
            placeholder="说点什么安慰一下 ta…"
            class="w-full resize-none rounded-md border border-border bg-surface px-3 py-2 text-sm text-text placeholder:text-mute focus:border-ai-1 focus:outline-none"
          />
          <div class="mt-2 flex items-center justify-between gap-2">
            <label class="flex items-center gap-2 text-sm text-dim">
              <input v-model="commentAnon" type="checkbox" class="accent-ai-1" />
              匿名评论
            </label>
            <button
              type="button"
              class="btn-primary px-4 py-2 text-sm font-semibold disabled:opacity-50"
              :disabled="posting"
              @click="postComment"
            >{{ posting ? '发表中…' : '发表评论' }}</button>
          </div>
          <p v-if="commentError" class="mt-2 text-sm text-danger">{{ commentError }}</p>
        </div>

        <!-- 评论列表 -->
        <div class="flex flex-col gap-3">
          <div
            v-if="!comments.length && !commentsLoading"
            class="glass-card p-6 text-center text-sm text-mute"
          >
            还没有评论，来第一个安慰 ta 吧。
          </div>
          <div v-for="c in comments" :key="c.id" class="glass-card flex gap-3 p-4">
            <div
              class="grid h-9 w-9 shrink-0 place-items-center rounded-md text-sm font-semibold text-white"
              :class="c.is_anonymous ? 'bg-surface-hover text-dim' : 'bg-grad-ai'"
            >
              {{ c.is_anonymous ? '匿' : (c.author?.nickname?.charAt(0) ?? '牛') }}
            </div>
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <span class="text-sm font-semibold text-text">{{ displayName(c) }}</span>
                <span class="text-xs text-mute">{{ formatDate(c.created_at) }}</span>
              </div>
              <p
                class="mt-1 whitespace-pre-wrap break-words text-sm text-dim"
                :class="{ 'italic text-mute': c.is_deleted }"
              >{{ c.content }}</p>
            </div>
          </div>

          <button
            v-if="hasMoreComments"
            type="button"
            class="rounded-md border border-border-strong py-2 text-sm text-dim transition hover:text-text disabled:opacity-50"
            :disabled="commentsLoading"
            @click="loadMoreComments"
          >{{ commentsLoading ? '加载中…' : '加载更多评论' }}</button>
        </div>
      </section>
    </template>
  </article>
</template>
