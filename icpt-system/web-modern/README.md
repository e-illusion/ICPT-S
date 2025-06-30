# ICPT 图像处理系统 - 现代化前端

<div align="center">

![Vue.js](https://img.shields.io/badge/Vue.js-3.4.0-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white)
![Element Plus](https://img.shields.io/badge/Element%20Plus-2.5.0-409EFF?style=for-the-badge&logo=element&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-4.5.0-646CFF?style=for-the-badge&logo=vite&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**企业级图像处理系统的现代化前端界面**

[English](README_EN.md) · [中文文档](README.md) · [演示预览](#演示预览) · [快速开始](#快速开始)

</div>

## 📋 目录

- [项目概述](#项目概述)
- [功能特性](#功能特性)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [功能模块](#功能模块)
- [开发指南](#开发指南)
- [部署指南](#部署指南)
- [故障排除](#故障排除)
- [性能优化](#性能优化)
- [贡献指南](#贡献指南)

## 🎯 项目概述

ICPT 现代化前端是一个基于 Vue.js 3 构建的企业级图像处理系统用户界面，提供直观、高效的图像管理和处理体验。系统采用最新的前端技术栈，支持实时通信、响应式设计和现代化的用户交互。

### ✨ 核心优势

- 🚀 **极速响应**：基于 Vite 的快速开发和构建
- 📱 **响应式设计**：完美适配桌面端、平板和移动设备
- 🎨 **现代化 UI**：Element Plus 组件库，支持明暗主题
- ⚡ **实时通信**：WebSocket 实时状态更新和通知
- 🔐 **安全认证**：JWT 令牌管理和路由权限控制
- 🖼️ **智能上传**：拖拽上传、进度显示、错误重试

## 🌟 功能特性

### 🔐 用户认证系统
- **安全登录**：JWT 令牌认证，支持记住登录状态
- **自动登录**：刷新页面保持登录状态
- **权限控制**：基于路由的访问权限管理
- **会话管理**：令牌过期自动处理和刷新

### 🖼️ 图像管理系统
- **智能上传**：拖拽上传、批量上传、格式验证
- **摄像头拍照**：实时摄像头预览、拍照上传、镜像模式
- **实时预览**：上传进度显示、缩略图预览
- **状态管理**：处理中、已完成、失败状态实时更新
- **批量操作**：多选删除、状态筛选、搜索功能

### 📊 数据可视化
- **仪表盘**：系统概览、统计图表、快捷操作
- **实时监控**：处理进度、成功率、平均处理时间
- **历史记录**：操作日志、处理历史查询

### 🌐 实时通信
- **WebSocket 连接**：自动重连、心跳检测
- **即时通知**：处理完成、错误提醒、系统通知
- **状态同步**：多页面状态实时同步更新

### 🎨 用户体验
- **主题切换**：明亮/暗黑主题一键切换
- **响应式布局**：适配各种屏幕尺寸
- **流畅动画**：页面切换、组件交互动画效果
- **键盘快捷键**：提高操作效率

## 🛠️ 技术栈

### 核心框架
- **[Vue.js 3](https://vuejs.org/)** - 渐进式 JavaScript 框架
- **[Composition API](https://vuejs.org/guide/extras/composition-api-faq.html)** - Vue 3 组合式 API
- **[Vue Router 4](https://router.vuejs.org/)** - 官方路由管理器
- **[Pinia](https://pinia.vuejs.org/)** - 状态管理库

### UI 组件库
- **[Element Plus](https://element-plus.org/)** - Vue 3 组件库
- **[Element Plus Icons](https://element-plus.org/en-US/component/icon.html)** - 图标库
- **[SCSS](https://sass-lang.com/)** - CSS 预处理器

### 开发工具
- **[Vite](https://vitejs.dev/)** - 下一代前端构建工具
- **[ESLint](https://eslint.org/)** - 代码质量检查
- **[Prettier](https://prettier.io/)** - 代码格式化

### 功能库
- **[Axios](https://axios-http.com/)** - HTTP 客户端
- **[ECharts](https://echarts.apache.org/)** - 数据可视化图表
- **[Day.js](https://dayjs.gitee.io/)** - 日期时间处理
- **[NProgress](https://ricostacruz.com/nprogress/)** - 页面加载进度条
- **[js-cookie](https://github.com/js-cookie/js-cookie)** - Cookie 管理

## 🚀 快速开始

### 环境要求

- **Node.js**: >= 16.0.0 (推荐 18.0.0+)
- **npm**: >= 8.0.0 或 **yarn**: >= 1.22.0
- **现代浏览器**: Chrome 90+, Firefox 88+, Safari 14+

### 安装依赖

```bash
# 克隆项目
git clone <repository-url>
cd icpt-system/web-modern

# 安装依赖
npm install
# 或使用 yarn
yarn install
```

### 环境配置

创建环境配置文件：

```bash
# 复制环境配置模板
cp .env.example .env

# 编辑环境配置
nano .env
```

```env
# ICPT 现代化前端环境配置
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_APP_TITLE=ICPT 图像处理系统
VITE_WS_URL=ws://localhost:8080/api/v1/ws
VITE_APP_VERSION=2.0.0
```

### 启动开发服务器

```bash
# 启动开发服务器
npm run dev
# 或
yarn dev

# 开发服务器将在 http://localhost:3000 启动
```

### 构建生产版本

```bash
# 构建生产版本
npm run build
# 或
yarn build

# 预览生产构建
npm run preview
# 或
yarn preview
```

## 📁 项目结构

```
web-modern/
├── public/                     # 静态资源
│   ├── favicon.ico            # 网站图标
│   └── ...                    # 其他静态文件
├── src/                       # 源代码
│   ├── api/                   # API 接口封装
│   │   ├── auth.js           # 认证相关接口
│   │   ├── images.js         # 图像管理接口
│   │   ├── request.js        # HTTP 请求基础配置
│   │   └── websocket.js      # WebSocket 服务
│   ├── layout/               # 布局组件
│   │   └── index.vue         # 主布局组件
│   ├── router/               # 路由配置
│   │   └── index.js          # 路由定义和守卫
│   ├── stores/               # 状态管理
│   │   └── auth.js           # 用户认证状态
│   ├── styles/               # 全局样式
│   │   ├── index.scss        # 样式入口
│   │   ├── variables.scss    # CSS 变量
│   │   └── themes.scss       # 主题配置
│   ├── utils/                # 工具函数
│   │   └── auth.js           # 认证工具函数
│   ├── views/                # 页面组件
│   │   ├── auth/             # 认证相关页面
│   │   │   └── Login.vue     # 登录页面
│   │   ├── dashboard/        # 仪表盘
│   │   │   └── index.vue     # 仪表盘主页
│   │   ├── images/           # 图像管理
│   │   │   ├── Upload.vue    # 图像上传
│   │   │   ├── Gallery.vue   # 图像列表
│   │   │   └── Detail.vue    # 图像详情
│   │   ├── error/            # 错误页面
│   │   │   ├── 404.vue       # 404 页面
│   │   │   ├── 403.vue       # 403 页面
│   │   │   └── 500.vue       # 500 页面
│   │   └── ...               # 其他页面
│   ├── App.vue               # 根组件
│   └── main.js               # 应用入口
├── .env.example              # 环境配置模板
├── package.json              # 项目配置
├── vite.config.js            # Vite 配置
├── quick-start.sh            # 快速启动脚本
└── README.md                 # 项目文档
```

## 🔧 功能模块

### 1. 用户认证模块

**登录功能**
```javascript
// 使用认证 store
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// 用户登录
await authStore.login({
  username: 'admin',
  password: 'admin123'
})

// 检查认证状态
if (authStore.isAuthenticated) {
  // 用户已登录
}
```

**权限控制**
```javascript
// 路由权限配置
{
  path: '/images/upload',
  component: Upload,
  meta: {
    requiresAuth: true,  // 需要登录访问
    title: '图像上传'
  }
}
```

### 2. 图像上传模块

**拖拽上传**
```vue
<template>
  <div 
    @drop="handleDrop" 
    @dragover="handleDragOver"
    class="upload-zone"
  >
    拖拽文件到这里上传
  </div>
</template>

<script setup>
const handleDrop = (event) => {
  const files = Array.from(event.dataTransfer.files)
  uploadFiles(files)
}
</script>
```

**摄像头拍照上传**
```vue
<template>
  <CameraCapture 
    @uploaded="handleCameraUpload"
    @error="handleCameraError"
    :quality="0.8"
    :width="800"
    :height="600"
  />
</template>

<script setup>
import CameraCapture from '@/components/CameraCapture.vue'

const handleCameraUpload = (data) => {
  console.log('摄像头照片上传成功:', data)
  // 处理上传成功逻辑
}
</script>
```

**进度监控**
```javascript
import { uploadImage } from '@/api/images'

// 上传文件并监控进度
await uploadImage(file, (progress) => {
  console.log(`上传进度: ${progress}%`)
})
```

### 3. WebSocket 通信

**连接管理**
```javascript
import webSocketService from '@/api/websocket'

// 连接 WebSocket
await webSocketService.connect()

// 监听事件
webSocketService.on('image_completed', (data) => {
  console.log('图像处理完成:', data)
})

// 断开连接
webSocketService.disconnect()
```

**事件处理**
```javascript
// 处理不同类型的 WebSocket 消息
const messageHandlers = {
  'image_processing': (data) => {
    // 处理中状态
  },
  'image_completed': (data) => {
    // 处理完成
  },
  'image_failed': (data) => {
    // 处理失败
  }
}
```

### 4. 主题系统

**主题切换**
```javascript
// 切换主题
const toggleTheme = () => {
  const newTheme = currentTheme === 'light' ? 'dark' : 'light'
  document.documentElement.setAttribute('data-theme', newTheme)
  localStorage.setItem('theme', newTheme)
}
```

**CSS 变量配置**
```scss
:root {
  --el-color-primary: #409eff;
  --el-color-success: #67c23a;
  --el-color-warning: #e6a23c;
  --el-color-danger: #f56c6c;
}

[data-theme='dark'] {
  --el-color-primary: #337ecc;
  --el-bg-color: #141414;
  --el-text-color-primary: #e5eaf3;
}
```

## 👨‍💻 开发指南

### 代码规范

项目遵循以下代码规范：

- **Vue 风格指南**：遵循 [Vue.js 官方风格指南](https://vuejs.org/style-guide/)
- **ESLint 规则**：基于 `@vue/eslint-config-standard`
- **Prettier 格式化**：统一代码格式
- **命名约定**：
  - 组件名：PascalCase (`UserProfile.vue`)
  - 函数名：camelCase (`getUserInfo`)
  - 常量名：SCREAMING_SNAKE_CASE (`API_BASE_URL`)

### 开发调试

**启用调试模式**
```bash
# 开发环境调试
npm run dev --debug

# 查看网络请求
# 浏览器开发者工具 -> Network 选项卡
```

**性能分析**
```bash
# 构建分析
npm run build --analyze

# 查看构建产物大小和依赖关系
```

**错误监控**
```javascript
// 全局错误处理
app.config.errorHandler = (error, instance, info) => {
  console.error('Global error:', error)
  // 发送错误报告到监控系统
}
```

### 组件开发

**组件模板**
```vue
<template>
  <div class="component-name">
    <!-- 组件内容 -->
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// 组件逻辑
const loading = ref(false)
const data = ref([])

const computedValue = computed(() => {
  return data.value.length
})

onMounted(() => {
  // 初始化逻辑
})
</script>

<style lang="scss" scoped>
.component-name {
  // 组件样式
}
</style>
```

**API 调用最佳实践**
```javascript
// 使用 try-catch 处理错误
const loadData = async () => {
  try {
    loading.value = true
    const response = await api.getData()
    data.value = response.data
  } catch (error) {
    console.error('加载数据失败:', error)
    ElMessage.error('加载失败，请重试')
  } finally {
    loading.value = false
  }
}
```

## 🚀 部署指南

### 生产构建

```bash
# 构建生产版本
npm run build

# 构建输出在 dist/ 目录
ls dist/
```

### Nginx 配置

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/dist;
    index index.html;

    # 处理 Vue Router 的 history 模式
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # API 代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket 代理
    location /api/v1/ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Docker 部署

**Dockerfile**
```dockerfile
# 构建阶段
FROM node:18-alpine as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build

# 生产阶段
FROM nginx:alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

**构建和运行**
```bash
# 构建镜像
docker build -t icpt-frontend .

# 运行容器
docker run -d -p 80:80 --name icpt-frontend icpt-frontend
```

## 🔍 故障排除

### 常见问题

**1. 登录后需要重新登录**
```javascript
// 检查 token 存储
console.log('Token:', localStorage.getItem('icpt_token'))

// 检查认证状态
import { useAuthStore } from '@/stores/auth'
const authStore = useAuthStore()
console.log('Is authenticated:', authStore.isAuthenticated)
```

**2. 图像上传后列表不更新**
```javascript
// 确保上传成功后触发列表刷新
window.dispatchEvent(new CustomEvent('image-uploaded', {
  detail: { imageId: response.data.imageId }
}))
```

**3. WebSocket 连接失败**
```javascript
// 检查 WebSocket URL 配置
console.log('WebSocket URL:', import.meta.env.VITE_WS_URL)

// 检查连接状态
import webSocketService from '@/api/websocket'
console.log('WebSocket state:', webSocketService.getState())
```

**4. Element Plus 组件警告**
```vue
<!-- 错误：使用已弃用的 type="text" -->
<el-button type="text">按钮</el-button>

<!-- 正确：使用新的 link 属性 -->
<el-button link>按钮</el-button>
```

**5. 摄像头无法启动**
```javascript
// 检查浏览器权限
if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
  console.error('浏览器不支持摄像头功能')
}

// 处理权限拒绝
try {
  const stream = await navigator.mediaDevices.getUserMedia({ video: true })
} catch (error) {
  if (error.name === 'NotAllowedError') {
    console.error('用户拒绝了摄像头权限')
  }
}
```

**6. 跨域问题**
```javascript
// vite.config.js 配置代理
export default {
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
}
```

### 调试技巧

**开启详细日志**
```javascript
// 在 main.js 中开启调试
if (import.meta.env.DEV) {
  window.__VUE_DEVTOOLS_GLOBAL_HOOK__ = true
  console.log('开发模式已启用')
}
```

**性能监控**
```javascript
// 监控组件渲染性能
import { nextTick } from 'vue'

const startTime = performance.now()
await nextTick()
console.log(`渲染耗时: ${performance.now() - startTime}ms`)
```

## ⚡ 性能优化

### 打包优化

**代码分割**
```javascript
// 路由懒加载
const Dashboard = () => import('@/views/dashboard/index.vue')

// 组件懒加载
const HeavyComponent = defineAsyncComponent(() => 
  import('@/components/HeavyComponent.vue')
)
```

**资源压缩**
```javascript
// vite.config.js
export default {
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['vue', 'vue-router', 'pinia'],
          elementPlus: ['element-plus']
        }
      }
    }
  }
}
```

### 运行时优化

**图片懒加载**
```vue
<template>
  <img 
    v-lazy="imageSrc" 
    alt="图片描述"
    @error="handleImageError"
  />
</template>
```

**虚拟滚动**
```vue
<template>
  <el-virtual-list
    :data="largeDataList"
    :height="400"
    :item-size="50"
  >
    <template #default="{ item }">
      <div>{{ item.name }}</div>
    </template>
  </el-virtual-list>
</template>
```

### 缓存策略

**HTTP 缓存**
```javascript
// 请求拦截器添加缓存控制
request.interceptors.request.use(config => {
  if (config.method === 'get') {
    config.headers['Cache-Control'] = 'max-age=300'
  }
  return config
})
```

**本地缓存**
```javascript
// 缓存 API 响应
const cacheKey = `api_${url}_${JSON.stringify(params)}`
const cached = localStorage.getItem(cacheKey)

if (cached && !isExpired(cached)) {
  return JSON.parse(cached).data
}
```

## 📱 移动端和摄像头支持

### 摄像头功能
系统支持在现代浏览器中直接使用摄像头拍照上传：

- **浏览器支持**: Chrome 90+, Firefox 88+, Safari 14+
- **权限要求**: 需要用户授权摄像头访问
- **功能特性**: 实时预览、拍照、镜像模式、照片上传

### 移动端体验
- **响应式设计**: 完美适配手机和平板
- **触摸友好**: 支持触摸操作和手势
- **摄像头调用**: 移动端可直接调用设备摄像头

### 客户端摄像头功能
CLI客户端也支持摄像头功能：

```bash
# 列出可用摄像头
./cli-client camera list

# 预览摄像头
./cli-client camera preview 0

# 拍照并保存
./cli-client camera capture 0

# 拍照并直接上传
./cli-client camera upload 0

# 录制视频
./cli-client camera record 10 0
```

## 🤝 贡献指南

### 开发流程

1. **Fork 项目**到你的 GitHub 账户
2. **创建功能分支**: `git checkout -b feature/new-feature`
3. **提交更改**: `git commit -am 'Add new feature'`
4. **推送分支**: `git push origin feature/new-feature`
5. **创建 Pull Request**

### 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```bash
# 功能增加
git commit -m "feat: 添加摄像头拍照上传功能"

# 问题修复
git commit -m "fix: 修复WebSocket连接问题"

# 文档更新
git commit -m "docs: 更新摄像头功能文档"

# 样式调整
git commit -m "style: 优化摄像头组件样式"

# 重构代码
git commit -m "refactor: 重构图像列表数据处理"

# 性能优化
git commit -m "perf: 优化图像加载性能"

# 测试相关
git commit -m "test: 添加摄像头组件单元测试"
```

### 代码审查

在提交 PR 前，请确保：

- [ ] 代码通过 ESLint 检查
- [ ] 组件具有适当的 TypeScript 类型（如果使用）
- [ ] 添加了必要的测试用例
- [ ] 更新了相关文档
- [ ] 功能在主流浏览器中测试通过

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 🙋‍♂️ 支持与反馈

- **GitHub Issues**: [提交问题和建议](https://github.com/your-repo/issues)
- **Email**: your-email@domain.com
- **文档**: [在线文档](https://your-docs-site.com)

## 🙏 致谢

感谢以下开源项目：

- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Element Plus](https://element-plus.org/) - Vue 3 组件库
- [Vite](https://vitejs.dev/) - 下一代前端构建工具

---

<div align="center">

**[⬆ 回到顶部](#icpt-图像处理系统---现代化前端)**

Made with ❤️ by ICPT Team

</div> 