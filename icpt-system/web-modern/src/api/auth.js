import { get, post } from './request'

// Auth API endpoints
const AUTH_ENDPOINTS = {
  LOGIN: '/auth/login',
  REGISTER: '/auth/register',
  PROFILE: '/profile',
  REFRESH: '/auth/refresh',
  LOGOUT: '/auth/logout',
}

/**
 * User login
 * @param {Object} credentials - Login credentials
 * @param {string} credentials.username - Username or email
 * @param {string} credentials.password - Password
 * @returns {Promise<Object>} Login response with token and user info
 */
export const login = (credentials) => {
  return post(AUTH_ENDPOINTS.LOGIN, {
    username: credentials.username,
    password: credentials.password,
  })
}

/**
 * User registration
 * @param {Object} userData - Registration data
 * @param {string} userData.username - Username
 * @param {string} userData.email - Email address
 * @param {string} userData.password - Password
 * @param {string} [userData.confirmPassword] - Password confirmation
 * @returns {Promise<Object>} Registration response with token and user info
 */
export const register = (userData) => {
  return post(AUTH_ENDPOINTS.REGISTER, {
    username: userData.username,
    email: userData.email,
    password: userData.password,
  })
}

/**
 * Get current user profile
 * @returns {Promise<Object>} User profile data
 */
export const getUserProfile = () => {
  return get(AUTH_ENDPOINTS.PROFILE)
}

/**
 * Update user profile
 * @param {Object} updates - Profile updates
 * @returns {Promise<Object>} Updated user profile
 */
export const updateUserProfile = (updates) => {
  return post(AUTH_ENDPOINTS.PROFILE, updates)
}

/**
 * Refresh authentication token
 * @returns {Promise<Object>} New token data
 */
export const refreshToken = () => {
  return post(AUTH_ENDPOINTS.REFRESH)
}

/**
 * User logout (if backend supports logout endpoint)
 * @returns {Promise<Object>} Logout response
 */
export const logout = () => {
  return post(AUTH_ENDPOINTS.LOGOUT)
}

/**
 * Validate if username is available
 * @param {string} username - Username to check
 * @returns {Promise<Object>} Availability check result
 */
export const checkUsernameAvailability = (username) => {
  return get('/auth/check-username', { username })
}

/**
 * Validate if email is available
 * @param {string} email - Email to check
 * @returns {Promise<Object>} Availability check result
 */
export const checkEmailAvailability = (email) => {
  return get('/auth/check-email', { email })
}

/**
 * Send password reset email
 * @param {string} email - Email address for password reset
 * @returns {Promise<Object>} Reset email response
 */
export const sendPasswordResetEmail = (email) => {
  return post('/auth/forgot-password', { email })
}

/**
 * Reset password with token
 * @param {Object} resetData - Reset password data
 * @param {string} resetData.token - Reset token
 * @param {string} resetData.password - New password
 * @returns {Promise<Object>} Password reset response
 */
export const resetPassword = (resetData) => {
  return post('/auth/reset-password', resetData)
}

/**
 * Change password (authenticated user)
 * @param {Object} passwordData - Password change data
 * @param {string} passwordData.currentPassword - Current password
 * @param {string} passwordData.newPassword - New password
 * @returns {Promise<Object>} Password change response
 */
export const changePassword = (passwordData) => {
  return post('/auth/change-password', passwordData)
} 