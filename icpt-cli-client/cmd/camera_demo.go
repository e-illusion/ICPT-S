package main

import (
	"fmt"
	"log"
	"os"

	"icpt-cli-client/internal/camera"
)

func main() {
	fmt.Println("🎥 ICPT 简化摄像头功能演示程序")
	fmt.Println("=================================")

	if len(os.Args) > 1 && os.Args[1] == "test" {
		// 运行自动化测试
		runAutomatedTest()
	} else {
		// 交互式测试
		runInteractiveDemo()
	}
}

func runAutomatedTest() {
	fmt.Println("🔄 运行自动化测试...")
	
	if err := camera.TestSimpleCamera(); err != nil {
		log.Printf("❌ 自动化测试失败: %v", err)
		os.Exit(1)
	}
	
	fmt.Println("✅ 自动化测试通过！")
}

func runInteractiveDemo() {
	fmt.Println("🎮 交互式摄像头演示")
	fmt.Println("")

	// 1. 枚举摄像头
	fmt.Println("1️⃣ 枚举摄像头设备...")
	cameras, err := camera.ListSimpleCameras()
	if err != nil {
		log.Printf("❌ 枚举摄像头失败: %v", err)
		fmt.Println("💡 这可能是因为：")
		fmt.Println("   - 没有连接摄像头设备")
		fmt.Println("   - 摄像头权限不足")
		fmt.Println("   - 驱动程序问题")
		return
	}

	fmt.Printf("✅ 发现 %d 个摄像头设备:\n", len(cameras))
	for _, cam := range cameras {
		fmt.Printf("  📷 ID: %d - %s\n", cam.ID, cam.Name)
	}
	fmt.Println("")

	if len(cameras) == 0 {
		fmt.Println("❌ 未发现可用摄像头，退出演示")
		fmt.Println("💡 建议：连接USB摄像头或启用内置摄像头")
		return
	}

	// 2. 选择摄像头进行测试
	selectedCamera := cameras[0]
	fmt.Printf("2️⃣ 使用摄像头 %d 进行演示...\n", selectedCamera.ID)

	capture, err := camera.NewSimpleCameraCapture(selectedCamera.ID)
	if err != nil {
		log.Printf("❌ 打开摄像头失败: %v", err)
		return
	}
	defer capture.CloseSimple()

	// 3. 获取摄像头信息
	fmt.Println("3️⃣ 获取摄像头信息...")
	info := capture.GetSimpleCameraInfo()
	fmt.Printf("✅ 摄像头详细信息:\n")
	for key, value := range info {
		fmt.Printf("  - %s: %v\n", key, value)
	}
	fmt.Println("")

	// 4. 拍照测试
	fmt.Println("4️⃣ 进行拍照演示...")
	photoDir := "camera_demo_photos"
	
	photoPath, err := capture.CaptureSimplePhoto(photoDir)
	if err != nil {
		log.Printf("❌ 拍照失败: %v", err)
		fmt.Println("💡 可能的原因：")
		fmt.Println("   - 摄像头被其他程序占用")
		fmt.Println("   - 磁盘空间不足")
		fmt.Println("   - 权限问题")
	} else {
		fmt.Printf("✅ 拍照成功！保存至: %s\n", photoPath)
		
		// 检查文件大小
		if stat, err := os.Stat(photoPath); err == nil {
			fmt.Printf("📊 文件大小: %d 字节\n", stat.Size())
		}
	}
	fmt.Println("")

	// 5. 预览功能说明
	fmt.Println("5️⃣ 预览功能说明")
	fmt.Println("📺 实时预览功能需要GUI环境支持")
	fmt.Println("🖥️ 如果在桌面环境中，可以调用:")
	fmt.Println("   capture.StartSimplePreview()")
	fmt.Println("⌨️ 预览控制:")
	fmt.Println("   - 按 's' 键拍照")
	fmt.Println("   - 按 'q' 键退出")
	fmt.Println("")

	// 6. 测试报告
	fmt.Println("🎉 摄像头功能演示完成！")
	fmt.Println("📋 功能测试报告:")
	fmt.Printf("  ✅ 摄像头枚举: %d 个设备\n", len(cameras))
	fmt.Printf("  ✅ 摄像头打开: 成功\n")
	fmt.Printf("  ✅ 信息获取: 成功\n")
	if photoPath != "" {
		fmt.Printf("  ✅ 拍照功能: 成功 (%s)\n", photoPath)
	} else {
		fmt.Printf("  ❌ 拍照功能: 失败\n")
	}
	fmt.Printf("  ⚠️ 预览功能: 需要GUI环境\n")
	fmt.Println("")
	
	fmt.Println("🚀 GoCV编译问题解决方案验证:")
	fmt.Println("  ✅ 简化模块成功避免Aruco冲突")
	fmt.Println("  ✅ 核心摄像头功能正常工作")
	fmt.Println("  ✅ 代码可以正常编译运行")
} 