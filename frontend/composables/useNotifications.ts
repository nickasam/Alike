/**
 * useNotifications — 通知未读数与列表（跨组件共享未读态）。
 * 铃铛角标与通知中心页共用同一 unread ref，标记已读后即时同步。
 */
export interface Notification {
  id: number
  type: 'mention' | 'empathy' | 'reply' | 'system'
  content: string
  ref_id?: number
  is_read: boolean
  created_at: string
}

interface NotificationListResp {
  list: Notification[]
  total: number
  unread: number
  page: number
  page_size: number
}

export function useNotifications() {
  const api = useApi()
  // useState 保证 unread 在组件间共享（铃铛 + 通知页）。
  const unread = useState<number>('notif-unread', () => 0)

  async function refreshUnread() {
    try {
      const res = await api.get<NotificationListResp>('/notifications?page=1&page_size=1')
      unread.value = res?.unread ?? 0
    } catch {
      // 静默：未登录或网络异常时不打扰用户
    }
  }

  async function fetchList(page = 1, pageSize = 20) {
    const res = await api.get<NotificationListResp>(
      `/notifications?page=${page}&page_size=${pageSize}`,
    )
    unread.value = res?.unread ?? 0
    return res
  }

  async function markRead(id: number) {
    await api.put(`/notifications/${id}/read`)
    if (unread.value > 0) unread.value -= 1
  }

  async function markAllRead() {
    await api.put('/notifications/read-all')
    unread.value = 0
  }

  return { unread, refreshUnread, fetchList, markRead, markAllRead }
}
