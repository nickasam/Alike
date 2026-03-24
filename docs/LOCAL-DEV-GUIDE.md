# Alike 本地开发环境运行指南

## ✅ 已完成的工作

### 1. 前端页面优化
- ✅ 优化了 `launcher.html` 页面
- ✅ 统一使用 Material Icons 图标
- ✅ 采用紫粉渐变配色方案
- ✅ 添加高质量背景图片
- ✅ 实现响应式布局（桌面端/移动端）
- ✅ 优化了 `login.html` 页面为全屏布局

### 2. 项目配置
- ✅ 创建了 `.env` 配置文件
- ✅ 创建了开发启动脚本 `scripts/start-dev.sh`

---

## 🚀 快速启动前端

### 方法一：使用启动脚本（推荐）

```bash
cd /Users/zhenghongfei6/go/src/github.com/Alike
./scripts/start-dev.sh
```

这个脚本会：
- 自动查找可用端口（8000-8081）
- 启动前端服务
- 在浏览器中打开页面
- 提供后端启动指南

### 方法二：手动启动

1. 找到可用端口：
```bash
lsof -i :8000 -i :8001 -i :8002 -i :3000
```

2. 启动服务（选择未被占用的端口）：
```bash
cd web/public
python3 -m http.server 8000
```

3. 在浏览器中访问：
- 首页: http://localhost:8000/launcher.html
- 登录: http://localhost:8000/login.html
- 主页: http://localhost:8000/index.html
- 全局聊天: http://localhost:8000/global-chat.html
- 个人资料: http://localhost:8000/profile.html

---

## 🔧 启动后端服务（可选）

由于 Docker Desktop 当前无法启动，有以下几种方案：

### 方案一：修复 Docker Desktop

1. **重启 Docker Desktop**：
```bash
killall Docker
open -a Docker
```

2. **如果仍然失败，尝试**：
   - 检查系统偏好设置中的 Docker
   - 重启计算机
   - 或重新安装 Docker Desktop

3. **启动数据库容器**：
```bash
cd /Users/zhenghongfei6/go/src/github.com/Alike
docker-compose -f deployments/docker/docker-compose.yml up -d postgres redis
```

### 方案二：使用本地 PostgreSQL

1. **安装 PostgreSQL**：
```bash
brew install postgresql
brew services start postgresql
```

2. **创建数据库**：
```bash
psql postgres
CREATE DATABASE alike_db;
CREATE USER alike_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE alike_db TO alike_user;
\q
```

3. **运行迁移**：
```bash
make migrate-up
```

4. **启动 API 服务**：
```bash
go run cmd/api/main.go
```

### 方案三：暂时只运行前端

如果只需要查看页面效果，可以只运行前端服务，暂时不启动后端：
- 所有页面都可以正常访问
- 登录功能会显示网络错误（这是正常的）
- 可以测试UI和交互效果

---

## 🔄 开发工作流程

### 1. 启动开发环境

#### 完整环境（前端 + 后端）

**前置要求**：
- Docker Desktop 已安装并运行
- Python 3.x 已安装
- Go 1.21+ 已安装

**启动步骤**：

```bash
# 1. 进入项目目录
cd /Users/zhenghongfei6/go/src/github.com/Alike

# 2. 启动数据库服务
docker-compose -f deployments/docker/docker-compose.yml up -d postgres redis

# 3. 运行数据库迁移
make migrate-up

# 4. 填充测试数据（可选）
psql -U alike_user -d alike_db -f db/seeds/seed.sql

# 5. 启动后端 API 服务
go run cmd/api/main.go

# 6. 启动前端服务（新终端窗口）
cd web/public
python3 -m http.server 8000

# 7. 访问应用
# 前端: http://localhost:8000/launcher.html
# 后端: http://localhost:8080/health
```

#### 仅前端开发（快速测试）

```bash
# 启动前端服务
cd /Users/zhenghongfei6/go/src/github.com/Alike
./scripts/start-dev.sh
```

### 2. 数据库操作

#### 运行迁移

```bash
# 升级到最新版本
make migrate-up

# 回滚一个版本
make migrate-down

# 查看迁移状态
make migrate-status

# 创建新迁移
migrate create -ext sql -dir db/migrations -seq create_new_table
```

