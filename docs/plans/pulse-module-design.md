# 最近发生（Pulse / 世界脉搏）— 首发聚焦版设计文档

> **状态：设计草案 + 视觉原型已收敛，待技术评审** ｜ 后端模块代号 `pulse` ｜ 前端一级入口文案「最近发生」
> 与「首页」「日记广场」「排行榜」同级导航（**放在日记广场与排行榜之间**）。**定位：纯资讯瀑布 · 只读展示**（不接入站内互动/评论/共情）。
> 首发只做 **2 个专题**先跑通抓取→展示的最短闭环：**GitHub 每日 Star 榜** + **AI 圈大事**。其余专题进附录路线图。
>
> **视觉原型：** `docs/design/07-pulse.html`（已按线上 http://39.107.58.169/ 的 design tokens 严格对齐 —— 复用 `.glass-card` / `.text-gradient` / `.btn-primary` 通用类，导航结构、`--max-w-content-wide (1040px)`、`xl:block (≥1280px)` 右栏断点、未登录态登录/注册按钮、logo 花瓣 icon、双主题变量均一致）。

---

## 电梯陈述

> **打工人的"世界正在发生什么"**——不是又一个门户资讯页，而是把散落在 GitHub、Hacker News、X 上"我今天想知道但没时间刷"的信号，**按打工人视角**聚合到一个地方。切了平台，但没切自己。

核心比喻：**脉搏（pulse）**。世界的脉搏、行业的脉搏，也是打工人"我还在这个圈子里"的那点心跳感。视觉上对应 `--grad-ai`（蓝紫→青）流光效果，与首页 hero 同频。

**一句话守则：最近发生不是"新闻聚合器"，是"打工人不切平台就能同步世界的一扇窗"——只读、精选、克制，永远不做评论/共情/线程，让它保留"看一眼就走"的干净感。**

> **内容边界守则（与其他板块的硬分野）：**
> - **VS 频道/日记广场**：那两处是站内 UGC；这里 100% 是**外部信号的镜像**，没有一条内容起源于站内用户，规避 UGC 审核成本。
> - **VS 排行榜**：排行榜排的是**站内牛马**（人）；这里排的是**外部事件**（事）。人榜关心谁在扛，事榜关心世界在动。
> - **VS 成长充电**：那里是"我在动"（动作/输出）；这里是"世界在动"（信号/事件）。二者一内一外，互补而不重叠。
>
> 手段是**"只读 + 结构化抓取"而非人工编辑**——每个专题绑定一个明确的**外部数据源 + 抓取器 + 排序公式**，Alike 只做"搬运 + 排序 + 视觉"，不做主观选题。这既过滤了 UGC 风险、又不用养编辑部（规避人肉成本与偏见争议）。

### 首发跑通一条最短闭环（"抓 → 存 → 排 → 展示" 四步走）

| 环节 | 功能 | 作用 |
|------|------|------|
| 抓取 | 定时任务按专题拉源站（GitHub Trending / HN AI 关键词） | 内容 = 外部信号，零 UGC 冷启动依赖 |
| 存储 | `pulse_items` 表按 `(topic,source,source_id)` 去重 + 打时间戳 | 幂等、可回溯、可回滚 |
| 排序 | 每专题一套 `score` 公式（今日增 star / HN points × 时间衰减）| 用规则替代编辑，规避人肉成本 |
| 展示 | 一级导航「最近发生」→ 专题 Tab → 事件卡片流（只读，点开跳外链）| 一屏看完，不留、不 UGC、不焦虑 |

**最小闭环**：用户点「最近发生」→ 默认 Tab 是「GitHub 每日 Star」→ 看到今日增 star Top20 卡片（含仓库名/描述/语言/今日 +N ⭐/链接）→ 点卡片新开外链去 GitHub → 切 Tab 到「AI 圈大事」看 HN 热帖 → 关掉页面。**全程不登录也能看**。

---

## 一、概述与背景

Alike 已有频道、消息线程、共情、情绪、日记、排行榜、通知、成长充电等以 UGC 为核心的能力。这些能力全部围绕**站内用户**运转。但打工人真实场景里，"我想知道最近发生了什么"这个需求会**把他切走到别的 App**（GitHub / 微博热搜 / 掘金 / X），这既是**流失口**也是**回流机会**——如果 Alike 有一个"顺手看一眼"的最近发生页，用户就多一个不切平台的理由。

设计原则来自四方评审等价推演的核心判断：**不做 UGC 的板块，冷启动反而最简单——内容不依赖用户密度，只依赖抓取器质量。** 因此首发砍掉一切"需要用户参与"的功能（评论/共情/收藏/订阅），只保留"抓 → 展示"的纯读通道，把风险和成本压到最低，等留存数据说话再谈是否加互动。

---

## 二、产品设计

### 2.1 定位与核心价值

**命名：最近发生**（英文 `Pulse`，后端代号 `pulse`）。主视觉隐喻"心跳/脉搏"，对应"世界在动、你也还在圈子里"的感觉，与 `--grad-ai` 渐变（蓝紫→青）契合。

