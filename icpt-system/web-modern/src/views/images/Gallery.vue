<template>
  <div class="gallery-page">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">图像列表</h1>
        <p class="page-subtitle">管理您的所有图像，支持筛选、搜索和批量操作</p>
      </div>
      
      <div class="header-actions">
        <el-button
          type="primary"
          :icon="Upload"
          @click="goToUpload"
        >
          上传图像
        </el-button>
      </div>
    </div>

    <div class="filters-section">
      <el-card>
        <div class="filters-content">
          <div class="filter-group">
            <el-input
              v-model="searchQuery"
              placeholder="搜索图像名称..."
              :prefix-icon="Search"
              clearable
              @input="handleSearch"
              @clear="handleSearch"
              style="width: 300px"
            />
          </div>
          
          <div class="filter-group">
            <el-select
              v-model="statusFilter"
              placeholder="筛选状态"
              clearable
              @change="handleFilterChange"
              @clear="handleFilterChange"
              style="width: 150px"
            >
              <el-option label="全部状态" value="" />
              <el-option label="处理中" value="processing" />
              <el-option label="已完成" value="completed" />
              <el-option label="处理失败" value="failed" />
            </el-select>
          </div>
          
          <div class="filter-group">
            <el-select
              v-model="sortField"
              @change="handleSortChange"
              style="width: 150px"
            >
              <el-option label="创建时间" value="created_at" />
              <el-option label="文件名" value="original_filename" />
              <el-option label="文件大小" value="file_size" />
            </el-select>
          </div>
          
          <div class="filter-group">
            <el-select
              v-model="sortOrder"
              @change="handleSortChange"
              style="width: 100px"
            >
              <el-option label="降序" value="desc" />
              <el-option label="升序" value="asc" />
            </el-select>
          </div>
          
          <div class="filter-group">
            <el-button
              v-if="selectedImages.length > 0"
              type="danger"
              :icon="Delete"
              @click="batchDelete"
            >
              删除选中 ({{ selectedImages.length }})
            </el-button>
          </div>
        </div>
      </el-card>
    </div>

    <div class="images-section">
      <div v-if="loading" class="loading-state">
        <el-skeleton :rows="3" animated />
      </div>
      
      <div v-else-if="images.length === 0" class="empty-state">
        <el-empty description="暂无图像数据">
          <el-button type="primary" @click="goToUpload">
            立即上传
          </el-button>
        </el-empty>
      </div>
      
      <div v-else class="images-grid">
        <div
          v-for="image in images"
          :key="image.id"
          class="image-item"
          :class="{ 'image-item--selected': selectedImages.includes(image.id) }"
        >
          <div class="image-selection">
            <el-checkbox
              :model-value="selectedImages.includes(image.id)"
              @change="toggleImageSelection(image.id)"
            />
          </div>
          
          <div class="image-display" @click="viewImageDetail(image.id)">
            <img
              v-if="image.thumbnailUrl"
              :src="getThumbnailUrl(image.thumbnailUrl)"
              :alt="image.fileName"
              @error="handleImageError"
            />
            <div v-else class="image-placeholder">
              <el-icon><Picture /></el-icon>
            </div>
            
            <div v-if="image.status !== 'completed'" class="status-overlay">
              <el-icon v-if="image.status === 'processing'">
                <Loading />
              </el-icon>
              <el-icon v-else-if="image.status === 'failed'">
                <Warning />
              </el-icon>
            </div>
          </div>
          
          <div class="image-info">
            <h4 class="image-name" :title="image.fileName">
              {{ image.fileName }}
            </h4>
            
            <div class="image-meta">
              <span class="meta-item">
                {{ formatFileSize(image.fileSize) }}
              </span>
              <span class="meta-item">
                {{ formatTime(image.createdAt) }}
              </span>
            </div>
            
            <div class="image-status">
              <el-tag
                :type="getStatusType(image.status)"
                size="small"
              >
                {{ getStatusText(image.status) }}
              </el-tag>
            </div>
          </div>
          
          <div class="image-actions">
            <el-button-group>
              <el-button
                type="primary"
                size="small"
                @click="viewImageDetail(image.id)"
              >
                查看
              </el-button>
              
              <el-button
                v-if="image.status === 'completed' && image.originalUrl"
                size="small"
                @click="downloadImage(image)"
              >
                下载
              </el-button>
              
              <el-button
                v-if="image.status === 'failed'"
                type="warning"
                size="small"
                @click="reprocessImage(image.id)"
              >
                重试
              </el-button>
              
              <el-button
                type="danger"
                size="small"
                @click="deleteImage(image.id)"
              >
                删除
              </el-button>
            </el-button-group>
          </div>
        </div>
      </div>
    </div>

    <div v-if="totalCount > 0" class="pagination-section">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="totalCount"
        :page-sizes="[12, 24, 48, 96]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handlePageSizeChange"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'
