// Package camera 提供简化的摄像头采集功能
// 避免使用有问题的Aruco模块，专注核心摄像头功能
package camera

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"

	"gocv.io/x/gocv"
)

// SimpleCameraDevice 简化的摄像头设备信息
type SimpleCameraDevice struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SimpleCameraCapture 简化的摄像头采集器
type SimpleCameraCapture struct {
	webcam    *gocv.VideoCapture
	deviceID  int
	isRunning bool
	window    *gocv.Window
}

// ListSimpleCameras 枚举可用的摄像头设备 (简化版)
func ListSimpleCameras() ([]SimpleCameraDevice, error) {
	var cameras []SimpleCameraDevice

	// 尝试打开前5个摄像头ID
	for i := 0; i < 5; i++ {
		webcam, err := gocv.OpenVideoCapture(i)
		if err != nil {
			continue
		}

		if webcam.IsOpened() {
			cameras = append(cameras, SimpleCameraDevice{
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

// NewSimpleCameraCapture 创建简化摄像头采集器
func NewSimpleCameraCapture(deviceID int) (*SimpleCameraCapture, error) {
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		return nil, fmt.Errorf("无法打开摄像头 %d: %w", deviceID, err)
	}

	if !webcam.IsOpened() {
		return nil, fmt.Errorf("摄像头 %d 未能正确打开", deviceID)
	}

	return &SimpleCameraCapture{
		webcam:   webcam,
		deviceID: deviceID,
	}, nil
}

// StartSimplePreview 开始简单实时预览
func (c *SimpleCameraCapture) StartSimplePreview() error {
	if c.isRunning {
		return fmt.Errorf("预览已在运行")
	}

	c.window = gocv.NewWindow("简化摄像头预览 - 按 's' 拍照，按 'q' 退出")
	c.window.ResizeWindow(800, 600)
	c.isRunning = true

	img := gocv.NewMat()
	defer img.Close()

	fmt.Println("📹 简化摄像头预览已启动")
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
			filename := c.captureSimplePhoto(&img)
			if filename != "" {
				fmt.Printf("📷 拍照成功: %s\n", filename)
			}
		case 'q', 'Q', 27: // 退出 (ESC)
			c.isRunning = false
		}
	}

	return nil
}

// CaptureSimplePhoto 简单拍照并保存
func (c *SimpleCameraCapture) CaptureSimplePhoto(outputDir string) (string, error) {
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
	filename := fmt.Sprintf("simple_photo_%s.jpg", timestamp)
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

// captureSimplePhoto 内部简单拍照方法
func (c *SimpleCameraCapture) captureSimplePhoto(img *gocv.Mat) string {
	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("simple_photo_%s.jpg", timestamp)

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

// GetSimpleCameraInfo 获取简化摄像头信息
func (c *SimpleCameraCapture) GetSimpleCameraInfo() map[string]interface{} {
	info := make(map[string]interface{})

	info["device_id"] = c.deviceID
	info["width"] = c.webcam.Get(gocv.VideoCaptureFrameWidth)
	info["height"] = c.webcam.Get(gocv.VideoCaptureFrameHeight)
	info["fps"] = c.webcam.Get(gocv.VideoCaptureFPS)
	info["is_opened"] = c.webcam.IsOpened()

	return info
}

// SetSimpleResolution 设置简化摄像头分辨率
func (c *SimpleCameraCapture) SetSimpleResolution(width, height int) error {
	c.webcam.Set(gocv.VideoCaptureFrameWidth, float64(width))
	c.webcam.Set(gocv.VideoCaptureFrameHeight, float64(height))

	// 验证设置是否成功
	actualWidth := c.webcam.Get(gocv.VideoCaptureFrameWidth)
	actualHeight := c.webcam.Get(gocv.VideoCaptureFrameHeight)

	fmt.Printf("分辨率设置为: %.0fx%.0f\n", actualWidth, actualHeight)
	return nil
}

// StopSimplePreview 停止简单预览
func (c *SimpleCameraCapture) StopSimplePreview() {
	c.isRunning = false
	if c.window != nil {
		c.window.Close()
	}
}

// CloseSimple 关闭简化摄像头
func (c *SimpleCameraCapture) CloseSimple() error {
	c.StopSimplePreview()

	if c.webcam != nil {
		return c.webcam.Close()
	}

	return nil
}

// TestSimpleCamera 测试简化摄像头功能
func TestSimpleCamera() error {
	fmt.Println("🔍 测试简化摄像头功能...")

	// 1. 枚举摄像头
	cameras, err := ListSimpleCameras()
	if err != nil {
		return fmt.Errorf("枚举摄像头失败: %w", err)
	}

	fmt.Printf("✅ 发现 %d 个摄像头设备:\n", len(cameras))
	for _, cam := range cameras {
		fmt.Printf("  - ID: %d, 名称: %s\n", cam.ID, cam.Name)
	}

	if len(cameras) == 0 {
		return fmt.Errorf("未发现可用摄像头")
	}

	// 2. 测试第一个摄像头
	capture, err := NewSimpleCameraCapture(cameras[0].ID)
	if err != nil {
		return fmt.Errorf("打开摄像头失败: %w", err)
	}
	defer capture.CloseSimple()

	// 3. 获取摄像头信息
	info := capture.GetSimpleCameraInfo()
	fmt.Printf("✅ 摄像头信息: %+v\n", info)

	// 4. 快速拍照测试
	photoPath, err := capture.CaptureSimplePhoto(".")
	if err != nil {
		return fmt.Errorf("拍照失败: %w", err)
	}

	fmt.Printf("✅ 测试拍照成功: %s\n", photoPath)
	fmt.Println("🎉 简化摄像头功能测试完成！")

	return nil
} 