**核心价值：**
1. **不切平台** — 想看 GitHub Trending / HN / X 热帖时不用离开 Alike，用户在站时长自然增加。
2. **打工人视角精选** — 不是全量新闻，只挑"程序员/产品/设计/AI 从业者会关心"的信号源，切掉娱乐八卦噪音。
3. **只读不打扰** — 没有评论、没有共情、没有@，看完关掉不留任何"你欠我一个回复"的心理负担。这是它区别于其他板块的**克制感**。

### 2.2 目标用户与典型场景

- **画像 A｜信息饥渴的技术小林**（26 岁前端，工龄 3 年）：早上刷 GitHub Trending 看今日热门仓库是习惯。以前必须打开 github.com，现在直接进 Alike「最近发生」→「GitHub 每日 Star」，顺便把手机的 GitHub App 卸了。
- **画像 B｜AI 焦虑症老王**（32 岁产品，工龄 8 年）：怕被 AI 淘汰，每天要看模型/公司动态。以前刷 X、订阅新闻信 letter，信息碎且英文多。现在进「AI 圈大事」看中文摘要 + 高赞 HN 帖子的一句话原意，10 分钟同步完当日圈内事。
- **画像 C｜通勤党阿May**（29 岁设计，工龄 5 年）：地铁 40 分钟。以前刷 B 站、微博，看完啥也没记住。有了「最近发生」，切成"看一眼行业信号"的模式，落地时脑子里有 3-5 个可以在午饭和同事聊的话题。

### 2.3 核心功能（首发 = 2 个专题，一条闭环）

> 原则：**只做搬运 + 排序，不做原创和互动**；专题即数据源的可插拔封装，首发 2 个跑通框架，后续加专题只需实现新 collector。

**① GitHub 每日 Star 榜（Topic: `github-trending`）**

- **数据源**：`https://github.com/trending?since=daily` HTML 抓取（GitHub 官方没有 Trending API，Trending 页 HTML 结构相对稳定；备选 `https://ghtrending.jyt.io/` 或自建镜像）。
- **抓取频率**：每 60 分钟一次（GitHub Trending 页本身每天更新一次为主，这里频率主要为容错），存最近 24 小时快照。
- **卡片字段**：`owner/repo` · 一句话描述 · 主语言（带颜色圆点）· **今日 +N ⭐** · 总 star 数 · fork 数 · 外链。
- **排序**：按 `today_stars` 降序，同分按 `total_stars` 降序。
- **展示**：默认 Top 25，右上角小字"截至 HH:MM 更新"。
- **点击行为**：整卡可点，`target="_blank" rel="noopener"` 打开 GitHub 仓库页。
- **打工人视角加料**：语言标签用 GitHub 官方颜色圆点保持熟悉感；仓库名右侧显示 `#GO` / `#PY` / `#TS` 等常见语言 Chip，便于快速筛选。

**② AI 圈大事（Topic: `ai-news`）**

- **数据源**：Hacker News Algolia API `https://hn.algolia.com/api/v1/search?query=<kw>&tags=story&numericFilters=points>50`，关键词轮询 `openai|anthropic|claude|gpt|llm|gemini|mistral|deepseek|ai model|foundation model` + 时间窗口"过去 72 小时"。（免费、稳定、有 CORS）
- **抓取频率**：每 30 分钟一次。
- **卡片字段**：标题 · 来源域名（如 `openai.com` / `arxiv.org` / `techcrunch.com`）· HN points · HN 评论数 · 发布时间（相对时间"3 小时前"）· 外链。
- **排序**：加权分 `score = hn_points × time_decay(hours)`（半衰期 24 小时，见 3.5）。
- **展示**：默认 Top 20。
- **点击行为**：整卡点击 → 打开原文外链；副按钮"HN 讨论 →"→ 打开 HN 帖子页。
- **打工人视角加料**：标题下方一行灰字标"发布方 / 讨论热度"，让"这条重不重要"一眼可判。

**注：**首发不做以下功能，见 2.6 反范围：
- 不做评论/点赞/收藏/共情。
- 不做订阅/推送/@提醒。
- 不做搜索/过滤/自定义 Tab（首发 2 Tab 固定顺序）。
- 不做站内评论区/线程讨论。
- 不做 AI 摘要（首发用原标题+一句话原始描述，规避幻觉与成本）。

### 2.4 核心用户故事（含验收）

- **US-1 一屏看完当日 GitHub Trending**：进入「最近发生」默认落在 GitHub Tab，看到今日 star 增长 Top 25。**验收**：卡片按 `today_stars` 降序；每卡显示 owner/repo、描述、主语言色点、今日 +N ⭐、总 star 数；顶部显示"截至 HH:MM 更新"；点整卡新开外链。
- **US-2 切 Tab 看 AI 圈大事**：Tab 切到「AI 圈大事」，看 HN 上 AI 相关的高分帖子。**验收**：Top 20 按加权分排序；卡片显示标题、来源域名、HN points、评论数、发布时间；点整卡开原文；有副按钮"HN 讨论 →"开 HN 页面。
- **US-3 未登录也能看**：不登录状态直接访问 `/pulse`，看到完整数据。**验收**：`/api/pulse/*` 全部路由都是公开可读（无 authMW）；页面 SSR/SSG 出得来（Nuxt3 mode）。
- **US-4 抓取失败不空白**：源站临时抽风时页面不能空。**验收**：collector 抓取失败时不覆写数据库，用户看到的是上一次成功抓取的快照；顶部小字仍显示"截至 HH:MM 更新"（真实抓取成功时间，不撒谎）。
- **US-5 时间新鲜度可判**：用户能判断当前看到的是"新鲜的"还是"过时的"。**验收**：每卡片有"3 小时前"相对时间；专题头有"截至 HH:MM 更新"绝对时间；超过 6 小时未更新则显示黄色警告标"数据可能不新鲜"。
- **US-6 移动端可读**：手机浏览器打开卡片流不错位。**验收**：卡片在 375px 宽下不横向滚动；语言色点/分数标签不换行错位；WCAG AA 对比度。

