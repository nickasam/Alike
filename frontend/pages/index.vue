<script setup lang="ts">
/**
 * 首页 — 热门频道 + 今日牛马榜。
 * 热门频道从后端 /api/channels 拉取真实数据（按成员数排序），
 * 保证卡片文字与点击后跳转的频道一致。
 */
import { useChannelStore, type Channel } from '~/stores/channel'

useHead({ title: 'Alike · 汇聚天下牛马' })

const api = useApi()
const channelStore = useChannelStore()

const hotChannels = ref<Channel[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await api.get<{ list: Channel[] }>('/channels?page=1&page_size=6')
    const list = res?.list ?? []
    hotChannels.value = list
    channelStore.setChannels(list)
  } catch {
    hotChannels.value = []
  } finally {
    loading.value = false
  }
})

/** 频道名可能已含 # 前缀，去掉后由模板统一加 #，避免 ## 重复。 */
function displayName(name: string): string {
  return name.replace(/^#/, '')
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Hero -->
    <section class="glass-card animate-rise-in p-8">
      <h1 class="text-gradient text-2xl font-extrabold leading-tight md:text-3xl">
        汇聚天下牛马
      </h1>
      <p class="mt-2 text-md text-dim">
        总有人懂你的辛苦，你不是一个人在扛。
      </p>
      <div class="mt-5 flex flex-wrap gap-3">
        <button class="btn-primary px-5 py-2 text-base font-semibold">
          进入频道吐槽
        </button>
        <button
          class="rounded-md border border-border-strong px-5 py-2 text-base text-dim transition hover:text-text"
        >
          今日牛马榜
        </button>
      </div>
    </section>

    <!-- 热门频道 -->
    <section>
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
          <div class="flex items-center justify-between">
            <h3 class="text-md font-semibold">#{{ displayName(ch.name) }}</h3>
            <span class="text-xs text-mute">{{ ch.member_count }} 位牛马</span>
          </div>
          <p v-if="ch.description" class="mt-1 text-sm text-dim">{{ ch.description }}</p>
        </NuxtLink>
      </div>
    </section>
  </div>
</template>
