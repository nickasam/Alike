<template>
  <div class="global-chat-container">
    <!-- 漂浮装饰 -->
    <div class="floating-circle circle-1"></div>
    <div class="floating-circle circle-2"></div>
    <div class="floating-circle circle-3"></div>
    <div class="floating-circle circle-4"></div>

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
      <!-- 导航栏 -->
      <div class="navbar">
        <div class="navbar-left">
          <button class="sidebar-toggle" @click="toggleSidebar" title="切换侧边栏">
            <span class="material-icons">menu</span>
          </button>
          <div class="navbar-avatar">{{ userInitial }}</div>
          <div class="navbar-info">
            <h2>Alike大家庭</h2>
            <div class="navbar-status">
              <div class="status-dot"></div>
              <span>{{ onlineCount }}</span> 人在线
            </div>
          </div>
        </div>
        <div class="navbar-actions">
          <button class="icon-btn" title="搜索" @click="showToast('搜索功能开发中...')">
            <span class="material-icons">search</span>
          </button>
          <button class="icon-btn" title="更多" @click="showToast('菜单功能开发中...')">
            <span class="material-icons">more_horiz</span>
          </button>
          <button class="icon-btn" title="退出" @click="logout">
            <span class="material-icons">logout</span>
          </button>
        </div>
      </div>

      <!-- 主应用区域 -->
      <div class="app-main">
        <!-- 侧边栏 -->
        <div class="sidebar" :class="{ collapsed: !sidebarOpen }">
          <div class="sidebar-header">
            <div class="sidebar-title">
              <span class="material-icons">people</span>
              在线用户
              <span class="sidebar-count">({{ onlineUsers.length }})</span>
            </div>
          </div>
          <div class="sidebar-users">
            <div 
              v-for="user in onlineUsers" 
              :key="user.id"
              class="user-item"
              @click="mentionUser(user)"
              :class="{ active: selectedUserId === user.id }"
            >
              <div class="user-avatar-small">
                {{ user.nickname ? user.nickname[0].toUpperCase() : '?' }}
                <div class="user-status-dot"></div>
              </div>
              <div class="user-info">
                <div class="user-name">{{ user.nickname }}</div>
                <div class="user-badge">{{ user.badge || '在线' }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 聊天区域 -->
        <div class="chat-area">
          <!-- 消息列表 -->
          <div class="messages-container" ref="messagesContainer">
            <!-- 空状态 -->
            <div class="empty-state" v-if="messages.length === 0">
              <div class="empty-icon-wrapper">
                <div class="empty-icon-glow"></div>
                <div class="empty-icon">
                  <span class="material-icons">chat_bubble</span>
                </div>
              </div>
              <div class="empty-text">暂无消息</div>
              <div class="empty-hint">成为第一个发言的人吧！</div>
            </div>

            <!-- 消息列表 -->
            <div 
              v-for="msg in sortedMessages" 
              :key="msg.id"
              class="message"
              :class="{ own: msg.user_id === currentUserId }"
            >
              <div class="message-avatar" v-if="msg.user_id !== currentUserId">
                {{ msg.username ? msg.username[0].toUpperCase() : '?' }}
              </div>
              <div class="message-content">
                <div class="message-header">
                  <span class="message-name" v-if="msg.user_id !== currentUserId">
                    {{ msg.username }}
                  </span>
                  <span class="message-time">{{ formatTime(msg.created_at) }}</span>
                </div>
                <div class="message-bubble">{{ msg.content }}</div>
              </div>
            </div>
          </div>

          <!-- 输入区域 -->
          <div class="input-container">
            <div class="input-wrapper">
              <div class="input-box">
                <textarea
                  class="message-input"
                  v-model="newMessage"
                  placeholder="输入消息... (@某人可快速提及)"
                  rows="1"
                  @keydown.exact.enter="sendMessage"
                  @input="handleInput"
                  ref="messageInput"
                ></textarea>
              </div>
              <div class="input-actions">
                <button 
                  class="action-btn" 
                  title="表情" 
                  @click="showEmojiPicker = !showEmojiPicker"
                >
                  <span class="material-icons">sentiment_satisfied</span>
                </button>
                <button class="send-btn" @click="sendMessage" :disabled="!newMessage.trim()">
                  <span class="material-icons">send</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useGlobalStore } from '@/stores/global'

