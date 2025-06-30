<template>
  <div id="app" :data-theme="theme">
    <!-- Global Loading Bar -->
    <div v-if="isGlobalLoading" class="global-loading">
      <div class="loading-content">
        <div class="loading-spinner"></div>
        <p class="loading-text">{{ loadingText }}</p>
      </div>
    </div>

    <!-- Router View -->
    <router-view v-slot="{ Component, route }">
      <transition name="fade" mode="out-in">
        <component :is="Component" :key="route.path" />
      </transition>
    </router-view>

    <!-- Global Message Container -->
    <teleport to="body">
      <div id="message-container"></div>
    </teleport>

    <!-- Global Modal Container -->
    <teleport to="body">
      <div id="modal-container"></div>
    </teleport>

    <!-- WebSocket Connection Status -->
    <div 
      v-if="showConnectionStatus" 
      class="connection-status"
      :class="connectionStatusClass"
    >
      <el-icon>
        <component :is="connectionIcon" />
      </el-icon>
      <span>{{ connectionStatusText }}</span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, getCurrentInstance } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import { Connection, Close, Loading } from '@element-plus/icons-vue'

// Stores
const authStore = useAuthStore()
const router = useRouter()

// Reactive state
const theme = ref('light')
const isGlobalLoading = ref(false)
const loadingText = ref('正在加载...')
const wsConnection = ref(null)
const isConnected = ref(false)
const showConnectionStatus = ref(false)
const retryCount = ref(0)
const maxRetries = 3

// Computed properties
const connectionStatusClass = computed(() => ({
  'connection-status--connected': isConnected.value,
  'connection-status--disconnected': !isConnected.value,
  'connection-status--visible': showConnectionStatus.value,
}))

const connectionIcon = computed(() => 
  isConnected.value ? Connection : Close
)

const connectionStatusText = computed(() => 
  isConnected.value ? 'WebSocket 已连接' : 'WebSocket 连接断开'
)

// Global loading methods
const setGlobalLoading = (loading, text = '正在加载...') => {
  isGlobalLoading.value = loading
  loadingText.value = text
}

