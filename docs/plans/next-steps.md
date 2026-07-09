# Alike 项目下一步计划

> **编写时间：** 2026-07-10 03:10
> **当前状态：** 阶段零（文档+设计稿）✅ 已完成，阶段一（脚手架）待开始
> **编写人：** Hermes（架构师）

---

## 一、阶段零交付物清单（已完成）

| 交付物 | 路径 | 状态 |
|--------|------|------|
| 架构设计实施计划 | `docs/plans/architecture-design.md` | ✅ 权威文档，不可改 |
| 设计系统 | `docs/design/design-system.md` | ✅ Aurora 极光风，暗色默认+亮色可切换 |
| 首页原型 | `docs/design/01-home.html` | ✅ Aurora 极光风 |
| 频道页原型 | `docs/design/02-channel.html` | ✅ Aurora 极光风 |
| 日记页原型 | `docs/design/03-diary.html` | ✅ Aurora 极光风 |
| 排行榜原型 | `docs/design/04-ranking.html` | ✅ Aurora 极光风 |
| 个人主页原型 | `docs/design/05-profile.html` | ✅ Aurora 极光风 |
| 系统架构图 | `docs/architecture-diagram.html` | ✅ Aurora 极光风 |
| 组件交互规范 | `docs/design/component-spec.md` | ✅ 8 个核心组件 |
| PRD | `docs/prd/PRD.md` | ✅ 645 行，含 P0-P2 功能、用户画像、竞品分析 |
| 用户故事 | `docs/stories/user-stories.md` | ✅ 23 个用户故事 |
| 测试计划 | `docs/qa/test-plan.md` | ✅ 7 模块测试用例 + 回归策略 |
| 团队角色定义 | `docs/team-roles.md` | ✅ 7 个 Claude 角色 |
| .gitignore | `.gitignore` | ✅ 已创建 |

---

## 二、下一步：阶段一 — 项目脚手架

**目标：** 搭建前后端项目骨架，Docker 开发环境可用，数据库迁移就绪。

**预计耗时：** 2-3 小时

### 任务分配

| # | 任务 | 负责角色 | 输出 | 验收标准 |
|---|------|---------|------|---------|
| 1.1 | 初始化 Go 后端 | `claude-backend` | `backend/` | `go mod init`, Gin 框架, 目录结构, 配置加载, `go build` 通过 |
| 1.2 | 初始化 Nuxt 前端 | `claude-frontend` | `frontend/` | `npx nuxi init`, Tailwind CSS, 目录结构, `npm run dev` 可启动 |
| 1.3 | Docker Compose | `claude-devops` | `docker-compose.yml` | PostgreSQL + Redis + MinIO + Nginx 容器可启动 |
| 1.4 | 数据库迁移 | `claude-backend` | `backend/migrations/001_init.sql` | 7 张核心表 + 索引创建成功 |
| 1.5 | Makefile | `claude-devops` | `Makefile` | `make dev/build/test/migrate/seed/lint` 命令可用 |
| 1.6 | .env 模板 | `claude-devops` | `.env.example` | 所有环境变量有默认值和注释 |
| 1.7 | Nginx 配置 | `claude-devops` | `nginx/nginx.conf` | HTTP + WebSocket 反代配置 |

### 执行顺序

```
1.3 Docker Compose（先搭环境）
  ├── 1.1 Go 后端初始化（并行）
  └── 1.2 Nuxt 前端初始化（并行）
       ↓
1.4 数据库迁移（依赖 1.1）
       ↓
1.5 + 1.6 + 1.7（并行收尾）
```

### 预置频道种子数据

迁移完成后灌入 15 个预置频道：

| 分类 | 频道 |
|------|------|
| industry | #互联网、#工厂、#外卖骑手、#医护、#教培、#服务业 |
| job | #程序员、#产品经理、#设计师、#流水线、#骑手 |
| topic | #吐槽大会、#打工日记、#薪资揭秘、#摸鱼技巧 |

---

## 三、阶段二 — 用户认证系统

**目标：** 用户可注册、登录、获取 JWT、访问受保护接口。

**预计耗时：** 2-3 小时

| # | 任务 | 负责角色 |
|---|------|---------|
| 2.1 | 注册接口（邮箱+bcrypt） | `claude-backend` |
| 2.2 | 登录接口（JWT 签发+Redis 会话） | `claude-backend` |
| 2.3 | JWT 中间件（请求鉴权） | `claude-backend` |
| 2.4 | 前端登录/注册页 | `claude-frontend` |
| 2.5 | 前端认证守卫（路由守卫+自动刷新） | `claude-frontend` |
| 2.6 | 单元测试 | `claude-qa` |

