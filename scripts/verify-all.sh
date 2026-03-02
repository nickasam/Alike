#!/bin/bash

# 加载测试函数
source scripts/verify-functions.sh

echo "🧪 Alike 功能测试"
echo "================="
echo ""

# 步骤1: 健康检查
echo "📍 步骤 1/10: 健康检查"
if curl -s http://localhost:8080/health | grep -q "ok"; then
    echo -e "${GREEN}✅ API健康检查通过${NC}"
else
    echo -e "${RED}❌ API未运行${NC}"
    echo "请先运行: bash scripts/setup-local-env.sh"
    exit 1
fi
echo ""

# 步骤2: 用户注册
echo "📍 步骤 2/10: 用户注册"
test_api "用户注册" \
    "POST" \
    "/auth/register" \
    '{"phone":"+8613800138001","verification_code":"123456","nickname":"测试用户","password":"password123"}'

if [ "${TEST_RESULTS[-1]}" == "PASS" ]; then
    # 从响应中提取token
    REGISTER_RESPONSE=$(curl -s -X POST "$API_BASE/auth/register" \
        -H "Content-Type: application/json" \
        -d '{"phone":"+8613800138002","verification_code":"123456","nickname":"验证用户","password":"password123"}')
    
    ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    USERNAME=$(echo "$REGISTER_RESPONSE" | grep -o '"nickname":"[^"]*"' | cut -d'"' -f4)
    USER_ID=$(echo "$REGISTER_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    
    echo "  Token: ${ACCESS_TOKEN:0:20}..."
    echo "  Username: $USERNAME"
    echo "  User ID: $USER_ID"
else
    echo -e "${RED}❌ 注册失败，无法继续测试${NC}"
    exit 1
fi
echo ""

# 步骤3: 用户登录
echo "📍 步骤 3/10: 用户登录"
test_api "用户登录" \
    "POST" \
    "/auth/login" \
    '{"phone":"+8613800138000","password":"password123"}'

if [ "${TEST_RESULTS[-1]}" == "PASS" ]; then
    # 获取token
    LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
        -H "Content-Type: application/json" \
        -d '{"phone":"+8613800138000","password":"password123"}')
    
    MAIN_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    echo "  Main Token: ${MAIN_TOKEN:0:20}..."
fi
echo ""

# 步骤4: 获取用户信息
echo "📍 步骤 4/10: 获取用户信息"
test_api "获取用户信息" \
    "GET" \
    "/auth/me" \
    "" \
    "$MAIN_TOKEN"
echo ""

# 步骤5: 获取附近用户
echo "📍 步骤 5/10: 获取附近用户"
test_api "获取附近用户" \
    "GET" \
    "/users/nearby?lat=31.2304&lng=121.4737&radius=10" \
    "" \
    "$MAIN_TOKEN"
echo ""

# 步骤6: 获取匹配列表
echo "📍 步骤 6/10: 获取匹配列表"
test_api "获取匹配列表" \
    "GET" \
    "/matches" \
    "" \
    "$MAIN_TOKEN"
echo ""

# 步骤7: 发送喜欢
echo "📍 步骤 7/10: 发送喜欢"
test_api "发送喜欢" \
    "POST" \
    "/matches/test-user-001/like" \
    "" \
    "$MAIN_TOKEN"
echo ""

# 步骤8: 获取聊天列表
echo "📍 步骤 8/10: 获取聊天列表"
test_api "获取聊天列表" \
    "GET" \
    "/chats" \
    "" \
    "$MAIN_TOKEN"
echo ""

# 步骤9: 获取全局聊天室
echo "📍 步骤 9/10: 获取全局聊天室"
test_api "获取全局聊天室" \
    "GET" \
    "/global/room" \
    "" \
    "$MAIN_TOKEN"
echo ""

# 步骤10: 获取全局消息
echo "📍 步骤 10/10: 获取全局消息"
test_api "获取全局消息" \
    "GET" \
    "/global/messages" \
    "" \
    "$MAIN_TOKEN"
echo ""

# 生成测试报告
echo "================================"
echo "📊 测试报告"
echo "================================"
echo ""

PASS_COUNT=0
FAIL_COUNT=0

for i in "${!TEST_NAMES[@]}"; do
    if [ "${TEST_RESULTS[$i]}" == "PASS" ]; then
        echo -e "${GREEN}✅ ${TEST_NAMES[$i]}${NC}"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo -e "${RED}❌ ${TEST_NAMES[$i]}${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
done

echo ""
echo "总计: $((PASS_COUNT + FAIL_COUNT)) 个测试"
echo -e "${GREEN}通过: $PASS_COUNT${NC}"
echo -e "${RED}失败: $FAIL_COUNT${NC}"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}🎉 所有测试通过！${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  有 $FAIL_COUNT 个测试失败${NC}"
    exit 1
fi
