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

	// 创建队列（使用 MySQL 队列）
	q, err := queue.NewMySQLQueue(config.DSN)
	if err != nil {
		fmt.Printf("Error creating MySQL queue: %v\n", err)
		return
	}

	// 业务处理配置
	handlers := map[string]handler.Handler{
		"user_registered": &handler.UserRegisteredHandler{},
		"user_subscribed": &handler.UserSubscribedHandler{},
		"user_purchased":  &handler.UserPurchasedHandler{},
	}

	// Worker 进程启动
	worker := worker.NewWorker(q, handlers)
	go worker.Start(config.WorkerCount)


	// 接口进程启动
	r := router.NewRouter(q)
	r.Start(config.HTTPPort)
}
