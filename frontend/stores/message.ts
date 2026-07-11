/**
 * message store — 频道消息流与线程面板状态。
 *
 * 消息按 created_at 升序存放（顶部最旧、底部最新），符合聊天视图习惯。
 * 后端 GET /channels/:id/messages 按 DESC 游标返回 { list, has_more, next_cursor }，
 * 本地翻转为升序追加到列表头部（loadMore 加载更早历史，before=当前最旧消息 id）。
 *
 * WebSocket 事件：
 *   - new_message   → 追加到当前频道列表尾部（含自己发的回显）；
 *   - thread_reply  → 若命中已打开线程则追加回复，并同步父消息 reply_count；
 *   - empathy       → 更新对应消息的 empathy_count / empathized。
 */
import { defineStore } from 'pinia'

/** 消息作者公开信息（匿名/软删除消息不含）。 */
export interface Author {
  id: number
  nickname: string
  avatar_url: string
}

/** 与后端 message.Message 对齐的领域模型。 */
export interface Message {
  id: number
  channel_id: number
  parent_id?: number | null
  content: string
  /** 情绪标签 key（tired/angry/...），软删除后清空 */
  emotion?: string
  is_anonymous: boolean
  empathy_count: number
  reply_count: number
  is_deleted: boolean
  /** 匿名或软删除消息不返回作者 */
  author?: Author | null
  created_at: string
  deleted_at?: string | null
  /** 前端本地态：当前用户是否已共情 */
  empathized?: boolean
  /** 发送者客户端幂等标识，用于乐观回显与广播去重合并 */
  client_msg_id?: string
  /** 前端本地态：乐观插入、尚未收到服务端回环确认 */
  pending?: boolean
}

/** 游标分页响应体（对齐后端 listData）。 */
export interface MessagePage {
  list: Message[]
  has_more: boolean
  next_cursor: number
}

/** thread_reply 事件的 data 结构（对齐后端 BroadcastThreadReply）。 */
export interface ThreadReplyPayload {
  parent_id: number
  reply: Message
}

interface ChannelState {
  list: Message[]
  hasMore: boolean
  loading: boolean
  /** 是否已首屏加载过 */
  initialized: boolean
}

function emptyChannel(): ChannelState {
  return { list: [], hasMore: false, loading: false, initialized: false }
}

/** 按 id 升序插入，若已存在则更新（去重，防重连/回显重复）。 */
function upsertAscending(list: Message[], msg: Message) {
  const idx = list.findIndex((m) => m.id === msg.id)
  if (idx !== -1) {
    list[idx] = msg
    return
  }
  // 新消息通常最新，直接尾插；否则按 id 有序插入。
  if (list.length === 0 || msg.id > list[list.length - 1].id) {
    list.push(msg)
    return
  }
  const at = list.findIndex((m) => m.id > msg.id)
  list.splice(at === -1 ? list.length : at, 0, msg)
}