---

## 四、阶段三 — 频道系统

**目标：** 用户可浏览频道、加入/离开频道、查看频道详情。

**预计耗时：** 2-3 小时

| # | 任务 | 负责角色 |
|---|------|---------|
| 3.1 | 频道 CRUD API | `claude-backend` |
| 3.2 | 频道成员管理（加入/离开/列表） | `claude-backend` |
| 3.3 | 种子数据脚本 | `claude-backend` |
| 3.4 | 前端频道侧边栏 | `claude-frontend` |
| 3.5 | 前端频道页 | `claude-frontend` |

---

## 五、阶段四 — 消息与线程 + WebSocket（核心）

**目标：** 实时聊天可用，消息发布→WebSocket 推送→实时显示。

**预计耗时：** 3-4 小时

| # | 任务 | 负责角色 |
|---|------|---------|
| 4.1 | 消息发布 API | `claude-backend` |
| 4.2 | 消息列表 API（分页） | `claude-backend` |
| 4.3 | 线程回复 API | `claude-backend` |
| 4.4 | WebSocket Hub（连接管理+频道广播） | `claude-backend` |
| 4.5 | Redis Pub/Sub 多实例广播 | `claude-backend` |
| 4.6 | 前端消息列表（分页加载） | `claude-frontend` |
| 4.7 | 前端消息输入（表情+情绪标签） | `claude-frontend` |
| 4.8 | 前端线程面板 | `claude-frontend` |
| 4.9 | WebSocket 客户端 | `claude-frontend` |

---

## 六、阶段五~九概览

| 阶段 | 内容 | 预计耗时 |
|------|------|---------|
| 五 | 情绪系统（情绪标签+看板+实时更新） | 2h |
| 六 | 共情机制（抱团取暖+实时通知+排行榜） | 2h |
| 七 | 打工日记（CRUD+打卡日历+日记广场） | 2h |
| 八 | 排行榜与通知（三个榜单+通知系统） | 2h |
| 九 | 搜索+文件上传+匿名模式+UI 打磨+部署 | 2-3h |

---

## 七、MVP 里程碑

**MVP v1.0 = 阶段一 ~ 阶段九全部完成**

验收标准：
- [ ] 用户可完成注册→登录→加入频道→发消息→选情绪→收到回复→收到共情→收到通知的完整流程
- [ ] WebSocket 实时推送正常（新消息、共情、通知）
- [ ] 情绪看板实时展示频道情绪分布
- [ ] API P95 响应时间 < 200ms
- [ ] 消息推送延迟 < 2 秒
- [ ] 首页首屏加载 < 3 秒
- [ ] 暗色主题视觉完整（Aurora 极光风格）
- [ ] PC + 移动端响应式适配完成

---

## 八、派活策略

| 模式 | 说明 |
|------|------|
| **并行** | 前后端可并行开发（API 契约已定义），DevOps 独立 |
| **Claude Code** | 委派给 Claude Code CLI 执行具体编码 |
| **Hermes 监工** | 派活→验收→打回/通过→汇总 |
| **最多 2 agent** | 同时最多 2 个 Claude Code session 并行 |

### 典型工作流

```
Hermes 编写任务卡 (.claude-task-N.md)
  ↓
启动 Claude Code session（print mode -p）
  ↓
Claude Code 读取任务卡 → 执行 → 输出结果
  ↓
Hermes 验收（编译/测试/浏览器检查）
  ├── 通过 → 进入下一任务
  └── 不通过 → 打回，附修改意见，重新执行
  ↓
阶段完成 → 更新 team-roles.md → 通知用户验收
```

---

## 九、风险与注意事项

| 风险 | 应对 |
|------|------|
| WebSocket 多实例广播 | 必须用 Redis Pub/Sub，不能仅靠内存 |
| 消息量大时查询慢 | 消息表按 channel_id + created_at 索引，分页不用 OFFSET |
| 匿名模式隐私泄露 | 匿名消息不返回 user_id，JWT 中不暴露匿名消息关联 |
| 情绪看板实时性 | WebSocket 推送情绪更新，Redis 缓存聚合数据 |
| 文件上传安全 | 限制文件类型和大小，图片压缩 |
| API 限流 | 单 IP 每分钟 ≤ 60 次，发消息 ≤ 20 条/分钟 |

---

> **下一步行动：** 用户确认后，启动阶段一——派 `claude-devops` 搭 Docker 环境 + `claude-backend` 初始化 Go 项目。
