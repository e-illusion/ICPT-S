#!/bin/bash

# ICPT系统完整HTTPS服务启动脚本
# 同时启动前端(3000)和后端(8080)的HTTPS服务

echo "🔒 ICPT系统 - 启动完整HTTPS服务套件"
echo "======================================"
echo ""

# 检查运行目录
if [ ! -f "go.mod" ] || [ ! -d "web-modern" ]; then
    echo "❌ 请在icpt-system项目根目录运行此脚本"
    exit 1
fi

# 停止现有服务
echo "🛑 停止现有服务..."
pkill -f "./api-server" 2>/dev/null || true
pkill -f "go run cmd/server/main.go" 2>/dev/null || true
pkill -f "vite" 2>/dev/null || true
sleep 2

# 1. 启动后端HTTPS服务
echo ""
echo "🔥 第一步：启动后端HTTPS服务 (端口8080)"
echo "----------------------------------------"

# 检查后端Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go语言环境未安装，请先安装Go"
    exit 1
fi

# 检查后端配置
if [ ! -f "config.yaml" ]; then
    echo "❌ 后端config.yaml配置文件不存在"
    exit 1
fi

# 创建后端证书目录
mkdir -p certs

# 检查后端证书
if [ ! -f "certs/server.crt" ] || [ ! -f "certs/server.key" ]; then
    echo "🔐 生成后端SSL证书..."
    openssl req -x509 -newkey rsa:4096 -keyout certs/server.key -out certs/server.crt -days 365 -nodes \
        -subj "/C=CN/ST=Beijing/L=Beijing/O=ICPT/OU=Backend/CN=114.55.58.3"
    echo "✅ 后端SSL证书生成完成"
fi

# 编译后端
echo "🔨 编译后端项目..."
go build -o api-server cmd/server/main.go
if [ $? -ne 0 ]; then
    echo "❌ 后端编译失败"
    exit 1
fi

# 启动后端服务
echo "🚀 启动后端HTTPS服务..."
nohup ./api-server > backend.log 2>&1 &
BACKEND_PID=$!

# 等待后端启动
echo "⏳ 等待后端服务启动..."
sleep 3

# 检查后端服务
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo "❌ 后端服务启动失败，请检查日志："
    tail -20 backend.log
    exit 1
fi

echo "✅ 后端HTTPS服务启动成功！进程ID: $BACKEND_PID"

# 2. 启动前端HTTPS服务
echo ""
echo "🎨 第二步：启动前端HTTPS服务 (端口3000)"
echo "----------------------------------------"

# 切换到前端目录
cd web-modern

# 检查Node.js环境
if ! command -v npm &> /dev/null; then
    echo "❌ Node.js/npm环境未安装，请先安装Node.js"
    cd ..
    pkill -P $BACKEND_PID 2>/dev/null || true
    exit 1
fi

# 创建前端证书目录
mkdir -p cert

# 检查前端证书
if [ ! -f "cert/cert.pem" ] || [ ! -f "cert/key.pem" ]; then
    echo "🔐 生成前端SSL证书..."
    openssl req -x509 -newkey rsa:4096 -keyout cert/key.pem -out cert/cert.pem -days 365 -nodes \
        -subj "/C=CN/ST=Beijing/L=Beijing/O=ICPT/OU=Frontend/CN=114.55.58.3"
    chmod 600 cert/key.pem
    chmod 644 cert/cert.pem
    echo "✅ 前端SSL证书生成完成"
fi

# 检查并创建.env文件
if [ ! -f ".env" ]; then
    echo "⚙️ 创建前端环境配置..."
    cat > .env << EOF
# HTTPS和WSS配置
VITE_API_BASE_URL=https://114.55.58.3:8080/api/v1
VITE_APP_TITLE=ICPT 图像处理系统  
VITE_WS_URL=wss://114.55.58.3:8080/api/v1/ws
VITE_APP_VERSION=2.0.0

# 开发服务器配置
VITE_DEV_SERVER_HOST=0.0.0.0
VITE_DEV_SERVER_PORT=3000
VITE_DEV_SERVER_HTTPS=true

# 后端服务器配置
VITE_BACKEND_HOST=114.55.58.3
VITE_BACKEND_PORT=8080
VITE_BACKEND_HTTPS=true
EOF
    echo "✅ 前端环境配置创建完成"
fi

# 安装前端依赖
if [ ! -d "node_modules" ]; then
    echo "📦 安装前端依赖..."
    npm install
    if [ $? -ne 0 ]; then
        echo "❌ 前端依赖安装失败"
        cd ..
        pkill -P $BACKEND_PID 2>/dev/null || true
        exit 1
    fi
fi

# 启动前端服务
echo "🚀 启动前端HTTPS开发服务器..."
nohup npm run dev > frontend.log 2>&1 &
FRONTEND_PID=$!

# 回到项目根目录
cd ..

# 等待前端启动
echo "⏳ 等待前端服务启动..."
sleep 5

# 检查前端服务
if ! kill -0 $FRONTEND_PID 2>/dev/null; then
    echo "❌ 前端服务启动失败，请检查日志："
    tail -20 web-modern/frontend.log
    pkill -P $BACKEND_PID 2>/dev/null || true
    exit 1
fi

echo "✅ 前端HTTPS服务启动成功！进程ID: $FRONTEND_PID"

# 3. 服务启动完成
echo ""
echo "🎉 ICPT系统HTTPS服务启动完成！"
echo "======================================"
echo ""
echo "🌐 服务地址："
echo "  • 前端界面:    https://114.55.58.3:3000"
echo "  • 后端API:    https://114.55.58.3:8080/api/v1"
echo "  • WebSocket:  wss://114.55.58.3:8080/api/v1/ws"
echo "  • 健康检查:    https://114.55.58.3:8080/ping"
echo ""
echo "📊 进程信息："
echo "  • 后端进程ID:  $BACKEND_PID"
echo "  • 前端进程ID:  $FRONTEND_PID"
echo ""
echo "📋 重要提醒："
echo "  1. 首次访问时需要信任自签名证书"
echo "  2. 前端日志: web-modern/frontend.log"
echo "  3. 后端日志: backend.log"
echo "  4. 停止服务: pkill -f api-server && pkill -f vite"
echo ""
echo "🔍 实时查看日志："
echo "  • 后端: tail -f backend.log"
echo "  • 前端: tail -f web-modern/frontend.log"
echo ""
echo "🚀 现在您可以访问: https://114.55.58.3:3000"
echo "⚠️ 如果浏览器提示证书不安全，请点击'高级' -> '继续访问'" 