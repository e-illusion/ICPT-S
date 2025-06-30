<template>
  <div class="layout-container" :class="{ 'layout-collapsed': isCollapsed }">
    <!-- Sidebar -->
    <div class="sidebar" :class="{ 'sidebar-collapsed': isCollapsed }">
      <div class="sidebar-header">
        <div class="logo" v-if="!isCollapsed">
          <el-icon class="logo-icon">
            <Picture />
          </el-icon>
          <span class="logo-text">ICPT</span>
        </div>
        <el-icon v-else class="logo-icon-collapsed">
          <Picture />
        </el-icon>
      </div>

      <el-menu
        :default-active="activeMenu"
        class="sidebar-menu"
        :collapse="isCollapsed"
        :unique-opened="false"
        background-color="var(--sidebar-bg)"
        text-color="var(--sidebar-text)"
        active-text-color="var(--sidebar-active-text)"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <template #title>仪表盘</template>
        </el-menu-item>

        <el-sub-menu index="/images">
          <template #title>
            <el-icon><Picture /></el-icon>
            <span>图像管理</span>
          </template>
          <el-menu-item index="/images/upload">
            <el-icon><Upload /></el-icon>
            <template #title>图像上传</template>
          </el-menu-item>
          <el-menu-item index="/images/gallery">
            <el-icon><FolderOpened /></el-icon>
            <template #title>图像列表</template>
          </el-menu-item>
        </el-sub-menu>

        <el-menu-item index="/analytics">
          <el-icon><TrendCharts /></el-icon>
          <template #title>数据分析</template>
        </el-menu-item>

        <el-menu-item index="/profile">
          <el-icon><User /></el-icon>
          <template #title>个人中心</template>
        </el-menu-item>

        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <template #title>系统设置</template>
        </el-menu-item>
      </el-menu>
    </div>

    <!-- Main Content -->
    <div class="main-container">
      <!-- Header -->
      <div class="header">
        <div class="header-left">
          <el-button
            link
            class="collapse-btn"
            @click="toggleCollapse"
          >
            <el-icon>
              <Fold v-if="!isCollapsed" />
              <Expand v-else />
            </el-icon>
          </el-button>

          <el-breadcrumb class="breadcrumb" separator="/">
            <el-breadcrumb-item
              v-for="item in breadcrumbs"
              :key="item.path"
              :to="item.path === currentRoute.path ? '' : item.path"
            >
              {{ item.meta.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>

        <div class="header-right">
          <!-- Theme Toggle -->
          <el-tooltip content="切换主题" placement="bottom">
            <el-button
              link
              class="theme-btn"
              @click="toggleTheme"
            >
              <el-icon>
                <Sunny v-if="isDark" />
                <Moon v-else />
              </el-icon>
            </el-button>
          </el-tooltip>

          <!-- User Menu -->
          <el-dropdown trigger="click" @command="handleUserCommand">
            <div class="user-info">
              <el-avatar
                :src="userStore.user?.avatar"
                :size="32"
                class="user-avatar"
              >
                <el-icon><User /></el-icon>
              </el-avatar>
              <span v-if="!isCollapsed" class="username">
                {{ userStore.user?.username || 'Admin' }}
              </span>
              <el-icon v-if="!isCollapsed" class="dropdown-icon">
                <ArrowDown />
              </el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人中心
                </el-dropdown-item>
                <el-dropdown-item command="settings">
                  <el-icon><Setting /></el-icon>
                  系统设置
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>

      <!-- Content -->
      <div class="content">
        <router-view v-slot="{ Component, route }">
          <transition name="fade-slide" mode="out-in">
            <component :is="Component" :key="route.path" />
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import {
  Picture,
  DataAnalysis,
  Upload,
  FolderOpened,
  TrendCharts,
  User,
  Setting,
  Fold,
  Expand,
  Sunny,
  Moon,
  ArrowDown,
  SwitchButton,
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

// Router and stores
const route = useRoute()
const router = useRouter()
const userStore = useAuthStore()

// Reactive state
const isCollapsed = ref(false)
const isDark = ref(false)

// Computed properties
const currentRoute = computed(() => route)
const activeMenu = computed(() => route.path)

const breadcrumbs = computed(() => {
  const matched = route.matched.filter(item => item.meta && item.meta.title)
  return matched
})

// Methods
const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
  localStorage.setItem('sidebar-collapsed', isCollapsed.value.toString())
}

const toggleTheme = () => {
  isDark.value = !isDark.value
  const theme = isDark.value ? 'dark' : 'light'
  document.documentElement.setAttribute('data-theme', theme)
  localStorage.setItem('theme', theme)
  
  // Add/remove Element Plus dark class
  if (isDark.value) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

const handleUserCommand = async (command) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'settings':
      router.push('/settings')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm(
          '确定要退出登录吗？',
          '退出确认',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
          }
        )
        await userStore.logout()
        ElMessage.success('已退出登录')
        router.push('/login')
      } catch {
        // User cancelled
      }
      break
  }
}

