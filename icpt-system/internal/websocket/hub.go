// Package websocket 提供WebSocket实时通知功能
package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Hub WebSocket连接管理器
type Hub struct {
	// 注册的客户端连接
	clients map[*Client]bool

	// 广播消息通道
	broadcast chan []byte

	// 注册客户端通道
	register chan *Client

	// 注销客户端通道
	unregister chan *Client

	// 用户订阅映射 (userID -> clients)
	userSubscriptions map[uint][]*Client

	// 互斥锁
	mutex sync.RWMutex
}

// Client WebSocket客户端
type Client struct {
	hub *Hub

	// WebSocket连接
	conn *websocket.Conn

	// 发送消息缓冲通道
	send chan []byte

	// 用户ID
	userID uint
}

// Message WebSocket消息结构
type Message struct {
	Type      string      `json:"type"`
	UserID    uint        `json:"user_id,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// NotificationType 通知类型
type NotificationType string

const (
	// 图像处理状态通知
	ImageProcessing NotificationType = "image_processing"
	ImageCompleted  NotificationType = "image_completed"
	ImageFailed     NotificationType = "image_failed"

	// 用户通知
	UserOnline  NotificationType = "user_online"
	UserOffline NotificationType = "user_offline"

	// 系统通知
	SystemNotice NotificationType = "system_notice"
)

// ImageNotification 图像处理通知
type ImageNotification struct {
	ImageID      uint   `json:"image_id"`
	Status       string `json:"status"`
	FileName     string `json:"file_name"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	ErrorInfo    string `json:"error_info,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决 an insecure WebSocket connection may not be initiated from a page loaded over HTTPS
	// 同时解决 a WebSocket handshake request but the client did not send an 'Upgrade' header
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		// 允许的来源列表，与main.go中的CORS配置保持一致
		allowedPatterns := []string{
			"http://localhost:",
			"https://localhost:",
			"http://127.0.0.1:",
			"https://127.0.0.1:",
			"http://[::1]:",
			"https://[::1]:",
			"http://114.55.58.3:",
			"https://114.55.58.3:",
			"http://0.0.0.0:",
			"http://114.55.58.3:3000",
		}

		for _, pattern := range allowedPatterns {
			if strings.HasPrefix(origin, pattern) {
				return true
			}
		}

		// 如果是空 origin（如直接访问或非浏览器客户端），也允许
		return origin == ""
	},
}

// NewHub 创建新的WebSocket Hub
func NewHub() *Hub {
	return &Hub{
		clients:           make(map[*Client]bool),
		broadcast:         make(chan []byte),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		userSubscriptions: make(map[uint][]*Client),
	}
}

// Run 启动Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient 注册客户端
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client] = true

	// 添加到用户订阅
	if client.userID > 0 {
		h.userSubscriptions[client.userID] = append(h.userSubscriptions[client.userID], client)
	}

	log.Printf("WebSocket客户端已连接，用户ID: %d，总连接数: %d", client.userID, len(h.clients))

	// 发送连接成功消息
	welcome := Message{
		Type: string(UserOnline),
		Data: map[string]interface{}{
			"message": "WebSocket连接成功",
			"user_id": client.userID,
		},
		Timestamp: getCurrentTimestamp(),
	}

	if data, err := json.Marshal(welcome); err == nil {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// unregisterClient 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)

		// 从用户订阅中移除
		if client.userID > 0 {
			clients := h.userSubscriptions[client.userID]
			for i, c := range clients {
				if c == client {
					h.userSubscriptions[client.userID] = append(clients[:i], clients[i+1:]...)
					break
				}
			}

			// 如果用户没有连接了，删除订阅
			if len(h.userSubscriptions[client.userID]) == 0 {
				delete(h.userSubscriptions, client.userID)
			}
		}

		log.Printf("WebSocket客户端已断开，用户ID: %d，总连接数: %d", client.userID, len(h.clients))
	}
}

// broadcastMessage 广播消息
func (h *Hub) broadcastMessage(message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// NotifyUser 向特定用户发送通知
func (h *Hub) NotifyUser(userID uint, notificationType NotificationType, data interface{}) {
	h.mutex.RLock()
	clients, exists := h.userSubscriptions[userID]
	h.mutex.RUnlock()

	if !exists || len(clients) == 0 {
		return
	}

	message := Message{
		Type:      string(notificationType),
		UserID:    userID,
		Data:      data,
		Timestamp: getCurrentTimestamp(),
	}

	messageData, err := json.Marshal(message)
	if err != nil {
		log.Printf("序列化消息失败: %v", err)
		return
	}

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for _, client := range clients {
		select {
		case client.send <- messageData:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}

	log.Printf("发送通知到用户 %d: %s", userID, notificationType)
}

// BroadcastAll 广播消息给所有用户
func (h *Hub) BroadcastAll(notificationType NotificationType, data interface{}) {
	message := Message{
		Type:      string(notificationType),
		Data:      data,
		Timestamp: getCurrentTimestamp(),
	}

	messageData, err := json.Marshal(message)
	if err != nil {
		log.Printf("序列化广播消息失败: %v", err)
		return
	}

	select {
	case h.broadcast <- messageData:
	default:
		log.Printf("广播通道已满，丢弃消息")
	}

	log.Printf("广播通知: %s", notificationType)
}

// GetConnectionCount 获取连接数统计
func (h *Hub) GetConnectionCount() map[string]int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return map[string]int{
		"total_connections":   len(h.clients),
		"authenticated_users": len(h.userSubscriptions),
	}
}

// NotifyImageProcessing 通知图像处理开始
func (h *Hub) NotifyImageProcessing(userID uint, imageID uint, fileName string) {
	notification := ImageNotification{
		ImageID:  imageID,
		Status:   "processing",
		FileName: fileName,
	}
	h.NotifyUser(userID, ImageProcessing, notification)
}

// NotifyImageCompleted 通知图像处理完成
func (h *Hub) NotifyImageCompleted(userID uint, imageID uint, fileName, thumbnailURL string) {
	notification := ImageNotification{
		ImageID:      imageID,
		Status:       "completed",
		FileName:     fileName,
		ThumbnailURL: thumbnailURL,
	}
	h.NotifyUser(userID, ImageCompleted, notification)
}

// NotifyImageFailed 通知图像处理失败
func (h *Hub) NotifyImageFailed(userID uint, imageID uint, fileName, errorInfo string) {
	notification := ImageNotification{
		ImageID:   imageID,
		Status:    "failed",
		FileName:  fileName,
		ErrorInfo: errorInfo,
	}
	h.NotifyUser(userID, ImageFailed, notification)
}

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// 全局Hub实例
var GlobalHub *Hub

// InitHub 初始化全局Hub
func InitHub() {
	GlobalHub = NewHub()
	go GlobalHub.Run()
	log.Println("WebSocket Hub 已启动")
}
