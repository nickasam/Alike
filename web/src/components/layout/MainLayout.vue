<template>
  <div class="main-layout">
    <!-- 背景装饰层 -->
    <div class="background-decoration">
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
      <div class="gradient-orb orb-3"></div>
      <div class="grid-pattern"></div>
    </div>

    <!-- 顶部导航栏 -->
    <nav class="navbar">
      <div class="navbar-brand">
        <div class="brand-logo">💜</div>
        <div class="brand-text">
          <div class="brand-name">Alike</div>
          <div class="brand-slogan">相似灵魂的相遇</div>
        </div>
      </div>

      <div class="nav-actions">
        <router-link 
          to="/" 
          class="nav-btn"
          :class="{ active: $route.path === '/' }"
        >
          <span class="btn-icon">📍</span>
          <span class="btn-text">附近用户</span>
        </router-link>
        <router-link 
          to="/global" 
          class="nav-btn"
          :class="{ active: $route.path === '/global' }"
        >
          <span class="btn-icon">💬</span>
          <span class="btn-text">全局聊天</span>
        </router-link>
        <router-link 
          to="/profile" 
          class="nav-btn"
          :class="{ active: $route.path === '/profile' }"
        >
          <span class="btn-icon">👤</span>
          <span class="btn-text">我的</span>
        </router-link>
        <button @click="handleLogout" class="logout-btn" title="退出登录">
          <span class="logout-icon">🚪</span>
          <span class="logout-text">退出</span>
        </button>
      </div>
    </nav>

    <!-- 主内容区 -->
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const handleLogout = async () => {
  if (confirm('确定要退出登录吗？')) {
    await userStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.main-layout {
  /* CSS 变量 */
  --primary: #8b5cf6;
  --primary-light: #a78bfa;
  --primary-dark: #7c3aed;
  --secondary: #ec4899;
  --accent: #06b6d4;

  /* 渐变 */
  --gradient-brand: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  --gradient-purple: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

  /* 玻璃态 */
  --glass-bg: rgba(255, 255, 255, 0.15);
  --glass-border: rgba(255, 255, 255, 0.2);

  /* 深色主题 */
  --bg-dark: #0a0a0f;
  --bg-darker: #050508;
  --text-primary: #ffffff;
  --text-secondary: rgba(255, 255, 255, 0.7);
  --text-tertiary: rgba(255, 255, 255, 0.5);

  /* 阴影 */
  --shadow-sm: 0 4px 15px rgba(0, 0, 0, 0.2);
  --shadow-md: 0 8px 25px rgba(0, 0, 0, 0.3);
  --shadow-lg: 0 12px 40px rgba(240, 147, 251, 0.4);

  /* 布局 */
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: var(--bg-dark);
  color: var(--text-primary);
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  overflow-x: hidden;
  position: relative;
}

/* ========== 背景装饰层 ========== */
.background-decoration {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
  pointer-events: none;
  z-index: 0;
}

/* 网格图案 */
.grid-pattern {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: 
    linear-gradient(rgba(139, 92, 246, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(139, 92, 246, 0.03) 1px, transparent 1px);
  background-size: 50px 50px;
  opacity: 0.5;
}

/* 渐变光球 */
.gradient-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.3;
  animation: orbFloat 20s ease-in-out infinite;
}

.orb-1 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, rgba(139, 92, 246, 0.4) 0%, transparent 70%);
  top: -200px;
  right: -150px;
  animation-delay: 0s;
}

.orb-2 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(236, 72, 153, 0.3) 0%, transparent 70%);
  bottom: -150px;
  left: -100px;
  animation-delay: 7s;
}

.orb-3 {
  width: 350px;
  height: 350px;
  background: radial-gradient(circle, rgba(6, 182, 212, 0.25) 0%, transparent 70%);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: 14s;
}

@keyframes orbFloat {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  25% {
    transform: translate(30px, -30px) scale(1.1);
  }
  50% {
    transform: translate(-20px, 20px) scale(0.95);
  }
  75% {
    transform: translate(20px, 30px) scale(1.05);
  }
}

/* ========== 顶部导航栏 ========== */
.navbar {
  background: var(--glass-bg);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--glass-border);
  padding: 0 40px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 70px;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: var(--shadow-sm);
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  text-decoration: none;
}

.brand-logo {
  width: 44px;
  height: 44px;
  background: var(--gradient-brand);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  box-shadow: var(--shadow-lg);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

.brand-text {
  display: flex;
  flex-direction: column;
}

.brand-name {
  font-size: 24px;
  font-weight: 900;
  background: linear-gradient(135deg, #ffffff 0%, #f0f4ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.5px;
}

.brand-slogan {
  font-size: 12px;
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
  padding: 12px 20px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-size: 15px;
  font-weight: 700;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
}

.nav-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
  transform: translateY(-1px);
}

.nav-btn.active {
  background: rgba(255, 255, 255, 0.95);
  color: #667eea;
  box-shadow: var(--shadow-sm);
}

.btn-icon {
  font-size: 18px;
}

.logout-btn {
  padding: 12px 28px;
  background: rgba(255, 255, 255, 0.95);
  color: #667eea;
  border: none;
  font-size: 15px;
  font-weight: 700;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 6px;
  box-shadow: var(--shadow-sm);
}

.logout-btn:hover {
  transform: translateY(-2px) scale(1.05);
  box-shadow: var(--shadow-md);
  color: #f5576c;
}

.logout-icon {
  font-size: 18px;
}

/* ========== 主内容区 ========== */
.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 0;
  max-width: 100%;
  position: relative;
  z-index: 1;
}

/* ========== 响应式 ========== */
@media (max-width: 1024px) {
  .navbar {
    padding: 0 24px;
  }
  
  .brand-slogan {
    display: none;
  }
  
  .brand-name {
    font-size: 20px;
  }
  
  .orb-1, .orb-2, .orb-3 {
    width: 300px;
    height: 300px;
  }
}

@media (max-width: 768px) {
  .navbar {
    padding: 0 16px;
    height: 60px;
  }

  .brand-logo {
    width: 38px;
    height: 38px;
    font-size: 20px;
  }

  .brand-name {
    font-size: 18px;
  }

  .nav-btn .btn-text,
  .logout-btn .logout-text {
    display: none;
  }

  .nav-btn {
    padding: 10px 14px;
  }

  .logout-btn {
    padding: 10px 16px;
  }
  
  .orb-1, .orb-2, .orb-3 {
    width: 200px;
    height: 200px;
  }
}

@media (max-width: 480px) {
  .navbar {
    padding: 0 12px;
  }
  
  .nav-actions {
    gap: 4px;
  }
  
  .nav-btn, .logout-btn {
    padding: 8px 12px;
  }
  
  .grid-pattern {
    background-size: 30px 30px;
  }
}
</style>