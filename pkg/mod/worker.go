package mod

import (
	"context"
	"fmt"
)

type Worker struct {
	ctx    context.Context
	cancel context.CancelFunc
	msg    *Message
	bus    *Bus
	ch     chan *Message
}

func NewWorker(bus *Bus, ctx context.Context, cancel context.CancelFunc) *Worker {
	return &Worker{
		bus:    bus,
		ctx:    ctx,
		cancel: cancel,
		// 10 个消息缓冲区
		ch: make(chan *Message, 10),
	}
}

func (w *Worker) GetContext() context.Context {
	return w.ctx
}

func (w *Worker) Initialize(c context.Context) error {
	// 订阅消息
	w.bus.Subscribe("Worker", w.ch)
	fmt.Println("Worker initialized with context:", w.GetContext())

	go w.run()

	return nil
}

func (w *Worker) DeInitialize() error {
	// 取消订阅消息
	w.bus.Unsubscribe("Worker", w.ch)
	fmt.Println("Worker deinitialized")
	return nil
}

func (w *Worker) run() {
	for {
		select {
		case msg := <-w.ch:
			fmt.Println("🐮🐴收到:", msg.Time.Format("2007-01-02 15:04:05"), "Name:", msg.Name, "Msg:", msg.Msg)
			var newMsg Message
			newMsg.Name = "System"
			newMsg.Msg = "嘿嘿 Worker 收到消息了"
			w.bus.CreateMsg(newMsg)

		case <-w.ctx.Done():
			fmt.Println("Worker context canceled")
			return
		}
	}
}
