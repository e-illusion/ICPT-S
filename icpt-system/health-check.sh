#!/bin/bash

# 高性能图像处理与传输系统 - 健康检查脚本
echo "🔍 系统健康检查..."
echo "================================"

# 检查颜色支持
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 成功/失败计数器
SUCCESS_COUNT=0
TOTAL_CHECKS=0

# 检查函数
check_status() {
    local description="$1"
    local command="$2"
    local expected_result="$3"
    
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
    printf "%-40s" "$description"
    
    if eval "$command" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 正常${NC}"
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    else
        echo -e "${RED}❌ 异常${NC}"
        if [ -n "$expected_result" ]; then
            echo "   预期: $expected_result"
        fi
    fi
}

# 开始检查
echo "📋 基础服务检查"
echo "--------------------------------"

# 1. MySQL服务检查
check_status "MySQL服务状态" "systemctl is-active --quiet mysql"

# 2. Redis服务检查
check_status "Redis服务状态" "pgrep -x redis-server"

# 3. API服务器进程检查
check_status "API服务器进程" "pgrep -f api-server"

# 4. Worker进程检查
check_status "Worker进程状态" "pgrep -f worker"

echo ""
echo "🌐 网络服务检查"
echo "--------------------------------"

# 5. API服务器端口检查
check_status "API服务器端口(8080)" "ss -tulpn | grep :8080"

# 6. API健康检查
check_status "API健康检查(/ping)" "curl -s -f http://localhost:8080/ping"

# 7. 首页访问检查
check_status "Web界面访问" "curl -s -f http://localhost:8080/ | grep -q html"

echo ""
echo "💾 数据存储检查"
echo "--------------------------------"

# 8. MySQL连接检查
check_status "MySQL数据库连接" "mysql -h 127.0.0.1 -u icpt_user -p123 -e 'SELECT 1' 2>/dev/null"

# 9. Redis连接检查
check_status "Redis缓存连接" "redis-cli ping"

# 10. 数据库表检查
check_status "数据库表结构" "mysql -h 127.0.0.1 -u icpt_user -p123 -e 'USE ICPT_System; SHOW TABLES;' 2>/dev/null | grep -q users"

echo ""
echo "📁 文件系统检查"
echo "--------------------------------"

# 11. 上传目录检查
check_status "上传目录权限" "test -w uploads/originals && test -w uploads/thumbnails"

# 12. 日志文件检查
check_status "日志文件存在" "test -f server.log && test -f worker.log"

# 13. 配置文件检查
check_status "配置文件格式" "python3 -c 'import yaml; yaml.safe_load(open(\"config.yaml\"))' 2>/dev/null || grep -q 'server:' config.yaml"

echo ""
echo "🔧 可执行文件检查"
echo "--------------------------------"

# 14. API服务器可执行文件
check_status "API服务器可执行文件" "test -x bin/api-server"

# 15. Worker可执行文件
check_status "Worker可执行文件" "test -x bin/worker"

# 16. 客户端可执行文件
check_status "客户端可执行文件" "test -x ../icpt-cli-client/bin/cli-client"

echo ""
echo "================================"
echo "📊 检查结果汇总"
echo "--------------------------------"

# 计算成功率
SUCCESS_RATE=$((SUCCESS_COUNT * 100 / TOTAL_CHECKS))

if [ $SUCCESS_RATE -ge 90 ]; then
    echo -e "总体状态: ${GREEN}优秀${NC} ($SUCCESS_COUNT/$TOTAL_CHECKS 通过, ${SUCCESS_RATE}%)"
    echo "🎉 系统运行正常，可以开始使用！"
elif [ $SUCCESS_RATE -ge 70 ]; then
    echo -e "总体状态: ${YELLOW}良好${NC} ($SUCCESS_COUNT/$TOTAL_CHECKS 通过, ${SUCCESS_RATE}%)"
    echo "⚠️  系统基本正常，但有一些小问题需要注意"
else
    echo -e "总体状态: ${RED}需要修复${NC} ($SUCCESS_COUNT/$TOTAL_CHECKS 通过, ${SUCCESS_RATE}%)"
    echo "❌ 系统存在较多问题，建议检查日志并修复"
fi

echo ""
echo "💡 快速修复建议:"
echo "  - 如果MySQL有问题: sudo systemctl start mysql"
echo "  - 如果Redis有问题: redis-server --daemonize yes"
echo "  - 如果进程未运行: ./start-services.sh"
echo "  - 查看详细日志: tail -f logs/api-server.log"
echo "  - 完整故障排查: 查看README.md的'故障排查'部分"

echo "================================" 