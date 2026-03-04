<div align="center">

# Alike

### 🌈 相似灵魂的相遇

**一个专为 LGBTQ+ 社区打造的现代化社交应用**

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/license-Proprietary-red.svg)](LICENSE)

[功能特性](#-功能特性) • [快速开始](#-快速开始) • [技术栈](#-技术栈) • [开发指南](#-开发指南) • [项目路线图](#-项目路线图)

</div>

---

## ✨ 项目简介

**Alike** 是一个基于位置的社交匹配平台，致力于帮助 LGBTQ+ 社区的用户发现附近的志同道合的朋友。我们相信，每个人都值得被看见，每段关系都值得被珍惜。

### 🎯 核心价值

- **真实性**: 真实的用户，真实的连接
- **安全性**: 隐私保护，安全交友
- **包容性**: 为 LGBTQ+ 社区打造
- **便捷性**: 基于位置的智能匹配

---

## 🎨 设计理念

### Material Design 设计系统
- ✨ 统一的 **Material Icons** 图标风格
- 🎨 精美的 **紫粉渐变** 配色方案
- 📱 **完全响应式** 设计（桌面端 + 移动端）
- 💫 流畅的 **动画和过渡效果**

### 页面展示
- **启动页** (`launcher.html`): 现代化的 Hero 区域，精美的功能卡片
- **登录页** (`login.html`): 全屏布局，优雅的表单设计
- **主页** (`index.html`): 用户匹配和附近用户发现
- **全局聊天** (`global-chat.html`): 实时社交互动
- **个人资料** (`profile.html`): 完善的用户信息展示

---

## 🚀 快速开始

### 方式 1: 使用开发启动脚本（推荐）

```bash
# 克隆项目
git clone https://github.com/nickasam/Alike.git
cd Alike

# 使用开发脚本启动
./scripts/start-dev.sh
```

启动脚本会：
- 🔍 自动查找可用端口（8000-8081）
- 🚀 启动前端服务
- 🌐 在浏览器中自动打开页面
- 📋 提供后端启动指南

### 方式 2: Docker Compose

```bash
# 启动所有服务（前端 + 后端 + 数据库）
docker-compose -f deployments/docker/docker-compose.yml up -d

# 访问应用
open http://localhost
```

### 方式 3: 手动启动

#### 前端服务
```bash
cd web/public
python3 -m http.server 8000
```

访问：**http://localhost:8000/launcher.html**

#### 后端服务
```bash
# 1. 启动 PostgreSQL
brew services start postgresql@14

# 2. 创建数据库
createdb alike_db

# 3. 运行数据库迁移
make migrate-up

# 4. 启动 API 服务
go run cmd/api/main.go
```

API 将运行在：**http://localhost:8080**

---

## 📱 功能特性

### ✅ 已实现功能

#### 用户系统
- [x] 用户注册/登录（手机号 + 验证码）
- [x] JWT Token 认证
- [x] 个人资料管理
- [x] 头像上传
- [x] 隐私设置

#### 社交功能
- [x] **全局聊天室**: 实时公共聊天，结识新朋友
- [x] **附近用户**: 基于位置的智能推荐
- [x] **喜欢匹配**: 双向喜欢即可开始聊天
- [x] **实时消息**: WebSocket 长连接

#### 界面设计
- [x] Material Design 设计语言
- [x] 完全响应式布局
- [x] 流畅的动画效果
- [x] 无障碍访问支持

### 🚧 开发中功能

- [ ] 一对一私聊
- [ ] 图片/语音消息
- [ ] 消息推送通知
- [ ] 用户关注系统
- [ ] 动态发布功能
- [ ] 内容审核系统

---

## 🛠️ 技术栈

### 后端
- **语言**: Go 1.23+
- **框架**: Gin (Web 框架)
- **数据库**: PostgreSQL 14+
- **缓存**: Redis 7+
- **认证**: JWT (JSON Web Tokens)
- **实时通信**: WebSocket

### 前端
- **技术**: HTML5 + CSS3 + JavaScript (ES6+)
- **设计系统**: Material Design
- **图标**: Material Icons
- **字体**: Inter (Google Fonts)

### 基础设施
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx
- **部署**: Kubernetes 准备就绪
- **数据库迁移**: golang-migrate

---

## 📚 项目结构

```
Alike/
├── cmd/                    # 应用入口
│   ├── api/               # API 服务
│   ├── migrate/           # 数据库迁移工具
│   └── worker/            # 后台任务处理器
├── internal/              # 私有应用代码
│   ├── api/              # API 处理器和中间件
│   ├── auth/             # 认证服务
│   ├── chat/             # 聊天服务
│   ├── domain/           # 领域模型
│   ├── match/            # 匹配服务
│   ├── notification/     # 通知服务
│   ├── repository/       # 数据访问层
│   └── user/             # 用户服务
├── web/                   # 前端资源
│   └── public/           # 静态文件
│       ├── launcher.html # 启动页 ✨
│       ├── login.html    # 登录页 ✨
│       ├── index.html    # 主页
│       ├── global-chat.html # 全局聊天
│       └── profile.html  # 个人资料
├── db/                    # 数据库相关
│   ├── migrations/       # 数据库迁移文件
│   └── seeds/            # 种子数据
├── deployments/           # 部署配置
│   ├── docker/           # Docker 配置
│   └── k8s/              # Kubernetes 配置
├── docs/                  # 项目文档
│   ├── ROADMAP.md        # 项目路线图
│   ├── LOCAL-DEV-GUIDE.md # 本地开发指南
│   ├── architecture.md   # 架构设计
│   └── design.md         # UI 设计文档
├── scripts/               # 实用脚本
│   ├── start-dev.sh      # 开发启动脚本 ✨
│   ├── start.sh          # 生产启动脚本
│   └── setup-db.sh       # 数据库设置脚本
└── config/               # 配置文件
```

---

## 📖 开发指南

### 环境要求

- Go 1.23+
- PostgreSQL 14+
- Redis 7+
- Docker & Docker Compose (可选)

### 本地开发

```bash
# 1. 克隆项目
git clone https://github.com/nickasam/Alike.git
cd Alike

# 2. 安装依赖
go mod download

# 3. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置数据库连接等信息

# 4. 运行数据库迁移
make migrate-up

# 5. 启动开发服务器
make dev

# 或使用开发脚本
./scripts/start-dev.sh
```

### 可用命令

```bash
# 运行 API
make run

# 运行测试
make test

# 构建项目
make build

# 数据库迁移
make migrate-up    # 执行迁移
make migrate-down  # 回滚迁移

# 数据库种子
make seed
```

### API 文档

API 文档请查看：[docs/api/README.md](docs/api/README.md)

主要端点：
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/users/nearby` - 获取附近用户
- `GET /api/v1/global/room` - 全局聊天室
- `POST /api/v1/global/messages` - 发送消息

---

## 🗺️ 项目路线图

### 🎯 Phase 1: 核心功能完善（进行中）
- [ ] SMS 验证码登录
- [ ] 完善用户资料系统
- [ ] GPS 定位集成
- [ ] 智能匹配算法

### 🔥 Phase 2: 社交互动功能
- [ ] WebSocket 实时聊天
- [ ] 一对一私聊
- [ ] 消息推送
- [ ] 关注/好友系统

### 🎨 Phase 3: 用户体验优化
- [ ] 性能优化
- [ ] 移动端适配
- [ ] PWA 支持
- [ ] 安全加固

### 📱 Phase 4: 移动端支持
- [ ] 移动端专用布局
- [ ] 原生应用准备

详细路线图请查看：[docs/ROADMAP.md](docs/ROADMAP.md)

---

## 🤝 贡献指南

我们欢迎各种形式的贡献！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 代码规范

- 遵循 [Effective Go](https://golang.org/doc/effective_go) 指南
- 使用 `gofmt` 格式化代码
- 编写单元测试
- 更新相关文档

---

## 📄 文档

- 📖 [本地开发指南](docs/LOCAL-DEV-GUIDE.md)
- 🗺️ [项目路线图](docs/ROADMAP.md)
- 🏗️ [架构设计](docs/architecture.md)
- 🎨 [UI 设计文档](docs/design.md)
- 🔌 [API 文档](docs/api/README.md)

---

## 🔐 安全

我们非常重视用户安全和隐私：

- ✅ 密码使用 bcrypt 加密存储
- ✅ JWT Token 认证
- ✅ HTTPS 加密传输
- ✅ SQL 注入防护
- ✅ XSS/CSRF 防护
- ✅ API 限流保护

---

## 📞 联系我们

- **Website**: https://alike.app
- **Email**: support@alike.app
- **Issues**: [GitHub Issues](https://github.com/nickasam/Alike/issues)

---

## 📄 许可证

Copyright © 2026 Alike. All rights reserved.

本项目为专有软件，未经授权不得复制、分发或修改。

---

<div align="center">

**让每一次相遇都成为可能** ❤️

Made with ❤️ by the Alike Team

</div>
