<template>
  <div class="global-chat-container">
    <!-- 登录界面 -->
    <div class="login-container" v-if="!isLoggedIn">
      <div class="login-header">
        <div class="login-logo-wrapper">
          <div class="login-logo-glow"></div>
          <div class="login-logo">
            <span class="material-icons">favorite</span>
          </div>
        </div>
        <div class="login-title">Alike</div>
        <div class="login-subtitle">相似灵魂的相遇</div>
      </div>
      <form class="login-form" @submit.prevent="handleLogin">
        <div class="form-group">
          <input 
            type="tel" 
            class="form-input" 
            v-model="loginForm.phone" 
            placeholder="手机号" 
            required
          >
        </div>
        <div class="form-group">
          <input 
            type="password" 
            class="form-input" 
            v-model="loginForm.password" 
            placeholder="密码" 
            required
          >
        </div>
        <button type="submit" class="login-btn" :disabled="isLoading">
          {{ isLoading ? '登录中...' : '进入聊天室' }}
        </button>
        <div class="error-msg" v-if="errorMessage">{{ errorMessage }}</div>
      </form>
    </div>

    <!-- 主应用 -->
    <div class="app-container" v-else>
      <!-- 左侧边栏 -->
      <aside class="sidebar" :class="{ open: sidebarOpen }" id="sidebar">
        <div class="sidebar-header">
          <div class="sidebar-title">
            频道
            <span class="channel-count">{{ channels.length }}</span>
          </div>

          <div 
            v-for="channel in channels" 
            :key="channel.id"
            class="channel-item"
            :class="{ active: currentChannelId === channel.id }"
            @click="switchChannel(channel)"
          >
            <span class="channel-icon">#</span>
            <span class="channel-name">{{ channel.name }}</span>
            <span class="channel-badge">{{ channel.badge }}</span>
          </div>
        </div>

        <div class="users-section">
          <div class="section-title">在线用户 — {{ formatNumber(onlineCount) }}</div>

          <div 
            v-for="user in onlineUsers" 
            :key="user.id"
            class="user-item"
            @click="mentionUser(user)"
          >
            <div 
              class="user-item-avatar"
              :style="{ background: user.gradient || getGradientForUser(user.nickname) }"
            >
              {{ user.nickname ? user.nickname[0].toUpperCase() : '?' }}
              <span 
                class="status"
                :class="{
                  idle: user.status === 'idle',
                  dnd: user.status === 'dnd'
                }"
              ></span>
            </div>
            <div class="user-item-info">
              <div class="user-item-name">{{ user.nickname }}</div>
              <div class="user-item-status">{{ user.statusText || '在线' }}</div>
            </div>
          </div>
        </div>
      </aside>

      <!-- 主内容区 -->
      <main class="main-content" @click="closeSidebarOnMobile">
        <div class="chat-header">
          <h1># {{ currentChannelName }}</h1>
          <p>与 {{ formatNumber(onlineCount) }} 位在线用户分享你的故事 💜</p>
        </div>

        <div class="messages-container" ref="messagesContainer">
          <!-- 空状态 -->
          <div class="empty-state" v-if="messages.length === 0">
            <div class="empty-icon">💬</div>
            <div class="empty-text">暂无消息</div>
            <div class="empty-hint">成为第一个发言的人吧！</div>
          </div>

          <!-- 消息列表 -->
          <div 
            v-for="msg in sortedMessages" 
            :key="msg.id"
            class="message-group"
            :class="{ 'message-own': msg.user_id === currentUserId }"
          >
            <div class="message-header">
              <div 
                class="message-avatar"
                :style="{ 
                  background: msg.user_id === currentUserId 
                    ? 'var(--gradient-brand)' 
                    : getGradientForUser(msg.username) 
                }"
              >
                {{ msg.username ? msg.username[0].toUpperCase() : '?' }}
              </div>
              <div class="message-info">
                <div class="message-name">{{ msg.username }}</div>
                <div class="message-time">{{ formatTime(msg.created_at) }}</div>
              </div>
            </div>
            <div class="message-content">{{ msg.content }}</div>
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="input-area">
          <div class="input-wrapper">
            <div class="input-container">
              <textarea
                class="message-input"
                v-model="newMessage"
                placeholder="输入消息... (Enter发送，Shift+Enter换行)"
                rows="1"
                @keydown.exact.enter="sendMessage"
                @input="autoResize"
                ref="messageInput"
              ></textarea>
              <div class="input-actions">
                <button 
                  class="action-btn" 
                  title="表情" 
                  @click="showEmojiPicker = !showEmojiPicker"
                  :class="{ active: showEmojiPicker }"
                >
                  😊
                </button>
                <button 
                  class="action-btn" 
                  title="图片"
                  @click="handleImageUpload"
                >
                  📷
                </button>
                <button 
                  class="action-btn" 
                  title="文件"
                  @click="handleFileUpload"
                >
                  📎
                </button>
              </div>

              <!-- 表情选择器 -->
              <div class="emoji-picker" v-if="showEmojiPicker" v-click-outside="() => showEmojiPicker = false">
                <div class="emoji-grid">
                  <span 
                    v-for="emoji in emojis" 
                    :key="emoji"
                    class="emoji-item"
                    @click="insertEmoji(emoji)"
                  >
                    {{ emoji }}
                  </span>
                </div>
              </div>
            </div>
            <button class="send-btn" @click="sendMessage" :disabled="!newMessage.trim()">
              ➤
            </button>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useGlobalChatStore } from '@/stores/global'

