package api

import (
	"net/http"
	"strconv"

	"icpt-system/internal/websocket"

	"github.com/gin-gonic/gin"
)

// WebSocketHandler 处理WebSocket连接
func WebSocketHandler(c *gin.Context) {
	// 获取用户ID（从认证中间件）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户未认证",
		})
		return
	}

	// 升级到WebSocket连接
	websocket.ServeWS(websocket.GlobalHub, c.Writer, c.Request, userID.(uint))
}

// WebSocketStatsHandler 获取WebSocket连接统计
func WebSocketStatsHandler(c *gin.Context) {
	stats := websocket.GlobalHub.GetConnectionCount()

	c.JSON(http.StatusOK, gin.H{
		"message": "WebSocket统计信息",
		"data":    stats,
	})
}

// NotifyTestHandler 测试通知功能（开发用）
func NotifyTestHandler(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	// 发送测试通知
	websocket.GlobalHub.NotifyUser(uint(userID), websocket.SystemNotice, map[string]interface{}{
		"title":   "测试通知",
		"message": "这是一条测试消息",
		"type":    "info",
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "测试通知已发送",
		"user_id": userID,
	})
}
