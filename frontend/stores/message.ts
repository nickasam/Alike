/**
 * message store — 频道消息流状态。骨架版，阶段四接入消息与线程。
 */
import { defineStore } from 'pinia'

export interface Message {
  id: number
  channel_id: number
  /** 匿名消息不返回 user_id */
  user_id?: number
  nickname?: string
  avatar?: string
  parent_id?: number | null
  content: string
  /** 情绪标签 key，如 tired / angry / cheer */
  emotion?: string | null
  anonymous: boolean
  empathy_count: number
  empathized?: boolean
  created_at: string
}

export const useMessageStore = defineStore('message', {
  state: () => ({
    // 以 channelId 为 key 缓存消息列表
    byChannel: {} as Record<number, Message[]>,
    loading: false,
  }),

  getters: {
    listOf: (state) => (channelId: number) =>
      state.byChannel[channelId] ?? [],
  },

  actions: {
    setMessages(channelId: number, list: Message[]) {
      this.byChannel[channelId] = list
    },
    prepend(channelId: number, msg: Message) {
      const arr = this.byChannel[channelId] ?? (this.byChannel[channelId] = [])
      arr.push(msg)
    },
  },
})
