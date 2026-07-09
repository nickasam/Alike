/**
 * auth store — 当前用户与 JWT 状态。
 * 阶段二接入真实登录流：access token + refresh token 双令牌，
 * 状态在客户端初始化时从 localStorage 恢复。
 */
import { defineStore } from 'pinia'

export interface AuthUser {
  id: number
  /** 仅 /auth/me 返回本人邮箱 */
  email?: string
  nickname: string
  avatar_url?: string
  bio?: string
  industry?: string
  job_title?: string
  /** 工龄（年） */
  work_years?: number
  /** 牛马等级 */
  level: number
  /** 被共情数 */
  empathy_received: number
  /** 给出共情数 */
  empathy_given: number
  /** 累计打卡天数 */
  total_check_in_days: number
  /** 是否默认匿名 */
  is_anonymous: boolean
}

/** localStorage key —— 与 useApi/useAuth 保持一致 */
export const TOKEN_KEY = 'alike_token'
export const REFRESH_KEY = 'alike_refresh_token'
export const USER_KEY = 'alike_user'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: '' as string,
    refreshToken: '' as string,
    user: null as AuthUser | null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token && !!state.user,
  },

  actions: {
    /** 写入完整登录态（access + refresh + user），并持久化到 localStorage。 */
    setAuth(token: string, refreshToken: string, user: AuthUser) {
      this.token = token
      this.refreshToken = refreshToken
      this.user = user
      if (import.meta.client) {
        localStorage.setItem(TOKEN_KEY, token)
        localStorage.setItem(REFRESH_KEY, refreshToken)
        localStorage.setItem(USER_KEY, JSON.stringify(user))
      }
    },

    /** 仅更新 access token（refresh 刷新后调用）。 */
    setToken(token: string) {
      this.token = token
      if (import.meta.client) localStorage.setItem(TOKEN_KEY, token)
    },

    setUser(user: AuthUser) {
      this.user = user
      if (import.meta.client) localStorage.setItem(USER_KEY, JSON.stringify(user))
    },

    /** 客户端挂载时从 localStorage 恢复登录态。 */
    init() {
      if (!import.meta.client) return
      const token = localStorage.getItem(TOKEN_KEY)
      const refreshToken = localStorage.getItem(REFRESH_KEY)
      const rawUser = localStorage.getItem(USER_KEY)
      if (token) this.token = token
      if (refreshToken) this.refreshToken = refreshToken
      if (rawUser) {
        try {
          this.user = JSON.parse(rawUser) as AuthUser
        } catch {
          this.user = null
        }
      }
    },

    clear() {
      this.token = ''
      this.refreshToken = ''
      this.user = null
      if (import.meta.client) {
        localStorage.removeItem(TOKEN_KEY)
        localStorage.removeItem(REFRESH_KEY)
        localStorage.removeItem(USER_KEY)
      }
    },
  },
})
