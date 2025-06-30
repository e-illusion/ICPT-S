<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1 class="page-title">仪表盘</h1>
      <p class="page-subtitle">欢迎使用 ICPT 图像处理系统</p>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">
          <el-icon><Picture /></el-icon>
        </div>
        <div class="stat-content">
          <h3 class="stat-value">
            <span v-if="loading">-</span>
            <span v-else>{{ stats.totalImages || 0 }}</span>
          </h3>
          <p class="stat-label">总图像数</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <el-icon><Clock /></el-icon>
        </div>
        <div class="stat-content">
          <h3 class="stat-value">
            <span v-if="loading">-</span>
            <span v-else>{{ stats.todayProcessed || 0 }}</span>
          </h3>
          <p class="stat-label">今日上传</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <el-icon><TrendCharts /></el-icon>
        </div>
        <div class="stat-content">
          <h3 class="stat-value">
            <span v-if="loading">-</span>
            <span v-else>{{ stats.successRate || 0 }}%</span>
          </h3>
          <p class="stat-label">成功率</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <el-icon><Timer /></el-icon>
        </div>
        <div class="stat-content">
          <h3 class="stat-value">
            <span v-if="loading">-</span>
            <span v-else>{{ formatAvgTime(stats.avgTime) }}</span>
          </h3>
          <p class="stat-label">平均处理时间</p>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="quick-actions">
      <h2 class="section-title">快捷操作</h2>
      <div class="actions-grid">
        <el-card class="action-card" @click="goToUpload">
          <div class="action-content">
            <el-icon class="action-icon"><Upload /></el-icon>
            <h3>上传图像</h3>
            <p>上传新的图像进行处理</p>
          </div>
        </el-card>

        <el-card class="action-card" @click="goToGallery">
          <div class="action-content">
            <el-icon class="action-icon"><FolderOpened /></el-icon>
            <h3>图像列表</h3>
            <p>查看所有已处理的图像</p>
          </div>
        </el-card>

        <el-card class="action-card" @click="goToAnalytics">
          <div class="action-content">
            <el-icon class="action-icon"><DataAnalysis /></el-icon>
            <h3>数据分析</h3>
            <p>查看处理统计和分析</p>
          </div>
        </el-card>

        <el-card class="action-card" @click="goToSettings">
          <div class="action-content">
            <el-icon class="action-icon"><Setting /></el-icon>
            <h3>系统设置</h3>
            <p>配置系统参数</p>
          </div>
        </el-card>
      </div>
    </div>

    <!-- Recent Activity -->
    <div class="recent-activity">
      <h2 class="section-title">最近活动</h2>
      <el-card>
        <el-table :data="recentActivity" style="width: 100%">
          <el-table-column prop="time" label="时间" width="180">
            <template #default="{ row }">
              {{ formatTime(row.time) }}
            </template>
          </el-table-column>
          <el-table-column prop="action" label="操作" width="120" />
          <el-table-column prop="image" label="图像" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag
                :type="row.status === 'success' ? 'success' : 'danger'"
                size="small"
              >
                {{ row.status === 'success' ? '成功' : '失败' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Picture,
  Clock,
  TrendCharts,
  Timer,
  Upload,
  FolderOpened,
  DataAnalysis,
  Setting,
} from '@element-plus/icons-vue'
import { getDashboardStats, getRecentActivity } from '@/api/stats'
import dayjs from 'dayjs'

// Router
const router = useRouter()

// Reactive data
const stats = ref({
  totalImages: 0,
  todayProcessed: 0,
  successRate: 0,
  avgTime: 0,
})

const recentActivity = ref([])
const loading = ref(false)

// Methods
const goToUpload = () => {
  router.push('/images/upload')
}

const goToGallery = () => {
  router.push('/images/gallery')
}

const goToAnalytics = () => {
  router.push('/analytics')
}

const goToSettings = () => {
  router.push('/settings')
}

const formatTime = (time) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

const formatAvgTime = (time) => {
  if (!time || time === 0) {
    return '0ms'
  }
  
  // time 现在是毫秒单位
  if (time < 1) {
    // 小于1毫秒，显示微秒
    return `${(time * 1000).toFixed(0)}μs`
  } else if (time < 1000) {
    // 小于1秒，显示毫秒
    return `${time.toFixed(1)}ms`
  } else {
    // 大于1秒，显示秒
    return `${(time / 1000).toFixed(2)}s`
  }
}

const loadStats = async () => {
  try {
    loading.value = true
    const response = await getDashboardStats()
    
    if (response && response.data) {
    stats.value = {
        totalImages: response.data.total_images || 0,
        todayProcessed: response.data.today_processed || 0,
        successRate: response.data.success_rate ? parseFloat(response.data.success_rate).toFixed(1) : '0.0',
        avgTime: response.data.avg_time ? parseFloat(response.data.avg_time) : 0, // 保留原始毫秒数值
      }
    }
  } catch (error) {
    console.error('Failed to load stats:', error)
    ElMessage.error('加载统计信息失败')
    
    // 如果API失败，显示默认值
    stats.value = {
      totalImages: 0,
      todayProcessed: 0,
      successRate: '0.0',
      avgTime: 0,
    }
  } finally {
    loading.value = false
  }
}

const loadRecentActivity = async () => {
  try {
    const response = await getRecentActivity()
    
    if (response && response.data && Array.isArray(response.data)) {
      recentActivity.value = response.data.map(item => ({
        time: new Date(item.time),
        action: item.action,
        image: item.image,
        status: item.status,
      }))
    }
  } catch (error) {
    console.error('Failed to load recent activity:', error)
    ElMessage.error('加载最近活动失败')
    
    // 如果API失败，显示空数组
    recentActivity.value = []
  }
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    loadStats(),
    loadRecentActivity(),
  ])
})
</script>

