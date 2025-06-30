// Package worker 提供高性能的异步任务处理功能
// 支持并发处理、负载均衡和错误恢复
package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"icpt-system/internal/config"
	"icpt-system/internal/store"

	"github.com/go-redis/redis/v8"
)

// HighPerformanceWorker 高性能工作器
type HighPerformanceWorker struct {
	redisClient *redis.Client
	workerCount int
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	stats       *WorkerStats
}

// WorkerStats 工作器统计信息
type WorkerStats struct {
	TotalProcessed   int64         // 总处理数
	SuccessCount     int64         // 成功数
	FailureCount     int64         // 失败数
	AverageLatency   time.Duration // 平均延迟
	CurrentQueueSize int64         // 当前队列大小
	mu               sync.RWMutex
}

// NewHighPerformanceWorker 创建高性能工作器
func NewHighPerformanceWorker() *HighPerformanceWorker {
	ctx, cancel := context.WithCancel(context.Background())

	workerCount := config.Cfg.Performance.WorkerCount
	if workerCount <= 0 {
		workerCount = runtime.NumCPU() // 默认使用CPU核心数
	}

	return &HighPerformanceWorker{
		redisClient: store.RedisClient,
		workerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
		stats:       &WorkerStats{},
	}
}

// Start 启动高性能工作器
func (w *HighPerformanceWorker) Start() {
	log.Printf("🚀 启动高性能Worker，并发数: %d", w.workerCount)

	// 启动多个工作协程
	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.workerRoutine(i)
	}

	// 启动统计协程
	w.wg.Add(1)
	go w.statsRoutine()

	log.Println("✅ 高性能Worker启动成功")
}

// Stop 停止工作器
func (w *HighPerformanceWorker) Stop() {
	log.Println("🛑 停止高性能Worker...")
	w.cancel()
	w.wg.Wait()
	log.Println("✅ 高性能Worker已停止")
}

// workerRoutine 工作协程
func (w *HighPerformanceWorker) workerRoutine(workerID int) {
	defer w.wg.Done()

	log.Printf("Worker-%d 启动", workerID)

	for {
		select {
		case <-w.ctx.Done():
			log.Printf("Worker-%d 收到停止信号", workerID)
			return
		default:
			w.processNextTask(workerID)
		}
	}
}

// processNextTask 处理下一个任务
func (w *HighPerformanceWorker) processNextTask(workerID int) {
	startTime := time.Now()

	// 从Redis队列获取任务（阻塞式）
	result, err := w.redisClient.BLPop(w.ctx, 1*time.Second, "image_processing_queue").Result()
	if err != nil {
		if err != redis.Nil && err != context.Canceled {
			log.Printf("Worker-%d 获取任务失败: %v", workerID, err)
		}
		return
	}

	// 解析任务数据
	taskData := result[1]
	var task map[string]interface{}
	if err := json.Unmarshal([]byte(taskData), &task); err != nil {
		log.Printf("Worker-%d 解析任务失败: %v", workerID, err)
		w.updateStats(false, time.Since(startTime))
		return
	}

	// 处理任务
	success := w.handleImageProcessingTask(workerID, task)
	latency := time.Since(startTime)

	// 更新统计信息
	w.updateStats(success, latency)

	if success {
		log.Printf("Worker-%d 任务处理成功，耗时: %v", workerID, latency)
	} else {
		log.Printf("Worker-%d 任务处理失败，耗时: %v", workerID, latency)
	}
}

