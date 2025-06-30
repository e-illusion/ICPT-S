package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"icpt-system/internal/models"
	"icpt-system/internal/store"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ImageListResponse 图像列表响应结构
type ImageListResponse struct {
	ID               uint    `json:"id"`
	OriginalFilename string  `json:"original_filename"`
	Status           string  `json:"status"`
	ThumbnailURL     string  `json:"thumbnail_url,omitempty"`
	OriginalURL      string  `json:"original_url,omitempty"`
	FileSize         int64   `json:"file_size"`
	CreatedAt        string  `json:"created_at"`
	ProcessedAt      *string `json:"processed_at,omitempty"` // 新增处理时间字段
	ErrorInfo        string  `json:"error_info,omitempty"`
}

// PaginatedResponse 分页响应结构
type PaginatedResponse struct {
	Data       []ImageListResponse `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

// GetUserImagesHandler 获取用户的图像列表（分页）
func GetUserImagesHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户未认证",
			"code":  "UNAUTHENTICATED",
		})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status") // 可选的状态过滤

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 构建查询
	query := store.DB.Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	if err := query.Model(&models.Image{}).Count(&total).Error; err != nil {
		log.Printf("查询图像总数错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询失败",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// 分页查询
	var images []models.Image
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&images).Error; err != nil {
		log.Printf("查询图像列表错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询失败",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// 转换为响应格式
	imageList := make([]ImageListResponse, len(images))
	for i, img := range images {
		imageList[i] = ImageListResponse{
			ID:               img.ID,
			OriginalFilename: img.OriginalFilename,
			Status:           img.Status,
			FileSize:         img.FileSize,
			CreatedAt:        img.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		// 添加处理时间（如果存在）
		if img.ProcessedAt != nil {
			processedTime := img.ProcessedAt.Format("2006-01-02 15:04:05")
			imageList[i].ProcessedAt = &processedTime
		}

		// 如果有缩略图，添加URL (去掉开头的uploads/以匹配静态文件配置)
		if img.ThumbnailPath != "" {
			// 数据库存储: uploads/thumbnails/thumb-xxx.jpg
			// 静态服务配置: /static -> ./uploads
			// 所以API应该返回: thumbnails/thumb-xxx.jpg
			// 前端访问: /static/thumbnails/thumb-xxx.jpg -> ./uploads/thumbnails/thumb-xxx.jpg ✓
			thumbnailPath := strings.TrimPrefix(img.ThumbnailPath, "uploads/")
			imageList[i].ThumbnailURL = thumbnailPath
		}

		// 原始图片URL (同样去掉uploads/前缀)
		if img.StoragePath != "" {
			originalPath := strings.TrimPrefix(img.StoragePath, "uploads/")
			imageList[i].OriginalURL = originalPath
		}

		// 错误信息
		if img.ErrorInfo != "" {
			imageList[i].ErrorInfo = img.ErrorInfo
		}
	}

	// 计算总页数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	// 直接返回扁平化的响应结构，不再嵌套在data字段中
	c.JSON(http.StatusOK, gin.H{
		"message":     "查询成功",
		"data":        imageList, // 直接返回图像数组
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}

// DeleteImageHandler 删除用户的图像
func DeleteImageHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户未认证",
			"code":  "UNAUTHENTICATED",
		})
		return
	}

	imageID := c.Param("id")

	var image models.Image
	// 确保只能删除自己的图像
	result := store.DB.Where("id = ? AND user_id = ?", imageID, userID).First(&image)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "图像未找到",
				"code":  "IMAGE_NOT_FOUND",
			})
			return
		}
		log.Printf("查询图像错误: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "服务器内部错误",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// 删除物理文件
	var deletionErrors []string
	
	// 删除原始文件
	if image.StoragePath != "" {
		if err := os.Remove(image.StoragePath); err != nil && !os.IsNotExist(err) {
			log.Printf("删除原始文件失败 %s: %v", image.StoragePath, err)
			deletionErrors = append(deletionErrors, fmt.Sprintf("原始文件: %v", err))
		} else {
			log.Printf("已删除原始文件: %s", image.StoragePath)
		}
	}
	
	// 删除缩略图文件
	if image.ThumbnailPath != "" {
		if err := os.Remove(image.ThumbnailPath); err != nil && !os.IsNotExist(err) {
			log.Printf("删除缩略图文件失败 %s: %v", image.ThumbnailPath, err)
			deletionErrors = append(deletionErrors, fmt.Sprintf("缩略图: %v", err))
		} else {
			log.Printf("已删除缩略图文件: %s", image.ThumbnailPath)
		}
	}

	// 删除数据库记录
	if err := store.DB.Delete(&image).Error; err != nil {
		log.Printf("删除图像记录错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除失败",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	log.Printf("用户 %v 删除了图像 %d", userID, image.ID)

	response := gin.H{
		"message": "删除成功",
		"data": gin.H{
			"id": image.ID,
		},
	}

	// 如果有文件删除错误，添加警告信息
	if len(deletionErrors) > 0 {
		response["warnings"] = fmt.Sprintf("部分文件删除失败: %s", strings.Join(deletionErrors, ", "))
	}

	c.JSON(http.StatusOK, response)
}

// BatchDeleteImagesHandler 批量删除图像
func BatchDeleteImagesHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户未认证",
			"code":  "UNAUTHENTICATED",
		})
		return
	}

	var req struct {
		ImageIDs []uint `json:"image_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求数据格式错误",
			"code":    "INVALID_REQUEST",
			"details": err.Error(),
		})
		return
	}

	if len(req.ImageIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "图像ID列表不能为空",
			"code":  "EMPTY_IMAGE_LIST",
		})
		return
	}

	// 首先查询要删除的图像信息（用于删除物理文件）
	var images []models.Image
	if err := store.DB.Where("id IN ? AND user_id = ?", req.ImageIDs, userID).Find(&images).Error; err != nil {
		log.Printf("查询待删除图像错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询失败",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	if len(images) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "未找到可删除的图像",
			"code":  "NO_IMAGES_FOUND",
		})
		return
	}

	// 删除物理文件
	var deletionErrors []string
	filesDeleted := 0

	for _, image := range images {
		// 删除原始文件
		if image.StoragePath != "" {
			if err := os.Remove(image.StoragePath); err != nil && !os.IsNotExist(err) {
				log.Printf("删除原始文件失败 %s: %v", image.StoragePath, err)
				deletionErrors = append(deletionErrors, fmt.Sprintf("图像%d原始文件: %v", image.ID, err))
			} else {
				log.Printf("已删除原始文件: %s", image.StoragePath)
				filesDeleted++
			}
		}
		
		// 删除缩略图文件
		if image.ThumbnailPath != "" {
			if err := os.Remove(image.ThumbnailPath); err != nil && !os.IsNotExist(err) {
				log.Printf("删除缩略图文件失败 %s: %v", image.ThumbnailPath, err)
				deletionErrors = append(deletionErrors, fmt.Sprintf("图像%d缩略图: %v", image.ID, err))
			} else {
				log.Printf("已删除缩略图文件: %s", image.ThumbnailPath)
				filesDeleted++
			}
		}
	}

	// 删除数据库记录
	result := store.DB.Where("id IN ? AND user_id = ?", req.ImageIDs, userID).Delete(&models.Image{})
	if result.Error != nil {
		log.Printf("批量删除图像错误: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除失败",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	log.Printf("用户 %v 批量删除了 %d 个图像，删除了 %d 个物理文件", userID, result.RowsAffected, filesDeleted)

	response := gin.H{
		"message": "批量删除成功",
		"data": gin.H{
			"deleted_count": result.RowsAffected,
			"files_deleted": filesDeleted,
		},
	}

	// 如果有文件删除错误，添加警告信息
	if len(deletionErrors) > 0 {
		response["warnings"] = fmt.Sprintf("部分文件删除失败: %s", strings.Join(deletionErrors, "; "))
	}

	c.JSON(http.StatusOK, response)
} 