<style lang="scss" scoped>
.dashboard {
  padding: 0;
}

.dashboard-header {
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

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
  margin-bottom: 32px;
}

.stat-card {
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-light);
  border-radius: 12px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.3s ease;

  &:hover {
    border-color: var(--el-color-primary);
    box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
    transform: translateY(-2px);
  }

  .stat-icon {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    background: linear-gradient(135deg, var(--el-color-primary), var(--el-color-primary-light-3));
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 24px;
  }

  .stat-content {
    flex: 1;

    .stat-value {
      font-size: 24px;
      font-weight: 600;
      color: var(--el-text-color-primary);
      margin: 0 0 4px 0;
    }

    .stat-label {
      font-size: 14px;
      color: var(--el-text-color-regular);
      margin: 0;
    }
  }
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 16px 0;
}

.quick-actions {
  margin-bottom: 32px;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.action-card {
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  }

  :deep(.el-card__body) {
    padding: 24px;
  }

  .action-content {
    text-align: center;

    .action-icon {
      font-size: 32px;
      color: var(--el-color-primary);
      margin-bottom: 12px;
    }

    h3 {
      font-size: 16px;
      font-weight: 600;
      color: var(--el-text-color-primary);
      margin: 0 0 8px 0;
    }

    p {
      font-size: 14px;
      color: var(--el-text-color-regular);
      margin: 0;
    }
  }
}

.recent-activity {
  :deep(.el-table) {
    background-color: transparent;
  }

  :deep(.el-table th) {
    background-color: var(--el-fill-color-light);
    color: var(--el-text-color-primary);
  }

  :deep(.el-table td) {
    color: var(--el-text-color-primary);
    border-bottom: 1px solid var(--el-border-color-lighter);
  }
  
  :deep(.el-table tr) {
    background-color: transparent;
  }
  
  :deep(.el-card__body) {
    background-color: var(--el-bg-color);
  }
}

// Dark theme adjustments
[data-theme='dark'] {
  .stat-card {
    background: var(--el-bg-color-page);
    border-color: var(--el-border-color);
    
    .stat-icon {
      background: linear-gradient(135deg, var(--el-color-primary), var(--el-color-primary-light-3));
    }
    
    .stat-content {
      .stat-value {
        color: var(--el-text-color-primary);
      }
      
      .stat-label {
        color: var(--el-text-color-regular);
      }
    }
    
    &:hover {
      background: var(--el-bg-color);
      border-color: var(--el-color-primary);
    }
  }
  
  .action-card {
    :deep(.el-card__body) {
      background: var(--el-bg-color-page);
    }
    
    &:hover {
      :deep(.el-card__body) {
        background: var(--el-bg-color);
      }
    }
  }
  
  .recent-activity {
    :deep(.el-card__body) {
      background: var(--el-bg-color-page);
    }
    
    :deep(.el-table) {
      background-color: var(--el-bg-color-page);
      
      .el-table__body-wrapper {
        background-color: var(--el-bg-color-page);
      }
    }
    
    :deep(.el-table th) {
      background-color: var(--el-fill-color-dark);
      color: var(--el-text-color-primary);
    }
    
    :deep(.el-table td) {
      background-color: var(--el-bg-color-page);
      color: var(--el-text-color-primary);
      border-bottom-color: var(--el-border-color);
    }
    
    :deep(.el-table tr) {
      background-color: var(--el-bg-color-page);
      
      &:hover {
        background-color: var(--el-fill-color-light) !important;
      }
    }
    
    :deep(.el-table__empty-block) {
      background-color: var(--el-bg-color-page);
    }
  }
}

// Mobile responsive
@media (max-width: 768px) {
  .dashboard-header {
    .page-title {
      font-size: 24px;
    }

    .page-subtitle {
      font-size: 14px;
    }
  }

  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  }

  .actions-grid {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }

  .stat-card {
    padding: 16px;

    .stat-icon {
      width: 40px;
      height: 40px;
      font-size: 20px;
    }

    .stat-content .stat-value {
      font-size: 20px;
    }
  }
}
</style> 