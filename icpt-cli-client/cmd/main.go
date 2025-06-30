package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"icpt-cli-client/internal/auth"
	"icpt-cli-client/internal/camera"
	"icpt-cli-client/internal/compress"
	"icpt-cli-client/internal/config"

	"golang.org/x/term"
)

var authClient *auth.AuthClient

func main() {
	// 加载配置
	config.LoadConfig("config.yaml")

	// 创建认证客户端
	serverURL := config.Cfg.Server.PublicHost // 使用配置中的服务器地址

	authConfig := auth.AuthConfig{
		ServerURL:     serverURL,
		Timeout:       30 * time.Second,
		SkipTLSVerify: true, // 开发环境跳过TLS验证
	}

	authClient = auth.NewAuthClient(authConfig)

	// 检查命令行参数
	if len(os.Args) < 2 {
		showUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "register":
		handleRegister()
	case "login":
		handleLogin()
	case "upload":
		handleUpload()
	case "batch-upload":
		handleBatchUpload()
	case "list":
		handleList()
	case "status":
		handleStatus()
	case "profile":
		handleProfile()
	case "delete":
		handleDelete()
	case "info":
		handleInfo()
	case "compress":
		handleCompress()
	case "camera":
		handleCamera()
	case "help":
		showUsage()
	default:
		fmt.Printf("未知命令: %s\n", command)
		showUsage()
	}
}

func showUsage() {
	fmt.Println("🚀 ICPT 图像处理系统客户端")
	fmt.Println("========================")
	fmt.Println("使用方法:")
	fmt.Println("  register                    - 用户注册")
	fmt.Println("  login                       - 用户登录")
	fmt.Println("  upload <文件路径>           - 上传单个图像文件（自动压缩）")
	fmt.Println("  batch-upload <目录路径>     - 批量上传图像文件（自动压缩）")
	fmt.Println("  list [page] [page_size]     - 查看图像列表（分页）")
	fmt.Println("  status <图像ID>             - 查看图像处理状态")
	fmt.Println("  delete <图像ID>             - 删除指定图像")
	fmt.Println("  profile                     - 查看用户信息")
	fmt.Println("  info <文件路径>             - 查看图像文件信息")
	fmt.Println("  compress <文件路径> [质量]  - 压缩图像文件（1-100，默认75）")
	fmt.Println("")
	fmt.Println("📹 摄像头功能:")
	fmt.Println("  camera list                 - 列出可用摄像头")
	fmt.Println("  camera preview [设备ID]     - 开始摄像头预览（默认设备0）")
	fmt.Println("  camera capture [设备ID]     - 快速拍照（默认设备0）")
	fmt.Println("  camera upload [设备ID]      - 拍照并直接上传（默认设备0）")
	fmt.Println("  camera record <时长秒> [设备ID] - 录制视频")
	fmt.Println("")
	fmt.Println("  help                        - 显示帮助信息")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  ./cli-client register")
	fmt.Println("  ./cli-client login")
	fmt.Println("  ./cli-client upload image.jpg")
	fmt.Println("  ./cli-client batch-upload ./images/")
	fmt.Println("  ./cli-client camera list")
	fmt.Println("  ./cli-client camera preview 0")
	fmt.Println("  ./cli-client camera capture 0")
	fmt.Println("  ./cli-client camera upload 0")
	fmt.Println("  ./cli-client camera record 10 0")
	fmt.Println("  ./cli-client list 1 10")
	fmt.Println("  ./cli-client delete 123")
	fmt.Println("  ./cli-client info image.jpg")
	fmt.Println("  ./cli-client compress image.jpg 60")
}

