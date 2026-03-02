#!/bin/bash

echo "🚀 Starting Alike..."

# Start API server
echo "📦 Starting API server..."
cd /Users/zhenghongfei6/go/src/github.com/Alike
go run cmd/api/main.go &
API_PID=$!

# Wait for API to start
sleep 3

# Check if API is running
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ API server running on http://localhost:8080"
else
    echo "❌ API server failed to start"
    exit 1
fi

# Start web server
echo "🌐 Starting web server..."
cd /Users/zhenghongfei6/go/src/github.com/Alike/web/public
python3 -m http.server 8000 &
WEB_PID=$!

echo ""
echo "✅ Both servers started!"
echo "   Web: http://localhost:8000"
echo "   API: http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop"

# Cleanup function
cleanup() {
    echo ""
    echo "🛑 Stopping servers..."
    kill $API_PID $WEB_PID 2>/dev/null
    exit 0
}

trap cleanup INT TERM

# Wait
wait
