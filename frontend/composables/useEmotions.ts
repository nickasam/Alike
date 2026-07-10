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
}

const EMOTIONS: EmotionMeta[] = [
  { key: 'tired', label: '疲惫', bg: 'rgba(148,163,184,.14)', color: '#cbd5e1' },
  { key: 'angry', label: '愤怒', bg: 'rgba(248,113,113,.14)', color: '#f87171' },
  { key: 'wronged', label: '委屈', bg: 'rgba(96,165,250,.14)', color: '#60a5fa' },
  { key: 'break', label: '崩溃', bg: 'rgba(167,139,250,.16)', color: '#a78bfa' },
  { key: 'numb', label: '麻木', bg: 'rgba(100,116,139,.18)', color: '#94a3b8' },
  { key: 'quit', label: '想润', bg: 'rgba(251,146,60,.14)', color: '#fb923c' },
  { key: 'anxious', label: '焦虑', bg: 'rgba(251,191,36,.14)', color: '#fbbf24' },
  { key: 'cheer', label: '加油', bg: 'rgba(52,211,153,.14)', color: '#34d399' },
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