### 2.5 成功指标

**北极星：最近发生页周活跃 PV（Pulse WAP）** — 每周至少浏览一次「最近发生」的独立用户数。

**辅助指标：**
1. **专题点击深度**（人均查看专题数） — 验证多专题是否有价值。
2. **外链点击率**（CTR = 外链点击 / 卡片曝光） — 验证卡片质量是否足够勾住用户。
3. **回访周期**（第二次访问的间隔天数）— 验证用户是否把这里当习惯性入口。

### 2.6 反范围（明确不做，守住克制感）

- ❌ **不做站内评论/共情/收藏/订阅** — 加互动就变频道，克制感立即崩塌。等留存验证后再讨论（附录 P1）。
- ❌ **不做 AI 摘要/自动翻译** — 首发用原标题 + 原描述，规避幻觉、翻译费与内容审核问题。
- ❌ **不做人工编辑/推荐位** — 内容 = 数据源 × 抓取器 × 排序公式，一切靠规则，规避主观偏见争议与编辑部成本。
- ❌ **不做用户自定义 Tab / 关键词订阅** — 是 P1 功能，首发不做；避免过度设计。
- ❌ **不做长时间历史归档/时间轴** — 只保留最近 7 天数据，过期归档 → 简单删除。这是"最近发生"不是"档案馆"。
- ❌ **不做推送通知/邮件订阅** — 是 P1，会大幅拉高工程复杂度。
- ❌ **不做非白名单数据源** — 首发只 GitHub Trending + HN，后加专题必须**产品拍板 + 有稳定官方 API 或稳定 HTML 结构**，绝不接色情/赌博/无版权源站。
- ❌ **不做无版权文章正文抓取** — 只抓标题 / 描述 / 元数据（GitHub 描述是仓库自有开源信息、HN 元数据是 HN 公开 API 提供），点开跳原文，不做全文镜像。

---

## 三、技术架构

> 遵循现有约定：`database/sql`+pgx、BIGSERIAL 主键、TIMESTAMPTZ DEFAULT NOW()、软删除 `deleted_at`（本模块用**硬删除**因为是外部数据快照，见 3.6）、冗余计数、游标分页。模块保持 `model.go/repository.go/service.go/handler.go` 四段式；抓取器单独一层 `collector/`。

### 3.1 新增后端模块

| 模块 | 路径 | 职责 | 与现有模块协作 |
|------|------|------|---------------|
| `pulse` | `internal/pulse` | 专题元数据 CRUD + 热点条目查询（只读对外）| 无 UGC 交互，独立成岛；不依赖 auth/emotion/notification |
| `pulse/collector` | `internal/pulse/collector` | 抓取器接口 `Collector` + `GitHubTrending` / `HackerNewsAI` 两个实现 | 依赖 `net/http` + `github.com/PuerkitoBio/goquery`（HTML 解析）|
| `pulse/scheduler` | `internal/pulse/scheduler` | 定时任务调度器：按每专题的 `refresh_interval_min` 周期性触发 collector | 应用启动时挂 goroutine，退出时优雅关闭；单实例部署无需 leader election |

> 首发不引入任何第三方定时任务库（cron），用 `time.Ticker` + goroutine 即可满足 2 个专题；如后续专题增多再评估。

### 3.2 数据库表设计（`migrations/007_pulse.sql`，编号接现有最大 006）

