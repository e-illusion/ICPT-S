package main

import (
	"context"
	"icpt-system/internal/config"
	"icpt-system/internal/models"
	"icpt-system/internal/services"
	"icpt-system/internal/store"
	"icpt-system/internal/websocket"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// 初始化所有组件，和 API 服务器一样
	config.LoadConfig("config.yaml")
	store.InitDB()
	store.InitRedis()

	// 初始化WebSocket Hub（用于发送通知）
	websocket.InitHub()

	os.MkdirAll("uploads/thumbnails", os.ModePerm) // 确保目录存在

	log.Println("🔧 后台 Worker 已启动，正在等待任务...")
	log.Println("✅ WebSocket通知已启用")

	for {
		// 使用 BRPOP 进行阻塞式读取，如果队列为空，会等待
		// "0" 表示无限期等待，直到有新任务
		result, err := store.Rdb.BRPop(context.Background(), 0, store.TaskQueueName).Result()
		if err != nil {
			log.Printf("错误: 从 Redis 队列获取任务失败: %v. 5秒后重试...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// result 是一个 []string, result[0] 是队列名, result[1] 是任务值 (我们的图片ID)
		imageIDStr := result[1]
		imageID, _ := strconv.ParseUint(imageIDStr, 10, 64)

		log.Printf("📋 接收到新任务, 图片 ID: %d", imageID)

		// ---- 执行真正的图像处理 ----
		var image models.Image
		// 根据 ID 从数据库中查找记录
		if err := store.DB.First(&image, imageID).Error; err != nil {
			log.Printf("错误: 无法在数据库中找到 ID 为 %d 的图片: %v", imageID, err)
			continue // 找不到记录，继续下一个任务
		}

		// 发送开始处理通知
		if websocket.GlobalHub != nil {
			websocket.GlobalHub.NotifyImageProcessing(image.UserID, image.ID, image.OriginalFilename)
		}

		// 调用我们之前写好的图像处理服务
		thumbPath, err := services.GenerateThumbnail(image.StoragePath, image.OriginalFilename)

		// ---- 更新数据库中的任务状态 ----
		now := time.Now()
		if err != nil {
			log.Printf("❌ 错误: 处理图片 (ID: %d) 失败: %v", imageID, err)
			// 更新数据库，将任务标记为失败
			store.DB.Model(&image).Updates(models.Image{
				Status:      "failed",
				ErrorInfo:   err.Error(),
				ProcessedAt: &now, // 设置处理完成时间
			})

			// 发送失败通知
			if websocket.GlobalHub != nil {
				websocket.GlobalHub.NotifyImageFailed(image.UserID, image.ID, image.OriginalFilename, err.Error())
			}
		} else {
			log.Printf("✅ 成功处理图片 (ID: %d), 缩略图路径: %s", imageID, thumbPath)
			// 更新数据库，写入缩略图路径并将状态标记为完成
			store.DB.Model(&image).Updates(models.Image{
				ThumbnailPath: thumbPath,
				Status:        "completed",
				ErrorInfo:     "",   // 清空错误信息
				ProcessedAt:   &now, // 设置处理完成时间
			})

			// 构建缩略图URL (去掉uploads/前缀以匹配静态文件配置)
			// thumbPath格式: uploads/thumbnails/thumb-xxx.jpg
			// 静态服务配置: /static -> ./uploads
			// 所以URL应该是: /static/thumbnails/thumb-xxx.jpg
			thumbnailPathForURL := strings.TrimPrefix(thumbPath, "uploads/")
			thumbnailURL := config.Cfg.Server.PublicHost + "/static/" + thumbnailPathForURL

			// 发送完成通知
			if websocket.GlobalHub != nil {
				websocket.GlobalHub.NotifyImageCompleted(image.UserID, image.ID, image.OriginalFilename, thumbnailURL)
			}
		}
	}
}
