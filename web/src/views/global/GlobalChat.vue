
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
      <!-- 主内容区 -->
      <main class="main-content">
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
const onlineCount = ref(0)
const onlineUsers = ref([])
const showEmojiPicker = ref(false)
const currentView = ref('chat')

const messagesContainer = ref(null)
const messageInput = ref(null)

let refreshInterval = null

const currentChannelName = computed(() => '全局聊天室')

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
    loadOnlineUsers()
    startAutoRefresh()
  }
}

const handleLogin = async () => {
  try {
    isLoading.value = true
    errorMessage.value = ''

    // 先尝试登录
    let result = await userStore.login({
      phone: loginForm.value.phone,
      password: loginForm.value.password
    })

    if (result.success) {
      currentUserId.value = userStore.userId
      currentUsername.value = userStore.nickname || ''
      isLoggedIn.value = true
      loadMessages()
      startAutoRefresh()
    } else {
      // 登录失败，尝试自动注册
      result = await userStore.register({
        phone: loginForm.value.phone,
        password: loginForm.value.password,
        nickname: '用户' + loginForm.value.phone.slice(-4)
      })

      if (result.success) {
        currentUserId.value = userStore.userId
        currentUsername.value = userStore.nickname || ''
        isLoggedIn.value = true
        loadMessages()
        startAutoRefresh()
      } else {
        errorMessage.value = result.message || '登录失败，请重试'
      }
    }
  } catch (error) {
    console.error('登录错误:', error)
    errorMessage.value = '网络错误，请检查连接'
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
    console.log('加载消息结果:', result)
    if (result.success && result.data) {
      messages.value = result.data
      console.log('当前消息列表:', messages.value)
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
    console.log('发送消息:', content)
    const result = await globalStore.sendMessage({ content })
    console.log('发送消息结果:', result)
    if (result.success) {
      newMessage.value = ''
      if (messageInput.value) {
        messageInput.value.style.height = 'auto'
      }
      // 立即加载消息列表
      await loadMessages()
    } else {
      showToast('发送失败：' + (result.message || '未知错误'))
    }
  } catch (error) {
    console.error('发送消息错误:', error)
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
  return num.toString()
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
      onlineUsers.value = result.data
    }
    // 单独获取在线人数
    const countResult = await globalStore.fetchOnlineCount()
    if (countResult.success && countResult.data !== undefined) {
      onlineCount.value = countResult.data
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

.global-chat-container {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  background: var(--bg-primary);
  color: var(--text-primary);

  /* 全屏布局：覆盖MainLayout的margin-left和侧边栏 */
  position: fixed;
  top: 64px;  /* 桌面端顶部导航栏高度 */
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 100;  /* 确保在侧边栏(z-index: 50)之上 */

  width: 100%;
  max-width: 100vw;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  margin: 0;
  padding: 0;
}

@media (max-width: 1023px) {
  .global-chat-container {
    top: 56px;  /* 平板/移动端顶部导航栏高度 */
    height: calc(100dvh - 56px);
  }
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
  padding: 14px 16px;  /* 增加内边距 */
  border: 1px solid var(--border-color);
  border-radius: 12px;
  font-size: 16px;  /* 修复：16px防止iOS自动缩放 */
  min-height: 44px;  /* 添加：触摸友好 */
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
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* ========== 顶部导航栏 ========== */
.navbar {
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


/* ========== 主内容区 ========== */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow: hidden;
  min-height: 0;
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
  background: rgba(255, 255, 255, 0.08);
  padding: 12px 16px;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
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
  background: linear-gradient(135deg, #8b5cf6 0%, #ec4899 100%);
  color: #ffffff;
  border: 2px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.4);
}

.message-own .message-name {
  color: var(--primary);
}

/* ========== 输入区域 ========== */
.input-area {
  padding: 20px 24px 24px;
  background: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
  position: relative;  /* 添加定位上下文 */
  z-index: 10;  /* 确保输入区域在最上层 */
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
  padding: 18px 60px 18px 20px;
  border: 3px solid rgba(255, 255, 255, 0.7);
  border-radius: 20px;
  font-size: 16px;
  font-weight: 500;
  resize: none;
  transition: all 0.3s ease;
  background: linear-gradient(135deg, #7c3aed 0%, #ec4899 100%);
  color: #ffffff;
  min-height: 60px;
  max-height: 150px;
  font-family: inherit;
  line-height: 1.5;
  box-shadow: 0 8px 24px rgba(139, 92, 246, 0.4), 0 0 0 1px rgba(255, 255, 255, 0.1);
}

.message-input:hover {
  border-color: rgba(255, 255, 255, 0.9);
  background: linear-gradient(135deg, #8b5cf6 0%, #f472b6 100%);
  box-shadow: 0 12px 32px rgba(139, 92, 246, 0.5), 0 0 0 1px rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.message-input:focus {
  outline: none;
  border-color: #ffffff;
  background: linear-gradient(135deg, #9b6cf6 0%, #f569b9 100%);
  box-shadow: 0 0 0 5px rgba(139, 92, 246, 0.5), 0 16px 40px rgba(236, 72, 153, 0.4);
  transform: translateY(-2px);
}

.message-input::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.input-actions {
  position: absolute;
  right: 10px;
  bottom: 10px;
  display: flex;
  gap: 6px;
}

.action-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: transparent;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 18px;
  color: var(--text-tertiary);
  transition: all 0.2s ease;
}

.action-btn:hover {
  background: rgba(139, 92, 246, 0.15);
  color: var(--primary);
  transform: scale(1.05);
}

.action-btn.active {
  background: rgba(139, 92, 246, 0.15);
  color: var(--primary);
}

.send-btn {
  width: 60px;
  height: 60px;
  border-radius: 16px;
  background: linear-gradient(135deg, #7c3aed 0%, #ec4899 100%);
  border: 3px solid rgba(255, 255, 255, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: white;
  font-size: 24px;
  box-shadow: 0 8px 24px rgba(139, 92, 246, 0.4);
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.send-btn:hover:not(:disabled) {
  transform: scale(1.1) translateY(-2px);
  box-shadow: 0 12px 32px rgba(139, 92, 246, 0.6);
  border-color: rgba(255, 255, 255, 0.8);
}

.send-btn:active:not(:disabled) {
  transform: scale(0.95);
}

.send-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  transform: none;
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
  .app-container {
    grid-template-columns: 1fr;
  }

  .main-content {
    grid-column: 1;
  }
}

@media (max-width: 767px) {  /* 修复：统一使用767px断点 */
  /* 容器已经在桌面端设置为fixed，移动端只需调整top */
  .global-chat-container {
    top: 56px;  /* 移动端顶部导航栏高度 */
    height: calc(100dvh - 56px);
  }

  .app-container {
    width: 100%;
    max-width: 100%;
    box-sizing: border-box;
  }

  .main-content {
    width: 100%;
    max-width: 100%;
    box-sizing: border-box;
  }

  .navbar {
    padding: 0 12px;
    height: 56px;  /* 移动端优化 */
  }

  .brand-slogan {
    display: none;
  }

  .nav-btn .btn-text {
    display: none;
  }

  .nav-btn {
    padding: 8px 12px;
    min-height: 44px;  /* 触摸友好 */
    min-width: 44px;
  }

  .chat-header {
    padding: 16px;
  }

  .chat-header h1 {
    font-size: 20px;  /* 响应式字体 */
  }

  .messages-container {
    padding: 12px;  /* 移动端减少padding */
    flex: 1;  /* 占据剩余空间 */
    overflow-y: auto;  /* 允许滚动 */
  }

  .message-group {
    margin-bottom: 12px;
  }

  .message-avatar {
    width: 36px;
    height: 36px;
    font-size: 14px;
  }

  .message-content {
    font-size: 14px;
    line-height: 1.4;
  }

  .input-area {
    padding: 12px 12px;  /* 移动端padding */
    padding-bottom: max(12px, env(safe-area-inset-bottom));  /* 底部安全区域 */
    box-sizing: border-box;
    position: relative;  /* 确保定位上下文 */
    z-index: 10;  /* 确保在最上层 */
    flex-shrink: 0;  /* 防止被压缩 */
  }

  .message-input {
    padding: 12px 50px 12px 16px;  /* 恢复合适的padding */
    font-size: 16px;  /* 防止缩放 */
    min-height: 44px;  /* 恢复触摸友好的高度 */
    height: auto;  /* 自适应高度 */
    line-height: 1.4;  /* 行高 */
    position: relative;  /* 确保输入框可点击 */
    z-index: 11;  /* 比input-area更高 */
  }

  .send-btn {
    width: 44px;  /* 触摸友好 */
    height: 44px;
    min-width: 44px;
    min-height: 44px;
  }

  .emoji-btn {
    width: 44px;
    height: 44px;
    min-width: 44px;
    min-height: 44px;
  }

  .emoji-picker {
    width: 280px;
    left: -10px;
    bottom: 60px;  /* 避免被输入框遮挡 */
  }

  .emoji-grid {
    grid-template-columns: repeat(7, 1fr);
  }

  /* 登录界面移动端优化 */
  .login-container {
    padding: 20px;
  }

  .login-title {
    font-size: 32px;  /* 响应式 */
  }

  .form-group {
    margin-bottom: 12px;
  }
}
</style>