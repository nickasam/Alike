#!/bin/bash

echo "🛑 停止所有Alike服务..."

# 停止API服务器
if [ -f /tmp/alike-api.pid ]; then
    API_PID=$(cat /tmp/alike-api.pid)
    if kill -0 $API_PID 2>/dev/null; then
        kill $API_PID
        echo "✅ API服务器已停止 (PID: $API_PID)"
    fi
    rm -f /tmp/alike-api.pid
fi

# 停止Web服务器
if [ -f /tmp/alike-web.pid ]; then
    WEB_PID=$(cat /tmp/alike-web.pid)
    if kill -0 $WEB_PID 2>/dev/null; then
        kill $WEB_PID
        echo "✅ Web服务器已停止 (PID: $WEB_PID)"
    fi
    rm -f /tmp/alike-web.pid
fi

# 清理端口
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
lsof -ti:8002 | xargs kill -9 2>/dev/null || true

# 停止Docker PostgreSQL（如果使用Docker）
if docker ps | grep -q alike-postgres; then
    docker stop alike-postgres
    docker rm alike-postgres
    echo "✅ PostgreSQL容器已停止"
fi

echo ""
echo "✅ 所有服务已停止"
