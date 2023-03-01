package data

import (
	"time"
)

// struct representing a chat room
type ChatRoom struct {
	Title     string             `json:"title"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
	ID        int                `json:"id"`
	Clients   map[string]*Client `json:"-"`
}
