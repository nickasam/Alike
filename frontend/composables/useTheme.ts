/**
 * useTheme — Aurora 双主题切换（暗色默认 / 亮色可切换）。
 * 主题通过 <html> 上的 dark / light class 控制，持久化到 localStorage。
 */
export type ThemeMode = 'dark' | 'light'

export function useTheme() {
  const theme = useState<ThemeMode>('theme', () => 'dark')

  function apply(mode: ThemeMode) {
    if (import.meta.client) {
      const el = document.documentElement
      el.classList.toggle('light', mode === 'light')
      el.classList.toggle('dark', mode === 'dark')
      localStorage.setItem('alike_theme', mode)
    }
    theme.value = mode
  }

  function toggle() {
    apply(theme.value === 'dark' ? 'light' : 'dark')
  }

  /** 从 localStorage 初始化（在客户端挂载时调用）。 */
  function init() {
    if (import.meta.client) {
      const saved = localStorage.getItem('alike_theme') as ThemeMode | null
      apply(saved === 'light' ? 'light' : 'dark')
    }
  }

  return { theme, toggle, apply, init }
}
