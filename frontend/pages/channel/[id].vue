<script setup lang="ts">
/**
 * 频道页 — 消息流 + 情绪看板 + 线程面板集成。
 *
 * 布局：
 *   - 左/主区：频道信息头 + MessageList + MessageInput；
 *   - 右侧（点击消息时）：ThreadPanel 覆盖式面板。
 *
 * 生命周期：
 *   - 进入：拉取频道信息、加载首屏消息、订阅 WS 事件、joinChannel；
 *   - 离开：leaveChannel、注销 WS 订阅、关闭线程面板。
 *
 * 消息发送经 WebSocket（send_message），后端广播 new_message 回显入列；
 * 线程回复经 REST（sendReply），后端广播 thread_reply 入列并同步计数。
 */
import { useMessageStore, type Message } from '~/stores/message'
import { useChannelStore, type Channel } from '~/stores/channel'

definePageMeta({ middleware: 'auth' })

const route = useRoute()
const channelId = computed(() => Number(route.params.id))

const messageStore = useMessageStore()
const channelStore = useChannelStore()
const ws = useWebSocket()
const api = useApi()
const auth = useAuthStore()

const channel = computed<Channel | null>(() => channelStore.current)
const threadOpen = computed(() => messageStore.threadOpen)
const sending = ref(false)
const replySending = ref(false)

useHead(() => ({ title: channel.value ? `#${channel.value.name} · Alike` : '频道 · Alike' }))

/** 拉取频道详情并写入 store（供 current getter 使用）。 */
async function loadChannel() {
  try {
    const data = await api.get<Channel>(`/channels/${channelId.value}`)
    const rest = channelStore.channels.filter((c) => c.id !== data.id)
    channelStore.setChannels([...rest, data])
  } catch {
    // 频道信息拉取失败不阻塞消息流
  } finally {
    channelStore.setCurrent(channelId.value)
  }
}

// —— WS 订阅解绑函数 —— //
let offMessage: (() => void) | null = null
let offThreadReply: (() => void) | null = null
let offEmpathy: (() => void) | null = null

function subscribeWs() {
  offMessage = ws.on<Message>('new_message', (msg) => messageStore.receiveMessage(msg))
  offThreadReply = ws.on('thread_reply', (payload) =>
    messageStore.receiveThreadReply(payload as any),
  )
  offEmpathy = ws.on('empathy', (payload) => messageStore.applyEmpathy(payload as any))
}

function unsubscribeWs() {
  offMessage?.()
  offThreadReply?.()
  offEmpathy?.()
  offMessage = offThreadReply = offEmpathy = null
}

/** 发送消息（经 WebSocket，new_message 回显入列）。 */
function onSend(payload: { content: string; emotion: string | null; isAnonymous: boolean }) {
  const clientMsgId =
    import.meta.client && 'randomUUID' in crypto ? crypto.randomUUID() : ''
  ws.sendMessage(channelId.value, payload.content, payload.emotion, clientMsgId, payload.isAnonymous)
}

/** 打开线程面板。 */
function onOpenThread(message: Message) {
  messageStore.openThread(message)
}

function onCloseThread() {
  messageStore.closeThread()
}

/** 线程回复（REST，thread_reply 广播入列）。 */
async function onReply(content: string) {
  const parent = messageStore.threadParent
  if (!parent || replySending.value) return
  replySending.value = true
  try {
    await messageStore.sendReply(parent.id, content, null, auth.user?.is_anonymous ?? false)
  } finally {
    replySending.value = false
  }
}

/** 共情（REST 增删，empathy 广播更新计数）。 */
async function onEmpathy({ message, action }: { message: Message; action: 'add' | 'remove' }) {
  try {
    if (action === 'add') {
      await api.post(`/messages/${message.id}/empathy`)
    } else {
      await api.del(`/messages/${message.id}/empathy`)
    }
  } catch {
    // 失败静默，计数以 WS empathy 广播为准
  }
}

onMounted(async () => {
  ws.connect()
  subscribeWs()
  await loadChannel()
  if (!messageStore.channelState(channelId.value).initialized) {
    await messageStore.loadInitial(channelId.value)
  }
  ws.joinChannel(channelId.value)
})

onBeforeUnmount(() => {
  ws.leaveChannel(channelId.value)
  unsubscribeWs()
  messageStore.closeThread()
})
</script>

<template>
  <div class="flex h-[calc(100vh-8rem)] gap-4">
    <!-- 主区：频道头 + 消息流 + 输入 -->
    <section class="flex min-w-0 flex-1 flex-col gap-3">
      <!-- 频道信息头 -->
      <header class="glass-card flex items-center justify-between gap-3 p-4">
        <div class="min-w-0">
          <h1 class="flex items-center gap-2 text-lg font-semibold">
            <AppIcon name="hash" :size="20" />
            <span class="truncate">{{ channel?.name ?? `频道 #${channelId}` }}</span>
          </h1>
          <p v-if="channel?.description" class="mt-0.5 truncate text-sm text-dim">
            {{ channel.description }}
          </p>
        </div>
        <span class="shrink-0 text-xs text-mute">
          {{ channel?.member_count ?? 0 }} 位牛马
        </span>
      </header>

      <!-- 消息流 -->
      <div class="glass-card flex min-h-0 flex-1 flex-col p-3">
        <MessageList
          :channel-id="channelId"
          @open-thread="onOpenThread"
          @empathy="onEmpathy"
        />
      </div>

      <!-- 输入框 -->
      <MessageInput
        :sending="sending"
        :default-anonymous="auth.user?.is_anonymous ?? false"
        @send="onSend"
      />
    </section>

    <!-- 线程面板：桌面右侧内联，窄屏覆盖 -->
    <aside
      v-if="threadOpen"
      class="fixed inset-0 z-40 bg-black/50 lg:static lg:z-auto lg:w-96 lg:bg-transparent"
      @click.self="onCloseThread"
    >
      <div class="ml-auto h-full w-full max-w-md lg:max-w-none">
        <ThreadPanel
          :parent-message="messageStore.threadParent"
          :replies="messageStore.threadReplies"
          :loading="messageStore.threadLoading"
          :sending="replySending"
          @close="onCloseThread"
          @reply="onReply"
        />
      </div>
    </aside>
  </div>
</template>
