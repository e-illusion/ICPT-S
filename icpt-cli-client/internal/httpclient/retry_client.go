// Package httpclient 提供带重试机制的HTTP客户端
// 支持网络中断自动恢复和HTTPS安全传输
package httpclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries    int           // 最大重试次数
	InitialDelay  time.Duration // 初始延迟
	BackoffFactor float64       // 退避因子
	MaxDelay      time.Duration // 最大延迟
	Timeout       time.Duration // 请求超时
}

// DefaultRetryConfig 返回默认重试配置
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:    3,
		InitialDelay:  1 * time.Second,
		BackoffFactor: 2.0,
		MaxDelay:      30 * time.Second,
		Timeout:       30 * time.Second,
	}
}

// RetryableHTTPClient 可重试的HTTP客户端
type RetryableHTTPClient struct {
	httpClient *http.Client
	config     RetryConfig
}

// NewRetryableHTTPClient 创建新的可重试HTTP客户端
func NewRetryableHTTPClient(config RetryConfig, skipTLSVerify bool) *RetryableHTTPClient {
	// 配置TLS（支持HTTPS）
	tlsConfig := &tls.Config{
		InsecureSkipVerify: skipTLSVerify, // 开发环境可以跳过证书验证
	}

	// 创建HTTP客户端
	httpClient := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return &RetryableHTTPClient{
		httpClient: httpClient,
		config:     config,
	}
}

// Do 执行HTTP请求（带重试机制）
func (c *RetryableHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var lastErr error
	delay := c.config.InitialDelay

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		// 如果不是第一次尝试，等待一段时间
		if attempt > 0 {
			fmt.Printf("🔄 第%d次重试，等待%v...\n", attempt, delay)
			time.Sleep(delay)

			// 计算下一次重试的延迟（指数退避）
			delay = time.Duration(float64(delay) * c.config.BackoffFactor)
			if delay > c.config.MaxDelay {
				delay = c.config.MaxDelay
			}
		}

		// 克隆请求以支持重试（Body需要重新设置）
		reqClone := cloneRequest(req)

		// 执行请求
		resp, err := c.httpClient.Do(reqClone)
		if err == nil {
			// 检查HTTP状态码是否需要重试
			if !shouldRetryStatus(resp.StatusCode) {
				return resp, nil
			}
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
		} else {
			lastErr = err
		}

		// 记录重试原因
		if attempt < c.config.MaxRetries {
			fmt.Printf("⚠️  请求失败: %v\n", lastErr)
		}
	}

	return nil, fmt.Errorf("请求在%d次尝试后仍然失败: %v", c.config.MaxRetries+1, lastErr)
}

// shouldRetryStatus 判断HTTP状态码是否应该重试
func shouldRetryStatus(statusCode int) bool {
	// 以下状态码应该重试：
	// 408 Request Timeout
	// 429 Too Many Requests
	// 500 Internal Server Error
	// 502 Bad Gateway
	// 503 Service Unavailable
	// 504 Gateway Timeout
	switch statusCode {
	case 408, 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}

// cloneRequest 克隆HTTP请求以支持重试
func cloneRequest(req *http.Request) *http.Request {
	// 克隆请求
	reqClone := req.Clone(req.Context())

	// 如果有Body，需要重新设置
	if req.Body != nil {
		// 注意：这里假设Body是可以重复读取的
		// 对于文件上传等场景，调用方需要确保Body可以重置
		reqClone.Body = req.Body
	}

	return reqClone
}

// HTTPClientStats HTTP客户端统计信息
type HTTPClientStats struct {
	TotalRequests   int           // 总请求数
	SuccessRequests int           // 成功请求数
	FailedRequests  int           // 失败请求数
	TotalRetries    int           // 总重试次数
	AverageLatency  time.Duration // 平均延迟
}

// StatsTracker 统计跟踪器
type StatsTracker struct {
	stats HTTPClientStats
}

// NewStatsTracker 创建新的统计跟踪器
func NewStatsTracker() *StatsTracker {
	return &StatsTracker{}
}

// TrackRequest 跟踪请求统计
func (s *StatsTracker) TrackRequest(success bool, retries int, latency time.Duration) {
	s.stats.TotalRequests++
	if success {
		s.stats.SuccessRequests++
	} else {
		s.stats.FailedRequests++
	}
	s.stats.TotalRetries += retries

	// 计算平均延迟
	totalLatency := time.Duration(s.stats.TotalRequests-1)*s.stats.AverageLatency + latency
	s.stats.AverageLatency = totalLatency / time.Duration(s.stats.TotalRequests)
}

// GetStats 获取统计信息
func (s *StatsTracker) GetStats() HTTPClientStats {
	return s.stats
}

// PrintStats 打印统计信息
func (s *StatsTracker) PrintStats() {
	stats := s.GetStats()
	fmt.Printf("\n📊 HTTP客户端统计信息:\n")
	fmt.Printf("  总请求数: %d\n", stats.TotalRequests)
	fmt.Printf("  成功请求: %d\n", stats.SuccessRequests)
	fmt.Printf("  失败请求: %d\n", stats.FailedRequests)
	fmt.Printf("  总重试次数: %d\n", stats.TotalRetries)
	fmt.Printf("  平均延迟: %v\n", stats.AverageLatency)
	if stats.TotalRequests > 0 {
		successRate := float64(stats.SuccessRequests) / float64(stats.TotalRequests) * 100
		fmt.Printf("  成功率: %.1f%%\n", successRate)
	}
}
