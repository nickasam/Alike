# Alike

一个专为男同性恋社区打造的社交应用。

## 🚀 快速开始

### 方法1: Docker（推荐 - 最简单）

```bash
cd /Users/zhenghongfei6/go/src/github.com/Alike
docker-compose -f deployments/docker/docker-compose.yml up -d
```

等待1-2分钟，然后访问：**http://localhost**

---

### 方法2: 本地开发

#### 步骤1: 安装PostgreSQL

```bash
# 安装
brew install postgresql@15

# 启动服务
brew services start postgresql@15

# 验证
psql --version
```

#### 步骤2: 设置数据库

```bash
cd /Users/zhenghongfei6/go/src/github.com/Alike
./scripts/setup-db.sh
```

#### 步骤3: 启动应用

```bash
./scripts/start.sh
```

访问：**http://localhost:8000**

---

## 🔑 默认账号

```
手机号: +8613800138000
密码: password123
```

---

## 📱 功能特性

- ✅ 用户注册/登录
- ✅ 查看附近用户
- ✅ 喜欢匹配功能
- ✅ 聊天功能
- ✅ 现代化Web界面
- ✅ RESTful API

---

## 🛠️ 开发

```bash
# 启动API
go run cmd/api/main.go

# 运行测试
make test

# 构建项目
make build
```

---

## 📚 技术栈

- **后端**: Go 1.23 + Gin + PostgreSQL + Redis
- **前端**: HTML + CSS + JavaScript
- **部署**: Docker + Docker Compose

---

## 📄 文档

- [架构设计](docs/architecture.md)
- [UI设计](docs/design.md)

---

## 📄 许可证

Copyright © 2026 Alike. All rights reserved.
