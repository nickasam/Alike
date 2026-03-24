<template>
  <div class="user-card" :class="{ compact: compact }">
    <div class="user-avatar" @click="$emit('click')">
      <img :src="user.avatar" :alt="user.nickname" />
      <span v-if="showStatus" class="online-status" :class="statusClass"></span>
    </div>
    <div v-if="!compact" class="user-info">
      <div class="user-name">{{ user.nickname }}</div>
      <div v-if="showStatus" class="user-status">{{ statusText }}</div>
    </div>
    <slot></slot>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  user: {
    type: Object,
    required: true,
    default: () => ({
      nickname: '用户',
      avatar: '/default-avatar.png'
    })
  },
  compact: {
    type: Boolean,
    default: false
  },
  showStatus: {
    type: Boolean,
    default: true
  },
  online: {
    type: Boolean,
    default: true
  }
})

defineEmits(['click'])

const statusClass = computed(() => ({
  'status-online': props.online,
  'status-offline': !props.online
}))

const statusText = computed(() => props.online ? '● 在线' : '○ 离线')
</script>

<style scoped>
.user-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 12px;
  transition: background 0.2s ease;
}

.user-card:hover {
  background: rgba(0, 0, 0, 0.03);
}

.user-card.compact {
  padding: 0;
}

.user-avatar {
  position: relative;
  width: 48px;
  height: 48px;
  flex-shrink: 0;
  cursor: pointer;
  border-radius: 12px;
  overflow: hidden;
}

.user-card.compact .user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 10px;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.online-status {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 14px;
  height: 14px;
  border: 2px solid white;
  border-radius: 50%;
}

.status-online {
  background: #22c55e;
}

.status-offline {
  background: #9ca3af;
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
  margin-top: 2px;
}

.status-online .user-status {
  color: #22c55e;
}

.status-offline .user-status {
  color: #9ca3af;
}
</style>
