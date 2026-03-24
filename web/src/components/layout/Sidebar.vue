<template>
  <aside class="sidebar">
    <!-- 用户卡片 -->
    <div class="user-card">
      <div class="user-avatar">
        <img :src="user.avatar" :alt="user.nickname" />
        <span class="online-status"></span>
      </div>
      <div class="user-info">
        <div class="user-name">{{ user.nickname }}</div>
        <div class="user-status">● 在线</div>
      </div>
    </div>

    <!-- 导航菜单 -->
    <nav class="sidebar-nav">
      <router-link
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        class="nav-item"
        :class="{ active: isActive(item.path) }"
      >
        <span class="material-icons nav-icon">{{ item.icon }}</span>
        <span class="nav-label">{{ item.label }}</span>
        <span v-if="item.badge > 0" class="nav-badge">{{ item.badge }}</span>
      </router-link>
    </nav>
  </aside>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const userStore = useUserStore()

// 用户信息
const user = computed(() => ({
  nickname: userStore.user?.nickname || '用户',
  avatar: userStore.user?.avatar_url || '/default-avatar.svg'
}))

// 导航菜单项
const navItems = computed(() => [
  { path: '/home', label: '首页', icon: 'home', badge: 0 },
  { path: '/global', label: '聊天', icon: 'chat', badge: 0 },
  { path: '/match', label: '匹配', icon: 'favorite', badge: 0 },
  { path: '/messages', label: '消息', icon: 'mail', badge: 3 },
  { path: '/notifications', label: '通知', icon: 'notifications', badge: 5 },
  { path: '/profile', label: '我的', icon: 'person', badge: 0 }
])

// 判断是否激活
const isActive = (path) => {
  return route.path === path || route.path.startsWith(path + '/')
}
</script>

<style scoped>
.sidebar {
  width: 240px;
  background: white;
  border-right: 1px solid rgba(0, 0, 0, 0.06);
  display: flex;
  flex-direction: column;
  position: fixed;
  left: 0;
  top: 64px;
  bottom: 0;
  z-index: 50;
}

/* 用户卡片 */
.user-card {
  padding: 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  position: relative;
  width: 48px;
  height: 48px;
  flex-shrink: 0;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  border-radius: 12px;
  object-fit: cover;
}

.online-status {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 14px;
  height: 14px;
  background: #22c55e;
  border: 2px solid white;
  border-radius: 50%;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 15px;
  font-weight: 600;
  color: #1a1a1a;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-status {
  font-size: 13px;
  color: #22c55e;
  margin-top: 2px;
}

/* 导航菜单 */
.sidebar-nav {
  flex: 1;
  padding: 16px 12px;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 10px;
  color: #6b7280;
  text-decoration: none;
  transition: all 0.2s ease;
  position: relative;
  margin-bottom: 4px;
}

.nav-item:hover {
  background: rgba(102, 126, 234, 0.08);
  color: #667eea;
}

.nav-item.active {
  background: rgba(102, 126, 234, 0.15);
  color: #667eea;
}

/* 左侧边条高亮 */
.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 24px;
  background: #667eea;
  border-radius: 0 2px 2px 0;
}

.nav-icon {
  font-size: 22px;
  flex-shrink: 0;
}

.nav-label {
  font-size: 15px;
  font-weight: 500;
  flex: 1;
}

.nav-badge {
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  background: #ef4444;
  color: white;
  font-size: 11px;
  font-weight: 700;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

/* 滚动条样式 */
.sidebar-nav::-webkit-scrollbar {
  width: 4px;
}

.sidebar-nav::-webkit-scrollbar-track {
  background: transparent;
}

.sidebar-nav::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 2px;
}

.sidebar-nav::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.2);
}
</style>
