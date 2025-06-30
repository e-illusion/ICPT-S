package api

import (
	"fmt"
	"icpt-system/internal/models"
	"icpt-system/internal/services"
	"icpt-system/internal/store"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// saveUploadedFile 负责保存原始文件并返回路径
func saveUploadedFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 沿用之前的唯一文件名逻辑
	uniqueFilename := fmt.Sprintf("%s-%s", time.Now().Format("20060102150405"), file.Filename)
	originalFilePath := filepath.Join("uploads/originals", uniqueFilename)

	dst, err := os.Create(originalFilePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return originalFilePath, err
}

func UploadImageHandler(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败: " + err.Error()})
		return
	}

	log.Printf("收到文件: %s, 大小: %d bytes", file.Filename, file.Size)

	// ---- 1. 只保存原始文件 ----
	originalPath, err := saveUploadedFile(file)
	if err != nil {
		log.Printf("错误: 保存原始文件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存文件"})
		return
	}

	// ---- 2. 获取当前用户ID ----
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// ---- 3. 在数据库中创建初始记录 ----
	imageRecord := models.Image{
		UserID:           userID.(uint),
		OriginalFilename: file.Filename,
		StoragePath:      originalPath,
		FileSize:         file.Size,    // 设置文件大小
		Status:           "processing", // 初始状态为 "处理中"
	}
	result := store.DB.Create(&imageRecord)
	if result.Error != nil {
		log.Printf("错误: 数据库创建初始记录失败: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建任务"})
		return
	}

	// ---- 4. 将任务推入 Redis 队列 ----
	err = store.Rdb.LPush(store.Ctx, store.TaskQueueName, imageRecord.ID).Err()
	if err != nil {
		log.Printf("错误: 推送任务到 Redis 失败: %v", err)
		// 注意：这里可以加入补偿逻辑，比如将任务标记为失败
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法调度任务"})
		return
	}

	log.Printf("图片记录 (ID: %d) 创建成功, 任务已推送到队列", imageRecord.ID)

	// ---- 5. 立即返回响应 ----
	c.JSON(http.StatusAccepted, gin.H{ // 返回 202 Accepted 表示请求已被接受，正在处理
		"message": "文件上传成功，正在后台处理中...",
		"data": gin.H{
			"imageId": imageRecord.ID,
			"status":  imageRecord.Status,
		},
	})
}

// UploadImageSyncHandlerForTest 是我们原始 MVP 处理器的副本，用于基准测试
func UploadImageSyncHandlerForTest(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败: " + err.Error()})
		return
	}

	// --- 这部分是整个同步工作负载 ---
	// 1. 保存原始文件
	originalPath, err := saveUploadedFile(file) // 我们可以复用我们的辅助函数
	if err != nil {
		log.Printf("同步测试错误：保存原始文件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "处理文件失败: " + err.Error()})
		return
	}

	// 2. 生成缩略图（慢部分）
	thumbPath, err := services.GenerateThumbnail(originalPath, file.Filename)
	if err != nil {
		log.Printf("同步测试错误：生成缩略图失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "处理文件失败: " + err.Error()})
		return
	}
	// --- 同步工作负载结束 ---

	imageRecord := models.Image{
		UserID:           1, // 硬编码的测试用户
		OriginalFilename: file.Filename,
		StoragePath:      originalPath,
		ThumbnailPath:    thumbPath,
		Status:           "completed_sync", // 使用不同的状态来标识这些记录
	}
	result := store.DB.Create(&imageRecord)
	if result.Error != nil {
		log.Printf("同步测试错误：保存到数据库失败: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存文件信息到数据库"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文件上传并处理成功！(同步)",
		"data": gin.H{
			"imageId": imageRecord.ID,
		},
	})
}
