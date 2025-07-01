<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <div class="logo">
          <el-icon class="logo-icon">
            <Picture />
          </el-icon>
          <h1 class="title">ICPT 图像处理系统</h1>
        </div>
        <p class="subtitle">请登录您的账户</p>
      </div>

      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @submit.prevent="handleSubmit"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="用户名"
            size="large"
            :prefix-icon="User"
            clearable
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            size="large"
            :prefix-icon="Lock"
            show-password
            clearable
            @keyup.enter="handleSubmit"
          />
        </el-form-item>

        <el-form-item>
          <el-checkbox v-model="loginForm.remember">
            记住登录状态
          </el-checkbox>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="login-button"
            :loading="loading"
            @click="handleSubmit"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="login-footer">
        <p class="register-tip">
          还没有账户？
          <el-link type="primary" @click="handleRegister">立即注册</el-link>
        </p>
      </div>
    </div>

    <!-- Background decoration -->
    <div class="background-decoration">
      <div class="decoration-circle circle-1"></div>
      <div class="decoration-circle circle-2"></div>
      <div class="decoration-circle circle-3"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Picture } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

// Router and stores
const router = useRouter()
const authStore = useAuthStore()

// Form reference
const loginFormRef = ref(null)

// Loading state
const loading = ref(false)

// Form data
const loginForm = reactive({
  username: '',
  password: '',
  remember: false,
})

// Form validation rules
const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' },
  ],
}

// Handle form submission
const handleSubmit = async () => {
  if (!loginFormRef.value) return

  try {
    await loginFormRef.value.validate()
    loading.value = true

    const loginData = {
      username: loginForm.username,
      password: loginForm.password,
    }

    await authStore.login(loginData)

    ElMessage.success('登录成功！')
    
    // Redirect to dashboard or previous page
    const redirect = router.currentRoute.value.query.redirect || '/dashboard'
    await router.push(redirect)
    
  } catch (error) {
    console.error('Login error:', error)
    
    if (error.response?.status === 401) {
      ElMessage.error('用户名或密码错误')
    } else if (error.response?.status === 429) {
      ElMessage.error('登录尝试过于频繁，请稍后再试')
    } else {
      ElMessage.error(error.message || '登录失败，请检查网络连接')
    }
  } finally {
    loading.value = false
  }
}

// Handle register
const handleRegister = () => {
  router.push('/register')
}

// Initialize default credentials in development
onMounted(() => {
  if (import.meta.env.DEV) {
    loginForm.username = 'admin'
    loginForm.password = 'admin123'
  }
})
</script>

<style lang="scss" scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

.login-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 48px 40px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 2;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;

  .logo {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    margin-bottom: 16px;

    .logo-icon {
      font-size: 32px;
      color: var(--el-color-primary);
    }

    .title {
      color: #2c3e50;
      font-size: 24px;
      font-weight: 600;
      margin: 0;
    }
  }

  .subtitle {
    color: #7f8c8d;
    font-size: 14px;
    margin: 0;
  }
}

.login-form {
  .login-button {
    width: 100%;
    height: 44px;
    font-size: 16px;
    border-radius: 12px;
  }

  :deep(.el-input__wrapper) {
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  :deep(.el-checkbox__label) {
    color: #7f8c8d;
    font-size: 14px;
  }
}

.login-footer {
  text-align: center;
  margin-top: 24px;

  .register-tip {
    color: #7f8c8d;
    font-size: 14px;
    margin: 0;
  }
}

.background-decoration {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 1;
}

.decoration-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 6s ease-in-out infinite;

  &.circle-1 {
    width: 200px;
    height: 200px;
    top: 10%;
    left: 10%;
    animation-delay: 0s;
  }

  &.circle-2 {
    width: 150px;
    height: 150px;
    top: 60%;
    right: 10%;
    animation-delay: 2s;
  }

  &.circle-3 {
    width: 100px;
    height: 100px;
    bottom: 20%;
    left: 20%;
    animation-delay: 4s;
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-20px);
  }
}

// Dark theme support
[data-theme='dark'] {
  .login-card {
    background: rgba(17, 24, 39, 0.95);
    
    .login-header {
      .title {
        color: #f8fafc;
      }

      .subtitle {
        color: #94a3b8;
      }
    }

    .login-footer .register-tip {
      color: #94a3b8;
    }
  }
}

// Mobile responsive
@media (max-width: 480px) {
  .login-card {
    margin: 20px;
    padding: 32px 24px;
  }

  .login-header .logo .title {
    font-size: 20px;
  }
}
</style> 