func handleRegister() {
	fmt.Println("📝 用户注册")
	fmt.Println("============")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("用户名: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("邮箱: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("密码: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalf("读取密码失败: %v", err)
	}
	password := string(passwordBytes)
	fmt.Println() // 换行

	// 执行注册
	fmt.Println("正在注册...")
	authResp, err := authClient.Register(username, email, password)
	if err != nil {
		log.Fatalf("注册失败: %v", err)
	}

	fmt.Printf("✅ 注册成功！欢迎 %s\n", authResp.User.Username)
	fmt.Printf("用户ID: %d\n", authResp.User.ID)
	fmt.Printf("邮箱: %s\n", authResp.User.Email)
	fmt.Println("您已自动登录，可以开始使用其他功能。")
}

func handleLogin() {
	fmt.Println("🔐 用户登录")
	fmt.Println("============")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("用户名或邮箱: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("密码: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalf("读取密码失败: %v", err)
	}
	password := string(passwordBytes)
	fmt.Println() // 换行

	// 执行登录
	fmt.Println("正在登录...")
	authResp, err := authClient.Login(username, password)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}

	fmt.Printf("✅ 登录成功！欢迎回来 %s\n", authResp.User.Username)
	fmt.Printf("Token: %s...\n", authResp.Token[:20])
}

func handleUpload() {
	if len(os.Args) < 3 {
		fmt.Println("❌ 请指定要上传的文件路径")
		fmt.Println("使用方法: ./cli-client upload <文件路径>")
		return
	}

	if !authClient.IsLoggedIn() {
		fmt.Println("❌ 请先登录")
		fmt.Println("使用: ./cli-client login")
		return
	}

	filePath := os.Args[2]

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("文件不存在: %s", filePath)
	}

	fmt.Printf("📤 上传文件: %s\n", filePath)

	// 这里复用原来的上传逻辑，但添加认证header
	uploadResp, err := uploadFileWithAuth(filePath)
	if err != nil {
		log.Fatalf("上传失败: %v", err)
	}

	fmt.Printf("✅ 文件已接收，图片ID: %d\n", uploadResp.Data.ImageID)
	fmt.Println("开始查询处理状态...")

	// 轮询状态
	pollStatus(uploadResp.Data.ImageID)
}

func handleBatchUpload() {
	if len(os.Args) < 3 {
		fmt.Println("❌ 请指定要上传的目录路径")
		fmt.Println("使用方法: ./cli-client batch-upload <目录路径>")
		return
	}

	if !authClient.IsLoggedIn() {
		fmt.Println("❌ 请先登录")
		fmt.Println("使用: ./cli-client login")
		return
	}

	dirPath := os.Args[2]

	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		log.Fatalf("目录不存在: %s", dirPath)
	}

	fmt.Printf("📂 上传目录: %s\n", dirPath)

	// 遍历目录中的所有文件
	err := filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fmt.Printf("📤 上传文件: %s\n", filePath)

		// 这里复用原来的上传逻辑，但添加认证header
		uploadResp, err := uploadFileWithAuth(filePath)
		if err != nil {
			log.Fatalf("上传失败: %v", err)
		}

		fmt.Printf("✅ 文件已接收，图片ID: %d\n", uploadResp.Data.ImageID)
		fmt.Println("开始查询处理状态...")

		// 轮询状态
		pollStatus(uploadResp.Data.ImageID)

		return nil
	})

	if err != nil {
		log.Fatalf("遍历目录失败: %v", err)
	}
}

func handleList() {
	if !authClient.IsLoggedIn() {
		fmt.Println("❌ 请先登录")
		return
	}

	page := 1
	pageSize := 10

	if len(os.Args) >= 3 {
		if p, err := strconv.Atoi(os.Args[2]); err == nil {
			page = p
		}
	}

	if len(os.Args) >= 4 {
		if ps, err := strconv.Atoi(os.Args[3]); err == nil {
			pageSize = ps
		}
	}

	fmt.Printf("📋 查看图像列表 (第%d页，每页%d条)\n", page, pageSize)

	imageList, err := getImageList(page, pageSize)
	if err != nil {
		log.Fatalf("获取图像列表失败: %v", err)
	}

	if len(imageList.Data) == 0 {
		fmt.Println("📭 暂无图像")
		return
	}

	fmt.Printf("\n总计: %d 张图像，第 %d/%d 页\n", imageList.Total, imageList.Page, imageList.TotalPages)
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-5s %-20s %-12s %-19s %-30s\n", "ID", "文件名", "状态", "创建时间", "缩略图URL")
	fmt.Println(strings.Repeat("-", 80))

	for _, img := range imageList.Data {
		thumbnailURL := img.ThumbnailURL
		if thumbnailURL == "" {
			thumbnailURL = "处理中..."
		} else if len(thumbnailURL) > 30 {
			thumbnailURL = thumbnailURL[:27] + "..."
		}

		filename := img.OriginalFilename
		if len(filename) > 20 {
			filename = filename[:17] + "..."
		}

		fmt.Printf("%-5d %-20s %-12s %-19s %-30s\n",
			img.ID,
			filename,
			img.Status,
			img.CreatedAt,
			thumbnailURL,
		)
	}
}

