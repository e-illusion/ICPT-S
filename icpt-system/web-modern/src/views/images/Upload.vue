<template>
  <div class="upload-page">
    <div class="page-header">
      <h1 class="page-title">图像上传</h1>
      <p class="page-subtitle">支持拖拽上传、批量上传，实时查看处理进度</p>
    </div>

    <!-- Upload Area -->
    <div class="upload-area">
      <el-card class="upload-card">
        <el-tabs v-model="activeTab" class="upload-tabs">
          <!-- 文件上传标签页 -->
          <el-tab-pane label="文件上传" name="file">
            <div
          class="upload-zone"
          :class="{
            'upload-zone--dragover': isDragOver,
            'upload-zone--uploading': isUploading,
          }"
          @drop="handleDrop"
          @dragover="handleDragOver"
          @dragenter="handleDragEnter"
          @dragleave="handleDragLeave"
          @click="triggerFileSelect"
        >
          <div class="upload-content">
            <el-icon class="upload-icon">
              <UploadFilled v-if="!isUploading" />
              <Loading v-else />
            </el-icon>
            
            <h3 class="upload-title">
              {{ isUploading ? '正在上传...' : '点击或拖拽文件到这里上传' }}
            </h3>
            
            <p class="upload-hint">
              支持 JPG、PNG、GIF、WEBP 格式，单个文件不超过 32MB
            </p>

            <!-- Progress Bar -->
            <div v-if="isUploading" class="upload-progress">
              <el-progress
                :percentage="uploadProgress"
                :stroke-width="8"
                status="success"
              />
              <p class="progress-text">
                已上传 {{ uploadedCount }}/{{ totalFiles }} 个文件
              </p>
            </div>
          </div>

          <!-- Hidden File Input -->
          <input
            ref="fileInputRef"
            type="file"
            multiple
            accept="image/*"
            style="display: none"
            @change="handleFileSelect"
          />
            </div>
          </el-tab-pane>
          
          <!-- 摄像头拍照标签页 -->
          <el-tab-pane name="camera">
            <template #label>
              <span class="tab-label">
                <el-icon><VideoCamera /></el-icon>
                摄像头拍照
              </span>
            </template>
            <CameraCapture 
              @uploaded="handleCameraUpload"
              @error="handleCameraError"
            />
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </div>

    <!-- Upload Queue -->
    <div v-if="uploadQueue.length > 0" class="upload-queue">
      <h2 class="section-title">上传队列</h2>
      <el-card>
        <div class="queue-list">
          <div
            v-for="item in uploadQueue"
            :key="item.id"
            class="queue-item"
          >
            <div class="item-info">
              <div class="item-icon">
                <el-icon v-if="item.status === 'pending'">
                  <Clock />
                </el-icon>
                <el-icon v-else-if="item.status === 'uploading'">
                  <Loading />
                </el-icon>
                <el-icon v-else-if="item.status === 'completed'">
                  <CircleCheck />
                </el-icon>
                <el-icon v-else-if="item.status === 'error'">
                  <CircleClose />
                </el-icon>
              </div>
              
              <div class="item-details">
                <h4 class="item-name">{{ item.file.name }}</h4>
                <p class="item-meta">
                  {{ formatFileSize(item.file.size) }} • 
                  {{ formatTime(item.createdAt) }}
                </p>
              </div>
            </div>

            <div class="item-status">
              <el-tag
                :type="getStatusType(item.status)"
                size="small"
              >
                {{ getStatusText(item.status) }}
              </el-tag>
              
              <div v-if="item.status === 'uploading'" class="item-progress">
                <el-progress
                  :percentage="item.progress"
                  :stroke-width="4"
                  :show-text="false"
                />
              </div>
            </div>

            <div class="item-actions">
              <el-button
                v-if="item.status === 'completed' && item.imageId"
                link
                size="small"
                @click="viewImage(item.imageId)"
              >
                查看
              </el-button>
              
              <el-button
                v-if="item.status === 'error'"
                link
                size="small"
                @click="retryUpload(item)"
              >
                重试
              </el-button>
              
              <el-button
                link
                size="small"
                @click="removeFromQueue(item.id)"
              >
                移除
              </el-button>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- Recent Uploads -->
    <div v-if="recentUploads.length > 0" class="recent-uploads">
      <h2 class="section-title">最近上传</h2>
      <el-card>
        <div class="recent-grid">
          <div
            v-for="upload in recentUploads"
            :key="upload.id"
            class="recent-item"
            @click="viewImage(upload.id)"
          >
            <div class="recent-image">
              <img
                v-if="upload.thumbnailUrl"
                :src="upload.thumbnailUrl"
                :alt="upload.fileName"
                @error="handleImageError"
              />
              <div v-else class="recent-placeholder">
                <el-icon><Picture /></el-icon>
              </div>
            </div>
            
            <div class="recent-info">
              <h4 class="recent-name">{{ upload.fileName }}</h4>
              <p class="recent-time">{{ formatTime(upload.createdAt) }}</p>
            </div>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElNotification } from 'element-plus'
