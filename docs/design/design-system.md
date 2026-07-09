# Alike 设计系统 — Aurora 视觉语言

> **"清新、现代感、AI 感、智能感"** — 让打工人社区拥有未来感，而非廉价感。
>
> **默认暗色主题（Aurora 极光风），支持切换亮色模式。**

---

## 一、设计理念

### 1.1 核心关键词

| 关键词 | 含义 | 体现方式 |
|--------|------|----------|
| **清新** | 不沉闷、不压抑 | 极光渐变背景、充足留白、呼吸感 |
| **现代感** | 2025+ 设计趋势 | 玻璃拟态、毛玻璃模糊、微动效 |
| **AI 感** | 智能生成的精致感 | 光晕、粒子、流动渐变、数据可视化 |
| **智能感** | 信息层次清晰 | 智能布局、上下文感知、实时反馈 |
| **温度** | 不忘产品温度 | 共情暖色点缀、圆润边角、柔和阴影 |

### 1.2 设计原则

1. **光感优先** — 每个界面都有光源方向（左上→右下），元素有光泽、有投影
2. **层次分明** — 通过模糊、透明度、阴影创造 Z 轴深度，不用线框堆砌
3. **流动而非静止** — 渐变是活的，hover 有弹性，数据有动画
4. **克制使用色彩** — 深色底 + 1 个主色 + 1 个辅助色，其余用透明度
5. **SVG 而非 Emoji** — 所有图标用 SVG 图标系统，禁用 emoji 当 UI 元素
6. **圆角统一** — 全局 16px 圆角体系，组件内 12px，小元素 8px
7. **双主题** — 暗色为默认（Aurora 极光风），亮色为可切换模式，两套配色通过 CSS 变量无缝切换

---

## 二、色彩体系

### 2.1 暗色主题 Design Tokens（默认）

```css
:root {
  /* 深空背景 — 多层极光渐变 */
  --bg:            #0a0e1a;
  --bg-2:          #0b1220;
  --surface:       rgba(26, 34, 54, 0.6);
  --surface-solid: #111729;
  --surface-hover: #1a2236;
  --border:        rgba(255, 255, 255, 0.08);
  --border-strong: rgba(255, 255, 255, 0.14);
  --shadow:        0 8px 32px rgba(0, 0, 0, 0.4);

  /* 文字色阶 */
  --text:          #f1f5f9;
  --text-dim:      #94a3b8;
  --text-mute:     #64748b;
  --text-disabled: #475569;

  /* AI 极光色系 — 紫 → 蓝 → 青 */
  --ai-1:    #6366f1;  /* 极光蓝紫 */
  --ai-2:    #22d3ee;  /* 极光青 */
  --ai-3:    #a78bfa;  /* 极光紫 */
  --grad-ai: linear-gradient(135deg, #6366f1, #22d3ee);

  /* 温度色 — 牛马色（产品色） */
  --warm:        #fb923c;
  --warm-deep:   #f97316;
  --grad-warm:   linear-gradient(135deg, #fb923c, #f97316);

  /* 共情色 — 翠绿 */
  --empathy:      #34d399;
  --empathy-soft: rgba(52, 211, 153, 0.12);
  --grad-empathy: linear-gradient(135deg, #34d399, #10b981);

  /* 功能色 */
  --danger: #f87171;
  --info:   #60a5fa;
  --gold:   #fbbf24;

  /* 尺寸 */
  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 16px;
  --radius-xl: 20px;
  --radius-full: 999px;
  --nav-h: 64px;
}
```

### 2.2 亮色主题 Design Tokens（可切换）

```css
.light {
  /* 亮色背景 — 近白、通透 */
  --bg:            #FAFBFF;
  --bg-2:          #F0F2FA;
  --surface:       rgba(255, 255, 255, 0.72);
  --surface-solid: #FFFFFF;
  --surface-hover: #F5F6FC;
  --border:        rgba(124, 111, 240, 0.14);
  --border-strong: rgba(124, 111, 240, 0.22);
  --shadow:        0 8px 30px rgba(80, 80, 160, 0.08);

  /* 文字色阶 — 深靛蓝而非纯黑 */
  --text:          #1E2340;
  --text-dim:      #5B6178;
  --text-mute:     #9AA0B8;
  --text-disabled: #C8CCE0;

  /* AI 渐变 — 亮色下稍调高饱和 */
  --ai-1:    #7C6FF0;
  --ai-2:    #4A90FF;
  --ai-3:    #22D3EE;
  --grad-ai: linear-gradient(120deg, #7C6FF0, #4A90FF, #22D3EE);

  /* 温度色 */
  --warm:        #fb923c;
  --warm-deep:   #f97316;
  --grad-warm:   linear-gradient(135deg, #fb923c, #f97316);

  /* 共情色 */
  --empathy:      #14B8A6;
  --empathy-soft: rgba(20, 184, 166, 0.10);
  --grad-empathy: linear-gradient(135deg, #14B8A6, #2DD4BF);

  /* 功能色 */
  --danger: #FB7185;
  --info:   #60A5FA;
  --gold:   #F59E0B;
}
```

