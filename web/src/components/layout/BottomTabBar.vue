<template>
  <nav class="bottom-tabbar">
    <router-link
      v-for="item in tabItems"
      :key="item.path"
      :to="item.path"
      class="tab-item"
      :class="{ active: isActive(item.path) }"
    >
      <span class="material-icons tab-icon">{{ item.icon }}</span>
      <span class="tab-label">{{ item.label }}</span>
      <span v-if="item.badge > 0" class="tab-badge">{{ item.badge }}</span>
    </router-link>
  </nav>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// 底部Tab项 (不包含通知，通知在"我的"页面内)
const tabItems = computed(() => [
  { path: '/home', label: '首页', icon: 'home', badge: 0 },
  { path: '/global', label: '聊天', icon: 'chat', badge: 0 },
  { path: '/match', label: '匹配', icon: 'favorite', badge: 0 },
  { path: '/messages', label: '消息', icon: 'mail', badge: 3 },
  { path: '/profile', label: '我的', icon: 'person', badge: 0 }
])

// 判断是否激活
const isActive = (path) => {
  return route.path === path || route.path.startsWith(path + '/')
}
</script>

<style scoped>
.bottom-tabbar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-top: 1px solid rgba(0, 0, 0, 0.08);
  display: flex;
  justify-content: space-around;
  align-items: center;
  z-index: 100;
  padding-bottom: env(safe-area-inset-bottom);
}

.tab-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 8px 4px;
  text-decoration: none;
  color: #9ca3af;
  position: relative;
  transition: all 0.2s ease;
  min-width: 0;
}

.tab-item:active {
  transform: scale(0.95);
}

.tab-item.active {
  color: #667eea;
}

.tab-icon {
  font-size: 24px;
  margin-bottom: 4px;
  transition: transform 0.2s ease;
}

.tab-item.active .tab-icon {
  transform: scale(1.1);
}

.tab-label {
  font-size: 11px;
  font-weight: 600;
  line-height: 1;
}

.tab-badge {
  position: absolute;
  top: 6px;
  right: 20%;
  min-width: 16px;
  height: 16px;
  padding: 0 4px;
  background: #ef4444;
  color: white;
  font-size: 10px;
  font-weight: 700;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

/* 小屏幕适配 */
@media (max-width: 374px) {
  .tab-icon {
    font-size: 22px;
  }

  .tab-label {
    font-size: 10px;
  }
}

/* 安全区域适配 */
@supports (padding: env(safe-area-inset-bottom)) {
  .bottom-tabbar {
    padding-bottom: calc(8px + env(safe-area-inset-bottom));
  }
}
</style>
