# Alike

一个专为男同性恋社区打造的社交应用。

## 快速开始

### 一键启动（推荐）

```bash
# 1. 克隆项目
git clone https://github.com/nickasam/Alike.git
cd Alike

# 2. 一键启动所有服务
./scripts/start.sh
```

访问：**http://localhost**

### 分步启动

```bash
# 1. 设置数据库
./scripts/setup-db.sh

# 2. 启动API
go run cmd/api/main.go

# 3. 启动Web服务器（新终端）
cd web/public && python3 -m http.server 8000
```

访问：**http://localhost:8000**

## 默认账号

- 手机号：`+8613800138000`
- 密码：`password123`

## 功能特性

- ✅ 用户注册/登录
- ✅ 查看附近用户
- ✅ 喜欢匹配
- ✅ 聊天功能
- ✅ Web界面
- ✅ RESTful API

## 技术栈

- **后端**: Go + Gin + PostgreSQL + Redis
- **前端**: HTML + CSS + JavaScript
- **部署**: Docker + Docker Compose

## 文档

- [架构设计](docs/architecture.md)
- [UI设计](docs/design.md)

## 许可证

Copyright © 2026 Alike. All rights reserved.
