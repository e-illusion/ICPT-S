import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { login, register, getUserProfile } from '@/api/auth'
import { removeToken, setToken, getToken } from '@/utils/auth'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const token = ref(getToken())
  const isLoading = ref(false)

  // Getters
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const userInfo = computed(() => user.value)
  const username = computed(() => user.value?.username || '')
  const avatar = computed(() => user.value?.avatar || '')
  const email = computed(() => user.value?.email || '')

  // Actions
  const setUserInfo = (userInfo) => {
    user.value = userInfo
  }

  const setAuthToken = (authToken) => {
    token.value = authToken
    setToken(authToken)
  }

  // Login action
  const loginUser = async (credentials) => {
    try {
      isLoading.value = true
      const response = await login(credentials)
      
      if (response.data) {
        const { token: authToken, user: userInfo } = response.data
        
        setAuthToken(authToken)
        setUserInfo(userInfo)
        
        ElMessage.success('登录成功！')
        return Promise.resolve(response)
      } else {
        throw new Error('登录响应数据格式错误')
      }
    } catch (error) {
      console.error('Login error:', error)
      const message = error.response?.data?.error || error.message || '登录失败'
      ElMessage.error(message)
      return Promise.reject(error)
    } finally {
      isLoading.value = false
    }
  }

  // Register action
  const registerUser = async (userData) => {
    try {
      isLoading.value = true
      const response = await register(userData)
      
      if (response.data) {
        const { token: authToken, user: userInfo } = response.data
        
        setAuthToken(authToken)
        setUserInfo(userInfo)
        
        ElMessage.success('注册成功！')
        return Promise.resolve(response)
      } else {
        throw new Error('注册响应数据格式错误')
      }
    } catch (error) {
      console.error('Register error:', error)
      const message = error.response?.data?.error || error.message || '注册失败'
      ElMessage.error(message)
      return Promise.reject(error)
    } finally {
      isLoading.value = false
    }
  }

  // Get user info action
  const getUserInfo = async () => {
    try {
      if (!token.value) {
        throw new Error('No token available')
      }

      const response = await getUserProfile()
      
      if (response.data) {
        setUserInfo(response.data)
        return Promise.resolve(response.data)
      } else {
        throw new Error('获取用户信息失败')
      }
    } catch (error) {
      console.error('Get user info error:', error)
      
      // If token is invalid, logout
      if (error.response?.status === 401) {
        logout()
      }
      
      return Promise.reject(error)
    }
  }

  // Logout action
  const logout = () => {
    user.value = null
    token.value = null
    removeToken()
    
    // Clear other stores if needed
    // const imageStore = useImageStore()
    // imageStore.$reset()
    
    ElMessage.success('已退出登录')
  }

  // Update user profile
  const updateUserProfile = (updates) => {
    if (user.value) {
      user.value = { ...user.value, ...updates }
    }
  }

  // Reset store
  const $reset = () => {
    user.value = null
    token.value = null
    isLoading.value = false
    removeToken()
  }

  // Auto login check
  const checkAuth = async () => {
    const savedToken = getToken()
    if (savedToken && !user.value) {
      try {
        token.value = savedToken
        await getUserInfo()
      } catch (error) {
        console.error('Auto login failed:', error)
        logout()
      }
    }
  }

  return {
    // State
    user,
    token,
    isLoading,
    
    // Getters
    isAuthenticated,
    userInfo,
    username,
    avatar,
    email,
    
    // Actions
    login: loginUser,        // 添加别名
    register: registerUser,  // 添加别名
    loginUser,
    registerUser,
    getUserInfo,
    logout,
    updateUserProfile,
    checkAuth,
    $reset,
  }
}) 