### 2.3 渐变系统

```css
/* 品牌渐变 — 按钮/Logo/关键强调 */
--grad-ai: linear-gradient(135deg, var(--ai-1), var(--ai-2));

/* 温度渐变 — 牛马主题 */
--grad-warm: linear-gradient(135deg, var(--warm), var(--warm-deep));

/* 共情渐变 */
--grad-empathy: linear-gradient(135deg, var(--empathy), var(--empathy-soft));

/* 文字渐变 — Hero 标题 */
--grad-text: linear-gradient(135deg, var(--ai-1), var(--ai-2), var(--ai-3));
```

### 2.4 玻璃拟态

```css
/* 暗色玻璃 */
--glass-bg:       rgba(26, 34, 54, 0.6);
--glass-border:   1px solid rgba(255, 255, 255, 0.08);
--glass-blur:     blur(20px);
--glass-shadow:   0 8px 32px rgba(0, 0, 0, 0.4);
--glass-highlight: inset 0 1px 0 rgba(255, 255, 255, 0.06);

/* 亮色玻璃（.light 下覆盖） */
.light {
  --glass-bg:       rgba(255, 255, 255, 0.72);
  --glass-border:   1px solid rgba(124, 111, 240, 0.14);
  --glass-shadow:   0 8px 30px rgba(80, 80, 160, 0.08);
  --glass-highlight: inset 0 1px 0 rgba(255, 255, 255, 0.8);
}
```

### 2.5 页面背景光晕

```css
/* 暗色背景 — 多层极光 radial 光晕 */
background:
  radial-gradient(ellipse 80% 60% at 20% 0%, rgba(99,102,241,.15), transparent 60%),
  radial-gradient(ellipse 60% 50% at 80% 10%, rgba(34,211,238,.10), transparent 60%),
  radial-gradient(ellipse 50% 40% at 50% 100%, rgba(167,139,250,.08), transparent 60%),
  #0a0e1a;

/* 亮色背景 — 柔和 radial 光晕网格 */
.light body {
  background:
    radial-gradient(ellipse 70% 50% at 15% 0%, rgba(124,111,240,.06), transparent 60%),
    radial-gradient(ellipse 60% 40% at 85% 10%, rgba(34,211,238,.05), transparent 60%),
    #FAFBFF;
}
```

---

## 三、字体体系

```css
/* 字体栈 */
--font-sans: 'Inter', -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Microsoft YaHei', sans-serif;
--font-mono: 'JetBrains Mono', 'Fira Code', monospace;

/* 字号 */
--text-xs:   11px;   /* 标签/徽章 */
--text-sm:   13px;   /* 辅助文字 */
--text-base: 14px;   /* 正文 */
--text-md:   16px;   /* 小标题 */
--text-lg:   18px;   /* 卡片标题 */
--text-xl:   22px;   /* 页面标题 */
--text-2xl:  28px;   /* Hero 标题 */
--text-3xl:  36px;   /* 数字大字 */

/* 字重 */
--fw-normal: 400;
--fw-medium: 500;
--fw-semibold: 600;
--fw-bold: 700;
--fw-black: 800;

/* 行高 */
--leading-tight: 1.3;
--leading-normal: 1.6;
--leading-relaxed: 1.8;
```

---

## 四、间距与圆角

```css
/* 间距 — 8px 栅格 */
--space-1: 4px;
--space-2: 8px;
--space-3: 12px;
--space-4: 16px;
--space-5: 20px;
--space-6: 24px;
--space-8: 32px;
--space-10: 40px;
--space-12: 48px;

/* 圆角 */
--radius-sm: 8px;     /* 小按钮/标签 */
--radius-md: 12px;    /* 输入框/小卡片 */
--radius-lg: 16px;    /* 卡片/面板（全局统一） */
--radius-xl: 20px;    /* 大卡片/Hero */
--radius-full: 999px; /* 胶囊/圆形 */
```

---

## 五、阴影系统

