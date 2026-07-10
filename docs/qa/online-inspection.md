# Alike 线上服务器巡检报告

> **目标主机：** http://39.107.58.169（Docker Compose + Nginx :80）
> **巡检方式：** 每 30 分钟自动复测（黑盒 API 探测 + 服务器资源/容器/日志检查）
> **说明：** 本文档由自动化巡检持续追加。每轮仅记录「需修复的 Bug」与「需加固项」两类结论。

---

## 当前加固清单（按优先级，滚动更新）

### 🔴 需要修复的 Bug
- 暂无。核心链路（health / 认证鉴权 / 限流 / 分页边界 / 错误码语义 / 防枚举 / 上传鉴权）连续 3 轮探测均正确。

### 🟠 需要加固

| 优先级 | 项 | 现状 | 建议 |
|--------|----|------|------|
| **P1** | 无 HTTPS | 仅 HTTP，:443 不通，JWT/密码明文传输 | Let's Encrypt + nginx 443 + HSTS |
| **P1** | 磁盘 74%(28G/40G) + Docker 日志无轮转 | 稳定未增长；主体为数据卷（Docker JSON 日志仅 208K，已排除） | 配 Docker 日志 `max-size`/`max-file` + Postgres 定时 pg_dump 备份清理 |
| **P2** | 无监控告警 | 磁盘/内存/错误率/daemon 健康无告警，靠人工发现 | 接 Prometheus + 告警 |
| **P2** | 无结构化日志 + request_id | 仅 `log.Printf` 文本，线上排障困难 | 引入 slog JSON + X-Request-ID 中间件 |
| **P2** | CI 未接部署门禁 | 部署仍手动 rsync+rebuild（踩过 rsync -az 漏传、nginx bind-mount 需 force-recreate） | CI 接 PR 门禁 + 部署脚本固定 `rsync -avc` |
| **P3** | 内存 1.8G 偏小 | 当前用 ~800M 有余量，无冗余 | 确认 swap；关注数据增长后余量 |
| **⚠️ 跟踪** | Docker daemon / 磁盘 I/O 重操作不稳定 | 连续 2 轮：`docker system df` 返回 `docker.sock EOF`、`du -xh /` 全盘扫描超时；轻命令(`df`/`docker compose ps`/HTTP API)全部正常 | 查 `dmesg`/`journalctl` I/O error 或 OOM；`docker system prune` 清构建缓存；关注云盘限速/坏道 |

---

## 巡检历史

### 第 3 轮（2026-07-10）
- ✅ health `database/redis/minio` 全 ok；首页/频道/搜索 200；鉴权 401；限流 429 正常。
- ✅ 6 容器全 healthy；磁盘 74%（稳定未增长）；内存用 ~800M/1.8G；负载低。
- ⚠️ **新增跟踪项**：`docker system df` 报 `docker.sock EOF`、`du -xh /` 超时——docker daemon / 磁盘 I/O 在重操作下不稳定（已连续 2 轮出现，当前最值得深挖的基础设施隐患）。
- Bug：无。

### 第 2 轮（2026-07-10）
- ✅ 核心链路复测全绿；澄清上轮疑似的登录 422 实为密码<6位的 binding 校验（非 bug），合法长度错误密码正确返回 401「邮箱或密码错误」。
- ✅ 6 容器 healthy；磁盘 74%；内存 807M/1.8G；后端日志无 error/panic。
- Bug：无。

### 第 1 轮（2026-07-10）
- ✅ health 全 ok；安全头齐全（X-Frame/X-Content-Type/X-XSS/Referrer-Policy）；未认证访问 401；参数边界（page_size 上限 50、非法/不存在 id → 404、空搜索 422）正确；限流触发 429；防枚举一致；11MB body 被拦截。
- ✅ 6 容器 healthy；内存 841M/1.8G；负载 0.4；磁盘 74%；后端日志无 error/panic。
- ⚠️ 首次标记：无 HTTPS、磁盘 74%、无监控/结构化日志/CI 门禁。
- Bug：无。
