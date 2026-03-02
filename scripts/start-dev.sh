#!/bin/bash

echo "🚀 Alike 本地开发环境启动脚本"
echo "================================"

# 检查端口占用
check_port() {
    local port=$1
    if lsof -i :$port > /dev/null 2>&1; then
        echo "⚠️  端口 $port 已被占用"
        return 1
    fi
    return 0
}

# 查找可用端口
find_free_port() {
    for port in 8000 8001 8002 3000 3001 8080 8081; do
        if ! lsof -i :$port > /dev/null 2>&1; then
            echo $port
            return 0
        fi
    done
    echo "未找到可用端口"
    return 1
}

# 启动前端服务
echo "📱 启动前端服务..."
FREE_PORT=$(find_free_port)

if [ "$FREE_PORT" == "未找到可用端口" ]; then
    echo "❌ 未找到可用端口，请先关闭占用端口的服务"
    exit 1
fi

echo "✅ 使用端口 $FREE_PORT"

cd "$(dirname "$0")/../web/public"
python3 -m http.server $FREE_PORT > /tmp/alike-web.log 2>&1 &
WEB_PID=$!
echo "   前端服务 PID: $WEB_PID"

# 等待服务启动
sleep 2

if lsof -i :$FREE_PORT > /dev/null 2>&1; then
    echo "✅ 前端服务已启动: http://localhost:$FREE_PORT"
    echo ""
    echo "🌐 可用页面:"
    echo "   - 首页/Launcher: http://localhost:$FREE_PORT/launcher.html"
    echo "   - 登录页面: http://localhost:$FREE_PORT/login.html"
    echo "   - 主页: http://localhost:$FREE_PORT/index.html"
    echo "   - 全局聊天: http://localhost:$FREE_PORT/global-chat.html"
    echo "   - 个人资料: http://localhost:$FREE_PORT/profile.html"
    echo ""
    echo "📝 启动后端服务（需要先启动数据库）:"
    echo "   1. 启动 PostgreSQL 和 Redis:"
    echo "      docker-compose -f deployments/docker/docker-compose.yml up -d postgres redis"
    echo ""
    echo "   2. 运行数据库迁移:"
    echo "      make migrate-up"
    echo ""
    echo "   3. 启动 API 服务:"
    echo "      go run cmd/api/main.go"
    echo ""
    echo "按 Ctrl+C 停止前端服务"

    # 捕获退出信号
    cleanup() {
        echo ""
        echo "🛑 停止服务..."
        kill $WEB_PID 2>/dev/null
        exit 0
    }

    trap cleanup INT TERM

    # 在浏览器中打开
    if command -v open > /dev/null; then
        sleep 1
        open http://localhost:$FREE_PORT/launcher.html
    fi

    # 保持运行
    wait
else
    echo "❌ 前端服务启动失败，请查看日志: /tmp/alike-web.log"
    kill $WEB_PID 2>/dev/null
    exit 1
fi