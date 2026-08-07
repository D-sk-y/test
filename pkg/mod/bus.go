package mod

import (
	"fmt"
	"sync"
)

type Bus struct {
	// 消息队列
	msgList map[string][]chan *Message
	mutex   sync.RWMutex
}

var (
	defaultBus *Bus
	busOnce    sync.Once
)

// GetBus returns the singleton instance of the Bus. It initializes the Bus if it hasn't been created yet.
func GetBus() *Bus {
	busOnce.Do(func() {
		defaultBus = &Bus{
			msgList: make(map[string][]chan *Message),
		}
	})
	return defaultBus
}

func (b *Bus) CreateMsg(name string, msg any) {
	// Create a message and send it to all registered channels for the given name
	b.createMsg(name, msg)
}

func (b *Bus) createMsg(name string, msg any) bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	flag := false
	message := NewMessage(name, msg)

	// Send the message to all registered channels for the given name
	if channels, exists := b.msgList[name]; exists {
		for _, ch := range channels {
			select {
			case ch <- message:
				flag = true
			default:
				fmt.Println("Warning: Channel for message", name, "is full. Message dropped.")
			}
		}
	}
	return flag
}
