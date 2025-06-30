package store

import (
	"fmt"
	"log"
	"icpt-system/internal/config"
	"icpt-system/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	var err error
	// 从全局配置中获取数据库连接信息
	c := config.Cfg.Database

	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)

	// 使用 GORM 连接数据库
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("错误: 无法连接到数据库: %v", err)
	}

	log.Println("数据库连接成功！")

	// 自动迁移，确保数据库表结构与我们的模型定义一致
	// 这行代码会自动创建或更新 'images' 和 'users' 表
	err = DB.AutoMigrate(&models.Image{}, &models.User{})
	if err != nil {
		log.Fatalf("错误: 数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移成功！")
}
