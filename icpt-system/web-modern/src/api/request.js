import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getToken, removeToken } from '@/utils/auth'
import router from '@/router'
import NProgress from 'nprogress'

// Create axios instance
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 15000,
  withCredentials: false,
})

// Request interceptor
request.interceptors.request.use(
  (config) => {
    // Start loading progress
    NProgress.start()
    
    // Add auth token to headers
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    // Add timestamp to prevent caching for GET requests
    if (config.method === 'get') {
      config.params = {
        ...config.params,
        _t: Date.now(),
      }
    }
    
    return config
  },
  (error) => {
    console.error('Request error:', error)
    NProgress.done()
    return Promise.reject(error)
  }
)

// Response interceptor
request.interceptors.response.use(
  (response) => {
    NProgress.done()
    
    // Extract data from response
    const { data, status } = response
    
    // Handle different response formats
    if (status === 200 || status === 201 || status === 202) {
      return data
    } else {
      console.error('Unexpected response status:', status)
      return Promise.reject(new Error(`HTTP ${status}`))
    }
  },
  async (error) => {
    NProgress.done()
    
    const { response, request, message } = error
    
    // Network error
    if (!response) {
      if (request) {
        ElMessage.error('网络连接失败，请检查网络设置')
      } else {
        ElMessage.error('请求配置错误')
      }
      return Promise.reject(error)
    }
    
    const { status, data } = response
    
    // Handle different error status codes
    switch (status) {
      case 400:
        ElMessage.error(data?.error || '请求参数错误')
        break
        
      case 401:
        // Unauthorized - clear token and redirect to login
        ElMessage.error('登录已过期，请重新登录')
        removeToken()
        
        // Avoid multiple redirects
        if (router.currentRoute.value.path !== '/login') {
          router.replace('/login')
        }
        break
        
      case 403:
        ElMessage.error('访问被禁止，权限不足')
        router.replace('/403')
        break
        
      case 404:
        ElMessage.error('请求的资源不存在')
        break
        
      case 409:
        ElMessage.error(data?.error || '资源冲突')
        break
        
      case 422:
        // Validation errors
        if (data?.errors) {
          const errorMessages = Object.values(data.errors).flat()
          ElMessage.error(errorMessages.join(', '))
        } else {
          ElMessage.error(data?.error || '数据验证失败')
        }
        break
        
      case 429:
        ElMessage.error('请求过于频繁，请稍后再试')
        break
        
      case 500:
        ElMessage.error('服务器内部错误')
        router.replace('/500')
        break
        
      case 502:
      case 503:
      case 504:
        ElMessage.error('服务器暂时不可用，请稍后再试')
        break
        
      default:
        ElMessage.error(data?.error || `请求失败 (${status})`)
    }
    
    return Promise.reject(error)
  }
)

// Request helpers
export const get = (url, params = {}, config = {}) => {
  return request.get(url, { params, ...config })
}

export const post = (url, data = {}, config = {}) => {
  return request.post(url, data, config)
}

export const put = (url, data = {}, config = {}) => {
  return request.put(url, data, config)
}

export const patch = (url, data = {}, config = {}) => {
  return request.patch(url, data, config)
}

export const del = (url, config = {}) => {
  return request.delete(url, config)
}

// Upload helper
export const upload = (url, formData, config = {}) => {
  const uploadConfig = {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    onUploadProgress: (progressEvent) => {
      if (config.onProgress) {
        const percentCompleted = Math.round(
          (progressEvent.loaded * 100) / progressEvent.total
        )
        config.onProgress(percentCompleted)
      }
    },
    ...config,
  }
  
  return request.post(url, formData, uploadConfig)
}

// Download helper
export const download = async (url, filename, config = {}) => {
  try {
    const response = await request.get(url, {
      responseType: 'blob',
      ...config,
    })
    
    // Create blob URL
    const blob = new Blob([response])
    const downloadUrl = window.URL.createObjectURL(blob)
    
    // Create download link
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = filename || 'download'
    document.body.appendChild(link)
    link.click()
    
    // Cleanup
    document.body.removeChild(link)
    window.URL.revokeObjectURL(downloadUrl)
    
    return response
  } catch (error) {
    console.error('Download error:', error)
    throw error
  }
}

export default request 