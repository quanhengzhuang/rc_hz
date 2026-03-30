package worker

import (
	"context"
	"fmt"
	"time"

	"rc_hz/handler"
	"rc_hz/queue"
)

// Worker 工作线程
type Worker struct {
	queue    queue.Queue
	handlers map[string]handler.Handler
}

// NewWorker 创建工作线程实例
func NewWorker(q queue.Queue, handlers map[string]handler.Handler) *Worker {
	return &Worker{
		queue:    q,
		handlers: handlers,
	}
}

// Start 启动工作线程
func (w *Worker) Start(workerCount int) {
	// 启动指定数量的工作线程
	for i := 0; i < workerCount; i++ {
		go func(id int) {
			fmt.Printf("Worker %d started\n", id)
			for {
				// 消费消息
				message, err := w.queue.Consume()
				if err != nil {
					fmt.Printf("Error consuming message: %v\n", err)
					time.Sleep(1 * time.Second)
					continue
				}

				// 检查是否有消息
				if message.ID == "" {
					time.Sleep(1 * time.Second)
					continue
				}

				// 处理消息
				w.processMessage(message)
			}
		}(i)
	}

	// 阻塞主线程
	select {}
}

// processMessage 处理消息
func (w *Worker) processMessage(message queue.Message) {
	// 获取对应的处理器
	h, ok := w.handlers[message.Type]
	if !ok {
		fmt.Printf("No handler found for message type: %s\n", message.Type)
		// 更新消息状态为处理失败
		w.updateMessageStatus(message, 2)
		return
	}

	// 处理消息
	err := h.Handle(context.Background(), message.Body)
	if err != nil {
		fmt.Printf("Error processing message %s: %v\n", message.ID, err)
		// 更新消息状态为处理失败
		w.updateMessageStatus(message, 2)
		return
	}

	// 更新消息状态为已处理
	w.updateMessageStatus(message, 1)
}

// updateMessageStatus 更新消息状态
func (w *Worker) updateMessageStatus(message queue.Message, status int8) {
	var nextRetryAt time.Time
	retryCount := message.RetryCount

	if status == 2 {
		// 处理失败，增加重试次数
		retryCount++
		// 计算下一次重试时间（指数退避）
		delay := time.Duration(1<<uint(retryCount)) * time.Second
		nextRetryAt = time.Now().Add(delay)
	} else {
		nextRetryAt = time.Now()
	}

	// 调用队列的 UpdateMessageStatus 方法
	w.queue.UpdateMessageStatus(message.ID, status, retryCount, nextRetryAt)
}
