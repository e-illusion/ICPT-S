package api

import (
	"icpt-system/internal/config"
	"icpt-system/internal/models"
	"icpt-system/internal/store"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ImageStatusResponse 定义查询响应的结构
type ImageStatusResponse struct {
	ID               uint   `json:"id"`
	Status           string `json:"status"`
	OriginalFilename string `json:"original_filename"`
	StoragePath      string `json:"storage_path"`
	ThumbnailPath    string `json:"thumbnail_path,omitempty"`
	ThumbnailURL     string `json:"thumbnail_url,omitempty"`
	ErrorInfo        string `json:"error_info,omitempty"`
	CreatedAt        string `json:"created_at"`
}

// GetImageStatusHandler 根据 ID 查询图片状态和信息
func GetImageStatusHandler(c *gin.Context) {
	// 从 URL 路径中获取 id
	id := c.Param("id")

	var image models.Image
	// 使用 First 方法根据主键查询
	result := store.DB.First(&image, id)

	// 检查查询结果
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 如果记录未找到，返回 404
			c.JSON(http.StatusNotFound, gin.H{
				"error": "图片未找到",
				"code":  "IMAGE_NOT_FOUND",
			})
		} else {
			// 如果是其他数据库错误，返回 500
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "查询数据库时发生错误",
				"code":  "DATABASE_ERROR",
			})
		}
		return
	}

	// 构建响应数据
	response := ImageStatusResponse{
		ID:               image.ID,
		Status:           image.Status,
		OriginalFilename: image.OriginalFilename,
		StoragePath:      image.StoragePath,
		CreatedAt:        image.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 如果有缩略图，添加缩略图信息和完整URL
	if image.ThumbnailPath != "" {
		response.ThumbnailPath = image.ThumbnailPath
		// 构建完整的访问URL (去掉uploads/前缀以匹配静态文件配置)
		baseURL := strings.TrimSuffix(config.Cfg.Server.PublicHost, "/")
		// 数据库存储: uploads/thumbnails/thumb-xxx.jpg
		// 静态服务配置: /static -> ./uploads
		// 所以URL应该是: /static/thumbnails/thumb-xxx.jpg
		thumbnailPath := strings.TrimPrefix(image.ThumbnailPath, "uploads/")
		response.ThumbnailURL = baseURL + "/static/" + thumbnailPath
	}

	// 如果有错误信息，添加错误信息
	if image.ErrorInfo != "" {
		response.ErrorInfo = image.ErrorInfo
	}

	// 成功找到记录，返回 200 和图片详细信息
	c.JSON(http.StatusOK, gin.H{
		"message": "查询成功",
		"data":    response,
	})
}
