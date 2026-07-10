<script setup lang="ts">
/**
 * WarmestBoard — 首页右侧栏「今日牛马榜」（最暖牛马）。
 * 拉 GET /api/ranking/warmest（按给出共情数降序），每条跳 /profile/:user_id。
 */
const api = useApi()

interface RankUser {
  user_id: number
  nickname: string
  avatar_url: string
  level: number
  metric: number
}

const list = ref<RankUser[]>([])
const loading = ref(true)

async function load() {
  try {
    const res = await api.get<{ list: RankUser[] }>('/ranking/warmest?limit=5')
    list.value = res?.list ?? []
  } catch {
    list.value = []
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <section class="glass-card p-4" aria-label="今日牛马榜">
    <h3 class="mb-2 flex items-center gap-2 text-md font-semibold">
      <AppIcon name="trophy" :size="18" />
      今日牛马榜
      <NuxtLink to="/ranking" class="ml-auto text-xs font-normal text-mute transition hover:text-text">
        查看全部
      </NuxtLink>
    </h3>

    <p v-if="loading" class="text-sm text-mute">加载中…</p>
    <p v-else-if="list.length === 0" class="text-sm text-mute">
      还没有人给出共情，第一个抱抱别人的就是你。
    </p>
    <ol v-else class="flex flex-col gap-1.5">
      <li v-for="(u, i) in list" :key="u.user_id">
        <NuxtLink
          :to="`/profile/${u.user_id}`"
          class="flex items-center gap-2.5 rounded-md px-1.5 py-1 transition hover:bg-surface-hover"
        >
          <span
            class="w-4 shrink-0 text-center text-xs font-bold"
            :class="i < 3 ? 'text-warm' : 'text-mute'"
          >
            {{ i + 1 }}
          </span>
          <img
            v-if="u.avatar_url"
            :src="u.avatar_url"
            :alt="u.nickname"
            class="h-7 w-7 rounded-md object-cover"
          />
          <span
            v-else
            class="grid h-7 w-7 shrink-0 place-items-center rounded-md bg-grad-ai text-xs font-semibold text-white"
          >
            {{ u.nickname?.[0] ?? '牛' }}
          </span>
          <span class="min-w-0 flex-1 truncate text-sm font-medium text-text">
            {{ u.nickname }}
          </span>
          <span class="shrink-0 text-xs text-empathy">{{ u.metric }}</span>
        </NuxtLink>
      </li>
    </ol>
  </section>
</template>
