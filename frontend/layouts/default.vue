<script setup lang="ts">
/**
 * default 布局 — 三列响应式骨架。
 * 断点（交互规范 §1.1）：
 *   桌面 ≥1280(xl)  三列：侧边栏 260 + 主内容 1fr + 右侧栏 320
 *   平板 768-1279(md) 两列：侧边栏 + 主内容，右侧栏隐藏
 *   手机 <768        单列：侧边栏转抽屉，主内容全宽
 */
const drawerOpen = ref(false)

// 频道页用频道级实时看板；其它页面（首页等）用全站今日聚合。
const route = useRoute()
const globalBoard = computed(() => !route.path.startsWith('/channel/'))

// 全局 WebSocket：登录后建立连接并消费 notification 事件，实时更新铃铛未读数。
// 频道页也复用这条单例连接（useWebSocket 为模块级单例）。
const auth = useAuth()
const ws = useWebSocket()
const { unread, refreshUnread } = useNotifications()
let offNotif: (() => void) | null = null

onMounted(() => {
  if (!auth.isAuthenticated.value) return
  ws.connect()
  refreshUnread()
  // 收到实时通知即未读数 +1（页面已在通知中心时由该页自行刷新）。
  offNotif = ws.on('notification', () => {
    if (route.path !== '/notifications') unread.value += 1
  })
})

onBeforeUnmount(() => offNotif?.())
</script>

<template>
  <div class="min-h-screen text-text">
    <TopNav @toggle-drawer="drawerOpen = !drawerOpen" />

    <!-- 移动端抽屉遮罩 -->
    <div
      v-if="drawerOpen"
      class="fixed inset-0 z-40 bg-black/50 md:hidden"
      @click="drawerOpen = false"
    />

    <div class="mx-auto flex max-w-app gap-5 px-4 py-5">
      <!-- 左侧频道栏：桌面/平板常驻(粘顶，不随主内容高度变化而跳动)，移动端抽屉 -->
      <aside
        class="glass-card fixed inset-y-0 left-0 z-50 w-sidebar -translate-x-full transition-transform duration-std ease-out md:static md:z-auto md:block md:w-sidebar md:translate-x-0 md:self-start md:sticky md:top-5 md:max-h-[calc(100vh-2.5rem)] md:overflow-y-auto"
        :class="{ 'translate-x-0': drawerOpen }"
      >
        <ChannelSidebar />
      </aside>

      <!-- 主内容 -->
      <main class="min-w-0 flex-1">
        <slot />
      </main>

      <!-- 右侧栏：仅桌面(≥1280)显示。首页等非频道页展示全站今日聚合。 -->
      <aside class="hidden w-aside flex-shrink-0 flex-col gap-5 xl:flex">
        <EmotionBoard :global="globalBoard" />
        <WarmestBoard />
      </aside>
    </div>
  </div>
</template>
