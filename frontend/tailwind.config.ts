import type { Config } from 'tailwindcss'

/**
 * Alike Tailwind 配置 — Aurora 极光设计系统落地。
 *
 * 颜色统一走 CSS 变量（见 assets/css/main.css 的 :root / .light），
 * 这样暗色（默认）与亮色主题可无缝切换，Tailwind 只做语义映射。
 * 用法示例：bg-surface / text-dim / text-ai-1 / rounded-lg / shadow-glow-ai
 */
export default <Partial<Config>>{
  darkMode: 'class',
  content: [
    './components/**/*.{vue,js,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './composables/**/*.{js,ts}',
    './stores/**/*.{js,ts}',
    './app.vue',
    './error.vue',
  ],
  theme: {
    extend: {
      colors: {
        // 背景 / 表面
        bg: 'var(--bg)',
        'bg-2': 'var(--bg-2)',
        surface: 'var(--surface)',
        'surface-solid': 'var(--surface-solid)',
        'surface-hover': 'var(--surface-hover)',
        border: 'var(--border)',
        'border-strong': 'var(--border-strong)',
        // 文字色阶
        text: 'var(--text)',
        dim: 'var(--text-dim)',
        mute: 'var(--text-mute)',
        disabled: 'var(--text-disabled)',
        // AI 极光色系
        'ai-1': 'var(--ai-1)',
        'ai-2': 'var(--ai-2)',
        'ai-3': 'var(--ai-3)',
        // 温度色（牛马色）
        warm: 'var(--warm)',
        'warm-deep': 'var(--warm-deep)',
        // 共情色
        empathy: 'var(--empathy)',
        'empathy-soft': 'var(--empathy-soft)',
        // 功能色
        danger: 'var(--danger)',
        info: 'var(--info)',
        gold: 'var(--gold)',
      },
      backgroundImage: {
        'grad-ai': 'var(--grad-ai)',
        'grad-warm': 'var(--grad-warm)',
        'grad-empathy': 'var(--grad-empathy)',
        'grad-text': 'var(--grad-text)',
      },
      fontFamily: {
        sans: [
          'Inter',
          '-apple-system',
          'BlinkMacSystemFont',
          'PingFang SC',
          'Microsoft YaHei',
          'sans-serif',
        ],
        mono: ['JetBrains Mono', 'Fira Code', 'monospace'],
      },
      fontSize: {
        xs: '11px',
        sm: '13px',
        base: '14px',
        md: '16px',
        lg: '18px',
        xl: '22px',
        '2xl': '28px',
        '3xl': '36px',
      },
      spacing: {
        1: '4px',
        2: '8px',
        3: '12px',
        4: '16px',
        5: '20px',
        6: '24px',
        8: '32px',
        10: '40px',
        12: '48px',
        'nav-h': 'var(--nav-h)',
        sidebar: '260px',
        'sidebar-collapsed': '72px',
        aside: '320px',
      },
      borderRadius: {
        sm: '8px',
        md: '12px',
        lg: '16px',
        xl: '20px',
        full: '999px',
      },
      boxShadow: {
        sm: 'var(--shadow-sm)',
        md: 'var(--shadow-md)',
        lg: 'var(--shadow-lg)',
        xl: 'var(--shadow-xl)',
        glass: 'var(--glass-shadow)',
        'glow-ai': 'var(--glow-ai)',
        'glow-cyan': 'var(--glow-cyan)',
        'glow-warm': 'var(--glow-warm)',
        'glow-empathy': 'var(--glow-empathy)',
      },
      backdropBlur: {
        glass: '20px',
      },
      maxWidth: {
        app: '1440px',
        content: '760px',
      },
      transitionTimingFunction: {
        out: 'cubic-bezier(0.16, 1, 0.3, 1)',
        'in-out': 'cubic-bezier(0.65, 0, 0.35, 1)',
        spring: 'cubic-bezier(0.34, 1.56, 0.64, 1)',
      },
      transitionDuration: {
        fast: '150ms',
        std: '250ms',
        enter: '400ms',
        focus: '600ms',
      },
      screens: {
        // 交互规范：md=768(平板)、xl=1280(桌面三列)
        md: '768px',
        xl: '1280px',
      },
      keyframes: {
        'rise-in': {
          '0%': { opacity: '0', transform: 'translateY(12px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        'glow-pulse': {
          '0%, 100%': { boxShadow: 'var(--glow-ai)' },
          '50%': { boxShadow: '0 0 40px rgba(99,102,241,0.6)' },
        },
      },
      animation: {
        'rise-in': 'rise-in 400ms var(--ease-out) both',
        'glow-pulse': 'glow-pulse 2s ease-in-out infinite',
      },
    },
  },
  plugins: [],
}
