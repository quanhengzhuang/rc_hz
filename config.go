package main

import (
	"os"
	"strconv"
)

// Config 配置结构体
type Config struct {
	// 数据库连接字符串
	DSN string
	// Worker 数量
	WorkerCount int
	// HTTP 服务端口
	HTTPPort string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 从环境变量加载配置
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/message_queue?charset=utf8mb4&parseTime=True&loc=Local"
	}

	workerCountStr := os.Getenv("WORKER_COUNT")
	workerCount := 5
	if workerCountStr != "" {
		if count, err := strconv.Atoi(workerCountStr); err == nil {
			workerCount = count
		}
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	return &Config{
		DSN:         dsn,
		WorkerCount: workerCount,
		HTTPPort:    httpPort,
	}
}
