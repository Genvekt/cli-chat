package model

import "time"

// Message repository model
type Message struct {
	ID        int64
	SenderID  int64
	ChatID    int64
	Content   string
	Timestamp time.Time
}