```css
/* 暗色阴影 */
--shadow-sm: 0 2px 8px rgba(0, 0, 0, 0.2);
--shadow-md: 0 4px 16px rgba(0, 0, 0, 0.3);
--shadow-lg: 0 8px 32px rgba(0, 0, 0, 0.4);
--shadow-xl: 0 16px 48px rgba(0, 0, 0, 0.5);

/* 亮色阴影（.light 下覆盖） */
.light {
  --shadow-sm: 0 2px 8px rgba(80, 80, 160, .06);
  --shadow-md: 0 4px 16px rgba(80, 80, 160, .08);
  --shadow-lg: 0 8px 30px rgba(80, 80, 160, .10);
  --shadow-xl: 0 16px 48px rgba(80, 80, 160, .12);
}

/* 光晕阴影 — AI 感核心 */
--glow-ai:      0 0 24px rgba(99, 102, 241, 0.4);
--glow-cyan:    0 0 24px rgba(34, 211, 238, 0.4);
--glow-warm:    0 0 24px rgba(251, 146, 60, 0.4);
--glow-empathy: 0 0 24px rgba(52, 211, 153, 0.4);

/* 内光 — 玻璃顶部高光 */
--inner-glow: inset 0 1px 0 rgba(255, 255, 255, 0.06);
.light { --inner-glow: inset 0 1px 0 rgba(255, 255, 255, 0.8); }
```

---

## 六、动效系统

### 6.1 缓动函数

```css
--ease-out: cubic-bezier(0.16, 1, 0.3, 1);       /* 弹性出场 */
--ease-in-out: cubic-bezier(0.65, 0, 0.35, 1);    /* 平滑过渡 */
--ease-spring: cubic-bezier(0.34, 1.56, 0.64, 1); /* 弹簧效果 */
```

### 6.2 时长

```
快速反馈    150ms    hover/active 状态
标准过渡    250ms    面板展开/折叠
入场动画    400ms    页面加载/卡片出现
重点动画    600ms    共情动画/数据更新
```

### 6.3 关键动效

| 动效 | 应用场景 | 实现 |
|------|---------|------|
| **渐入上浮** | 卡片/面板出现 | `opacity 0→1, translateY 12px→0, 400ms ease-out` |
| **光晕脉冲** | AI 特征元素 | `box-shadow 脉冲, 2s infinite` |
| **弹性缩放** | 按钮 hover | `scale(1.03), 200ms ease-spring` |
| **渐变流动** | Hero 背景 | `background-position 动画, 8s linear infinite` |
| **共情涟漪** | 抱团取暖点击 | `ripple + scale + glow, 600ms` |

---

## 七、SVG 图标系统

**禁用 emoji 作为 UI 元素。** 所有图标使用内联 SVG，统一描边风格。

### 图标规范

- **风格**: 线性图标 (Lucide / Heroicons 风格)
- **描边**: `stroke-width: 1.5`，`stroke-linecap: round`
- **尺寸**: 16px / 20px / 24px / 32px
- **颜色**: `stroke: currentColor`

### 核心图标清单

| 用途 | 图标名 | 替代原 emoji |
|------|--------|-------------|
| 搜索 | search | 🔍 |
| 通知 | bell | 🔔 |
| 消息 | mail / message-circle | ✉️ |
| 首页 | home | — |
| 频道 | hash | # |
| 日记 | book-open | 📓 |
| 排行榜 | trophy | 🏆 |
| 个人 | user | — |
| 共情 | heart-handshake | 🫂 |
| 情绪 | sparkles | ✨ |
| 发送 | send | ➤ |
| 回复 | corner-up-left | ↩ |
| 匿名 | eye-off | 🕶️ |
| 表情 | smile | 😊 |
| 图片 | image | 🖼️ |
| 火焰 | flame | 🔥 |
| 新建 | plus / pencil | ✍️ |
| 主题切换 | sun / moon | ☀️/🌙 |

---

## 八、组件视觉规范

### 8.1 玻璃卡片 (GlassCard)

```css
.glass-card {
  background: var(--glass-bg);
  backdrop-filter: var(--glass-blur);
  border: var(--glass-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--glass-shadow), var(--inner-glow);
  transition: all 250ms var(--ease-out);
}
.glass-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg), var(--inner-glow);
}
```

### 8.2 渐变按钮

```css
/* 主按钮 — AI 渐变 */
.btn-primary {
  background: var(--grad-ai);
  color: white;
  border-radius: var(--radius-md);
  box-shadow: var(--glow-ai);
  transition: all 250ms var(--ease-out);
}
.btn-primary:hover {
  transform: translateY(-2px) scale(1.02);
  box-shadow: 0 0 32px rgba(99, 102, 241, 0.6);
}

/* 温度按钮 (牛马主题) */
.btn-warm {
  background: var(--grad-warm);
  color: #1a0f00;
  box-shadow: var(--glow-warm);
}
```

