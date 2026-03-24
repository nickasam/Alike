# Alike 主应用布局设计文档

**项目**: Alike - LGBTQ+ 社交应用
**日期**: 2026-03-24
**状态**: 设计阶段

---

## 1. 概述

设计 Alike 应用的主布局架构，采用响应式平衡设计，桌面端使用侧边栏导航，移动端使用底部 Tab 栏导航。

### 1.1 目标

- 提供统一的应用布局框架
- 支持桌面端和移动端的响应式适配
- 包含 6 个核心功能区块的导航

### 1.2 导航区块

1. 🏠 首页 - 附近用户发现
2. 💬 聊天 - 全局聊天室
3. 💝 匹配 - 用户匹配系统
4. 📨 消息 - 消息列表
5. 🔔 通知 - 通知中心
6. 👤 我的 - 个人中心

---

## 2. 布局结构

### 2.1 桌面端 (>1024px)

```
┌─────────────────────────────────────────────────────────────────┐
│ 顶部导航栏 (64px)                                                │
│  Logo + 搜索框 + 消息/通知/用户头像                              │
├────────┬────────────────────────────────────────────────────────┤
│ 侧边栏 │         主内容区域                                     │
│ 240px  │         (可滚动)                                      │
│        │                                                        │
│ 用户卡 │         各页面内容渲染                                 │
│ ├────┤│                                                        │
│ 导航项││                                                        │
│      ││                                                        │
└────────┴────────────────────────────────────────────────────────┘
```

### 2.2 平板端 (768px-1024px)

```
┌─────────────────────────────────────────┐
│ 顶部导航栏 (56px)                        │
├─────────────────────────────────────────┤
│        主内容区域 (可滚动)               │
│                                         │
└─────────────────────────────────────────┘
│  [🏠][💬][💝][📨][👤] ← 底部Tab栏        │
└─────────────────────────────────────────┘
```

### 2.3 移动端 (<768px)

```
┌─────────────────────────────┐
│ 顶部导航栏 (56px)            │
├─────────────────────────────┤
│      主内容区域 (可滚动)     │
│                             │
└─────────────────────────────┘
│   [🏠][💬][💝][📨][👤]       │
└─────────────────────────────┘
```

---

## 3. 组件设计

### 3.1 MainLayout.vue

主布局容器，负责：
- 根据屏幕尺寸显示/隐藏侧边栏和底部Tab栏
- 管理全局状态（用户信息、未读数量等）
- 提供页面切换的过渡动画

**Props**: 无

**State**:
- `user`: 当前用户信息
- `unreadMessageCount`: 未读消息数
- `notificationCount`: 通知数

### 3.2 TopNavBar.vue

顶部导航栏组件

**元素**:
- Logo (心形图标 + "Alike" 文字)
- 搜索框
- 消息按钮 (带未读角标)
- 通知按钮 (带未读角标)
- 用户头像 (点击跳转个人中心)

**尺寸**: 桌面端 64px, 移动端 56px

### 3.3 Sidebar.vue

侧边栏组件（仅桌面端）

**元素**:
- 用户卡片 (头像 + 昵称 + 在线状态)
- 导航菜单 (图标 + 文字 + 徽章)
- 选中状态：左侧边条高亮

**宽度**: 240px (固定)

### 3.4 BottomTabBar.vue

底部Tab栏组件（仅移动端）

**元素**:
- 5个Tab项 (图标 + 文字 + 徽章)
- 通知功能移至"我的"页面内

**高度**: 64px

**导航项**:
1. 🏠 首页
2. 💬 聊天
3. 💝 匹配
4. 📨 消息
5. 👤 我的

### 3.5 UserCard.vue

用户卡片组件（复用于侧边栏和独立场景）

**元素**:
- 用户头像
- 昵称
- 在线状态指示器

---

## 4. 颜色系统

```css
/* 主色调 */
--primary-gradient: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
--accent-gradient: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);

/* 状态颜色 */
--success: #22c55e;
--warning: #f59e0b;
--error: #ef4444;

/* 导航栏 */
--navbar-bg: rgba(255, 255, 255, 0.95);
--navbar-border: rgba(0, 0, 0, 0.08);

/* 侧边栏 */
--sidebar-bg: #ffffff;
--sidebar-border: rgba(0, 0, 0, 0.06);
--sidebar-item-hover: rgba(102, 126, 234, 0.08);
--sidebar-item-active: rgba(102, 126, 234, 0.15);
--sidebar-indicator: #667eea;

/* 文字 */
--text-primary: #1a1a1a;
--text-secondary: #6b7280;
--text-tertiary: #9ca3af;
```

---

## 5. 响应式断点

```css
/* 移动端 (< 768px) */
- 隐藏侧边栏
- 显示底部Tab栏
- 顶部导航栏高度: 56px

/* 平板端 (768px - 1023px) */
- 隐藏侧边栏
- 显示底部Tab栏
- 顶部导航栏高度: 56px

/* 桌面端 (≥ 1024px) */
- 显示侧边栏 (240px固定宽度)
- 隐藏底部Tab栏
- 顶部导航栏高度: 64px
```

---

## 6. 路由配置

```javascript
const routes = [
  {
    path: '/',
    component: MainLayout,
    children: [
      { path: '', redirect: '/home' },
      { path: 'home', name: 'Home', component: () => import('@/views/Home.vue') },
      { path: 'chat', name: 'Chat', component: () => import('@/views/ChatRoom.vue') },
      { path: 'match', name: 'Match', component: () => import('@/views/Match.vue') },
      { path: 'messages', name: 'Messages', component: () => import('@/views/Messages.vue') },
      { path: 'notifications', name: 'Notifications', component: () => import('@/views/Notifications.vue') },
      { path: 'profile', name: 'Profile', component: () => import('@/views/Profile.vue') }
    ]
  }
]
```

---

## 7. 数据流

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  Pinia      │────▶│  MainLayout  │────▶│  子组件      │
│  userStore  │     │              │     │  (NavBar,   │
│             │     │              │     │   Sidebar)  │
└─────────────┘     └──────────────┘     └─────────────┘
       ▲                     │
       │                     │
       └─────────────────────┘
          用户操作更新状态
```

---

## 8. 文件结构

```
web/src/
├── components/
│   └── layout/
│       ├── MainLayout.vue        # 主布局容器
│       ├── TopNavBar.vue         # 顶部导航栏
│       ├── Sidebar.vue           # 侧边栏 (桌面端)
│       ├── BottomTabBar.vue      # 底部Tab栏 (移动端)
│       └── UserCard.vue          # 用户卡片组件
├── views/
│   ├── Home.vue                  # 首页 - 附近用户
│   ├── ChatRoom.vue              # 全局聊天室
│   ├── Match.vue                 # 匹配页面
│   ├── Messages.vue              # 消息列表
│   ├── Notifications.vue         # 通知
│   └── Profile.vue               # 我的
├── stores/
│   └── layout.js                 # 布局相关状态 (可选)
└── router/
    └── index.js                  # 路由配置
```

---

## 9. 实现顺序

1. 创建基础布局组件 (MainLayout, TopNavBar, Sidebar, BottomTabBar, UserCard)
2. 更新路由配置
3. 创建占位页面 (Home, ChatRoom, Match, Messages, Notifications, Profile)
4. 更新 App.vue 使用 MainLayout
5. 添加全局样式
6. 测试响应式行为

---

*文档版本: 1.0*
*创建日期: 2026-03-24*