```sql
-- 专题：一个专题 = 一个数据源封装
CREATE TABLE IF NOT EXISTS pulse_topics (
    id                    BIGSERIAL PRIMARY KEY,
    slug                  VARCHAR(32) UNIQUE NOT NULL,   -- 'github-trending', 'ai-news'
    name                  VARCHAR(64) NOT NULL,          -- 'GitHub 每日 Star 榜'
    emoji                 VARCHAR(8),                    -- '🌟' '🤖'
    description           TEXT,                          -- 卡片下方灰字副标题
    collector_kind        VARCHAR(32) NOT NULL,          -- 'github_trending' / 'hackernews_ai'（对应 collector 注册键）
    collector_config      JSONB,                         -- collector 的可调参数（如 HN 关键词、GitHub 时间窗口）
    sort_order            INT DEFAULT 0,                 -- 前端 Tab 顺序
    is_active             BOOLEAN DEFAULT TRUE,          -- 关闭抓取但保留数据
    refresh_interval_min  INT DEFAULT 60,                -- 抓取间隔（分钟）
    last_fetched_at       TIMESTAMPTZ,                   -- 最近一次成功抓取时间（页面顶部"截至 HH:MM 更新"）
    last_error            TEXT,                          -- 最近一次抓取错误信息（成功则清空）
    created_at            TIMESTAMPTZ DEFAULT NOW(),
    updated_at            TIMESTAMPTZ DEFAULT NOW()
);

-- 热点条目：一条 = 一个外部事件（一个 GitHub 仓库今日快照 / 一条 HN 帖子）
CREATE TABLE IF NOT EXISTS pulse_items (
    id             BIGSERIAL PRIMARY KEY,
    topic_id       BIGINT NOT NULL REFERENCES pulse_topics(id) ON DELETE CASCADE,
    source         VARCHAR(32) NOT NULL,               -- 'github' / 'hackernews'（同一 topic 可能有多个 source，暂无需但预留）
    source_id      VARCHAR(255) NOT NULL,              -- 去重键：'owner/repo' / HN story_id
    title          TEXT NOT NULL,                      -- 仓库名 / HN 标题
    summary        TEXT,                               -- 仓库描述 / HN 帖子摘要（HN 有 story_text 时用）
    url            TEXT NOT NULL,                      -- 外链
    author         VARCHAR(128),                       -- 仓库 owner / HN 作者
    score          INT DEFAULT 0,                      -- 展示主分：今日 +N stars / HN points
    score_delta    INT DEFAULT 0,                      -- 增量：相对上次抓取的增量（GitHub 增 star 用得上）
    extra          JSONB,                              -- 源特有字段：{"lang":"Go","total_stars":12345,"forks":234} / {"comments":128,"hn_id":39012345}
    published_at   TIMESTAMPTZ,                        -- 原始发布/首次出现时间
    captured_at    TIMESTAMPTZ DEFAULT NOW(),          -- 本次抓取时间
    UNIQUE(topic_id, source, source_id)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_pulse_items_topic_score
    ON pulse_items(topic_id, score DESC, captured_at DESC);
CREATE INDEX IF NOT EXISTS idx_pulse_items_topic_captured
    ON pulse_items(topic_id, captured_at DESC);

-- 种子专题
INSERT INTO pulse_topics (slug, name, emoji, description, collector_kind, collector_config, sort_order, refresh_interval_min) VALUES
  ('github-trending', 'GitHub 每日 Star 榜', '🌟', '今天全世界程序员都在给谁星星', 'github_trending',
   '{"since":"daily","limit":25}'::jsonb, 1, 60),
  ('ai-news',         'AI 圈大事',         '🤖', 'HN 上关于 AI 的最新高赞讨论',    'hackernews_ai',
   '{"keywords":["openai","anthropic","claude","gpt","llm","gemini","mistral","deepseek"],"min_points":50,"window_hours":72,"limit":20}'::jsonb,
   2, 30)
ON CONFLICT (slug) DO NOTHING;
```

> **数据保留策略**：`pulse_items` 用**硬删除 + 定时清理**，而不是软删除。理由：这里是外部数据镜像，不是 UGC，用户没有"我发的东西被删"的心理预期；保留 7 天窗口足够（"最近发生"语义即最近），旧数据直接 `DELETE WHERE captured_at < NOW() - INTERVAL '7 days'`，定时任务每天跑一次。

### 3.3 抓取器（Collector）设计

**接口（`internal/pulse/collector/collector.go`）：**

```go
package collector

import (
    "context"
    "encoding/json"
    "time"
)

// Item 抓取器返回的标准条目（对应 pulse_items 表的核心字段）
type Item struct {
    Source      string          // 'github' / 'hackernews'
    SourceID    string          // 去重键
    Title       string
    Summary     string
    URL         string
    Author      string
    Score       int             // 主分：today_stars / hn_points
    ScoreDelta  int             // 增量（可选）
    Extra       json.RawMessage // JSONB
    PublishedAt *time.Time
}

// Collector 单个专题的抓取器，无状态、并发安全
type Collector interface {
    Kind() string                                            // 注册键 'github_trending' / 'hackernews_ai'
    Fetch(ctx context.Context, config json.RawMessage) ([]Item, error)
}

// Registry Collector 注册表（应用启动时装配）
var Registry = map[string]Collector{}

func Register(c Collector) { Registry[c.Kind()] = c }
```

**GitHub Trending Collector（`github_trending.go`）：**

- 抓取 `https://github.com/trending?since={config.since}`（`daily` / `weekly` / `monthly`）。
- 用 `goquery` 解析：每个 `.Box-row` 是一个仓库，字段抽取：
  - `a[href^="/"]` → `owner/repo`
  - `p.color-fg-muted` → 描述
  - `span[itemprop="programmingLanguage"]` → 主语言
  - `.d-inline-block.ml-0.mr-3 > a[href$="/stargazers"]` → 总 star 数
  - `.d-inline-block.float-sm-right` → **今日 +N stars**（这就是我们要排的 `score`）
- User-Agent 伪装成 `Mozilla/5.0 ...`（GitHub 会挡明显的爬虫 UA）。
- 失败重试：3 次 + 指数退避（500ms/1s/2s）；全失败则返回 error（scheduler 层不覆盖旧数据）。

**Hacker News AI Collector（`hackernews_ai.go`）：**

