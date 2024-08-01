package chat

import "time"

// Member repository model
type Member struct {
	UserID   int64     `json:"user_id"`
	ChatID   int64     `json:"chat_id"`
	JoinedAt time.Time `json:"joined_at"`
}
