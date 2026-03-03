<template>
  <div class="container">
    <!-- 左侧图片区域（桌面端） -->
    <div class="left-panel">
      <img src="https://images.unsplash.com/photo-1523240795612-9a054b0db644?w=800&q=80" alt="Friends at cafe">
      <div class="left-overlay">
        <div class="left-text">
          <h1>Where Like Souls Meet</h1>
          <p>相似灵魂的相遇<br>发现志同道合的朋友</p>
        </div>
      </div>
    </div>

    <!-- 右侧表单区域 -->
    <div class="right-panel">
      <div class="logo">
        <div class="logo-icon">
          <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
          </svg>
        </div>
        <div class="logo-text">Alike</div>
      </div>

      <div class="welcome-text">欢迎回来，开始你的社交之旅</div>

      <div class="tabs">
        <button class="tab" :class="{ active: activeTab === 'login' }" @click="activeTab = 'login'">登录</button>
        <button class="tab" :class="{ active: activeTab === 'register' }" @click="activeTab = 'register'">注册</button>
      </div>

      <!-- 登录表单 -->
      <form v-if="activeTab === 'login'" @submit.prevent="handleLogin">
        <div class="form-group">
          <input type="tel" class="form-input" v-model="loginForm.phone" placeholder="手机号" required>
        </div>
        <div class="form-group">
          <input type="password" class="form-input" v-model="loginForm.password" placeholder="密码" required>
        </div>
        <button type="submit" class="submit-btn" :disabled="isLoading">
          {{ isLoading ? '登录中...' : '登录' }}
        </button>
        <div class="error-msg" v-if="errorMessage">{{ errorMessage }}</div>
      </form>

      <!-- 注册表单 -->
      <form v-else @submit.prevent="handleRegister">
        <div class="form-group">
          <input type="tel" class="form-input" v-model="registerForm.phone" placeholder="手机号" required>
        </div>
        <div class="form-group">
          <input type="text" class="form-input" v-model="registerForm.nickname" placeholder="昵称" required>
        </div>
        <div class="form-group">
          <input type="password" class="form-input" v-model="registerForm.password" placeholder="密码（至少8位）" required>
        </div>
        <button type="submit" class="submit-btn" :disabled="isLoading">
          {{ isLoading ? '注册中...' : '注册' }}
        </button>
        <div class="error-msg" v-if="errorMessage">{{ errorMessage }}</div>
      </form>

      <div class="divider"><span>其他登录方式</span></div>

      <div class="social-login">
        <button class="social-btn" title="微信" @click="showToast('微信登录开发中...')">
          <span class="material-icons">chat_bubble</span>
        </button>
        <button class="social-btn" title="QQ" @click="showToast('QQ登录开发中...')">
          <span class="material-icons">groups</span>
        </button>
        <button class="social-btn" title="微博" @click="showToast('微博登录开发中...')">
          <span class="material-icons">public</span>
        </button>
      </div>

      <div class="footer-links">
        <a @click="router.push('/launcher')">← 返回首页</a>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const activeTab = ref('login')
const isLoading = ref(false)
const errorMessage = ref('')

const loginForm = ref({
  phone: '',
  password: ''
})

const registerForm = ref({
  phone: '',
  nickname: '',
  password: ''
})

const handleLogin = async () => {
  try {
    isLoading.value = true
    errorMessage.value = ''
    
    const result = await userStore.login({
      phone: loginForm.value.phone,
      password: loginForm.value.password
    })
    
    if (result.success) {
      router.push('/')
    } else {
      errorMessage.value = result.message || '登录失败'
    }
  } catch (error) {
    errorMessage.value = '网络错误'
  } finally {
    isLoading.value = false
  }
}

const handleRegister = async () => {
  try {
    isLoading.value = true
    errorMessage.value = ''
    
    const result = await userStore.register({
      phone: registerForm.value.phone,
      password: registerForm.value.password,
      nickname: registerForm.value.nickname
    })
    
    if (result.success) {
      router.push('/')
    } else {
      errorMessage.value = result.message || '注册失败'
    }
  } catch (error) {
    errorMessage.value = '网络错误'
  } finally {
    isLoading.value = false
  }
}

const showToast = (message) => {
  alert(message)
}
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.container {
  width: 100%;
  height: 100vh;
  max-width: none;
  background: white;
  border-radius: 0;
  box-shadow: none;
  overflow: hidden;
  display: flex;
  min-height: 100vh;
}

/* 左侧图片区域（仅桌面端显示） */
.left-panel {
  flex: 1;
  position: relative;
  display: none;
  overflow: hidden;
}

.left-panel img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.left-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(180deg, rgba(102, 126, 234, 0.3) 0%, rgba(118, 75, 162, 0.6) 100%);
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding: 40px;
}

