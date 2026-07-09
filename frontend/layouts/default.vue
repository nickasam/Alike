<script setup lang="ts">
/**
 * default 布局 — 三列响应式骨架。
 * 断点（交互规范 §1.1）：
 *   桌面 ≥1280(xl)  三列：侧边栏 260 + 主内容 1fr + 右侧栏 320
 *   平板 768-1279(md) 两列：侧边栏 + 主内容，右侧栏隐藏
 *   手机 <768        单列：侧边栏转抽屉，主内容全宽
 */
const drawerOpen = ref(false)
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
      <!-- 左侧频道栏：桌面/平板常驻，移动端抽屉 -->
      <aside
        class="glass-card fixed inset-y-0 left-0 z-50 w-sidebar -translate-x-full transition-transform duration-std ease-out md:static md:z-auto md:block md:w-sidebar md:translate-x-0"
        :class="{ 'translate-x-0': drawerOpen }"
      >
        <ChannelSidebar />
      </aside>

      <!-- 主内容 -->
      <main class="min-w-0 flex-1">
        <slot />
      </main>

      <!-- 右侧栏：仅桌面(≥1280)显示 -->
      <aside class="hidden w-aside flex-shrink-0 flex-col gap-5 xl:flex">
        <EmotionBoard />
        <div class="glass-card p-4">
          <h3 class="mb-2 flex items-center gap-2 text-md font-semibold">
            <AppIcon name="trophy" :size="18" />
            今日牛马榜
          </h3>
          <p class="text-sm text-mute">榜单占位（阶段八实现）</p>
        </div>
      </aside>
    </div>
  </div>
</template>