import {
  Upload,
  Search,
  Delete,
  Picture,
  Loading,
  Warning,
} from '@element-plus/icons-vue'
import {
  getImagesList,
  deleteImage as deleteImageApi,
  batchDeleteImages,
  reprocessImage as reprocessImageApi,
  getThumbnailUrl as getImageThumbnailUrl,
  getOriginalUrl,
} from '@/api/images'
import webSocketService from '@/api/websocket'
import dayjs from 'dayjs'

// Router
const router = useRouter()

// State
const loading = ref(false)
const images = ref([])
const selectedImages = ref([])
const searchQuery = ref('')
const statusFilter = ref('')
const sortField = ref('created_at')
const sortOrder = ref('desc')
const currentPage = ref(1)
const pageSize = ref(24)
const totalCount = ref(0)

// Computed
const hasSelection = computed(() => selectedImages.value.length > 0)

// Methods
const loadImages = async () => {
  try {
    loading.value = true
    
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      sort: sortField.value,
      order: sortOrder.value,
    }
    
    if (searchQuery.value.trim()) {
      params.search = searchQuery.value.trim()
    }
    
    if (statusFilter.value) {
      params.status = statusFilter.value
    }
    
    const response = await getImagesList(params)
    
    // The response from getImagesList is already the "data" part of the API response,
    // which should be in the format: { data: [...], total: xxx, page: xxx, ... }
    if (response && Array.isArray(response.data)) {
      let mappedImages = response.data.map(item => ({
        id: item.id,
        fileName: item.original_filename,
        fileSize: item.file_size,
        status: item.status,
        createdAt: item.created_at,
        // The URL is now a relative path, and getThumbnailUrl will prepend /static/
        thumbnailUrl: item.thumbnail_url,
        originalUrl: item.original_url,
      }))
      
      // ✨ FIX: 添加前端排序逻辑，确保排序正确
      mappedImages = applySorting(mappedImages, sortField.value, sortOrder.value)
      
      images.value = mappedImages
      totalCount.value = response.total || 0
      
      // ✨ FIX: Dynamically fetch file sizes for images with zero file_size
      // Use try-catch to prevent this from breaking the main function
      try {
        await loadMissingFileSizes()
      } catch (error) {
        console.warn('Failed to load missing file sizes:', error)
        // This error shouldn't prevent the main loading from succeeding
      }
    } else {
      console.warn('Unexpected API response format:', response)
      images.value = []
      totalCount.value = 0
    }
  } catch (error) {
    console.error('Failed to load images:', error)
    ElMessage.error('加载图像列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = debounce(() => {
  currentPage.value = 1
  loadImages()
}, 500)

const handleFilterChange = () => {
  currentPage.value = 1
  loadImages()
}

const handleSortChange = () => {
  currentPage.value = 1
  loadImages()
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadImages()
}

const handlePageSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadImages()
}

const toggleImageSelection = (imageId) => {
  const index = selectedImages.value.indexOf(imageId)
  if (index > -1) {
    selectedImages.value.splice(index, 1)
  } else {
    selectedImages.value.push(imageId)
  }
}

const clearSelection = () => {
  selectedImages.value = []
}

const viewImageDetail = (imageId) => {
  router.push(`/images/detail/${imageId}`)
}

const goToUpload = () => {
  router.push('/images/upload')
}

const deleteImage = async (imageId) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这张图像吗？此操作不可逆。',
      '删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await deleteImageApi(imageId)
    
    ElMessage.success('图像删除成功')
    
    // Remove from local list
    const index = images.value.findIndex(img => img.id === imageId)
    if (index > -1) {
      images.value.splice(index, 1)
      totalCount.value--
    }
    
    // Remove from selection
    const selectionIndex = selectedImages.value.indexOf(imageId)
    if (selectionIndex > -1) {
      selectedImages.value.splice(selectionIndex, 1)
    }
    
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Delete image failed:', error)
      ElMessage.error('删除图像失败')
    }
  }
}

