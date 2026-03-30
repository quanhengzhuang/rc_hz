package handler

import "context"

// Handler 处理器接口
type Handler interface {
	Handle(ctx context.Context, body string) error
}
