<script setup lang="ts">
/**
 * ThreadPanel — 线程回复面板。
 *
 * - 顶部展示父消息（头像/昵称/情绪徽标/内容/时间）；
 * - 中部线程回复列表（滚动，匿名/软删除占位处理）；
 * - 底部回复输入框 + 发送按钮（Ctrl/⌘+Enter 发送）；
 * - 右上角关闭按钮 emit close；提交回复 emit reply(content)。
 *
 * 断点适配：桌面内联、平板/移动由父级用样式覆盖（本组件填满容器高度）。
 */
import type { Message } from '~/stores/message'
import { useEmotions } from '~/composables/useEmotions'

const props = withDefaults(
  defineProps<{
    parentMessage: Message | null
    replies?: Message[]
    loading?: boolean
    sending?: boolean
  }>(),
  { replies: () => [], loading: false, sending: false },
)

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'reply', content: string): void
}>()

const { find: findEmotion } = useEmotions()

const draft = ref('')
const canSend = computed(() => draft.value.trim().length > 0 && !props.sending)

function avatarChar(m: Message): string {
  if (m.is_anonymous) return '匿'
  return m.author?.nickname?.charAt(0) ?? '牛'
}

function displayName(m: Message): string {
  if (m.is_anonymous) return '匿名牛马'
  return m.author?.nickname ?? '牛马'
}

function formatTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  return `${hh}:${mm}`
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
    e.preventDefault()
    submit()
  }
}

function submit() {
  const text = draft.value.trim()
  if (!text || props.sending) return
  emit('reply', text)
  draft.value = ''
}
</script>

<template>
  <aside
    class="glass-card flex h-full flex-col"
    role="dialog"
    aria-modal="true"
    aria-label="线程回复"
  >
    <!-- 头部 -->
    <header class="flex items-center justify-between border-b border-border px-4 py-3">
      <h2 class="flex items-center gap-2 text-md font-semibold">
        <AppIcon name="hash" :size="18" />
        线程
      </h2>
      <button
        type="button"
        class="rounded-md p-1 text-mute transition hover:text-text"
        aria-label="关闭线程面板"
        @click="emit('close')"
      >
        <AppIcon name="plus" :size="18" class="rotate-45" />
      </button>
    </header>

    <div class="flex flex-1 flex-col overflow-y-auto px-4 py-3">
      <!-- 父消息 -->
      <article v-if="parentMessage" class="flex gap-3 border-b border-border pb-4">
        <div
          class="grid h-10 w-10 shrink-0 place-items-center rounded-md text-sm font-semibold text-white"
          :class="parentMessage.is_anonymous ? 'bg-surface-hover text-dim' : 'bg-grad-ai'"
          aria-hidden="true"
        >
          <template v-if="!parentMessage.is_deleted">{{ avatarChar(parentMessage) }}</template>
        </div>
        <div class="min-w-0 flex-1">
          <p
            v-if="parentMessage.is_deleted"
            class="rounded-md border border-dashed border-border px-3 py-2 text-sm italic text-mute"
          >
            该消息已被删除
          </p>
          <template v-else>
            <div class="flex flex-wrap items-center gap-2">
              <span class="text-base font-semibold text-text">{{ displayName(parentMessage) }}</span>
              <span
                v-if="findEmotion(parentMessage.emotion)"
                class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
                :style="{
                  background: findEmotion(parentMessage.emotion)!.bg,
                  color: findEmotion(parentMessage.emotion)!.color,
                }"
              >
                {{ findEmotion(parentMessage.emotion)!.label }}
              </span>
              <span class="text-xs text-mute">{{ formatTime(parentMessage.created_at) }}</span>
            </div>
            <p class="mt-1 whitespace-pre-wrap break-words text-base leading-relaxed text-dim">
              {{ parentMessage.content }}
            </p>
          </template>
        </div>
      </article>

      <!-- 回复列表 -->
      <div class="flex flex-1 flex-col gap-4 pt-4">
        <p v-if="loading && !replies.length" class="py-2 text-center text-xs text-mute">
          加载中…
        </p>
        <p
          v-else-if="!replies.length"
          class="py-6 text-center text-sm text-mute"
        >
          还没有人回复，来说点什么陪陪 TA 吧。
        </p>

        <article
          v-for="r in replies"
          :key="r.id"
          class="flex gap-3"
        >
          <div
            class="grid h-8 w-8 shrink-0 place-items-center rounded-md text-xs font-semibold text-white"
            :class="r.is_anonymous ? 'bg-surface-hover text-dim' : 'bg-grad-ai'"
            aria-hidden="true"
          >
            <template v-if="!r.is_deleted">{{ avatarChar(r) }}</template>
          </div>
          <div class="min-w-0 flex-1">
            <p
              v-if="r.is_deleted"
              class="rounded-md border border-dashed border-border px-3 py-1.5 text-sm italic text-mute"
            >
              该回复已被删除
            </p>
            <template v-else>
              <div class="flex flex-wrap items-center gap-2">
                <span class="text-sm font-semibold text-text">{{ displayName(r) }}</span>
                <span
                  v-if="findEmotion(r.emotion)"
                  class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
                  :style="{ background: findEmotion(r.emotion)!.bg, color: findEmotion(r.emotion)!.color }"
                >
                  {{ findEmotion(r.emotion)!.label }}
                </span>
                <span class="text-xs text-mute">{{ formatTime(r.created_at) }}</span>
              </div>
              <p class="mt-0.5 whitespace-pre-wrap break-words text-sm leading-relaxed text-dim">
                {{ r.content }}
              </p>
            </template>
          </div>
        </article>
      </div>
    </div>

    <!-- 回复输入 -->
    <footer class="border-t border-border p-3">
      <div class="flex items-end gap-2">
        <textarea
          v-model="draft"
          rows="1"
          aria-label="回复线程"
          placeholder="回复这条线程…"
          class="max-h-28 min-h-[40px] flex-1 resize-none rounded-md border border-border bg-surface px-3 py-2 text-sm leading-relaxed text-text outline-none placeholder:text-mute focus:border-ai-1"
          @keydown="onKeydown"
        />
        <button
          type="button"
          class="btn-primary shrink-0 px-4 py-2 text-sm disabled:cursor-not-allowed disabled:opacity-50"
          :disabled="!canSend"
          @click="submit"
        >
          {{ sending ? '发送中…' : '回复' }}
        </button>
      </div>
    </footer>
  </aside>
</template>
