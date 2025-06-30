// Package middleware 提供Gin中间件功能
package middleware

import (
	"net/http"
	"sync"
	"github.com/gin-gonic/gin"
	"fmt"
	"time"
)

// ConcurrencyLimitMiddleware 并发限制中间件
// 限制同时处理的请求数量，防止服务器过载
func ConcurrencyLimitMiddleware(maxConcurrency int) gin.HandlerFunc {
	// 创建一个带缓冲的channel作为信号量
	semaphore := make(chan struct{}, maxConcurrency)
	
	// 用于跟踪活跃连接数的计数器
	var activeConnections int64
	var mu sync.Mutex

	return func(c *gin.Context) {
		// 尝试获取信号量
		select {
		case semaphore <- struct{}{}:
			// 成功获取信号量，增加活跃连接数
			mu.Lock()
			activeConnections++
			current := activeConnections
			mu.Unlock()

			// 设置响应头显示当前负载
			c.Header("X-Active-Connections", fmt.Sprintf("%d", current))
			c.Header("X-Max-Concurrency", fmt.Sprintf("%d", maxConcurrency))

			// 继续处理请求
			c.Next()

			// 请求完成后释放信号量和减少连接数
			<-semaphore
			mu.Lock()
			activeConnections--
			mu.Unlock()

		default:
			// 无法获取信号量，返回503 Service Unavailable
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "服务器忙碌",
				"message": "当前并发请求过多，请稍后重试",
				"code":    "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}
	}
}

// RateLimitMiddleware 速率限制中间件 
// 基于令牌桶算法实现请求频率限制
func RateLimitMiddleware(requestsPerSecond int) gin.HandlerFunc {
	// 使用令牌桶算法
	type TokenBucket struct {
		tokens    int
		capacity  int
		rate      int // 每秒补充的令牌数
		lastRefill time.Time
		mu        sync.Mutex
	}

	bucket := &TokenBucket{
		tokens:     requestsPerSecond,
		capacity:   requestsPerSecond,
		rate:       requestsPerSecond,
		lastRefill: time.Now(),
	}

	return func(c *gin.Context) {
		bucket.mu.Lock()
		defer bucket.mu.Unlock()

		now := time.Now()
		elapsed := now.Sub(bucket.lastRefill)

		// 补充令牌
		tokensToAdd := int(elapsed.Seconds()) * bucket.rate
		if tokensToAdd > 0 {
			bucket.tokens = min(bucket.capacity, bucket.tokens+tokensToAdd)
			bucket.lastRefill = now
		}

		// 检查是否有足够的令牌
		if bucket.tokens > 0 {
			bucket.tokens--
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", bucket.tokens))
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", bucket.capacity))
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "请求频率过高",
				"message": "请求速度过快，请稍后重试",
				"code":    "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
		}
	}
}

// min 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
} 