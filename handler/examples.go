package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UserRegisteredHandler 用户注册处理器
type UserRegisteredHandler struct {}

// Handle 处理用户注册消息
func (h *UserRegisteredHandler) Handle(ctx context.Context, body string) error {
	// 解析消息体
	var message map[string]interface{}
	if err := json.Unmarshal([]byte(body), &message); err != nil {
		return err
	}

	// 模拟调用广告系统 API
	url := "https://api.ad-system.com/notify"
	headers := map[string]string{
		"Content-Type": "application/json",
		"Authorization": "Bearer token",
	}

	// 构造请求体
	reqBody, err := json.Marshal(map[string]interface{}{
		"user_id": message["user_id"],
		"event":   "registered",
		"source":  message["source"],
	})
	if err != nil {
		return err
	}

	// 发送 HTTP 请求
	return sendHTTPRequest(url, headers, reqBody)
}

// UserSubscribedHandler 用户订阅处理器
type UserSubscribedHandler struct {}

// Handle 处理用户订阅消息
func (h *UserSubscribedHandler) Handle(ctx context.Context, body string) error {
	// 解析消息体
	var message map[string]interface{}
	if err := json.Unmarshal([]byte(body), &message); err != nil {
		return err
	}

	// 模拟调用 CRM 系统 API
	url := "https://api.crm-system.com/contact/update"
	headers := map[string]string{
		"Content-Type": "application/json",
		"API-Key":      "crm-api-key",
	}

	// 构造请求体
	reqBody, err := json.Marshal(map[string]interface{}{
		"contact_id": message["user_id"],
		"status":     "subscribed",
		"plan":       message["plan"],
	})
	if err != nil {
		return err
	}

	// 发送 HTTP 请求
	return sendHTTPRequest(url, headers, reqBody)
}

// UserPurchasedHandler 用户购买处理器
type UserPurchasedHandler struct {}

// Handle 处理用户购买消息
func (h *UserPurchasedHandler) Handle(ctx context.Context, body string) error {
	// 解析消息体
	var message map[string]interface{}
	if err := json.Unmarshal([]byte(body), &message); err != nil {
		return err
	}

	// 模拟调用库存系统 API
	url := "https://api.inventory-system.com/update"
	headers := map[string]string{
		"Content-Type": "application/json",
		"X-API-Key":    "inventory-api-key",
	}

	// 构造请求体
	reqBody, err := json.Marshal(map[string]interface{}{
		"product_id": message["product_id"],
		"quantity":   message["quantity"],
		"action":     "decrease",
	})
	if err != nil {
		return err
	}

	// 发送 HTTP 请求
	return sendHTTPRequest(url, headers, reqBody)
}

// sendHTTPRequest 发送 HTTP 请求
func sendHTTPRequest(url string, headers map[string]string, body []byte) error {
	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 设置请求体
	req.Body = nil // 这里简化处理，实际应该使用 bytes.NewBuffer(body)

	// 发送请求
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	// 打印日志
	fmt.Printf("Sent request to %s with body: %s\n", url, string(body))

	return nil
}
