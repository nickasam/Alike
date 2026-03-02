#!/bin/bash

echo "🚀 Launching Alike..."

# 检查API是否运行
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "📦 Starting API server..."
    cd /Users/zhenghongfei6/go/src/github.com/Alike
    go run cmd/api/main.go > /tmp/alike-api.log 2>&1 &
    API_PID=$!
    echo $API_PID > /tmp/alike-api.pid
    
    echo "⏳ Waiting for API to start..."
    sleep 3
    
    if curl -s http://localhost:8080/health > /dev/null; then
        echo "✅ API server started"
    else
        echo "❌ API server failed to start"
        cat /tmp/alike-api.log
        exit 1
    fi
fi

# 启动Web服务器
echo "🌐 Starting web server..."
cd /Users/zhenghongfei6/go/src/github.com/Alike/web/public
python3 -m http.server 8002 > /tmp/alike-web.log 2>&1 &
WEB_PID=$!
echo $WEB_PID > /tmp/alike-web.pid

echo ""
echo "✅ All servers started!"
echo ""
echo "📱 Open launcher at: http://localhost:8002/launcher.html"
echo "💬 Or global chat: http://localhost:8002/global-chat.html"
echo ""
echo "Press Ctrl+C to stop"

# Trap Ctrl+C
trap cleanup INT TERM

cleanup() {
    echo ""
    echo "🛑 Stopping servers..."
    kill $API_PID $WEB_PID 2>/dev/null
    exit 0
}

wait
