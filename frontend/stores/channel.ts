/**
 * channel store — 频道列表与当前频道状态。骨架版，阶段三接入真实数据。
 */
import { defineStore } from 'pinia'

export type ChannelCategory = 'industry' | 'job' | 'topic' | 'custom'

export interface Channel {
  id: number
  name: string
  slug: string
  description?: string
  category: ChannelCategory
  icon?: string
  member_count: number
  /** 当前用户是否已加入 */
  joined?: boolean
  /** 未读数 */
  unread?: number
}

export const useChannelStore = defineStore('channel', {
  state: () => ({
    /** 侧边栏全量频道列表 */
    channels: [] as Channel[],
    /** 首页热门频道（独立状态，不与侧边栏全量互相覆盖） */
    hotChannels: [] as Channel[],
    currentId: null as number | null,
    loading: false,
  }),

  getters: {
    current: (state) =>
      state.channels.find((c) => c.id === state.currentId) ??
      state.hotChannels.find((c) => c.id === state.currentId) ??
      null,
    byCategory: (state) => (cat: ChannelCategory) =>
      state.channels.filter((c) => c.category === cat),
  },

  actions: {
    setChannels(list: Channel[]) {
      this.channels = list
    },
    setHotChannels(list: Channel[]) {
      this.hotChannels = list
    },
    setCurrent(id: number | null) {
      this.currentId = id
    },
  },
})
