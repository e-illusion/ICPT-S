<template>
  <div class="error-page">
    <h1>404</h1>
    <h2>页面未找到</h2>
    <p>抱歉，您访问的页面不存在。</p>
    <el-button type="primary" @click="goHome">返回首页</el-button>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<style lang="scss" scoped>
.error-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  
  h1 {
    font-size: 72px;
    margin: 0;
    color: var(--el-color-primary);
  }
  
  h2 {
    font-size: 24px;
    margin: 16px 0 8px;
  }
  
  p {
    color: var(--el-text-color-regular);
    margin-bottom: 24px;
  }
}
</style> 