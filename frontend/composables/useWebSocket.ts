/**
 * useWebSocket — WebSocket 连接封装（骨架）。
 *
 * 端点：WS /ws?token=<JWT>
 * 客户端 → 服务端事件：join_channel / leave_channel / typing / send_message
 * 服务端 → 客户端事件：new_message / thread_reply / empathy / user_joined /
 *                       emotion_update / notification
 *
 * 说明：这是阶段一骨架，提供连接/断开/发送/订阅的最小可用 API，
 * 具体业务事件处理在各功能阶段接入。
 */

export type WsEvent =
  | 'new_message'
  | 'thread_reply'
  | 'empathy'
  | 'user_joined'
  | 'emotion_update'
  | 'notification'

export type WsClientEvent =
  | 'join_channel'
  | 'leave_channel'
  | 'typing'
  | 'send_message'

export interface WsMessage<T = unknown> {
  event: WsEvent | WsClientEvent
  data: T
}

type Handler<T = unknown> = (data: T) => void

export function useWebSocket() {
  const config = useRuntimeConfig()
  const socket = ref<WebSocket | null>(null)
  const connected = ref(false)
  const handlers = new Map<string, Set<Handler>>()

  function resolveUrl(token: string): string {
    const base = config.public.wsBase
    // 若为相对路径则拼接当前 host，并按协议映射 ws/wss
    if (base.startsWith('ws')) return `${base}?token=${token}`
    if (import.meta.client) {
      const proto = location.protocol === 'https:' ? 'wss' : 'ws'
      return `${proto}://${location.host}${base}?token=${token}`
    }
    return `${base}?token=${token}`
  }

  function connect(token: string) {
    if (import.meta.server || socket.value) return
    const ws = new WebSocket(resolveUrl(token))
    socket.value = ws

    ws.onopen = () => {
      connected.value = true
    }
    ws.onclose = () => {
      connected.value = false
      socket.value = null
    }
    ws.onmessage = (ev: MessageEvent) => {
      try {
        const msg = JSON.parse(ev.data) as WsMessage
        handlers.get(msg.event)?.forEach((fn) => fn(msg.data))
      } catch {
        // 忽略非 JSON 消息
      }
    }
  }

  function disconnect() {
    socket.value?.close()
    socket.value = null
    connected.value = false
  }

  function send(event: WsClientEvent, data: unknown) {
    if (socket.value?.readyState === WebSocket.OPEN) {
      socket.value.send(JSON.stringify({ event, data }))
    }
  }

  /** 订阅服务端事件，返回取消订阅函数。 */
  function on<T = unknown>(event: WsEvent, handler: Handler<T>): () => void {
    if (!handlers.has(event)) handlers.set(event, new Set())
    handlers.get(event)!.add(handler as Handler)
    return () => handlers.get(event)?.delete(handler as Handler)
  }

  return { socket, connected, connect, disconnect, send, on }
}
