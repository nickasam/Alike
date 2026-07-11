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
  <div class="flex h-screen flex-col text-text">
    <TopNav @toggle-drawer="drawerOpen = !drawerOpen" />

    <!-- 移动端抽屉遮罩 -->
    <div
      v-if="drawerOpen"
      class="fixed inset-0 z-40 bg-black/50 md:hidden"
      @click="drawerOpen = false"
    />

    <!-- 满屏三栏 shell：顶栏下方占满剩余高度，侧栏/右栏贴边毛玻璃 -->
    <div class="flex min-h-0 flex-1">
      <!-- 左侧频道栏：贴边毛玻璃 + 右边框；移动端抽屉 -->
      <aside
        class="fixed inset-y-0 left-0 z-50 w-sidebar -translate-x-full overflow-y-auto border-r border-border bg-surface backdrop-blur-glass transition-transform duration-std ease-out md:static md:z-auto md:block md:w-sidebar md:translate-x-0"
        :class="{ 'translate-x-0': drawerOpen }"
      >
        <ChannelSidebar />
      </aside>

      <!-- 主内容：可滚动 -->
      <main class="min-w-0 flex-1 overflow-y-auto px-6 py-6">
        <div class="mx-auto max-w-content-wide">
          <slot />
        </div>
      </main>

      <!-- 右侧栏：贴边毛玻璃 + 左边框，仅桌面(≥1280)显示 -->
      <aside
        class="hidden w-aside flex-shrink-0 overflow-y-auto border-l border-border bg-surface px-4 py-5 backdrop-blur-glass xl:block"
      >
        <div class="flex flex-col gap-5">
          <EmotionBoard :global="globalBoard" />
          <WarmestBoard />
        </div>
      </aside>
    </div>
  </div>
</template>
