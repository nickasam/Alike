# Alike 测试计划文档

> **项目：** Alike — 汇聚天下牛马，总有人懂你的辛苦
>
> **版本：** v1.0 | **日期：** 2026-07-10 | **编写人：** claude-qa

---

## 一、测试策略概述

### 1.1 目标

本测试计划旨在为 Alike 聊天社区平台提供全面、系统化的质量保障方案，覆盖后端 API、前端组件、端到端用户流程、WebSocket 实时通信以及数据库集成等各个层面，确保产品在功能正确性、实时性、并发安全、数据一致性等方面达到上线标准。

### 1.2 测试原则

- **分层测试金字塔**：单元测试 > 集成测试 > E2E 测试，底层测试数量最多、执行最快，顶层测试数量最少但覆盖最关键用户路径。
- **自动化优先**：后端 Go 测试全部使用 `*_test.go` 自动化执行；前端使用 Vitest + Vue Test Utils；E2E 使用 Playwright；CI/CD 中自动触发。
- **数据隔离**：每个测试用例使用独立测试数据库，测试结束后自动清理，避免数据污染。
- **边界覆盖**：每个接口除正常路径外，必须覆盖参数缺失、类型错误、权限不足、并发冲突等异常路径。
- **实时性验证**：WebSocket 相关功能必须有专门的连接、断线重连、消息推送验证。

### 1.3 测试工具链

| 层 | 工具 | 说明 |
|----|------|------|
| 后端单元/集成测试 | Go `testing` + `testify` | `go test ./...` |
| HTTP API 测试 | `net/http/httptest` + `testify` | 模拟 Gin 路由 |
| 数据库测试 | `testcontainers-go` | 临时 PostgreSQL 实例 |
| Redis 测试 | `miniredis` 或 `testcontainers` | 内存 Redis 模拟 |
| 前端组件测试 | Vitest + `@vue/test-utils` | 组件挂载与交互 |
| E2E 测试 | Playwright | 浏览器自动化 |
| WebSocket 测试 | `gorilla/websocket` 客户端 | Go 测试内拨号 |
| CI/CD | GitHub Actions | PR 触发自动测试 |

---

## 二、测试范围

### 2.1 后端 API 测试

覆盖全部 12 个后端模块的 REST API 端点，包括：

- 请求参数校验（必填项、类型、长度、范围）
- JWT 认证与权限校验（未登录、token 过期、无权限访问他人资源）
- 业务逻辑正确性（注册去重、频道成员唯一约束、共情幂等等）
- 错误处理与统一响应格式（`code`/`message`/`data` 结构）
- 分页查询（首页、末页、越界页码）
- 并发安全（同时加入频道、同时共情同一消息）

### 2.2 前端组件测试

覆盖 6 个页面和核心组件：

- 组件渲染正确性（props 传入、slot 内容、条件渲染）
- 用户交互（点击、输入、表单提交、滚动加载）
- 状态管理（Pinia store 的状态变更、响应式更新）
- API 调用 mock（`useApi` composable 的请求/错误处理）
- WebSocket 事件驱动 UI 更新（新消息追加、共情数实时刷新）

### 2.3 E2E 测试

使用 Playwright 模拟真实用户操作流程：

- 完整注册→登录→加入频道→发消息→收到回复→共情→查看通知
- 日记创建→公开/私密切换→日记广场展示→评论
- 频道切换、消息分页加载、情绪看板实时更新
- 排行榜页面数据展示与排序验证
- 匿名模式发帖与身份隐藏验证

### 2.4 WebSocket 测试

- 连接建立与 JWT 鉴权
- `join_channel` / `leave_channel` 事件
- `send_message` → `new_message` 实时推送
- `typing` 事件广播
- `empathy` 事件推送
- 断线重连机制
- 多用户并发连接与消息广播一致性

### 2.5 DB 集成测试

