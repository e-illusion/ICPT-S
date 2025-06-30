import { get, post, del, upload } from './request'

// Image API endpoints
const IMAGE_ENDPOINTS = {
    UPLOAD: '/upload',
    LIST: '/images',
    DETAIL: '/images',
    DELETE: '/images',
    BATCH_DELETE: '/images/batch-delete',
    STATS: '/images/stats',
}

/**
 * Upload image file
 * @param {File} file - Image file to upload
 * @param {Function} onProgress - Upload progress callback
 * @returns {Promise<Object>} Upload response with image ID
 */
export const uploadImage = (file, onProgress) => {
    const formData = new FormData()
    formData.append('image', file)

    return upload(IMAGE_ENDPOINTS.UPLOAD, formData, {
        onProgress,
    })
}

/**
 * Upload multiple images
 * @param {FileList|Array} files - Image files to upload
 * @param {Function} onProgress - Upload progress callback
 * @returns {Promise<Array>} Array of upload responses
 */
export const uploadMultipleImages = async (files, onProgress) => {
    const uploads = Array.from(files).map((file, index) => {
        return uploadImage(file, (progress) => {
            if (onProgress) {
                onProgress(index, progress)
            }
        })
    })

    return Promise.all(uploads)
}

/**
 * Get user images list with pagination and filters
 * @param {Object} params - Query parameters
 * @param {number} [params.page=1] - Page number
 * @param {number} [params.page_size=10] - Items per page
 * @param {string} [params.status] - Filter by status (processing, completed, failed)
 * @param {string} [params.search] - Search term for filename
 * @param {string} [params.sort] - Sort field (created_at, filename, status)
 * @param {string} [params.order] - Sort order (asc, desc)
 * @returns {Promise<Object>} Paginated images list
 */
export const getImagesList = (params = {}) => {
    const defaultParams = {
        page: 1,
        page_size: 10,
        sort: 'created_at',
        order: 'desc',
    }

    return get(IMAGE_ENDPOINTS.LIST, { ...defaultParams, ...params })
}

/**
 * Get image details by ID
 * @param {number|string} imageId - Image ID
 * @returns {Promise<Object>} Image details
 */
export const getImageDetail = (imageId) => {
    return get(`${IMAGE_ENDPOINTS.DETAIL}/${imageId}`)
}

/**
 * Delete image by ID
 * @param {number|string} imageId - Image ID to delete
 * @returns {Promise<Object>} Delete response
 */
export const deleteImage = (imageId) => {
    return del(`${IMAGE_ENDPOINTS.DELETE}/${imageId}`)
}

/**
 * Batch delete multiple images
 * @param {Array<number>} imageIds - Array of image IDs to delete
 * @returns {Promise<Object>} Batch delete response
 */
export const batchDeleteImages = (imageIds) => {
    return post(IMAGE_ENDPOINTS.BATCH_DELETE, {
        image_ids: imageIds,
    })
}

/**
 * Get image processing statistics
 * @param {Object} params - Query parameters
 * @param {string} [params.period] - Time period (day, week, month, year)
 * @param {string} [params.start_date] - Start date (YYYY-MM-DD)
 * @param {string} [params.end_date] - End date (YYYY-MM-DD)
 * @returns {Promise<Object>} Statistics data
 */
export const getImageStats = (params = {}) => {
    return get(IMAGE_ENDPOINTS.STATS, params)
}

/**
 * Get image processing status count
 * @returns {Promise<Object>} Status count data
 */
export const getStatusCount = () => {
    return get('/images/status-count')
}

/**
 * Get recent activity
 * @param {number} limit - Number of recent activities to fetch
 * @returns {Promise<Array>} Recent activities
 */
export const getRecentActivity = (limit = 10) => {
    return get('/images/recent-activity', { limit })
}

/**
 * Search images by filename or metadata
 * @param {Object} searchParams - Search parameters
 * @param {string} searchParams.query - Search query
 * @param {Array<string>} [searchParams.filters] - Additional filters
 * @param {number} [searchParams.page] - Page number
 * @param {number} [searchParams.page_size] - Items per page
 * @returns {Promise<Object>} Search results
 */
export const searchImages = (searchParams) => {
    return get('/images/search', searchParams)
}

/**
 * Get image download URL
 * @param {number|string} imageId - Image ID
 * @param {string} [size] - Image size (original, thumbnail)
 * @returns {string} Download URL
 */
export const getImageUrl = (imageId, size = 'original') => {
    return `/static/images/${imageId}?size=${size}`
}

/**
 * Get image thumbnail URL from a relative path
 * @param {string} thumbnailPath - Thumbnail path from API response (e.g., "thumbnails/thumb-...")
 * @returns {string} Full, usable thumbnail URL
 */
export const getThumbnailUrl = (thumbnailPath) => {
    if (!thumbnailPath) return null
    // The backend now provides a relative path like "thumbnails/thumb-foo.jpg"
    // We prepend the '/static/' base to make it a valid URL
    return `/static/${thumbnailPath}`
}

/**
 * Get original image URL from a relative path
 * @param {string} originalPath - Original image path from API response (e.g., "originals/foo.jpg")
 * @returns {string} Full, usable original image URL
 */
export const getOriginalUrl = (originalPath) => {
    if (!originalPath) return null
    return `/static/${originalPath}`
}

/**
 * Update image metadata
 * @param {number|string} imageId - Image ID
 * @param {Object} metadata - New metadata
 * @returns {Promise<Object>} Update response
 */
export const updateImageMetadata = (imageId, metadata) => {
    return post(`/images/${imageId}/metadata`, metadata)
}

/**
 * Reprocess image (retry failed processing)
 * @param {number|string} imageId - Image ID
 * @returns {Promise<Object>} Reprocess response
 */
export const reprocessImage = (imageId) => {
    return post(`/images/${imageId}/reprocess`)
}

/**
 * Get image processing logs
 * @param {number|string} imageId - Image ID
 * @returns {Promise<Array>} Processing logs
 */
export const getImageLogs = (imageId) => {
    return get(`/images/${imageId}/logs`)
} 