const router = useRouter()
const userStore = useUserStore()
const globalStore = useGlobalChatStore()

// 状态
const isLoggedIn = ref(false)
const isLoading = ref(false)
const errorMessage = ref('')
const loginForm = ref({
  phone: '',
  password: ''
})

const currentUserId = ref(null)
const currentUsername = ref('')
const messages = ref([])
const newMessage = ref('')
const onlineCount = ref(1234)
const sidebarOpen = ref(false)
const onlineUsers = ref([])
const showEmojiPicker = ref(false)
const currentView = ref('chat')
const currentChannelId = ref(1)

const messagesContainer = ref(null)
const messageInput = ref(null)

let refreshInterval = null

// 频道数据
const channels = ref([
  { id: 1, name: '全局聊天室', badge: '1.2k' },
  { id: 2, name: '技术交流', badge: '356' },
  { id: 3, name: '音乐分享', badge: '189' }
])

const currentChannelName = computed(() => {
  const channel = channels.value.find(c => c.id === currentChannelId.value)
  return channel ? channel.name : '全局聊天室'
})

// 表情列表
const emojis = [
  '😀', '😃', '😄', '😁', '😅', '😂', '🤣', '😊',
  '😇', '🙂', '🙃', '😉', '😌', '😍', '🥰', '😘',
  '😗', '😙', '😚', '😋', '😛', '😝', '😜', '🤪',
  '🤨', '🧐', '🤓', '😎', '🤩', '🥳', '😏', '😒',
  '😞', '😔', '😟', '😕', '🙁', '😣', '😖', '😫',
  '😩', '🥺', '😢', '😭', '😤', '😠', '😡', '🤬',
  '👍', '👎', '👌', '✌️', '🤞', '🤝', '🙏', '💪',
  '❤️', '💕', '💖', '💗', '💙', '💚', '💛', '🧡',
  '🎉', '🎊', '🎈', '🎁', '🏆', '🥇', '⭐', '✨',
  '🌟', '💫', '🔥', '💯', '✅', '❌', '⚡', '🎯'
]

// 计算属性
const userInitial = computed(() => {
  return currentUsername.value ? currentUsername.value[0].toUpperCase() : '?'
})

const sortedMessages = computed(() => {
  return [...messages.value].sort((a, b) => 
    new Date(a.created_at) - new Date(b.created_at)
  )
})

// 渐变色生成器
const gradients = [
  'linear-gradient(135deg, #06b6d4 0%, #3b82f6 100%)',
  'linear-gradient(135deg, #22c55e 0%, #06b6d4 100%)',
  'linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)',
  'linear-gradient(135deg, #ec4899 0%, #8b5cf6 100%)',
  'linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%)',
  'linear-gradient(135deg, #14b8a6 0%, #22c55e 100%)',
  'linear-gradient(135deg, #f97316 0%, #f59e0b 100%)',
  'linear-gradient(135deg, #8b5cf6 0%, #ec4899 100%)'
]

