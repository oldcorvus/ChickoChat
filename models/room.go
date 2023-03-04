package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// struct representing a chat room
type ChatRoom struct {
	ID        primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Title     string               `json:"title" bson:"title"`
	CreatedAt time.Time            `json:"createdAt" bson:"created"`
	UpdatedAt time.Time            `json:"updatedAt" bson:"updated"`
	Clients   []primitive.ObjectID `json:"users" bson:"users"`
}
