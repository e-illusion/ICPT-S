package models

import "time"

// Image 结构体对应 'images' 表
type Image struct {
	ID               uint       `gorm:"primaryKey"`
	UserID           uint       `gorm:"index"`
	OriginalFilename string     `gorm:"type:varchar(255);not null"`
	StoragePath      string     `gorm:"type:varchar(1024);not null"`
	ThumbnailPath    string     `gorm:"type:varchar(1024)"`
	Status           string     `gorm:"type:varchar(50);not null;default:'processing'"` // <-- 新增
	ErrorInfo        string     `gorm:"type:text"`                                      // <-- 新增
	FileSize         int64      `gorm:"type:bigint;default:0"`                          // <-- 新增文件大小字段
	ProcessedAt      *time.Time `gorm:"index"`                                          // <-- 新增处理完成时间字段
	CreatedAt        time.Time  `gorm:"autoCreateTime"`
}

// TableName 指定了此模型对应的数据库表名
func (Image) TableName() string {
	return "images"
}