const getGradientForUser = (username) => {
  if (!username) return gradients[0]
  const index = username.charCodeAt(0) % gradients.length
  return gradients[index]
}

// 方法
const checkAuth = () => {
  const token = localStorage.getItem('alike_access_token')
  const userId = localStorage.getItem('alike_user_id')
  const username = localStorage.getItem('alike_username')
  
  if (token && userId) {
    currentUserId.value = userId
    currentUsername.value = username || ''
    isLoggedIn.value = true
    loadMessages()
    startAutoRefresh()
  }
}

const handleLogin = async () => {
  try {
    isLoading.value = true
    errorMessage.value = ''
    
    let result = await userStore.login({
      phone: loginForm.value.phone,
      password: loginForm.value.password
    })
    
    if (result.success) {
      currentUserId.value = userStore.user.id
      currentUsername.value = userStore.user.nickname
      isLoggedIn.value = true
      loadMessages()
      startAutoRefresh()
    } else {
      result = await userStore.register({
        phone: loginForm.value.phone,
        password: loginForm.value.password,
        nickname: loginForm.value.phone.slice(-4) + '用户'
      })
      
      if (result.success) {
        currentUserId.value = userStore.user.id
        currentUsername.value = userStore.user.nickname
        isLoggedIn.value = true
        loadMessages()
        startAutoRefresh()
      } else {
        errorMessage.value = result.message || '登录失败'
      }
    }
  } catch (error) {
    errorMessage.value = '网络错误：' + error.message
  } finally {
    isLoading.value = false
  }
}

const logout = () => {
  if (confirm('确定要退出吗？')) {
    userStore.logout()
    isLoggedIn.value = false
    messages.value = []
    if (refreshInterval) {
      clearInterval(refreshInterval)
      refreshInterval = null
    }
  }
}

const loadMessages = async () => {
  try {
    const result = await globalStore.fetchMessages()
    if (result.success && result.data) {
      messages.value = result.data
      await nextTick()
      scrollToBottom()
    }
  } catch (error) {
    console.error('加载消息失败:', error)
  }
}

const sendMessage = async () => {
  const content = newMessage.value.trim()
  if (!content) return

  try {
    const result = await globalStore.sendMessage({ content })
    if (result.success) {
      newMessage.value = ''
      if (messageInput.value) {
        messageInput.value.style.height = 'auto'
      }
      loadMessages()
    } else {
      showToast('发送失败：' + (result.message || '未知错误'))
    }
  } catch (error) {
    showToast('网络错误')
  }
}

const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const formatTime = (dateStr) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now - date

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  
  const hours = date.getHours().toString().padStart(2, '0')
  const minutes = date.getMinutes().toString().padStart(2, '0')
  return `今天 ${hours}:${minutes}`
}

const formatNumber = (num) => {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
}

const closeSidebarOnMobile = () => {
  if (window.innerWidth <= 1024) {
    sidebarOpen.value = false
  }
}

const switchView = (view) => {
  currentView.value = view
  if (view === 'chat') {
    console.log('切换到聊天')
  } else if (view === 'nearby') {
    console.log('切换到附近的人')
    router.push('/nearby')
  } else if (view === 'profile') {
    console.log('切换到个人中心')
    router.push('/profile')
  }
}

const switchChannel = (channel) => {
  currentChannelId.value = channel.id
  console.log('切换到频道:', channel.name)
}

const mentionUser = (user) => {
  const mention = `@${user.nickname} `
  newMessage.value += mention
  if (messageInput.value) {
    messageInput.value.focus()
  }
}

const showUserProfile = () => {
  console.log('显示用户资料')
  switchView('profile')
}

const autoResize = (event) => {
  const textarea = event.target
  textarea.style.height = 'auto'
  textarea.style.height = Math.min(textarea.scrollHeight, 120) + 'px'
}

const insertEmoji = (emoji) => {
  newMessage.value += emoji
  if (messageInput.value) {
    messageInput.value.focus()
  }
  showEmojiPicker.value = false
}

const handleImageUpload = () => {
  showToast('图片上传功能开发中...')
}

const handleFileUpload = () => {
  showToast('文件上传功能开发中...')
}

