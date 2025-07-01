import { createRouter, createWebHistory } from 'vue-router'
import NProgress from 'nprogress'
import { useAuthStore } from '@/stores/auth'
import { getToken } from '@/utils/auth'

// Layout components
const Layout = () => import('@/layout/index.vue')

// Page components (lazy loaded)
const Login = () => import('@/views/auth/Login.vue')
const Register = () => import('@/views/auth/Register.vue')
const Dashboard = () => import('@/views/dashboard/index.vue')
const ImageUpload = () => import('@/views/images/Upload.vue')
const ImageGallery = () => import('@/views/images/Gallery.vue')
const ImageDetail = () => import('@/views/images/Detail.vue')
const Profile = () => import('@/views/profile/index.vue')
const Settings = () => import('@/views/settings/index.vue')
const Analytics = () => import('@/views/analytics/index.vue')

// Error pages
const NotFound = () => import('@/views/error/404.vue')
const Forbidden = () => import('@/views/error/403.vue')
const ServerError = () => import('@/views/error/500.vue')

const routes = [
    {
        path: '/',
        redirect: '/dashboard',
    },
    {
        path: '/login',
        name: 'Login',
        component: Login,
        meta: {
            title: '登录',
            requiresAuth: false,
            hideInMenu: true,
        },
    },
    {
        path: '/register',
        name: 'Register',
        component: Register,
        meta: {
            title: '注册',
            requiresAuth: false,
            hideInMenu: true,
        },
    },
    {
        path: '/dashboard',
        component: Layout,
        redirect: '/dashboard/index',
        children: [
            {
                path: 'index',
                name: 'Dashboard',
                component: Dashboard,
                meta: {
                    title: '仪表盘',
                    icon: 'DataAnalysis',
                    requiresAuth: true,
                },
            },
        ],
    },
    {
        path: '/images',
        component: Layout,
        meta: {
            title: '图像管理',
            icon: 'Picture',
            requiresAuth: true,
        },
        children: [
            {
                path: 'upload',
                name: 'ImageUpload',
                component: ImageUpload,
                meta: {
                    title: '图像上传',
                    icon: 'Upload',
                    requiresAuth: true,
                },
            },
            {
                path: 'gallery',
                name: 'ImageGallery',
                component: ImageGallery,
                meta: {
                    title: '图像列表',
                    icon: 'FolderOpened',
                    requiresAuth: true,
                },
            },
            {
                path: 'detail/:id',
                name: 'ImageDetail',
                component: ImageDetail,
                meta: {
                    title: '图像详情',
                    hideInMenu: true,
                    requiresAuth: true,
                },
            },
        ],
    },
    {
        path: '/analytics',
        component: Layout,
        redirect: '/analytics/index',
        children: [
            {
                path: 'index',
                name: 'Analytics',
                component: Analytics,
                meta: {
                    title: '数据分析',
                    icon: 'TrendCharts',
                    requiresAuth: true,
                },
            },
        ],
    },
    {
        path: '/profile',
        component: Layout,
        redirect: '/profile/index',
        children: [
            {
                path: 'index',
                name: 'Profile',
                component: Profile,
                meta: {
                    title: '个人中心',
                    icon: 'User',
                    requiresAuth: true,
                },
            },
        ],
    },
    {
        path: '/settings',
        component: Layout,
        redirect: '/settings/index',
        children: [
            {
                path: 'index',
                name: 'Settings',
                component: Settings,
                meta: {
                    title: '系统设置',
                    icon: 'Setting',
                    requiresAuth: true,
                },
            },
        ],
    },
    // Error pages
    {
        path: '/403',
        name: 'Forbidden',
        component: Forbidden,
        meta: {
            title: '访问禁止',
            hideInMenu: true,
        },
    },
    {
        path: '/500',
        name: 'ServerError',
        component: ServerError,
        meta: {
            title: '服务器错误',
            hideInMenu: true,
        },
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: NotFound,
        meta: {
            title: '页面不存在',
            hideInMenu: true,
        },
    },
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
    scrollBehavior(to, from, savedPosition) {
        if (savedPosition) {
            return savedPosition
        } else {
            return { top: 0 }
        }
    },
})

// Navigation guards
router.beforeEach(async (to, from, next) => {
    NProgress.start()

    // Set page title
    if (to.meta.title) {
        document.title = `${to.meta.title} - ICPT 图像处理系统`
    }

    const authStore = useAuthStore()
    const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

    if (requiresAuth) {
        if (!authStore.isAuthenticated) {
            // Check if token exists using proper token key
            const token = getToken()
            if (token) {
                try {
                    // Validate token and get user info
                    authStore.token = token  // Set token first
                    await authStore.getUserInfo()
                    next()
                } catch (error) {
                    console.error('Token validation failed:', error)
                    authStore.logout()
                    next('/login')
                }
            } else {
                next('/login')
            }
        } else {
            next()
        }
    } else {
        // If user is authenticated and tries to access login or register page
        if ((to.path === '/login' || to.path === '/register') && authStore.isAuthenticated) {
            next('/dashboard')
        } else {
            next()
        }
    }
})

router.afterEach(() => {
    NProgress.done()
})

// Error handling
router.onError((error) => {
    console.error('Router error:', error)
    NProgress.done()
})

export default router 