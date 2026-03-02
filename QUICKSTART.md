# ⚡ Alike - 快速开始（3分钟）

## 🚀 最简单的方式

### 方法1: 一键启动（推荐）

```bash
cd /Users/zhenghongfei6/go/src/github.com/Alike

# 一键启动所有服务
./scripts/start-all.sh
```

这会自动启动：
- ✅ PostgreSQL数据库
- ✅ Redis缓存
- ✅ API服务器
- ✅ Nginx Web服务器

然后访问：**http://localhost**

---

### 方法2: 分步启动

#### 步骤1: 设置数据库

```bash
cd /Users/zhenghongfei6/go/src/github.com/Alike
./scripts/setup-db.sh
```

#### 步骤2: 启动API服务器

```bash
# 新终端窗口1
cd /Users/zhenghongfei6/go/src/github.com/Alike
go run cmd/api/main.go
```

#### 步骤3: 启动Web服务器

```bash
# 新终端窗口2
cd /Users/zhenghongfei6/go/src/github.com/Alike/web/public
python3 -m http.server 8000
```

#### 步骤4: 打开浏览器

访问：**http://localhost:8000**

---

### 方法3: 仅启动API（无数据库）

```bash
cd /Users/zhenghongfei6/go/src/github.com/Alike
go run cmd/api/main.go
```

然后在新终端：
```bash
cd web/public
python3 -m http.server 8000
```

访问：**http://localhost:8000**

---

## 📱 默认登录信息

**测试账号：**
- 手机号：`+8613800138000`
- 密码：`password123`

**或注册新账号：**
- 验证码填写任意值，如 `123456`

---

## 🔧 服务端口

- **Web界面**: http://localhost:80 或 http://localhost:8000
- **API服务器**: http://localhost:8080
- **API健康检查**: http://localhost:8080/health

---

## 🛑 停止服务

### 方法1: Ctrl+C
在启动终端按 `Ctrl+C`

### 方法2: 手动停止
```bash
# 停止Docker服务
docker-compose -f deployments/docker/docker-compose.yml down

# 或杀死占用端口的进程
lsof -ti:8080 | xargs kill -9
lsof -ti:8000 | xargs kill -9
```

---

## 📚 详细文档

- **完整使用指南**: 见项目根目录
- **API文档**: `API.md`
- **数据库设置**: 见项目文档

---

## ⚡ 快速测试

```bash
# 1. 启动所有服务
cd /Users/zhenghongfei6/go/src/github.com/Alike
./scripts/start-all.sh

# 2. 打开浏览器
open http://localhost

# 3. 使用默认账号登录
# 手机: +8613800138000
# 密码: password123
```

---

**推荐**: 使用 `./scripts/start-all.sh` 一键启动所有服务！

*最后更新: 2026-03-03*
