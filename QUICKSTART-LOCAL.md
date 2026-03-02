# 🚀 Alike 本地调试指南

## 快速开始（3分钟）

### 方式1: 自动设置（推荐）

```bash
cd /Users/zhenghongfei6/go/src/github.com/Alike
bash scripts/setup-local-env.sh
```

**这会自动：**
1. ✅ 检查Go环境
2. ✅ 检查Python环境
3. ✅ 安装/检查PostgreSQL
4. ✅ 测试数据库连接
5. ✅ 运行数据库迁移
6. ✅ 导入测试数据
7. ✅ 启动API服务器
8. ✅ 启动Web服务器

---

## 手动设置

### 步骤1: 安装PostgreSQL

```bash
# 使用Homebrew安装
brew install postgresql@15

# 启动服务
brew services start postgresql@15

# 验证安装
psql --version
```

### 步骤2: 创建数据库

```bash
# 创建数据库
createdb alike_db

# 创建用户
psql -d postgres -c "CREATE USER alike_user WITH PASSWORD 'alike_password';"
psql -d postgres -c "GRANT ALL PRIVILEGES ON DATABASE alike_db TO alike_user;"

# 测试连接
PGPASSWORD=alike_password psql -U alike_user -d alike_db -c "SELECT 1;"
```

### 步骤3: 运行迁移

```bash
cd /Users/zhenghongfei6/go/src/github.com/Alike
go run cmd/migrate/main.go up
```

### 步骤4: 导入测试数据

```bash
PGPASSWORD=alike_password psql -U alike_user -d alike_db -f db/seeds/seed.sql
```

### 步骤5: 启动服务

```bash
# 终端1: 启动API
go run cmd/api/main.go

# 终端2: 启动Web
cd web/public && python3 -m http.server 8002
```

---

## 运行测试

### 自动测试所有功能

```bash
bash scripts/verify-all.sh
```

**测试内容：**
1. ✅ API健康检查
2. ✅ 用户注册
3. ✅ 用户登录
4. ✅ 获取用户信息
5. ✅ 获取附近用户
6. ✅ 获取匹配列表
7. ✅ 发送喜欢
8. ✅ 获取聊天列表
9. ✅ 获取全局聊天室
10. ✅ 获取全局消息

---

## 访问应用

### 启动页
```
http://localhost:8002/launcher.html
```

### 全局聊天室
```
http://localhost:8002/global-chat.html
```

### 附近用户
```
http://localhost:8002/index.html
```

---

## 测试账号

```
手机: +8613800138000
密码: password123
```

---

## 停止服务

### 停止所有服务

```bash
bash scripts/stop-all.sh
```

### 或手动停止

```bash
# 停止API
kill $(cat /tmp/alike-api.pid)

# 停止Web
kill $(cat /tmp/alike-web.pid)

# 停止PostgreSQL
brew services stop postgresql@15
```

---

## 故障排除

### 问题1: PostgreSQL连接失败

```bash
# 检查PostgreSQL状态
brew services list | grep postgres

# 启动服务
brew services start postgresql@15

# 检查连接
PGPASSWORD=alike_password psql -U alike_user -d alike_db
```

### 问题2: 端口被占用

```bash
# 查看占用
lsof -ti:8080
lsof -ti:8002

# 杀死进程
lsof -ti:8080 | xargs kill -9
lsof -ti:8002 | xargs kill -9
```

### 问题3: 数据库迁移失败

```bash
# 手动运行迁移
cd /Users/zhenghongfei6/go/src/github.com/Alike
go run cmd/migrate/main.go up

# 查看错误
PGPASSWORD=alike_password psql -U alike_user -d alike_db
\dt  # 查看表
```

---

## 日志文件

```bash
# API日志
cat /tmp/alike-api.log

# Web日志
cat /tmp/alike-web.log

# 实时查看
tail -f /tmp/alike-api.log
```

---

## 开发工作流

### 典型开发流程

```bash
# 1. 启动环境
bash scripts/setup-local-env.sh

# 2. 修改代码
vim internal/...

# 3. 重启API
kill $(cat /tmp/alike-api.pid)
go run cmd/api/main.go > /tmp/alike-api.log 2>&1 &
echo $! > /tmp/alike-api.pid

# 4. 测试功能
bash scripts/verify-all.sh

# 5. 访问应用
open http://localhost:8002/launcher.html
```

---

## 🔗 API端点

### 认证
- `POST /api/v1/auth/register` - 注册
- `POST /api/v1/auth/login` - 登录
- `GET /api/v1/auth/me` - 获取当前用户

### 用户
- `GET /api/v1/users/me` - 获取我的信息
- `PUT /api/v1/users/me` - 更新我的信息
- `GET /api/v1/users/nearby` - 获取附近用户

### 匹配
- `GET /api/v1/matches` - 获取匹配列表
- `POST /api/v1/matches/:id/like` - 发送喜欢

### 聊天
- `GET /api/v1/chats` - 获取聊天列表
- `POST /api/v1/chats/:id/messages` - 发送消息

### 全局聊天
- `GET /api/v1/global/room` - 获取聊天室信息
- `GET /api/v1/global/messages` - 获取消息列表
- `POST /api/v1/global/messages` - 发送消息

---

**脚本版本**: v1.0  
**最后更新**: 2026-03-03  
**GitHub**: https://github.com/nickasam/Alike
