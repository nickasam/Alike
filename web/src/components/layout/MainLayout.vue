<template>
  <div class="main-layout">
    <!-- 顶部导航栏 -->
    <TopNavBar />

    <!-- 主容器 -->
    <div class="main-container">
      <!-- 侧边栏 (桌面端) -->
      <Sidebar v-if="isDesktop" />

      <!-- 主内容区域 -->
      <main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <!-- 底部Tab栏 (移动端) -->
    <BottomTabBar v-if="!isDesktop" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import TopNavBar from './TopNavBar.vue'
import Sidebar from './Sidebar.vue'
import BottomTabBar from './BottomTabBar.vue'

// 响应式断点检测
const windowWidth = ref(window.innerWidth)

const isDesktop = computed(() => windowWidth.value >= 1024)

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
  background: #f5f7fa;
  display: flex;
  flex-direction: column;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

.main-container {
  display: flex;
  flex: 1;
  padding-top: 64px; /* 桌面端顶部导航栏高度 */
}

@media (max-width: 1023px) {
  .main-container {
    padding-top: 56px; /* 移动端顶部导航栏高度 */
    padding-bottom: 64px; /* 底部Tab栏高度 */
  }
}

.main-content {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

/* 页面切换过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
