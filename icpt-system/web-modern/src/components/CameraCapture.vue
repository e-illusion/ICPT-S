<template>
  <div class="camera-capture">
    <el-card class="camera-card">
      <template #header>
        <div class="camera-header">
          <span class="camera-title">
            <el-icon><Camera /></el-icon>
            摄像头拍照
          </span>
          <el-button 
            v-if="isActive" 
            @click="stopCamera" 
            type="danger" 
            size="small"
          >
            关闭摄像头
          </el-button>
        </div>
      </template>

      <div class="camera-content">
        <!-- 摄像头视频流 -->
        <div class="video-container" v-show="isActive && !capturedImage">
          <video 
            ref="videoRef" 
            autoplay 
            playsinline
            class="camera-video"
            :class="{ 'mirrored': mirrorVideo }"
          ></video>
          <div class="video-overlay">
            <div class="capture-controls">
              <el-button-group>
                <el-button 
                  @click="capturePhoto" 
                  type="primary" 
                  size="large"
                  :loading="capturing"
                >
                  <el-icon><Camera /></el-icon>
                  拍照
                </el-button>
                <el-button 
                  @click="toggleMirror" 
                  size="large"
                >
                  <el-icon><RefreshLeft /></el-icon>
                  镜像
                </el-button>
              </el-button-group>
            </div>
          </div>
        </div>

        <!-- 拍摄的照片预览 -->
        <div class="photo-preview" v-if="capturedImage">
          <img :src="capturedImage" alt="拍摄的照片" class="captured-photo" />
          <div class="photo-actions">
            <el-button-group>
              <el-button @click="retakePhoto" size="large">
                <el-icon><RefreshRight /></el-icon>
                重新拍照
              </el-button>
              <el-button 
                @click="uploadPhoto" 
                type="primary" 
                size="large"
                :loading="uploading"
              >
                <el-icon><Upload /></el-icon>
                上传照片
              </el-button>
            </el-button-group>
          </div>
        </div>

        <!-- 启动摄像头按钮 -->
        <div class="camera-start" v-if="!isActive && !capturedImage">
          <el-button 
            @click="startCamera" 
            type="primary" 
            size="large"
            :loading="starting"
          >
            <el-icon><VideoCamera /></el-icon>
            启动摄像头
          </el-button>
          <p class="camera-tip">
            点击按钮启动摄像头进行拍照上传
          </p>
        </div>

        <!-- 错误提示 -->
        <div class="camera-error" v-if="errorMessage">
          <el-alert
            :title="errorMessage"
            type="error"
            show-icon
            :closable="false"
          />
          <el-button @click="retryCamera" style="margin-top: 12px;">
            重试
          </el-button>
        </div>
      </div>

      <!-- 隐藏的canvas用于拍照 -->
      <canvas ref="canvasRef" style="display: none;"></canvas>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElNotification } from 'element-plus'
import { Camera, VideoCamera, Upload, RefreshLeft, RefreshRight } from '@element-plus/icons-vue'
import { uploadImage } from '@/api/images'

// Props
const props = defineProps({
  quality: {
    type: Number,
    default: 0.8
  },
  width: {
    type: Number,
    default: 800
  },
  height: {
    type: Number,
    default: 600
  }
})

// Emits
const emit = defineEmits(['uploaded', 'error'])

// Refs
const videoRef = ref(null)
const canvasRef = ref(null)

// State
const isActive = ref(false)
const starting = ref(false)
const capturing = ref(false)
const uploading = ref(false)
const capturedImage = ref(null)
const errorMessage = ref('')
const mirrorVideo = ref(true)
const mediaStream = ref(null)

// Methods
const startCamera = async () => {
  starting.value = true
  errorMessage.value = ''

  try {
    // 检查浏览器支持
    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
      throw new Error('摄像头功能需要安全的HTTPS环境或在localhost下运行')
    }

    // 请求摄像头权限
    const constraints = {
      video: {
        width: { ideal: props.width },
        height: { ideal: props.height },
        facingMode: 'user' // 前置摄像头
      },
      audio: false
    }

    const stream = await navigator.mediaDevices.getUserMedia(constraints)
    mediaStream.value = stream
    
    if (videoRef.value) {
      videoRef.value.srcObject = stream
      await new Promise((resolve) => {
        videoRef.value.onloadedmetadata = resolve
      })
    }

    isActive.value = true
    ElMessage.success('摄像头启动成功')
  } catch (error) {
    console.error('启动摄像头失败:', error)
    
    let message = '启动摄像头失败'
    if (error.name === 'NotAllowedError') {
      message = '请允许访问摄像头权限'
    } else if (error.name === 'NotFoundError') {
      message = '未检测到摄像头设备'
    } else if (error.name === 'NotReadableError') {
      message = '摄像头被其他应用占用'
    } else if (error.message) {
      message = error.message
    }
    
    errorMessage.value = message
    emit('error', error)
  } finally {
    starting.value = false
  }
}

