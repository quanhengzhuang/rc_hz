package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"rc_hz/queue"
)

// Router 路由器
type Router struct {
	queue queue.Queue
}

// NewRouter 创建路由器实例
func NewRouter(q queue.Queue) *Router {
	return &Router{queue: q}
}

// Start 启动 HTTP 服务
func (r *Router) Start(port string) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 生产消息接口
	router.POST("/message", r.produceMessage)

	return router.Run(":" + port)
}

// produceMessage 生产消息
func (r *Router) produceMessage(c *gin.Context) {
	var request struct {
		Type string `json:"type" binding:"required"`
		Body string `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建消息
	message := queue.Message{
		ID:          uuid.New().String(),
		Type:        request.Type,
		Body:        request.Body,
		Status:      0, // 待处理
		CreateAt:    time.Now(),
		RetryCount:  0,
		NextRetryAt: time.Now(),
	}

	// 生产消息
	if err := r.queue.Produce(message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message produced successfully"})
}
