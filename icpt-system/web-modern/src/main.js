import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

import App from './App.vue'
import router from './router'
import './styles/index.scss'

// Configure NProgress
NProgress.configure({ 
  showSpinner: false,
  minimum: 0.1,
  easing: 'ease',
  speed: 500
})

const app = createApp(App)
const pinia = createPinia()

// Register Element Plus with Chinese locale
app.use(ElementPlus, {
  locale: zhCn,
})

// Register all Element Plus icons
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// Use plugins
app.use(pinia)
app.use(router)

// Global error handler
app.config.errorHandler = (err, vm, info) => {
  console.error('Vue Error:', err, info)
  // You can integrate with error reporting service here
}

// Global warning handler
app.config.warnHandler = (msg, vm, trace) => {
  console.warn('Vue Warning:', msg, trace)
}

// Global properties
app.config.globalProperties.$APP_NAME = 'ICPT 图像处理系统'
app.config.globalProperties.$APP_VERSION = '2.0.0'

// Mount the app
app.mount('#app')

// Hide loading screen after app is mounted
nextTick(() => {
  const loadingScreen = document.getElementById('loading-screen')
  if (loadingScreen) {
    document.body.classList.add('app-loaded')
    setTimeout(() => {
      loadingScreen.remove()
    }, 300)
  }
})

export default app 