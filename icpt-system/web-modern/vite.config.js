import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import { resolve } from 'path'
import fs from 'fs'

export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [ElementPlusResolver()],
      imports: [
        'vue',
        'vue-router',
        'pinia'
      ],
      dts: true,
    }),
    Components({
      resolvers: [ElementPlusResolver()],
      dts: true,
    }),
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 3000,
    https: {
      key: fs.readFileSync(resolve(__dirname, 'cert/key.pem')),
      cert: fs.readFileSync(resolve(__dirname, 'cert/cert.pem')),
    },
    proxy: {
      '/api': {
        target: 'https://114.55.58.3:8080',
        changeOrigin: true,
        secure: false, // 允许自签名证书
        ws: true, // 支持WebSocket代理
        configure: (proxy, options) => {
          // 自定义代理配置
          proxy.on('error', (err, req, res) => {
            console.log('Proxy error:', err);
          });
          proxy.on('proxyReq', (proxyReq, req, res) => {
            console.log('Proxying request:', req.method, req.url);
          });
        },
      },
      '/static': {
        target: 'https://114.55.58.3:8080',
        changeOrigin: true,
        secure: false, // 允许自签名证书
      },
      // WebSocket代理配置 - 专门处理WebSocket
      '/api/v1/ws': {
        target: 'wss://114.55.58.3:8080',
        ws: true,
        secure: false,
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    sourcemap: false,
    rollupOptions: {
      output: {
        manualChunks: {
          'element-plus': ['element-plus'],
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
          'charts': ['echarts', 'vue-echarts'],
        },
      },
    },
  },
  // 优化开发体验
  optimizeDeps: {
    include: ['element-plus', 'vue', 'vue-router', 'pinia'],
  },
  // 定义全局常量
  define: {
    __VUE_OPTIONS_API__: true,
    __VUE_PROD_DEVTOOLS__: false,
  },
}) 