package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config 结构体映射 config.yaml 文件中的内容
type Config struct {
	Server struct {
		PublicHost string `yaml:"public_host"`
	} `yaml:"server"`
}

// Cfg 是一个用于存储配置的全局变量
var Cfg Config

// LoadConfig 从指定的路径加载配置
func LoadConfig(configPath string) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("错误: 无法读取配置文件 '%s': %v", configPath, err)
	}

	err = yaml.Unmarshal(file, &Cfg)
	if err != nil {
		log.Fatalf("错误: 无法解析配置文件: %v", err)
	}
}
