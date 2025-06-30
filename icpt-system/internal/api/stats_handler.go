package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"icpt-system/internal/models"
)

// DashboardStats 仪表盘统计数据结构
type DashboardStats struct {
	TotalImages    int64   `json:"total_images"`
	TodayProcessed int64   `json:"today_processed"`
	SuccessRate    float64 `json:"success_rate"`
	AvgTime        float64 `json:"avg_time"` // 平均处理时间(毫秒)
}

// RecentActivity 最近活动数据结构
type RecentActivity struct {
	Time   time.Time `json:"time"`
	Action string    `json:"action"`
	Image  string    `json:"image"`
	Status string    `json:"status"`
}

// GetDashboardStats 获取仪表盘统计信息
func GetDashboardStats(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		var stats DashboardStats

		// 1. 总图像数
		db.Model(&models.Image{}).Where("user_id = ?", userID).Count(&stats.TotalImages)

		// 2. 今日上传数量（创建的图像）
		todayStart := time.Now().Truncate(24 * time.Hour)
		todayEnd := todayStart.Add(24 * time.Hour)
		db.Model(&models.Image{}).
			Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, todayStart, todayEnd).
			Count(&stats.TodayProcessed)

		// 3. 成功率计算
		var totalProcessed, totalCompleted int64

		// 获取已完成处理的图像数量
		db.Model(&models.Image{}).
			Where("user_id = ? AND status = ?", userID, "completed").
			Count(&totalCompleted)

		// 获取总处理数量（已完成+失败）
		db.Model(&models.Image{}).
			Where("user_id = ? AND status IN ?", userID, []string{"completed", "failed"}).
			Count(&totalProcessed)

		if totalProcessed > 0 {
			stats.SuccessRate = (float64(totalCompleted) / float64(totalProcessed)) * 100
		} else {
			stats.SuccessRate = 0
		}

		// 4. 平均处理时间（单位：毫秒）
		type AvgResult struct {
			AvgTime float64
		}
		var avgResult AvgResult

		// 只计算已完成的图像的平均处理时间，使用微秒精度然后转换为毫秒
		err := db.Model(&models.Image{}).
			Select("AVG(TIMESTAMPDIFF(MICROSECOND, created_at, processed_at)) as avg_time").
			Where("user_id = ? AND status = ? AND processed_at IS NOT NULL", userID, "completed").
			Scan(&avgResult).Error

		if err == nil && avgResult.AvgTime > 0 {
			// 将微秒转换为毫秒（除以1000）
			stats.AvgTime = avgResult.AvgTime / 1000.0
		} else {
			stats.AvgTime = 0
		}

		c.JSON(http.StatusOK, gin.H{
			"data":    stats,
			"message": "获取统计信息成功",
		})
	}
}

// GetRecentActivity 获取最近活动
func GetRecentActivity(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		var images []models.Image
		db.Where("user_id = ? AND status IN ?", userID, []string{"completed", "failed"}).
			Order("processed_at DESC").
			Limit(10).
			Find(&images)

		var activities []RecentActivity
		for _, img := range images {
			status := "success"
			action := "图像处理"

			if img.Status == "failed" {
				status = "failed"
			}

			// 处理ProcessedAt为nil的情况
			activityTime := img.CreatedAt
			if img.ProcessedAt != nil {
				activityTime = *img.ProcessedAt
			}

			activities = append(activities, RecentActivity{
				Time:   activityTime,
				Action: action,
				Image:  img.OriginalFilename,
				Status: status,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"data":    activities,
			"message": "获取最近活动成功",
		})
	}
}

// GetImageStatusCount 获取图像状态统计
func GetImageStatusCount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")

		type StatusCount struct {
			Status string `json:"status"`
			Count  int64  `json:"count"`
		}

		var statusCounts []StatusCount
		db.Model(&models.Image{}).
			Select("status, COUNT(*) as count").
			Where("user_id = ?", userID).
			Group("status").
			Scan(&statusCounts)

		c.JSON(http.StatusOK, gin.H{
			"data":    statusCounts,
			"message": "获取状态统计成功",
		})
	}
}
 