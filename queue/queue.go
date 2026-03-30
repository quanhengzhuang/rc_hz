package queue

import "time"

// Queue 队列接口
type Queue interface {
	// 生产消息
	Produce(message Message) (string, error)
	// 消费消息
	Consume() (message Message, err error)
	// 更新消息状态
	UpdateMessageStatus(id string, status int8, retryCount int, nextRetryAt time.Time) error
}