.left-text {
  color: white;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.left-text h1 {
  font-size: 42px;
  font-weight: 700;
  margin-bottom: 12px;
  line-height: 1.2;
}

.left-text p {
  font-size: 18px;
  opacity: 0.95;
  line-height: 1.6;
}

/* 右侧表单区域 */
.right-panel {
  flex: 1;
  padding: 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  max-width: 480px;
  margin: 0 auto;
}

.logo {
  text-align: center;
  margin-bottom: 24px;
}

.logo-icon {
  width: 56px;
  height: 56px;
  margin: 0 auto 12px;
  background: linear-gradient(135deg, #8b5cf6 0%, #ec4899 100%);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(139, 92, 246, 0.3);
  position: relative;
  overflow: hidden;
}

.logo-icon svg {
  width: 30px;
  height: 30px;
  fill: white;
}

.logo-icon::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: linear-gradient(
    45deg,
    transparent,
    rgba(255, 255, 255, 0.15),
    transparent
  );
  transform: rotate(45deg);
  animation: logoShine 3s infinite;
}

@keyframes logoShine {
  0% {
    transform: translateX(-100%) translateY(-100%) rotate(45deg);
  }
  100% {
    transform: translateX(100%) translateY(100%) rotate(45deg);
  }
}

.logo-text {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a1a;
}

.welcome-text {
  text-align: center;
  margin-bottom: 24px;
  color: #666;
  font-size: 14px;
}

.tabs {
  display: flex;
  margin-bottom: 24px;
  background: #f5f5f5;
  border-radius: 10px;
  padding: 4px;
}

.tab {
  flex: 1;
  padding: 10px;
  border: none;
  background: transparent;
  color: #999;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.tab.active {
  background: white;
  color: #8b5cf6;
  font-weight: 600;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
}

.form-group {
  margin-bottom: 16px;
}

.form-input {
  width: 100%;
  padding: 12px 16px;
  border: 1.5px solid #e5e5e5;
  border-radius: 8px;
  font-size: 14px;
  transition: all 0.3s ease;
  background: #fafafa;
}

.form-input:focus {
  outline: none;
  border-color: #8b5cf6;
  background: white;
  box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
}

.submit-btn {
  width: 100%;
  padding: 12px;
  background: linear-gradient(135deg, #8b5cf6 0%, #ec4899 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-top: 8px;
  box-shadow: 0 3px 10px rgba(139, 92, 246, 0.25);
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(139, 92, 246, 0.4);
}

.submit-btn:active:not(:disabled) {
  transform: translateY(0);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.divider {
  display: flex;
  align-items: center;
  margin: 28px 0;
  color: #ccc;
  font-size: 13px;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: #e5e5e5;
}

.divider span {
  padding: 0 16px;
}

.social-login {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  justify-content: center;
}

.social-btn {
  width: 50px;
  height: 50px;
  border: 2px solid #e5e5e5;
  border-radius: 12px;
  background: white;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: #666;
}

.social-btn .material-icons {
  font-size: 26px;
}

.social-btn:hover {
  border-color: #8b5cf6;
  background: #fafafa;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.2);
}

.footer-links {
  text-align: center;
  font-size: 14px;
  color: #999;
}

.footer-links a {
  color: #8b5cf6;
  text-decoration: none;
  font-weight: 500;
  cursor: pointer;
}

.error-msg {
  color: #ef4444;
  font-size: 13px;
  margin-top: 12px;
  text-align: center;
  min-height: 20px;
}

/* 平板设备 */
@media (min-width: 769px) {
  .left-panel {
    display: block;
  }
}

/* 移动设备 */
@media (max-width: 768px) {
  .container {
    max-width: 100%;
    border-radius: 0;
    box-shadow: none;
    min-height: 100vh;
    flex-direction: column;
  }

  .left-panel {
    display: none;
  }

  .right-panel {
    padding: 30px 20px;
    max-height: none;
  }

  .logo-icon {
    width: 60px;
    height: 60px;
    margin-bottom: 12px;
  }
  
  .logo-icon svg {
    width: 32px;
    height: 32px;
  }

  .logo-text {
    font-size: 26px;
  }

  .welcome-text {
    font-size: 14px;
    margin-bottom: 24px;
  }

  .tabs {
    margin-bottom: 24px;
    padding: 4px;
  }

  .tab {
    padding: 12px;
    font-size: 15px;
  }

  .form-input {
    padding: 14px 16px;
    border-width: 1px;
  }

  .submit-btn {
    padding: 14px;
  }

  .social-btn {
    width: 48px;
    height: 48px;
  }
}

/* 小屏幕手机 */
@media (max-width: 375px) {
  .right-panel {
    padding: 24px 16px;
  }

  .logo-text {
    font-size: 24px;
  }

  .tab {
    padding: 10px;
    font-size: 14px;
  }

  .form-input {
    padding: 12px 14px;
    font-size: 14px;
  }

  .social-login {
    gap: 10px;
  }

  .social-btn {
    width: 44px;
    height: 44px;
  }
}
</style>