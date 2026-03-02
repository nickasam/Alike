#!/bin/bash

echo "🚀 Starting Alike Full Stack..."

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "🛑 Stopping all services..."
    
    # Stop docker-compose if running
    if docker-compose -f deployments/docker/docker-compose.yml ps | grep -q "Up"; then
        echo "Stopping Docker services..."
        docker-compose -f deployments/docker/docker-compose.yml down
    fi
    
    # Kill any processes on port 8080 and 8000
    lsof -ti:8080 | xargs kill -9 2>/dev/null || true
    lsof -ti:8000 | xargs kill -9 2>/dev/null || true
    
    echo "✅ All services stopped"
    exit 0
}

# Trap Ctrl+C
trap cleanup INT TERM

# Start services with Docker Compose
echo "📦 Starting services with Docker Compose..."
docker-compose -f deployments/docker/docker-compose.yml up -d

echo ""
echo "⏳ Waiting for services to start..."
sleep 10

# Check services
echo "🔍 Checking services..."

# Check API
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ API server running on http://localhost:8080"
else
    echo "⚠️  API server not ready yet"
fi

# Check Nginx
if curl -s http://localhost > /dev/null; then
    echo "✅ Web server running on http://localhost"
else
    echo "⚠️  Web server not ready yet"
fi

echo ""
echo "🎉 All services started!"
echo ""
echo "📱 Open browser at: http://localhost"
echo "🔧 API docs: http://localhost:8080/health"
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

# Keep script running
wait
