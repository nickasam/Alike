<script setup lang="ts">
/**
 * MessageInput — 消息输入框。
 *
 * - 自适应高度 textarea（随内容增长，上限约 6 行）；
 * - 情绪标签选择器：点击展开 8 种情绪单选（useEmotions 元数据），再点清除；
 * - 匿名切换按钮：显示当前匿名/实名状态；
 * - 发送按钮：内容为空或发送中禁用；
 * - Ctrl/⌘+Enter 发送，Enter 换行；
 * - 发送后清空内容与情绪选择（匿名态保留，符合用户习惯）；
 * - emit send: { content, emotion, isAnonymous }。
 */
import { useEmotions } from '~/composables/useEmotions'

const props = withDefaults(
  defineProps<{
    /** 外部发送态：为 true 时禁用发送并显示发送中 */
    sending?: boolean
    /** 默认匿名（取自用户设置） */
    defaultAnonymous?: boolean
    placeholder?: string
  }>(),
  {
    sending: false,
    defaultAnonymous: false,
    placeholder: '说点什么，总有人懂你的辛苦...',
  },
)

const emit = defineEmits<{
  (e: 'send', payload: { content: string; emotion: string | null; isAnonymous: boolean }): void
  (e: 'typing'): void
}>()

const { emotions, find: findEmotion } = useEmotions()

const content = ref('')
const emotion = ref<string | null>(null)
const isAnonymous = ref(props.defaultAnonymous)
const pickerOpen = ref(false)
const textarea = ref<HTMLTextAreaElement | null>(null)

const selectedEmotion = computed(() => findEmotion(emotion.value))
const canSend = computed(() => content.value.trim().length > 0 && !props.sending)

/** 随内容自适应高度（上限 ~6 行 = 144px）。 */
function autoGrow() {
  const el = textarea.value
  if (!el) return
  el.style.height = 'auto'
  el.style.height = `${Math.min(el.scrollHeight, 144)}px`
}

function onInput() {
  autoGrow()
  emit('typing')
}

function onKeydown(e: KeyboardEvent) {
  // Ctrl/⌘+Enter 发送；纯 Enter 换行（浏览器默认行为）。
  if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
    e.preventDefault()
    submit()
  }
}

function pickEmotion(key: string) {
  emotion.value = emotion.value === key ? null : key
  pickerOpen.value = false
}

function clearEmotion() {
  emotion.value = null
}

function submit() {
  const text = content.value.trim()
  if (!text || props.sending) return
  emit('send', { content: text, emotion: emotion.value, isAnonymous: isAnonymous.value })
  content.value = ''
  emotion.value = null
  pickerOpen.value = false
  nextTick(autoGrow)
}
</script>

<template>
  <div class="glass-card flex flex-col gap-2 p-3">
    <!-- 顶部工具条：情绪 + 匿名 -->
    <div class="relative flex flex-wrap items-center gap-2">
      <!-- 情绪选择触发 -->
      <button
        type="button"
        class="inline-flex items-center gap-1 rounded-full border border-border px-3 py-1 text-sm transition hover:border-ai-1"
        :class="selectedEmotion ? 'text-text' : 'text-mute'"
        :style="selectedEmotion ? { background: selectedEmotion.bg, color: selectedEmotion.color, borderColor: 'transparent' } : undefined"
        :aria-expanded="pickerOpen"
        aria-haspopup="listbox"
        @click="pickerOpen = !pickerOpen"
      >
        <AppIcon name="sparkles" :size="14" />
        {{ selectedEmotion ? selectedEmotion.label : '加个情绪' }}
      </button>
      <button
        v-if="selectedEmotion"
        type="button"
        class="text-xs text-mute transition hover:text-text"
        aria-label="清除情绪"
        @click="clearEmotion"
      >
        清除
      </button>

      <!-- 匿名切换 -->
      <button
        type="button"
        class="ml-auto inline-flex items-center gap-1 rounded-full border px-3 py-1 text-sm transition"
        :class="isAnonymous ? 'border-ai-1 text-ai-1' : 'border-border text-mute hover:text-text'"
        :aria-pressed="isAnonymous"
        @click="isAnonymous = !isAnonymous"
      >
        <AppIcon name="user" :size="14" />
        {{ isAnonymous ? '匿名发送' : '实名发送' }}
      </button>

      <!-- 情绪选择弹层 -->
      <div
        v-if="pickerOpen"
        class="glass-card absolute bottom-full left-0 z-10 mb-2 flex max-w-xs flex-wrap gap-2 p-3"
        role="listbox"
        aria-label="情绪选择器"
      >
        <button
          v-for="e in emotions"
          :key="e.key"
          type="button"
          role="option"
          :aria-selected="emotion === e.key"
          class="rounded-full px-2.5 py-1 text-xs font-medium transition"
          :class="emotion === e.key ? 'ring-2 ring-ai-1' : ''"
          :style="{ background: e.bg, color: e.color }"
          @click="pickEmotion(e.key)"
        >
          {{ e.label }}
        </button>
      </div>
    </div>

    <!-- 输入区 -->
    <div class="flex items-end gap-2">
      <textarea
        ref="textarea"
        v-model="content"
        rows="1"
        aria-label="发送消息"
        :placeholder="placeholder"
        class="max-h-36 min-h-[40px] flex-1 resize-none rounded-md border border-border bg-surface px-3 py-2 text-base leading-relaxed text-text outline-none placeholder:text-mute focus:border-ai-1"
        @input="onInput"
        @keydown="onKeydown"
      />
      <button
        type="button"
        class="btn-primary shrink-0 px-4 py-2 text-base disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="!canSend"
        @click="submit"
      >
        {{ sending ? '发送中…' : '发送' }}
      </button>
    </div>
    <p class="text-xs text-mute">Ctrl/⌘ + Enter 发送 · Enter 换行</p>
  </div>
</template>