const router = useRouter()
const userStore = useUserStore()
const globalStore = useGlobalStore()

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
const sidebarOpen = ref(true)
const onlineUsers = ref([])
const selectedUserId = ref(null)
const showEmojiPicker = ref(false)

const messagesContainer = ref(null)
const messageInput = ref(null)

let refreshInterval = null

// 计算属性
const userInitial = computed(() => {
  return currentUsername.value ? currentUsername.value[0].toUpperCase() : '?'
})

const sortedMessages = computed(() => {
  return [...messages.value].sort((a, b) => 
    new Date(a.created_at) - new Date(b.created_at)
  )
})

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
      // 尝试注册
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
    const result = await globalStore.getMessages()
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
      loadMessages()
    } else {
      showToast('发送失败：' + (result.message || '未知错误'))
    }
  } catch (error) {
    showToast('网络错误')
  }
}

const startAutoRefresh = () => {
  if (refreshInterval) return
  refreshInterval = setInterval(() => {
    loadMessages()
  }, 3000)
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
  return date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})
}

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
}

const mentionUser = (user) => {
  selectedUserId.value = user.id
  const mention = `@${user.nickname} `
  newMessage.value += mention
  if (messageInput.value) {
    messageInput.value.focus()
  }
}

const handleInput = (event) => {
  const textarea = event.target
  textarea.style.height = 'auto'
  textarea.style.height = Math.min(textarea.scrollHeight, 120) + 'px'
}

const showToast = (message) => {
  alert(message)
}

onMounted(() => {
  checkAuth()
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.global-chat-container {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
  background-size: 400% 400%;
  animation: gradientShift 15s ease infinite;
  min-height: 100vh;
  overflow: hidden;
  position: relative;
}

@keyframes gradientShift {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

/* 漂浮装饰圆圈 */
.floating-circle {
  position: fixed;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 20s infinite;
  pointer-events: none;
  z-index: 0;
}

.circle-1 {
  width: 300px;
  height: 300px;
  top: -100px;
  left: -100px;
  animation-delay: 0s;
}

.circle-2 {
  width: 200px;
  height: 200px;
  top: 20%;
  right: -50px;
  animation-delay: 5s;
}

.circle-3 {
  width: 150px;
  height: 150px;
  bottom: 30%;
  left: 10%;
  animation-delay: 10s;
}

.circle-4 {
  width: 250px;
  height: 250px;
  bottom: -80px;
  right: -80px;
  animation-delay: 15s;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0) rotate(0deg); }
  25% { transform: translate(30px, -30px) rotate(90deg); }
  50% { transform: translate(-20px, 20px) rotate(180deg); }
  75% { transform: translate(20px, 30px) rotate(270deg); }
}

.material-icons {
  font-family: 'Material Icons';
  font-weight: normal;
  font-style: normal;
  line-height: 1;
  letter-spacing: normal;
  text-transform: none;
  white-space: nowrap;
  word-wrap: normal;
  direction: ltr;
  -webkit-font-smoothing: antialiased;
  text-rendering: optimizeLegibility;
  -moz-osx-font-smoothing: grayscale;
}

/* 登录界面 */
.login-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  z-index: 9999;
  background: transparent;
}

.login-header {
  padding: 60px 20px 40px;
  text-align: center;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: relative;
  z-index: 1;
}

.login-logo-wrapper {
  position: relative;
  display: inline-block;
  margin-bottom: 24px;
}

.login-logo {
  width: 110px;
  height: 110px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  border-radius: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 16px 50px rgba(240, 147, 251, 0.5);
  position: relative;
  animation: bounce 2s ease-in-out infinite;
  z-index: 1;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.login-logo .material-icons {
  font-size: 56px;
  color: white;
}

.login-logo-glow {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 28px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  filter: blur(40px);
  opacity: 0.6;
  animation: glow 2s ease-in-out infinite;
  top: 0;
  left: 0;
}

@keyframes glow {
  0%, 100% { opacity: 0.6; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.1); }
}