#### 数据库连接

```bash
# 连接到 PostgreSQL
psql -U alike_user -d alike_db

# 或使用 Docker
docker exec -it alike-postgres psql -U alike_user -d alike_db
```

#### 常用数据库查询

```sql
-- 查看所有表
\dt

-- 查看用户数据
SELECT id, phone, nickname, created_at FROM users LIMIT 10;

-- 查看最近的匹配
SELECT * FROM matches ORDER BY created_at DESC LIMIT 10;

-- 查看聊天消息
SELECT * FROM messages ORDER BY created_at DESC LIMIT 10;

-- 查看在线用户
SELECT * FROM users WHERE last_online_at > NOW() - INTERVAL '5 minutes';

-- 统计数据
SELECT 
  (SELECT COUNT(*) FROM users) as total_users,
  (SELECT COUNT(*) FROM matches) as total_matches,
  (SELECT COUNT(*) FROM messages) as total_messages;
```

### 3. API 测试

#### 使用 curl 测试

**认证相关**：

```bash
# 注册新用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "password": "password123",
    "nickname": "测试用户"
  }'

# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "password": "password123"
  }'

# 获取当前用户信息（需要替换 YOUR_JWT_TOKEN）
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**用户相关**：

```bash
# 获取当前用户信息
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 更新用户信息
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "新昵称",
    "bio": "这是我的个人简介"
  }'

# 获取附近用户
curl -X GET "http://localhost:8080/api/v1/users/nearby?lat=31.2304&lng=121.4737&radius=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**全局聊天相关**：

```bash
# 获取全局聊天室信息
curl -X GET http://localhost:8080/api/v1/global/room \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 获取全局消息
curl -X GET "http://localhost:8080/api/v1/global/messages?limit=20&offset=0" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# 发送全局消息
curl -X POST http://localhost:8080/api/v1/global/messages \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "大家好！"
  }'
```

#### 使用 Postman 测试

1. 导入 API 集合（如果提供了 Postman collection）
2. 设置环境变量：
   - `base_url`: `http://localhost:8080`
   - `token`: 登录后获取的 JWT token
3. 运行请求测试

### 4. 代码提交

```bash
# 查看修改状态
git status

# 添加所有修改
git add .

# 提交代码（使用规范的提交信息）
git commit -m "feat: 添加全局聊天室功能"
git commit -m "fix: 修复用户登录问题"
git commit -m "docs: 更新 API 文档"
git commit -m "style: 格式化代码"
git commit -m "refactor: 重构用户服务"

# 推送到远程仓库
git push origin main
```

**提交信息规范**：
- `feat:` 新功能
- `fix:` 修复bug
- `docs:` 文档更新
- `style:` 代码格式（不影响功能）
- `refactor:` 重构
- `test:` 测试相关
- `chore:` 构建/工具相关

---

## 📁 项目结构

```
Alike/
├── web/public/              # 前端静态文件
│   ├── launcher.html        # 首页（已优化）✨
│   ├── login.html          # 登录页（已优化）✨
│   ├── index.html          # 主页
│   ├── global-chat.html    # 全局聊天
│   └── profile.html        # 个人资料
├── scripts/
│   └── start-dev.sh        # 开发启动脚本 ✨
├── .env                    # 环境配置（已创建）✨
└── deployments/docker/
    └── docker-compose.yml  # Docker 配置
```

---

## 🎨 页面设计特点

### Launcher Page (launcher.html)
- 🎯 **现代 Hero 区域**：大标题 + 渐变效果
- 🖼️ **精美卡片设计**：使用 Unsplash 高质量图片
- 📱 **完全响应式**：桌面端双列，移动端单列
- ✨ **Material Icons**：统一的图标风格
- 🎨 **紫粉渐变主题**：与登录页面保持一致

### Login Page (login.html)
- 📐 **全屏布局**：左右分栏（桌面端）
- 💫 **Logo 动画**：闪光效果
- 🎯 **紧凑表单**：优化的尺寸和间距
- 🔐 **登录/注册切换**：标签页设计

---

## ⚠️ 常见问题排查

### 问题 1：数据库连接失败

**症状**：
```
connection refused
connection timeout
could not connect to server
```

**原因**：PostgreSQL 未启动或配置错误