- 用 Algolia API：`https://hn.algolia.com/api/v1/search?query={kw}&tags=story&numericFilters=points>{min_points},created_at_i>{unix_ts}`。
- 对 `config.keywords` 数组每个关键词并发查询，Merge 去重（按 `story_id`）后统一排序。
- 直接拿到 JSON，无 HTML 解析，稳定性远高于 GitHub。
- 字段抽取：
  - `title` → Title
  - `url` → URL（HN 帖子可能是纯讨论无外链，此时 URL 用 HN 帖子页 `https://news.ycombinator.com/item?id={id}`）
  - `points` → Score
  - `num_comments` → extra.comments
  - `author` → Author
  - `created_at` → PublishedAt

**注册（`internal/pulse/collector/init.go`）：**

```go
func init() {
    Register(&GitHubTrending{})
    Register(&HackerNewsAI{})
}
```

### 3.4 REST API 契约（`/api/pulse/*`）

统一响应 `{ code, message, data }`；列表用普通 `{list, total}`（首发不需要游标分页——每专题 ≤ 25 条固定）。**全部公开可读**（无 authMW）。

| 方法 | 路径 | 说明 | 返回 |
|------|------|------|------|
| GET | `/api/pulse/topics` | 专题列表（前端渲染 Tab） | `[{id,slug,name,emoji,description,last_fetched_at,sort_order}]`（`is_active=true`）|
| GET | `/api/pulse/topics/:slug/items` | 某专题的条目列表（默认 Top 25，按 score DESC）| `{topic:{...}, list:[...], last_fetched_at, stale:bool}`；`stale=true` 当 `NOW() - last_fetched_at > 6h` |
| GET | `/api/pulse/topics/:slug` | 单专题元数据（可选，前端可从 `/topics` 里筛出） | `{...}` |

**响应示例（`GET /api/pulse/topics/github-trending/items`）：**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "topic": {
      "slug": "github-trending",
      "name": "GitHub 每日 Star 榜",
      "emoji": "🌟",
      "description": "今天全世界程序员都在给谁星星"
    },
    "last_fetched_at": "2026-07-23T09:15:00+08:00",
    "stale": false,
    "list": [
      {
        "id": 1,
        "title": "microsoft/vscode",
        "summary": "Visual Studio Code",
        "url": "https://github.com/microsoft/vscode",
        "author": "microsoft",
        "score": 823,
        "extra": {"lang": "TypeScript", "lang_color": "#3178c6", "total_stars": 158234, "forks": 27943},
        "published_at": "2015-09-03T00:00:00Z",
        "captured_at": "2026-07-23T09:15:00+08:00"
      }
    ]
  }
}
```

### 3.5 排序公式与时间衰减

**GitHub Trending**：直接用 `today_stars` 降序。GitHub 已经做了排序，我们信它。

**HN AI**：加权分公式：

```
score_ranked = hn_points × exp(-ln(2) × age_hours / 24)
```

- `age_hours = (NOW() - published_at).Hours`
- 半衰期 24 小时：新帖 24 小时后热度砍半，48 小时砍到 1/4。
- 这是 HN 自家排序公式的简化版（HN 用的是 Gravity=1.8 的 Reddit 公式），半衰期 24h 对"AI 圈大事"这个场景（多为 24-72 小时的新闻）合适。

**实现层**：`pulse_items.score` 存原始分（`hn_points`），`score_ranked` **不入库**，查询时在 SQL 里现算：

```sql
SELECT *, score * EXP(-0.693 * EXTRACT(EPOCH FROM (NOW() - published_at))/86400) AS ranked
FROM pulse_items WHERE topic_id = $1
ORDER BY ranked DESC LIMIT 20;
```

这样调半衰期不用回填历史数据。

### 3.6 定时任务与失败处理（Scheduler）

**结构**：应用启动时启动一个 supervisor goroutine，为每个 `is_active=true` 的专题起一个 worker goroutine（`time.Ticker` + `ctx.Done()`）。

```go
// internal/pulse/scheduler/scheduler.go
func (s *Scheduler) Start(ctx context.Context) {
    topics := s.repo.ListActiveTopics(ctx)
    for _, t := range topics {
        go s.runTopic(ctx, t)  // 每个专题独立 goroutine
    }
}