func handleStatus() {
	if len(os.Args) < 3 {
		fmt.Println("❌ 请指定图像ID")
		fmt.Println("使用方法: ./cli-client status <图像ID>")
		return
	}

	if !authClient.IsLoggedIn() {
		fmt.Println("❌ 请先登录")
		return
	}

	imageID, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		log.Fatalf("无效的图像ID: %s", os.Args[2])
	}

	fmt.Printf("🔍 查询图像状态 (ID: %d)\n", imageID)

	status, err := getImageStatus(uint(imageID))
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	fmt.Printf("文件名: %s\n", status.OriginalFilename)
	fmt.Printf("状态: %s\n", status.Status)
	fmt.Printf("创建时间: %s\n", status.CreatedAt)

	if status.ThumbnailURL != "" {
		fmt.Printf("缩略图URL: %s\n", status.ThumbnailURL)
	}

	if status.ErrorInfo != "" {
		fmt.Printf("错误信息: %s\n", status.ErrorInfo)
	}
}

func handleProfile() {
	if !authClient.IsLoggedIn() {
		fmt.Println("❌ 请先登录")
		return
	}

	fmt.Println("👤 用户信息")
	fmt.Println("============")

	user, err := authClient.GetProfile()
	if err != nil {
		log.Fatalf("获取用户信息失败: %v", err)
	}

	fmt.Printf("用户ID: %d\n", user.ID)
	fmt.Printf("用户名: %s\n", user.Username)
	fmt.Printf("邮箱: %s\n", user.Email)
	fmt.Printf("状态: %s\n", user.Status)
}

func handleDelete() {
	if len(os.Args) < 3 {
		fmt.Println("❌ 请指定图像ID")
		fmt.Println("使用方法: ./cli-client delete <图像ID>")
		return
	}

	if !authClient.IsLoggedIn() {
		fmt.Println("❌ 请先登录")
		return
	}

	imageID, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		log.Fatalf("无效的图像ID: %s", os.Args[2])
	}

	fmt.Printf("🗑 删除图像 (ID: %d)\n", imageID)

	err = deleteImage(uint(imageID))
	if err != nil {
		log.Fatalf("删除失败: %v", err)
	}

	fmt.Println("✅ 图像删除成功")
}

func handleInfo() {
	if len(os.Args) < 3 {
		fmt.Println("❌ 请指定文件路径")
		fmt.Println("使用方法: ./cli-client info <文件路径>")
		return
	}

	filePath := os.Args[2]

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("文件不存在: %s", filePath)
	}

	fmt.Printf("📋 查看图像文件信息: %s\n", filePath)

	info, err := getImageInfo(filePath)
	if err != nil {
		log.Fatalf("获取图像文件信息失败: %v", err)
	}

	fmt.Printf("文件名: %s\n", info.OriginalFilename)
	fmt.Printf("文件大小: %d bytes\n", info.FileSize)
	fmt.Printf("创建时间: %s\n", info.CreatedAt)
	fmt.Printf("图像宽度: %d pixels\n", info.Width)
	fmt.Printf("图像高度: %d pixels\n", info.Height)
	fmt.Printf("图像格式: %s\n", info.Format)
}

func handleCompress() {
	if len(os.Args) < 3 {
		fmt.Println("❌ 请指定文件路径和质量")
		fmt.Println("使用方法: ./cli-client compress <文件路径> [质量]")
		return
	}

	filePath := os.Args[2]

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("文件不存在: %s", filePath)
	}

	fmt.Printf("📋 压缩图像文件: %s\n", filePath)

	// 解析质量参数
	quality := 75
	if len(os.Args) >= 4 {
		if q, err := strconv.Atoi(os.Args[3]); err == nil {
			quality = q
		}
	}

	// 执行压缩
	fmt.Println("正在压缩图像...")
	compressedFilePath, err := compressImage(filePath, quality)
	if err != nil {
		log.Fatalf("压缩图像失败: %v", err)
	}

	fmt.Printf("✅ 压缩成功！压缩后的文件保存到: %s\n", compressedFilePath)
}

