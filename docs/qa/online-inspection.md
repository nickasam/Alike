# Alike 线上服务器巡检报告

> **目标主机：** http://39.107.58.169（Docker Compose + Nginx :80）
> **巡检方式：** 每 30 分钟自动复测（黑盒 API 探测 + 服务器资源/容器/日志检查）
> **说明：** 本文档由自动化巡检持续追加。每轮仅记录「需修复的 Bug」与「需加固项」两类结论。

---

## 当前加固清单（按优先级，滚动更新）

### 🔴 需要修复的 Bug
- 暂无应用层 bug。但第 4 轮发生**线上宕机事故**（见下），根因为 OOM，属基础设施缺陷，加固项已升级。

### 🟠 需要加固

| 优先级 | 项 | 现状 | 建议 |
|--------|----|------|------|
| **🔴P0** | **内存不足致 OOM 宕机** | 第4轮实测：内核 OOM 杀死 dockerd，6 容器全 Exited、服务中断 ≥13 分钟无人知 | ①compose 全服务加 `restart: always` / `unless-stopped` 确保 OOM 后自愈；②各容器设 `mem_limit`（尤其 Postgres/MinIO）；③加 swap；④清理服务器上无关容器（openclaw 等挤占内存）；⑤升级内存 |
| **P1** | 无 HTTPS | 仅 HTTP，:443 不通，JWT/密码明文传输 | Let's Encrypt + nginx 443 + HSTS |
| **P1** | 磁盘 74%(28G/40G) + Docker 日志无轮转 | 稳定未增长；主体为数据卷 | 配 Docker 日志 `max-size`/`max-file` + Postgres 定时备份清理 |
| **P1** | 无监控告警 | OOM 宕机 13 分钟无告警，靠巡检偶然发现 | 接 Prometheus + 内存/存活/错误率告警（本次事故直接证明其必要性） |
| **P2** | 无结构化日志 + request_id | 仅 `log.Printf` 文本，线上排障困难 | slog JSON + X-Request-ID 中间件 |
| **P2** | CI 未接部署门禁 | 部署仍手动 rsync+rebuild | CI 接 PR 门禁 + 部署脚本固定 `rsync -avc` |
| **⚠️ 跟踪** | 服务器混部无关容器 | 存在 `dk_openclaw-*`、`condescending_chaplygin` 等非本项目容器 | 评估迁出或限制其资源，避免挤占 Alike |

---

## 巡检历史

### 第 4 轮（2026-07-10）🔴 线上宕机事故
- ❌ **首次探测全部 000（连接被拒）——服务中断**。诊断：主机 ping/SSH 通、负载正常，但 :80 无监听、`docker compose ps` 空。
- 根因：`docker ps -a` 显示 6 容器 **13 分钟前同时 Exited(0)**；`dmesg` 确认 **内核 OOM 杀死 dockerd**（`Out of memory: Killed process dockerd`）。内存打满 → OOM 杀 dockerd → daemon 重启后容器未自动拉起（compose 未配 restart 策略）。
- 加剧因素：服务器混部了 `dk_openclaw-*`、`condescending_chaplygin` 等非本项目容器，挤占 1.8G 内存。
- 处置：`docker compose up -d` 手动拉起，验证 health ok、接口 200、6 容器 healthy，服务恢复。
- **教训**：前几轮标记的"内存偏小(P3)""docker daemon 不稳定(跟踪)"两隐患合流爆发 → 升级为 **P0**；无 restart 策略 + 无告警使 13 分钟宕机无人知。

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
