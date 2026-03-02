# Alike Web Application

这是一个简单的Web前端，可以让你在浏览器中使用Alike应用。

## 功能

- ✅ 用户注册/登录
- ✅ 查看附近用户
- ✅ 喜欢用户
- ✅ 查看匹配
- ✅ 查看聊天
- ✅ 响应式设计（手机/平板/桌面）

## 使用方法

### 开发模式

1. 启动后端API：
```bash
cd /Users/zhenghongfei6/go/src/github.com/Alike
go run cmd/api/main.go
```

2. 打开Web应用：
- 直接在浏览器中打开 `web/public/index.html`
- 或者使用HTTP服务器：
```bash
cd web/public
python3 -m http.server 8000
```
- 然后访问 http://localhost:8000

### 生产部署

使用Nginx提供静态文件服务：

```nginx
server {
    listen 80;
    server_name alike.app;

    # Web应用
    location / {
        root /path/to/Alike/web/public;
        try_files $uri $uri/ /index.html;
    }

    # API代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 测试账号

由于使用了测试数据，你可以使用以下账号登录：

- 手机号: +8613800138000
- 密码: password123

或者注册新账号（验证码填写任意值，如：123456）

## 技术栈

- 纯HTML + CSS + JavaScript
- 无框架依赖
- 响应式设计
- RESTful API集成

## 功能演示

### 1. 注册/登录
- 填写手机号、密码等信息
- 验证码可填写任意值（开发模式）

### 2. 查看附近用户
- 基于地理位置推荐
- 显示用户头像、昵称、简介
- 可以"喜欢"或"打招呼"

### 3. 匹配功能
- 双向喜欢自动匹配
- 查看所有匹配

### 4. 聊天功能
- 查看聊天列表
- 发送消息（API已实现）

## 下一步

- [ ] 实时聊天（WebSocket）
- [ ] 图片上传
- [ ] 个人资料编辑
- [ ] 消息推送
- [ ] 更丰富的UI交互

---

**提示**: 这是Web版本的MVP，主要功能完整但界面简洁。移动端体验会更佳！
