// Package compress 提供图像压缩功能
// 支持JPEG质量压缩，减少传输数据量
package compress

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png" // 导入png解码器
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

// CompressConfig 压缩配置
type CompressConfig struct {
	Quality    int  // JPEG质量 (1-100)
	MaxWidth   uint // 最大宽度（像素）
	MaxHeight  uint // 最大高度（像素）
	EnableSize bool // 是否启用尺寸压缩
}

// DefaultConfig 返回默认压缩配置
func DefaultConfig() CompressConfig {
	return CompressConfig{
		Quality:    75,   // 75%质量，平衡质量和大小
		MaxWidth:   1920, // 最大1920像素宽度
		MaxHeight:  1080, // 最大1080像素高度
		EnableSize: true, // 启用尺寸压缩
	}
}

// CompressImage 压缩图像文件
// 输入：原始图像文件路径
// 输出：压缩后的临时文件路径和错误信息
func CompressImage(inputPath string, config CompressConfig) (string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return "", fmt.Errorf("文件不存在: %s", inputPath)
	}

	// 打开原始图像文件
	file, err := os.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	// 解码图像
	img, format, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("无法解码图像: %w", err)
	}

	// 获取原始尺寸
	bounds := img.Bounds()
	width := uint(bounds.Dx())
	height := uint(bounds.Dy())

	// 尺寸压缩（如果启用且超过限制）
	var resizedImg image.Image = img
	if config.EnableSize {
		if width > config.MaxWidth || height > config.MaxHeight {
			// 计算缩放比例，保持宽高比
			ratioW := float64(config.MaxWidth) / float64(width)
			ratioH := float64(config.MaxHeight) / float64(height)
			ratio := ratioW
			if ratioH < ratioW {
				ratio = ratioH
			}

			newWidth := uint(float64(width) * ratio)
			newHeight := uint(float64(height) * ratio)

			// 使用高质量插值算法调整尺寸
			resizedImg = resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
			fmt.Printf("图像尺寸已调整: %dx%d -> %dx%d\n", width, height, newWidth, newHeight)
		}
	}

	// 生成压缩后的临时文件路径
	ext := strings.ToLower(filepath.Ext(inputPath))
	outputPath := strings.TrimSuffix(inputPath, ext) + "_compressed.jpg"

	// 创建输出文件
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("无法创建输出文件: %w", err)
	}
	defer outFile.Close()

	// 根据原始格式和配置进行压缩
	switch format {
	case "jpeg":
		// JPEG格式，应用质量压缩
		options := &jpeg.Options{Quality: config.Quality}
		err = jpeg.Encode(outFile, resizedImg, options)
		if err != nil {
			return "", fmt.Errorf("JPEG编码失败: %w", err)
		}
	case "png":
		// PNG格式，转换为JPEG以减少大小
		options := &jpeg.Options{Quality: config.Quality}
		err = jpeg.Encode(outFile, resizedImg, options)
		if err != nil {
			return "", fmt.Errorf("PNG转JPEG编码失败: %w", err)
		}
	default:
		// 其他格式，尝试转换为JPEG
		options := &jpeg.Options{Quality: config.Quality}
		err = jpeg.Encode(outFile, resizedImg, options)
		if err != nil {
			return "", fmt.Errorf("图像编码失败: %w", err)
		}
	}

	// 获取文件大小信息
	originalStat, _ := os.Stat(inputPath)
	compressedStat, _ := os.Stat(outputPath)

	originalSize := originalStat.Size()
	compressedSize := compressedStat.Size()
	ratio := float64(compressedSize) / float64(originalSize) * 100

	fmt.Printf("压缩完成: %d bytes -> %d bytes (%.1f%%)\n",
		originalSize, compressedSize, ratio)

	return outputPath, nil
}

// CompressImageWithQuality 使用指定质量压缩图像（快捷方法）
func CompressImageWithQuality(inputPath string, quality int) (string, error) {
	config := DefaultConfig()
	config.Quality = quality
	return CompressImage(inputPath, config)
}

// GetImageInfo 获取图像文件信息
func GetImageInfo(imagePath string) (*ImageInfo, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	// 解码图像获取信息（不加载完整图像数据）
	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, fmt.Errorf("无法解码图像配置: %w", err)
	}

	// 获取文件大小
	stat, err := os.Stat(imagePath)
	if err != nil {
		return nil, fmt.Errorf("无法获取文件信息: %w", err)
	}

	return &ImageInfo{
		Width:  config.Width,
		Height: config.Height,
		Format: format,
		Size:   stat.Size(),
	}, nil
}

// ImageInfo 图像信息结构
type ImageInfo struct {
	Width  int    `json:"width"`  // 图像宽度
	Height int    `json:"height"` // 图像高度
	Format string `json:"format"` // 图像格式 (jpeg, png, etc.)
	Size   int64  `json:"size"`   // 文件大小（字节）
}

// String 返回图像信息的字符串表示
func (info *ImageInfo) String() string {
	return fmt.Sprintf("%dx%d %s (%.1f KB)",
		info.Width, info.Height, strings.ToUpper(info.Format),
		float64(info.Size)/1024)
}
