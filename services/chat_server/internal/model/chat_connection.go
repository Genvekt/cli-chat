package model

import (
	"sync"
)

// ChatConnection represents active chat with user connections and inner buffer
type ChatConnection struct {
	Mx          sync.RWMutex
	ChatID      int64
	Buffer      chan *Message
	Connections map[string]chan *Message
}