- 数据库迁移脚本正确性（`migrations/001_init.sql`）
- 表结构约束验证（唯一约束、外键级联删除、非空约束）
- 索引有效性验证（`EXPLAIN ANALYZE` 关键查询）
- 事务回滚与数据一致性
- Redis Pub/Sub 消息发布与订阅可靠性

---

## 三、测试环境要求

### 3.1 开发/测试环境

| 组件 | 版本要求 | 说明 |
|------|----------|------|
| Go | ≥ 1.22 | 后端运行时 |
| Node.js | ≥ 20 LTS | 前端构建 |
| PostgreSQL | ≥ 15 | 测试数据库 |
| Redis | ≥ 7 | 缓存与 Pub/Sub |
| MinIO | 最新稳定版 | 对象存储（可用 minio server 本地运行） |
| Docker | ≥ 24 | 容器化测试环境 |
| Docker Compose | ≥ 2.20 | 多容器编排 |

### 3.2 环境配置

```bash
# 测试专用环境变量
DB_HOST=localhost
DB_PORT=5432
DB_NAME=alike_test        # 独立测试数据库
DB_USER=alike_test
DB_PASSWORD=alike_test_pwd

REDIS_ADDR=localhost:6379
REDIS_DB=1                # 测试专用 DB 编号

MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=test
MINIO_SECRET_KEY=test123
MINIO_BUCKET=alike-test

JWT_SECRET=test-jwt-secret-key
JWT_EXPIRE_HOURS=24
```

### 3.3 测试数据准备

- 使用 `make seed` 灌入种子数据（频道、测试用户）
- 每个测试套件运行前清空相关表，运行后清理
- 密码统一使用 bcrypt 加密后的测试密码
- 准备至少 3 个测试用户（不同行业/岗位）用于多用户交互测试

---

## 四、各模块测试用例要点

### 4.1 Auth 模块（认证）

| # | 测试用例 | 要点 |
|---|----------|------|
| A1 | 邮箱注册成功 | POST `/api/auth/register`，合法邮箱+密码，返回 JWT，用户写入 DB，密码为 bcrypt hash |
| A2 | 重复邮箱注册失败 | 同一邮箱二次注册，返回 `code != 0`，错误信息提示邮箱已存在 |
| A3 | 注册参数校验 | 缺少邮箱/密码、密码过短（<6位）、邮箱格式非法，均返回 400 |
| A4 | 登录成功 | 正确邮箱+密码，返回 access_token + refresh_token，token 可解析出 user_id |
| A5 | 登录失败 | 错误密码返回 401；不存在的邮箱返回 404；密码错误次数过多触发限流 |
| A6 | JWT 刷新 | 使用 refresh_token 调用 `/api/auth/refresh`，获取新 token；过期 refresh_token 返回 401 |
| A7 | 获取当前用户 | GET `/api/auth/me` 携带有效 JWT，返回用户完整信息；无 token 返回 401 |

### 4.2 Channel 模块（频道）

| # | 测试用例 | 要点 |
|---|----------|------|
| C1 | 获取频道列表 | GET `/api/channels`，按 category 分类返回，包含 member_count，支持分页 |
| C2 | 频道详情查询 | GET `/api/channels/:id`，返回频道完整信息；不存在的 ID 返回 404 |
| C3 | 加入频道 | POST `/api/channels/:id/join`，member_count +1，channel_members 写入记录，重复加入返回幂等成功 |
| C4 | 离开频道 | POST `/api/channels/:id/leave`，member_count -1，channel_members 删除记录；未加入时离开返回 200 |
| C5 | 频道成员列表 | GET `/api/channels/:id/members`，返回成员列表含角色(member/admin)，支持分页 |
| C6 | 频道情绪看板 | GET `/api/channels/:id/emotion-board`，返回各情绪标签计数，数据与 messages 表 emotion 字段一致 |
| C7 | 自建频道申请 | POST 创建频道（需审核逻辑），name/slug 唯一性校验，slug 重复返回 409 |

