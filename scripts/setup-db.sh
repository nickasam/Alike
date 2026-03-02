#!/bin/bash

set -e

echo "🚀 Alike Database Setup"
echo "======================="
echo ""

# 检查PostgreSQL
if ! command -v psql &> /dev/null; then
    echo "❌ PostgreSQL 未安装"
    echo ""
    echo "请先安装 PostgreSQL："
    echo "  brew install postgresql@15"
    echo "  brew services start postgresql@15"
    echo ""
    echo "或者使用 Docker（推荐）："
    echo "  cd /Users/zhenghongfei6/go/src/github.com/Alike"
    echo "  docker-compose -f deployments/docker/docker-compose.yml up -d"
    echo ""
    exit 1
fi

# 检查服务状态
if ! pg_isready -q 2>/dev/null; then
    echo "🔧 启动 PostgreSQL 服务..."
    brew services start postgresql@15
    echo "⏳ 等待服务启动..."
    sleep 5
fi

# 再次检查
if ! pg_isready -q 2>/dev/null; then
    echo "❌ PostgreSQL 服务未启动"
    echo ""
    echo "请手动启动："
    echo "  brew services start postgresql@15"
    echo ""
    echo "检查状态："
    echo "  brew services list | grep postgres"
    echo ""
    exit 1
fi

echo "✅ PostgreSQL 服务运行中"
echo ""

# 创建数据库
echo "📊 创建数据库..."
if createdb alike_db 2>/dev/null; then
    echo "✅ 数据库创建成功: alike_db"
else
    echo "ℹ️  数据库已存在: alike_db"
fi
echo ""

# 创建用户
echo "👤 创建数据库用户..."
psql -d postgres >/dev/null 2>&1 << SQL || true
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_user WHERE usename = 'alike_user') THEN
        CREATE USER alike_user WITH PASSWORD 'alike_password';
    END IF;
    GRANT ALL PRIVILEGES ON DATABASE alike_db TO alike_user;
END
\$\$;
SQL

echo "✅ 用户创建成功: alike_user"
echo ""

# 测试连接
echo "🔗 测试连接..."
if PGPASSWORD=alike_password psql -U alike_user -d alike_db -c "SELECT 1;" >/dev/null 2>&1; then
    echo "✅ 数据库连接成功"
else
    echo "❌ 数据库连接失败"
    echo ""
    echo "请检查："
    echo "  PGPASSWORD=alike_password psql -U alike_user -d alike_db"
    echo ""
    exit 1
fi
echo ""

# 运行迁移
echo "📋 运行数据库迁移..."
cd /Users/zhenghongfei6/go/src/github.com/Alike
if go run cmd/migrate/main.go up 2>&1 | grep -q "Migration completed"; then
    echo "✅ 迁移完成"
else
    echo "⚠️  迁移可能失败，请检查输出"
fi
echo ""

# 导入测试数据
echo "🌱 导入测试数据..."
if PGPASSWORD=alike_password psql -U alike_user -d alike_db -f db/seeds/seed.sql >/dev/null 2>&1; then
    echo "✅ 测试数据导入成功"
    echo ""
    echo "测试账号："
    echo "  手机: +8613800138000"
    echo "  密码: password123"
else
    echo "ℹ️  测试数据已存在"
fi
echo ""

echo "✅ 数据库设置完成！"
echo ""
echo "📊 连接信息："
echo "   Host: localhost"
echo "   Port: 5432"
echo "   Database: alike_db"
echo "   User: alike_user"
echo "   Password: alike_password"
echo ""
echo "🧪 测试连接："
echo "   PGPASSWORD=alike_password psql -U alike_user -d alike_db"
echo ""
echo "🚀 启动应用："
echo "   ./scripts/start.sh"
echo ""