const loadOnlineUsers = async () => {
  try {
    const result = await globalStore.fetchOnlineUsers()
    if (result.success && result.data) {
      onlineUsers.value = result.data.map(user => ({
        ...user,
        status: ['online', 'idle', 'dnd'][Math.floor(Math.random() * 3)],
        statusText: ['在线', '离开 - 15分钟前', '忙碌'][Math.floor(Math.random() * 3)]
      }))
      onlineCount.value = result.data.length
    }
  } catch (error) {
    console.error('加载在线用户失败:', error)
  }
}

const startAutoRefresh = () => {
  if (refreshInterval) return
  refreshInterval = setInterval(() => {
    loadMessages()
    loadOnlineUsers()
  }, 3000)
}

const showToast = (message) => {
  alert(message)
}

const handleClickOutside = (event) => {
  if (showEmojiPicker.value && !event.target.closest('.emoji-picker') && !event.target.closest('.action-btn')) {
    showEmojiPicker.value = false
  }
}

onMounted(() => {
  checkAuth()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
/* ========== CSS 变量 ========== */
:root {
  /* Alike品牌色 */
  --primary: #8b5cf6;
  --primary-light: #a78bfa;
  --primary-dark: #7c3aed;
  --secondary: #ec4899;
  --accent: #06b6d4;

  /* 渐变 */
  --gradient-brand: linear-gradient(135deg, #8b5cf6 0%, #ec4899 100%);
  --gradient-accent: linear-gradient(135deg, #06b6d4 0%, #8b5cf6 100%);

  /* 深色主题 */
  --bg-primary: #0f0f0f;
  --bg-secondary: #1a1a1a;
  --bg-tertiary: #262626;
  --bg-card: #1f1f1f;

  /* 文字色 */
  --text-primary: #ffffff;
  --text-secondary: #a1a1aa;
  --text-tertiary: #71717a;

  /* 边框和分割线 */
  --border-color: #262626;
  --divider-color: #2a2a2a;

  /* 阴影 */
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.3);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.4);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.5);
  --shadow-brand: 0 4px 16px rgba(139, 92, 246, 0.3);
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.global-chat-container {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  background: var(--bg-primary);
  color: var(--text-primary);
  min-height: 100vh;
}

/* ========== 登录界面 ========== */
.login-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: var(--bg-primary);
  z-index: 9999;
  padding: 20px;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.login-logo-wrapper {
  position: relative;
  display: inline-block;
  margin-bottom: 24px;
}

.login-logo {
  width: 100px;
  height: 100px;
  background: var(--gradient-brand);
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-brand);
  position: relative;
  z-index: 1;
  animation: brandPulse 3s ease-in-out infinite;
}

.login-logo .material-icons {
  font-size: 48px;
  color: white;
}

@keyframes brandPulse {
  0%, 100% {
    box-shadow: var(--shadow-brand);
  }
  50% {
    box-shadow: 0 6px 20px rgba(139, 92, 246, 0.5);
  }
}

.login-logo-glow {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 24px;
  background: var(--gradient-brand);
  filter: blur(40px);
  opacity: 0.4;
  top: 0;
  left: 0;
}

.login-title {
  font-size: 42px;
  font-weight: 900;
  background: var(--gradient-brand);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 8px;
}

.login-subtitle {
  font-size: 16px;
  color: var(--text-tertiary);
  font-weight: 500;
}

.login-form {
  width: 100%;
  max-width: 400px;
}

.form-group {
  margin-bottom: 16px;
}

.form-input {
  width: 100%;
  padding: 16px 20px;
  border: 1px solid var(--border-color);
  border-radius: 12px;
  font-size: 15px;
  font-weight: 500;
  background: var(--bg-secondary);
  color: var(--text-primary);
  transition: all 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: var(--primary);
  background: var(--bg-card);
  box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.2);
}

.form-input::placeholder {
  color: var(--text-tertiary);
}

.login-btn {
  width: 100%;
  padding: 16px;
  background: var(--gradient-brand);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-top: 24px;
  box-shadow: var(--shadow-brand);
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(139, 92, 246, 0.4);
}

