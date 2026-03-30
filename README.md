# rc_hz

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

## 思考
类似要实现一个「事件消息 -> 单一订阅者」的消息队列，可分为 producer、consumer、broker 三个视角来考虑。

producer:
- 企业内部业务系统的事件，如：用户注册成功、订阅付款成功、购买商品等
- producer 发消息给 MessageRouter，得到成功回复即可
- 针对 producer 的接口协议可用固定格式

consumer:
- consumer 端为各个外部系统，如广告系统、CRM 系统等
- consumer 端的接口格式不同，所以最好有一层适配层，用于将 producer 端的消息转为其请求参数，并识别 producer 端的返回是否成功

broker 要做的事:
- 为 producer 端提供一个 https 接口
- 适配层（可配置），用于将 producer 端的消息转为 consumer 端需要的格式，并处理返回值
- 可靠存储 producer 过来的消息，可使用消息队列或数据库
- 可靠投递给 consumer 端的消息，此处需要考虑失败重试
- 可支持 producer 端的事务消息（AI 提示）
- 增加唯一标识，用于给 consumer 端做幂等判断（AI 提示）

以上是必要的设计。未来可考虑对 consumer 的限流，批量处理，长连接等。

## 详细设计
### 服务
./message-router

1. 提供一个对外接口，用于生产消息：
- 接口名：POST /message
- 请求参数：消息类型、消息体
- 返回：消息 ID (用于状态查询)

2. 批量起消费 Worker：
- 用于从队列读取消息，并调用第三方接口。

### 消息队列
消息队列是个代码级模块，不是服务。是个抽象，底层可以用多种实现，如 MySQL，RocketMQ 等。

// 队列
type Queue interface {
    // 生产消息
    Produce(message Message) error
    // 消费消息
    Consume() (message Message, err error)
    // 更新消息状态
    UpdateMessageStatus(id string, status int8, retryCount int, nextRetryAt time.Time) error
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

## 快速开始

### 构建

```bash
go build -o rc_hz
```

### 运行

```bash
./rc_hz
```

服务将同时启动 HTTP 服务器和 Worker 线程。

## 配置

通过环境变量配置：

- `DSN`: MySQL 数据库连接字符串（可选，默认使用内存队列）
- `WORKER_COUNT`: Worker 线程数量（默认：5）
- `HTTP_PORT`: HTTP 服务端口（默认：8080）

## API 文档

### 生产消息

**接口地址：** `POST /message`

**请求头：**
```
Content-Type: application/json
```

**请求体：**
```json
{
  "type": "user_registered",
  "body": "{\"user_id\":123,\"source\":\"ad\",\"email\":\"test@example.com\"}"
}
```

**参数说明：**
- `type` (必填): 消息类型，支持以下类型：
  - `user_registered`: 用户注册
  - `user_subscribed`: 用户订阅
  - `user_purchased`: 用户购买
- `body` (必填): 消息体，JSON 字符串格式

**响应：**
```json
{
  "message_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**参数说明：**
- `message_id`: 消息 ID，用于追踪消息状态

**示例：**

```bash
curl -X POST http://localhost:8080/message \
  -H "Content-Type: application/json" \
  -d '{
    "type": "user_registered",
    "body": "{\"user_id\":123,\"source\":\"ad\",\"email\":\"test@example.com\"}"
  }'
```
