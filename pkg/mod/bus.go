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

func (b *Bus) CreateMsg(msg Message) {
	// Create a message and send it to all registered channels for the given name
	b.createMsg(msg)
}

func (b *Bus) createMsg(msg Message) bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	flag := false
	message := NewMessage(msg.Name, msg.Msg)

	fmt.Println("CreateMsg: ", message, "createMsg")

	// Send the message to all registered channels for the given name
	if channels, exists := b.msgList[msg.Name]; exists {
		for _, ch := range channels {
			select {
			case ch <- message:
				flag = true
			default:
				fmt.Println("Warning: Channel for message", msg.Name, "is full. Message dropped.")
			}
		}
	}
	return flag
}

// Subscribe 注册一个 channel，用于接收指定 name（target）的消息
func (b *Bus) Subscribe(name string, ch chan *Message) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.msgList[name] = append(b.msgList[name], ch)
}

// Unsubscribe 取消订阅指定 name（target）的消息 channel
func (b *Bus) Unsubscribe(name string, ch chan *Message) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if channels, exists := b.msgList[name]; exists {
		for i, c := range channels {
			if c == ch {
				b.msgList[name] = append(channels[:i], channels[i+1:]...)
				break
			}
		}
	}
}
