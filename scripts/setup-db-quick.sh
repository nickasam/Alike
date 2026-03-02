#!/bin/bash

echo "🚀 Quick Database Setup"

# 使用 Docker 快速启动
if command -v docker &> /dev/null; then
    echo "📦 Using Docker to start PostgreSQL..."
    
    docker run -d \
        --name alike-postgres \
        -e POSTGRES_DB=alike_db \
        -e POSTGRES_USER=alike_user \
        -e POSTGRES_PASSWORD=alike_password \
        -p 5432:5432 \
        postgres:15-alpine
    
    echo "⏳ Waiting for PostgreSQL to start..."
    sleep 8
    
    # 检查是否成功
    if docker ps | grep -q alike-postgres; then
        echo "✅ PostgreSQL started successfully"
        
        # 运行迁移
        echo "📋 Running migrations..."
        cd /Users/zhenghongfei6/go/src/github.com/Alike
        go run cmd/migrate/main.go up
        
        # 导入测试数据
        echo "🌱 Seeding test data..."
        PGPASSWORD=alike_password psql -h localhost -U alike_user -d alike_db -f db/seeds/seed.sql
        
        echo ""
        echo "✅ Database setup complete!"
        echo "📊 Test credentials:"
        echo "   Phone: +8613800138000"
        echo "   Password: password123"
        echo ""
        echo "🚀 Now restart the API:"
        echo "   kill \$(cat /tmp/alike-api.pid)"
        echo "   go run cmd/api/main.go > /tmp/alike-api.log 2>&1 &"
    else
        echo "❌ Failed to start PostgreSQL"
    fi
else
    echo "❌ Docker not found"
    echo "Please install Docker or use local PostgreSQL"
fi