**解决方案**：

```bash
# 检查 PostgreSQL 状态
brew services list

# 启动 PostgreSQL
brew services start postgresql

# 检查 Docker 容器状态（如果使用 Docker）
docker ps | grep postgres

# 启动 Docker 容器
docker-compose -f deployments/docker/docker-compose.yml up -d postgres

# 验证连接
psql -U alike_user -d alike_db -c "SELECT version();"
```

**检查配置**：
```bash
# 查看 .env 文件中的数据库配置
cat .env | grep DB_

# 确认以下配置正确：
# DB_HOST=localhost
# DB_PORT=5432
# DB_NAME=alike_db
# DB_USER=alike_user
# DB_PASSWORD=your_password
```

---

### 问题 2：端口已被占用

**症状**：
```
bind: address already in use
Error: listen tcp :8000: bind: address already in use
```

**解决方案**：

```bash
# 查看占用端口的进程
lsof -i :8000
lsof -i :8080

# 杀死占用端口的进程
kill -9 <PID>

# 或使用其他端口
python3 -m http.server 8001

# 修改后端端口（在 .env 或代码中）
SERVER_PORT=8081
```

**批量查找可用端口**：
```bash
# 查找 8000-8080 之间的可用端口
for port in {8000..8080}; do
  lsof -i :$port > /dev/null 2>&1 || echo "Port $port is available"
done
```

---

### 问题 3：CORS 错误

**症状**：
```
Access to XMLHttpRequest at 'http://localhost:8080/api/v1/auth/login' 
from origin 'http://localhost:8000' has been blocked by CORS policy
```

**原因**：前后端跨域请求被阻止

**解决方案**：

```bash
# 确认后端已启用 CORS 中间件
# 检查 internal/api/middleware/common.go

# 或在开发环境中临时禁用浏览器安全策略
# Chrome (Mac)
open -n -a /Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome \
  --args --user-data-dir="/tmp/chrome_dev_test" \
  --disable-web-security \
  --disable-features=BlockInsecurePrivateNetworkRequests

# Chrome (Windows)
chrome.exe --user-data-dir="C:/chrome_dev" --disable-web-security
```

---

### 问题 4：JWT Token 无效

**症状**：
```
invalid token
token expired
unauthorized
```

**解决方案**：

```bash
# 1. 重新登录获取新 token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","password":"password123"}'

# 2. 检查 token 过期时间
# 使用 jwt.io 解码查看 exp 字段

# 3. 刷新 token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Authorization: Bearer YOUR_REFRESH_TOKEN"

# 4. 检查 JWT_SECRET 配置
cat .env | grep JWT_SECRET
```

---

### 问题 5：前端页面无法访问

**症状**：
- 浏览器显示 404 Not Found
- 页面空白
- 资源加载失败

**解决方案**：

```bash
# 1. 确认前端服务正在运行
lsof -i :8000

# 2. 检查文件路径
ls -la web/public/

# 3. 确认访问的 URL 正确
# 正确: http://localhost:8000/launcher.html
# 错误: http://localhost:8000/launcher

# 4. 检查浏览器控制台
# 打开开发者工具 (F12) 查看 Console 和 Network 标签

# 5. 清除浏览器缓存
# Chrome: Cmd+Shift+Delete (Mac) / Ctrl+Shift+Delete (Windows)
```

---

### 问题 6：Go 模块依赖问题

**症状**：
```
cannot find package
module not found
go: cannot find main module
```

**解决方案**：

```bash
# 1. 初始化 Go 模块
go mod tidy

# 2. 下载依赖
go mod download

# 3. 验证依赖
go mod verify

# 4. 清理缓存
go clean -modcache

# 5. 重新下载
go mod download
```

---

### 问题 7：Docker 相关问题

**症状**：
- Docker 无法启动
- 容器无法创建
- 网络连接问题

**解决方案**：

```bash
# 1. 重启 Docker Desktop
killall Docker
open -a Docker

# 2. 检查 Docker 状态
docker info
docker ps

# 3. 清理 Docker 资源
docker system prune -a

# 4. 重新构建镜像
docker-compose -f deployments/docker/docker-compose.yml build

# 5. 查看容器日志
docker-compose -f deployments/docker/docker-compose.yml logs

# 6. 重启容器
docker-compose -f deployments/docker/docker-compose.yml restart
```

