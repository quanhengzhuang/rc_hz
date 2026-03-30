AI 使用记录

## 设计思路对比
将需求分别提交给 deepseek 和 gemini，对比和自己的设计思路，查漏补缺。

prompt:
```
请理解以下需求，说一下实现思路，未来将会用 Go 语言实现。

## 需求
企业内部多个业务系统在关键事件发生时，需要调用外部系统供应商提供的 HTTP(S) API 进行通知。例如：
- 用户通过第三方广告系统引流并成功注册后，通知对应的广告系统
- 用户订阅付款成功后，通知 CRM 系统更改 Contact 状态
- 用户购买商品后，通知库存系统进行库存变更

不同供应商的 API：
- 请求地址不同
- Header / Body 格式不同

业务系统本身：
- 不需要关心外部 API 的返回值
- 只需确保通知请求能够被稳定、可靠地送达
```

## 通过设计让 AI 给出代码实现
参考 AI 给出的代码，调整设计文档。
prompt:
```
以上需求，用下面的设计，请实现具体的代码：
## 详细设计
### 消息生产服务
/bin/message-router

提供一个对外接口：
- 生产消息：POST /message

### 消费者服务
/bin/message-worker

用于从队列读取消息，并调用第三方接口。

### 消息队列
消息队列是个代码级模块，不是服务。是个抽象，底层可以用多种实现，如 MySQL，RocketMQ 等。

// 队列
type Queue interface {
    // 生产消息
    Produce(message Message) error
    // 消费消息
    Consume() (message Message, err error)
}

// 消息
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

### 业务适配层
// 根据 Message.Type 来路由到不同的 Handler 对象处理
type Handler struct {
    ...
}
```

## 调整设计后，再次由 AI 给出代码实现
prompt:
```
以下是我最新的设计，请给出实现代码，代码要可运行。

## 详细设计
### 消息生产服务
/bin/message-router

提供一个对外接口：
- 生产消息：POST /message

### 消费者服务
/bin/message-worker

用于从队列读取消息，并调用第三方接口。

### 消息队列
消息队列是个代码级模块，不是服务。是个抽象，底层可以用多种实现，如 MySQL，RocketMQ 等。

// 队列
type Queue interface {
    // 生产消息
    Produce(message Message) error
    // 消费消息
    Consume() (message Message, err error)
}

// 消息
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

### 业务适配层
// 根据 Message.Type 来路由到不同的 Handler 对象处理
type Handler struct {
    Handle(ctx context.Context, body string) error
}

// Handler 配置
HandlerConfig := map[string]Handler{
    "user_registered": UserRegisteredHandler{},
    "user_subscribed": UserSubscribedHandler{},
    "user_purchased": UserPurchasedHandler{},
}

### 代码目录
rc_hz/
├── go.mod
├── go.sum
├── main.go               # 入口，根据命令行参数启动 router 或 worker
├── config.go             # 配置（数据库连接、Worker 数量等）
├── queue/
│   ├── queue.go          # Queue 接口定义
│   ├── mysql_queue.go    # MySQL 实现
│   └── message.go        # Message 结构体
├── handler/
│   ├── handler.go        # Handler 接口
│   └── examples.go       # 示例 Handler（UserRegistered 等）
├── router/
│   └── router.go         # HTTP 路由及处理器
└── worker/
    └── worker.go         # Worker 逻辑
```