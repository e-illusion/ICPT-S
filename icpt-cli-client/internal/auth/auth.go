package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"icpt-cli-client/internal/httpclient"
)

// AuthConfig 认证配置
type AuthConfig struct {
	ServerURL     string
	Timeout       time.Duration
	SkipTLSVerify bool // 是否跳过TLS证书验证（开发环境）
}

// AuthClient 认证客户端
type AuthClient struct {
	config       AuthConfig
	httpClient   *httpclient.RetryableHTTPClient
	statsTracker *httpclient.StatsTracker
	token        string
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	User  UserInfo `json:"user"`
	Token string   `json:"token"`
}

// APIResponse 通用API响应
type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"`
}

// NewAuthClient 创建新的认证客户端
func NewAuthClient(config AuthConfig) *AuthClient {
	// 创建重试配置
	retryConfig := httpclient.DefaultRetryConfig()
	retryConfig.Timeout = config.Timeout

	// 创建可重试的HTTP客户端
	httpClient := httpclient.NewRetryableHTTPClient(retryConfig, config.SkipTLSVerify)

	return &AuthClient{
		config:       config,
		httpClient:   httpClient,
		statsTracker: httpclient.NewStatsTracker(),
	}
}

// Register 用户注册
func (ac *AuthClient) Register(username, email, password string) (*AuthResponse, error) {
	reqData := map[string]string{
		"username": username,
		"email":    email,
		"password": password,
	}

	return ac.doAuthRequest("POST", "/api/v1/auth/register", reqData)
}

// Login 用户登录
func (ac *AuthClient) Login(username, password string) (*AuthResponse, error) {
	reqData := map[string]string{
		"username": username,
		"password": password,
	}

	authResp, err := ac.doAuthRequest("POST", "/api/v1/auth/login", reqData)
	if err == nil && authResp != nil {
		// 保存token
		ac.token = authResp.Token
	}
	return authResp, err
}

// GetProfile 获取用户信息
func (ac *AuthClient) GetProfile() (*UserInfo, error) {
	if ac.token == "" {
		return nil, fmt.Errorf("未登录，请先登录")
	}

	startTime := time.Now()
	req, err := http.NewRequest("GET", ac.config.ServerURL+"/api/v1/profile", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+ac.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ac.httpClient.Do(req)
	latency := time.Since(startTime)

	if err != nil {
		ac.statsTracker.TrackRequest(false, 0, latency)
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	ac.statsTracker.TrackRequest(true, 0, latency)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("服务器错误: %s", apiResp.Error)
	}

	// 解析用户信息
	userBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("解析用户数据失败: %w", err)
	}

	var user UserInfo
	if err := json.Unmarshal(userBytes, &user); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %w", err)
	}

	return &user, nil
}

// SetToken 设置认证令牌
func (ac *AuthClient) SetToken(token string) {
	ac.token = token
}

// GetToken 获取当前令牌
func (ac *AuthClient) GetToken() string {
	return ac.token
}

// IsLoggedIn 检查是否已登录
func (ac *AuthClient) IsLoggedIn() bool {
	return ac.token != ""
}

// GetStats 获取HTTP客户端统计信息
func (ac *AuthClient) GetStats() httpclient.HTTPClientStats {
	return ac.statsTracker.GetStats()
}

// PrintStats 打印HTTP客户端统计信息
func (ac *AuthClient) PrintStats() {
	ac.statsTracker.PrintStats()
}

// doAuthRequest 执行认证相关的HTTP请求
func (ac *AuthClient) doAuthRequest(method, path string, data interface{}) (*AuthResponse, error) {
	startTime := time.Now()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %w", err)
	}

	req, err := http.NewRequest(method, ac.config.ServerURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := ac.httpClient.Do(req)
	latency := time.Since(startTime)

	if err != nil {
		ac.statsTracker.TrackRequest(false, 0, latency)
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	ac.statsTracker.TrackRequest(true, 0, latency)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("服务器错误 (%d): %s", resp.StatusCode, apiResp.Error)
	}

	// 解析认证响应
	authBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("解析认证数据失败: %w", err)
	}

	var authResp AuthResponse
	if err := json.Unmarshal(authBytes, &authResp); err != nil {
		return nil, fmt.Errorf("解析认证响应失败: %w", err)
	}

	return &authResp, nil
}