const batchDelete = async () => {
  if (selectedImages.value.length === 0) return
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedImages.value.length} 张图像吗？此操作不可逆。`,
      '批量删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await batchDeleteImages(selectedImages.value)
    
    ElMessage.success(`成功删除 ${selectedImages.value.length} 张图像`)
    
    // Reload images
    selectedImages.value = []
    await loadImages()
    
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Batch delete failed:', error)
      ElMessage.error('批量删除失败')
    }
  }
}

const reprocessImage = async (imageId) => {
  try {
    await reprocessImageApi(imageId)
    ElMessage.success('已重新提交处理')
    
    // Update local status
    const image = images.value.find(img => img.id === imageId)
    if (image) {
      image.status = 'processing'
    }
  } catch (error) {
    console.error('Reprocess image failed:', error)
    ElMessage.error('重新处理失败')
  }
}

const downloadImage = async (image) => {
  try {
    const url = getOriginalUrl(image.originalUrl)
    
    // 显示下载开始提示
    ElMessage.info('开始下载图像...')
    
    // 使用fetch下载文件，这样可以更好地处理HTTPS和错误
    const response = await fetch(url, {
      method: 'GET',
      credentials: 'same-origin',
      headers: {
        'Accept': 'image/*'
      }
    })
    
    if (!response.ok) {
      throw new Error(`下载失败: ${response.status} ${response.statusText}`)
    }
    
    // 获取文件blob
    const blob = await response.blob()
    
    // 创建下载链接
    const downloadUrl = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = image.fileName || 'image.jpg'
    link.style.display = 'none'
    
    // 添加到DOM并触发下载
    document.body.appendChild(link)
    link.click()
    
    // 清理
    document.body.removeChild(link)
    window.URL.revokeObjectURL(downloadUrl)
    
    ElMessage.success('图像下载成功！')
    
  } catch (error) {
    console.error('Download error:', error)
    ElMessage.error(`下载失败: ${error.message}`)
  }
}

const getThumbnailUrl = (thumbnailPath) => {
  return getImageThumbnailUrl(thumbnailPath)
}

// ✨ FIX: 前端排序函数，确保排序正确
const applySorting = (imageList, sortField, sortOrder) => {
  if (!imageList || imageList.length === 0) return imageList
  
  const sortedImages = [...imageList].sort((a, b) => {
    let valueA, valueB
    
    switch (sortField) {
      case 'created_at':
        valueA = new Date(a.createdAt).getTime()
        valueB = new Date(b.createdAt).getTime()
        break
      case 'original_filename':
        valueA = (a.fileName || '').toLowerCase()
        valueB = (b.fileName || '').toLowerCase()
        break
      case 'file_size':
        valueA = a.fileSize || 0
        valueB = b.fileSize || 0
        break
      default:
        valueA = new Date(a.createdAt).getTime()
        valueB = new Date(b.createdAt).getTime()
    }
    
    // 处理字符串和数字比较
    if (typeof valueA === 'string' && typeof valueB === 'string') {
      if (sortOrder === 'asc') {
        return valueA.localeCompare(valueB)
      } else {
        return valueB.localeCompare(valueA)
      }
    } else {
      if (sortOrder === 'asc') {
        return valueA - valueB
      } else {
        return valueB - valueA
      }
    }
  })
  
  return sortedImages
}

// Utility functions
const formatFileSize = (bytes) => {
  // ✨ FIX: Better handling for zero or missing file size
  if (bytes === null || bytes === undefined) return '未知大小'
  if (bytes === 0) return '计算中...'  // Show "Computing..." for zero values as we're loading them dynamically
  if (bytes === -1) return '获取失败' // Show "Failed to get" for failed loading
  if (bytes < 0) return '无效大小'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatTime = (time) => {
  return dayjs(time).format('MM-DD HH:mm')
}

const getStatusType = (status) => {
  const types = {
    processing: 'warning',
    completed: 'success',
    failed: 'danger',
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    processing: '处理中',
    completed: '已完成',
    failed: '处理失败',
  }
  return texts[status] || '未知'
}

const handleImageError = (event) => {
  event.target.style.display = 'none'
}

// Debounce function with error handling
function debounce(func, wait) {
  let timeout
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout)
      try {
        func(...args)
      } catch (error) {
        console.error('Debounced function error:', error)
      }
    }
    clearTimeout(timeout)
    timeout = setTimeout(later, wait)
  }
}

// WebSocket event handlers
const handleImageCompleted = (data) => {
  // Update local image status
  const image = images.value.find(img => img.id === data.image_id)
  if (image) {
    image.status = 'completed'
    image.thumbnailUrl = data.thumbnail_url
    
    // ✨ FIX: Update the file size from the WebSocket event data
    if (data.file_size !== undefined) {
      image.fileSize = data.file_size
    }
  }
  
  ElNotification.success({
    title: '处理完成',
    message: `${data.file_name} 处理完成！`,
  })
}

const handleImageFailed = (data) => {
  // Update local image status
  const image = images.value.find(img => img.id === data.image_id)
  if (image) {
    image.status = 'failed'
  }
  
  ElNotification.error({
    title: '处理失败',
    message: `${data.file_name} 处理失败`,
  })
}

const handleImageProcessing = (data) => {
  // Update local image status
  const image = images.value.find(img => img.id === data.image_id)
  if (image) {
    image.status = 'processing'
  }
}

// Handle image upload events from other pages
const handleImageUploaded = () => {
  // Refresh images when new upload occurs
  loadImages()
}

// Lifecycle
onMounted(async () => {
  await loadImages()
  
  // Setup WebSocket listeners
  webSocketService.on('image_completed', handleImageCompleted)
  webSocketService.on('image_failed', handleImageFailed)
  webSocketService.on('image_processing', handleImageProcessing)
  
  // Listen for upload events from other pages
  window.addEventListener('image-uploaded', handleImageUploaded)
  
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
  webSocketService.off('image_completed', handleImageCompleted)
  webSocketService.off('image_failed', handleImageFailed)
  webSocketService.off('image_processing', handleImageProcessing)
  
  // Remove upload event listener
  window.removeEventListener('image-uploaded', handleImageUploaded)
})

// Watch for route changes to refresh data
watch(() => router.currentRoute.value.path, (newPath) => {
  if (newPath === '/images/gallery') {
    loadImages()
  }
})

// ✨ NEW: Function to dynamically load missing file sizes
const loadMissingFileSizes = async () => {
  const imagesWithoutSize = images.value.filter(img => 
    (!img.fileSize || img.fileSize === 0) && img.originalUrl && img.status === 'completed'
  )
  
  if (imagesWithoutSize.length === 0) return
  
  console.log(`Loading file sizes for ${imagesWithoutSize.length} images...`)
  
  // Load file sizes in parallel, but limit concurrent requests to avoid overwhelming the server
  const concurrency = 3 // Reduced concurrency to be more conservative
  for (let i = 0; i < imagesWithoutSize.length; i += concurrency) {
    const batch = imagesWithoutSize.slice(i, i + concurrency)
    await Promise.allSettled(batch.map(async (image) => {
      try {
        const fileSize = await getFileSize(image.originalUrl)
        if (fileSize > 0) {
          image.fileSize = fileSize
        }
      } catch (error) {
        console.warn(`Failed to get file size for ${image.fileName}:`, error)
        // Set a default message for failed file size loading
        image.fileSize = -1 // Use -1 to indicate failed loading
      }
    }))
  }
}

// ✨ NEW: Function to get file size via HTTP HEAD request
const getFileSize = async (originalUrl) => {
  if (!originalUrl) return 0
  
  try {
    const url = getOriginalUrl(originalUrl)
    
    // Create abort controller for timeout
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), 5000) // 5 second timeout
    
    const response = await fetch(url, { 
      method: 'HEAD',
      mode: 'cors',
      credentials: 'same-origin',
      signal: controller.signal
    })
    
    clearTimeout(timeoutId) // Clear timeout if request completes
    
    if (response.ok) {
      const contentLength = response.headers.get('content-length')
      if (contentLength && contentLength !== '0') {
        return parseInt(contentLength, 10)
      }
    }
  } catch (error) {
    // Don't log timeout errors as warnings, they're expected
    if (error.name !== 'AbortError') {
      console.warn('Failed to fetch file size:', error)
    }
  }
  
  return 0
}
</script>

<style lang="scss" scoped>
.gallery-page {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  
  .header-content {
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
}

.filters-section {
  margin-bottom: 24px;
  
  .filters-content {
    display: flex;
    align-items: center;
    gap: 16px;
    flex-wrap: wrap;
  }
  
  .filter-group {
    display: flex;
    align-items: center;
  }
}

.images-section {
  margin-bottom: 24px;
}

.loading-state {
  padding: 40px;
}

.empty-state {
  padding: 80px 40px;
  text-align: center;
}

.images-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.image-item {
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-light);
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
  position: relative;
  
  &:hover {
    border-color: var(--el-color-primary);
    box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
    transform: translateY(-2px);
  }
  
  &--selected {
    border-color: var(--el-color-primary);
    background-color: var(--el-color-primary-light-9);
  }
}

.image-selection {
  position: absolute;
  top: 8px;
  left: 8px;
  z-index: 2;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 4px;
  padding: 4px;
}

.image-display {
  aspect-ratio: 16/9;
  background-color: var(--el-fill-color-light);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  cursor: pointer;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  .image-placeholder {
    color: var(--el-text-color-placeholder);
    font-size: 48px;
  }
  
  .status-overlay {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 8px;
    border-radius: 50%;
    font-size: 20px;
    
    .el-icon {
      animation: pulse 2s infinite;
    }
  }
}

.image-info {
  padding: 16px;
  
  .image-name {
    font-size: 14px;
    font-weight: 500;
    color: var(--el-text-color-primary);
    margin: 0 0 8px 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .image-meta {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
    
    .meta-item {
      font-size: 12px;
      color: var(--el-text-color-regular);
    }
  }
  
  .image-status {
    margin-bottom: 12px;
  }
}

.image-actions {
  padding: 0 16px 16px;
  
  :deep(.el-button-group) {
    width: 100%;
    
    .el-button {
      flex: 1;
      font-size: 12px;
    }
  }
}

.pagination-section {
  display: flex;
  justify-content: center;
  padding: 32px 0;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

// Mobile responsive
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .filters-content {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
    
    .filter-group {
      width: 100%;
      
      .el-input,
      .el-select {
        width: 100% !important;
      }
    }
  }
  
  .images-grid {
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 16px;
  }
  
  .pagination-section {
    :deep(.el-pagination) {
      flex-wrap: wrap;
      justify-content: center;
    }
  }
}

// Dark theme adjustments
[data-theme='dark'] {
  .gallery-page {
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
  
  .filters-section {
    :deep(.el-card) {
      background-color: var(--el-bg-color);
      border-color: var(--el-border-color);
    }
    
    :deep(.el-input) {
      .el-input__wrapper {
        background-color: var(--el-fill-color);
        border-color: var(--el-border-color);
      }
      
      .el-input__inner {
        color: var(--el-text-color-primary);
        
        &::placeholder {
          color: var(--el-text-color-placeholder);
        }
      }
    }
    
    :deep(.el-select) {
      .el-select__wrapper {
        background-color: var(--el-fill-color);
        border-color: var(--el-border-color);
      }
    }
    
    :deep(.el-button) {
      &.el-button--danger {
        background-color: var(--el-color-danger);
        border-color: var(--el-color-danger);
      }
    }
  }
  
  .image-item {
    background-color: var(--el-bg-color);
    border-color: var(--el-border-color);
    
    &:hover {
      background-color: var(--el-fill-color-light);
      border-color: var(--el-color-primary);
    }
    
    &--selected {
      background-color: var(--el-color-primary-light-9);
      border-color: var(--el-color-primary);
    }
  }
  
  .image-selection {
    background: rgba(0, 0, 0, 0.7);
  }
  
  .image-display {
    background-color: var(--el-fill-color);
    
    .image-placeholder {
      color: var(--el-text-color-placeholder);
    }
  }
  
  .image-info {
    .image-name {
      color: var(--el-text-color-primary);
    }
    
    .meta-item {
      color: var(--el-text-color-regular);
    }
  }
  
  .empty-state {
    color: var(--el-text-color-regular);
  }
  
  :deep(.el-pagination) {
    .el-pager li {
      background-color: var(--el-fill-color);
      color: var(--el-text-color-primary);
      
      &.is-active {
        background-color: var(--el-color-primary);
        color: white;
      }
    }
    
    .btn-prev,
    .btn-next {
      background-color: var(--el-fill-color);
      color: var(--el-text-color-primary);
    }
    
    .el-select .el-select__wrapper {
      background-color: var(--el-fill-color);
      border-color: var(--el-border-color);
    }
  }
}
</style>