// 实现带认证的文件上传功能（支持压缩）
func uploadFileWithAuth(filePath string) (*UploadResponse, error) {
	// 检查文件是否为图像格式
	if !isImageFile(filePath) {
		return nil, fmt.Errorf("不支持的文件格式，仅支持 JPEG 和 PNG 格式")
	}

	// 显示原始文件信息
	fmt.Printf("🖼️  分析图像文件: %s\n", filepath.Base(filePath))
	originalInfo, err := compress.GetImageInfo(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取图像信息失败: %w", err)
	}
	fmt.Printf("   原始信息: %s\n", originalInfo.String())

	// 自动压缩图像（如果文件大于1MB或尺寸过大）
	var finalFilePath string
	shouldCompress := originalInfo.Size > 1024*1024 || // 大于1MB
		originalInfo.Width > 1920 || originalInfo.Height > 1080 // 尺寸过大

	if shouldCompress {
		fmt.Println("🔧 正在压缩图像以优化传输...")
		compressedPath, err := compress.CompressImage(filePath, compress.DefaultConfig())
		if err != nil {
			fmt.Printf("⚠️  压缩失败，使用原始文件: %v\n", err)
			finalFilePath = filePath
		} else {
			finalFilePath = compressedPath
			// 清理函数，上传完成后删除临时压缩文件
			defer func() {
				if compressedPath != filePath {
					os.Remove(compressedPath)
				}
			}()
		}
	} else {
		fmt.Println("📁 文件大小合适，无需压缩")
		finalFilePath = filePath
	}

	// 检查最终文件是否存在
	file, err := os.Open(finalFilePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件 '%s': %w", finalFilePath, err)
	}
	defer file.Close()

	// 创建multipart form数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filepath.Base(filePath)) // 使用原始文件名
	if err != nil {
		return nil, fmt.Errorf("创建 form-data 失败: %w", err)
	}

	// 拷贝文件内容
	if _, err = io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("拷贝文件数据失败: %w", err)
	}
	writer.Close()

	// 创建HTTP请求
	fmt.Printf("🚀 上传文件到服务器: %s\n", config.Cfg.Server.PublicHost)
	uploadURL := config.Cfg.Server.PublicHost + "/api/v1/upload"
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return nil, fmt.Errorf("创建 HTTP 请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+authClient.GetToken())

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求到 '%s' 失败: %w", uploadURL, err)
	}
	defer resp.Body.Close()

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("服务器返回错误状态 %d: %s", resp.StatusCode, string(responseBody))
	}

	// 解析JSON响应
	var uploadResp UploadResponse
	if err := json.NewDecoder(bytes.NewReader(responseBody)).Decode(&uploadResp); err != nil {
		return nil, fmt.Errorf("无法解析上传响应的JSON: %w", err)
	}

	return &uploadResp, nil
}

// 检查文件是否为支持的图像格式
func isImageFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

// 获取本地图像文件信息（用于info命令）
func getImageInfo(filePath string) (*LocalImageInfo, error) {
	// 获取文件基本信息
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 获取图像详细信息
	imageInfo, err := compress.GetImageInfo(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取图像信息失败: %w", err)
	}

	return &LocalImageInfo{
		OriginalFilename: stat.Name(),
		FileSize:         stat.Size(),
		CreatedAt:        stat.ModTime().Format("2006-01-02 15:04:05"),
		Width:            imageInfo.Width,
		Height:           imageInfo.Height,
		Format:           imageInfo.Format,
	}, nil
}

// 压缩图像文件（用于compress命令）
func compressImage(filePath string, quality int) (string, error) {
	return compress.CompressImageWithQuality(filePath, quality)
}