func (s *Scheduler) runTopic(ctx context.Context, t Topic) {
    // 启动时先跑一次，别等第一个 tick
    s.fetchOnce(ctx, t)
    tick := time.NewTicker(time.Duration(t.RefreshIntervalMin) * time.Minute)
    defer tick.Stop()
    for {
        select {
        case <-ctx.Done(): return
        case <-tick.C: s.fetchOnce(ctx, t)
        }
    }
}
```

**幂等 upsert**：`pulse_items` 用 `INSERT ... ON CONFLICT (topic_id, source, source_id) DO UPDATE SET title=..., score=..., score_delta=score-EXCLUDED.score, captured_at=NOW()`。

**失败处理（三层）：**

1. **单次抓取失败**（网络错误 / 源站 5xx / 结构变化解析失败）：collector 内部重试 3 次（500ms/1s/2s 退避）。
2. **整轮失败**：collector 三次全败 → scheduler **不覆盖数据库**，只写 `pulse_topics.last_error` 和监控日志；用户看到的仍是上一次成功的快照。
3. **数据过期告警**：`NOW() - last_fetched_at > 6h` 时 API 响应里 `stale=true`，前端头部显示黄色警告"数据可能不新鲜"。

**清理任务**：每天凌晨跑一次 `DELETE FROM pulse_items WHERE captured_at < NOW() - INTERVAL '7 days'`，作为额外的 worker goroutine（周期 24h）。

**关闭**：应用退出时 `ctx.Cancel()`，所有 worker 收到信号后清理 Ticker 并退出，`sync.WaitGroup` 兜底最多等 5 秒。

### 3.7 缓存

- **API 层缓存**：`GET /api/pulse/topics/:slug/items` 加 Redis 缓存 60s TTL，key=`pulse:items:{slug}:v1`，Cache-Control `public, max-age=60`。这个板块是公共只读数据，缓存命中率极高。
- **单实例部署阶段**（当前）：可以直接用进程内 `sync.Map`+过期时间，规避 Redis 依赖。切多实例时切 Redis。

### 3.8 WebSocket

**首发 0 新增 WS 事件。** 这个板块**从设计上就不做实时推送**——保留"看一眼就走"的克制感。

---

## 四、视觉与交互（Aurora 极光风）

> **高保真原型已完成，见 `docs/design/07-pulse.html`。** 已按线上 http://39.107.58.169/ 的实际 DOM 严格对齐 —— 复用 `.glass-card` / `.text-gradient` / `.btn-primary` 通用类、`--max-w-content-wide (1040px)`、`xl:block (≥1280px)` 右栏断点、未登录态"登录/注册"按钮、logo 花瓣 icon、双主题变量均一致。实施时前端直接照着做。

### 4.1 导航与信息架构

`TopNav` / `LeftSidebar` 新增一级入口 **最近发生**，图标 `activity`（心跳线，`<polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>`），**放在「日记广场」与「排行榜」之间**（同为"信号看板"类，但排行榜排的是人、这里排的是事）。

**导航项呼吸绿点**：入口右侧带一个 `pulse-dot`（呼吸动画），暗示"这里正在跳动"，是全站唯一带跳动元素的导航项。

**移动端**：并入侧边抽屉（沿用线上 `<md` 断点隐藏 sidebar 的既有规则）。

### 4.2 页面布局

**桌面（≥1280px，三列 shell 复用线上约定）**：左侧 260px sidebar + 中间 max 1040px 主区 + 右侧 320px aside（`xl:block`）。主区从上到下：**Hero（`glass-card animate-rise-in`，含渐变大标题 + 心电图水印 + 三 chip：🟢正在跳动 / 截至 HH:MM 更新 / 数据来源）** → stale 警告条（仅数据不新鲜时显示）→ 「专题 · 每小时更新」板块标题 → Tab（🌟GitHub 每日 Star / 🤖AI 圈大事）→ 卡片列表（Top 25/20）→ 底部克制感 `.quiet-strip`。

**移动端（<768px）**：隐藏左右栏，卡片全宽；Tab 顶部粘性；卡片可点整块；分数块从右侧移到卡片底部（`flex-wrap`）。

### 4.3 首发前端组件

| 组件 | 路径 | 职责 |
|------|------|------|
| `PulseIndex` | `pages/pulse/index.vue` | 页面主体：调 `/api/pulse/topics` 渲染 Tab + 侧栏专题项，切 Tab 请求 `/items` |
| `PulseHero` | `components/pulse/PulseHero.vue` | 页面头：渐变大标题 + 心电图水印 + 三 chip（live 状态 · 更新时间 · 数据源），复用 `.glass-card animate-rise-in` |
| `PulseTabs` | `components/pulse/PulseTabs.vue` | 专题 Tab 切换（复用现有 tab pill 样式）|
| `PulseFeed` | `components/pulse/PulseFeed.vue` | 卡片列表容器 + 空/错/加载态 + `stale` 警告条 |
| `RepoCard` | `components/pulse/RepoCard.vue` | GitHub 专题卡（`.row.repo`）：左侧橙彩条 · 仓库名 · 语言色点 · 星数 · fork · 今日增星 · 外链 |
| `NewsCard` | `components/pulse/NewsCard.vue` | HN 专题卡（`.row.news`）：左侧 HN 橙彩条 · 标题 · 域名 · points · 评论数 · 相对时间 · 外链 · HN 讨论副按钮 |
| `SideTopicNav` | `components/pulse/SideTopicNav.vue` | 侧栏专题项 + 呼吸绿点 + 候补专题灰态 |
| `QuietStrip` | `components/pulse/QuietStrip.vue` | 底部/右栏的克制感声明卡（渐变淡 AI 色边框）|

**卡片视觉规范**（沿用线上 design token）：

- **容器基类** `.row`：`background=--glass-bg`、`border=--glass-border`、`radius=--radius-lg`、`padding=16px 18px`、`hover:translateY(-2px) + shadow-lg`。
- **左侧彩条** `::before`：`width=3px`，GitHub 用 `--grad-warm`（暖橙），HN 用 `linear-gradient(180deg,#ff8a3d,#ff6600)`（HN 橙）——一眼区分类型。
- **排名色**：#1 `--gold` / #2 `#cbd5e1` / #3 `--warm`，与排行榜领奖台一致。
- **语言色点**：GitHub 官方色，见 3.3 常量表；前端维护在 `constants/langColors.ts`。
- **分数**：GitHub 今日增星用 `--warm` + 星星图标；HN points 用 HN 橙 `#ff6600` + 火苗图标。**不复用 `--empathy` 绿**（避免语义混淆——共情是站内 UGC 语义，不能借来标外部数据；原型已改用 HN 橙）。
- **外链图标**：右下角 `lucide external-link`，`aria-label="打开外部链接"`；hover 时字体色变 `--ai-2`。

### 4.4 状态与更新时间

- **"正在跳动"绿丸**：`meta-chip.live` 常驻显示，只要页面能出内容就跳；不代表"实时"，代表"抓取器还活着"。
- **更新时间**：`截至 09:15 更新 · 3 分钟前`，绝对时间 + 相对时间双显示，用户能一眼判断新鲜度。
- **stale 警告条**：`NOW() - last_fetched_at > 6h` 时显示黄色警告条 —— 明确告诉用户"数据可能不新鲜"，不撒谎。
- **Tab 切换联动**：切 Tab 时同步 —— 主区 feed 显隐、侧栏 active 状态、Hero 里的更新时间 chip、数据源 chip、stale 警告 —— 让用户随时明白当前看的是"哪个源的什么时候的数据"。

### 4.5 情感化与克制感（原 4.4）

**立场：这里不是社区，是"世界的窗户"——克制、干净、看一眼就走。**

- **零红点 / 零徽章 / 零"未读"** — 这不是 IM，不制造回访焦虑。
- **零站内互动按钮** — 从 DOM 层面就没有点赞/评论/共情按钮的存在（不是 disabled，是压根没渲染）。
- **零动画干扰** — 卡片列表用最基础的 `fade + translateY(4px)` 入场；Hero 有一个 4s 周期的柔脉冲晕光，`prefers-reduced-motion` 下直接 0ms。
- **数据陈旧不撒谎** — `stale=true` 时头部黄色警告 + 具体上次抓取时间，不做"永远显示刚更新"的假象。
- **失败态友好** — 抓取器全挂时页面显示"最近 6 小时没成功抓到新数据，我们正在处理"，不空白也不 500 页。
- **打工人语气** — 空态文案"世界今天很安静，或者服务器在偷懒 🐢"，不用严肃 tech 词。
- **底部克制感声明卡** — `.quiet-strip` 明确告诉用户"为什么这里没有点赞、评论、共情按钮"，把产品性格前置说清楚，避免用户困惑。
- **双主题** — 沿用现有 `.light` class 双主题变量，不新增 token。

---

## 五、分阶段实施计划（首发 = M0~M3）

| 阶段 | 交付物 | 验收标准 |
|------|--------|---------|
| **M0 地基** | `007_pulse.sql` 迁移（pulse_topics + pulse_items + 种子 2 专题）；`internal/pulse` 模块骨架（model/repo/service/handler）；`internal/pulse/collector` 接口 + Registry；`internal/pulse/scheduler` 骨架；`AppIcon` 加 `activity`；TopNav/Sidebar 加"最近发生"入口（放在日记广场与排行榜之间，带 pulse-dot 呼吸绿点），先落到空白占位页 | 迁移可执行/回滚；模块编译通过；导航入口可点进空白页；`/api/pulse/topics` 返回 2 条种子专题；导航位置与呼吸绿点视觉与 `docs/design/07-pulse.html` 一致 |
| **M1 GitHub Trending 打通** | `GitHubTrending` collector（goquery 解析 + UA + 重试 + 失败不覆写）；scheduler 起来能每 60 分钟拉一次并写库；`/api/pulse/topics/:slug/items` 只读接口；`RepoCard` + `PulseFeed` + `PulseIndex` 页面首个 Tab | 后端能稳定抓取并入库；前端页面能显示 Top 25 且点开跳 GitHub 外链；`stale` 判定生效；抓取失败时旧数据保留 |
| **M2 AI 圈大事打通** | `HackerNewsAI` collector（Algolia API + 关键词并发 + 去重 + 时间衰减排序）；`NewsCard`；Tab 切换（联动 Hero 更新时间/数据源 chip）；语言色点常量 `constants/langColors.ts`；stale 警告条 / 呼吸绿丸 / 心电图水印；暗色亮色双主题走查；底部 `.quiet-strip` 克制感声明卡 | AI Tab Top 20 稳定；HN 讨论副按钮打开正确页面；缓存 60s 生效；移动端 375px 可读；WCAG AA；视觉与 `docs/design/07-pulse.html` 一致 |
| **M3 打磨与部署（= 首发完成）** | Redis/进程缓存 60s；每日清理 goroutine（7 天窗口）；空态/错态/骨架屏；API 与前端 e2e 测试；Nginx 反代 + Docker Compose 部署 | 缓存命中率 ≥90%；清理任务不误删；`make test` 通过；上线 http://39.107.58.169/pulse 可直接访问（未登录也能看） |

> 首发 = M0+M1+M2+M3，跑通"抓 → 排 → 展示"这一条最短闭环。验证 CTR/回访后再进 P1。

---

## 附录 A：P1 / P2 路线图（首发不做，验证后按序切入）

**P1（若北极星达标，即"最近发生"有稳定周活，增强"看得爽"）**

- **新增专题**（每个新专题 = 一个 Collector）：
  - 🔥 **HN Top Stories**（`hn-top`）：不限 AI 的 HN 全站高分帖
  - 💰 **技术公司动态**（`tech-news`）：从 TechCrunch/The Verge/InfoQ RSS 抓摘要
  - 📉 **裁员/融资动态**（`layoffs`）：抓 layoffs.fyi + 猎聘/Boss 公开信号（若源站合规）
  - 🇨🇳 **国内技术圈**（`domestic-tech`）：掘金/InfoQ 中国站 RSS
  - 📢 **AI 官方公告**（`ai-official`）：OpenAI/Anthropic/Google DeepMind 官方博客 RSS 白名单
- **卡片交互增强**：整卡懒加载 GitHub OpenGraph 图 / HN 帖首屏预览截图（成本较高，需评估）。
- **筛选**：GitHub Tab 加"语言筛选"chip（Go/Python/TS/Rust…）。
- **时间维度切换**：GitHub 加 daily / weekly / monthly Tab 切换（对应 GitHub 官方 `?since=` 参数）。
- **失败降级**：抓取失败时兜底展示某个静态"精选"（人工兜底一批），避免完全空。

**P2（内容留存证明后的深度化，可能改变克制感定位，需重新评审）**

- **收藏 / "我已经看过"标记**：登录用户可收藏和标记已读，跨设备同步。这是**克制感失守的第一道口子**，一定要小心；建议做成"不打扰的书签"（无红点、无提醒、无社交暴露）。
- **RSS 订阅输出**：把 `/api/pulse/topics/:slug/items` 输出成 RSS 2.0 供第三方阅读器订阅，反哺"打工人不切平台"的定位。
- **周报邮件**（`weekly digest`）：每周一早上给订阅用户发一封"上周最近发生 Top 20"邮件；纯 opt-in、可一键退订。
- **AI 中文摘要**（可选）：给 HN AI 专题的英文帖用 Claude Haiku 生成 1 句中文摘要 + 免责声明"AI 摘要仅供参考"。上线前评估幻觉率与成本。
- **专题订阅 / 关键词自定义**：登录用户可自定义关注哪些专题（而不是全部一次性看）、加自己的关键词生成"我的最近发生"。这是社区化的深水区，见 P3。
- **P3（社区化，需另行评审是否要做）**：评论区 / 站内讨论线程。**目前立场：不做。** 一旦做了，"最近发生"就从"世界的窗户"退化成又一个吐槽频道，与「频道」板块严重重叠。

---

## 附录 B：开放问题与待决策

1. **GitHub Trending 抓取稳定性**：GitHub Trending 页面 HTML 结构历史上有过变更，一旦变了 collector 就挂。**建议**：M1 阶段并行观察 30 天，同时评估备选方案（自建 Star 增量监控 —— 定时轮询 GitHub Search API `stars:>1000 pushed:>{yesterday}` 拿变化）。方案二成本更高但更稳。
2. **HN AI 关键词覆盖度**：初始关键词表 `openai|anthropic|claude|gpt|llm|gemini|mistral|deepseek` 是否够？会不会漏掉如"大模型/foundation model"这类中文/宽泛表达？**建议**：M2 上线后每周产品复盘一次覆盖度，按需扩充；不做全自动关键词发现（成本高、噪音大）。
3. **多实例部署时的抓取重复**：目前单实例部署，`time.Ticker` 直接跑。若切多实例部署（如 K8s），需要**任一时刻只有一个副本在抓**——最简单是用 Redis 分布式锁 `SET NX PX`；或独立出一个抓取 sidecar 单副本部署。**当前阶段不做，见 3.6**。
4. **数据源版权与 robots.txt**：GitHub Trending 页无 robots 禁抓且 UI 上鼓励镜像；HN Algolia 是 HN 官方鼓励使用的 API。二者合规无虞。**未来加专题**（尤其国内媒体）必须**先看 robots.txt + ToS**；如遇动态渲染/反爬严重的源站（如今日头条/微信公众号），直接不做，绕过风险。
5. **命名统一**：中文"最近发生"、英文 `Pulse`、后端代号 `pulse`、迁移 `007_pulse.sql`、路由前缀 `/api/pulse/`、前端页面 `pages/pulse/index.vue`、组件目录 `components/pulse/`、设计稿 `docs/design/07-pulse.html`——多处命名需一致。
6. **暗示 AI 摘要边界**：一旦上 AI 摘要（P2），必须**明示标注"AI 生成，可能有误"** + **保留跳原文链接**，规避 AI 幻觉误导用户造成社区信任下滑。
7. **入口分流风险**：加了"最近发生"后，会不会分掉「频道/日记」的注意力反而降低核心 UGC 活跃？**建议**上线 4 周做一次数据回顾：如果 UGC 活跃率同比未下降（甚至因用户在站时长增加而上升），则确认协同关系；若明显下滑，需重新评估"最近发生"在首页导航的位置或折叠优先级。

---

> **本文档为首发聚焦版设计草案；视觉原型 `docs/design/07-pulse.html` 已完成并对齐线上，后端 / 前端代码尚未实现。** 待评审后再拆分任务进入开发；实施前需产品拍板 2.3 中的两个专题最终字段清单，以及 3.5 的时间衰减半衰期是否需调整。
