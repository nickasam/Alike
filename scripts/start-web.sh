#!/bin/bash

echo "🌐 Starting Alike Web Application..."

# Check if API server is running
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "⚠️  API server is not running!"
    echo "📦 Starting API server first..."
    cd /Users/zhenghongfei6/go/src/github.com/Alike
    
    # Start API in background
    go run cmd/api/main.go &
    API_PID=$!
    
    echo "⏳ Waiting for API to start..."
    sleep 3
    
    if curl -s http://localhost:8080/health > /dev/null; then
        echo "✅ API server started (PID: $API_PID)"
    else
        echo "❌ API server failed to start"
        exit 1
    fi
fi

# Start web server
cd /Users/zhenghongfei6/go/src/github.com/Alike/web/public

echo "🚀 Starting web server on http://localhost:8000"
echo "📱 Open browser at: http://localhost:8000"
echo ""
echo "Press Ctrl+C to stop both servers"

# Trap Ctrl+C to kill both processes
trap "echo 'Stopping servers...'; kill $API_PID 2>/dev/null; exit" INT TERM

# Start Python HTTP server
python3 -m http.server 8000 &
WEB_PID=$!

echo ""
echo "✅ Both servers running!"
echo "   API: http://localhost:8080"
echo "   Web: http://localhost:8000"
echo ""

# Wait for processes
wait $API_PID $WEB_PID
