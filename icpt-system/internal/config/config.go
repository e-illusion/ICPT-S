package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 结构体严格对应 yaml 文件的结构
type Config struct {
	Server struct {
		Port       string `yaml:"port"`
		PublicHost string `yaml:"public_host"` // <-- 新增此行
		HTTPS      struct {
			Enabled  bool   `yaml:"enabled"`   // 是否启用HTTPS
			CertFile string `yaml:"cert_file"` // 证书文件路径
			KeyFile  string `yaml:"key_file"`  // 私钥文件路径
			Port     string `yaml:"port"`      // HTTPS端口
		} `yaml:"https"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	JWT struct {
		SecretKey   string `yaml:"secret_key"`
		ExpireHours int    `yaml:"expire_hours"`
	} `yaml:"jwt"`
	Performance struct {
		WorkerCount          int  `yaml:"worker_count"`           // Worker进程数量
		MaxRequestSize       int  `yaml:"max_request_size"`       // 最大请求大小（MB）
		EnableGzip           bool `yaml:"enable_gzip"`            // 启用Gzip压缩
		EnableFileCache      bool `yaml:"enable_file_cache"`      // 启用文件缓存
		EnableConcurrency    bool `yaml:"enable_concurrency"`     // 启用并发处理
		MaxConcurrentUploads int  `yaml:"max_concurrent_uploads"` // 最大并发上传数
	} `yaml:"performance"`
}

var Cfg *Config

// LoadConfig 函数负责从指定的路径加载配置文件并解析
func LoadConfig(configPath string) {
	// 读取 yaml 文件内容
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("错误: 无法读取配置文件 %s: %v", configPath, err)
	}

	// 创建一个 Config 类型的变量
	var config Config

	// 将 yaml 内容解析到 config 变量中
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("错误: 解析配置文件失败: %v", err)
	}

	// 将解析后的配置赋值给全局变量 Cfg，方便项目其他地方使用
	Cfg = &config
	log.Println("配置文件加载成功！")
}
