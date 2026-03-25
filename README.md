<div align="center">

# Alike

### 🌈 相似灵魂的相遇

**一个专为 LGBTQ+ 社区打造的现代化社交应用**

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-3.5-4FC08D?style=flat&logo=vue.js&logoColor=white)](https://vuejs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)

</div>

---

## 项目简介

**Alike** 是一个基于位置的社交匹配平台，致力于帮助 LGBTQ+ 社区的用户发现附近的志同道合的朋友。

---

## 快速开始

### 后端服务

```bash
# 1. 启动 PostgreSQL
brew services start postgresql@15

# 2. 运行数据库迁移
go run cmd/migrate/main.go up

# 3. 启动 API 服务
go run cmd/api/main.go
```

API 将运行在：**http://localhost:8080**

### 前端服务

```bash
cd web
npm install
npm run dev
```

前端将运行在：**http://localhost:3000**

---

## 技术栈

### 后端
- **语言**: Go 1.23+
- **框架**: Gin
- **数据库**: PostgreSQL 14+
- **认证**: JWT

### 前端
- **框架**: Vue 3
- **构建工具**: Vite
- **路由**: Vue Router
- **状态管理**: Pinia
- **HTTP 客户端**: Axios

---

## 项目结构

```
Alike/
├── cmd/                    # 应用入口
│   ├── api/               # API 服务
│   └── migrate/           # 数据库迁移
├── internal/              # 私有应用代码
│   ├── api/              # API 处理器
│   ├── domain/           # 领域模型
│   └── ...
├── web/                   # Vue 3 前端
│   ├── src/
│   │   ├── components/   # 组件
│   │   ├── views/        # 页面
│   │   ├── router/       # 路由
│   │   └── stores/       # 状态管理
│   └── package.json
├── db/                    # 数据库相关
│   ├── migrations/       # 迁移文件
│   └── seeds/            # 种子数据
└── docs/                  # 文档
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

### 全局聊天
- `GET /api/v1/global/messages` - 获取消息列表
- `POST /api/v1/global/messages` - 发送消息

---

## 文档

- [本地开发指南](QUICKSTART-LOCAL.md)
- [项目路线图](docs/ROADMAP.md)

---

## 许可证

Copyright © 2026 Alike. All rights reserved.

---

<div align="center">

**让每一次相遇都成为可能** ❤️

</div>
