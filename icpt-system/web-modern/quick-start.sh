#!/bin/bash

# ICPT 现代化前端快速启动脚本
# Quick start script for ICPT Modern Frontend

echo "🎨 ICPT 图像处理系统 - 现代化前端"
echo "=================================="
echo ""

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${PURPLE}=== $1 ===${NC}"
}

# Check if we're in the right directory
if [ ! -f "package.json" ]; then
    print_error "package.json not found. Please run this script from the web-modern directory."
    exit 1
fi

# Check Node.js version
print_header "检查环境"
if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    print_status "Node.js version: $NODE_VERSION"
    
    # Extract major version number
    NODE_MAJOR=$(echo $NODE_VERSION | cut -d. -f1 | sed 's/v//')
    if [ "$NODE_MAJOR" -lt 16 ]; then
        print_warning "Node.js version should be >= 16.0.0. Current: $NODE_VERSION"
        print_warning "Please upgrade Node.js for best compatibility."
    fi
else
    print_error "Node.js is not installed. Please install Node.js >= 16.0.0"
    exit 1
fi

# Check package manager
PACKAGE_MANAGER=""
if command -v yarn &> /dev/null; then
    PACKAGE_MANAGER="yarn"
    print_status "Using Yarn package manager"
elif command -v npm &> /dev/null; then
    PACKAGE_MANAGER="npm"
    print_status "Using npm package manager"
else
    print_error "No package manager found. Please install npm or yarn."
    exit 1
fi

# Check if backend is running
print_header "检查后端服务"
if curl -s http://localhost:8080/ping > /dev/null 2>&1; then
    print_success "后端服务运行正常 (http://localhost:8080)"
else
    print_warning "后端服务未运行或不可访问"
    print_warning "请确保后端服务在 http://localhost:8080 上运行"
    echo ""
    echo "启动后端服务的命令:"
    echo "cd ../icpt-system && ./api-server"
    echo ""
    read -p "是否继续启动前端? (y/N): " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Install dependencies
print_header "安装依赖包"
if [ ! -d "node_modules" ] || [ ! -f "node_modules/.installed" ]; then
    print_status "正在安装项目依赖..."
    
    if [ "$PACKAGE_MANAGER" = "yarn" ]; then
        yarn install
    else
        npm install
    fi
    
    if [ $? -eq 0 ]; then
        touch node_modules/.installed
        print_success "依赖安装完成"
    else
        print_error "依赖安装失败"
        exit 1
    fi
else
    print_status "依赖已安装，跳过安装步骤"
fi

# Create .env file if it doesn't exist
print_header "配置环境变量"
if [ ! -f ".env" ]; then
    print_status "创建环境配置文件..."
    cat > .env << EOF
# ICPT 现代化前端环境配置
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_APP_TITLE=ICPT 图像处理系统
VITE_WS_URL=ws://localhost:8080/api/v1/ws
VITE_APP_VERSION=2.0.0
EOF
    print_success "环境配置文件已创建"
else
    print_status "环境配置文件已存在"
fi

# Show project information
print_header "项目信息"
echo -e "${CYAN}项目名称:${NC} ICPT 图像处理系统 - 现代化前端"
echo -e "${CYAN}技术栈:${NC} Vue.js 3 + Element Plus + Vite"
echo -e "${CYAN}开发地址:${NC} http://localhost:3000"
echo -e "${CYAN}API 地址:${NC} http://localhost:8080/api/v1"
echo -e "${CYAN}文档地址:${NC} ./README.md"
echo ""

# Show features
print_header "功能特性"
echo "✨ 现代化 Vue.js 3 Composition API"
echo "🎨 Element Plus 组件库"
echo "🚀 Vite 快速构建"
echo "📱 响应式设计"
echo "🌙 暗色主题支持"
echo "🔄 实时 WebSocket 通信"
echo "📊 数据可视化图表"
echo "🖼️ 拖拽图片上传"
echo "🔐 JWT 认证管理"
echo "⚡ 自动 API 代理"
echo ""

# Ask if user wants to start dev server
print_header "启动开发服务器"
read -p "是否立即启动开发服务器? (Y/n): " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Nn]$ ]]; then
    print_status "跳过启动，您可以稍后运行以下命令:"
    if [ "$PACKAGE_MANAGER" = "yarn" ]; then
        echo "yarn dev"
    else
        echo "npm run dev"
    fi
    echo ""
    print_success "项目配置完成！"
    exit 0
fi

# Start development server
print_status "正在启动开发服务器..."
echo ""
echo -e "${GREEN}🚀 启动中...${NC}"
echo ""
echo "开发服务器信息:"
echo "- 前端地址: http://localhost:3000"
echo "- 后端代理: http://localhost:8080"
echo "- 热重载: 已启用"
echo "- 源码映射: 已启用"
echo ""
echo "快捷键:"
echo "- Ctrl+C: 停止服务器"
echo "- Ctrl+Cmd+\\: 切换主题 (在浏览器中)"
echo ""
echo -e "${YELLOW}提示: 首次启动可能需要较长时间进行依赖优化${NC}"
echo ""

# Start the dev server
if [ "$PACKAGE_MANAGER" = "yarn" ]; then
    yarn dev
else
    npm run dev
fi 