// 实现状态轮询功能
func pollStatus(imageID uint) {
	const maxRetries = 15                 // 最多轮询15次
	const retryInterval = 2 * time.Second // 每次轮询间隔2秒

	for i := 0; i < maxRetries; i++ {
		time.Sleep(retryInterval) // 等待一下再查询

		statusResp, err := getImageStatus(imageID)
		if err != nil {
			fmt.Printf("警告: 第 %d 次查询状态失败: %v\n", i+1, err)
			continue
		}

		switch statusResp.Status {
		case "processing":
			fmt.Print(".") // 打印一个点表示仍在处理中
		case "completed":
			fmt.Printf("\n✅ 成功! 图像处理完成。\n")
			if statusResp.ThumbnailURL != "" {
				fmt.Printf("缩略图访问地址: %s\n", statusResp.ThumbnailURL)
			} else {
				fmt.Printf("缩略图路径: %s\n", statusResp.ThumbnailPath)
			}
			return // 任务完成，退出
		case "failed":
			fmt.Printf("\n❌ 失败! 图像处理失败。\n")
			if statusResp.ErrorInfo != "" {
				fmt.Printf("错误信息: %s\n", statusResp.ErrorInfo)
			}
			return // 任务失败，退出
		default:
			fmt.Printf("\n收到未知状态: %s\n", statusResp.Status)
		}
	}

	fmt.Printf("\n⚠️  超时: 在 %d 次尝试后仍未获得最终结果，请稍后手动查询。\n", maxRetries)
}

// 实现图像列表获取功能
func getImageList(page, pageSize int) (*PaginatedResponse, error) {
	// 构建请求URL
	listURL := fmt.Sprintf("%s/api/v1/images?page=%d&page_size=%d",
		config.Cfg.Server.PublicHost, page, pageSize)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置认证头
	req.Header.Set("Authorization", "Bearer "+authClient.GetToken())
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		var errorResp auth.APIResponse
		json.Unmarshal(body, &errorResp)
		return nil, fmt.Errorf("服务器错误 (%d): %s", resp.StatusCode, errorResp.Error)
	}

	// 解析响应
	var apiResp auth.APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 解析分页数据
	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("解析分页数据失败: %w", err)
	}

	var paginatedResp PaginatedResponse
	if err := json.Unmarshal(dataBytes, &paginatedResp); err != nil {
		return nil, fmt.Errorf("解析分页响应失败: %w", err)
	}

	return &paginatedResp, nil
}

// 实现图像状态查询功能
func getImageStatus(imageID uint) (*ImageStatusResponse, error) {
	// 构建请求URL
	statusURL := fmt.Sprintf("%s/api/v1/images/%d", config.Cfg.Server.PublicHost, imageID)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置认证头
	req.Header.Set("Authorization", "Bearer "+authClient.GetToken())
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		var errorResp auth.APIResponse
		json.Unmarshal(body, &errorResp)
		return nil, fmt.Errorf("服务器错误 (%d): %s", resp.StatusCode, errorResp.Error)
	}

	// 解析响应
	var apiResp auth.APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 解析图像状态数据
	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("解析状态数据失败: %w", err)
	}

	var statusResp ImageStatusResponse
	if err := json.Unmarshal(dataBytes, &statusResp); err != nil {
		return nil, fmt.Errorf("解析状态响应失败: %w", err)
	}

	return &statusResp, nil
}

// 实现图像删除功能
func deleteImage(imageID uint) error {
	// 构建请求URL
	deleteURL := fmt.Sprintf("%s/api/v1/images/%d", config.Cfg.Server.PublicHost, imageID)

	// 创建HTTP请求
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置认证头
	req.Header.Set("Authorization", "Bearer "+authClient.GetToken())
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		var errorResp auth.APIResponse
		json.Unmarshal(body, &errorResp)
		return fmt.Errorf("服务器错误 (%d): %s", resp.StatusCode, errorResp.Error)
	}

	return nil
}

