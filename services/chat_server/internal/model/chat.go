package model

import "time"

// Chat service model
type Chat struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	Members   []*ChatMember `json:"members"`
	CreatedAt time.Time     `json:"created_at"`
}
