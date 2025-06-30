import { get } from './request'

/**
 * 获取仪表盘统计信息
 * @returns {Promise<Object>} 统计数据
 */
export const getDashboardStats = () => {
    return get('/stats/dashboard')
}

/**
 * 获取最近活动
 * @returns {Promise<Array>} 最近活动列表
 */
export const getRecentActivity = () => {
    return get('/activity/recent')
}

/**
 * 获取图像状态统计
 * @returns {Promise<Array>} 状态统计数据
 */
export const getImageStatusCount = () => {
    return get('/stats/status-count')
} 