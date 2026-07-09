/**
 * auth store — 当前用户与 JWT 状态。骨架版，阶段二接入真实登录流。
 */
import { defineStore } from 'pinia'

export interface AuthUser {
  id: number
  nickname: string
  avatar?: string
  industry?: string
  job?: string
  /** 牛马等级 */
  level: number
  /** 累计共情数 */
  empathy_count: number
  /** 是否默认匿名 */
  anonymous: boolean
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: '' as string,
    user: null as AuthUser | null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token && !!state.user,
  },

  actions: {
    setAuth(token: string, user: AuthUser) {
      this.token = token
      this.user = user
      if (import.meta.client) localStorage.setItem('alike_token', token)
    },
    setUser(user: AuthUser) {
      this.user = user
    },
    clear() {
      this.token = ''
      this.user = null
      if (import.meta.client) localStorage.removeItem('alike_token')
    },
  },
})
