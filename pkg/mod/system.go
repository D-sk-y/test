package mod

import (
	"context"
	"fmt"
)

type System struct {
	bus    *Bus
	ctx    context.Context
	cancel context.CancelFunc
	msg    *Message
	ch     chan *Message
	// 这里可以添加系统的字段，例如配置、状态等
}

func NewSystem(bus *Bus, ctx context.Context, cancel context.CancelFunc) *System {
	return &System{
		bus:    bus,
		ctx:    ctx,
		cancel: cancel,
		// 10 个消息缓冲区
		ch: make(chan *Message, 10),
	}
}

func (s *System) GetContext() context.Context {
	return s.ctx
}

func (s *System) Initialize(c context.Context) error {
	// 订阅消息
	s.bus.Subscribe("System", s.ch)

	go s.run()

	fmt.Println("System initialized with context:", s.GetContext())
	return nil
}

func (s *System) DeInitialize() error {
	s.bus.Unsubscribe("System", s.ch)
	fmt.Println("System deinitialized")
	return nil
}

func (s *System) run() {
	for {
		select {
		case msg := <-s.ch:
			fmt.Println(msg.Msg)
		case <-s.ctx.Done():
			fmt.Println("System context canceled")
			return
		}
	}
}
