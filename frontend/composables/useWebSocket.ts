/**
 * useWebSocket — 全局唯一 WebSocket 客户端。
 *
 * 端点：GET /api/ws（升级为 WS，JWT 不放 URL query，避免泄漏到访问日志）。
 *
 * 协议信封（与后端 ws.Envelope 对齐）：
 *   { type: string, data?: object, channel_id?: number }
 *
 * 连接流程：
 *   1. onopen 后立即发送首帧 { type: "auth", data: { token } }；
 *   2. 收到 { type: "auth_ok" } 后置 authed=true，并重放此前请求的频道订阅；
 *   3. 收到 { type: "ping" } 回 { type: "pong" }（叠加协议层 ping 控制帧，双保险）；
 *   4. 断线后指数退避重连（1s→2s→4s…上限 30s），重连后重新 auth 并重放订阅。
 *
 * 客户端 → 服务端：auth / join_channel / leave_channel / typing / send_message / pong
 * 服务端 → 客户端：auth_ok / new_message / thread_reply / empathy / user_joined /
 *                   emotion_update / notification / error / ping
 *
 * 该 composable 返回进程内单例（模块级状态），任意组件调用共享同一连接。
 */
import { useAuthStore } from '~/stores/auth'

/** 服务端 → 客户端可订阅的事件类型。 */
export type WsServerEvent =
  | 'auth_ok'
  | 'new_message'
  | 'thread_reply'
  | 'message_deleted'
  | 'empathy'
  | 'user_joined'
  | 'emotion_update'
  | 'notification'
  | 'typing'
  | 'error'

/** 客户端 → 服务端事件类型。 */
export type WsClientEvent =
  | 'auth'
  | 'join_channel'
  | 'leave_channel'
  | 'typing'
  | 'send_message'
  | 'pong'

/** WebSocket 帧信封。 */
export interface WsEnvelope<T = unknown> {
  type: string
  data?: T
  channel_id?: number
}

type Handler<T = unknown> = (data: T, channelId?: number) => void

const RECONNECT_BASE = 1000
const RECONNECT_MAX = 30000

// —— 模块级单例状态（跨组件共享一条连接）——
let socket: WebSocket | null = null
const connected = ref(false)
const authed = ref(false)
const handlers = new Map<string, Set<Handler>>()
// 已请求订阅的频道集合，重连后重放。
const joinedChannels = new Set<number>()
// 出站待发队列：连接/鉴权未就绪时暂存帧，auth_ok 后按序补发，避免静默丢弃。
const outboundQueue: string[] = []
let reconnectAttempts = 0
let reconnectTimer: ReturnType<typeof setTimeout> | null = null
// 等待 token 恢复的重试（auth 插件异步 init 时，connect 可能早于 token 就绪）。
let tokenRetryAttempts = 0
let tokenRetryTimer: ReturnType<typeof setTimeout> | null = null
// 主动断开标记：为 true 时 onclose 不触发重连。
let manualClose = false

/** 读取当前 access token（登录后写入）。 */
function currentToken(): string {
  return useAuthStore().token
}

function resolveUrl(): string {
  const base = useRuntimeConfig().public.wsBase as string
  if (base.startsWith('ws')) return base
  if (import.meta.client) {
    const proto = location.protocol === 'https:' ? 'wss' : 'ws'
    return `${proto}://${location.host}${base}`
  }
  return base
}

/** 发送一帧。auth 帧或连接已 OPEN 时直接发；否则入队待 auth_ok 后补发，
 *  并确保连接已发起，避免"连接未就绪 → 消息静默丢弃"。 */
function sendRaw(type: WsClientEvent, data?: unknown, channelId?: number) {
  const env: WsEnvelope = { type }
  if (data !== undefined) env.data = data
  if (channelId !== undefined) env.channel_id = channelId
  const frame = JSON.stringify(env)

  // auth 帧在 onopen 时发送，直接走；其余帧需连接 OPEN。
  if (socket?.readyState === WebSocket.OPEN) {
    socket.send(frame)
    return
  }
  // 未就绪：入队（auth 帧不入队，onopen 会重发），并确保连接已发起。
  if (type !== 'auth') {
    outboundQueue.push(frame)
    if (!socket || socket.readyState === WebSocket.CLOSED) open()
  }
}

/** auth_ok 后按序补发排队的帧。 */
function flushQueue() {
  if (socket?.readyState !== WebSocket.OPEN) return
  while (outboundQueue.length > 0) {
    socket.send(outboundQueue.shift() as string)
  }
}

function dispatch(env: WsEnvelope) {
  handlers.get(env.type)?.forEach((fn) => fn(env.data, env.channel_id))
}

