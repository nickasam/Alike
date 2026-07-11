import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useMessageStore, type Message } from '~/stores/message'

function makeMessage(over: Partial<Message> = {}): Message {
  return {
    id: 1,
    channel_id: 10,
    parent_id: null,
    content: '今天又加班到十点',
    is_anonymous: false,
    empathy_count: 0,
    reply_count: 0,
    is_deleted: false,
    created_at: '2026-07-11T10:00:00Z',
    ...over,
  }
}

function makeReply(id: number, parentId: number): Message {
  return makeMessage({
    id,
    parent_id: parentId,
    content: `回复 ${id}`,
    created_at: '2026-07-11T10:01:00Z',
  })
}

describe('message store — reply_count 同步', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('打开线程后回复一次，父消息 reply_count 只 +1（不因 threadParent 与列表项同引用而 +2）', () => {
    const store = useMessageStore()
    const parent = makeMessage({ id: 1, reply_count: 0 })
    // 列表项与 threadParent 指向同一对象（openThread 直接传入列表项）
    store.ensureChannel(10).list.push(parent)
    store.threadOpen = true
    store.threadParent = parent

    store.addThreadReply(1, makeReply(100, 1))

    expect(parent.reply_count).toBe(1)
    expect(store.threadParent?.reply_count).toBe(1)
    expect(store.threadReplies).toHaveLength(1)
  })

  it('同一条回复重复入列（REST 响应 + WS 广播）只计数一次', () => {
    const store = useMessageStore()
    const parent = makeMessage({ id: 1, reply_count: 0 })
    store.ensureChannel(10).list.push(parent)
    store.threadOpen = true
    store.threadParent = parent

    const reply = makeReply(100, 1)
    store.addThreadReply(1, reply) // REST 响应即时入列
    store.receiveThreadReply({ parent_id: 1, reply }) // WS 广播回环

    expect(parent.reply_count).toBe(1)
    expect(store.threadReplies).toHaveLength(1)
  })

  it('线程面板未打开时收到广播回复，父消息仍正确 +1', () => {
    const store = useMessageStore()
    const parent = makeMessage({ id: 1, reply_count: 2 })
    store.ensureChannel(10).list.push(parent)
    // 未打开线程：threadParent 为 null，threadOpen=false

    store.receiveThreadReply({ parent_id: 1, reply: makeReply(100, 1) })

    expect(parent.reply_count).toBe(3)
    expect(store.threadReplies).toHaveLength(0)
  })

  it('乐观发送的消息被 WS 回环按 client_msg_id 替换，不重复', () => {
    const store = useMessageStore()
    store.addOptimistic(10, makeMessage({ id: -1, client_msg_id: 'abc', pending: true }))
    store.receiveMessage(makeMessage({ id: 5, client_msg_id: 'abc' }))

    const list = store.listOf(10)
    expect(list).toHaveLength(1)
    expect(list[0].id).toBe(5)
    expect(list[0].pending).toBeUndefined()
  })
})