const stopCamera = () => {
  if (mediaStream.value) {
    mediaStream.value.getTracks().forEach(track => track.stop())
    mediaStream.value = null
  }
  isActive.value = false
  capturedImage.value = null
  errorMessage.value = ''
}

const capturePhoto = () => {
  if (!videoRef.value || !canvasRef.value) {
    ElMessage.error('摄像头未就绪')
    return
  }

  capturing.value = true

  try {
    const video = videoRef.value
    const canvas = canvasRef.value
    const context = canvas.getContext('2d')

    // 设置canvas尺寸
    canvas.width = video.videoWidth
    canvas.height = video.videoHeight

    // 如果是镜像模式，翻转画布
    if (mirrorVideo.value) {
      context.scale(-1, 1)
      context.drawImage(video, -canvas.width, 0, canvas.width, canvas.height)
    } else {
      context.drawImage(video, 0, 0, canvas.width, canvas.height)
    }

    // 获取图片数据
    capturedImage.value = canvas.toDataURL('image/jpeg', props.quality)
    
    ElMessage.success('拍照成功')
  } catch (error) {
    console.error('拍照失败:', error)
    ElMessage.error('拍照失败')
  } finally {
    capturing.value = false
  }
}

const retakePhoto = () => {
  capturedImage.value = null
}

const uploadPhoto = async () => {
  if (!capturedImage.value) {
    ElMessage.error('没有可上传的照片')
    return
  }

  uploading.value = true

  try {
    // 将base64转换为Blob
    const response = await fetch(capturedImage.value)
    const blob = await response.blob()
    
    // 创建File对象
    const fileName = `camera_${Date.now()}.jpg`
    const file = new File([blob], fileName, { type: 'image/jpeg' })

    // 上传文件
    const uploadResponse = await uploadImage(file, (progress) => {
      console.log(`上传进度: ${progress}%`)
    })

    ElNotification.success({
      title: '上传成功',
      message: `照片 ${fileName} 已上传成功`,
    })

    // 触发上传成功事件
    emit('uploaded', {
      file: file,
      response: uploadResponse,
      imageId: uploadResponse.data?.imageId
    })

    // 重置状态
    capturedImage.value = null
    
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('照片上传失败')
    emit('error', error)
  } finally {
    uploading.value = false
  }
}

const toggleMirror = () => {
  mirrorVideo.value = !mirrorVideo.value
}

const retryCamera = () => {
  errorMessage.value = ''
  startCamera()
}

// Lifecycle
onMounted(() => {
  // 自动检查摄像头权限
})

onUnmounted(() => {
  stopCamera()
})

// 暴露方法给父组件
defineExpose({
  startCamera,
  stopCamera,
  capturePhoto,
  isActive
})
</script>

<style lang="scss" scoped>
.camera-capture {
  width: 100%;
}

.camera-card {
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
}

.camera-header {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .camera-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 500;
  }
}

.camera-content {
  min-height: 300px;
}

.video-container {
  position: relative;
  width: 100%;
  border-radius: 8px;
  overflow: hidden;
  background: #000;
}

.camera-video {
  width: 100%;
  height: auto;
  display: block;
  min-height: 300px;

  &.mirrored {
    transform: scaleX(-1);
  }
}

.video-overlay {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);

  .capture-controls {
    background: rgba(0, 0, 0, 0.7);
    padding: 12px;
    border-radius: 8px;
    backdrop-filter: blur(4px);
  }
}

.photo-preview {
  text-align: center;

  .captured-photo {
    width: 100%;
    max-height: 400px;
    object-fit: contain;
    border-radius: 8px;
    border: 1px solid var(--el-border-color);
  }

  .photo-actions {
    margin-top: 16px;
  }
}

.camera-start {
  text-align: center;
  padding: 40px 20px;

  .camera-tip {
    margin-top: 12px;
    color: var(--el-text-color-regular);
    font-size: 14px;
  }
}

.camera-error {
  padding: 20px;
  text-align: center;
}

// 响应式设计
@media (max-width: 768px) {
  .camera-card {
    max-width: 100%;
  }

  .camera-video {
    min-height: 250px;
  }

  .video-overlay {
    bottom: 10px;

    .capture-controls {
      padding: 8px;
    }
  }

  .capture-controls :deep(.el-button) {
    padding: 8px 12px;
  }
}
</style> 