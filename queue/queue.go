package queue

// Queue 队列接口
type Queue interface {
	// 生产消息
	Produce(message Message) error
	// 消费消息
	Consume() (message Message, err error)
}
