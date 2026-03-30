package queue

import (
	"errors"
	"sync"
	"time"
)

// MySQLQueue MySQL 队列实现（mock 版本）
type MySQLQueue struct {
	messages map[string]Message
	mutex    sync.Mutex
}

// NewMySQLQueue 创建 MySQL 队列实例（mock 版本）
func NewMySQLQueue(dsn string) (*MySQLQueue, error) {
	// 不需要真正连接数据库，直接返回 mock 实例
	return &MySQLQueue{
		messages: make(map[string]Message),
	}, nil
}

// Produce 生产消息
func (q *MySQLQueue) Produce(message Message) (string, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// 存储消息到内存
	q.messages[message.ID] = message
	return message.ID, nil
}

// Consume 消费消息
func (q *MySQLQueue) Consume() (Message, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	var selectedMessage Message
	var selectedID string
	now := time.Now()

	// 查找待处理或需要重试的消息
	for id, msg := range q.messages {
		if (msg.Status == 0 || msg.Status == 2) && msg.NextRetryAt.Before(now) {
			// 找到最早需要处理的消息
			if selectedID == "" || msg.NextRetryAt.Before(selectedMessage.NextRetryAt) {
				selectedMessage = msg
				selectedID = id
			}
		}
	}

	if selectedID == "" {
		return Message{}, nil // 没有消息
	}

	// 更新消息状态为处理中
	selectedMessage.Status = 1
	q.messages[selectedID] = selectedMessage

	return selectedMessage, nil
}

// UpdateMessageStatus 更新消息状态
func (q *MySQLQueue) UpdateMessageStatus(id string, status int8, retryCount int, nextRetryAt time.Time) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// 检查消息是否存在
	if msg, exists := q.messages[id]; exists {
		// 更新消息状态
		msg.Status = status
		msg.RetryCount = retryCount
		msg.NextRetryAt = nextRetryAt
		q.messages[id] = msg
		return nil
	}

	return errors.New("message not found")
}
