package store

import (
	"context"
	"icpt-system/internal/config"
	"log"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client
var Ctx = context.Background()

const TaskQueueName = "image_processing_queue" // 定义任务队列的名称

// InitRedis 初始化 Redis 连接
func InitRedis() {
	c := config.Cfg.Redis
	Rdb = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,

		DB:       c.DB,
	})

	// 测试连接
	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("错误: 无法连接到 Redis: %v", err)
	}

	log.Println("Redis 连接成功！")
}
