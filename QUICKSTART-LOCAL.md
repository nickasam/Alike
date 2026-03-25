# 🚀 Alike 本地开发指南

## 快速开始

### 环境要求

- Go 1.23+
- Node.js 18+
- PostgreSQL 14+

### 步骤 1: 启动数据库

```bash
# macOS
brew services start postgresql@15

# 或使用 Docker
docker run -d --name alike-postgres \
  -e POSTGRES_DB=alike_db \
  -e POSTGRES_USER=alike_user \
  -e POSTGRES_PASSWORD=alike_password \
  -p 5432:5432 \
  postgres:15
```

### 步骤 2: 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库连接
```

### 步骤 3: 运行数据库迁移

```bash
# 运行迁移
go run cmd/migrate/main.go up

# 导入测试数据（可选）
psql -U alike_user -d alike_db -f db/seeds/seed.sql
```

### 步骤 4: 启动后端服务

```bash
go run cmd/api/main.go
```

API 将运行在：**http://localhost:8080**

### 步骤 5: 启动前端服务

```bash
cd web
npm install
npm run dev
```

前端将运行在：**http://localhost:3000**

---

## 访问应用

- **启动页**: http://localhost:3000/launcher
- **登录页**: http://localhost:3000/login
- **首页**: http://localhost:3000/home
- **全局聊天**: http://localhost:3000/global
- **个人中心**: http://localhost:3000/profile

---

## 测试账号

```
手机: 13800138000
密码: password123
```

---

## API 端点

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

## 故障排除

### 问题1: PostgreSQL 连接失败

```bash
# 检查 PostgreSQL 状态
brew services list | grep postgres

# 启动服务
brew services start postgresql@15

# 检查连接
psql -U alike_user -d alike_db
```

### 问题2: 端口被占用

```bash
# 查看占用
lsof -ti:8080  # 后端端口
lsof -ti:3000  # 前端端口

# 杀死进程
lsof -ti:8080 | xargs kill -9
lsof -ti:3000 | xargs kill -9
```

### 问题3: 前端依赖安装失败

```bash
cd web
rm -rf node_modules package-lock.json
npm install
```

---

**最后更新**: 2026-03-24
