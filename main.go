package main

import (
	"fmt"

	"rc_hz/handler"
	"rc_hz/queue"
	"rc_hz/router"
	"rc_hz/worker"
)

func main() {
	// 加载配置
	config := LoadConfig()

	// 创建队列（优先使用 MySQL，失败则使用内存队列）
	var q queue.Queue
	q, err := queue.NewMySQLQueue(config.DSN)
	if err != nil {
		fmt.Printf("Warning: Failed to create MySQL queue: %v\n", err)
		fmt.Println("Using memory queue instead")
		q = queue.NewMemoryQueue()
	}

	// 初始化处理器
	handlers := map[string]handler.Handler{
		"user_registered": &handler.UserRegisteredHandler{},
		"user_subscribed": &handler.UserSubscribedHandler{},
		"user_purchased":  &handler.UserPurchasedHandler{},
	}

	// 同时启动路由器和工作线程
	worker := worker.NewWorker(q, handlers)
	go worker.Start(config.WorkerCount)
	r := router.NewRouter(q)
	r.Start(config.HTTPPort)
}