### 4.3 Message 模块（消息）

| # | 测试用例 | 要点 |
|---|----------|------|
| M1 | 发布消息 | POST `/api/channels/:id/messages`，内容写入 messages 表，emotion 字段可选，返回消息完整对象 |
| M2 | 消息分页查询 | GET `/api/channels/:id/messages`，按 created_at DESC 排序，page/page_size 参数，返回 total/page/page_size |
| M3 | 线程回复 | POST `/api/messages/:id/replies`，parent_id 关联正确，GET `/api/messages/:id/threads` 返回回复列表 |
| M4 | 删除消息 | DELETE `/api/messages/:id`，仅本人或管理员可删除；删除后 attachments 级联删除；他人删除返回 403 |
| M5 | 匿名消息 | is_anonymous=true 时，返回中 user 信息脱敏（昵称隐藏或显示"匿名牛马"） |
| M6 | 消息内容校验 | 空内容返回 400；超长内容（>5000字符）返回 400；XSS 内容需转义存储 |
| M7 | @提及功能 | 消息内容中 @用户 生成 notification 记录，type=mention，被提及用户可收到通知 |

### 4.4 Emotion 模块（情绪）

| # | 测试用例 | 要点 |
|---|----------|------|
| E1 | 情绪标签列表 | 获取系统预设情绪标签列表（疲惫/愤怒/委屈/崩溃/麻木/想润等），含 emoji 和 key |
| E2 | 发消息带情绪 | POST 消息时 emotion="tired"，DB 存储正确，返回对象包含 emotion 字段 |
| E3 | 频道情绪看板聚合 | 多条不同情绪消息后，GET 情绪看板返回各情绪计数正确，百分比计算无误 |
| E4 | 情绪趋势统计 | 按时间范围（日/周/月）统计情绪趋势，返回时间序列数据，空数据返回空数组 |
| E5 | 无效情绪标签 | 传入不存在的 emotion key，返回 400 错误，消息不写入 DB |
| E6 | 情绪看板实时更新 | 新消息发布后，WebSocket `emotion_update` 事件推送更新后的看板数据 |

### 4.5 Empathy 模块（共情）

| # | 测试用例 | 要点 |
|---|----------|------|
| P1 | 共情消息 | POST `/api/messages/:id/empathy`，empathies 表写入，messages.empathy_count +1，用户 empathy_received +1 |
| P2 | 重复共情幂等 | 同一用户对同一消息二次共情，因 UNIQUE(message_id, user_id) 约束返回幂等成功，count 不重复增加 |
| P3 | 取消共情 | DELETE `/api/messages/:id/empathy`，删除记录，empathy_count -1，count 不为负 |
| P4 | 共情用户列表 | GET `/api/messages/:id/empathy-users`，返回共情过该消息的用户列表，支持分页 |
| P5 | 共情排行榜 | GET `/api/ranking/empathy`，按 empathy_count DESC 排序，返回最受共情帖子列表 |
| P6 | 共情实时推送 | 共情后 WebSocket 推送 `empathy` 事件，包含 message_id、from_user、更新后 count |
| P7 | 匿名共情 | 匿名用户的共情行为不暴露真实身份，from_user 显示"匿名牛马" |

### 4.6 Diary 模块（打工日记）

| # | 测试用例 | 要点 |
|---|----------|------|
| D1 | 创建日记 | POST `/api/diaries`，写入 diaries 表，title 可选、content 必填，mood 字段可选，streak_days 自动计算 |
| D2 | 日记广场 | GET `/api/diaries`，仅返回 is_public=true 的日记，按 created_at DESC 排序，支持分页 |
| D3 | 日记详情 | GET `/api/diaries/:id`，返回完整内容；私密日记仅本人可查看，他人访问返回 403 |
| D4 | 打卡日历 | GET `/api/diaries/streak/:user_id`，返回用户打卡记录，streak_days 连续天数计算正确（跨天中断归零） |
| D5 | 日记公开/私密切换 | 更新 is_public 字段，切换后广场列表实时反映，私密日记从广场消失 |
| D6 | 连续打卡逻辑 | 连续多天写日记，streak_days 递增；间隔一天后重新写，streak_days 归零重新计数 |
| D7 | 日记内容校验 | 空内容返回 400；超长内容校验；mood 字段值合法性校验 |

