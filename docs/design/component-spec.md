# Alike 前端组件规范文档

> **项目：** Alike — 汇聚天下牛马，总有人懂你的辛苦
>
> **技术栈：** Vue 3 + Nuxt 3 + TypeScript + Tailwind CSS + Pinia
>
> **编码约定：** `<script setup lang="ts">` 组合式 API，Props 用 `defineProps<T>()`，Emits 用 `defineEmits<T>()`，样式用 Tailwind CSS（默认暗色主题 Aurora 极光风，支持 `dark:` / `.light` 切换亮色）。视觉规范见 `docs/design/design-system.md`。

---

## 目录

| # | 组件 | 路径 |
|---|------|------|
| 1 | [MessageList](#1-messagelist) | `components/chat/MessageList.vue` |
| 2 | [MessageInput](#2-messageinput) | `components/chat/MessageInput.vue` |
| 3 | [ThreadPanel](#3-threadpanel) | `components/chat/ThreadPanel.vue` |
| 4 | [EmotionPicker](#4-emotionpicker) | `components/emotion/EmotionPicker.vue` |
| 5 | [EmotionBoard](#5-emotionboard) | `components/emotion/EmotionBoard.vue` |
| 6 | [EmpathyButton](#6-empathybutton) | `components/empathy/EmpathyButton.vue` |
| 7 | [ChannelSidebar](#7-channelsidebar) | `components/layout/ChannelSidebar.vue` |
| 8 | [TopNav](#8-topnav) | `components/layout/TopNav.vue` |

---

## 1. MessageList

**组件名：** `MessageList`
**路径：** `frontend/components/chat/MessageList.vue`

消息列表是频道页面的核心展示组件，负责渲染频道内的消息流，支持向上滚动分页加载历史消息、新消息实时追加、自动滚动到底部等功能。该组件接收频道 ID 和消息数组作为 Props，通过 WebSocket 推送或 REST API 分页拉取数据。每条消息项展示发送者头像、昵称（或"匿名牛马"标识）、消息内容、情绪标签、共情按钮入口以及线程回复数。组件内部维护滚动位置状态，在加载历史消息后保持视觉位置不跳变。消息列表还处理"有新消息"提示条，当用户不在底部时显示未读提示，点击可平滑滚动至最新消息。

### Props 定义

| 名称 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `channelId` | `number` | — (必填) | 当前频道 ID，用于拉取该频道的消息 |
| `messages` | `Message[]` | `[]` | 消息数组，每项含 `id/userId/nickname/avatar/content/emotion/isAnonymous/empathyCount/createdAt/parentId/threadCount` |
| `loading` | `boolean` | `false` | 是否正在加载历史消息（显示顶部加载指示器） |
| `hasMore` | `boolean` | `true` | 是否还有更多历史消息可加载 |
| `currentUserId` | `number \| null` | `null` | 当前登录用户 ID，用于判断是否为自己发的消息（右侧对齐） |

### Emits 定义

| 事件名 | 载荷 | 说明 |
|--------|------|------|
| `load-more` | `void` | 滚动到顶部时触发，父组件应拉取更早的历史消息 |
| `open-thread` | `{ messageId: number }` | 点击消息的线程区域时触发，打开线程面板 |
| `empathy` | `{ messageId: number }` | 点击消息上的共情按钮时触发 |
| `scroll-to-bottom` | `void` | 用户手动点击"新消息"提示条，滚动到底部 |

### Slots

| 名称 | 说明 |
|------|------|
| `message-item` | 自定义单条消息的渲染（作用域插槽，提供 `message` 和 `index`） |
| `empty` | 频道暂无消息时的空状态占位内容 |

### 内部状态

| 状态 | 类型 | 说明 |
|------|------|------|
| `scrollTop` | `number` | 当前滚动位置，用于判断是否到达顶部触发分页 |
| `isAtBottom` | `boolean` | 是否滚动到底部，决定新消息提示条的显示 |
| `unreadCount` | `number` | 不在底部时累计的新消息数 |
| `scrollContainerRef` | `HTMLElement \| null` | 滚动容器 DOM 引用 |

### UI 结构概要

```html
<div class="message-list-container" ref="scrollContainerRef" @scroll="handleScroll">
  <!-- 顶部加载指示器 -->
  <div v-if="loading" class="loading-indicator">加载中...</div>

  <!-- 空状态 -->
  <slot v-if="messages.length === 0 && !loading" name="empty">
    <div class="empty-state">还没有人说话，来当第一个牛马吧 🐮</div>
  </slot>

  <!-- 消息列表 -->
  <div v-for="msg in messages" :key="msg.id" class="message-item" :class="{ 'is-self': msg.userId === currentUserId }">
    <slot name="message-item" :message="msg">
      <img :src="msg.avatar" class="avatar" />
      <div class="message-body">
        <span class="nickname">{{ msg.isAnonymous ? '匿名牛马' : msg.nickname }}</span>
        <span v-if="msg.emotion" class="emotion-tag">{{ msg.emotion }}</span>
        <p class="content">{{ msg.content }}</p>
        <div class="message-footer">
          <button @click="$emit('empathy', { messageId: msg.id })">🤝 {{ msg.empathyCount }}</button>
          <button v-if="msg.threadCount > 0" @click="$emit('open-thread', { messageId: msg.id })">
            💬 {{ msg.threadCount }} 回复
          </button>
          <span class="timestamp">{{ formatTime(msg.createdAt) }}</span>
        </div>
      </div>
    </slot>
  </div>

  <!-- 新消息提示条 -->
  <div v-if="!isAtBottom && unreadCount > 0" class="new-message-bar" @click="scrollToBottom">
    ↓ {{ unreadCount }} 条新消息
  </div>
</div>
```

### 交互行为说明

1. **向上滚动分页**：当 `scrollTop` 接近 0 且 `hasMore` 为 `true` 时，触发 `load-more` 事件。加载完成后记录原 scrollHeight，在新消息插入后恢复滚动位置，避免视觉跳动。
2. **新消息追加**：通过 Props 的 `messages` 变化监听新消息。若用户当前在底部，自动平滑滚动到底部；否则 `unreadCount++` 并显示提示条。
3. **点击共情**：消息项上的共情按钮触发 `empathy` 事件，由父组件调用 API 并更新计数。
4. **打开线程**：点击消息的回复区域触发 `open-thread` 事件，父组件打开 `ThreadPanel`。
5. **自动滚动**：首次加载和切换频道时自动滚动到底部。

---

## 2. MessageInput

**组件名：** `MessageInput`
**路径：** `frontend/components/chat/MessageInput.vue`

消息输入组件是用户发布消息的入口，集成了多行文本输入框、表情快捷面板、情绪标签选择器入口、匿名发送开关以及附件上传按钮。该组件支持 Enter 发送、Shift+Enter 换行的快捷键操作，输入内容实时双向绑定。发送前可选择当前心情情绪标签（如"疲惫""愤怒""崩溃"等），情绪标签会随消息一起发送并在消息流中展示。组件还提供输入字数统计和超过限制时的提示。匿名开关允许用户隐藏身份发送，保护打工人隐私。组件通过 `emit('send', payload)` 将消息内容、情绪标签、匿名状态传递给父组件处理。

### Props 定义

| 名称 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `channelId` | `number` | — (必填) | 目标频道 ID |
| `disabled` | `boolean` | `false` | 是否禁用输入（如未加入频道时） |
| `placeholder` | `string` | `'说点什么吧，牛马...'` | 输入框占位文本 |
| `maxLength` | `number` | `2000` | 最大输入字符数 |
| `allowAnonymous` | `boolean` | `true` | 是否允许匿名发送 |

### Emits 定义

| 事件名 | 载荷 | 说明 |
|--------|------|------|
| `send` | `{ content: string; emotion: string \| null; isAnonymous: boolean; attachments: File[] }` | 发送消息时触发 |
| `typing` | `{ isTyping: boolean }` | 用户开始/停止输入时触发（用于 WebSocket typing 事件） |
| `toggle-anonymous` | `{ isAnonymous: boolean }` | 匿名开关切换时触发 |

### Slots

| 名称 | 说明 |
|------|------|
| `toolbar` | 自定义工具栏区域（在输入框上方） |
| `send-button` | 自定义发送按钮 |

### 内部状态

| 状态 | 类型 | 说明 |
|------|------|------|
| `content` | `string` | 输入框文本内容 |
| `selectedEmotion` | `string \| null` | 当前选中的情绪标签 |
| `isAnonymous` | `boolean` | 是否匿名发送 |
| `attachments` | `File[]` | 待上传的附件文件列表 |
| `showEmotionPicker` | `boolean` | 是否显示情绪选择器弹层 |
| `isTyping` | `boolean` | 当前是否正在输入（防抖控制） |

### UI 结构概要

```html
<div class="message-input-container">
  <!-- 工具栏 -->
  <div class="toolbar">
    <slot name="toolbar">
      <button @click="showEmotionPicker = !showEmotionPicker" title="选择情绪">😊 情绪</button>
      <button @click="triggerFileUpload" title="上传附件">📎</button>
      <label v-if="allowAnonymous" class="anonymous-toggle">
        <input type="checkbox" v-model="isAnonymous" /> 匿名
      </label>
    </slot>
  </div>

  <!-- 情绪选择器弹层 -->
  <EmotionPicker
    v-if="showEmotionPicker"
    :selected="selectedEmotion"
    @select="onEmotionSelect"
    @close="showEmotionPicker = false"
  />

  <!-- 输入区域 -->
  <div class="input-area">
    <textarea
      v-model="content"
      :placeholder="placeholder"
      :maxlength="maxLength"
      :disabled="disabled"
      @keydown.enter="handleEnter"
      @input="handleInput"
      rows="1"
    />
    <span class="char-count">{{ content.length }}/{{ maxLength }}</span>
    <slot name="send-button">
      <button :disabled="!canSend" @click="handleSend">发送</button>
    </slot>
  </div>

  <!-- 附件预览 -->
  <div v-if="attachments.length > 0" class="attachment-preview">
    <div v-for="(file, i) in attachments" :key="i" class="attachment-item">
      <img v-if="isImage(file)" :src="previewUrl(file)" />
      <span v-else>{{ file.name }}</span>
      <button @click="removeAttachment(i)">✕</button>
    </div>
  </div>
</div>
```

### 交互行为说明

1. **快捷键发送**：Enter 键直接发送消息；Shift+Enter 插入换行。发送后清空 `content`、`selectedEmotion`、`attachments`，但保留 `isAnonymous` 状态。
2. **输入防抖 typing**：用户开始输入时触发 `typing({ isTyping: true })`，停止输入 2 秒后触发 `typing({ isTyping: false })`。
3. **情绪选择**：点击"情绪"按钮弹出 `EmotionPicker`，选择后显示已选标签，可再次点击取消。
4. **匿名切换**：勾选匿名开关后发送的消息将隐藏用户身份，触发 `toggle-anonymous` 事件。
5. **附件上传**：点击附件按钮触发文件选择，图片类型自动生成缩略图预览，支持删除已选附件。
6. **字数限制**：超过 `maxLength` 时输入框阻止继续输入并显示红色提示。

---

## 3. ThreadPanel

**组件名：** `ThreadPanel`
**路径：** `frontend/components/chat/ThreadPanel.vue`

线程面板组件用于展示某条主消息下的所有回复，类似 Slack 的线程对话功能。面板从右侧滑入，顶部显示被回复的原始消息（含发送者、内容、情绪标签），下方列出所有线程回复，底部提供回复输入框。线程回复支持共情按钮，但不支持再嵌套线程（仅一级）。该组件接收父消息 ID 作为 Prop，通过 API 拉取线程回复列表，并监听 WebSocket 的 `thread_reply` 事件实时追加新回复。面板可通过关闭按钮或点击遮罩层关闭，关闭时保留滚动位置以便再次打开时恢复。

### Props 定义

| 名称 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `parentMessage` | `Message \| null` | `null` | 被回复的父消息对象，为 null 时面板不显示 |
| `replies` | `Message[]` | `[]` | 线程回复列表 |
| `loading` | `boolean` | `false` | 是否正在加载回复 |
| `currentUserId` | `number \| null` | `null` | 当前用户 ID，用于判断自己的回复 |
| `visible` | `boolean` | `false` | 面板是否可见 |

### Emits 定义

| 事件名 | 载荷 | 说明 |
|--------|------|------|
| `close` | `void` | 关闭面板时触发 |
| `send-reply` | `{ parentId: number; content: string; isAnonymous: boolean }` | 发送线程回复时触发 |
| `empathy` | `{ messageId: number }` | 点击回复的共情按钮时触发 |

### Slots

| 名称 | 说明 |
|------|------|
| `header` | 自定义面板头部（标题区域） |
| `reply-item` | 自定义单条回复的渲染（作用域插槽，提供 `reply` 和 `index`） |

### 内部状态

| 状态 | 类型 | 说明 |
|------|------|------|
| `replyContent` | `string` | 回复输入框内容 |
| `isAnonymous` | `boolean` | 回复是否匿名 |
| `scrollRef` | `HTMLElement \| null` | 回复列表滚动容器引用 |

### UI 结构概要

```html
<!-- 遮罩层 -->
<div v-if="visible" class="thread-overlay" @click="$emit('close')"></div>

<!-- 面板主体 -->
<transition name="slide-right">
  <div v-if="visible" class="thread-panel">
    <!-- 头部 -->
    <div class="thread-header">
      <slot name="header">
        <h3>线程</h3>
        <button @click="$emit('close')">✕</button>
      </slot>
    </div>

    <!-- 父消息 -->
    <div v-if="parentMessage" class="parent-message">
      <img :src="parentMessage.avatar" class="avatar" />
      <div class="message-body">
        <span class="nickname">{{ parentMessage.isAnonymous ? '匿名牛马' : parentMessage.nickname }}</span>
        <span v-if="parentMessage.emotion" class="emotion-tag">{{ parentMessage.emotion }}</span>
        <p class="content">{{ parentMessage.content }}</p>
      </div>
    </div>

    <!-- 分隔线 -->
    <div class="divider">{{ replies.length }} 条回复</div>

    <!-- 回复列表 -->
    <div class="reply-list" ref="scrollRef">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-for="reply in replies" :key="reply.id" class="reply-item">
        <slot name="reply-item" :reply="reply">
          <img :src="reply.avatar" class="avatar" />
          <div class="reply-body">
            <span class="nickname">{{ reply.isAnonymous ? '匿名牛马' : reply.nickname }}</span>
            <p class="content">{{ reply.content }}</p>
            <div class="reply-footer">
              <button @click="$emit('empathy', { messageId: reply.id })">🤝 {{ reply.empathyCount }}</button>
              <span class="timestamp">{{ formatTime(reply.createdAt) }}</span>
            </div>
          </div>
        </slot>
      </div>
    </div>

    <!-- 回复输入 -->
    <div class="reply-input">
      <textarea
        v-model="replyContent"
        placeholder="回复这条消息..."
        @keydown.enter.prevent="handleSendReply"
      />
      <label class="anonymous-toggle">
        <input type="checkbox" v-model="isAnonymous" /> 匿名
      </label>
      <button :disabled="!replyContent.trim()" @click="handleSendReply">回复</button>
    </div>
  </div>
</transition>
```

### 交互行为说明

1. **滑入动画**：面板从右侧滑入，使用 Vue `<transition>` 组件实现 `slide-right` 过渡效果。
2. **发送回复**：Enter 键或点击按钮发送，触发 `send-reply` 事件，携带 `parentId`、`content`、`isAnonymous`。发送后清空输入框。
3. **实时追加**：父组件监听 WebSocket `thread_reply` 事件后更新 `replies` Prop，面板自动滚动到底部显示新回复。
4. **关闭面板**：点击关闭按钮或遮罩层触发 `close` 事件，父组件将 `visible` 设为 `false`。
5. **共情**：每条回复上的共情按钮触发 `empathy` 事件，与主消息列表逻辑一致。

---

## 4. EmotionPicker

**组件名：** `EmotionPicker`
**路径：** `frontend/components/emotion/EmotionPicker.vue`

情绪选择器组件是一个弹出式面板，供用户在发送消息时选择一个情绪标签。Alike 的核心特色之一就是情绪标签系统，预设了"😮‍💨疲惫、😡愤怒、😢委屈、🤯崩溃、😴麻木、🔥想润、💪加油、😊开心"等情绪选项。组件以网格布局展示所有可选情绪，点击后高亮选中项并立即通知父组件。已选情绪可再次点击取消选择。组件支持键盘导航（方向键切换、Enter 确认、Escape 关闭）。弹层带有指向触发按钮的小箭头，位置自动计算避免溢出视口。

### Props 定义

| 名称 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `selected` | `string \| null` | `null` | 当前已选中的情绪标识（如 `'tired'`） |
| `emotions` | `EmotionOption[]` | 预设列表 | 可选情绪列表，每项含 `key/label/emoji` |
| `multiple` | `boolean` | `false` | 是否允许多选（默认单选） |

### Emits 定义

| 事件名 | 载荷 | 说明 |
|--------|------|------|
| `select` | `{ emotion: string \| null }` | 选择或取消选择情绪时触发 |
| `close` | `void` | 关闭选择器面板时触发 |

### Slots

| 名称 | 说明 |
|------|------|
| `trigger` | 自定义触发按钮（若组件自带触发器） |
| `emotion-item` | 自定义单个情绪选项的渲染（作用域插槽，提供 `emotion` 和 `isSelected`） |

### 内部状态

| 状态 | 类型 | 说明 |
|------|------|------|
| `activeIndex` | `number` | 键盘导航时当前聚焦的情绪索引 |
| `panelPosition` | `{ top: number; left: number }` | 面板弹出位置坐标 |

### UI 结构概要

```html
<div class="emotion-picker-wrapper">
  <!-- 触发按钮 -->
  <slot name="trigger">
    <button @click="togglePanel" :class="{ active: selected }">
      {{ selected ? getEmotion(selected).emoji + ' ' + getEmotion(selected).label : '😊 情绪' }}
    </button>
  </slot>

  <!-- 弹层面板 -->
  <transition name="fade">
    <div v-if="isPanelOpen" class="emotion-panel" :style="panelPosition" @keydown="handleKeydown">
      <div class="emotion-grid">
        <slot
          v-for="(emotion, i) in emotions"
          :key="emotion.key"
          name="emotion-item"
          :emotion="emotion"
          :isSelected="isSelected(emotion.key)"
        >
          <button
            class="emotion-option"
            :class="{ selected: isSelected(emotion.key), active: i === activeIndex }"
            @click="toggleEmotion(emotion.key)"
          >
            <span class="emoji">{{ emotion.emoji }}</span>
            <span class="label">{{ emotion.label }}</span>
          </button>
        </slot>
      </div>

      <!-- 清除选择 -->
      <button v-if="selected" class="clear-btn" @click="clearSelection">清除选择</button>
    </div>
  </transition>
</div>
```

### 交互行为说明

1. **弹出/收起**：点击触发按钮切换面板显示状态。面板弹出时自动计算位置，确保不溢出视口边缘。
2. **选择情绪**：点击某个情绪选项，若为单选模式则高亮该项并触发 `select` 事件；若该情绪已选中则取消选择（`select` 载荷为 `null`）。
3. **键盘导航**：方向键左右上下切换 `activeIndex`，Enter 确认选择，Escape 关闭面板并触发 `close`。
4. **点击外部关闭**：监听 document click 事件，点击面板外部时关闭面板。
5. **多选模式**：`multiple` 为 `true` 时，选择不互斥，`select` 事件返回已选数组。

---

## 5. EmotionBoard

**组件名：** `EmotionBoard`
**路径：** `frontend/components/emotion/EmotionBoard.vue`

情绪看板组件以可视化方式展示当前频道内所有成员的实时情绪分布，是 Alike 的特色功能之一。看板以热力图/环形图/条形图形式展示各情绪标签的占比和数量，让用户一眼看出"今天大家都在崩溃还是都在摸鱼"。数据通过 REST API 初始加载，随后通过 WebSocket 的 `emotion_update` 事件实时更新。看板支持按时间范围筛选（今日/本周/全部），鼠标悬停各项显示详细数据（情绪标签、人数、占比）。组件采用轻量级 CSS 动画实现数据变化的过渡效果，无需引入额外图表库。

### Props 定义

| 名称 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `channelId` | `number` | — (必填) | 频道 ID |
| `boardData` | `EmotionBoardData` | — | 情绪看板数据，含 `total/emotions[{key,label,emoji,count,percentage}]` |
| `timeRange` | `'today' \| 'week' \| 'all'` | `'today'` | 时间范围筛选 |
| `loading` | `boolean` | `false` | 是否正在加载 |
| `compact` | `boolean` | `false` | 是否紧凑模式（侧边栏展示时使用） |

### Emits 定义

| 事件名 | 载荷 | 说明 |
|--------|------|------|
| `change-range` | `{ range: 'today' \| 'week' \| 'all' }` | 切换时间范围时触发 |
| `refresh` | `void` | 手动刷新时触发 |

### Slots

| 名称 | 说明 |
|------|------|
| `title` | 自定义看板标题区域 |
| `chart` | 自定义图表渲染（作用域插槽，提供 `boardData`） |
| `legend` | 自定义图例区域 |

### 内部状态

| 状态 | 类型 | 说明 |
|------|------|------|
| `activeRange` | `'today' \| 'week' \| 'all'` | 当前选中的时间范围 |
| `hoveredEmotion` | `string \| null` | 鼠标悬停的情绪 key |
| `showTooltip` | `boolean` | 是否显示 tooltip |

### UI 结构概要

```html
<div class="emotion-board" :class="{ compact }">
  <!-- 标题 -->
  <div class="board-header">
    <slot name="title">
      <h3>📊 情绪看板</h3>
    </slot>
    <!-- 时间范围切换 -->
    <div class="range-tabs">
      <button
        v-for="r in ['today', 'week', 'all']"
        :key="r"
        :class="{ active: activeRange === r }"
        @click="changeRange(r)"
      >
        {{ rangeLabel(r) }}
      </button>
    </div>
  </div>

  <!-- 加载状态 -->
  <div v-if="loading" class="loading">加载中...</div>

  <!-- 空状态 -->
  <div v-else-if="boardData.total === 0" class="empty">
    还没有情绪数据，快来第一个发表心情吧
  </div>

  <!-- 图表区域 -->
  <slot name="chart" :boardData="boardData">
    <!-- 条形图 -->
    <div class="bar-chart">
      <div
        v-for="emotion in boardData.emotions"
        :key="emotion.key"
        class="bar-item"
        @mouseenter="hoveredEmotion = emotion.key"
        @mouseleave="hoveredEmotion = null"
      >
        <span class="emoji">{{ emotion.emoji }}</span>
        <div class="bar-track">
          <div class="bar-fill" :style="{ width: emotion.percentage + '%' }"></div>
        </div>
        <span class="count">{{ emotion.count }} 人</span>
        <!-- Tooltip -->
        <div v-if="hoveredEmotion === emotion.key" class="tooltip">
          {{ emotion.label }}: {{ emotion.count }} 人 ({{ emotion.percentage }}%)
        </div>
      </div>
    </div>
  </slot>

  <!-- 总计 -->
  <div class="board-footer">
    <span>共 {{ boardData.total }} 位牛马表达了情绪</span>
    <button @click="$emit('refresh')">🔄 刷新</button>
  </div>
</div>
```

### 交互行为说明：

1. **时间范围切换**：点击 today/week/all 标签触发 `change-range` 事件，父组件重新拉取对应范围的数据。
2. **实时更新**：父组件通过 WebSocket `emotion_update` 事件接收新数据后更新 `boardData` Prop，条形图宽度变化使用 CSS transition 平滑过渡。
3. **悬停详情**：鼠标悬停某条情绪时显示 tooltip，展示情绪名称、人数和百分比。
4. **手动刷新**：点击刷新按钮触发 `refresh` 事件，父组件通过 REST API 重新拉取完整数据。
5. **紧凑模式**：`compact` 为 `true` 时缩小间距、隐藏标题，适用于频道页侧边栏。

---

## 6. EmpathyButton

**组件名：** `EmpathyButton`
**路径：** `frontend/components/empathy/EmpathyButton.vue`

共情按钮（抱团取暖）是 Alike 最核心的交互组件之一，不同于普通的"点赞"，它传达的是"我懂你的辛苦"的情感共鸣。按钮展示当前消息的共情次数，点击后触发共情动画（心形/握手图标弹跳 + 粒子扩散效果），并实时更新计数。已共情状态下按钮高亮，再次点击可取消共情。按钮状态通过 WebSocket 的 `empathy` 事件实时同步——当其他用户共情了同一条消息时，计数会自动增加并伴随轻微脉冲动画。组件内置防抖处理，避免快速连续点击导致重复请求。

### Props 定义

| 名称 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `messageId` | `number` | — (必填) | 关联的消息 ID |
| `count` | `number` | `0` | 当前共情次数 |
| `isEmpathized` | `boolean` | `false` | 当前用户是否已共情 |
| `disabled` | `boolean` | `false` | 是否禁用（未登录时） |
| `size` | `'sm' \| 'md' \| 'lg'` | `'md'` | 按钮尺寸 |
| `showLabel` | `boolean` | `true` | 是否显示文字标签"我懂你" |

### Emits 定义

| 事件名 | 载荷 | 说明 |
|--------|------|------|
| `empathy` | `{ messageId: number; action: 'add' \| 'remove' }` | 点击共情/取消共情时触发 |
| `animation-end` | `void` | 共情动画播放结束时触发 |

### Slots

| 名称 | 说明 |
|------|------|
| `icon` | 自定义图标区域 |
| `label` | 自定义文字标签 |

### 内部状态

| 状态 | 类型 | 说明 |
|------|------|------|
| `isAnimating` | `boolean` | 是否正在播放共情动画 |
| `particles` | `Particle[]` | 粒子动画数据数组 |
| `localCount` | `number` | 本地缓存的计数（乐观更新） |
| `localEmpathized` | `boolean` | 本地缓存的共情状态（乐观更新） |

### UI 结构概要

```html
<div class="empathy-button-wrapper">
  <button
    class="empathy-btn"
    :class="[size, { active: localEmpathized, animating: isAnimating, disabled }]"
    :disabled="disabled"
    @click="handleClick"
  >
    <slot name="icon">
      <span class="icon">{{ localEmpathized ? '🤝' : '🫂' }}</span>
    </slot>
    <slot name="label" v-if="showLabel">
      <span class="label">{{ localEmpathized ? '已懂你' : '我懂你' }}</span>
    </slot>
    <span class="count">{{ localCount }}</span>
  </button>

  <!-- 粒子动画 -->
  <transition-group name="particle" tag="div" class="particle-container">
    <div
      v-for="particle in particles"
      :key="particle.id"
      class="particle"
      :style="particle.style"
    >
      {{ particle.emoji }}
    </div>
  </transition-group>

  <!-- 脉冲效果（其他用户共情时） -->
  <div v-if="pulse" class="pulse-ring"></div>
</div>
```

### 交互行为说明

1. **点击共情**：点击按钮执行乐观更新——立即切换 `localEmpathized` 状态并调整 `localCount`，同时触发 `empathy` 事件。若 API 返回失败则回滚状态。
2. **共情动画**：点击后播放图标弹跳动画，同时生成 6-8 个心形/握手 emoji 粒子向外扩散，动画结束后触发 `animation-end`。
3. **防抖处理**：500ms 内重复点击只触发一次 API 请求，避免重复共情。
4. **实时同步**：父组件监听 WebSocket `empathy` 事件后更新 `count` Prop。若变化非当前用户操作，则播放轻微脉冲动画提示"有人也懂你了"。
5. **禁用状态**：未登录时 `disabled` 为 `true`，按钮置灰，点击跳转登录页。
6. **尺寸适配**：`sm` 用于线程回复，`md` 用于消息列表，`lg` 用于日记详情页。

---

## 7. ChannelSidebar

**组件名：** `ChannelSidebar`
**路径：** `frontend/components/layout/ChannelSidebar.vue`

频道侧边栏是应用左侧的全局导航组件，展示所有可用频道并按分类（行业、岗位、主题、自建）分组展示。用户可以浏览、搜索和切换频道，已加入的频道标记高亮并置顶显示。侧边栏顶部提供频道搜索框，支持按名称模糊匹配过滤。每个频道项显示频道图标、名称和未读消息计数红点。侧边栏支持折叠/展开，折叠时仅显示图标。对于频道管理员，频道项上显示管理入口。组件还展示"创建频道"按钮，点击后弹出创建频道对话框。

### Props 定义

| 名称 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `channels` | `Channel[]` | `[]` | 全部频道列表，每项含 `id/name/slug/description/category/icon/memberCount/unreadCount/isJoined` |
| `activeChannelId` | `number \| null` | `null` | 当前选中的频道 ID |
| `collapsed` | `boolean` | `false` | 侧边栏是否折叠（仅显示图标） |
| `categories` | `string[]` | `['industry', 'job', 'topic', 'custom']` | 频道分类顺序 |

### Emits 定义

| 事件名 | 载荷 | 说明 |
|--------|------|------|
| `select-channel` | `{ channelId: number }` | 点击频道项时触发 |
| `join-channel` | `{ channelId: number }` | 加入频道时触发 |
| `create-channel` | `void` | 点击"创建频道"按钮时触发 |
| `toggle-collapse` | `{ collapsed: boolean }` | 折叠/展开侧边栏时触发 |
| `search` | `{ keyword: string }` | 搜索框输入时触发 |

### Slots

| 名称 | 说明 |
|------|------|
| `header` | 自定义侧边栏头部区域 |
| `channel-item` | 自定义频道项渲染（作用域插槽，提供 `channel`） |
| `footer` | 自定义底部区域 |

### 内部状态

| 状态 | 类型 | 说明 |
|------|------|------|
| `searchKeyword` | `string` | 搜索关键词 |
| `expandedCategories` | `Set<string>` | 展开的分类集合（可折叠分类组） |
| `localCollapsed` | `boolean` | 本地折叠状态（用于动画过渡） |

### UI 结构概要

```html
<div class="channel-sidebar" :class="{ collapsed: localCollapsed }">
  <!-- 头部 -->
  <div class="sidebar-header">
    <slot name="header">
      <button class="collapse-btn" @click="toggleCollapse">
        <span v-if="localCollapsed">☰</span>
        <span v-else>◀</span>
      </button>
      <h2 v-if="!localCollapsed">频道</h2>
    </slot>
  </div>

  <!-- 搜索框 -->
  <div v-if="!localCollapsed" class="search-box">
    <input
      v-model="searchKeyword"
      placeholder="搜索频道..."
      @input="$emit('search', { keyword: searchKeyword })"
    />
  </div>

  <!-- 频道列表（按分类分组） -->
  <div class="channel-list">
    <div v-for="category in filteredCategories" :key="category" class="channel-group">
      <!-- 分类标题 -->
      <button class="category-header" @click="toggleCategory(category)">
        <span class="arrow">{{ expandedCategories.has(category) ? '▼' : '▶' }}</span>
        <span v-if="!localCollapsed">{{ categoryLabel(category) }}</span>
      </button>

      <!-- 分类下的频道 -->
      <div v-if="expandedCategories.has(category)" class="channel-items">
        <div
          v-for="channel in getChannelsByCategory(category)"
          :key="channel.id"
          class="channel-item"
          :class="{ active: channel.id === activeChannelId, joined: channel.isJoined }"
          @click="$emit('select-channel', { channelId: channel.id })"
        >
          <slot name="channel-item" :channel="channel">
            <span class="channel-icon">{{ channel.icon }}</span>
            <span v-if="!localCollapsed" class="channel-name">{{ channel.name }}</span>
            <span v-if="channel.unreadCount > 0 && !localCollapsed" class="unread-badge">
              {{ channel.unreadCount }}
            </span>
          </slot>
        </div>
      </div>
    </div>
  </div>

  <!-- 底部 -->
  <div class="sidebar-footer">
    <slot name="footer">
      <button class="create-channel-btn" @click="$emit('create-channel')">
        ＋ 创建频道
      </button>
    </slot>
  </div>
</div>
```

### 交互行为说明

1. **选择频道**：点击频道项触发 `select-channel` 事件，父组件切换路由到对应频道页。已选频道高亮显示。
2. **搜索过滤**：在搜索框输入关键词后，频道列表实时过滤匹配的频道，不匹配的分类组自动隐藏。
3. **分类折叠**：点击分类标题可展开/收起该分类下的频道列表，状态保存在 `expandedCategories` 中。
4. **侧边栏折叠**：点击折叠按钮将侧边栏收窄为仅图标模式，节省空间，鼠标悬停图标时显示频道名称 tooltip。
5. **未读计数**：频道有未读消息时显示红色数字徽标，点击进入频道后清除。
6. **创建频道**：点击底部"创建频道"按钮触发 `create-channel` 事件，父组件弹出创建频道对话框。

---

## 8. TopNav

**组件名：** `TopNav`
**路径：** `frontend/components/layout/TopNav.vue`

顶部导航栏是应用的全局导航组件，固定在页面顶部。左侧显示 Alike 品牌Logo和标语，中间是主导航链接（首页、频道、日记广场、排行榜），右侧是搜索框、通知图标（带未读红点）、用户头像下拉菜单。通知图标点击后展开通知下拉面板，显示最近的@提及、共情、回复通知。用户头像下拉菜单包含"个人主页""我的日记""设置""退出登录"等入口。导航栏在移动端自适应为汉堡菜单布局。组件还展示当前用户的牛马等级徽章，增强归属感。

### Props 定义

| 名称 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `user` | `UserInfo \| null` | `null` | 当前登录用户信息，含 `id/nickname/avatarUrl/level`；为 null 表示未登录 |
| `notifications` | `Notification[]` | `[]` | 通知列表，每项含 `id/type/content/isRead/createdAt/refId` |
| `unreadCount` | `number` | `0` | 未读通知数量 |
| `activeRoute` | `string` | `'/'` | 当前路由路径，用于高亮导航项 |

### Emits 定义

| 事件名 | 载荷 | 说明 |
|--------|------|------|
| `navigate` | `{ path: string }` | 点击导航链接时触发 |
| `search` | `{ keyword: string }` | 搜索框提交时触发 |
| `mark-notification-read` | `{ notificationId: number }` | 标记单条通知已读 |
| `mark-all-read` | `void` | 标记全部通知已读 |
| `logout` | `void` | 退出登录时触发 |

### Slots

| 名称 | 说明 |
|------|------|
| `logo` | 自定义 Logo 区域 |
| `nav-items` | 自定义导航链接区域 |
| `actions` | 自定义右侧操作区域（通知、用户菜单之前） |
| `user-menu` | 自定义用户下拉菜单内容 |

### 内部状态

| 状态 | 类型 | 说明 |
|------|------|------|
| `searchKeyword` | `string` | 搜索框关键词 |
| `showNotifications` | `boolean` | 是否展开通知下拉面板 |
| `showUserMenu` | `boolean` | 是否展开用户下拉菜单 |
| `showMobileMenu` | `boolean` | 移动端汉堡菜单是否展开 |

### UI 结构概要

```html
<nav class="top-nav">
  <!-- 左侧：Logo -->
  <div class="nav-left">
    <slot name="logo">
      <NuxtLink to="/" class="logo">
        <span class="logo-icon">🐮</span>
        <span class="logo-text">Alike</span>
        <span class="tagline">总有人懂你的辛苦</span>
      </NuxtLink>
    </slot>
  </div>

  <!-- 中间：导航链接（桌面端） -->
  <div class="nav-center desktop-only">
    <slot name="nav-items">
      <NuxtLink
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        class="nav-link"
        :class="{ active: activeRoute === item.path }"
      >
        {{ item.label }}
      </NuxtLink>
    </slot>
  </div>

  <!-- 右侧：搜索 + 通知 + 用户 -->
  <div class="nav-right">
    <slot name="actions">
      <!-- 搜索框 -->
      <div class="search-box desktop-only">
        <input
          v-model="searchKeyword"
          placeholder="搜索消息、日记..."
          @keydown.enter="$emit('search', { keyword: searchKeyword })"
        />
      </div>
    </slot>

    <!-- 通知图标 -->
    <div class="notification-wrapper">
      <button class="icon-btn" @click="showNotifications = !showNotifications">
        🔔
        <span v-if="unreadCount > 0" class="badge">{{ unreadCount }}</span>
      </button>
      <!-- 通知下拉面板 -->
      <transition name="dropdown">
        <div v-if="showNotifications" class="notification-panel">
          <div class="panel-header">
            <span>通知</span>
            <button v-if="unreadCount > 0" @click="$emit('mark-all-read')">全部已读</button>
          </div>
          <div class="notification-list">
            <div
              v-for="notif in notifications"
              :key="notif.id"
              class="notification-item"
              :class="{ unread: !notif.isRead }"
              @click="$emit('mark-notification-read', { notificationId: notif.id })"
            >
              <span class="notif-icon">{{ notifIcon(notif.type) }}</span>
              <span class="notif-content">{{ notif.content }}</span>
              <span class="notif-time">{{ formatTime(notif.createdAt) }}</span>
            </div>
          </div>
        </div>
      </transition>
    </div>

    <!-- 用户菜单 -->
    <div v-if="user" class="user-wrapper">
      <button class="user-btn" @click="showUserMenu = !showUserMenu">
        <img :src="user.avatarUrl" class="avatar" />
        <span class="level-badge">Lv.{{ user.level }}</span>
      </button>
      <transition name="dropdown">
        <div v-if="showUserMenu" class="user-menu">
          <slot name="user-menu">
            <NuxtLink :to="`/profile/${user.id}`">个人主页</NuxtLink>
            <NuxtLink :to="`/diary?user=${user.id}`">我的日记</NuxtLink>
            <NuxtLink to="/settings">设置</NuxtLink>
            <button @click="$emit('logout')">退出登录</button>
          </slot>
        </div>
      </transition>
    </div>
    <!-- 未登录 -->
    <div v-else class="auth-buttons">
      <NuxtLink to="/login">登录</NuxtLink>
      <NuxtLink to="/register">注册</NuxtLink>
    </div>

    <!-- 移动端汉堡菜单 -->
    <button class="hamburger mobile-only" @click="showMobileMenu = !showMobileMenu">
      ☰
    </button>
  </div>

  <!-- 移动端菜单展开 -->
  <transition name="slide-down">
    <div v-if="showMobileMenu" class="mobile-menu">
      <NuxtLink
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        @click="showMobileMenu = false"
      >
        {{ item.label }}
      </NuxtLink>
    </div>
  </transition>
</nav>
```

### 交互行为说明

1. **导航跳转**：点击导航链接触发 `navigate` 事件并路由跳转，当前页对应的导航项高亮显示。
2. **搜索**：在搜索框输入关键词后按 Enter 触发 `search` 事件，父组件跳转到搜索结果页。
3. **通知面板**：点击通知图标展开下拉面板，展示最近通知。点击单条通知触发 `mark-notification-read` 并跳转到关联页面。点击"全部已读"触发 `mark-all-read`。
4. **用户菜单**：点击头像展开下拉菜单，包含个人主页、我的日记、设置、退出登录等入口。点击外部区域自动收起菜单。
5. **移动端适配**：窗口宽度小于 768px 时隐藏中间导航和搜索框，显示汉堡菜单按钮。点击汉堡菜单展开垂直导航列表。
6. **牛马等级徽章**：头像旁显示当前用户等级徽章，等级越高徽章颜色越醒目，增强用户的归属感和成就感。
7. **未登录状态**：`user` 为 null 时右侧显示"登录/注册"按钮，点击跳转对应页面。

---

## 附录：通用类型定义

```typescript
// 消息类型
interface Message {
  id: number
  userId: number | null
  nickname: string
  avatar: string
  content: string
  emotion: string | null
  isAnonymous: boolean
  empathyCount: number
  threadCount: number
  parentId: number | null
  createdAt: string
}

// 频道类型
interface Channel {
  id: number
  name: string
  slug: string
  description: string
  category: 'industry' | 'job' | 'topic' | 'custom'
  icon: string
  memberCount: number
  unreadCount: number
  isJoined: boolean
}

// 情绪选项
interface EmotionOption {
  key: string
  label: string
  emoji: string
}

// 情绪看板数据
interface EmotionBoardData {
  total: number
  emotions: Array<{
    key: string
    label: string
    emoji: string
    count: number
    percentage: number
  }>
}

// 用户信息
interface UserInfo {
  id: number
  nickname: string
  avatarUrl: string
  level: number
  industry: string
  jobTitle: string
}

// 通知
interface Notification {
  id: number
  type: 'mention' | 'empathy' | 'reply' | 'system'
  content: string
  isRead: boolean
  createdAt: string
  refId: number
}
```

---

> **文档版本：** v1.0
>
> **最后更新：** 2026-07-10
>
> **关联文档：** [架构设计实施计划](../plans/architecture-design.md) | [CLAUDE.md](../../CLAUDE.md)
