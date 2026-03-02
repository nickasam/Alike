# Alike

一个专为男同性恋社区打造的社交应用。

## 项目简介

Alike 是一个基于 Go 语言开发的 LGBTQ+ 社交平台，旨在为男同性恋群体提供一个安全、真实、友好的交友环境。

## 功能特性

### 核心功能
- 🔐 **安全认证** - 手机号/邮箱注册，JWT 认证
- 👤 **个人资料** - 完善的用户信息，照片墙，标签系统
- 📍 **附近的人** - 基于地理位置的推荐算法
- 💬 **即时聊天** - WebSocket 实时通信
- 🔔 **推送通知** - 支持 APNs 和 FCM
- 🏷️ **兴趣匹配** - 基于标签和偏好的智能推荐

### 安全与隐私
- 🔒 **隐私控制** - 屏蔽、举报、隐身模式
- ✅ **身份认证** - 可选真人认证，提高安全性
- 🛡️ **内容审核** - 图片和文字内容过滤
- 🔐 **数据加密** - 传输和存储加密

## 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: PostgreSQL 15+
- **缓存**: Redis 7+
- **实时通信**: WebSocket
- **ORM**: GORM
- **验证**: go-playground/validator

### 基础设施
- **容器**: Docker
- **编排**: Docker Compose
- **反向代理**: Nginx
- **监控**: Prometheus + Grafana
- **CI/CD**: GitHub Actions

## 项目结构

```
Alike/
├── cmd/              # 应用入口
├── internal/         # 私有应用代码
│   ├── api/         # HTTP handlers
│   ├── auth/        # 认证服务
│   ├── chat/        # 聊天服务
│   ├── match/       # 匹配服务
│   ├── user/        # 用户服务
│   ├── notification/# 推送通知
│   └── domain/      # 领域模型
├── pkg/             # 公共库
│   ├── database/    # 数据库连接
│   ├── config/      # 配置管理
│   └── utils/       # 工具函数
├── docs/            # 文档
└── scripts/         # 脚本
```

## 快速开始

### 环境要求
- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker (可选)

### 安装

```bash
# 克隆项目
git clone https://github.com/Alike/alike.git
cd alike

# 安装依赖
go mod download

# 复制配置文件
cp config/config.yaml.example config/config.yaml

# 编辑配置文件
vim config/config.yaml

# 运行数据库迁移
go run cmd/migrate/main.go up

# 启动服务
go run cmd/api/main.go
```

### 使用 Docker

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

## API 文档

API 文档请查看 [docs/api/README.md](docs/api/README.md)

主要端点：
- 认证: `/api/v1/auth/*`
- 用户: `/api/v1/users/*`
- 匹配: `/api/v1/matches/*`
- 聊天: `/api/v1/chats/*`
- WebSocket: `/api/v1/chats/:id/ws`

## 开发路线图

### MVP (v0.1.0)
- [x] 项目初始化
- [x] UI/UX 设计
- [ ] 数据库模型
- [ ] 认证模块
- [ ] 用户模块
- [ ] 匹配算法
- [ ] 聊天功能
- [ ] 推送通知

### v0.2.0
- [ ] 群组功能
- [ ] 动态/朋友圈
- [ ] 视频通话
- [ ] 消息加密

### v0.3.0
- [ ] 会员系统
- [ ] 高级筛选
- [ ] 数据分析
- [ ] 性能优化

## 贡献指南

我们欢迎任何形式的贡献！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 行为准则

- 尊重每个人
- 包容多样性
- 友好交流
- 拒绝歧视

## 许可证

Copyright © 2026 Alike. All rights reserved.

## 联系我们

- 网站: [alike.app](https://alike.app)
- 邮箱: support@alike.app
- Discord: [Alike Community](https://discord.gg/alike)

---

**❤️ 为社区而生，为真实而建。**
