/**
 * auth 路由中间件 —— 保护需登录的页面。
 *
 * 用法：在页面 <script setup> 中
 *   definePageMeta({ middleware: 'auth' })
 *
 * 规则：
 *   - 未登录访问受保护页 → 跳 /login（带 redirect 参数）
 *   - 已登录访问 /login  → 跳首页
 *
 * 登录态在客户端由 stores/auth 的 init() 从 localStorage 恢复
 * （见 plugins/auth.client.ts）。SSR 阶段无 localStorage，
 * 认证判断延迟到客户端，避免 SSR 误判导致的闪跳。
 */
export default defineNuxtRouteMiddleware((to) => {
  const store = useAuthStore()

  // SSR 阶段无法读取 localStorage，交由客户端再判断
  if (import.meta.server) return

  const authed = store.isAuthenticated

  if (to.path === '/login') {
    if (authed) return navigateTo('/')
    return
  }

  if (!authed) {
    return navigateTo({ path: '/login', query: { redirect: to.fullPath } })
  }
})
