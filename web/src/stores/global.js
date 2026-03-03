import { defineStore } from 'pinia'
import { getGlobalMessages, sendGlobalMessage, getOnlineUsers, getOnlineCount } from '@/api/global'

export const useGlobalChatStore = defineStore('globalChat', {
  state: () => ({
    messages: [],
    onlineUsers: [],
    onlineCount: 0,
    isLoading: false,
    error: null,
    refreshInterval: null
  }),

  getters: {
    sortedMessages: (state) => {
      return [...state.messages].sort((a, b) => 
        new Date(a.created_at) - new Date(b.created_at)
      )
    }
  },

  actions: {
    // 获取全局消息
    async fetchMessages(params = {}) {
      try {
        this.isLoading = true
        const response = await getGlobalMessages(params)
        this.messages = response.data || []
        return { success: true }
      } catch (error) {
        console.error('获取全局消息失败:', error)
        this.error = error.message
        return { success: false, message: error.message }
      } finally {
        this.isLoading = false
      }
    },

    // 发送全局消息
    async sendMessage(content) {
      try {
        const response = await sendGlobalMessage({ content })
        this.messages.push(response.data)
        return { success: true }
      } catch (error) {
        console.error('发送消息失败:', error)
        return { success: false, message: error.message }
      }
    },

    // 获取在线用户
    async fetchOnlineUsers() {
      try {
        const response = await getOnlineUsers()
        this.onlineUsers = response.data || []
        return { success: true }
      } catch (error) {
        console.error('获取在线用户失败:', error)
        return { success: false, message: error.message }
      }
    },

    // 获取在线人数
    async fetchOnlineCount() {
      try {
        const response = await getOnlineCount()
        this.onlineCount = response.data?.count || 0
        return { success: true }
      } catch (error) {
        console.error('获取在线人数失败:', error)
        return { success: false, message: error.message }
      }
    },

    // 添加新消息（用于实时更新）
    addMessage(message) {
      this.messages.push(message)
    },

    // 更新在线用户列表
    updateOnlineUsers(users) {
      this.onlineUsers = users
    },

    // 更新在线人数
    updateOnlineCount(count) {
      this.onlineCount = count
    },

    // 开始自动刷新
    startAutoRefresh(interval = 3000) {
      this.stopAutoRefresh()
      
      this.refreshInterval = setInterval(async () => {
        await this.fetchMessages()
        await this.fetchOnlineCount()
      }, interval)
    },

    // 停止自动刷新
    stopAutoRefresh() {
      if (this.refreshInterval) {
        clearInterval(this.refreshInterval)
        this.refreshInterval = null
      }
    },

    // 清空消息
    clearMessages() {
      this.messages = []
    }
  }
})