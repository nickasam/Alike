<script setup lang="ts">
/**
 * 通知中心 — 展示当前用户的通知列表，支持标记单条/全部已读。
 * 共情/回复通知携带 ref_id（消息 ID），点击暂只标已读（跳转留待线程深链）。
 */
import { useNotifications, type Notification } from '~/composables/useNotifications'

definePageMeta({ middleware: 'auth' })
useHead({ title: '通知 · Alike' })

const { fetchList, markRead, markAllRead } = useNotifications()

const list = ref<Notification[]>([])
const loading = ref(true)
const total = ref(0)

const typeLabel: Record<string, string> = {
  empathy: '抱抱',
  reply: '回复',
  mention: '提及',
  system: '系统',
}

const typeIcon: Record<string, string> = {
  empathy: 'sparkles',
  reply: 'hash',
  mention: 'user',
  system: 'bell',
}

onMounted(load)

async function load() {
  loading.value = true
  try {
    const res = await fetchList(1, 50)
    list.value = res?.list ?? []
    total.value = res?.total ?? 0
  } catch {
    list.value = []
  } finally {
    loading.value = false
  }
}

async function onRead(n: Notification) {
  if (n.is_read) return
  try {
    await markRead(n.id)
    n.is_read = true
  } catch {
    // 静默，用户可重试
  }
}

async function onReadAll() {
  try {
    await markAllRead()
    list.value.forEach((n) => (n.is_read = true))
  } catch {
    // 静默
  }
}

/** 简易相对时间：今天显示时:分，否则显示月-日。 */
function fmtTime(iso: string): string {
  const d = new Date(iso)
  const now = new Date()
  const sameDay = d.toDateString() === now.toDateString()
  return sameDay
    ? `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
    : `${d.getMonth() + 1}-${d.getDate()}`
}

const hasUnread = computed(() => list.value.some((n) => !n.is_read))
</script>

<template>
  <div class="mx-auto w-full max-w-2xl">
    <section class="glass-card animate-rise-in p-6">
      <div class="mb-4 flex items-center gap-2">
        <AppIcon name="bell" :size="22" />
        <h1 class="text-xl font-semibold">通知</h1>
        <button
          v-if="hasUnread"
          class="ml-auto rounded-md border border-border-strong px-3 py-1 text-sm text-dim transition hover:text-text"
          @click="onReadAll"
        >
          全部已读
        </button>
      </div>

      <p v-if="loading" class="text-sm text-mute">加载中…</p>
      <p v-else-if="list.length === 0" class="py-8 text-center text-sm text-mute">
        还没有通知，去发条心声、抱抱别人吧。
      </p>
      <ul v-else class="flex flex-col gap-1">
        <li
          v-for="n in list"
          :key="n.id"
          class="flex cursor-pointer items-start gap-3 rounded-md px-3 py-3 transition hover:bg-surface-hover"
          :class="{ 'opacity-60': n.is_read }"
          @click="onRead(n)"
        >
          <span
            class="mt-0.5 grid h-8 w-8 shrink-0 place-items-center rounded-md"
            :class="n.is_read ? 'bg-surface-hover text-mute' : 'bg-grad-ai text-white'"
          >
            <AppIcon :name="typeIcon[n.type] ?? 'bell'" :size="16" />
          </span>
          <div class="min-w-0 flex-1">
            <p class="text-sm text-text">
              <span class="mr-1 text-xs text-mute">[{{ typeLabel[n.type] ?? n.type }}]</span>
              {{ n.content }}
            </p>
            <p class="mt-0.5 text-xs text-mute">{{ fmtTime(n.created_at) }}</p>
          </div>
          <span
            v-if="!n.is_read"
            class="mt-1 h-2 w-2 shrink-0 rounded-full bg-danger"
            aria-label="未读"
          />
        </li>
      </ul>
    </section>
  </div>
</template>
