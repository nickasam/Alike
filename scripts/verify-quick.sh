#!/bin/bash

export PATH="/opt/homebrew/opt/postgresql@15/bin:$PATH"

API_BASE="http://localhost:8080/api/v1"

echo "🧪 快速功能测试"
echo "================="
echo ""

# 颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

PASSED=0
FAILED=0

# 测试函数
test_endpoint() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local token=$5
    
    echo -n "测试 $name ... "
    
    if [ -n "$data" ]; then
        response=$(curl -s -X $method "$API_BASE$endpoint" \
            -H "Content-Type: application/json" \
            ${token:+-H "Authorization: Bearer $token"} \
            -d "$data")
    else
        response=$(curl -s -X $method "$API_BASE$endpoint" \
            ${token:+-H "Authorization: Bearer $token"})
    fi
    
    if echo "$response" | grep -q '"success":true\|"status":"ok"'; then
        echo -e "${GREEN}✅ PASS${NC}"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}❌ FAIL${NC}"
        echo "  Response: $response"
        FAILED=$((FAILED + 1))
    fi
}

# 1. 健康检查
test_endpoint "健康检查" "GET" "/health"

# 2. 注册
test_endpoint "用户注册" "POST" "/auth/register" '{"phone":"+8613800138005","verification_code":"123456","nickname":"测试用户5","password":"password123"}'

# 3. 登录
LOGIN_RESP=$(curl -s -X POST "$API_BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"phone":"+8613800138000","password":"password123"}')

if echo "$LOGIN_RESP" | grep -q '"success":true'; then
    TOKEN=$(echo "$LOGIN_RESP" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✅ 用户登录 PASS${NC}"
    echo "  Token: ${TOKEN:0:20}..."
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}❌ 用户登录 FAIL${NC}"
    echo "  Response: $LOGIN_RESP"
    FAILED=$((FAILED + 1))
fi

# 4-10. 需要认证的测试
if [ -n "$TOKEN" ]; then
    test_endpoint "获取用户信息" "GET" "/auth/me" "" "$TOKEN"
    test_endpoint "获取附近用户" "GET" "/users/nearby?lat=31.2304&lng=121.4737&radius=10" "" "$TOKEN"
    test_endpoint "获取匹配列表" "GET" "/matches" "" "$TOKEN"
    test_endpoint "获取聊天列表" "GET" "/chats" "" "$TOKEN"
    test_endpoint "获取全局聊天室" "GET" "/global/room" "" "$TOKEN"
    test_endpoint "获取全局消息" "GET" "/global/messages" "" "$TOKEN"
fi

echo ""
echo "================================"
echo "📊 测试结果"
echo "================================"
echo -e "${GREEN}通过: $PASSED${NC}"
echo -e "${RED}失败: $FAILED${NC}"
echo "总计: $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}⚠️  有测试失败${NC}"
    exit 1
fi