export const useMessageStore = defineStore('message', {
  state: () => ({
    byChannel: {} as Record<number, ChannelState>,
    // 线程面板状态
    threadOpen: false,
    threadParent: null as Message | null,
    threadReplies: [] as Message[],
    threadLoading: false,
    threadHasMore: false,
  }),

  getters: {
    listOf: (state) => (channelId: number) =>
      state.byChannel[channelId]?.list ?? [],
    channelState: (state) => (channelId: number) =>
      state.byChannel[channelId] ?? emptyChannel(),
  },

  actions: {
    ensureChannel(channelId: number): ChannelState {
      if (!this.byChannel[channelId]) {
        this.byChannel[channelId] = emptyChannel()
      }
      return this.byChannel[channelId]
    },

    /** 首屏加载最新一页消息（DESC → 翻转为 ASC）。 */
    async loadInitial(channelId: number) {
      const ch = this.ensureChannel(channelId)
      if (ch.loading) return
      ch.loading = true
      try {
        const page = await useApi().get<MessagePage>(
          `/channels/${channelId}/messages`,
          { limit: 20 },
        )
        // 后端 DESC（新→旧），翻转为升序展示。
        ch.list = [...page.list].reverse()
        ch.hasMore = page.has_more
        ch.initialized = true
      } finally {
        ch.loading = false
      }
    },

    /** 加载更早的历史消息（触顶时），before=当前最旧消息 id。 */
    async loadMore(channelId: number) {
      const ch = this.ensureChannel(channelId)
      if (ch.loading || !ch.hasMore || ch.list.length === 0) return
      ch.loading = true
      try {
        const before = ch.list[0].id
        const page = await useApi().get<MessagePage>(
          `/channels/${channelId}/messages`,
          { limit: 20, before },
        )
        const older = [...page.list].reverse()
        // 历史消息插到列表头部。
        ch.list = [...older, ...ch.list]
        ch.hasMore = page.has_more
      } finally {
        ch.loading = false
      }
    },

    /** 乐观插入：发送即上屏（临时负数 id + pending），WS 回环到达后被替换。 */
    addOptimistic(channelId: number, msg: Message) {
      const ch = this.ensureChannel(channelId)
      ch.list.push(msg)
    },

    /** WebSocket new_message：追加到对应频道。
     *  若命中本地乐观条目（同 client_msg_id 且 pending），则替换之，避免重复。 */
    receiveMessage(msg: Message) {
      if (!msg || !msg.channel_id) return
      const ch = this.ensureChannel(msg.channel_id)
      if (msg.client_msg_id) {
        const idx = ch.list.findIndex(
          (m) => m.pending && m.client_msg_id === msg.client_msg_id,
        )
        if (idx !== -1) {
          ch.list[idx] = msg // 用服务端真实消息替换乐观占位
          return
        }
      }
      upsertAscending(ch.list, msg)
    },

    /** 线程回复入列（幂等）：REST 响应与 WS thread_reply 广播共用此路径。
     *  仅当该回复 id 首次出现时才追加并对父消息 reply_count +1，
     *  避免 REST 响应与广播先后到达导致重复显示或计数翻倍。 */
    addThreadReply(parentId: number, reply: Message) {
      if (!reply) return
      const isNew =
        this.threadReplies.findIndex((m) => m.id === reply.id) === -1
      if (this.threadOpen && this.threadParent?.id === parentId) {
        upsertAscending(this.threadReplies, reply)
      }
      // 仅首次出现的回复才计数（无论线程面板是否打开）。
      if (isNew) this.bumpReplyCount(parentId)
    },

    /** WebSocket thread_reply：委托 addThreadReply 幂等入列 + 计数。 */
    receiveThreadReply(payload: ThreadReplyPayload) {
      if (!payload?.reply) return
      this.addThreadReply(payload.parent_id, payload.reply)
    },

    /** 同步某父消息的 reply_count（+1），并更新线程面板父引用。
     *  注意 threadParent 常与列表中的父消息是同一对象引用（openThread 直接传入列表项），
     *  故仅当二者为不同对象时才对 threadParent 补加，避免同一条消息被 +2。 */
    bumpReplyCount(parentId: number) {
      const ch = this.byChannel[this.threadParent?.channel_id ?? -1]
      const target =
        ch?.list.find((m) => m.id === parentId) ??
        Object.values(this.byChannel)
          .flatMap((c) => c.list)
          .find((m) => m.id === parentId)
      if (target) target.reply_count += 1
      if (this.threadParent?.id === parentId && this.threadParent !== target) {
        this.threadParent.reply_count += 1
      }
    },

    /** WebSocket empathy：更新共情数与本地已共情态。 */
    applyEmpathy(payload: { message_id: number; empathy_count: number; empathized?: boolean }) {
      const found =
        Object.values(this.byChannel)
          .flatMap((c) => c.list)
          .find((m) => m.id === payload.message_id) ??
        (this.threadParent?.id === payload.message_id ? this.threadParent : null) ??
        this.threadReplies.find((m) => m.id === payload.message_id)
      if (found) {
        found.empathy_count = payload.empathy_count
        if (payload.empathized !== undefined) found.empathized = payload.empathized
      }
    },

    /** WebSocket message_deleted：就地将消息置为已删除占位。 */
    markDeleted(payload: { message_id: number; channel_id?: number }) {
      if (!payload?.message_id) return
      const apply = (m: Message) => {
        m.is_deleted = true
        m.content = '该消息已删除'
        m.emotion = ''
        m.author = null
      }
      for (const c of Object.values(this.byChannel)) {
        const found = c.list.find((m) => m.id === payload.message_id)
        if (found) apply(found)
      }
      const reply = this.threadReplies.find((m) => m.id === payload.message_id)
      if (reply) apply(reply)
    },

    /** 打开线程面板并加载回复。 */
    async openThread(parent: Message) {
      this.threadOpen = true
      this.threadParent = parent
      this.threadReplies = []
      this.threadHasMore = false
      this.threadLoading = true
      try {
        const page = await useApi().get<MessagePage>(
          `/messages/${parent.id}/threads`,
          { limit: 20 },
        )
        // 线程回复后端已按 ASC 返回。
        this.threadReplies = page.list
        this.threadHasMore = page.has_more
      } finally {
        this.threadLoading = false
      }
    },

    /** 加载更多线程回复（after=当前最新回复 id）。 */
    async loadMoreThread() {
      if (
        !this.threadParent ||
        this.threadLoading ||
        !this.threadHasMore ||
        this.threadReplies.length === 0
      ) {
        return
      }
      this.threadLoading = true
      try {
        const after = this.threadReplies[this.threadReplies.length - 1].id
        const page = await useApi().get<MessagePage>(
          `/messages/${this.threadParent.id}/threads`,
          { limit: 20, after },
        )
        page.list.forEach((r) => upsertAscending(this.threadReplies, r))
        this.threadHasMore = page.has_more
      } finally {
        this.threadLoading = false
      }
    },

    closeThread() {
      this.threadOpen = false
      this.threadParent = null
      this.threadReplies = []
      this.threadHasMore = false
    },

    /**
     * REST 发送线程回复，返回落库后的回复对象。
     * 调用方应将返回值传入 addThreadReply 即时入列（幂等）；
     * 后端随后广播的 thread_reply 命中同 id 会被去重，不重复计数。
     */
    async sendReply(parentId: number, content: string, emotion?: string | null, isAnonymous = false) {
      return useApi().post<Message>(`/messages/${parentId}/replies`, {
        content,
        emotion: emotion ?? '',
        is_anonymous: isAnonymous,
      })
    },
  },
})
