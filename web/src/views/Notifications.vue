<template>
  <div class="notifications-page">
    <div class="page-header">
      <h1 class="page-title">通知</h1>
    </div>

    <div class="notifications-list">
      <div v-if="notifications.length === 0" class="empty-state">
        <span class="material-icons empty-icon">notifications_none</span>
        <p>暂无通知</p>
      </div>

      <div
        v-for="notification in notifications"
        :key="notification.id"
        class="notification-item"
        :class="{ unread: !notification.read }"
      >
        <div class="notification-icon" :class="`type-${notification.type}`">
          <span class="material-icons">{{ notification.icon }}</span>
        </div>
        <div class="notification-content">
          <div class="notification-title">{{ notification.title }}</div>
          <div class="notification-text">{{ notification.text }}</div>
          <div class="notification-time">{{ notification.time }}</div>
        </div>
        <div v-if="!notification.read" class="unread-dot"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

// 示例数据
const notifications = ref([
  {
    id: 1,
    type: 'like',
    icon: 'favorite',
    title: '新的喜欢',
    text: '小明 喜欢了你',
    time: '5分钟前',
    read: false
  },
  {
    id: 2,
    type: 'match',
    icon: 'people',
    title: '新的匹配',
    text: '你与 Alex 匹配成功',
    time: '1小时前',
    read: false
  },
  {
    id: 3,
    type: 'message',
    icon: 'mail',
    title: '新消息',
    text: 'Kevin 给你发了一条消息',
    time: '2小时前',
    read: true
  },
  {
    id: 4,
    type: 'view',
    icon: 'visibility',
    title: '谁看过我',
    text: 'Mike 查看了你的资料',
    time: '昨天',
    read: true
  }
])
</script>

<style scoped>
.notifications-page {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a1a;
}

.notifications-list {
  background: white;
  border-radius: 16px;
  overflow: hidden;
}

.empty-state {
  padding: 60px 20px;
  text-align: center;
  color: #9ca3af;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.notification-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  position: relative;
  transition: background 0.2s ease;
}

.notification-item:last-child {
  border-bottom: none;
}

.notification-item:hover {
  background: rgba(102, 126, 234, 0.03);
}

.notification-item.unread {
  background: rgba(102, 126, 234, 0.05);
}

.notification-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.notification-icon .material-icons {
  font-size: 24px;
  color: white;
}

.type-like { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.type-match { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.type-message { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
.type-view { background: #22c55e; }

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 15px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.notification-text {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 4px;
}

.notification-time {
  font-size: 12px;
  color: #9ca3af;
}

.unread-dot {
  width: 10px;
  height: 10px;
  background: #ef4444;
  border-radius: 50%;
  flex-shrink: 0;
}

@media (max-width: 767px) {
  .notifications-page {
    padding: 16px;
  }

  .page-title {
    font-size: 24px;
  }

  .notification-item {
    padding: 12px 16px;
    gap: 12px;
  }

  .notification-icon {
    width: 40px;
    height: 40px;
  }
}
</style>
