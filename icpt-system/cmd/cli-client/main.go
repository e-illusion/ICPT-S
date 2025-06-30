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
	"icpt-system/internal/config"
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
		ID            uint   `json:"ID"`
		Status        string `json:"Status"`
		ThumbnailPath string `json:"ThumbnailPath"`
		ErrorInfo     string `json:"ErrorInfo"`
	} `json:"data"`
}

const (
	uploadEndpoint = "/api/v1/upload"
	statusEndpoint = "/api/v1/images/"
	maxRetries     = 10 // 最多轮询 10 次
	retryInterval  = 2 * time.Second // 每次轮询间隔 2 秒
)

func main() {
	config.LoadConfig("config.yaml")

	if len(os.Args) < 2 {
		log.Fatal("请提供要上传的图片文件路径作为参数")
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
			// 构造可公开访问的 URL
			fullURL := serverBaseURL + "/static/" + statusResp.Data.ThumbnailPath
			log.Printf("缩略图访问地址: %s", fullURL)
			return // 任务完成，退出程序
		case "failed":
			log.Printf("\n失败! 图像处理失败。")
			log.Printf("错误信息: %s", statusResp.Data.ErrorInfo)
			return // 任务失败，退出程序
		default:
			log.Printf("\n收到未知状态: %s", statusResp.Data.Status)
		}
	}

	log.Printf("\n错误: 在 %d 次尝试后仍未获得最终结果，请稍后重试。", maxRetries)
}

// uploadFile 负责上传文件并解析初始响应
func uploadFile(filePath string, serverBaseURL string) (*UploadResponse, error) {
	// ... (这部分逻辑和之前基本一样，但返回类型变了)
	// (为了简洁，这里省略了 multipart/form-data 的构造代码，因为它和之前的版本一样)
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

    req, _ := http.NewRequest("POST", serverBaseURL+"/api/v1/upload", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var uploadResp UploadResponse
    if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
        return nil, fmt.Errorf("无法解析上传响应: %v", err)
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
		return nil, fmt.Errorf("无法解析状态响应: %v", err)
	}
	return &statusResp, nil
}