function scheduleReconnect() {
  if (manualClose || reconnectTimer) return
  const delay = Math.min(RECONNECT_BASE * 2 ** reconnectAttempts, RECONNECT_MAX)
  reconnectAttempts += 1
  reconnectTimer = setTimeout(() => {
    reconnectTimer = null
    open()
  }, delay)
}

/** token 尚未就绪时的重试：每 300ms 一次，最多 ~20 次(6s)。
 *  覆盖 auth 插件异步恢复登录态的窗口；超时仍无 token 视为未登录，停止。 */
function scheduleTokenRetry() {
  if (tokenRetryTimer || tokenRetryAttempts >= 20) return
  tokenRetryAttempts += 1
  tokenRetryTimer = setTimeout(() => {
    tokenRetryTimer = null
    open()
  }, 300)
}

function open() {
  if (import.meta.server) return
  if (socket && socket.readyState !== WebSocket.CLOSED) return
  const token = currentToken()
  if (!token) {
    // token 可能尚未从 localStorage 恢复（auth 插件异步 init）。
    // 不要静默放弃——短延迟后重试，否则 WS 永不连接、收不到任何实时消息。
    if (!manualClose) scheduleTokenRetry()
    return
  }

  manualClose = false
  const ws = new WebSocket(resolveUrl())
  socket = ws

  ws.onopen = () => {
    connected.value = true
    tokenRetryAttempts = 0
    // 首帧鉴权：JWT 走消息体，不放 URL。
    sendRaw('auth', { token: currentToken() })
  }

  ws.onmessage = (ev: MessageEvent) => {
    let env: WsEnvelope
    try {
      env = JSON.parse(ev.data) as WsEnvelope
    } catch {
      return // 忽略非 JSON
    }
    switch (env.type) {
      case 'auth_ok':
        authed.value = true
        reconnectAttempts = 0
        // 重连后重放订阅，再补发排队的出站帧（如连接未就绪时发的消息）。
        joinedChannels.forEach((id) => sendRaw('join_channel', { channel_id: id }))
        flushQueue()
        break
      case 'ping':
        sendRaw('pong')
        break
    }
    dispatch(env)
  }

  ws.onclose = () => {
    connected.value = false
    authed.value = false
    socket = null
    if (!manualClose) scheduleReconnect()
  }

  ws.onerror = () => {
    // 关闭交由 onclose 统一处理（触发重连）。
    ws.close()
  }
}

function close() {
  manualClose = true
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  if (tokenRetryTimer) {
    clearTimeout(tokenRetryTimer)
    tokenRetryTimer = null
  }
  tokenRetryAttempts = 0
  reconnectAttempts = 0
  joinedChannels.clear()
  outboundQueue.length = 0
  socket?.close()
  socket = null
  connected.value = false
  authed.value = false
}

export function useWebSocket() {
  /** 建立连接（幂等）。需在登录且 token 就绪后调用。 */
  function connect() {
    manualClose = false
    tokenRetryAttempts = 0
    open()
  }

  /** 主动断开并停止重连。 */
  function disconnect() {
    close()
  }

  /** 订阅服务端事件，返回取消订阅函数。 */
  function on<T = unknown>(event: WsServerEvent, handler: Handler<T>): () => void {
    if (!handlers.has(event)) handlers.set(event, new Set())
    handlers.get(event)!.add(handler as Handler)
    return () => handlers.get(event)?.delete(handler as Handler)
  }

  /** 加入频道（记录以便重连重放）。若已鉴权立即发送。 */
  function joinChannel(channelId: number) {
    joinedChannels.add(channelId)
    if (authed.value) sendRaw('join_channel', { channel_id: channelId })
  }

  /** 离开频道。 */
  function leaveChannel(channelId: number) {
    joinedChannels.delete(channelId)
    sendRaw('leave_channel', { channel_id: channelId })
  }

  /** 广播「正在输入」。 */
  function sendTyping(channelId: number) {
    sendRaw('typing', { channel_id: channelId })
  }

  /**
   * 经 WebSocket 发送消息。client_msg_id 用于幂等去重（后端按连接去重）。
   */
  function sendMessage(
    channelId: number,
    content: string,
    emotion?: string | null,
    clientMsgId?: string,
    isAnonymous = false,
  ) {
    sendRaw('send_message', {
      channel_id: channelId,
      content,
      emotion: emotion ?? '',
      is_anonymous: isAnonymous,
      client_msg_id: clientMsgId ?? '',
    })
  }

  return {
    connected,
    authed,
    connect,
    disconnect,
    on,
    joinChannel,
    leaveChannel,
    sendTyping,
    sendMessage,
  }
}