.login-title {
  font-size: 42px;
  font-weight: 900;
  color: white;
  margin-bottom: 8px;
  text-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  letter-spacing: -1px;
}

.login-subtitle {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.95);
  font-weight: 500;
}

.login-form {
  padding: 0 20px 60px;
  max-width: 420px;
  margin: 0 auto;
  width: 100%;
  position: relative;
  z-index: 1;
}

.form-group {
  margin-bottom: 16px;
}

.form-input {
  width: 100%;
  padding: 18px 24px;
  border: none;
  border-radius: 16px;
  font-size: 16px;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  font-weight: 500;
}

.form-input:focus {
  outline: none;
  background: white;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.2);
  transform: translateY(-3px);
}

.form-input::placeholder {
  color: #9ca3af;
}

.login-btn {
  width: 100%;
  padding: 18px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
  border: none;
  border-radius: 16px;
  font-size: 17px;
  font-weight: 800;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-top: 24px;
  box-shadow: 0 12px 40px rgba(240, 147, 251, 0.4);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-3px) scale(1.02);
  box-shadow: 0 16px 50px rgba(240, 147, 251, 0.5);
}

.login-btn:active:not(:disabled) {
  transform: translateY(-1px) scale(0.98);
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-msg {
  color: #ffebee;
  font-size: 14px;
  margin-top: 16px;
  text-align: center;
  min-height: 24px;
  background: rgba(255, 255, 255, 0.15);
  padding: 12px;
  border-radius: 12px;
  backdrop-filter: blur(20px);
  font-weight: 600;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

/* 主界面 */
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  position: relative;
  z-index: 1;
}

/* 导航栏 */
.navbar {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px);
  padding: 16px 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 0;
  z-index: 100;
}

.navbar-left {
  display: flex;
  align-items: center;
  gap: 14px;
}

.navbar-avatar {
  width: 48px;
  height: 48px;
  border-radius: 16px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 900;
  font-size: 20px;
  box-shadow: 0 8px 24px rgba(240, 147, 251, 0.4);
}

.navbar-info h2 {
  font-size: 20px;
  font-weight: 800;
  color: #1a1a1a;
  letter-spacing: -0.5px;
}

.navbar-status {
  font-size: 14px;
  color: #22c55e;
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}

.status-dot {
  width: 8px;
  height: 8px;
  background: #22c55e;
  border-radius: 50%;
  display: inline-block;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.6; transform: scale(1.2); }
}

.navbar-actions {
  display: flex;
  gap: 8px;
}

.icon-btn {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: transparent;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #666;
}

.icon-btn .material-icons {
  font-size: 24px;
}

.icon-btn:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  transform: scale(1.05);
}

.sidebar-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: rgba(102, 126, 234, 0.1);
  border: none;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #667eea;
}

.sidebar-toggle:hover {
  background: rgba(102, 126, 234, 0.2);
  transform: scale(1.05);
}

/* 主应用容器 */
.app-main {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* 侧边栏 */
.sidebar {
  width: 280px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-right: 1px solid rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
  transition: transform 0.3s ease;
  position: relative;
  z-index: 50;
}

.sidebar.collapsed {
  transform: translateX(-100%);
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}

.sidebar-title {
  font-size: 16px;
  font-weight: 700;
  color: #1a1a1a;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.sidebar-count {
  font-size: 13px;
  color: #667eea;
  font-weight: 600;
}

.sidebar-users {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.user-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 4px;
}

.user-item:hover {
  background: rgba(102, 126, 234, 0.08);
  transform: translateX(4px);
}

.user-item.active {
  background: rgba(102, 126, 234, 0.12);
}

.user-avatar-small {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 800;
  font-size: 16px;
  flex-shrink: 0;
  position: relative;
}

.user-status-dot {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 12px;
  height: 12px;
  background: #22c55e;
  border: 2px solid white;
  border-radius: 50%;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 14px;
  font-weight: 600;
  color: #1a1a1a;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-badge {
  font-size: 11px;
  color: #9ca3af;
  margin-top: 2px;
}

/* 聊天区域 */
.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

/* 消息列表 */
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 24px 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  background: rgba(245, 247, 250, 0.5);
}

