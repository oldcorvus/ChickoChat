package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// struct representing a chat room
type ChatRoom struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Clients map[string]*Client `json:"-"`
}
