package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// Subscribe is used to broadcast a message indicating user has joined ChatRoom
	Subscribe = "join"
	// Broadcast is used to broadcast messages to all subscribed users
	Broadcast = "send"
	// Unsubscribe is used to broadcast a message indicating user has left ChatRoom
	Unsubscribe = "leave"
)

// struct representing a message event in an  ChatRoom
type ChatEvent struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	EventType string             `json:"type" bson:"type ,omitempty"`
	UserID    primitive.ObjectID `json:"user_id,omitempty"`
	RoomID    primitive.ObjectID `json:"room_id,omitempty"`
	Message   string             `json:"message,omitempty"`
	Timestamp time.Time          `json:"time"`
}