// Initialize settings
onMounted(() => {
  // Initialize sidebar state
  const savedCollapsed = localStorage.getItem('sidebar-collapsed')
  if (savedCollapsed !== null) {
    isCollapsed.value = savedCollapsed === 'true'
  }

  // Initialize theme
  const savedTheme = localStorage.getItem('theme') || 'light'
  isDark.value = savedTheme === 'dark'
  document.documentElement.setAttribute('data-theme', savedTheme)
  
  // Apply Element Plus dark class
  if (savedTheme === 'dark') {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
})

// Watch for route changes to update active menu
watch(
  () => route.path,
  (newPath) => {
    // Additional logic if needed
  }
)
</script>

<style lang="scss" scoped>
.layout-container {
  display: flex;
  height: 100vh;
  background-color: var(--el-bg-color-page);
}

.sidebar {
  width: 240px;
  background-color: var(--sidebar-bg, #001529);
  transition: width 0.3s ease;
  overflow: hidden;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);

  &-collapsed {
    width: 64px;
  }

  &-header {
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 16px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);

    .logo {
      display: flex;
      align-items: center;
      gap: 8px;
      color: white;

      &-icon {
        font-size: 24px;
        color: #1890ff;
      }

      &-text {
        font-size: 18px;
        font-weight: 600;
      }
    }

    .logo-icon-collapsed {
      font-size: 24px;
      color: #1890ff;
    }
  }

  &-menu {
    border: none;
    height: calc(100vh - 64px);

    :deep(.el-menu-item) {
      height: 48px;
      line-height: 48px;
      margin: 4px 8px;
      border-radius: 6px;

      &:hover {
        background-color: rgba(255, 255, 255, 0.1);
      }

      &.is-active {
        background-color: var(--el-color-primary);
        color: white;
      }
    }

    :deep(.el-sub-menu__title) {
      height: 48px;
      line-height: 48px;
      margin: 4px 8px;
      border-radius: 6px;

      &:hover {
        background-color: rgba(255, 255, 255, 0.1);
      }
    }
  }
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.header {
  height: 64px;
  background-color: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 1001;

  &-left {
    display: flex;
    align-items: center;
    gap: 16px;

    .collapse-btn {
      padding: 8px;
      
      .el-icon {
        font-size: 18px;
      }
    }

    .breadcrumb {
      :deep(.el-breadcrumb__item:last-child .el-breadcrumb__inner) {
        color: var(--el-text-color-primary);
      }
    }
  }

  &-right {
    display: flex;
    align-items: center;
    gap: 16px;

    .theme-btn {
      padding: 8px;
      
      .el-icon {
        font-size: 18px;
      }
    }

    .user-info {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
      padding: 8px;
      border-radius: 6px;
      transition: background-color 0.3s;

      &:hover {
        background-color: var(--el-fill-color-light);
      }

      .username {
        font-size: 14px;
        color: var(--el-text-color-primary);
      }

      .dropdown-icon {
        font-size: 12px;
        color: var(--el-text-color-regular);
      }
    }
  }
}

.content {
  flex: 1;
  padding: 16px;
  overflow: auto;
  background-color: var(--el-bg-color-page);
}

// Page transition animations
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

// Dark theme variables
:root {
  --sidebar-bg: #001529;
  --sidebar-text: rgba(255, 255, 255, 0.85);
  --sidebar-active-text: #ffffff;
}

[data-theme='dark'] {
  --sidebar-bg: #141414;
  
  .layout-container {
    background-color: var(--el-bg-color-page);
  }
  
  .sidebar {
    background-color: var(--sidebar-bg);
    box-shadow: 2px 0 8px rgba(0, 0, 0, 0.3);
    
    &-header {
      border-bottom-color: rgba(255, 255, 255, 0.1);
      
      .logo {
        color: var(--el-text-color-primary);
        
        &-text {
          color: var(--el-text-color-primary);
        }
      }
    }
    
    &-menu {
      :deep(.el-menu-item) {
        color: var(--sidebar-text);
        
        &:hover {
          background-color: rgba(255, 255, 255, 0.1);
          color: var(--el-text-color-primary);
        }
        
        &.is-active {
          background-color: var(--el-color-primary);
          color: white;
        }
      }
      
      :deep(.el-sub-menu__title) {
        color: var(--sidebar-text);
        
        &:hover {
          background-color: rgba(255, 255, 255, 0.1);
          color: var(--el-text-color-primary);
        }
      }
      
      :deep(.el-sub-menu.is-opened .el-sub-menu__title) {
        color: var(--el-text-color-primary);
      }
    }
  }
  
  .header {
    background-color: var(--el-bg-color);
    border-bottom-color: var(--el-border-color);
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
    
    &-left {
      .breadcrumb {
        :deep(.el-breadcrumb__item) {
          .el-breadcrumb__inner {
            color: var(--el-text-color-regular);
            
            &:hover {
              color: var(--el-color-primary);
            }
          }
          
          &:last-child .el-breadcrumb__inner {
            color: var(--el-text-color-primary);
          }
        }
        
        :deep(.el-breadcrumb__separator) {
          color: var(--el-text-color-placeholder);
        }
      }
    }
    
    &-right {
      .user-info {
        &:hover {
          background-color: var(--el-fill-color);
        }
        
        .username {
          color: var(--el-text-color-primary);
        }
        
        .dropdown-icon {
          color: var(--el-text-color-regular);
        }
      }
    }
  }
  
  .content {
    background-color: var(--el-bg-color-page);
  }
  
  // Fix dropdown menu in dark mode
  :deep(.el-dropdown-menu) {
    background-color: var(--el-bg-color-overlay);
    border-color: var(--el-border-color);
    
    .el-dropdown-menu__item {
      color: var(--el-text-color-primary);
      
      &:hover {
        background-color: var(--el-fill-color);
        color: var(--el-color-primary);
      }
    }
  }
}

// Responsive design
@media (max-width: 768px) {
  .layout-container {
    &.layout-collapsed {
      .sidebar {
        transform: translateX(-100%);
      }
    }
  }

  .sidebar {
    position: absolute;
    z-index: 1000;
    height: 100vh;
  }

  .main-container {
    margin-left: 0;
  }

  .header-left .breadcrumb {
    display: none;
  }
}
</style> 