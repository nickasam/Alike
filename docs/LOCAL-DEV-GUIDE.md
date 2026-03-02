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

## ⚠️ 常见问题

### 1. 端口被占用
```bash
# 查看占用端口的进程
lsof -i :8000

# 杀死进程
kill -9 <PID>
```

### 2. Docker 无法启动
- 方案一：重启 Docker Desktop
- 方案二：使用本地 PostgreSQL
- 方案三：暂时只运行前端

### 3. Python http.server 无法启动
- 确认已安装 Python 3：`python3 --version`
- 检查防火墙设置
- 尝试其他端口

---

## 📝 下一步建议

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