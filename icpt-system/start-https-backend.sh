#!/bin/bash

# ICPT后端HTTPS服务启动脚本
# 支持HTTPS API和WSS WebSocket

echo "🔒 启动ICPT后端HTTPS服务..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go语言环境未安装，请先安装Go"
    exit 1
fi

# 检查当前目录
if [ ! -f "go.mod" ]; then
    echo "❌ 请在项目根目录运行此脚本"
    exit 1
fi

# 停止旧的服务
echo "🛑 停止现有服务..."
pkill -f "./api-server" 2>/dev/null || true
pkill -f "go run cmd/server/main.go" 2>/dev/null || true

# 检查证书文件
echo "🔐 检查SSL证书..."
if [ ! -f "certs/server.crt" ] || [ ! -f "certs/server.key" ]; then
    echo "⚠️ SSL证书文件不存在，将在启动时自动生成"
fi

# 检查配置文件
echo "⚙️ 检查配置文件..."
if [ ! -f "config.yaml" ]; then
    echo "❌ config.yaml配置文件不存在"
    exit 1
fi

# 显示当前配置
echo ""
echo "📋 当前HTTPS配置："
echo "  • 后端API:  https://114.55.58.3:8080/api/v1"
echo "  • WebSocket: wss://114.55.58.3:8080/api/v1/ws"
echo "  • 证书路径: certs/server.crt"
echo "  • 静态文件: https://114.55.58.3:8080/static/"
echo ""

# 编译项目
echo "🔨 编译项目..."
go build -o api-server cmd/server/main.go
if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi

echo "✅ 编译成功"

# 启动服务
echo "🚀 启动HTTPS API服务器..."
echo "📝 日志将写入 server.log"
echo ""

# 后台启动服务并重定向日志
nohup ./api-server > server.log 2>&1 &
SERVER_PID=$!

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 3

# 检查服务是否正常启动
if kill -0 $SERVER_PID 2>/dev/null; then
    echo "✅ 后端HTTPS服务启动成功！"
    echo "📊 进程ID: $SERVER_PID"
    echo ""
    echo "🌐 服务地址："
    echo "  • API接口: https://114.55.58.3:8080/api/v1"
    echo "  • WebSocket: wss://114.55.58.3:8080/api/v1/ws"
    echo "  • 健康检查: https://114.55.58.3:8080/ping"
    echo ""
    echo "📋 重要提醒："
    echo "  1. 首次访问时需要信任自签名证书"
    echo "  2. 日志文件: server.log"
    echo "  3. 停止服务: pkill -f './api-server'"
    echo ""
    echo "🔍 实时查看日志: tail -f server.log"
else
    echo "❌ 服务启动失败，请检查日志："
    tail -20 server.log
    exit 1
fi 