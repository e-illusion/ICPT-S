#!/bin/bash

# 高性能图像处理与传输系统 - 服务停止脚本
echo "🛑 停止高性能图像处理与传输系统..."

# 检查PID文件是否存在
if [ -f "logs/api-server.pid" ]; then
    API_PID=$(cat logs/api-server.pid)
    if ps -p $API_PID > /dev/null; then
        echo "🌐 停止API服务器 (PID: $API_PID)..."
        kill $API_PID
        # 等待进程结束
        sleep 2
        if ps -p $API_PID > /dev/null; then
            echo "⚠️ 强制停止API服务器..."
            kill -9 $API_PID
        fi
        echo "✅ API服务器已停止"
    else
        echo "⚠️ API服务器进程不存在"
    fi
    rm -f logs/api-server.pid
else
    echo "⚠️ 未找到API服务器PID文件，尝试按名称停止..."
    pkill -f "api-server"
fi

# 停止Worker进程
if [ -f "logs/worker.pid" ]; then
    WORKER_PID=$(cat logs/worker.pid)
    if ps -p $WORKER_PID > /dev/null; then
        echo "🔧 停止Worker进程 (PID: $WORKER_PID)..."
        kill $WORKER_PID
        # 等待进程结束
        sleep 2
        if ps -p $WORKER_PID > /dev/null; then
            echo "⚠️ 强制停止Worker进程..."
            kill -9 $WORKER_PID
        fi
        echo "✅ Worker进程已停止"
    else
        echo "⚠️ Worker进程不存在"
    fi
    rm -f logs/worker.pid
else
    echo "⚠️ 未找到Worker进程PID文件，尝试按名称停止..."
    pkill -f "worker"
fi

echo "✅ 所有服务已停止" 