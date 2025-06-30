import Cookies from 'js-cookie'

const TOKEN_KEY = 'icpt_token'
const REFRESH_TOKEN_KEY = 'icpt_refresh_token'
const TOKEN_EXPIRES = 7 // days

/**
 * Get authentication token from localStorage or cookies
 * @returns {string|null} Authentication token
 */
export const getToken = () => {
    // Try localStorage first (for SPA)
    let token = localStorage.getItem(TOKEN_KEY)

    // Fallback to cookies (for SSR or shared devices)
    if (!token) {
        token = Cookies.get(TOKEN_KEY)
    }

    return token
}

/**
 * Set authentication token to localStorage and cookies
 * @param {string} token - Authentication token
 */
export const setToken = (token) => {
    if (!token) return

    // Store in localStorage for SPA
    localStorage.setItem(TOKEN_KEY, token)

    // Store in secure cookie as backup
    Cookies.set(TOKEN_KEY, token, {
        expires: TOKEN_EXPIRES,
        secure: window.location.protocol === 'https:',
        sameSite: 'strict',
    })
}

/**
 * Remove authentication token from localStorage and cookies
 */
export const removeToken = () => {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
    Cookies.remove(TOKEN_KEY)
    Cookies.remove(REFRESH_TOKEN_KEY)
}

/**
 * Get refresh token
 * @returns {string|null} Refresh token
 */
export const getRefreshToken = () => {
    let refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)

    if (!refreshToken) {
        refreshToken = Cookies.get(REFRESH_TOKEN_KEY)
    }

    return refreshToken
}

/**
 * Set refresh token
 * @param {string} refreshToken - Refresh token
 */
export const setRefreshToken = (refreshToken) => {
    if (!refreshToken) return

    localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)

    Cookies.set(REFRESH_TOKEN_KEY, refreshToken, {
        expires: TOKEN_EXPIRES * 2, // Refresh token lives longer
        secure: window.location.protocol === 'https:',
        sameSite: 'strict',
    })
}

/**
 * Check if user is authenticated
 * @returns {boolean} Authentication status
 */
export const isAuthenticated = () => {
    const token = getToken()
    return !!token && !isTokenExpired(token)
}

/**
 * Check if token is expired
 * @param {string} token - JWT token
 * @returns {boolean} Whether token is expired
 */
export const isTokenExpired = (token) => {
    if (!token) return true

    try {
        // Decode JWT payload (simple base64 decode for expiry check)
        const payload = JSON.parse(atob(token.split('.')[1]))
        const currentTime = Date.now() / 1000

        return payload.exp < currentTime
    } catch (error) {
        console.error('Token decode error:', error)
        return true
    }
}

/**
 * Get token payload
 * @param {string} token - JWT token
 * @returns {Object|null} Token payload
 */
export const getTokenPayload = (token) => {
    if (!token) return null

    try {
        const payload = JSON.parse(atob(token.split('.')[1]))
        return payload
    } catch (error) {
        console.error('Token decode error:', error)
        return null
    }
}

/**
 * Get user ID from token
 * @returns {number|null} User ID
 */
export const getUserIdFromToken = () => {
    const token = getToken()
    const payload = getTokenPayload(token)
    return payload?.user_id || payload?.sub || null
}

/**
 * Get username from token
 * @returns {string|null} Username
 */
export const getUsernameFromToken = () => {
    const token = getToken()
    const payload = getTokenPayload(token)
    return payload?.username || null
}

/**
 * Calculate token remaining time in minutes
 * @param {string} token - JWT token
 * @returns {number} Remaining time in minutes
 */
export const getTokenRemainingTime = (token) => {
    if (!token) return 0

    try {
        const payload = getTokenPayload(token)
        const currentTime = Date.now() / 1000
        const remainingSeconds = payload.exp - currentTime

        return Math.max(0, Math.floor(remainingSeconds / 60))
    } catch (error) {
        console.error('Token time calculation error:', error)
        return 0
    }
}

/**
 * Clear all authentication data
 */
export const clearAuthData = () => {
    removeToken()

    // Clear other auth-related data
    localStorage.removeItem('user_info')
    localStorage.removeItem('user_preferences')

    // Clear auth-related cookies
    Cookies.remove('user_session')
}

/**
 * Store user preferences
 * @param {Object} preferences - User preferences
 */
export const setUserPreferences = (preferences) => {
    localStorage.setItem('user_preferences', JSON.stringify(preferences))
}

/**
 * Get user preferences
 * @returns {Object} User preferences
 */
export const getUserPreferences = () => {
    try {
        const preferences = localStorage.getItem('user_preferences')
        return preferences ? JSON.parse(preferences) : {}
    } catch (error) {
        console.error('Get user preferences error:', error)
        return {}
    }
} 