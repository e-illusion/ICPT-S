<template>
  <div class="register-container">
    <div class="register-card">
      <div class="register-header">
        <div class="logo">
          <el-icon class="logo-icon">
            <Picture />
          </el-icon>
          <h1 class="title">ICPT 图像处理系统</h1>
        </div>
        <p class="subtitle">创建您的账户</p>
      </div>

      <el-form
        ref="registerFormRef"
        :model="registerForm"
        :rules="registerRules"
        class="register-form"
        @submit.prevent="handleSubmit"
      >
        <el-form-item prop="username">
          <el-input
            v-model="registerForm.username"
            placeholder="用户名 (3-20字符)"
            size="large"
            :prefix-icon="User"
            clearable
          />
        </el-form-item>

        <el-form-item prop="email">
          <el-input
            v-model="registerForm.email"
            placeholder="邮箱地址"
            size="large"
            :prefix-icon="Message"
            clearable
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="registerForm.password"
            type="password"
            placeholder="密码 (至少6位)"
            size="large"
            :prefix-icon="Lock"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item prop="confirmPassword">
          <el-input
            v-model="registerForm.confirmPassword"
            type="password"
            placeholder="确认密码"
            size="large"
            :prefix-icon="Lock"
            show-password
            clearable
            @keyup.enter="handleSubmit"
          />
        </el-form-item>

        <el-form-item prop="agreement">
          <el-checkbox v-model="registerForm.agreement">
            我已阅读并同意
            <el-link type="primary" @click="showAgreement">《用户协议》</el-link>
            和
            <el-link type="primary" @click="showPrivacy">《隐私政策》</el-link>
          </el-checkbox>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="register-button"
            :loading="loading"
            @click="handleSubmit"
          >
            {{ loading ? '注册中...' : '立即注册' }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="register-footer">
        <p class="login-tip">
          已有账户？
          <el-link type="primary" @click="handleLogin">立即登录</el-link>
        </p>
      </div>
    </div>

    <!-- Background decoration -->
    <div class="background-decoration">
      <div class="decoration-circle circle-1"></div>
      <div class="decoration-circle circle-2"></div>
      <div class="decoration-circle circle-3"></div>
    </div>

    <!-- User Agreement Dialog -->
    <el-dialog
      v-model="agreementVisible"
      title="用户协议"
      width="600px"
      @close="agreementVisible = false"
    >
      <div class="agreement-content">
        <h3>ICPT图像处理系统用户协议</h3>
        <p>欢迎使用ICPT图像处理系统！请仔细阅读以下条款：</p>
        
        <h4>1. 服务说明</h4>
        <p>本系统提供图像上传、处理、存储等服务，致力于为用户提供高效、安全的图像处理体验。</p>
        
        <h4>2. 用户责任</h4>
        <p>用户应确保上传内容合法合规，不得上传违法、有害或侵犯他人权益的内容。</p>
        
        <h4>3. 隐私保护</h4>
        <p>我们承诺保护用户隐私，不会未经授权分享用户数据。</p>
        
        <h4>4. 服务变更</h4>
        <p>我们保留随时修改、暂停或终止服务的权利，并会提前通知用户。</p>
      </div>
      <template #footer>
        <el-button @click="agreementVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- Privacy Policy Dialog -->
    <el-dialog
      v-model="privacyVisible"
      title="隐私政策"
      width="600px"
      @close="privacyVisible = false"
    >
      <div class="privacy-content">
        <h3>隐私政策</h3>
        <p>本政策说明我们如何收集、使用和保护您的个人信息：</p>
        
        <h4>1. 信息收集</h4>
        <p>我们收集您主动提供的信息（如用户名、邮箱）和使用服务时产生的信息。</p>
        
        <h4>2. 信息使用</h4>
        <p>收集的信息仅用于提供服务、改善用户体验和系统安全维护。</p>
        
        <h4>3. 信息保护</h4>
        <p>我们采用加密存储、访问控制等技术措施保护您的信息安全。</p>
        
        <h4>4. 信息共享</h4>
        <p>除法律要求外，我们不会向第三方分享您的个人信息。</p>
      </div>
      <template #footer>
        <el-button @click="privacyVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Picture, Message } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

// Router and stores
const router = useRouter()
const authStore = useAuthStore()

// Form reference
const registerFormRef = ref(null)

// Loading state
const loading = ref(false)

// Dialog visibility
const agreementVisible = ref(false)
const privacyVisible = ref(false)

// Form data
const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  agreement: false,
})