// Theme management
const toggleTheme = () => {
  theme.value = theme.value === 'light' ? 'dark' : 'light'
  localStorage.setItem('theme', theme.value)
  document.documentElement.setAttribute('data-theme', theme.value)
  
  // Toggle Element Plus dark class
  if (theme.value === 'dark') {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

const initTheme = () => {
  const savedTheme = localStorage.getItem('theme') || 'light'
  theme.value = savedTheme
  document.documentElement.setAttribute('data-theme', savedTheme)
  
  // Apply Element Plus dark class
  if (savedTheme === 'dark') {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

// WebSocket management
const connectWebSocket = async () => {
  if (!authStore.isAuthenticated) return

  try {
    const webSocketService = (await import('@/api/websocket')).default
    await webSocketService.connect()
    
    isConnected.value = true
    showConnectionStatus.value = true
    
    // Hide status after 3 seconds if connected
    setTimeout(() => {
      if (isConnected.value) {
        showConnectionStatus.value = false
      }
    }, 3000)

    // Setup event listeners
    webSocketService.on('connected', () => {
      isConnected.value = true
      retryCount.value = 0
    })

    webSocketService.on('disconnected', () => {
      isConnected.value = false
      showConnectionStatus.value = true
    })

    webSocketService.on('error', (error) => {
      console.error('WebSocket error:', error)
      isConnected.value = false
      showConnectionStatus.value = true
    })

  } catch (error) {
    console.error('WebSocket connection error:', error)
  }
}

const disconnectWebSocket = async () => {
  try {
    const webSocketService = (await import('@/api/websocket')).default
    webSocketService.disconnect()
  } catch (error) {
    console.error('WebSocket disconnect error:', error)
  }
  
  isConnected.value = false
  showConnectionStatus.value = false
}

// Handle WebSocket messages (now handled by individual pages)
// This function is kept for compatibility but individual pages handle their own events

// Global error handler
const handleGlobalError = (error, instance, info) => {
  console.error('Global error:', error, info)
  
  ElMessage.error({
    message: '系统发生错误，请刷新页面重试',
    duration: 5000,
    showClose: true,
  })
}

// Global warning handler
const handleGlobalWarning = (msg, instance, trace) => {
  if (process.env.NODE_ENV === 'development') {
    console.warn('Global warning:', msg, trace)
  }
}

// Watch for authentication changes
watch(
  () => authStore.isAuthenticated,
  (isAuth) => {
    if (isAuth) {
      connectWebSocket()
    } else {
      disconnectWebSocket()
    }
  },
  { immediate: true }
)

// Lifecycle hooks
onMounted(async () => {
  // Initialize theme
  initTheme()
  
  // Check authentication
  try {
    setGlobalLoading(true, '检查登录状态...')
    await authStore.checkAuth()
  } catch (error) {
    console.error('Auth check failed:', error)
  } finally {
    setGlobalLoading(false)
  }
  
  // Set up global error handlers
  const app = getCurrentInstance()?.appContext.app
  if (app) {
    app.config.errorHandler = handleGlobalError
    app.config.warnHandler = handleGlobalWarning
  }
  
  // Global keyboard shortcuts
  document.addEventListener('keydown', handleKeyboardShortcuts)
})

onUnmounted(() => {
  disconnectWebSocket()
  document.removeEventListener('keydown', handleKeyboardShortcuts)
})

// Keyboard shortcuts
const handleKeyboardShortcuts = (event) => {
  // Ctrl/Cmd + K for search
  if ((event.ctrlKey || event.metaKey) && event.key === 'k') {
    event.preventDefault()
    // Open search modal or focus search input
    // This can be implemented later
  }
  
  // Ctrl/Cmd + \ for theme toggle
  if ((event.ctrlKey || event.metaKey) && event.key === '\\') {
    event.preventDefault()
    toggleTheme()
  }
}

// Expose methods to global scope for use in other components
defineExpose({
  setGlobalLoading,
  toggleTheme,
  connectWebSocket,
  disconnectWebSocket,
})
</script>

<style lang="scss" scoped>
#app {
  min-height: 100vh;
  position: relative;
}

.global-loading {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  
  .loading-content {
    text-align: center;
    
    .loading-spinner {
      width: 40px;
      height: 40px;
      border: 3px solid var(--color-gray-200);
      border-top: 3px solid var(--color-primary);
      border-radius: 50%;
      animation: spin 1s linear infinite;
      margin: 0 auto 16px;
    }
    
    .loading-text {
      color: var(--color-gray-600);
      font-size: 14px;
      margin: 0;
    }
  }
}

.connection-status {
  position: fixed;
  bottom: 20px;
  right: 20px;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
  z-index: 1000;
  transform: translateY(100px);
  opacity: 0;
  transition: all 0.3s ease;
  
  &--visible {
    transform: translateY(0);
    opacity: 1;
  }
  
  &--connected {
    background: var(--color-success);
    color: white;
  }
  
  &--disconnected {
    background: var(--color-danger);
    color: white;
  }
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

// Route transition animations
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

// Dark theme adjustments
[data-theme='dark'] {
  .global-loading {
    background: rgba(17, 24, 39, 0.9);
    
    .loading-text {
      color: var(--color-gray-400);
    }
  }
}

// Global CSS Variables for Dark Theme
:root[data-theme='dark'] {
  --el-bg-color: #1a1a1a;
  --el-bg-color-page: #141414;
  --el-bg-color-overlay: #1d1e1f;
  --el-text-color-primary: #e5eaf3;
  --el-text-color-regular: #cfd3dc;
  --el-text-color-secondary: #a3a6ad;
  --el-text-color-placeholder: #8d9095;
  --el-text-color-disabled: #6c6e72;
  --el-border-color: #4c4d4f;
  --el-border-color-light: #414243;
  --el-border-color-lighter: #363637;
  --el-border-color-extra-light: #2b2b2c;
  --el-border-color-dark: #58585b;
  --el-border-color-darker: #636466;
  --el-fill-color: #303133;
  --el-fill-color-light: #262727;
  --el-fill-color-lighter: #1d1d1d;
  --el-fill-color-extra-light: #191919;
  --el-fill-color-dark: #39393a;
  --el-fill-color-darker: #424243;
  --el-fill-color-blank: transparent;
  --el-box-shadow: 0px 12px 32px 4px rgba(0, 0, 0, 0.36), 0px 8px 20px rgba(0, 0, 0, 0.72);
  --el-box-shadow-light: 0px 0px 12px rgba(0, 0, 0, 0.72);
  --el-box-shadow-lighter: 0px 0px 6px rgba(0, 0, 0, 0.72);
  --el-box-shadow-dark: 0px 16px 48px 16px rgba(0, 0, 0, 0.72), 0px 12px 32px rgba(0, 0, 0, 0.72), 0px 8px 16px -8px rgba(0, 0, 0, 0.72);
}

// Light theme variables (explicit)
:root[data-theme='light'] {
  --el-bg-color: #ffffff;
  --el-bg-color-page: #f2f3f5;
  --el-bg-color-overlay: #ffffff;
  --el-text-color-primary: #303133;
  --el-text-color-regular: #606266;
  --el-text-color-secondary: #909399;
  --el-text-color-placeholder: #a8abb2;
  --el-text-color-disabled: #c0c4cc;
  --el-border-color: #dcdfe6;
  --el-border-color-light: #e4e7ed;
  --el-border-color-lighter: #ebeef5;
  --el-border-color-extra-light: #f2f6fc;
  --el-border-color-dark: #d4d7de;
  --el-border-color-darker: #cdd0d6;
  --el-fill-color: #f0f2f5;
  --el-fill-color-light: #f5f7fa;
  --el-fill-color-lighter: #fafafa;
  --el-fill-color-extra-light: #fafcff;
  --el-fill-color-dark: #ebedf0;
  --el-fill-color-darker: #e6e8eb;
  --el-fill-color-blank: #ffffff;
}
</style> 