### 4.7 WS 模块（WebSocket）

| # | 测试用例 | 要点 |
|---|----------|------|
| W1 | WebSocket 连接建立 | WS `/ws?token=<JWT>`，有效 token 连接成功；无 token 或无效 token 连接被拒（401 关闭） |
| W2 | 加入/离开频道 | 发送 `join_channel` 后，该连接收到频道内 `new_message` 事件；`leave_channel` 后不再收到 |
| W3 | 实时消息推送 | 用户 A 发送 `send_message`，同频道用户 B 收到 `new_message` 事件，data 内容与 DB 一致 |
| W4 | 输入中状态广播 | 用户 A 发送 `typing` 事件，同频道其他用户收到 typing 通知，is_typing=false 时停止 |
| W5 | 断线重连 | 主动断开连接后重新连接，JWT 仍有效则连接成功，重新加入之前的频道，消息不丢失（补拉 REST 历史） |
| W6 | 多连接并发 | 同一用户多设备同时连接，消息推送到所有连接；连接超时自动清理 |
| W7 | 共情事件推送 | 用户 A 对消息共情，被共情者收到 `empathy` 类型 WebSocket 事件，包含 message_id 和更新后的 count |
| W8 | 通知推送 | @提及、回复、共情触发 `notification` 事件，客户端实时收到通知数据 |

---

## 五、回归测试策略

### 5.1 回归触发时机

| 时机 | 策略 |
|------|------|
| 每次代码合并到 main 分支 | 运行全部自动化测试套件 |
| 每次创建 PR | 运行受影响模块的测试 + 核心冒烟测试 |
| 版本发布前 | 全量回归 + 手工探索测试 |
| 紧急热修复后 | 受影响模块测试 + 冒烟测试 + 相关 E2E |

### 5.2 回归测试集划分

**冒烟测试集（Smoke Test）— 快速验证核心链路可用：**
1. 注册 → 登录 → 获取用户信息
2. 获取频道列表 → 加入频道 → 发送消息 → 查询消息列表
3. WebSocket 连接 → 接收新消息
4. 共情消息 → 查看排行榜
5. 创建日记 → 日记广场展示

**核心回归集（Core Regression）— 覆盖全部 7 个模块的主路径：**
- auth: A1, A4, A7
- channel: C1, C3, C5
- message: M1, M2, M3
- emotion: E2, E3
- empathy: P1, P3, P5
- diary: D1, D2, D4
- ws: W1, W2, W3

**完整回归集（Full Regression）— 全部测试用例：**
- 所有模块全部测试用例 + E2E 全流程 + DB 集成测试 + 性能基线

### 5.3 自动化回归流程

```
PR 提交 → GitHub Actions 触发
  ├─ 后端：go test ./... -race -cover
  ├─ 前端：npm run test:unit
  ├─ Lint：golangci-lint + eslint
  └─ 冒烟 E2E：Playwright smoke suite

合并到 main → 触发完整回归
  ├─ 全量后端测试（含 testcontainers DB 集成）
  ├─ 全量前端组件测试
  ├─ 全量 E2E 测试
  └─ 生成覆盖率报告（目标 ≥ 70%）
```

### 5.4 覆盖率要求

| 模块 | 最低行覆盖率 | 分支覆盖率 |
|------|-------------|-----------|
| auth | 90% | 80% |
| channel | 85% | 75% |
| message | 85% | 75% |
| emotion | 80% | 70% |
| empathy | 85% | 75% |
| diary | 80% | 70% |
| ws | 80% | 70% |
| middleware | 90% | 80% |
| **整体** | **≥ 75%** | **≥ 65%** |

