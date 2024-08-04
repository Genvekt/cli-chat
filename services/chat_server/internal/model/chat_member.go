package model

import "time"

// ChatMember service model
type ChatMember struct {
	ID       int64     `json:"id"`
	JoinedAt time.Time `json:"joined_at"`
}
