package mod

import (
	"context"
	"fmt"
)

type System struct {
	bus    *Bus
	ctx    context.Context
	cancel context.CancelFunc
	// 这里可以添加系统的字段，例如配置、状态等
}

func NewSystem(bus *Bus, ctx context.Context, cancel context.CancelFunc) *System {
	return &System{
		bus:    bus,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *System) GetContext() context.Context {
	return s.ctx
}

func (s *System) Initialize(c context.Context) error {
	fmt.Println("System initialized with context:", s.GetContext())
	return nil
}

func (s *System) DeInitialize() error {
	fmt.Println("System deinitialized")
	return nil
}
