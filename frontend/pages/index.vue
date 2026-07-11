<script setup lang="ts">
/**
 * 首页 — Hero + 热门频道。
 * 热门频道从 /api/channels 拉真实数据；今日情绪看板 / 今日牛马榜见右侧栏（全站聚合）。
 */
import { useChannelStore, type Channel } from '~/stores/channel'

useHead({ title: 'Alike · 汇聚天下牛马' })

const api = useApi()
const channelStore = useChannelStore()

const hotChannels = ref<Channel[]>([])
const loading = ref(true)

onMounted(loadHotChannels)

async function loadHotChannels() {
  try {
    const res = await api.get<{ list: Channel[] }>('/channels?page=1&page_size=6')
    const list = res?.list ?? []
    hotChannels.value = list
    channelStore.setHotChannels(list)
  } catch {
    hotChannels.value = []
  } finally {
    loading.value = false
  }
}

/** 频道名可能已含 # 前缀，去掉后由模板统一加 #，避免 ## 重复。 */
function displayName(name: string): string {
  return name.replace(/^#/, '')
}

/** 「进入频道吐槽」目标：优先热门首个频道，其次侧边栏全量首个频道。 */
const firstChannelId = computed<number | null>(
  () => hotChannels.value[0]?.id ?? channelStore.channels[0]?.id ?? null,
)

function enterChannel() {
  if (firstChannelId.value != null) {
    navigateTo(`/channel/${firstChannelId.value}`)
  } else if (import.meta.client) {
    document.getElementById('hot-channels')?.scrollIntoView({ behavior: 'smooth' })
  }
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Hero -->
    <section class="glass-card animate-rise-in relative overflow-hidden p-8">
      <div
        class="pointer-events-none absolute -right-16 -top-24 h-80 w-80 rounded-full"
        style="background: radial-gradient(circle, rgba(99,102,241,.22), transparent 65%)"
        aria-hidden="true"
      />
      <h1 class="text-gradient relative text-2xl font-extrabold leading-tight md:text-3xl">
        汇聚天下牛马
      </h1>
      <p class="relative mt-2 text-md text-dim">
        总有人懂你的辛苦，你不是一个人在扛。
      </p>
      <div class="relative mt-5 flex flex-wrap gap-3">
        <button
          class="btn-primary inline-flex items-center gap-2 px-5 py-2 text-base font-semibold"
          @click="enterChannel"
        >
          <AppIcon name="hash" :size="18" />
          进入频道吐槽
        </button>
        <NuxtLink
          to="/ranking"
          class="inline-flex items-center gap-2 rounded-md border border-border-strong bg-surface-solid px-5 py-2 text-base text-text shadow-md transition hover:border-warm hover:text-warm hover:shadow-glow-warm"
        >
          <AppIcon name="trophy" :size="18" />
          今日牛马榜
        </NuxtLink>
      </div>
    </section>

    <!-- 热门频道 -->
    <section id="hot-channels">
      <h2 class="mb-3 flex items-center gap-2 text-lg font-semibold">
        <AppIcon name="hash" :size="20" />
        热门频道
      </h2>
      <p v-if="loading" class="glass-card p-5 text-sm text-mute">加载中…</p>
      <p v-else-if="hotChannels.length === 0" class="glass-card p-5 text-sm text-mute">
        暂无频道
      </p>
      <div v-else class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <NuxtLink
          v-for="ch in hotChannels"
          :key="ch.id"
          :to="`/channel/${ch.id}`"
          class="glass-card block p-5"
        >
          <div class="flex items-center gap-3">
            <span
              class="grid h-11 w-11 shrink-0 place-items-center rounded-md bg-grad-ai text-white shadow-glow-ai"
              aria-hidden="true"
            >
              <AppIcon name="hash" :size="22" />
            </span>
            <h3 class="min-w-0 flex-1 truncate text-md font-semibold">#{{ displayName(ch.name) }}</h3>
            <span class="shrink-0 text-xs text-mute">{{ ch.member_count }} 位牛马</span>
          </div>
          <p v-if="ch.description" class="mt-3 text-sm text-dim">{{ ch.description }}</p>
        </NuxtLink>
      </div>
    </section>
  </div>
</template>
