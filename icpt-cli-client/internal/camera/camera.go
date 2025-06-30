// Package camera 提供摄像头采集功能
// 支持摄像头枚举、实时预览、拍照和录制
package camera

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"

	"gocv.io/x/gocv"
)

// CameraDevice 摄像头设备信息
type CameraDevice struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CameraCapture 摄像头采集器
type CameraCapture struct {
	webcam    *gocv.VideoCapture
	deviceID  int
	isRunning bool
	window    *gocv.Window
}

// ListCameras 枚举可用的摄像头设备
func ListCameras() ([]CameraDevice, error) {
	var cameras []CameraDevice

	// 尝试打开前5个摄像头ID
	for i := 0; i < 5; i++ {
		webcam, err := gocv.OpenVideoCapture(i)
		if err != nil {
			continue
		}

		if webcam.IsOpened() {
			cameras = append(cameras, CameraDevice{
				ID:   i,
				Name: fmt.Sprintf("Camera %d", i),
			})
		}
		webcam.Close()
	}

	if len(cameras) == 0 {
		return nil, fmt.Errorf("未检测到可用摄像头")
	}

	return cameras, nil
}

// NewCameraCapture 创建摄像头采集器
func NewCameraCapture(deviceID int) (*CameraCapture, error) {
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		return nil, fmt.Errorf("无法打开摄像头 %d: %w", deviceID, err)
	}

	if !webcam.IsOpened() {
		return nil, fmt.Errorf("摄像头 %d 未能正确打开", deviceID)
	}

	return &CameraCapture{
		webcam:   webcam,
		deviceID: deviceID,
	}, nil
}

// StartPreview 开始实时预览
func (c *CameraCapture) StartPreview() error {
	if c.isRunning {
		return fmt.Errorf("预览已在运行")
	}

	c.window = gocv.NewWindow("摄像头预览 - 按 's' 拍照，按 'q' 退出")
	c.window.ResizeWindow(800, 600)
	c.isRunning = true

	img := gocv.NewMat()
	defer img.Close()

	fmt.Println("📹 摄像头预览已启动")
	fmt.Println("提示：")
	fmt.Println("  - 按 's' 键拍照")
	fmt.Println("  - 按 'q' 键退出")

	for c.isRunning {
		if ok := c.webcam.Read(&img); !ok {
			return fmt.Errorf("无法从摄像头读取图像")
		}

		if img.Empty() {
			continue
		}

		c.window.IMShow(img)
		key := c.window.WaitKey(1)

		switch key {
		case 's', 'S': // 拍照
			filename := c.capturePhoto(&img)
			if filename != "" {
				fmt.Printf("📷 拍照成功: %s\n", filename)
			}
		case 'q', 'Q', 27: // 退出 (ESC)
			c.isRunning = false
		}
	}

	return nil
}

// StopPreview 停止预览
func (c *CameraCapture) StopPreview() {
	c.isRunning = false
	if c.window != nil {
		c.window.Close()
	}
}

// CapturePhoto 拍照并保存
func (c *CameraCapture) CapturePhoto(outputDir string) (string, error) {
	img := gocv.NewMat()
	defer img.Close()

	if ok := c.webcam.Read(&img); !ok {
		return "", fmt.Errorf("无法从摄像头捕获图像")
	}

	if img.Empty() {
		return "", fmt.Errorf("捕获的图像为空")
	}

	// 确保输出目录存在
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("photo_%s.jpg", timestamp)
	filepath := filepath.Join(outputDir, filename)

	// 转换为Go标准图像格式
	goImg, err := img.ToImage()
	if err != nil {
		return "", fmt.Errorf("图像格式转换失败: %w", err)
	}

	// 保存为JPEG文件
	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	if err := jpeg.Encode(file, goImg, &jpeg.Options{Quality: 95}); err != nil {
		return "", fmt.Errorf("保存图像失败: %w", err)
	}

	return filepath, nil
}

