#!/bin/bash

set -e

echo "🚀 Alike 本地调试环境设置"
echo "======================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 步骤1: 检查Go环境
echo "📦 步骤 1/8: 检查Go环境..."
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✅ Go已安装: $GO_VERSION${NC}"
else
    echo -e "${RED}❌ Go未安装${NC}"
    echo "请访问 https://golang.org/dl/ 下载安装"
    exit 1
fi
echo ""

# 步骤2: 检查Python
echo "📦 步骤 2/8: 检查Python环境..."
if command -v python3 &> /dev/null; then
    PYTHON_VERSION=$(python3 --version)
    echo -e "${GREEN}✅ Python已安装: $PYTHON_VERSION${NC}"
else
    echo -e "${RED}❌ Python3未安装${NC}"
    echo "请安装Python3"
    exit 1
fi
echo ""

# 步骤3: 检查/安装PostgreSQL
echo "📦 步骤 3/8: 检查PostgreSQL..."
if command -v psql &> /dev/null; then
    echo -e "${GREEN}✅ PostgreSQL已安装${NC}"
    PG_VERSION=$(psql --version)
    echo "   $PG_VERSION"
else
    echo -e "${YELLOW}⚠️  PostgreSQL未安装${NC}"
    echo "使用Docker启动PostgreSQL..."
    
    if command -v docker &> /dev/null; then
        # 使用Docker启动PostgreSQL
        if docker ps | grep -q alike-postgres; then
            echo "PostgreSQL容器已在运行"
        else
            echo "启动PostgreSQL容器..."
            docker run -d \
                --name alike-postgres \
                -e POSTGRES_DB=alike_db \
                -e POSTGRES_USER=alike_user \
                -e POSTGRES_PASSWORD=alike_password \
                -p 5432:5432 \
                postgres:15-alpine
            
            echo "等待PostgreSQL启动..."
            sleep 8
        fi
        echo -e "${GREEN}✅ PostgreSQL已通过Docker启动${NC}"
    else
        echo -e "${RED}❌ Docker未安装，无法启动PostgreSQL${NC}"
        echo "请安装Docker或手动安装PostgreSQL:"
        echo "  brew install postgresql@15"
        echo "  brew services start postgresql@15"
        exit 1
    fi
fi
echo ""

# 步骤4: 测试数据库连接
echo "🔗 步骤 4/8: 测试数据库连接..."
MAX_ATTEMPTS=5
ATTEMPT=0

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
    if PGPASSWORD=alike_password psql -h localhost -U alike_user -d alike_db -c "SELECT 1;" &> /dev/null; then
        echo -e "${GREEN}✅ 数据库连接成功${NC}"
        break
    else
        ATTEMPT=$((ATTEMPT + 1))
        if [ $ATTEMPT -lt $MAX_ATTEMPTS ]; then
            echo "等待数据库启动... (尝试 $ATTEMPT/$MAX_ATTEMPTS)"
            sleep 3
        fi
    fi
done

if [ $ATTEMPT -eq $MAX_ATTEMPTS ]; then
    echo -e "${RED}❌ 数据库连接失败${NC}"
    echo "请检查PostgreSQL是否正在运行"
    exit 1
fi
echo ""

# 步骤5: 运行数据库迁移
echo "📋 步骤 5/8: 运行数据库迁移..."
cd /Users/zhenghongfei6/go/src/github.com/Alike

if go run cmd/migrate/main.go up 2>&1 | grep -q "Migration completed"; then
    echo -e "${GREEN}✅ 数据库迁移完成${NC}"
else
    echo -e "${YELLOW}⚠️  迁移可能失败，继续尝试...${NC}"
    if go run cmd/migrate/main.go up 2>&1; then
        echo -e "${GREEN}✅ 数据库迁移完成${NC}"
    else
        echo -e "${RED}❌ 数据库迁移失败${NC}"
        exit 1
    fi
fi
echo ""

# 步骤6: 导入测试数据
echo "🌱 步骤 6/8: 导入测试数据..."
if PGPASSWORD=alike_password psql -h localhost -U alike_user -d alike_db -f db/seeds/seed.sql &> /dev/null; then
    echo -e "${GREEN}✅ 测试数据导入成功${NC}"
else
    echo -e "${YELLOW}⚠️  测试数据可能已存在${NC}"
fi
echo ""

# 步骤7: 启动API服务器
echo "🚀 步骤 7/8: 启动API服务器..."

# 检查是否已有API在运行
if lsof -ti:8080 &> /dev/null; then
    echo -e "${YELLOW}⚠️  端口8080已被占用，正在停止...${NC}"
    lsof -ti:8080 | xargs kill -9 2>/dev/null || true
    sleep 2
fi

nohup go run cmd/api/main.go > /tmp/alike-api.log 2>&1 &
API_PID=$!
echo $API_PID > /tmp/alike-api.pid

echo "等待API启动..."
sleep 5

if curl -s http://localhost:8080/health > /dev/null; then
    echo -e "${GREEN}✅ API服务器启动成功 (PID: $API_PID)${NC}"
else
    echo -e "${RED}❌ API服务器启动失败${NC}"
    cat /tmp/alike-api.log
    exit 1
fi
echo ""

# 步骤8: 启动Web服务器
echo "🌐 步骤 8/8: 启动Web服务器..."

# 检查端口8002
if lsof -ti:8002 &> /dev/null; then
    echo -e "${YELLOW}⚠️  端口8002已被占用，正在停止...${NC}"
    lsof -ti:8002 | xargs kill -9 2>/dev/null || true
    sleep 2
fi

cd web/public
nohup python3 -m http.server 8002 > /tmp/alike-web.log 2>&1 &
WEB_PID=$!
echo $WEB_PID > /tmp/alike-web.pid

sleep 2

if curl -s http://localhost:8002 > /dev/null; then
    echo -e "${GREEN}✅ Web服务器启动成功 (PID: $WEB_PID)${NC}"
else
    echo -e "${RED}❌ Web服务器启动失败${NC}"
    cat /tmp/alike-web.log
    exit 1
fi
echo ""

# 完成
echo "================================"
echo -e "${GREEN}✅ 本地调试环境设置完成！${NC}"
echo "================================"
echo ""
echo "📊 服务状态："
echo "   API服务器: http://localhost:8080 (PID: $API_PID)"
echo "   Web服务器: http://localhost:8002 (PID: $WEB_PID)"
echo ""
echo "📱 访问地址："
echo "   启动页: http://localhost:8002/launcher.html"
echo "   全局聊天: http://localhost:8002/global-chat.html"
echo "   附近用户: http://localhost:8002/index.html"
echo ""
echo "🧪 测试账号："
echo "   手机: +8613800138000"
echo "   密码: password123"
echo ""
echo "📝 日志文件："
echo "   API日志: /tmp/alike-api.log"
echo "   Web日志: /tmp/alike-web.log"
echo ""
echo "🛑 停止服务："
echo "   kill $API_PID  # 停止API"
echo "   kill $WEB_PID  # 停止Web"
echo "   或运行: bash scripts/stop-all.sh"
echo ""
echo "🚀 开始测试吧！"
echo ""
