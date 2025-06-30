package main

import (
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
	"time"

	"golang.org/x/term"
	"gopkg.in/yaml.v2"
)

// Config 配置结构
type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		PublicHost string `yaml:"public_host"`
	} `yaml:"server"`
}

var config Config

// 响应结构体
type UploadResponse struct {
	Data struct {
		ImageID uint   `json:"imageId"`
		Status  string `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}

type StatusResponse struct {
	Data struct {
		ID            uint   `json:"id"`
		Status        string `json:"status"`
		ThumbnailURL  string `json:"thumbnail_url"`
		ErrorInfo     string `json:"error_info"`
	} `json:"data"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Message string `json:"message"`
}

type ImageListResponse struct {
	Data struct {
		Data []struct {
			ID               uint   `json:"id"`
			OriginalFilename string `json:"original_filename"`
			Status           string `json:"status"`
			ThumbnailURL     string `json:"thumbnail_url"`
			CreatedAt        string `json:"created_at"`
		} `json:"data"`
		Total int `json:"total"`
	} `json:"data"`
}

var authToken string

func main() {
	// 加载配置
	loadConfig()

	if len(os.Args) < 2 {
		showUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "login":
		handleLogin()
	case "upload":
		handleUpload()
	case "list":
		handleList()
	case "status":
		handleStatus()
	case "delete":
		handleDelete()
	case "help":
		showUsage()
	default:
		fmt.Printf("未知命令: %s\n", command)
		showUsage()
	}
}

func showUsage() {
	fmt.Println("ICPT 图像处理系统 - 简化客户端")
	fmt.Println("")
	fmt.Println("使用方法:")
	fmt.Println("  ./simple-client <命令> [参数...]")
	fmt.Println("")
	fmt.Println("可用命令:")
	fmt.Println("  login                       - 用户登录")
	fmt.Println("  upload <文件路径>           - 上传图像文件")
	fmt.Println("  list                        - 列出所有图像")
	fmt.Println("  status <图像ID>             - 查询图像处理状态")
	fmt.Println("  delete <图像ID>             - 删除图像")
	fmt.Println("  help                        - 显示帮助信息")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  ./simple-client login")
	fmt.Println("  ./simple-client upload image.jpg")
	fmt.Println("  ./simple-client list")
	fmt.Println("  ./simple-client status 1")
	fmt.Println("  ./simple-client delete 1")
}

func loadConfig() {
	configData, err := os.ReadFile("config.yaml")
	if err != nil {
		// 使用默认配置
		config.Server.Host = "http://localhost"
		config.Server.Port = ":8080"
		config.Server.PublicHost = "http://localhost:8080"
		return
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	// 构建完整的服务器地址
	if config.Server.PublicHost == "" {
		config.Server.PublicHost = config.Server.Host + config.Server.Port
	}
}

func handleLogin() {
	fmt.Print("用户名: ")
	var username string
	fmt.Scanln(&username)

	fmt.Print("密码: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("读取密码失败: %v", err)
	}
	password := string(passwordBytes)
	fmt.Println() // 换行

	// 构建登录请求
	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, _ := json.Marshal(loginData)
	resp, err := http.Post(config.Server.PublicHost+"/api/v1/auth/login", 
		"application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("登录请求失败: %v", err)
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		log.Fatalf("解析登录响应失败: %v", err)
	}

	if resp.StatusCode == 200 {
		authToken = loginResp.Token
		fmt.Printf("✅ 登录成功: %s\n", loginResp.Message)
		
		// 保存token到文件
		err = os.WriteFile(".auth_token", []byte(authToken), 0600)
		if err != nil {
			log.Printf("保存token失败: %v", err)
		} else {
			fmt.Println("🔐 认证令牌已保存")
		}
	} else {
		fmt.Printf("❌ 登录失败: %s\n", loginResp.Message)
	}
}

func loadToken() {
	if authToken != "" {
		return
	}
	
	tokenBytes, err := os.ReadFile(".auth_token")
	if err != nil {
		fmt.Println("❌ 未找到认证令牌，请先登录")
		os.Exit(1)
	}
	authToken = string(tokenBytes)
}

func makeAuthenticatedRequest(method, url string, body io.Reader) (*http.Response, error) {
	loadToken()
	
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+authToken)
	if method == "POST" && body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	client := &http.Client{Timeout: 30 * time.Second}
	return client.Do(req)
}

func handleUpload() {
	if len(os.Args) < 3 {
		fmt.Println("使用方法: ./simple-client upload <文件路径>")
		return
	}

	filePath := os.Args[2]
	
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("❌ 文件不存在: %s\n", filePath)
		return
	}

	fmt.Printf("📤 开始上传文件: %s\n", filePath)

	// 上传文件
	uploadResp, err := uploadFile(filePath)
	if err != nil {
		fmt.Printf("❌ 上传失败: %v\n", err)
		return
	}

	fmt.Printf("✅ %s\n", uploadResp.Message)
	fmt.Printf("📋 图像ID: %d\n", uploadResp.Data.ImageID)
	
	// 轮询处理状态
	fmt.Println("⏳ 正在处理中...")
	pollStatus(uploadResp.Data.ImageID)
}