// 数据结构定义
type UploadResponse struct {
	Data struct {
		ImageID uint   `json:"imageId"`
		Status  string `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}

type ImageListResponse struct {
	ID               uint   `json:"id"`
	OriginalFilename string `json:"original_filename"`
	Status           string `json:"status"`
	ThumbnailURL     string `json:"thumbnail_url,omitempty"`
	OriginalURL      string `json:"original_url,omitempty"`
	CreatedAt        string `json:"created_at"`
	ErrorInfo        string `json:"error_info,omitempty"`
}

type PaginatedResponse struct {
	Data       []ImageListResponse `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

type ImageStatusResponse struct {
	ID               uint   `json:"id"`
	Status           string `json:"status"`
	ThumbnailURL     string `json:"thumbnail_url"`
	ThumbnailPath    string `json:"thumbnail_path"`
	OriginalFilename string `json:"original_filename"`
	CreatedAt        string `json:"created_at"`
	ErrorInfo        string `json:"error_info"`
}

// LocalImageInfo 本地图像文件信息
type LocalImageInfo struct {
	OriginalFilename string `json:"original_filename"` // 文件名
	FileSize         int64  `json:"file_size"`         // 文件大小（字节）
	CreatedAt        string `json:"created_at"`        // 创建时间
	Width            int    `json:"width"`             // 图像宽度
	Height           int    `json:"height"`            // 图像高度
	Format           string `json:"format"`            // 图像格式
}

func handleCamera() {
	if len(os.Args) < 3 {
		showCameraUsage()
		return
	}

	subCommand := os.Args[2]

	switch subCommand {
	case "list":
		handleCameraList()
	case "preview":
		handleCameraPreview()
	case "capture":
		handleCameraCapture()
	case "upload":
		handleCameraUpload()
	case "record":
		handleCameraRecord()
	default:
		fmt.Printf("❌ 未知的摄像头命令: %s\n", subCommand)
		showCameraUsage()
	}
}

func showCameraUsage() {
	fmt.Println("📹 摄像头功能使用说明:")
	fmt.Println("  camera list                 - 列出可用摄像头")
	fmt.Println("  camera preview [设备ID]     - 开始摄像头预览（默认设备0）")
	fmt.Println("  camera capture [设备ID]     - 快速拍照（默认设备0）")
	fmt.Println("  camera upload [设备ID]      - 拍照并直接上传（默认设备0）")
	fmt.Println("  camera record <时长秒> [设备ID] - 录制视频")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  ./cli-client camera list")
	fmt.Println("  ./cli-client camera preview 0")
	fmt.Println("  ./cli-client camera capture 0")
	fmt.Println("  ./cli-client camera upload 0")
	fmt.Println("  ./cli-client camera record 10 0")
}

func handleCameraList() {
	fmt.Println("🔍 正在扫描可用摄像头...")

	cameras, err := camera.ListCameras()
	if err != nil {
		log.Fatalf("❌ 扫描摄像头失败: %v", err)
	}

	fmt.Printf("✅ 发现 %d 个摄像头设备:\n", len(cameras))
	fmt.Println("ID | 设备名称")
	fmt.Println("---|----------")
	for _, cam := range cameras {
		fmt.Printf("%2d | %s\n", cam.ID, cam.Name)
	}
}

func handleCameraPreview() {
	deviceID := 0 // 默认设备ID

	if len(os.Args) >= 4 {
		if id, err := strconv.Atoi(os.Args[3]); err == nil {
			deviceID = id
		}
	}

	fmt.Printf("📹 正在启动摄像头预览 (设备ID: %d)...\n", deviceID)

	capture, err := camera.NewCameraCapture(deviceID)
	if err != nil {
		log.Fatalf("❌ 无法打开摄像头: %v", err)
	}
	defer capture.Close()

	// 显示摄像头信息
	info := capture.GetCameraInfo()
	fmt.Printf("摄像头信息: %.0fx%.0f @ %.1f FPS\n",
		info["width"], info["height"], info["fps"])

	// 开始预览
	if err := capture.StartPreview(); err != nil {
		log.Fatalf("❌ 预览失败: %v", err)
	}

	fmt.Println("✅ 预览已结束")
}

func handleCameraCapture() {
	deviceID := 0 // 默认设备ID

	if len(os.Args) >= 4 {
		if id, err := strconv.Atoi(os.Args[3]); err == nil {
			deviceID = id
		}
	}

	fmt.Printf("📷 正在从摄像头拍照 (设备ID: %d)...\n", deviceID)

	capture, err := camera.NewCameraCapture(deviceID)
	if err != nil {
		log.Fatalf("❌ 无法打开摄像头: %v", err)
	}
	defer capture.Close()

	// 拍照
	outputDir := "./captures"
	photoPath, err := capture.CaptureQuickPhoto(outputDir)
	if err != nil {
		log.Fatalf("❌ 拍照失败: %v", err)
	}

	fmt.Printf("✅ 拍照成功: %s\n", photoPath)

	// 询问是否上传
	fmt.Print("是否要上传这张照片？(y/N): ")
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
		// 检查用户是否已登录
		if !authClient.IsLoggedIn() {
			fmt.Println("❌ 请先登录后再上传图片")
			fmt.Println("使用 './cli-client login' 命令登录")
			return
		}

		fmt.Println("🚀 正在上传照片...")
		uploadResp, err := uploadFileWithAuth(photoPath)
		if err != nil {
			log.Fatalf("❌ 上传失败: %v", err)
		}

		fmt.Printf("✅ 上传成功，图像ID: %d\n", uploadResp.Data.ImageID)
		fmt.Println("开始查询处理状态...")
		pollStatus(uploadResp.Data.ImageID)
	}
}

func handleCameraRecord() {
	if len(os.Args) < 4 {
		fmt.Println("❌ 请指定录制时长（秒）")
		fmt.Println("使用方法: ./cli-client camera record <时长秒> [设备ID]")
		return
	}

	duration, err := strconv.Atoi(os.Args[3])
	if err != nil || duration <= 0 {
		fmt.Println("❌ 录制时长必须是正整数")
		return
	}

	deviceID := 0 // 默认设备ID
	if len(os.Args) >= 5 {
		if id, err := strconv.Atoi(os.Args[4]); err == nil {
			deviceID = id
		}
	}

	fmt.Printf("🎥 正在录制视频 (设备ID: %d, 时长: %d秒)...\n", deviceID, duration)

	capture, err := camera.NewCameraCapture(deviceID)
	if err != nil {
		log.Fatalf("❌ 无法打开摄像头: %v", err)
	}
	defer capture.Close()

	// 设置分辨率
	capture.SetResolution(1280, 720)

	// 录制视频
	timestamp := time.Now().Format("20060102_150405")
	outputPath := fmt.Sprintf("./captures/video_%s.avi", timestamp)

	// 确保输出目录存在
	os.MkdirAll("./captures", os.ModePerm)

	if err := capture.RecordVideo(outputPath, duration); err != nil {
		log.Fatalf("❌ 录制失败: %v", err)
	}

	fmt.Printf("✅ 录制成功: %s\n", outputPath)
}

func handleCameraUpload() {
	deviceID := 0 // 默认设备ID

	if len(os.Args) >= 4 {
		if id, err := strconv.Atoi(os.Args[3]); err == nil {
			deviceID = id
		}
	}

	// 检查用户是否已登录
	if !authClient.IsLoggedIn() {
		fmt.Println("❌ 请先登录后再上传图片")
		fmt.Println("使用 './cli-client login' 命令登录")
		return
	}

	fmt.Printf("📷 正在从摄像头拍照并上传 (设备ID: %d)...\n", deviceID)

	capture, err := camera.NewCameraCapture(deviceID)
	if err != nil {
		log.Fatalf("❌ 无法打开摄像头: %v", err)
	}
	defer capture.Close()

	// 拍照
	outputDir := "./captures"
	photoPath, err := capture.CaptureQuickPhoto(outputDir)
	if err != nil {
		log.Fatalf("❌ 拍照失败: %v", err)
	}

	fmt.Printf("✅ 拍照成功: %s\n", photoPath)

	fmt.Println("🚀 正在上传照片...")
	uploadResp, err := uploadFileWithAuth(photoPath)
	if err != nil {
		log.Fatalf("❌ 上传失败: %v", err)
	}

	fmt.Printf("✅ 上传成功，图像ID: %d\n", uploadResp.Data.ImageID)
	fmt.Println("开始查询处理状态...")
	pollStatus(uploadResp.Data.ImageID)
}
