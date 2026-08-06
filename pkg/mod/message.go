package mod

import "time"

// Message represents a message with a name, content, and timestamp.
type Message struct {
	Name string    `json:"name"`
	Msg  any       `json:"msg"`
	Time time.Time `json:"time"`
}

func NewMessage(name string, msg any) *Message {
	return &Message{
		Name: name,
		Msg:  msg,
		Time: time.Now(),
	}
}