### 8.3 头像

```css
.avatar {
  border-radius: var(--radius-md);
  background: var(--grad-ai);
  position: relative;
}
.avatar::after {
  content: '';
  position: absolute;
  bottom: -2px; right: -2px;
  width: 12px; height: 12px;
  border-radius: 50%;
  background: var(--ai-2);
  box-shadow: var(--glow-cyan);
  border: 2px solid var(--bg);
}
```

### 8.4 情绪标签

```css
.emo-tag {
  padding: 4px 12px;
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: var(--fw-semibold);
  background: rgba(99, 102, 241, 0.12);
  border: 1px solid rgba(99, 102, 241, 0.25);
  color: #a5b4fc;
  transition: all 200ms var(--ease-out);
}
.emo-tag:hover {
  transform: translateY(-2px);
  border-color: rgba(99, 102, 241, 0.5);
  box-shadow: var(--glow-ai);
}
```

### 8.5 共情按钮（核心特色）

```css
.empathy-btn {
  background: var(--empathy-soft);
  border: 1px solid rgba(52, 211, 153, 0.3);
  border-radius: var(--radius-full);
  color: var(--empathy);
  transition: all 250ms var(--ease-out);
}
.empathy-btn:hover {
  background: rgba(52, 211, 153, 0.2);
  box-shadow: var(--glow-empathy);
  transform: scale(1.05);
}
.empathy-btn.active {
  background: var(--grad-empathy);
  color: #052e20;
  box-shadow: 0 0 20px rgba(52, 211, 153, 0.5);
}
```

### 8.6 主题切换按钮

```css
.theme-toggle {
  width: 40px; height: 40px;
  border-radius: var(--radius-md);
  background: var(--surface);
  border: var(--glass-border);
  display: grid; place-items: center;
  cursor: pointer;
  transition: all 200ms var(--ease-out);
}
.theme-toggle:hover {
  background: var(--surface-hover);
  box-shadow: var(--glow-ai);
}
/* 暗色模式显示太阳图标（点击切亮色），亮色模式显示月亮图标（点击切暗色） */
```

---

## 九、布局系统

### 9.1 全局布局

- **最大宽度**: 1440px 居中
- **导航栏高度**: 64px
- **侧边栏宽度**: 260px (可折叠到 72px)
- **右侧栏宽度**: 320px

### 9.2 响应式断点

| 断点 | 宽度 | 布局变化 |
|------|------|---------|
| 桌面 | ≥1280px | 三列完整布局 |
| 平板 | 768-1279px | 两列，右侧栏隐藏 |
| 手机 | <768px | 单列，侧边栏抽屉 |

---

## 十、主题切换实现

### 10.1 CSS 变量切换

所有颜色通过 CSS 自定义属性（`:root` / `.light`）定义，切换主题只需在 `<html>` 或 `<body>` 上切换 `.light` class：

```css
/* 默认暗色 */
:root { ... }

/* 亮色覆盖 */
.light { ... }
```

### 10.2 JavaScript 切换

```js
// 切换主题
function toggleTheme() {
  document.documentElement.classList.toggle('light');
  // 持久化到 localStorage
  const isLight = document.documentElement.classList.contains('light');
  localStorage.setItem('theme', isLight ? 'light' : 'dark');
}

// 初始化时读取偏好
const saved = localStorage.getItem('theme');
if (saved === 'light') {
  document.documentElement.classList.add('light');
}
```

### 10.3 HTML 原型中的切换

每个 HTML 原型页面导航栏右侧包含一个主题切换按钮（太阳/月亮 SVG 图标），点击切换 `.light` class。

---

## 十一、与旧设计的对比

| 维度 | 旧设计（丑深色） | 新设计（Aurora 极光风） |
|------|-----------------|----------------------|
| 背景 | 纯色 #0f172a | 极光多层渐变 #0a0e1a + radial 光晕 |
| 图标 | Emoji | SVG 线性图标 |
| 卡片 | 实色 + 边框 | 玻璃拟态 + 模糊 |
| 按钮 | 橙色扁平 | 品牌渐变 + 光晕 |
| 阴影 | 通用阴影 | 光晕 + 内光 |
| 动效 | 仅颜色变化 | 弹性 + 光晕 + 流动 |
| 层次 | 2D 平面 | Z 轴深度 |
| 品牌色 | 橙色单色 | 蓝青紫极光色系 |
| 文字 | 纯色 | 渐变标题文字 |
| 主题 | 仅暗色 | 暗色默认 + 亮色可切换 |
| 整体 | 2020教程风 | 2025 AI产品风 |
