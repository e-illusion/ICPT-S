package services

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/nfnt/resize"
)

const (
	thumbnailPath = "uploads/thumbnails"
	thumbWidth    = 400
)

// GenerateThumbnail 接收一个原始文件路径，为其生成缩略图
// 返回: 缩略图存储路径, 错误
func GenerateThumbnail(originalFilePath string, originalFilename string) (string, error) {
	// 打开原始文件
	file, err := os.Open(originalFilePath)
	if err != nil {
		log.Printf("错误: 打开原始文件 %s 失败: %v", originalFilePath, err)
		return "", fmt.Errorf("无法打开原始文件")
	}
	defer file.Close()
	
	// 解码图像
	img, _, err := image.Decode(file)
	if err != nil {
		log.Printf("错误: 解码图像 %s 失败: %v", originalFilePath, err)
		return "", fmt.Errorf("无法解码图像，可能是不支持的格式")
	}

	// 生成缩略图
	thumbnail := resize.Resize(thumbWidth, 0, img, resize.Lanczos3)
	
	// 构造缩略图文件名和路径
	uniqueThumbFilename := fmt.Sprintf("thumb-%s-%s", time.Now().Format("20060102150405"), originalFilename)
	thumbnailFilePath := filepath.Join(thumbnailPath, uniqueThumbFilename)

	// 创建缩略图文件
	out, err := os.Create(thumbnailFilePath)
	if err != nil {
		log.Printf("错误: 创建缩略图文件 %s 失败: %v", thumbnailFilePath, err)
		return "", fmt.Errorf("无法保存缩略图")
	}
	defer out.Close()

	// 将缩略图以 JPEG 格式写入文件
	err = jpeg.Encode(out, thumbnail, nil)
	if err != nil {
		log.Printf("错误: 编码缩略图 %s 失败: %v", thumbnailFilePath, err)
		return "", fmt.Errorf("无法保存缩略图")
	}
	
	return thumbnailFilePath, nil
}
