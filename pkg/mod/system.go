package mod

import "context"

type System struct {
	ctx    context.Context
	cancel context.CancelFunc
	// 这里可以添加系统的字段，例如配置、状态等
}

func NewSystem(ctx context.Context, cancel context.CancelFunc) *System {
	return &System{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *System) GetContext() context.Context {
	return s.ctx
}