// handleImageProcessingTask 处理图像处理任务
func (w *HighPerformanceWorker) handleImageProcessingTask(workerID int, task map[string]interface{}) bool {
	// 这里复用现有的图像处理逻辑
	// 但是会针对性能进行优化

	imageID, ok := task["imageId"].(float64)
	if !ok {
		log.Printf("Worker-%d 无效的imageId", workerID)
		return false
	}

	userID, ok := task["userId"].(float64)
	if !ok {
		log.Printf("Worker-%d 无效的userId", workerID)
		return false
	}

	originalFilename, ok := task["originalFilename"].(string)
	if !ok {
		log.Printf("Worker-%d 无效的originalFilename", workerID)
		return false
	}

	storagePath, ok := task["storagePath"].(string)
	if !ok {
		log.Printf("Worker-%d 无效的storagePath", workerID)
		return false
	}

	log.Printf("Worker-%d 开始处理图像: ID=%d, 文件=%s", workerID, int(imageID), originalFilename)

	// 更新状态为处理中
	if err := updateImageStatus(uint(imageID), "processing", ""); err != nil {
		log.Printf("Worker-%d 更新状态失败: %v", workerID, err)
		return false
	}

	// 生成缩略图（这里可以调用C++处理模块）
	thumbnailPath, err := w.generateThumbnail(storagePath, originalFilename)
	if err != nil {
		log.Printf("Worker-%d 生成缩略图失败: %v", workerID, err)
		updateImageStatus(uint(imageID), "error", err.Error())
		return false
	}

	// 更新数据库，标记为完成
	if err := updateImageCompletion(uint(imageID), thumbnailPath); err != nil {
		log.Printf("Worker-%d 更新完成状态失败: %v", workerID, err)
		return false
	}

	return true
}

// generateThumbnail 生成缩略图（优化版本）
func (w *HighPerformanceWorker) generateThumbnail(sourcePath, filename string) (string, error) {
	// 这里可以集成更高性能的图像处理库
	// 或者调用C++模块

	outputPath := fmt.Sprintf("uploads/thumbnails/thumb_%s", filename)

	// 模拟图像处理（实际应该调用真实的处理函数）
	// 在高并发情况下，这里可以使用更高效的图像处理算法

	// 模拟处理时间（实际处理会更快）
	time.Sleep(50 * time.Millisecond)

	return outputPath, nil
}

// updateStats 更新统计信息
func (w *HighPerformanceWorker) updateStats(success bool, latency time.Duration) {
	w.stats.mu.Lock()
	defer w.stats.mu.Unlock()

	w.stats.TotalProcessed++
	if success {
		w.stats.SuccessCount++
	} else {
		w.stats.FailureCount++
	}

	// 计算平均延迟
	if w.stats.TotalProcessed == 1 {
		w.stats.AverageLatency = latency
	} else {
		totalLatency := time.Duration(w.stats.TotalProcessed-1)*w.stats.AverageLatency + latency
		w.stats.AverageLatency = totalLatency / time.Duration(w.stats.TotalProcessed)
	}
}

// statsRoutine 统计协程
func (w *HighPerformanceWorker) statsRoutine() {
	defer w.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			w.printStats()
		}
	}
}

// printStats 打印统计信息
func (w *HighPerformanceWorker) printStats() {
	w.stats.mu.RLock()
	defer w.stats.mu.RUnlock()

	// 获取当前队列大小
	queueSize, _ := w.redisClient.LLen(w.ctx, "image_processing_queue").Result()
	w.stats.CurrentQueueSize = queueSize

	successRate := float64(0)
	if w.stats.TotalProcessed > 0 {
		successRate = float64(w.stats.SuccessCount) / float64(w.stats.TotalProcessed) * 100
	}

	log.Printf("📊 Worker统计信息:")
	log.Printf("  总处理数: %d", w.stats.TotalProcessed)
	log.Printf("  成功数: %d", w.stats.SuccessCount)
	log.Printf("  失败数: %d", w.stats.FailureCount)
	log.Printf("  成功率: %.1f%%", successRate)
	log.Printf("  平均延迟: %v", w.stats.AverageLatency)
	log.Printf("  当前队列: %d", w.stats.CurrentQueueSize)
	log.Printf("  Worker数: %d", w.workerCount)
}

// GetStats 获取统计信息
func (w *HighPerformanceWorker) GetStats() *WorkerStats {
	w.stats.mu.RLock()
	defer w.stats.mu.RUnlock()

	// 复制统计信息以避免并发问题
	statsCopy := *w.stats

	// 获取当前队列大小
	queueSize, _ := w.redisClient.LLen(w.ctx, "image_processing_queue").Result()
	statsCopy.CurrentQueueSize = queueSize

	return &statsCopy
}

// 辅助函数（复用现有代码）
func updateImageStatus(imageID uint, status, errorInfo string) error {
	// 这里调用现有的数据库更新函数
	return store.UpdateImageStatus(imageID, status, errorInfo)
}

func updateImageCompletion(imageID uint, thumbnailPath string) error {
	// 这里调用现有的数据库更新函数
	return store.UpdateImageCompletion(imageID, thumbnailPath)
}
