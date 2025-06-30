#!/bin/bash

# ICPT前端HTTPS服务启动脚本
# 支持HTTPS和WSS连接

echo "🚀 启动ICPT前端HTTPS服务..."

# 检查证书文件
if [ ! -f "cert/cert.pem" ] || [ ! -f "cert/key.pem" ]; then
    echo "❌ SSL证书文件不存在，正在生成自签名证书..."
    mkdir -p cert
    openssl req -x509 -newkey rsa:4096 -keyout cert/key.pem -out cert/cert.pem -days 365 -nodes \
        -subj "/C=CN/ST=Beijing/L=Beijing/O=ICPT/OU=Development/CN=114.55.58.3"
    echo "✅ SSL证书生成完成"
fi

# 检查证书权限
chmod 600 cert/key.pem
chmod 644 cert/cert.pem

# 检查环境变量
if [ ! -f ".env" ]; then
    echo "❌ .env文件不存在，正在创建默认配置..."
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
    echo "✅ .env文件创建完成"
fi

# 检查依赖
echo "📦 检查依赖包..."
if [ ! -d "node_modules" ]; then
    echo "⬇️ 安装依赖包..."
    npm install
fi

# 显示配置信息
echo ""
echo "📋 当前配置："
echo "  • 前端服务: https://114.55.58.3:3000"
echo "  • 后端API:  https://114.55.58.3:8080/api/v1"
echo "  • WebSocket: wss://114.55.58.3:8080/api/v1/ws"
echo "  • SSL证书:  自签名证书 (cert/cert.pem)"
echo ""

# 显示警告信息
echo "⚠️ 重要提醒："
echo "  1. 首次访问时，浏览器会提示证书不安全"
echo "  2. 请点击'高级' -> '继续访问' 来信任证书"
echo "  3. 确保后端服务 (8080端口) 已经启动并支持HTTPS"
echo "  4. 如果连接失败，请检查防火墙和端口配置"
echo ""

# 启动开发服务器
echo "🌟 启动HTTPS开发服务器..."
npm run dev

echo ""
echo "🎉 如果启动成功，请访问: https://114.55.58.3:3000" 