import {
  UploadFilled,
  Loading,
  Clock,
  CircleCheck,
  CircleClose,
  Picture,
  VideoCamera,
} from '@element-plus/icons-vue'
import { uploadImage, getImagesList, getThumbnailUrl } from '@/api/images'
import webSocketService from '@/api/websocket'
import dayjs from 'dayjs'
import CameraCapture from '@/components/CameraCapture.vue'

// Router
const router = useRouter()

// Refs
const fileInputRef = ref(null)

// State
const activeTab = ref('file')
const isDragOver = ref(false)
const isUploading = ref(false)
const uploadProgress = ref(0)
const uploadedCount = ref(0)
const totalFiles = ref(0)
const uploadQueue = ref([])
const recentUploads = ref([])

// Methods
const triggerFileSelect = () => {
  if (!isUploading.value) {
    fileInputRef.value?.click()
  }
}

const handleFileSelect = (event) => {
  const files = Array.from(event.target.files)
  addFilesToQueue(files)
  event.target.value = '' // Clear input
}

const handleDrop = (event) => {
  event.preventDefault()
  isDragOver.value = false
  
  const files = Array.from(event.dataTransfer.files).filter(file => 
    file.type.startsWith('image/')
  )
  
  if (files.length === 0) {
    ElMessage.warning('请选择图像文件')
    return
  }
  
  addFilesToQueue(files)
}

const handleDragOver = (event) => {
  event.preventDefault()
}

const handleDragEnter = (event) => {
  event.preventDefault()
  isDragOver.value = true
}

const handleDragLeave = (event) => {
  event.preventDefault()
  if (!event.currentTarget.contains(event.relatedTarget)) {
    isDragOver.value = false
  }
}

const addFilesToQueue = (files) => {
  const newItems = files.map(file => ({
    id: Date.now() + Math.random(),
    file,
    status: 'pending',
    progress: 0,
    imageId: null,
    createdAt: new Date(),
  }))
  
  uploadQueue.value.push(...newItems)
  startUpload()
}

const startUpload = async () => {
  if (isUploading.value) return
  
  const pendingItems = uploadQueue.value.filter(item => item.status === 'pending')
  if (pendingItems.length === 0) return
  
  isUploading.value = true
  totalFiles.value = pendingItems.length
  uploadedCount.value = 0
  uploadProgress.value = 0
  
  for (const item of pendingItems) {
    try {
      item.status = 'uploading'
      
      const response = await uploadImage(item.file, (progress) => {
        item.progress = progress
        updateOverallProgress()
      })
      
      item.status = 'completed'
      item.imageId = response.data?.imageId
      uploadedCount.value++
      
      ElNotification.success({
        title: '上传成功',
        message: `${item.file.name} 已上传成功，正在处理中...`,
      })
      
      // Emit event for other components to refresh
      window.dispatchEvent(new CustomEvent('image-uploaded', {
        detail: { imageId: response.data?.imageId, fileName: item.file.name }
      }))
      
    } catch (error) {
      item.status = 'error'
      item.error = error.message
      
      ElNotification.error({
        title: '上传失败',
        message: `${item.file.name} 上传失败: ${error.message}`,
      })
    }
    
    updateOverallProgress()
  }
  
  isUploading.value = false
  loadRecentUploads()
}

