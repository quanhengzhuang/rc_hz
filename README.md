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

## 思路
类似一个「事件消息 -> 单一订阅者」的消息队列，可以抽象出一个消息路由，叫 `MessageRouter`，可分为 producer、consumer、broker 三个视角来考虑。

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
- 可靠存储 producer 过来的消息，可使用消息队列
- 可靠投递给 consumer 端的消息，此处需要考虑失败重试
- 可支持 producer 端的事务消息（AI 提示）
- 增加唯一标识，用于给 consumer 端做幂等判断（AI 提示）

以上是必要的设计。另外可额外考虑对 consumer 的限流，批量处理，长连接等，这些不是必要的设计，后面可考虑。

设计思路一般都会对比多家，自己的是一种，不同 AI 可提供多种。如果自己不预想思路，直接找 AI，后面就很难判断哪里的是过度设计，或者哪里漏了。

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
}