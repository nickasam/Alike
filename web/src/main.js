import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

// 导入全局响应式样式
import './styles/responsive.css'

// 创建应用实例
const app = createApp(App)

// 创建 Pinia 实例
const pinia = createPinia()

// 使用插件
app.use(pinia)
app.use(router)

// 挂载应用
app.mount('#app')

// ============================================================================
// 移动端优化
// ============================================================================

// 1. 防止双击缩放
let lastTouchEnd = 0;
document.addEventListener('touchend', function (event) {
  const now = Date.now();
  if (now - lastTouchEnd <= 300) {
    event.preventDefault();
  }
  lastTouchEnd = now;
}, false);

// 2. 动态设置视口高度（解决 100vh 问题）
function setAppHeight() {
  const vh = window.innerHeight * 0.01;
  document.documentElement.style.setProperty('--vh', `${vh}px`);
}

// 初始化
setAppHeight();

// 监听变化
window.addEventListener('resize', setAppHeight);
window.addEventListener('orientationchange', () => {
  setAppHeight();
  setTimeout(setAppHeight, 100); // 延迟再次设置，确保准确
});

// 3. 滚动锁定工具
let scrollPosition = 0;
let isLocked = false;

function lockScroll() {
  if (isLocked) return;

  scrollPosition = window.pageYOffset;
  document.body.style.overflow = 'hidden';
  document.body.style.position = 'fixed';
  document.body.style.top = `-${scrollPosition}px`;
  document.body.style.width = '100%';
  isLocked = true;
}

function unlockScroll() {
  if (!isLocked) return;

  document.body.style.removeProperty('overflow');
  document.body.style.removeProperty('position');
  document.body.style.removeProperty('top');
  document.body.style.removeProperty('width');
  window.scrollTo(0, scrollPosition);
  isLocked = false;
}

// 导出为全局方法
window.lockScroll = lockScroll;
window.unlockScroll = unlockScroll;
window.isScrollLocked = () => isLocked;

// 4. 检测移动端
const isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
window.isMobile = isMobile;

// 5. 检测触摸设备
const isTouch = 'ontouchstart' in window || navigator.maxTouchPoints > 0;
window.isTouch = isTouch;

// 6. 添加设备类型到 body
if (isMobile) {
  document.body.classList.add('is-mobile');
}
if (isTouch) {
  document.body.classList.add('is-touch');
}

// 7. 监听视口大小变化
let resizeTimeout;
window.addEventListener('resize', () => {
  clearTimeout(resizeTimeout);
  resizeTimeout = setTimeout(() => {
    // 触发自定义事件
    window.dispatchEvent(new CustomEvent('viewportChange'));
  }, 250);
});

// 8. 检测网络状态
if (navigator.connection) {
  const updateNetworkStatus = () => {
    const connection = navigator.connection;
    document.body.classList.toggle('slow-network', connection.effectiveType.includes('2g'));
  };

  navigator.connection.addEventListener('change', updateNetworkStatus);
  updateNetworkStatus();
}

console.log('📱 移动端优化已加载');