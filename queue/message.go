package queue

import "time"

// Message 消息结构体
type Message struct {
	// 消息 ID，使用分布式的唯一标识，不要依赖数据库主键
	ID string `json:"id"`
	// 消息类型，用于区分不同业务消息
	Type string `json:"type"`
	// 消息体，用 Json 格式存储
	Body string `json:"body"`
	// 状态，0 表示待处理，1 表示已处理，2 表示处理失败
	Status int8 `json:"status"`
	// 创建时间
	CreateAt time.Time `json:"create_at"`
	// 消息的重试次数，默认 0，每次重试 +1
	RetryCount int `json:"retry_count"`
	// 下一次重试时间，默认当前时间，重试后为：当前时间 + RetryCount 指数倍
	NextRetryAt time.Time `json:"next_retry_at"`
}
