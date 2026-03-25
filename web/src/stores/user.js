import { defineStore } from 'pinia'
import { login, register, logout as apiLogout } from '@/api/auth'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('alike_access_token') || '',
    refreshToken: localStorage.getItem('alike_refresh_token') || '',
    userInfo: JSON.parse(localStorage.getItem('alike_user_info') || 'null'),
    isLoggedIn: !!localStorage.getItem('alike_access_token')
  }),

  getters: {
    userId: (state) => state.userInfo?.id,
    nickname: (state) => state.userInfo?.nickname,
    avatar: (state) => state.userInfo?.avatar,
    bio: (state) => state.userInfo?.bio
  },

  actions: {
    // 登录
    async login(credentials) {
      try {
        const response = await login(credentials)
        const { user, tokens } = response.data

        this.token = tokens.access_token
        this.refreshToken = tokens.refresh_token
        this.userInfo = user
        this.isLoggedIn = true

        // 保存到 localStorage
        localStorage.setItem('alike_access_token', tokens.access_token)
        localStorage.setItem('alike_refresh_token', tokens.refresh_token)
        localStorage.setItem('alike_user_id', user.id)
        localStorage.setItem('alike_username', user.nickname) // 保存用户名
        localStorage.setItem('alike_user_info', JSON.stringify(user))

        return { success: true }
      } catch (error) {
        console.error('登录失败:', error)
        return {
          success: false,
          message: error.response?.data?.error?.message || error.response?.data?.message || '登录失败'
        }
      }
    },

    // 注册
    async register(userData) {
      try {
        const response = await register(userData)
        const { user, tokens } = response.data

        this.token = tokens.access_token
        this.refreshToken = tokens.refresh_token
        this.userInfo = user
        this.isLoggedIn = true

        // 保存到 localStorage
        localStorage.setItem('alike_access_token', tokens.access_token)
        localStorage.setItem('alike_refresh_token', tokens.refresh_token)
        localStorage.setItem('alike_user_id', user.id)
        localStorage.setItem('alike_username', user.nickname) // 保存用户名
        localStorage.setItem('alike_user_info', JSON.stringify(user))

        return { success: true }
      } catch (error) {
        console.error('注册失败:', error)
        return {
          success: false,
          message: error.response?.data?.error?.message || error.response?.data?.message || '注册失败'
        }
      }
    },

    // 登出
    async logout() {
      try {
        await apiLogout()
      } catch (error) {
        console.error('登出请求失败:', error)
      } finally {
        // 清除状态
        this.token = ''
        this.refreshToken = ''
        this.userInfo = null
        this.isLoggedIn = false
        
        // 清除 localStorage
        localStorage.removeItem('alike_access_token')
        localStorage.removeItem('alike_refresh_token')
        localStorage.removeItem('alike_user_id')
        localStorage.removeItem('alike_user_info')
      }
    },

    // 更新用户信息
    updateUserInfo(userInfo) {
      this.userInfo = { ...this.userInfo, ...userInfo }
      localStorage.setItem('alike_user_info', JSON.stringify(this.userInfo))
    }
  }
})