.message {
  display: flex;
  gap: 12px;
  animation: messageIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes messageIn {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.message-avatar {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 900;
  font-size: 20px;
  flex-shrink: 0;
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.3);
}

.message-content {
  flex: 1;
  max-width: 75%;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.message-name {
  font-size: 15px;
  font-weight: 800;
  color: #667eea;
}

.message-time {
  font-size: 13px;
  color: #9ca3af;
  font-weight: 600;
}

.message-bubble {
  padding: 14px 18px;
  border-radius: 20px;
  font-size: 15px;
  line-height: 1.6;
  word-wrap: break-word;
  background: white;
  color: #1a1a1a;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border-top-left-radius: 6px;
  font-weight: 500;
}

.message.own {
  flex-direction: row-reverse;
}

.message.own .message-content {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.message.own .message-bubble {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
  border-top-right-radius: 6px;
  border-top-left-radius: 20px;
  box-shadow: 0 6px 20px rgba(240, 147, 251, 0.4);
}

.message.own .message-header {
  flex-direction: row-reverse;
}

.message.own .message-name {
  color: #f5576c;
}

/* 输入区域 */
.input-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  padding: 16px 20px;
  border-top: 1px solid rgba(0, 0, 0, 0.08);
  padding-bottom: calc(16px + env(safe-area-inset-bottom));
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.08);
}

.input-wrapper {
  display: flex;
  gap: 12px;
  align-items: flex-end;
  max-width: 1400px;
  margin: 0 auto;
}

.input-box {
  flex: 1;
  position: relative;
}

.message-input {
  width: 100%;
  padding: 14px 18px;
  border: 2px solid rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  font-size: 15px;
  resize: none;
  transition: all 0.3s ease;
  background: rgba(249, 250, 251, 0.8);
  max-height: 120px;
  font-family: inherit;
  font-weight: 500;
}

.message-input:focus {
  outline: none;
  border-color: #f093fb;
  background: white;
  box-shadow: 0 0 0 4px rgba(240, 147, 251, 0.1);
}

.input-actions {
  display: flex;
  gap: 10px;
}

.action-btn {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  background: rgba(245, 245, 245, 0.8);
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 24px;
  color: #666;
  flex-shrink: 0;
}

.action-btn:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  transform: scale(1.05);
}

.send-btn {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  color: white;
  font-size: 22px;
  flex-shrink: 0;
  box-shadow: 0 6px 20px rgba(240, 147, 251, 0.4);
}

.send-btn:hover:not(:disabled) {
  transform: scale(1.08);
  box-shadow: 0 8px 28px rgba(240, 147, 251, 0.5);
}

.send-btn:active:not(:disabled) {
  transform: scale(0.95);
}

.send-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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

.empty-icon-wrapper {
  position: relative;
  display: inline-block;
  margin-bottom: 24px;
}

.empty-icon {
  width: 90px;
  height: 90px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 12px 40px rgba(240, 147, 251, 0.5);
  position: relative;
  z-index: 1;
  animation: floatIcon 3s ease-in-out infinite;
}

@keyframes floatIcon {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.empty-icon .material-icons {
  font-size: 46px;
  color: white;
}

.empty-icon-glow {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 24px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  filter: blur(30px);
  opacity: 0.5;
  animation: glow 3s ease-in-out infinite;
}

.empty-text {
  font-size: 20px;
  margin-bottom: 8px;
  color: #1a1a1a;
  font-weight: 700;
}

.empty-hint {
  font-size: 15px;
  color: #6b7280;
  font-weight: 500;
}

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, #f5576c 0%, #f093fb 100%);
}

/* 响应式 */
@media (max-width: 768px) {
  .message-content {
    max-width: 80%;
  }

  .message-bubble {
    padding: 12px 16px;
    font-size: 14px;
  }

  .login-title {
    font-size: 36px;
  }

  .login-logo {
    width: 100px;
    height: 100px;
  }

  .login-logo .material-icons {
    font-size: 50px;
  }

  .sidebar {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    transform: translateX(-100%);
  }

  .sidebar:not(.collapsed) {
    transform: translateX(0);
  }
}
</style>