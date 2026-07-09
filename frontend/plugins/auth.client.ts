/**
 * 客户端启动时从 localStorage 恢复登录态，
 * 使刷新页面后 Pinia store 与路由守卫可正确判断登录状态。
 */
import { useAuthStore } from '~/stores/auth'

export default defineNuxtPlugin(() => {
  useAuthStore().init()
})
