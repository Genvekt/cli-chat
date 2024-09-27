package model

import "time"

type Message struct {
	SenderID  int64
	ChatID    int64
	Text      string
	Timestamp time.Time
}