const updateOverallProgress = () => {
  const completedCount = uploadQueue.value.filter(
    item => item.status === 'completed' || item.status === 'error'
  ).length
  
  uploadProgress.value = Math.round((completedCount / totalFiles.value) * 100)
}

const retryUpload = async (item) => {
  item.status = 'pending'
  item.error = null
  startUpload()
}

const removeFromQueue = (itemId) => {
  const index = uploadQueue.value.findIndex(item => item.id === itemId)
  if (index > -1) {
    uploadQueue.value.splice(index, 1)
  }
}

const viewImage = (imageId) => {
  router.push(`/images/detail/${imageId}`)
}

const loadRecentUploads = async () => {
  try {
    const response = await getImagesList({
      page: 1,
      page_size: 12,
      sort: 'created_at',
      order: 'desc',
    })
    
    if (response && Array.isArray(response.data)) {
      // ✨ FIX: Map API data to frontend format with proper URL handling
      recentUploads.value = response.data
        .map(item => ({
          id: item.id,
          fileName: item.original_filename,
          fileSize: item.file_size,
          status: item.status,
          createdAt: item.created_at,
          // ✨ FIX: Use getThumbnailUrl helper to build proper thumbnail URL
          thumbnailUrl: item.thumbnail_url ? getThumbnailUrl(item.thumbnail_url) : null,
          originalUrl: item.original_url,
        }))
        // ✨ FIX: 确保按创建时间降序排序
        .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
    }
  } catch (error) {
    console.error('Failed to load recent uploads:', error)
  }
}