// Custom validator for password confirmation
const validatePasswordConfirm = (rule, value, callback) => {
  if (value !== registerForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

// Custom validator for email
const validateEmail = (rule, value, callback) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(value)) {
    callback(new Error('请输入有效的邮箱地址'))
  } else {
    callback()
  }
}

// Custom validator for username
const validateUsername = (rule, value, callback) => {
  const usernameRegex = /^[a-zA-Z0-9_-]+$/
  if (!usernameRegex.test(value)) {
    callback(new Error('用户名只能包含字母、数字、下划线和连字符'))
  } else {
    callback()
  }
}

// Form validation rules
const registerRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' },
    { validator: validateUsername, trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { validator: validateEmail, trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validatePasswordConfirm, trigger: 'blur' },
  ],
  agreement: [
    { 
      type: 'boolean', 
      required: true, 
      message: '请同意用户协议和隐私政策', 
      trigger: 'change',
      validator: (rule, value, callback) => {
        if (!value) {
          callback(new Error('请同意用户协议和隐私政策'))
        } else {
          callback()
        }
      }
    },
  ],
}

// Handle form submission
const handleSubmit = async () => {
  if (!registerFormRef.value) return

  try {
    await registerFormRef.value.validate()
    loading.value = true

    const registerData = {
      username: registerForm.username,
      email: registerForm.email,
      password: registerForm.password,
    }

    await authStore.registerUser(registerData)

    ElMessage.success('注册成功！正在跳转到仪表盘...')
    
    // Redirect to dashboard
    setTimeout(() => {
      router.push('/dashboard')
    }, 1000)
    
  } catch (error) {
    console.error('Register error:', error)
    
    if (error.response?.status === 409) {
      ElMessage.error('用户名或邮箱已存在，请使用其他用户名或邮箱')
    } else if (error.response?.status === 400) {
      ElMessage.error('请检查输入信息是否正确')
    } else {
      ElMessage.error(error.message || '注册失败，请检查网络连接')
    }
  } finally {
    loading.value = false
  }
}

// Handle login navigation
const handleLogin = () => {
  router.push('/login')
}

// Show agreement dialog
const showAgreement = () => {
  agreementVisible.value = true
}

// Show privacy dialog
const showPrivacy = () => {
  privacyVisible.value = true
}
</script>

<style lang="scss" scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

.register-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 48px 40px;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 2;
}

.register-header {
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

.register-form {
  .register-button {
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
    line-height: 1.4;
  }

  :deep(.el-form-item) {
    margin-bottom: 20px;
  }
}

.register-footer {
  text-align: center;
  margin-top: 24px;

  .login-tip {
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

.agreement-content,
.privacy-content {
  max-height: 400px;
  overflow-y: auto;
  
  h3 {
    color: #2c3e50;
    margin-bottom: 16px;
  }

  h4 {
    color: #34495e;
    margin: 16px 0 8px 0;
  }

  p {
    color: #7f8c8d;
    line-height: 1.6;
    margin-bottom: 12px;
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
  .register-card {
    background: rgba(17, 24, 39, 0.95);
    
    .register-header {
      .title {
        color: #f8fafc;
      }

      .subtitle {
        color: #94a3b8;
      }
    }

    .register-footer .login-tip {
      color: #94a3b8;
    }
  }

  .agreement-content,
  .privacy-content {
    h3 {
      color: #f8fafc;
    }

    h4 {
      color: #e2e8f0;
    }

    p {
      color: #94a3b8;
    }
  }
}

// Mobile responsive
@media (max-width: 480px) {
  .register-card {
    margin: 20px;
    padding: 32px 24px;
  }

  .register-header .logo .title {
    font-size: 20px;
  }
}
</style>
