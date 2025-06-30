package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"icpt-system/internal/api" // <-- 导入 api 包
	"icpt-system/internal/config"
	"icpt-system/internal/middleware"
	"icpt-system/internal/store"
	"icpt-system/internal/websocket"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	// 确保上传目录存在
	os.MkdirAll("uploads/originals", os.ModePerm)
	os.MkdirAll("uploads/thumbnails", os.ModePerm)

	// 1. 加载配置
	config.LoadConfig("config.yaml")

	// 2. 初始化数据库和Redis连接
	store.InitDB()
	store.InitRedis()

	// 3. 初始化WebSocket Hub
	websocket.InitHub()

	// 4. 初始化 Gin 引擎
	r := gin.Default()

	// 为 multipart forms 设置一个较低的内存限制 (默认是 32 MiB)
	// 这意味着大于 8 MiB 的文件会临时存储在磁盘上，而不是完全加载到内存中。
	maxMemory := int64(config.Cfg.Performance.MaxRequestSize) << 20 // 配置的最大请求大小（MB）
	r.MaxMultipartMemory = maxMemory

	// CORS中间件配置 - 允许跨域请求（动态判断源）
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// 开发环境：允许localhost、127.0.0.1、::1的任何端口
			// 生产环境：允许指定的域名和端口
			allowedPatterns := []string{
				"http://localhost:",
				"https://localhost:",
				"http://127.0.0.1:",
				"https://127.0.0.1:",
				"http://[::1]:",
				"https://[::1]:",
				"http://114.55.58.3:",
				"https://114.55.58.3:",
				"http://0.0.0.0:",
				"http://114.55.58.3:3000", // 新增：允许特定的前端开发服务器
			}

			for _, pattern := range allowedPatterns {
				if strings.HasPrefix(origin, pattern) {
					return true
				}
			}

			// 如果是空 origin（如直接访问），也允许
			return origin == ""
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Length", "Content-Type", "Authorization",
			"Accept", "X-Requested-With", "Cache-Control",
		},
		ExposeHeaders: []string{
			"Content-Length", "Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	log.Println("✅ 启用CORS跨域支持（动态源判断）")

	// 性能优化中间件
	if config.Cfg.Performance.EnableGzip {
		r.Use(gzip.Gzip(gzip.DefaultCompression))
		log.Println("✅ 启用Gzip压缩")
	}

	// 添加限流中间件
	if config.Cfg.Performance.EnableConcurrency {
		r.Use(middleware.ConcurrencyLimitMiddleware(config.Cfg.Performance.MaxConcurrentUploads))
		log.Printf("✅ 启用并发限制: %d", config.Cfg.Performance.MaxConcurrentUploads)
	}

	// 新增：配置静态文件服务
	// 将 URL /static/ 映射到本地的 ./uploads 目录
	r.Static("/static", "./uploads") // 静态文件服务

	// 配置Web前端静态文件服务
	r.Static("/web", "./web")
	r.StaticFile("/", "./web/index.html") // 首页重定向

	// 5. 定义 API 路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API路由组
	v1 := r.Group("/api/v1")
	{
		// 公开接口（无需认证）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", api.RegisterHandler)
			auth.POST("/login", api.LoginHandler)
		}

		// 需要认证的接口
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// 用户相关
			protected.GET("/profile", api.GetProfileHandler)

			// 图像上传和管理
			protected.POST("/upload", api.UploadImageHandler)
			protected.GET("/images/:id", api.GetImageStatusHandler)
			protected.GET("/images", api.GetUserImagesHandler)
			protected.DELETE("/images/:id", api.DeleteImageHandler)
			protected.POST("/images/batch-delete", api.BatchDeleteImagesHandler)

			// 统计信息相关
			protected.GET("/stats/dashboard", api.GetDashboardStats(store.DB))
			protected.GET("/activity/recent", api.GetRecentActivity(store.DB))
			protected.GET("/stats/status-count", api.GetImageStatusCount(store.DB))

			// WebSocket相关
			protected.GET("/ws", api.WebSocketHandler)
			protected.GET("/ws/stats", api.WebSocketStatsHandler)
			protected.POST("/notify/:userID", api.NotifyTestHandler) // 测试用
		}

		// 测试接口（保留用于性能测试）
		v1.POST("/upload-sync", api.UploadImageSyncHandlerForTest)
	}

	// 6. 启动服务器（支持HTTP和HTTPS）
	if config.Cfg.Server.HTTPS.Enabled {
		// 启动HTTPS服务器
		startHTTPSServer(r)
	} else {
		// 启动HTTP服务器
		startHTTPServer(r)
	}
}

// startHTTPServer 启动HTTP服务器
func startHTTPServer(r *gin.Engine) {
	listenAddr := "0.0.0.0" + config.Cfg.Server.Port // 变为 "0.0.0.0:8080"
	log.Printf("HTTP API 服务器启动中，监听地址: %s (对外开放)", listenAddr)
	if err := r.Run(listenAddr); err != nil {
		log.Fatalf("错误: HTTP服务器启动失败: %v", err)
	}
}

// startHTTPSServer 启动HTTPS服务器
func startHTTPSServer(r *gin.Engine) {
	// 确保证书目录存在
	os.MkdirAll("certs", os.ModePerm)

	certFile := config.Cfg.Server.HTTPS.CertFile
	keyFile := config.Cfg.Server.HTTPS.KeyFile

	// 检查证书文件是否存在，如果不存在则生成自签名证书
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Println("SSL证书不存在，正在生成自签名证书...")
		if err := generateSelfSignedCert(certFile, keyFile); err != nil {
			log.Fatalf("生成自签名证书失败: %v", err)
		}
	}

	// 配置TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12, // 最低TLS 1.2
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	server := &http.Server{
		Addr:      "0.0.0.0" + config.Cfg.Server.HTTPS.Port,
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	log.Printf("HTTPS API 服务器启动中，监听地址: %s (对外开放)", server.Addr)
	log.Printf("证书文件: %s", certFile)
	log.Printf("私钥文件: %s", keyFile)

	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Fatalf("错误: HTTPS服务器启动失败: %v", err)
	}
}

// generateSelfSignedCert 生成自签名SSL证书
func generateSelfSignedCert(certFile, keyFile string) error {
	// 生成私钥
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// 创建证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"ICPT System"},
			Country:       []string{"CN"},
			Province:      []string{""},
			Locality:      []string{"Beijing"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour), // 1年有效期
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{
			net.IPv4(127, 0, 0, 1), 
			net.IPv6loopback,
			net.IPv4(114, 55, 58, 3), // 添加外网IP地址
		},
		DNSNames: []string{"localhost", "*.localhost", "114.55.58.3"},
	}

	// 创建证书
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// 保存证书文件
	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return err
	}

	// 保存私钥文件
	keyOut, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	defer keyOut.Close()

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return err
	}

	log.Printf("自签名证书生成成功: %s, %s", certFile, keyFile)
	return nil
}
