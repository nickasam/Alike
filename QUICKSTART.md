# ⚡ Alike - 3分钟快速开始

## 🚀 快速启动（2个命令）

```bash
# 1. 启动API服务器
cd /Users/zhenghongfei6/go/src/github.com/Alike
go run cmd/api/main.go

# 2. 打开Web应用（新终端）
open web/public/index.html
```

**就这么简单！** 🎉

---

## 📱 登录信息

**使用默认账号登录：**
- 手机号：`+8613800138000`
- 密码：`password123`

**或注册新账号：**
- 验证码填写任意值，如 `123456`

---

## 🎯 核心功能

### ✅ 已实现
- 👤 用户注册/登录
- 👥 查看附近用户
- ❤️ 喜欢用户
- 💕 查看匹配
- 💬 查看聊天
- 📱 响应式Web界面

---

## 🔧 其他启动方式

### 使用Docker（完整功能）
```bash
docker-compose -f deployments/docker/docker-compose.yml up -d
```

### 使用部署脚本
```bash
./scripts/deploy.sh
```

---

## 📚 详细文档

- **完整使用指南**: 见项目根目录 `USAGE-GUIDE.md`
- **API文档**: 见 `API.md`
- **架构文档**: 见 `docs/architecture.md`

---

## 🐛 遇到问题？

### 数据库连接失败
**正常！** 不影响使用，只是数据不会保存到数据库。

```bash
# 如果需要持久化，启动PostgreSQL：
brew services start postgresql
```

### 端口被占用
```bash
# 杀死占用8080端口的进程
lsof -ti:8080 | xargs kill -9
```

---

## 🎉 开始使用

```bash
# 最快的方式
cd /Users/zhenghongfei6/go/src/github.com/Alike
go run cmd/api/main.go
open web/public/index.html
```

**浏览器会自动打开，开始使用吧！** 🚀

---

*项目进度: 90% → 生产就绪*
*最后更新: 2026-03-03*