---

## 六、Bug 严重等级定义

### P0 — 阻塞级（Blocker）

**定义：** 核心功能完全不可用，系统崩溃或数据丢失，无法绕过，阻碍发布。

**示例：**
- 注册/登录功能完全不可用
- WebSocket 连接无法建立，实时消息功能瘫痪
- 数据库连接失败导致全部 API 不可用
- 消息发布后数据丢失（写 DB 失败但返回成功）
- 安全漏洞：JWT 可伪造、密码明文存储、SQL 注入

**处理时效：** 立即修复，2 小时内提交修复，阻塞所有其他开发。

---

### P1 — 严重级（Critical）

**定义：** 核心功能严重受损或重要业务流程中断，有 workaround 但影响极大。

**示例：**
- 消息发送偶尔失败（>5% 概率）
- 频道加入后无法收到实时消息推送
- 共情次数统计错误（重复计数或不计数）
- 分页查询返回错误数据（重复/遗漏消息）
- 匿名模式下泄露用户真实身份
- 并发场景下 channel_members 出现脏数据

**处理时效：** 24 小时内修复，当前版本必须修复。

---

### P2 — 一般级（Major）

**定义：** 非核心功能异常或核心功能有轻微问题，不影响主流程但有用户体验影响。

**示例：**
- 情绪看板数据更新有延迟（>3秒）
- 排行榜排序偶尔不准确
- 日记打卡连续天数计算有边界 case 错误
- 通知列表分页第二页为空
- 前端组件在特定分辨率下布局错乱
- WebSocket 断线重连需要手动刷新页面
- 错误提示信息不准确或不友好

**处理时效：** 当前迭代内修复，不阻塞发布。

---

### P3 — 轻微级（Minor）

**定义：** UI 细节问题、文案错误、非功能性的体验优化建议，不影响任何业务功能。

**示例：**
- 按钮文字/提示文案有错别字
- 情绪 emoji 在部分浏览器显示不一致
- 页面加载时短暂闪烁（FOUC）
- 暗色主题下个别元素对比度不足
- 移动端适配的轻微间距问题
- 非关键接口响应略慢（>500ms 但 <2s）

**处理时效：** 排入 backlog，有空时修复或下个版本处理。

---

## 七、测试通过标准

| 标准 | 要求 |
|------|------|
| P0 Bug | 0 个 |
| P1 Bug | 0 个 |
| P2 Bug | ≤ 3 个（需评估确认可延期） |
| P3 Bug | 不限制（记录 backlog） |
| 冒烟测试 | 100% 通过 |
| 核心回归 | 100% 通过 |
| 代码覆盖率 | 满足第五章 5.4 节要求 |
| E2E 关键路径 | 全部通过 |
| 性能基线 | API P95 响应 < 200ms；WebSocket 消息延迟 < 500ms |

---

## 八、风险与应对

| 风险 | 影响 | 应对措施 |
|------|------|----------|
| testcontainers 启动慢导致 CI 时间过长 | CI 超时 | 使用本地 Docker PostgreSQL，复用容器 |
| WebSocket 测试不稳定（网络抖动） | 假阳性失败 | 增加超时容忍 + 重试机制 |
| MinIO 测试环境不可用 | 文件上传测试阻塞 | 使用 mock storage 接口，MinIO 集成测试独立标记 |
| 前端组件依赖 SSR 环境 | 组件测试困难 | 使用 `@vue/test-utils` 的 `mount` + mock Nuxt 上下文 |
| 数据库迁移脚本变更 | 测试数据不一致 | 每次测试前 `migrate down → up` 全量重建 |

---

> **附注：** 本测试计划为活文档，随项目迭代持续更新。新增功能模块时需同步补充对应测试用例。所有测试用例代码存放于 `backend/**/*_test.go`（后端）和 `frontend/test/`（前端），E2E 测试存放于 `e2e/` 目录。
