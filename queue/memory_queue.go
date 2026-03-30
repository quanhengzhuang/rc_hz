package queue

import (
	"sync"
	"time"
)

// MemoryQueue 内存队列实现
type MemoryQueue struct {
	messages []Message
	mutex    sync.Mutex
}

// NewMemoryQueue 创建内存队列实例
func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{
		messages: make([]Message, 0),
	}
}

// Produce 生产消息
func (q *MemoryQueue) Produce(message Message) (string, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.messages = append(q.messages, message)
	return message.ID, nil
}

// Consume 消费消息
func (q *MemoryQueue) Consume() (Message, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	now := time.Now()
	for i, msg := range q.messages {
		if (msg.Status == 0 || msg.Status == 2) && !msg.NextRetryAt.After(now) {
			// 找到可处理的消息
			message := msg
			// 从队列中移除
			q.messages = append(q.messages[:i], q.messages[i+1:]...)
			// 更新状态为处理中
			message.Status = 1
			return message, nil
		}
	}

	// 没有可处理的消息
	return Message{}, nil
}

// UpdateMessageStatus 更新消息状态
func (q *MemoryQueue) UpdateMessageStatus(id string, status int8, retryCount int, nextRetryAt time.Time) error {
	// 内存队列不需要持久化状态，直接返回
	return nil
}
