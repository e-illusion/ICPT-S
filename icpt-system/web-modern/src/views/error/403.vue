<template>
  <div class="error-page">
    <h1>403</h1>
    <h2>访问被禁止</h2>
    <p>抱歉，您没有权限访问此页面。</p>
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
    color: var(--el-color-danger);
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