// Utility functions
const formatFileSize = (bytes) => {
  if (bytes === null || bytes === undefined) return '未知大小'
  if (bytes === 0) return '计算中...'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatTime = (time) => {
  return dayjs(time).format('MM-DD HH:mm')
}

const getStatusType = (status) => {
  const types = {
    pending: 'info',
    uploading: 'warning',
    completed: 'success',
    error: 'danger',
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    pending: '等待上传',
    uploading: '上传中',
    completed: '上传完成',
    error: '上传失败',
  }
  return texts[status] || '未知'
}

const handleImageError = (event) => {
  event.target.style.display = 'none'
}

// 摄像头上传处理
const handleCameraUpload = (data) => {
  console.log('摄像头照片上传成功:', data)
  
  // 添加到上传队列显示
  const uploadItem = {
    id: Date.now(),
    file: data.file,
    status: 'completed',
    progress: 100,
    imageId: data.imageId,
    createdAt: new Date()
  }
  
  uploadQueue.value.unshift(uploadItem)
  
  // 触发其他页面更新
  window.dispatchEvent(new CustomEvent('image-uploaded', {
    detail: { imageId: data.imageId, fileName: data.file.name }
  }))
  
  // 刷新最近上传
  loadRecentUploads()
}

const handleCameraError = (error) => {
  console.error('摄像头错误:', error)
}

// WebSocket event handlers
const handleImageProcessing = (data) => {
  // Find the corresponding upload item
  const item = uploadQueue.value.find(item => item.imageId === data.image_id)
  if (item) {
    ElNotification.info({
      title: '开始处理',
      message: `${data.file_name} 开始处理`,
    })
  }
}

const handleImageCompleted = (data) => {
  // Refresh recent uploads
  loadRecentUploads()
  
  ElNotification.success({
    title: '处理完成',
    message: `${data.file_name} 处理完成！`,
  })
}

const handleImageFailed = (data) => {
  ElNotification.error({
    title: '处理失败',
    message: `${data.file_name} 处理失败`,
  })
}

// Lifecycle
onMounted(async () => {
  await loadRecentUploads()
  
  // Setup WebSocket listeners
  webSocketService.on('image_processing', handleImageProcessing)
  webSocketService.on('image_completed', handleImageCompleted)
  webSocketService.on('image_failed', handleImageFailed)
  
  // Connect WebSocket if not connected
  if (!webSocketService.isConnected()) {
    try {
      await webSocketService.connect()
    } catch (error) {
      console.error('WebSocket connection failed:', error)
    }
  }
})

onUnmounted(() => {
  // Remove WebSocket listeners
  webSocketService.off('image_processing', handleImageProcessing)
  webSocketService.off('image_completed', handleImageCompleted)
  webSocketService.off('image_failed', handleImageFailed)
})
</script>

<style lang="scss" scoped>
.upload-page {
  padding: 0;
}

.page-header {
  margin-bottom: 24px;
  
  .page-title {
    font-size: 28px;
    font-weight: 600;
    color: var(--el-text-color-primary);
    margin: 0 0 8px 0;
  }
  
  .page-subtitle {
    font-size: 16px;
    color: var(--el-text-color-regular);
    margin: 0;
  }
}

.upload-area {
  margin-bottom: 32px;
}

.upload-card {
  :deep(.el-card__body) {
    padding: 0;
  }
}

.upload-tabs {
  :deep(.el-tabs__content) {
    padding: 0;
  }
  
  .tab-label {
    display: flex;
    align-items: center;
    gap: 6px;
  }
}

.upload-zone {
  min-height: 200px;
  border: 2px dashed var(--el-border-color);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  
  &:hover, &--dragover {
    border-color: var(--el-color-primary);
    background-color: var(--el-color-primary-light-9);
  }
  
  &--uploading {
    cursor: not-allowed;
    
    &:hover {
      border-color: var(--el-border-color);
      background-color: transparent;
    }
  }
}

.upload-content {
  text-align: center;
  padding: 32px;
  
  .upload-icon {
    font-size: 48px;
    color: var(--el-color-primary);
    margin-bottom: 16px;
  }
  
  .upload-title {
    font-size: 18px;
    font-weight: 500;
    color: var(--el-text-color-primary);
    margin: 0 0 8px 0;
  }
  
  .upload-hint {
    font-size: 14px;
    color: var(--el-text-color-regular);
    margin: 0 0 16px 0;
  }
}

.upload-progress {
  width: 100%;
  max-width: 300px;
  margin: 0 auto;
  
  .progress-text {
    font-size: 14px;
    color: var(--el-text-color-regular);
    margin-top: 8px;
  }
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 16px 0;
}

.upload-queue {
  margin-bottom: 32px;
}

.queue-list {
  .queue-item {
    display: flex;
    align-items: center;
    padding: 16px 0;
    border-bottom: 1px solid var(--el-border-color-lighter);
    
    &:last-child {
      border-bottom: none;
    }
  }
  
  .item-info {
    display: flex;
    align-items: center;
    flex: 1;
    
    .item-icon {
      margin-right: 12px;
      font-size: 18px;
      
      .el-icon {
        color: var(--el-color-primary);
      }
    }
    
    .item-details {
      .item-name {
        font-size: 14px;
        font-weight: 500;
        color: var(--el-text-color-primary);
        margin: 0 0 4px 0;
      }
      
      .item-meta {
        font-size: 12px;
        color: var(--el-text-color-regular);
        margin: 0;
      }
    }
  }
  
  .item-status {
    margin-right: 16px;
    
    .item-progress {
      width: 100px;
      margin-top: 8px;
    }
  }
  
  .item-actions {
    .el-button + .el-button {
      margin-left: 8px;
    }
  }
}

.recent-uploads {
  .recent-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 16px;
  }
  
  .recent-item {
    cursor: pointer;
    border-radius: 8px;
    overflow: hidden;
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    }
    
    .recent-image {
      aspect-ratio: 16/9;
      background-color: var(--el-fill-color-light);
      display: flex;
      align-items: center;
      justify-content: center;
      
      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
      
      .recent-placeholder {
        color: var(--el-text-color-placeholder);
        font-size: 32px;
      }
    }
    
    .recent-info {
      padding: 12px;
      
      .recent-name {
        font-size: 14px;
        font-weight: 500;
        color: var(--el-text-color-primary);
        margin: 0 0 4px 0;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
      
      .recent-time {
        font-size: 12px;
        color: var(--el-text-color-regular);
        margin: 0;
      }
    }
  }
}

