<template>
  <header class="top-navbar">
    <!-- Logo -->
    <div class="navbar-logo" @click="$router.push('/')">
      <div class="logo-icon">
        <span class="material-icons">favorite</span>
      </div>
      <span class="logo-text">Alike</span>
    </div>

    <!-- 搜索框 (桌面端显示) -->
    <div class="navbar-search">
      <span class="material-icons search-icon">search</span>
      <input type="text" placeholder="搜索用户、消息..." />
    </div>

    <!-- 右侧操作区 -->
    <div class="navbar-actions">
      <button class="action-btn" aria-label="消息" @click="$router.push('/messages')">
        <span class="material-icons">mail</span>
        <span class="badge" v-if="unreadCount > 0">{{ unreadCount > 99 ? '99+' : unreadCount }}</span>
      </button>
      <button class="action-btn" aria-label="通知" @click="$router.push('/notifications')">
        <span class="material-icons">notifications</span>
        <span class="badge" v-if="notificationCount > 0">{{ notificationCount > 99 ? '99+' : notificationCount }}</span>
      </button>
      <button class="user-btn" @click="$router.push('/profile')" :title="userName">
        <img :src="userAvatar" :alt="userName" />
      </button>
    </div>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

// 用户信息
const userName = computed(() => userStore.user?.nickname || '用户')
const userAvatar = computed(() => userStore.user?.avatar_url || '/default-avatar.png')

// 未读数量 (示例数据，后续从API获取)
const unreadCount = ref(3)
const notificationCount = ref(5)
</script>

<style scoped>
.top-navbar {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

/* Logo */
.navbar-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  user-select: none;
}

.logo-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(240, 147, 251, 0.4);
}

.logo-icon .material-icons {
  font-size: 24px;
  color: white;
}

.logo-text {
  font-size: 22px;
  font-weight: 800;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.5px;
}

/* 搜索框 */
.navbar-search {
  flex: 1;
  max-width: 400px;
  margin: 0 40px;
  position: relative;
  display: none;
}

@media (min-width: 768px) {
  .navbar-search {
    display: block;
  }
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 20px;
  color: #9ca3af;
  pointer-events: none;
}

.navbar-search input {
  width: 100%;
  height: 40px;
  padding: 0 40px 0 40px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 10px;
  background: rgba(249, 250, 251, 0.8);
  font-size: 14px;
  color: #1a1a1a;
  transition: all 0.2s ease;
}

.navbar-search input:focus {
  outline: none;
  border-color: #667eea;
  background: white;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.navbar-search input::placeholder {
  color: #9ca3af;
}

/* 右侧操作区 */
.navbar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-btn {
  width: 42px;
  height: 42px;
  border: none;
  background: transparent;
  border-radius: 10px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  transition: all 0.2s ease;
}

.action-btn:hover {
  background: rgba(102, 126, 234, 0.1);
}

.action-btn .material-icons {
  font-size: 22px;
  color: #6b7280;
}

.action-btn:hover .material-icons {
  color: #667eea;
}

.badge {
  position: absolute;
  top: 4px;
  right: 4px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  background: #ef4444;
  color: white;
  font-size: 11px;
  font-weight: 700;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

/* 用户头像按钮 */
.user-btn {
  width: 42px;
  height: 42px;
  border: 2px solid transparent;
  border-radius: 10px;
  cursor: pointer;
  padding: 0;
  overflow: hidden;
  transition: all 0.2s ease;
}

.user-btn:hover {
  border-color: #667eea;
  transform: scale(1.05);
}

.user-btn img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* 响应式 */
@media (max-width: 767px) {
  .top-navbar {
    padding: 0 16px;
    height: 56px;
  }

  .logo-icon {
    width: 36px;
    height: 36px;
  }

  .logo-icon .material-icons {
    font-size: 20px;
  }

  .logo-text {
    font-size: 18px;
  }

  .action-btn,
  .user-btn {
    width: 38px;
    height: 38px;
  }

  .action-btn .material-icons {
    font-size: 20px;
  }
}
</style>
