#!/bin/bash

# 高性能图像处理与传输系统 - 服务启动脚本
echo "🚀 启动高性能图像处理与传输系统..."

# 检查依赖
echo "📋 检查服务依赖..."

# 检查MySQL服务
if systemctl is-active --quiet mysql; then
    echo "✅ MySQL服务正在运行"
else
    echo "❌ MySQL服务未运行，正在启动..."
    sudo systemctl start mysql
    if [ $? -eq 0 ]; then
        echo "✅ MySQL服务启动成功"
    else
        echo "❌ MySQL服务启动失败，请检查配置"
        exit 1
    fi
fi

# 检查Redis服务
if pgrep -x "redis-server" > /dev/null; then
    echo "✅ Redis服务正在运行"
else
    echo "⚠️ Redis服务未运行，正在启动..."
    if command -v redis-server >/dev/null 2>&1; then
        redis-server --daemonize yes
        echo "✅ Redis服务启动成功"
    else
        echo "❌ Redis未安装，请先安装Redis"
        echo "Ubuntu/Debian: sudo apt-get install redis-server"
        echo "CentOS/RHEL: sudo yum install redis"
        exit 1
    fi
fi

# 检查可执行文件
if [ ! -f "./bin/api-server" ]; then
    echo "❌ API服务器可执行文件不存在，请先编译项目"
    echo "执行: go build -o bin/api-server ./cmd/server"
    exit 1
fi

if [ ! -f "./bin/worker" ]; then
    echo "❌ Worker可执行文件不存在，请先编译项目"
    echo "执行: go build -o bin/worker ./cmd/worker"
    exit 1
fi

# 创建必要的目录
echo "📁 创建必要的目录..."
mkdir -p uploads/originals
mkdir -p uploads/thumbnails
mkdir -p logs

# 检查配置文件
if [ ! -f "config.yaml" ]; then
    echo "❌ 配置文件不存在，请确保config.yaml文件存在"
    exit 1
fi

echo "✅ 所有依赖检查完成"

# 启动API服务器
echo "🌐 启动API服务器..."
nohup ./bin/api-server > logs/api-server.log 2>&1 &
API_PID=$!
echo "✅ API服务器已启动 (PID: $API_PID)"

# 等待API服务器启动
sleep 3

# 检查API服务器是否启动成功
if ps -p $API_PID > /dev/null; then
    echo "✅ API服务器运行正常"
else
    echo "❌ API服务器启动失败，请检查logs/api-server.log"
    exit 1
fi

# 启动Worker进程
echo "🔧 启动Worker进程..."
nohup ./bin/worker > logs/worker.log 2>&1 &
WORKER_PID=$!
echo "✅ Worker进程已启动 (PID: $WORKER_PID)"

# 等待Worker进程启动
sleep 2

# 检查Worker进程是否启动成功
if ps -p $WORKER_PID > /dev/null; then
    echo "✅ Worker进程运行正常"
else
    echo "❌ Worker进程启动失败，请检查logs/worker.log"
    exit 1
fi

# 保存PID到文件
echo $API_PID > logs/api-server.pid
echo $WORKER_PID > logs/worker.pid

echo ""
echo "🎉 系统启动完成！"
echo "================================"
echo "📊 服务状态："
echo "  • API服务器: http://localhost:8080 (PID: $API_PID)"
echo "  • Worker进程: 正在运行 (PID: $WORKER_PID)"
echo ""
echo "📋 可用接口："
echo "  • 健康检查: http://localhost:8080/ping"
echo "  • Web界面: http://localhost:8080/"
echo "  • API文档: http://localhost:8080/api/v1/"
echo ""
echo "📝 日志文件："
echo "  • API服务器: logs/api-server.log"
echo "  • Worker进程: logs/worker.log"
echo ""
echo "🛑 停止服务："
echo "  • 执行: ./stop-services.sh"
echo "  • 或手动停止: kill $API_PID $WORKER_PID"
echo "================================" 