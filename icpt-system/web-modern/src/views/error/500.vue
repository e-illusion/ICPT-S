<template>
  <div class="error-page">
    <h1>500</h1>
    <h2>服务器错误</h2>
    <p>抱歉，服务器出现了错误，请稍后再试。</p>
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
    color: var(--el-color-warning);
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