// Mobile responsive
@media (max-width: 768px) {
  .upload-content {
    padding: 24px 16px;
    
    .upload-icon {
      font-size: 36px;
    }
    
    .upload-title {
      font-size: 16px;
    }
  }
  
  .recent-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 12px;
  }
  
  .queue-item {
    flex-direction: column;
    align-items: flex-start;
    
    .item-status {
      margin: 8px 0;
    }
    
    .item-actions {
      margin-top: 8px;
    }
  }
}

// Dark theme adjustments
[data-theme='dark'] {
  .upload-page {
    background-color: var(--el-bg-color-page);
    color: var(--el-text-color-primary);
  }
  
  .page-header {
    .page-title {
      color: var(--el-text-color-primary);
    }
    
    .page-subtitle {
      color: var(--el-text-color-regular);
    }
  }
  
  .upload-card {
    :deep(.el-card) {
      background-color: var(--el-bg-color);
      border-color: var(--el-border-color);
    }
    
    :deep(.el-tabs__header) {
      background-color: var(--el-bg-color);
    }
    
    :deep(.el-tabs__nav) {
      background-color: var(--el-bg-color);
    }
    
    :deep(.el-tabs__item) {
      color: var(--el-text-color-regular);
      
      &.is-active {
        color: var(--el-color-primary);
      }
    }
    
    :deep(.el-tabs__active-bar) {
      background-color: var(--el-color-primary);
    }
  }
  
  .upload-zone {
    border-color: var(--el-border-color);
    background-color: var(--el-fill-color-extra-light);
    
    &:hover, &--dragover {
      border-color: var(--el-color-primary);
      background-color: var(--el-color-primary-light-9);
    }
  }
  
  .upload-content {
    .upload-icon {
      color: var(--el-color-primary);
    }
    
    .upload-title {
      color: var(--el-text-color-primary);
    }
    
    .upload-hint {
      color: var(--el-text-color-regular);
    }
    
    .progress-text {
      color: var(--el-text-color-regular);
    }
  }
  
  .section-title {
    color: var(--el-text-color-primary);
  }
  
  .upload-queue {
    :deep(.el-card) {
      background-color: var(--el-bg-color);
      border-color: var(--el-border-color);
    }
  }
  
  .queue-item {
    border-bottom-color: var(--el-border-color-lighter);
    
    .item-details {
      .item-name {
        color: var(--el-text-color-primary);
      }
      
      .item-meta {
        color: var(--el-text-color-regular);
      }
    }
    
    .item-icon {
      :deep(.el-icon) {
        color: var(--el-color-primary);
      }
    }
  }
  
  .recent-uploads {
    :deep(.el-card) {
      background-color: var(--el-bg-color);
      border-color: var(--el-border-color);
    }
  }
  
  .recent-item {
    &:hover {
      background-color: var(--el-fill-color-light);
    }
    
    .recent-image {
      background-color: var(--el-fill-color);
      
      .recent-placeholder {
        color: var(--el-text-color-placeholder);
      }
    }
    
    .recent-info {
      background-color: var(--el-bg-color);
      
      .recent-name {
        color: var(--el-text-color-primary);
      }
      
      .recent-time {
        color: var(--el-text-color-regular);
      }
    }
  }
}
</style> 