.login-btn:active:not(:disabled) {
  transform: translateY(0);
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-msg {
  color: #ef4444;
  font-size: 14px;
  margin-top: 16px;
  text-align: center;
  min-height: 24px;
  background: rgba(239, 68, 68, 0.1);
  padding: 12px;
  border-radius: 8px;
  border: 1px solid rgba(239, 68, 68, 0.2);
}

/* ========== 主容器 ========== */
.app-container {
  display: grid;
  grid-template-columns: 280px 1fr;
  height: 100vh;
  width: 100vw;
}

/* ========== 顶部导航栏 ========== */
.navbar {
  grid-column: 1 / -1;
  grid-row: 1;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
}

.sidebar-toggle {
  display: none;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: var(--bg-tertiary);
  border: none;
  color: var(--text-secondary);
  font-size: 18px;
  cursor: pointer;
  transition: all 0.2s ease;
  align-items: center;
  justify-content: center;
}

.sidebar-toggle:hover {
  background: var(--bg-card);
  color: var(--text-primary);
}

.brand-logo {
  width: 40px;
  height: 40px;
  background: var(--gradient-brand);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-brand);
  animation: brandPulse 3s ease-in-out infinite;
}

.brand-logo .material-icons {
  font-size: 22px;
  color: white;
}

.brand-text {
  display: flex;
  flex-direction: column;
}

.brand-name {
  font-size: 18px;
  font-weight: 900;
  background: var(--gradient-brand);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1.2;
}

.brand-slogan {
  font-size: 10px;
  color: var(--text-tertiary);
  font-weight: 500;
  line-height: 1;
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.nav-btn {
  padding: 8px 16px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-size: 13px;
  font-weight: 600;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.nav-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.nav-btn.active {
  background: var(--gradient-brand);
  color: white;
  box-shadow: var(--shadow-brand);
}

.btn-icon {
  font-size: 16px;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--gradient-brand);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.user-avatar:hover {
  transform: scale(1.05);
  box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.3);
}

/* ========== 左侧边栏 ========== */
.sidebar {
  grid-column: 1;
  grid-row: 1;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: transform 0.3s ease;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid var(--border-color);
}

.sidebar-title {
  font-size: 12px;
  font-weight: 700;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.channel-count {
  font-size: 11px;
  color: var(--primary);
  background: rgba(139, 92, 246, 0.15);
  padding: 2px 8px;
  border-radius: 10px;
}

.channel-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 4px;
}

.channel-item:hover {
  background: var(--bg-tertiary);
}

.channel-item.active {
  background: rgba(139, 92, 246, 0.15);
  color: var(--primary);
}

.channel-icon {
  font-size: 18px;
  color: var(--text-tertiary);
}

.channel-item.active .channel-icon {
  color: var(--primary);
}

.channel-name {
  flex: 1;
  font-size: 14px;
  font-weight: 600;
}

.channel-badge {
  background: var(--primary);
  color: white;
  font-size: 11px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 10px;
  min-width: 18px;
  text-align: center;
}

.users-section {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.section-title {
  font-size: 12px;
  font-weight: 700;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 12px;
  padding: 0 4px;
}

.user-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 2px;
}

.user-item:hover {
  background: var(--bg-tertiary);
}

.user-item-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--gradient-brand);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 14px;
  flex-shrink: 0;
  position: relative;
}

.user-item-avatar .status {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 10px;
  height: 10px;
  background: #22c55e;
  border: 2px solid var(--bg-secondary);
  border-radius: 50%;
}

.user-item-avatar .status.idle {
  background: #f59e0b;
}

.user-item-avatar .status.dnd {
  background: #ef4444;
}

.user-item-info {
  flex: 1;
  min-width: 0;
}

