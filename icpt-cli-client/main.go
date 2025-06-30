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
	"time"

	// <-- 修改 #1: 更新为本地模块的导入路径
	"icpt-cli-client/internal/config"
)

// 用于解析上传响应的结构体
type UploadResponse struct {
	Data struct {
		ImageID uint   `json:"imageId"`
		Status  string `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}

// 用于解析状态查询响应的结构体
type StatusResponse struct {
	Data struct {
		ID               uint   `json:"id"`
		Status           string `json:"status"`
		ThumbnailPath    string `json:"thumbnail_path"`
		ThumbnailURL     string `json:"thumbnail_url"`
		ErrorInfo        string `json:"error_info"`
		OriginalFilename string `json:"original_filename"`
		CreatedAt        string `json:"created_at"`
	} `json:"data"`
}

const (
	maxRetries    = 15              // 最多轮询 15 次
	retryInterval = 2 * time.Second // 每次轮询间隔 2 秒
)

func main() {
	// 加载配置
	config.LoadConfig("config.yaml")

	if len(os.Args) < 2 {
		log.Fatal("使用方法: go run . <要上传的图片路径>")
	}
	filePath := os.Args[1]
	serverBaseURL := config.Cfg.Server.PublicHost

	// ---- 1. 上传文件 ----
	uploadResp, err := uploadFile(filePath, serverBaseURL)
	if err != nil {
		log.Fatalf("错误: 文件上传失败: %v", err)
	}
	log.Printf("文件已接收，服务器消息: %s, 图片ID: %d", uploadResp.Message, uploadResp.Data.ImageID)
	log.Println("开始查询处理状态...")

	// ---- 2. 轮询状态 ----
	for i := 0; i < maxRetries; i++ {
		time.Sleep(retryInterval) // 等待一下再查询

		statusResp, err := queryStatus(uploadResp.Data.ImageID, serverBaseURL)
		if err != nil {
			log.Printf("警告: 第 %d 次查询状态失败: %v", i+1, err)
			continue
		}

		switch statusResp.Data.Status {
		case "processing":
			fmt.Printf(".") // 打印一个点表示仍在处理中
		case "completed":
			log.Printf("\n成功! 图像处理完成。")
			// 使用服务器返回的完整URL
			if statusResp.Data.ThumbnailURL != "" {
				log.Printf("缩略图访问地址: %s", statusResp.Data.ThumbnailURL)
			} else {
				log.Printf("缩略图路径: %s", statusResp.Data.ThumbnailPath)
			}
			return // 任务完成，退出程序
		case "failed":
			log.Printf("\n失败! 图像处理失败。")
			log.Printf("错误信息: %s", statusResp.Data.ErrorInfo)
			return // 任务失败，退出程序
		default:
			log.Printf("\n收到未知状态: %s", statusResp.Data.Status)
		}
	}

	log.Printf("\n错误: 在 %d 次尝试后仍未获得最终结果，请稍后或检查服务器日志。", maxRetries)
}

// uploadFile 负责上传文件并解析初始响应
func uploadFile(filePath string, serverBaseURL string) (*UploadResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件 '%s': %w", filePath, err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("创建 form-data 失败: %w", err)
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("拷贝文件数据失败: %w", err)
	}
	writer.Close()

	uploadURL := serverBaseURL + "/api/v1/upload"
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return nil, fmt.Errorf("创建 HTTP 请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求到 '%s' 失败: %w", uploadURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("服务器返回非 202 状态码: %s", resp.Status)
	}

	var uploadResp UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, fmt.Errorf("无法解析上传响应的JSON: %w", err)
	}
	return &uploadResp, nil
}

// queryStatus 负责查询单个图片的状态
func queryStatus(imageID uint, serverBaseURL string) (*StatusResponse, error) {
	queryURL := fmt.Sprintf("%s/api/v1/images/%d", serverBaseURL, imageID)
	resp, err := http.Get(queryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("服务器返回非 200 状态码: %s", resp.Status)
	}

	var statusResp StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("无法解析状态响应的JSON: %w", err)
	}
	return &statusResp, nil
}