---

### 问题 8：迁移失败

**症状**：
```
migration failed
error applying migration
database schema mismatch
```

**解决方案**：

```bash
# 1. 查看迁移状态
make migrate-status

# 2. 回滚到上一个版本
make migrate-down

# 3. 手动执行 SQL
psql -U alike_user -d alike_db -f db/migrations/000001_init_schema.up.sql

# 4. 检查迁移文件
ls -la db/migrations/

# 5. 查看错误日志
# 在迁移命令后添加 --verbose 标志
make migrate-up --verbose
```

---

### 问题 9：环境变量未生效

**症状**：
- 配置的值未被读取
- 使用了默认值
- 程序运行异常

**解决方案**：

```bash
# 1. 确认 .env 文件存在
ls -la .env

# 2. 检查 .env 文件权限
chmod 644 .env

# 3. 验证环境变量
source .env
echo $DB_NAME
echo $JWT_SECRET

# 4. 在 Go 代码中加载 .env
# 确保使用了 godotenv 包
import "github.com/joho/godotenv"
godotenv.Load()

# 5. 重启应用
# 环境变量只在启动时加载
```

---

### 问题 10：API 响应慢或超时

**症状**：
- 请求响应时间长
- 请求超时
- 页面加载缓慢

**解决方案**：

```bash
# 1. 检查数据库连接
psql -U alike_user -d alike_db -c "SELECT version();"

# 2. 查看数据库性能
psql -U alike_user -d alike_db -c "
  SELECT schemaname, tablename, seq_scan, idx_scan 
  FROM pg_stat_user_tables 
  ORDER BY seq_scan DESC;
"

# 3. 检查索引
psql -U alike_user -d alike_db -c "\d users"

# 4. 分析慢查询
# 启用日志记录
# 修改 postgresql.conf:
# log_min_duration_statement = 1000

# 5. 添加数据库索引
CREATE INDEX CONCURRENTLY idx_users_created_at ON users(created_at DESC);

# 6. 使用缓存
# 考虑使用 Redis 缓存频繁查询的数据
```

---

## 🔍 调试技巧

### 启用详细日志

```bash
# Go 应用
export LOG_LEVEL=debug
go run cmd/api/main.go

# PostgreSQL
# 修改 postgresql.conf:
log_min_duration_statement = 1000
log_statement = 'all'

# Docker
docker-compose -f deployments/docker/docker-compose.yml logs -f
```

### 使用开发者工具

```bash
# 前端调试
# 打开 Chrome DevTools (F12)
# - Console: 查看 JavaScript 错误
# - Network: 查看网络请求
# - Application: 查看 LocalStorage 和 Cookies

# 后端调试
# 使用 delve
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug cmd/api/main.go
```

### 性能分析

```bash
# Go pprof
go tool pprof http://localhost:8080/debug/pprof/profile

# 数据库分析
psql -U alike_user -d alike_db -c "EXPLAIN ANALYZE SELECT * FROM users;"
```

---

## 📞 获取帮助

如果以上方案都无法解决问题：

1. **查看日志**：
   ```bash
   # 应用日志
   tail -f logs/app.log
   
   # Docker 日志
   docker-compose -f deployments/docker/docker-compose.yml logs
   ```

2. **搜索已知问题**：
   - GitHub Issues
   - 项目文档

3. **提交问题**：
   - 提供详细的错误信息
   - 说明环境配置
   - 附上复现步骤

4. **联系维护者**：
   - 通过 GitHub Issues
   - 或其他联系方式

---

**最后更新**: 2026-03-03  
**版本**: 1.1.0

---

## � 下一步建议

1. **修复 Docker Desktop**：以便能够运行完整的开发环境
2. **安装 PostgreSQL**：本地数据库支持
3. **启动后端服务**：实现完整的登录和聊天功能
4. **优化其他页面**：global-chat.html 和 index.html

---

## 📞 联系支持

如果遇到问题，请查看：
- 项目文档：`docs/` 目录
- API 文档：`docs/api/` 目录
- 或提交 Issue

---

**最后更新**: 2026-03-03
**版本**: 1.0.0