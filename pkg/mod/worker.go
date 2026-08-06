package mod

import "context"

type Worker struct {
	ctx    context.Context
	cancel context.CancelFunc
	msg    *Message
	bus    *Bus
}