.user-item-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-item-status {
  font-size: 12px;
  color: var(--text-tertiary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* ========== 主内容区 ========== */
.main-content {
  grid-column: 2;
  grid-row: 1;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow: hidden;
}

.chat-header {
  padding: 16px 24px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-secondary);
}

.chat-header h1 {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.chat-header p {
  font-size: 13px;
  color: var(--text-tertiary);
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* 空状态 */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 20px;
  opacity: 0.5;
}

.empty-text {
  font-size: 18px;
  margin-bottom: 8px;
  color: var(--text-primary);
  font-weight: 700;
}

.empty-hint {
  font-size: 14px;
  color: var(--text-tertiary);
}

/* 消息组 */
.message-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-group:hover {
  background: rgba(255, 255, 255, 0.02);
  margin: -8px -12px;
  padding: 8px 12px;
  border-radius: 8px;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.message-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--gradient-brand);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 16px;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.message-info {
  flex: 1;
}

.message-name {
  font-size: 14px;
  font-weight: 700;
  color: var(--text-primary);
}

.message-time {
  font-size: 12px;
  color: var(--text-tertiary);
  font-weight: 500;
}

.message-content {
  margin-left: 52px;
  font-size: 15px;
  line-height: 1.5;
  color: var(--text-secondary);
  word-wrap: break-word;
}

.message-own {
  flex-direction: row-reverse;
}

.message-own .message-header {
  flex-direction: row-reverse;
}

.message-own .message-content {
  margin-left: 0;
  margin-right: 52px;
  text-align: right;
}

.message-own .message-name {
  color: var(--primary);
}

/* ========== 输入区域 ========== */
.input-area {
  padding: 16px 24px 24px;
  background: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
}

.input-wrapper {
  display: flex;
  gap: 12px;
  align-items: flex-end;
}

.input-container {
  flex: 1;
  position: relative;
}

.message-input {
  width: 100%;
  padding: 12px 50px 12px 16px;
  border: 1px solid var(--border-color);
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  resize: none;
  transition: all 0.2s ease;
  background: var(--bg-tertiary);
  color: var(--text-primary);
  max-height: 120px;
  font-family: inherit;
  line-height: 1.5;
}

.message-input:focus {
  outline: none;
  border-color: var(--primary);
  background: var(--bg-card);
  box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.2);
}

.message-input::placeholder {
  color: var(--text-tertiary);
}

.input-actions {
  position: absolute;
  right: 8px;
  bottom: 8px;
  display: flex;
  gap: 4px;
}

.action-btn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: transparent;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 16px;
  color: var(--text-tertiary);
  transition: all 0.2s ease;
}

.action-btn:hover {
  background: var(--bg-tertiary);
  color: var(--primary);
}

.action-btn.active {
  background: var(--bg-tertiary);
  color: var(--primary);
}

.send-btn {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  background: var(--gradient-brand);
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: white;
  font-size: 18px;
  box-shadow: var(--shadow-brand);
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.send-btn:hover:not(:disabled) {
  transform: scale(1.05);
  box-shadow: 0 6px 20px rgba(139, 92, 246, 0.4);
}

.send-btn:active:not(:disabled) {
  transform: scale(0.95);
}

.send-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* ========== 表情选择器 ========== */
.emoji-picker {
  position: absolute;
  bottom: 60px;
  left: 0;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 12px;
  box-shadow: var(--shadow-lg);
  z-index: 1000;
  width: 320px;
  max-height: 300px;
  overflow-y: auto;
}

.emoji-grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 4px;
}

.emoji-item {
  font-size: 20px;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s ease;
  user-select: none;
}

.emoji-item:hover {
  background: var(--bg-tertiary);
  transform: scale(1.1);
}

/* ========== 滚动条美化 ========== */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--text-tertiary);
}

/* ========== 响应式 ========== */
@media (max-width: 1024px) {
  .sidebar {
    position: fixed;
    left: -280px;
    top: 0;
    height: 100vh;
    z-index: 100;
    transition: left 0.3s ease;
  }

  .sidebar.open {
    left: 0;
  }

  .sidebar-toggle {
    display: flex;
  }

  .app-container {
    grid-template-columns: 1fr;
  }

  .main-content {
    grid-column: 1;
  }
}

@media (max-width: 768px) {
  .navbar {
    padding: 0 12px;
  }

  .brand-slogan {
    display: none;
  }

  .nav-btn .btn-text {
    display: none;
  }

  .nav-btn {
    padding: 8px 12px;
  }

  .messages-container {
    padding: 16px;
  }

  .input-area {
    padding: 12px 16px;
  }

  .emoji-picker {
    width: 280px;
    left: -10px;
  }

  .emoji-grid {
    grid-template-columns: repeat(7, 1fr);
  }
}
</style>