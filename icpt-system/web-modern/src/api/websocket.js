import { ElMessage } from 'element-plus'
import { getToken } from '@/utils/auth'

class WebSocketService {
    constructor() {
        this.ws = null
        this.reconnectAttempts = 0
        this.maxReconnectAttempts = 5
        this.reconnectInterval = 3000
        this.heartbeatInterval = 30000
        this.heartbeatTimer = null
        this.isConnecting = false
        this.listeners = new Map()
    }

    /**
     * Connect to WebSocket server
     * @param {string} [url] - WebSocket URL
     * @returns {Promise} Connection promise
     */
    connect(url) {
        if (this.isConnecting || (this.ws && this.ws.readyState === WebSocket.CONNECTING)) {
            return Promise.resolve()
        }

        this.isConnecting = true
        const token = getToken()

        if (!token) {
            console.warn('No authentication token available for WebSocket connection')
            this.isConnecting = false
            return Promise.reject(new Error('No authentication token'))
        }

        // Default WebSocket URL
        const wsUrl = url || this.getWebSocketUrl()
        const wsUrlWithAuth = `${wsUrl}?token=${encodeURIComponent(token)}`

        return new Promise((resolve, reject) => {
            try {
                this.ws = new WebSocket(wsUrlWithAuth)

                this.ws.onopen = () => {
                    console.log('WebSocket connected')
                    this.isConnecting = false
                    this.reconnectAttempts = 0
                    this.startHeartbeat()
                    this.emit('connected')
                    resolve()
                }

                this.ws.onmessage = (event) => {
                    try {
                        const data = JSON.parse(event.data)
                        this.handleMessage(data)
                    } catch (error) {
                        console.error('WebSocket message parse error:', error)
                    }
                }

                this.ws.onclose = (event) => {
                    console.log('WebSocket disconnected:', event.code, event.reason)
                    this.isConnecting = false
                    this.stopHeartbeat()
                    this.emit('disconnected', { code: event.code, reason: event.reason })

                    // Attempt to reconnect if not a clean close
                    if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
                        this.scheduleReconnect()
                    }
                }

                this.ws.onerror = (error) => {
                    console.error('WebSocket error:', error)
                    this.isConnecting = false
                    this.emit('error', error)
                    reject(error)
                }
            } catch (error) {
                this.isConnecting = false
                reject(error)
            }
        })
    }

    /**
     * Disconnect WebSocket
     */
    disconnect() {
        if (this.ws) {
            this.ws.close(1000, 'Client disconnect')
            this.ws = null
        }
        this.stopHeartbeat()
        this.reconnectAttempts = this.maxReconnectAttempts // Prevent auto-reconnect
    }

    /**
     * Send message to server
     * @param {Object} message - Message to send
     */
    send(message) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message))
        } else {
            console.warn('WebSocket is not connected')
        }
    }

    /**
     * Add event listener
     * @param {string} event - Event name
     * @param {Function} callback - Event callback
     */
    on(event, callback) {
        if (!this.listeners.has(event)) {
            this.listeners.set(event, [])
        }
        this.listeners.get(event).push(callback)
    }

    /**
     * Remove event listener
     * @param {string} event - Event name
     * @param {Function} callback - Event callback
     */
    off(event, callback) {
        if (this.listeners.has(event)) {
            const callbacks = this.listeners.get(event)
            const index = callbacks.indexOf(callback)
            if (index > -1) {
                callbacks.splice(index, 1)
            }
        }
    }

    /**
     * Emit event to listeners
     * @param {string} event - Event name
     * @param {*} data - Event data
     */
    emit(event, data) {
        if (this.listeners.has(event)) {
            this.listeners.get(event).forEach(callback => {
                try {
                    callback(data)
                } catch (error) {
                    console.error('WebSocket event listener error:', error)
                }
            })
        }
    }

    /**
     * Handle incoming WebSocket message
     * @param {Object} data - Message data
     */
    handleMessage(data) {
        const { type, data: messageData } = data

        switch (type) {
            case 'image_processing':
                ElMessage.info(`图像 ${messageData.file_name} 正在处理中...`)
                this.emit('image_processing', messageData)
                break

            case 'image_completed':
                ElMessage.success(`图像 ${messageData.file_name} 处理完成！`)
                this.emit('image_completed', messageData)
                break

            case 'image_failed':
                ElMessage.error(`图像 ${messageData.file_name} 处理失败`)
                this.emit('image_failed', messageData)
                break

            case 'notification':
                const messageType = messageData.level || 'info'
                ElMessage[messageType](messageData.message)
                this.emit('notification', messageData)
                break

            case 'system_notice':
                ElMessage.warning(messageData.message)
                this.emit('system_notice', messageData)
                break

            case 'heartbeat':
                // Server heartbeat response
                break

            case 'user_online':
                console.log('User online:', messageData)
                this.emit('user_online', messageData)
                break

            case 'user_offline':
                console.log('User offline:', messageData)
                this.emit('user_offline', messageData)
                break

            default:
                console.log('Unknown WebSocket message type:', type)
                this.emit('message', data)
        }
    }

    /**
     * Get WebSocket URL based on current location
     * @returns {string} WebSocket URL
     */
    getWebSocketUrl() {
        const baseUrl = import.meta.env.VITE_WS_URL

        if (baseUrl) {
            return baseUrl
        }

        // Default to current host with wss protocol (for HTTPS sites)
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const host = window.location.host
        return `${protocol}//${host}/api/v1/ws`
    }

    /**
     * Start heartbeat to keep connection alive
     */
    startHeartbeat() {
        this.stopHeartbeat()
        this.heartbeatTimer = setInterval(() => {
            if (this.ws && this.ws.readyState === WebSocket.OPEN) {
                this.send({ type: 'heartbeat', timestamp: Date.now() })
            }
        }, this.heartbeatInterval)
    }

    /**
     * Stop heartbeat timer
     */
    stopHeartbeat() {
        if (this.heartbeatTimer) {
            clearInterval(this.heartbeatTimer)
            this.heartbeatTimer = null
        }
    }

    /**
     * Schedule reconnection attempt
     */
    scheduleReconnect() {
        this.reconnectAttempts++
        const delay = this.reconnectInterval * this.reconnectAttempts

        console.log(`Scheduling WebSocket reconnect attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts} in ${delay}ms`)

        setTimeout(() => {
            if (this.reconnectAttempts < this.maxReconnectAttempts) {
                this.connect()
            }
        }, delay)
    }

    /**
     * Get connection status
     * @returns {boolean} Connection status
     */
    isConnected() {
        return this.ws && this.ws.readyState === WebSocket.OPEN
    }

    /**
     * Get connection state
     * @returns {string} Connection state
     */
    getState() {
        if (!this.ws) return 'disconnected'

        switch (this.ws.readyState) {
            case WebSocket.CONNECTING: return 'connecting'
            case WebSocket.OPEN: return 'connected'
            case WebSocket.CLOSING: return 'closing'
            case WebSocket.CLOSED: return 'disconnected'
            default: return 'unknown'
        }
    }
}

// Create singleton instance
const webSocketService = new WebSocketService()

export default webSocketService 