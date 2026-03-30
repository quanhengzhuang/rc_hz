package main

import (
	"fmt"
	"os"

	"rc_hz/handler"
	"rc_hz/queue"
	"rc_hz/router"
	"rc_hz/worker"
)

func main() {
	// 加载配置
	config := LoadConfig()

	// 创建 MySQL 队列
	q, err := queue.NewMySQLQueue(config.DSN)
	if err != nil {
		fmt.Printf("Error creating MySQL queue: %v\n", err)
		os.Exit(1)
	}

	// 初始化处理器
	handlers := map[string]handler.Handler{
		"user_registered": &handler.UserRegisteredHandler{},
		"user_subscribed": &handler.UserSubscribedHandler{},
		"user_purchased":  &handler.UserPurchasedHandler{},
	}

	// 解析命令行参数
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./rc_hz [router|worker]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "router":
		// 启动路由器
		r := router.NewRouter(q)
		fmt.Printf("Starting router on port %s...\n", config.HTTPPort)
		if err := r.Start(config.HTTPPort); err != nil {
			fmt.Printf("Error starting router: %v\n", err)
			os.Exit(1)
		}
	case "worker":
		// 启动工作线程
		fmt.Printf("Starting %d workers...\n", config.WorkerCount)
		for i := 0; i < config.WorkerCount; i++ {
			go func(id int) {
				fmt.Printf("Worker %d started\n", id)
				w := worker.NewWorker(q, handlers)
				w.Start()
			}(i)
		}
		// 阻塞主线程
		select {}
	default:
		fmt.Println("Usage: ./rc_hz [router|worker]")
		os.Exit(1)
	}
}
