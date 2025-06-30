# ICPT CLI 客户端

<div align="center">

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![OpenCV](https://img.shields.io/badge/OpenCV-4.0+-5C3EE8?style=for-the-badge&logo=opencv&logoColor=white)
![CLI](https://img.shields.io/badge/CLI-Tool-brightgreen?style=for-the-badge)

**ICPT 系统的命令行客户端工具**

</div>

## 📋 目录

- [项目概述](#项目概述)
- [功能特性](#功能特性)
- [安装指南](#安装指南)
- [配置说明](#配置说明)
- [使用手册](#使用手册)
- [摄像头功能](#摄像头功能)
- [故障排除](#故障排除)
- [开发指南](#开发指南)

## 🎯 项目概述

ICPT CLI 客户端是一个基于 Go 语言开发的命令行工具，提供与 ICPT 图像处理系统的完整交互功能。支持用户认证、图像上传、摄像头操作、批量处理等核心功能。

### ✨ 核心特性

- **🔐 用户认证**: 注册、登录、JWT令牌管理
- **📤 图像上传**: 单文件/批量上传，自动压缩优化
- **📷 摄像头操作**: 设备枚举、实时预览、拍照录制
- **📊 状态管理**: 实时查询处理状态、历史记录
- **⚡ 高性能**: 并发上传、智能重试、断点续传
- **🛠️ 易用性**: 交互式界面、详细帮助、进度显示

## 🌟 功能特性

### 🔐 用户管理
- **注册登录**: 支持用户名/邮箱注册和登录
- **令牌管理**: 自动存储和刷新JWT令牌
- **会话保持**: 记住登录状态，支持长期使用
- **权限验证**: 自动验证用户权限和访问控制

### 🖼️ 图像处理
- **智能上传**: 自动检测图像格式，优化传输
- **压缩算法**: 内置图像压缩，减少上传时间
- **批量操作**: 目录扫描，并发批量上传
- **进度显示**: 实时显示上传进度和处理状态
- **错误重试**: 自动重试机制，确保上传成功

### 📷 摄像头功能
- **设备检测**: 自动扫描可用摄像头设备
- **实时预览**: 摄像头实时画面预览
- **快速拍照**: 一键拍照并自动上传
- **视频录制**: 支持定时录制和手动控制
- **格式支持**: 多种图像和视频格式输出

### 📊 数据管理
- **列表查询**: 分页查看用户图像列表
- **状态监控**: 实时查询处理状态和进度
- **搜索过滤**: 按时间、状态、文件名搜索
- **统计信息**: 上传统计、成功率、平均时间

## 🛠️ 技术栈

- **开发语言**: Go 1.21+
- **图像处理**: GoCV + OpenCV 4.0+
- **HTTP客户端**: 内置HTTP客户端，支持认证
- **摄像头支持**: 跨平台摄像头驱动
- **配置管理**: YAML配置文件
- **日志系统**: 结构化日志记录

## 🚀 安装指南

### 环境要求

- **Go**: 1.21+ (如需从源码编译)
- **OpenCV**: 4.0+ (摄像头功能)
- **操作系统**: Linux/macOS/Windows

### 📦 预编译版本安装

```bash
# 1. 下载预编译版本
cd icpt-cli-client

# 2. 添加执行权限
chmod +x bin/cli-client

# 3. 测试安装
./bin/cli-client --version
```

### 🔨 从源码编译

```bash
# 1. 克隆项目
git clone <repository-url>
cd my-icpt-system/icpt-cli-client

# 2. 安装依赖（摄像头功能）
# Ubuntu/Debian:
sudo apt-get install libopencv-dev pkg-config

# CentOS/RHEL:
sudo yum install opencv-devel

# macOS:
brew install opencv pkg-config

# 3. 编译程序
go mod tidy
go build -o bin/cli-client ./cmd

# 4. 验证编译
./bin/cli-client help
```

### 🐳 Docker 使用

```bash
# 使用Docker运行（无摄像头功能）
docker run -it --rm \
  -v $(pwd)/config.yaml:/app/config.yaml \
  icpt-cli:latest help
```

## ⚙️ 配置说明

### 配置文件位置

- **默认配置**: `./config.yaml`
- **用户配置**: `~/.icpt/config.yaml`
- **系统配置**: `/etc/icpt/config.yaml`

### 配置文件示例

```yaml
# ICPT CLI 客户端配置
server:
  public_host: "http://localhost:8080"  # 服务器地址
  timeout: 30s                          # 请求超时时间
  
client:
  auto_compress: true                   # 自动压缩图像
  compression_quality: 85               # 压缩质量 (1-100)
  max_file_size: 50MB                   # 最大文件大小
  concurrent_uploads: 5                 # 并发上传数
  retry_attempts: 3                     # 重试次数
  
camera:
  default_device: 0                     # 默认摄像头设备ID
  capture_width: 1280                   # 拍照宽度
  capture_height: 720                   # 拍照高度
  video_fps: 30                         # 视频帧率
  
auth:
  token_file: "~/.icpt/token"          # 令牌存储文件
  auto_refresh: true                    # 自动刷新令牌
  
logging:
  level: "info"                         # 日志级别
  file: "~/.icpt/client.log"           # 日志文件
```

### 环境变量配置

```bash
# 设置环境变量
export ICPT_SERVER_HOST="https://your-server.com:8080"
export ICPT_TOKEN_FILE="/path/to/token"
export ICPT_LOG_LEVEL="debug"

# 或使用 .env 文件
echo "ICPT_SERVER_HOST=https://your-server.com:8080" > .env
```

## 📖 使用手册

### 基础命令

```bash
# 显示帮助信息
./bin/cli-client help
./bin/cli-client help upload  # 显示特定命令帮助

# 查看版本信息
./bin/cli-client --version

# 检查服务器连接
./bin/cli-client ping
```

### 🔐 用户认证

```bash
# 用户注册
./bin/cli-client register
# 输入提示：
# 用户名: myuser
# 邮箱: user@example.com
# 密码: ********

# 用户登录
./bin/cli-client login
# 输入提示：
# 用户名或邮箱: myuser
# 密码: ********

# 查看用户信息
./bin/cli-client profile

# 退出登录
./bin/cli-client logout
```

### 📤 图像上传

```bash
# 单文件上传
./bin/cli-client upload image.jpg
./bin/cli-client upload /path/to/image.png

# 批量上传目录
./bin/cli-client batch-upload ./photos/
./bin/cli-client batch-upload /home/user/images/

# 指定压缩质量上传
./bin/cli-client upload --quality 90 image.jpg

# 跳过压缩上传原图
./bin/cli-client upload --no-compress large-image.png

# 上传并设置标签
./bin/cli-client upload --tags "nature,landscape" photo.jpg
```

### 📊 图像管理

```bash
# 查看图像列表
./bin/cli-client list                    # 查看所有图像
./bin/cli-client list --page 2          # 查看第2页
./bin/cli-client list --size 20         # 每页20个
./bin/cli-client list --status pending  # 只看处理中的

# 查看图像详情
./bin/cli-client show 123              # 按ID查看
./bin/cli-client status 123            # 查看处理状态

# 搜索图像
./bin/cli-client search "sunset"       # 按文件名搜索
./bin/cli-client search --tags "nature" # 按标签搜索

# 删除图像
./bin/cli-client delete 123            # 删除指定图像
./bin/cli-client delete --all          # 删除所有图像（需确认）
```

### 📊 统计信息

```bash
# 查看用户统计
./bin/cli-client stats

# 输出示例：
# 📊 用户统计信息
# ==================
# 总上传数量: 156 张
# 处理成功: 152 张 (97.4%)
# 处理失败: 4 张 (2.6%)
# 平均处理时间: 1.8 秒
# 总文件大小: 245.6 MB
# 今日上传: 23 张
```

## 📷 摄像头功能

### 摄像头管理

```bash
# 列出可用摄像头
./bin/cli-client camera list

# 输出示例：
# 🔍 扫描可用摄像头...
# ✅ 发现 2 个摄像头设备:
# ID | 设备名称        | 分辨率      | 状态
# ---|----------------|-------------|------
#  0 | USB Camera     | 1280x720    | 可用
#  1 | Built-in Camera| 1920x1080   | 可用

# 查看摄像头详细信息
./bin/cli-client camera info 0
```

### 实时预览

```bash
# 启动摄像头预览
./bin/cli-client camera preview 0

# 预览控制：
# - 按 's' 键拍照并保存
# - 按 'u' 键拍照并上传
# - 按 'q' 键退出预览
# - 按 'h' 键显示帮助

# 设置预览参数
./bin/cli-client camera preview 0 --width 1920 --height 1080
```

### 拍照功能

```bash
# 快速拍照（保存到本地）
./bin/cli-client camera capture 0

# 拍照并自动上传
./bin/cli-client camera capture 0 --upload

# 连续拍照
./bin/cli-client camera capture 0 --count 5 --interval 2s

# 定时拍照
./bin/cli-client camera capture 0 --timer 5s

# 设置拍照参数
./bin/cli-client camera capture 0 \
  --width 1920 \
  --height 1080 \
  --quality 95 \
  --output ./captures/
```

### 视频录制

```bash
# 录制视频（10秒）
./bin/cli-client camera record 0 --duration 10s

# 录制到指定文件
./bin/cli-client camera record 0 \
  --duration 30s \
  --output video.mp4 \
  --fps 30

# 手动控制录制（按q停止）
./bin/cli-client camera record 0 --manual

# 录制完成后自动上传
./bin/cli-client camera record 0 --duration 10s --upload
```

## 🎛️ 高级功能

### 配置管理

```bash
# 查看当前配置
./bin/cli-client config show

# 设置配置项
./bin/cli-client config set server.public_host "https://new-server.com"
./bin/cli-client config set client.auto_compress false

# 重置配置
./bin/cli-client config reset

# 生成配置文件模板
./bin/cli-client config init
```

### 批量操作

```bash
# 批量上传并监控进度
./bin/cli-client batch-upload ./photos/ --progress

# 批量设置标签
./bin/cli-client batch-tag --tag "vacation2024" 101 102 103

# 批量删除
./bin/cli-client batch-delete --status failed

# 导出图像列表
./bin/cli-client export --format csv --output images.csv
./bin/cli-client export --format json --output images.json
```

### 脚本自动化

```bash
# 定时上传脚本
#!/bin/bash
# upload-cron.sh
./bin/cli-client login --username user --password-file .passwd
./bin/cli-client batch-upload /path/to/watch/folder/
./bin/cli-client logout

# 摄像头监控脚本
#!/bin/bash
# camera-monitor.sh
while true; do
  ./bin/cli-client camera capture 0 --upload
  sleep 300  # 每5分钟拍照一次
done
```

## 🔧 故障排除

### 常见问题

**1. 认证失败**
```bash
# 检查服务器连接
./bin/cli-client ping

# 清除本地令牌
rm ~/.icpt/token
./bin/cli-client login

# 检查服务器地址配置
./bin/cli-client config show | grep server
```

**2. 摄像头无法使用**
```bash
# 检查摄像头设备
ls /dev/video*  # Linux
system_profiler SPCameraDataType  # macOS

# 检查OpenCV安装
pkg-config --modversion opencv4

# 重新编译客户端
go clean -cache
go build -o bin/cli-client ./cmd
```

**3. 上传失败**
```bash
# 检查文件权限
ls -la image.jpg

# 检查文件大小限制
./bin/cli-client config show | grep max_file_size

# 检查网络连接
curl -I http://localhost:8080/ping

# 查看详细错误信息
./bin/cli-client upload image.jpg --verbose
```

**4. 配置问题**
```bash
# 检查配置文件路径
./bin/cli-client config show

# 重新生成配置
./bin/cli-client config init --force

# 使用默认配置
mv config.yaml config.yaml.bak
./bin/cli-client config init
```

### 调试模式

```bash
# 启用详细日志
./bin/cli-client --verbose upload image.jpg

# 启用调试模式
./bin/cli-client --debug camera list

# 查看日志文件
tail -f ~/.icpt/client.log

# 环境变量调试
export ICPT_LOG_LEVEL=debug
export ICPT_DEBUG=true
./bin/cli-client upload image.jpg
```

### 性能优化

```bash
# 调整并发数
./bin/cli-client config set client.concurrent_uploads 10

# 优化压缩设置
./bin/cli-client config set client.compression_quality 75
./bin/cli-client config set client.auto_compress true

# 增加超时时间
./bin/cli-client config set server.timeout 60s

# 启用本地缓存
./bin/cli-client config set client.enable_cache true
```

## 👨‍💻 开发指南

### 项目结构

```
icpt-cli-client/
├── cmd/                    # 应用程序入口
│   └── main.go            # 主程序
├── internal/              # 内部包
│   ├── auth/              # 认证管理
│   ├── camera/            # 摄像头操作
│   ├── compress/          # 图像压缩
│   ├── config/            # 配置管理
│   └── httpclient/        # HTTP客户端
├── bin/                   # 编译输出
│   └── cli-client         # 可执行文件
├── config.yaml            # 配置文件
├── go.mod                 # Go模块
└── README.md              # 项目文档
```

### 编译构建

```bash
# 开发环境编译
go build -o bin/cli-client ./cmd

# 生产环境编译（优化）
go build -ldflags="-s -w" -o bin/cli-client ./cmd

# 跨平台编译
GOOS=linux GOARCH=amd64 go build -o bin/cli-client-linux ./cmd
GOOS=windows GOARCH=amd64 go build -o bin/cli-client.exe ./cmd
GOOS=darwin GOARCH=amd64 go build -o bin/cli-client-mac ./cmd

# 静态链接编译
CGO_ENABLED=0 go build -a -ldflags="-s -w" -o bin/cli-client ./cmd
```

### 代码规范

```bash
# 代码格式化
go fmt ./...

# 代码检查
go vet ./...
golint ./...

# 安全检查
gosec ./...

# 依赖管理
go mod tidy
go mod verify
```

### 测试

```bash
# 运行单元测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成测试报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 运行基准测试
go test -bench=. ./...
```

## 📚 相关文档

- **[项目主文档](../README.md)** - 整体项目介绍
- **[后端文档](../icpt-system/README.md)** - 后端API和服务
- **[前端文档](../icpt-system/web-modern/README.md)** - 前端界面使用

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

### 贡献流程
1. Fork 项目
2. 创建功能分支
3. 提交代码
4. 创建 Pull Request

### 开发约定
- 遵循 Go 代码规范
- 添加适当的测试
- 更新相关文档
- 提交信息使用英文

## 📄 许可证

本项目采用 [MIT 许可证](../LICENSE)。

---

<div align="center">

**🚀 高效的命令行工具，让图像处理更简单！🚀**

</div> 