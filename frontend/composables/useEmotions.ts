/**
 * useEmotions — 情绪标签元数据（key / 中文标签 / 展示样式）。
 *
 * key 与后端 emotion.Tag 枚举严格一致（8 种），颜色取自 docs/design/02-channel.html
 * 的 .emo-* 规则，供 MessageInput、MessageList 情绪徽标复用。
 */

export interface EmotionMeta {
  /** 与后端一致的情绪 key */
  key: string
  /** 中文标签 */
  label: string
  /** 徽标背景色（含透明度） */
  bg: string
  /** 徽标文字/边框色 */
  color: string
  /** 图标 SVG path（viewBox 0 0 24 24，Lucide 风格），取自 02-channel 设计稿 */
  icon: string[]
}

const EMOTIONS: EmotionMeta[] = [
  { key: 'tired', label: '疲惫', bg: 'rgba(148,163,184,.14)', color: '#cbd5e1', icon: ['M17.7 7.7a2.5 2.5 0 1 1 1.8 4.3H2', 'M9.6 4.6A2 2 0 1 1 11 8H2', 'M12.6 19.4A2 2 0 1 0 14 16H2'] },
  { key: 'angry', label: '愤怒', bg: 'rgba(248,113,113,.14)', color: '#f87171', icon: ['M8.5 14.5A2.5 2.5 0 0 0 11 12c0-1.38-.5-2-1-3-1.072-2.143-.224-4.054 2-6 .5 2.5 2 4.9 4 6.5 2 1.6 3 3.5 3 5.5a7 7 0 1 1-14 0c0-1.153.433-2.294 1-3a2.5 2.5 0 0 0 2.5 2.5z'] },
  { key: 'wronged', label: '委屈', bg: 'rgba(96,165,250,.14)', color: '#60a5fa', icon: ['M12 22a7 7 0 0 0 7-7c0-2-1-3.9-3-5.5s-3.5-4-4-6.5c-.5 2.5-2 4.5-4 6.5S5 13 5 15a7 7 0 0 0 7 7z'] },
  { key: 'break', label: '崩溃', bg: 'rgba(167,139,250,.16)', color: '#a78bfa', icon: ['M13 2 3 14l9 0-1 8 10-12-9 0 1-8z'] },
  { key: 'numb', label: '麻木', bg: 'rgba(100,116,139,.18)', color: '#94a3b8', icon: ['M22 12h-4l-3 9L9 3l-3 9H2'] },
  { key: 'quit', label: '想润', bg: 'rgba(251,146,60,.14)', color: '#fb923c', icon: ['M17.8 19.2 16 11l3.5-3.5C21 6 21.5 4 21 3c-1-.5-3 0-4.5 1.5L13 8 4.8 6.2c-.5-.1-.9.1-1.1.5l-.3.5c-.2.5-.1 1 .3 1.3L9 12l-2 3H4l-1 1 3 2 2 3 1-1v-3l3-2 3.5 5.3c.3.4.8.5 1.3.3l.5-.2c.4-.3.6-.7.5-1.2z'] },
  { key: 'anxious', label: '焦虑', bg: 'rgba(251,191,36,.14)', color: '#fbbf24', icon: ['M22 12h-4l-3 9L9 3l-3 9H2'] },
  { key: 'cheer', label: '加油', bg: 'rgba(52,211,153,.14)', color: '#34d399', icon: ['M9 11H5a2 2 0 0 0-2 2v7a2 2 0 0 0 2 2h4m5-11 1-4a3 3 0 0 0-3-3l-4 8v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3z'] },
]

const byKey: Record<string, EmotionMeta> = Object.fromEntries(
  EMOTIONS.map((e) => [e.key, e]),
)

export function useEmotions() {
  /** 依据 key 查询元数据，未知 key 返回 undefined。 */
  function find(key?: string | null): EmotionMeta | undefined {
    if (!key) return undefined
    return byKey[key]
  }

  return { emotions: EMOTIONS, find }
}
