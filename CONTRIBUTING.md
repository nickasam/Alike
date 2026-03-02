# Contributing to Alike

感谢您对 Alike 项目的关注！我们欢迎任何形式的贡献。

## 如何贡献

### 报告 Bug
如果您发现了 bug，请创建 Issue 并提供：
- 清晰的标题
- 问题描述
- 复现步骤
- 环境信息

### 提交代码

1. Fork 项目
2. 创建功能分支：`git checkout -b feature/your-feature`
3. 提交代码：`git commit -m "feat: add your feature"`
4. 推送到 GitHub：`git push origin feature/your-feature`
5. 创建 Pull Request

### Commit 规范
遵循 Conventional Commits：
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `refactor`: 重构

## 开发环境

### 要求
- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### 设置
\`\`\`bash
go mod download
cp config/config.yaml.example config/config.yaml
# 编辑配置文件
make migrate-up
make run
\`\`\`

---

**再次感谢您的贡献！** ❤️
