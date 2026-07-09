/**
 * 客户端启动时恢复登录态：
 *   1. 从 localStorage 恢复 token / refreshToken / 缓存的 user，
 *      使刷新页面后 Pinia store 与路由守卫可立即判断登录状态；
 *   2. 若存在 token，异步调用 GET /auth/me 校验并刷新用户信息，
 *      失败（token 失效且无法刷新）时清除本地登录态。
 */
import { useAuthStore } from '~/stores/auth'

export default defineNuxtPlugin(async () => {
  const store = useAuthStore()
  store.init()

  if (!store.token) return

  try {
    const user = await useApi().get<import('~/stores/auth').AuthUser>('/auth/me')
    store.setUser(user)
  } catch {
    // token 失效且刷新失败：useApi 内部已尝试刷新，此处兜底清除
    store.clear()
  }
})
