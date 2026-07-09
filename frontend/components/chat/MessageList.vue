<script setup lang="ts">
/**
 * MessageList — 频道消息流。
 *
 * - 从 message store 读取按 created_at 升序的消息列表；
 * - 渲染头像/昵称/情绪徽标/内容/时间/共情数/回复数；
 * - 匿名消息隐藏头像与昵称，显示「匿名牛马」；软删除消息显示占位文案；
 * - 触顶时调用 store.loadMore 加载更早历史（保持滚动位置）；
 * - 新消息（列表尾部增长且用户贴近底部）自动滚到底部；
 * - 点击消息 emit open-thread，供父级打开线程面板。
 */
import { useMessageStore, type Message } from '~/stores/message'
import { useEmotions } from '~/composables/useEmotions'

const props = defineProps<{ channelId: number }>()

const emit = defineEmits<{
  (e: 'open-thread', message: Message): void
  (e: 'empathy', payload: { message: Message; action: 'add' | 'remove' }): void
}>()

const store = useMessageStore()
const { find: findEmotion } = useEmotions()

const messages = computed(() => store.listOf(props.channelId))
const chState = computed(() => store.channelState(props.channelId))

const scroller = ref<HTMLElement | null>(null)
/** 贴近底部阈值（px），用户上翻查看历史时不强制滚底。 */
const NEAR_BOTTOM = 120

function isNearBottom(): boolean {
  const el = scroller.value
  if (!el) return true
  return el.scrollHeight - el.scrollTop - el.clientHeight < NEAR_BOTTOM
}

function scrollToBottom() {
  const el = scroller.value
  if (el) el.scrollTop = el.scrollHeight
}

/** 触顶加载更早历史，加载后补偿滚动位置避免跳动。 */
async function onScroll() {
  const el = scroller.value
  if (!el || el.scrollTop > 8) return
  if (!chState.value.hasMore || chState.value.loading) return
  const prevHeight = el.scrollHeight
  await store.loadMore(props.channelId)
  await nextTick()
  el.scrollTop = el.scrollHeight - prevHeight
}

// 新消息到达且用户贴近底部时自动滚底。
watch(
  () => messages.value.length,
  async (len, prev) => {
    if (len > (prev ?? 0) && isNearBottom()) {
      await nextTick()
      scrollToBottom()
    }
  },
)

onMounted(async () => {
  if (!chState.value.initialized) {
    await store.loadInitial(props.channelId)
  }
  await nextTick()
  scrollToBottom()
})

/** 头像首字（昵称首字符），匿名固定「匿」。 */
function avatarChar(m: Message): string {
  if (m.is_anonymous) return '匿'
  return m.author?.nickname?.charAt(0) ?? '牛'
}

/** HH:mm 展示时间。 */
function formatTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  return `${hh}:${mm}`
}

function displayName(m: Message): string {
  if (m.is_anonymous) return '匿名牛马'
  return m.author?.nickname ?? '牛马'
}

function onEmpathy(m: Message, payload: { action: 'add' | 'remove' }) {
  emit('empathy', { message: m, action: payload.action })
}
</script>

<template>
  <div
    ref="scroller"
    class="flex flex-1 flex-col gap-4 overflow-y-auto px-1 py-2"
    aria-label="消息列表"
    @scroll="onScroll"
  >
    <!-- 加载更早历史提示 -->
    <p
      v-if="chState.loading && messages.length"
      class="py-1 text-center text-xs text-mute"
    >
      加载中…
    </p>
    <p
      v-else-if="!chState.hasMore && messages.length"
      class="py-1 text-center text-xs text-mute"
    >
      — 已经到顶了 —
    </p>

    <!-- 空态 -->
    <div
      v-if="chState.initialized && !messages.length"
      class="flex flex-1 flex-col items-center justify-center gap-2 text-center text-mute"
    >
      <AppIcon name="sparkles" :size="28" />
      <p class="text-sm">还没有人说话，来发第一条吧，总有人懂你的辛苦。</p>
    </div>

    <!-- 消息项 -->
    <article
      v-for="m in messages"
      :key="m.id"
      class="msg group flex gap-3"
    >
      <!-- 头像：软删除/匿名不展示真实头像 -->
      <div
        class="grid h-10 w-10 shrink-0 place-items-center rounded-md text-sm font-semibold text-white"
        :class="m.is_anonymous ? 'bg-surface-hover text-dim' : 'bg-grad-ai'"
        aria-hidden="true"
      >
        <template v-if="!m.is_deleted">{{ avatarChar(m) }}</template>
      </div>

      <div class="min-w-0 flex-1">
        <!-- 软删除占位 -->
        <p
          v-if="m.is_deleted"
          class="rounded-md border border-dashed border-border px-3 py-2 text-sm italic text-mute"
        >
          该消息已被删除
        </p>

        <template v-else>
          <div class="flex flex-wrap items-center gap-2">
            <span class="text-base font-semibold text-text">{{ displayName(m) }}</span>
            <span
              v-if="m.is_anonymous"
              class="inline-flex items-center gap-1 text-xs text-mute"
            >
              <AppIcon name="user" :size="12" />已隐身
            </span>
            <!-- 情绪徽标 -->
            <span
              v-if="findEmotion(m.emotion)"
              class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
              :style="{
                background: findEmotion(m.emotion)!.bg,
                color: findEmotion(m.emotion)!.color,
              }"
            >
              {{ findEmotion(m.emotion)!.label }}
            </span>
            <span class="text-xs text-mute">{{ formatTime(m.created_at) }}</span>
          </div>

          <p class="mt-1 whitespace-pre-wrap break-words text-base leading-relaxed text-dim">
            {{ m.content }}
          </p>

          <div class="mt-2 flex flex-wrap items-center gap-3">
            <EmpathyButton
              :count="m.empathy_count"
              :empathized="m.empathized"
              size="sm"
              @empathy="onEmpathy(m, $event)"
            />
            <button
              type="button"
              class="inline-flex items-center gap-1 rounded-full px-2 py-1 text-sm text-mute transition hover:text-ai-1"
              :aria-label="`查看线程，${m.reply_count} 条回复`"
              @click="emit('open-thread', m)"
            >
              <AppIcon name="hash" :size="14" />
              {{ m.reply_count }} 条回复
            </button>
          </div>
        </template>
      </div>
    </article>
  </div>
</template>
