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
}

func NewWorker(bus *Bus, ctx context.Context, cancel context.CancelFunc) *Worker {
	return &Worker{
		bus:    bus,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (w *Worker) GetContext() context.Context {
	return w.ctx
}

func (w *Worker) Initialize(c context.Context) error {
	fmt.Println("Worker initialized with context:", w.GetContext())
	return nil
}

func (w *Worker) DeInitialize() error {
	fmt.Println("Worker deinitialized")
	return nil
}
