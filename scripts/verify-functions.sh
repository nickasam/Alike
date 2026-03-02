#!/bin/bash

# API基础URL
API_BASE="http://localhost:8080/api/v1"

# 测试结果存储
declare -a TEST_RESULTS
declare -a TEST_NAMES

# 颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 测试函数
test_api() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local token=$5
    
    TEST_NAMES+=("$name")
    
    echo "测试: $name"
    
    if [ -n "$data" ]; then
        if [ -n "$token" ]; then
            response=$(curl -s -X $method "$API_BASE$endpoint" \
                -H "Authorization: Bearer $token" \
                -H "Content-Type: application/json" \
                -d "$data")
        else
            response=$(curl -s -X $method "$API_BASE$endpoint" \
                -H "Content-Type: application/json" \
                -d "$data")
        fi
    else
        if [ -n "$token" ]; then
            response=$(curl -s -X $method "$API_BASE$endpoint" \
                -H "Authorization: Bearer $token")
        else
            response=$(curl -s -X $method "$API_BASE$endpoint")
        fi
    fi
    
    # 检查响应
    if echo "$response" | grep -q '"success":true'; then
        echo -e "  ${GREEN}✅ PASS${NC}"
        TEST_RESULTS+=("PASS")
        return 0
    else
        echo -e "  ${RED}❌ FAIL${NC}"
        echo "  Response: $response"
        TEST_RESULTS+=("FAIL")
        return 1
    fi
}

wait_for_api() {
    echo "等待API启动..."
    for i in {1..30}; do
        if curl -s http://localhost:8080/health > /dev/null; then
            echo -e "${GREEN}✅ API已就绪${NC}"
            return 0
        fi
        sleep 1
    done
    echo -e "${RED}❌ API启动超时${NC}"
    return 1
}
