<script setup lang="ts">
/**
 * 个人主页 — 用户信息 + 统计 + 最近日记。
 *
 * - GET /api/users/:id 拉取公开资料（头像/昵称/行业/岗位/工龄/等级）；
 * - 统计卡：被共情数 / 给出共情数 / 累计打卡天数；
 * - GET /api/users/:id/diaries 拉取最近日记（后端返回列表时渲染）；
 * - 编辑资料按钮仅本人（登录 uid === 路由 id）可见。
 */
import { useAuthStore } from '~/stores/auth'

const route = useRoute()
const api = useApi()
const authStore = useAuthStore()

const userId = computed(() => Number(route.params.id))

interface PublicUser {
  id: number
  nickname: string
  avatar_url: string
  bio: string
  industry: string
  job_title: string
  work_years: number
  level: number
  empathy_received: number
  empathy_given: number
  total_check_in_days: number
  created_at: string
}

interface Diary {
  id: number
  title?: string
  content: string
  mood?: string
  created_at: string
}

const user = ref<PublicUser | null>(null)
const diaries = ref<Diary[]>([])
const loading = ref(true)
const error = ref('')

const isSelf = computed(
  () => authStore.isAuthenticated && authStore.user?.id === userId.value,
)

useHead(() => ({
  title: user.value ? `${user.value.nickname} · Alike` : '个人主页 · Alike',
}))

async function loadUser() {
  loading.value = true
  error.value = ''
  try {
    user.value = await api.get<PublicUser>(`/users/${userId.value}`)
  } catch {
    error.value = '用户不存在或加载失败'
  } finally {
    loading.value = false
  }
}

async function loadDiaries() {
  try {
    const res = await api.get<unknown>(`/users/${userId.value}/diaries`)
    // 后端返回列表则渲染；stub 阶段返回对象时视为空。
    const list = Array.isArray(res)
      ? res
      : (res as { list?: Diary[] })?.list
    diaries.value = Array.isArray(list) ? (list as Diary[]) : []
  } catch {
    diaries.value = []
  }
}

const stats = computed(() => [
  { label: '被共情', value: user.value?.empathy_received ?? 0, color: 'text-empathy' },
  { label: '给出共情', value: user.value?.empathy_given ?? 0, color: 'text-ai-1' },
  { label: '累计打卡', value: user.value?.total_check_in_days ?? 0, color: 'text-warm' },
])

function excerpt(text: string, n = 80): string {
  return text.length > n ? `${text.slice(0, n)}…` : text
}

function formatDate(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

watch(userId, () => {
  loadUser()
  loadDiaries()
}, { immediate: true })
</script>

<template>
  <div class="mx-auto flex max-w-content flex-col gap-5">
    <p v-if="loading" class="glass-card p-8 text-center text-sm text-mute">加载中…</p>
    <p v-else-if="error" class="glass-card p-8 text-center text-sm text-danger">{{ error }}</p>

    <template v-else-if="user">
      <!-- 资料卡 -->
      <section class="glass-card animate-rise-in p-6">
        <div class="flex items-start gap-4">
          <div class="grid h-16 w-16 shrink-0 place-items-center rounded-lg bg-grad-ai text-2xl font-bold text-white">
            {{ user.nickname.charAt(0) }}
          </div>
          <div class="min-w-0 flex-1">
            <div class="flex flex-wrap items-center gap-2">
              <h1 class="text-xl font-bold text-text">{{ user.nickname }}</h1>
              <span class="rounded-full bg-warm/15 px-2 py-0.5 text-xs font-semibold text-warm">
                Lv.{{ user.level }} 牛马
              </span>
            </div>
            <p class="mt-1 text-sm text-dim">{{ user.bio || '这只牛马还没有留下签名。' }}</p>
            <div class="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-xs text-mute">
              <span v-if="user.industry">行业：{{ user.industry }}</span>
              <span v-if="user.job_title">岗位：{{ user.job_title }}</span>
              <span v-if="user.work_years">工龄：{{ user.work_years }} 年</span>
            </div>
          </div>
          <NuxtLink
            v-if="isSelf"
            to="/settings"
            class="shrink-0 rounded-md border border-border-strong px-3 py-1.5 text-sm text-dim transition hover:text-text"
          >
            编辑资料
          </NuxtLink>
        </div>
      </section>

      <!-- 统计 -->
      <section class="grid grid-cols-3 gap-4">
        <div v-for="s in stats" :key="s.label" class="glass-card p-4 text-center">
          <p class="text-2xl font-extrabold" :class="s.color">{{ s.value }}</p>
          <p class="mt-1 text-xs text-mute">{{ s.label }}</p>
        </div>
      </section>

      <!-- 最近日记 -->
      <section>
        <h2 class="mb-3 flex items-center gap-2 text-lg font-semibold">
          <AppIcon name="book-open" :size="20" />
          最近日记
        </h2>
        <div v-if="!diaries.length" class="glass-card p-6 text-center text-sm text-mute">
          还没有公开的打工日记。
        </div>
        <div v-else class="flex flex-col gap-3">
          <NuxtLink
            v-for="d in diaries"
            :key="d.id"
            :to="`/diary/${d.id}`"
            class="glass-card block p-4"
          >
            <div class="flex items-center justify-between gap-2">
              <h3 class="truncate text-md font-semibold text-text">
                {{ d.title || '无题日记' }}
              </h3>
              <span class="shrink-0 text-xs text-mute">{{ formatDate(d.created_at) }}</span>
            </div>
            <p class="mt-1 text-sm text-dim">{{ excerpt(d.content) }}</p>
          </NuxtLink>
        </div>
      </section>
    </template>
  </div>
</template>
