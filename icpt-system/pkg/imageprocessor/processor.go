// Package imageprocessor 提供高性能C++图像处理功能
// 通过CGO调用OpenCV实现的C++库
package imageprocessor

/*
#cgo CXXFLAGS: -std=c++11 -I/usr/include/opencv4
#cgo LDFLAGS: -lopencv_core -lopencv_imgproc -lopencv_imgcodecs -lstdc++
#include "image_processor.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// ImageInfo 图像信息结构体
type ImageInfo struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Channels int    `json:"channels"`
	DataSize int    `json:"data_size"`
	Format   string `json:"format"`
}

// CompressConfig 压缩配置结构体
type CompressConfig struct {
	Quality      int  `json:"quality"`       // JPEG质量 (1-100)
	MaxWidth     int  `json:"max_width"`     // 最大宽度
	MaxHeight    int  `json:"max_height"`    // 最大高度
	EnableResize bool `json:"enable_resize"` // 是否启用尺寸调整
}

// DefaultCompressConfig 返回默认压缩配置
func DefaultCompressConfig() CompressConfig {
	return CompressConfig{
		Quality:      85,
		MaxWidth:     1920,
		MaxHeight:    1080,
		EnableResize: true,
	}
}

// HighQualityConfig 返回高质量配置
func HighQualityConfig() CompressConfig {
	return CompressConfig{
		Quality:      95,
		MaxWidth:     2560,
		MaxHeight:    1440,
		EnableResize: true,
	}
}

// ThumbnailConfig 返回缩略图配置
func ThumbnailConfig() CompressConfig {
	return CompressConfig{
		Quality:      75,
		MaxWidth:     400,
		MaxHeight:    300,
		EnableResize: true,
	}
}

// CompressImage 压缩图像文件
func CompressImage(inputPath, outputPath string, config CompressConfig) error {
	cInputPath := C.CString(inputPath)
	cOutputPath := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cInputPath))
	defer C.free(unsafe.Pointer(cOutputPath))

	cConfig := C.CompressConfig{
		quality:       C.int(config.Quality),
		max_width:     C.int(config.MaxWidth),
		max_height:    C.int(config.MaxHeight),
		enable_resize: C.bool(config.EnableResize),
	}

	result := C.compress_image(cInputPath, cOutputPath, &cConfig)
	if result != 0 {
		return fmt.Errorf("压缩失败，错误码: %d", int(result))
	}

	return nil
}

// GenerateThumbnail 生成缩略图
func GenerateThumbnail(inputPath, outputPath string, thumbWidth int) error {
	cInputPath := C.CString(inputPath)
	cOutputPath := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cInputPath))
	defer C.free(unsafe.Pointer(cOutputPath))

	result := C.generate_thumbnail(cInputPath, cOutputPath, C.int(thumbWidth))
	if result != 0 {
		return fmt.Errorf("生成缩略图失败，错误码: %d", int(result))
	}

	return nil
}

// GetImageInfo 获取图像信息
func GetImageInfo(inputPath string) (*ImageInfo, error) {
	cInputPath := C.CString(inputPath)
	defer C.free(unsafe.Pointer(cInputPath))

	var cInfo C.ImageInfo
	result := C.get_image_info(cInputPath, &cInfo)
	if result != 0 {
		return nil, fmt.Errorf("获取图像信息失败，错误码: %d", int(result))
	}

	info := &ImageInfo{
		Width:    int(cInfo.width),
		Height:   int(cInfo.height),
		Channels: int(cInfo.channels),
		DataSize: int(cInfo.data_size),
		Format:   C.GoString(&cInfo.format[0]),
	}

	return info, nil
}

// BatchProcessImages 批量处理图像
func BatchProcessImages(inputPaths, outputPaths []string, config CompressConfig) (int, error) {
	if len(inputPaths) != len(outputPaths) {
		return 0, fmt.Errorf("输入路径和输出路径数量不匹配")
	}

	count := len(inputPaths)
	if count == 0 {
		return 0, nil
	}

	// 转换为C字符串数组
	cInputPaths := make([]*C.char, count)
	cOutputPaths := make([]*C.char, count)

	for i, path := range inputPaths {
		cInputPaths[i] = C.CString(path)
	}
	for i, path := range outputPaths {
		cOutputPaths[i] = C.CString(path)
	}

	// 释放内存
	defer func() {
		for i := 0; i < count; i++ {
			C.free(unsafe.Pointer(cInputPaths[i]))
			C.free(unsafe.Pointer(cOutputPaths[i]))
		}
	}()

	cConfig := C.CompressConfig{
		quality:       C.int(config.Quality),
		max_width:     C.int(config.MaxWidth),
		max_height:    C.int(config.MaxHeight),
		enable_resize: C.bool(config.EnableResize),
	}

	result := C.batch_process_images(
		(**C.char)(unsafe.Pointer(&cInputPaths[0])),
		(**C.char)(unsafe.Pointer(&cOutputPaths[0])),
		C.int(count),
		&cConfig,
	)

	return int(result), nil
}

// ProcessImageMemory 在内存中处理图像
func ProcessImageMemory(inputData []byte, config CompressConfig) ([]byte, error) {
	if len(inputData) == 0 {
		return nil, fmt.Errorf("输入数据为空")
	}

	var outputData *C.uchar
	var outputSize C.size_t

	cConfig := C.CompressConfig{
		quality:       C.int(config.Quality),
		max_width:     C.int(config.MaxWidth),
		max_height:    C.int(config.MaxHeight),
		enable_resize: C.bool(config.EnableResize),
	}

	result := C.process_image_memory(
		(*C.uchar)(unsafe.Pointer(&inputData[0])),
		C.size_t(len(inputData)),
		&outputData,
		&outputSize,
		&cConfig,
	)

	if result != 0 {
		return nil, fmt.Errorf("内存中图像处理失败，错误码: %d", int(result))
	}

	// 转换C内存到Go切片
	output := C.GoBytes(unsafe.Pointer(outputData), C.int(outputSize))

	// 释放C分配的内存
	C.free_image_data(outputData)

	return output, nil
}

// GetVersion 获取库版本信息
func GetVersion() string {
	return C.GoString(C.get_version())
}

// GetOpenCVVersion 获取OpenCV版本信息
func GetOpenCVVersion() string {
	return C.GoString(C.get_opencv_version())
}

// IsAvailable 检查C++库是否可用
func IsAvailable() bool {
	version := GetVersion()
	return version != ""
}

// CompressImageWithQuality 使用指定质量压缩图像（快捷方法）
func CompressImageWithQuality(inputPath, outputPath string, quality int) error {
	config := DefaultCompressConfig()
	config.Quality = quality
	return CompressImage(inputPath, outputPath, config)
}

// GenerateThumbnailWithConfig 使用配置生成缩略图
func GenerateThumbnailWithConfig(inputPath, outputPath string, config CompressConfig) error {
	return CompressImage(inputPath, outputPath, config)
}

// ProcessorStats 处理器统计信息
type ProcessorStats struct {
	Version       string `json:"version"`
	OpenCVVersion string `json:"opencv_version"`
	Available     bool   `json:"available"`
}

// GetStats 获取处理器统计信息
func GetStats() ProcessorStats {
	return ProcessorStats{
		Version:       GetVersion(),
		OpenCVVersion: GetOpenCVVersion(),
		Available:     IsAvailable(),
	}
}

// ValidateConfig 验证压缩配置
func ValidateConfig(config CompressConfig) error {
	if config.Quality < 1 || config.Quality > 100 {
		return fmt.Errorf("质量参数必须在1-100之间")
	}
	if config.MaxWidth <= 0 || config.MaxHeight <= 0 {
		return fmt.Errorf("最大宽度和高度必须大于0")
	}
	return nil
}