func uploadFile(filePath string) (*UploadResponse, error) {
	loadToken()

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", filepath.Base(filePath))
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest("POST", config.Server.PublicHost+"/api/v1/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var uploadResp UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if resp.StatusCode != 202 {
		return nil, fmt.Errorf("服务器返回错误: %s", uploadResp.Message)
	}

	return &uploadResp, nil
}

func pollStatus(imageID uint) {
	for i := 0; i < 30; i++ { // 最多等待30次
		time.Sleep(2 * time.Second)
		
		resp, err := makeAuthenticatedRequest("GET", 
			fmt.Sprintf("%s/api/v1/images/%d", config.Server.PublicHost, imageID), nil)
		if err != nil {
			fmt.Printf("⚠️ 查询状态失败: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		var statusResp StatusResponse
		if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
			fmt.Printf("⚠️ 解析状态响应失败: %v\n", err)
			continue
		}

		switch statusResp.Data.Status {
		case "processing":
			fmt.Print(".")
		case "completed":
			fmt.Printf("\n✅ 处理完成！\n")
			if statusResp.Data.ThumbnailURL != "" {
				fmt.Printf("🖼️ 缩略图: %s\n", statusResp.Data.ThumbnailURL)
			}
			return
		case "failed":
			fmt.Printf("\n❌ 处理失败: %s\n", statusResp.Data.ErrorInfo)
			return
		default:
			fmt.Printf("\n⚠️ 未知状态: %s\n", statusResp.Data.Status)
		}
	}
	fmt.Printf("\n⏰ 处理超时，请稍后使用 status 命令查询\n")
}

func handleList() {
	resp, err := makeAuthenticatedRequest("GET", config.Server.PublicHost+"/api/v1/images", nil)
	if err != nil {
		fmt.Printf("❌ 获取图像列表失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var listResp ImageListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		fmt.Printf("❌ 解析响应失败: %v\n", err)
		return
	}

	if len(listResp.Data.Data) == 0 {
		fmt.Println("📭 暂无图像")
		return
	}

	fmt.Printf("📋 共 %d 个图像:\n", listResp.Data.Total)
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-4s %-20s %-12s %-40s %s\n", "ID", "文件名", "状态", "缩略图", "创建时间")
	fmt.Println(strings.Repeat("-", 80))

	for _, img := range listResp.Data.Data {
		thumbnailURL := img.ThumbnailURL
		if len(thumbnailURL) > 35 {
			thumbnailURL = thumbnailURL[:35] + "..."
		}
		fmt.Printf("%-4d %-20s %-12s %-40s %s\n", 
			img.ID, img.OriginalFilename, img.Status, thumbnailURL, img.CreatedAt)
	}
}

func handleStatus() {
	if len(os.Args) < 3 {
		fmt.Println("使用方法: ./simple-client status <图像ID>")
		return
	}

	imageID, err := strconv.ParseUint(os.Args[2], 10, 32)
	if err != nil {
		fmt.Printf("❌ 无效的图像ID: %s\n", os.Args[2])
		return
	}

	resp, err := makeAuthenticatedRequest("GET", 
		fmt.Sprintf("%s/api/v1/images/%d", config.Server.PublicHost, uint(imageID)), nil)
	if err != nil {
		fmt.Printf("❌ 查询状态失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var statusResp StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		fmt.Printf("❌ 解析响应失败: %v\n", err)
		return
	}

	fmt.Printf("📋 图像 ID: %d\n", statusResp.Data.ID)
	fmt.Printf("📊 状态: %s\n", statusResp.Data.Status)
	if statusResp.Data.ThumbnailURL != "" {
		fmt.Printf("🖼️ 缩略图: %s\n", statusResp.Data.ThumbnailURL)
	}
	if statusResp.Data.ErrorInfo != "" {
		fmt.Printf("❌ 错误信息: %s\n", statusResp.Data.ErrorInfo)
	}
}

func handleDelete() {
	if len(os.Args) < 3 {
		fmt.Println("使用方法: ./simple-client delete <图像ID>")
		return
	}

	imageID, err := strconv.ParseUint(os.Args[2], 10, 32)
	if err != nil {
		fmt.Printf("❌ 无效的图像ID: %s\n", os.Args[2])
		return
	}

	fmt.Printf("⚠️ 确认删除图像 %d? (y/N): ", uint(imageID))
	var confirm string
	fmt.Scanln(&confirm)
	
	if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
		fmt.Println("❌ 取消删除")
		return
	}

	resp, err := makeAuthenticatedRequest("DELETE", 
		fmt.Sprintf("%s/api/v1/images/%d", config.Server.PublicHost, uint(imageID)), nil)
	if err != nil {
		fmt.Printf("❌ 删除失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("✅ 图像 %d 删除成功\n", uint(imageID))
	} else {
		fmt.Printf("❌ 删除失败，状态码: %d\n", resp.StatusCode)
	}
}
 