// capturePhoto 内部拍照方法（用于预览时拍照）
func (c *CameraCapture) capturePhoto(img *gocv.Mat) string {
	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("photo_%s.jpg", timestamp)

	// 转换为Go标准图像格式
	goImg, err := img.ToImage()
	if err != nil {
		fmt.Printf("❌ 图像转换失败: %v\n", err)
		return ""
	}

	// 保存为JPEG文件
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("❌ 创建文件失败: %v\n", err)
		return ""
	}
	defer file.Close()

	if err := jpeg.Encode(file, goImg, &jpeg.Options{Quality: 95}); err != nil {
		fmt.Printf("❌ 保存图像失败: %v\n", err)
		return ""
	}

	return filename
}

// CaptureQuickPhoto 快速拍照（不启动预览）
func (c *CameraCapture) CaptureQuickPhoto(outputDir string) (string, error) {
	return c.CapturePhoto(outputDir)
}

// GetCameraInfo 获取摄像头信息
func (c *CameraCapture) GetCameraInfo() map[string]interface{} {
	info := make(map[string]interface{})

	info["device_id"] = c.deviceID
	info["width"] = c.webcam.Get(gocv.VideoCaptureFrameWidth)
	info["height"] = c.webcam.Get(gocv.VideoCaptureFrameHeight)
	info["fps"] = c.webcam.Get(gocv.VideoCaptureFPS)
	info["is_opened"] = c.webcam.IsOpened()

	return info
}

// SetResolution 设置摄像头分辨率
func (c *CameraCapture) SetResolution(width, height int) error {
	c.webcam.Set(gocv.VideoCaptureFrameWidth, float64(width))
	c.webcam.Set(gocv.VideoCaptureFrameHeight, float64(height))

	// 验证设置是否成功
	actualWidth := c.webcam.Get(gocv.VideoCaptureFrameWidth)
	actualHeight := c.webcam.Get(gocv.VideoCaptureFrameHeight)

	fmt.Printf("分辨率设置为: %.0fx%.0f\n", actualWidth, actualHeight)
	return nil
}

// Close 关闭摄像头
func (c *CameraCapture) Close() error {
	c.StopPreview()

	if c.webcam != nil {
		return c.webcam.Close()
	}

	return nil
}

// RecordVideo 录制视频（基础实现）
func (c *CameraCapture) RecordVideo(outputPath string, durationSeconds int) error {
	// 获取摄像头参数
	width := int(c.webcam.Get(gocv.VideoCaptureFrameWidth))
	height := int(c.webcam.Get(gocv.VideoCaptureFrameHeight))
	fps := c.webcam.Get(gocv.VideoCaptureFPS)

	if fps <= 0 {
		fps = 30 // 默认30fps
	}

	// 创建视频写入器
	writer, err := gocv.VideoWriterFile(outputPath, "MJPG", fps, width, height, true)
	if err != nil {
		return fmt.Errorf("创建视频写入器失败: %w", err)
	}
	defer writer.Close()

	img := gocv.NewMat()
	defer img.Close()

	fmt.Printf("🎥 开始录制视频: %s\n", outputPath)
	fmt.Printf("录制时长: %d 秒\n", durationSeconds)

	startTime := time.Now()
	frameCount := 0

	for time.Since(startTime).Seconds() < float64(durationSeconds) {
		if ok := c.webcam.Read(&img); !ok {
			break
		}

		if img.Empty() {
			continue
		}

		writer.Write(img)
		frameCount++

		// 每秒显示一次进度
		if frameCount%int(fps) == 0 {
			elapsed := time.Since(startTime).Seconds()
			fmt.Printf("录制进度: %.1f/%.0f 秒\n", elapsed, float64(durationSeconds))
		}
	}

	fmt.Printf("✅ 录制完成，共录制 %d 帧\n", frameCount)
	return nil
}
