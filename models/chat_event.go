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
	EventType string             `json:"event_type,omitempty"`
	UserID    primitive.ObjectID `json:"name,omitempty"`
	RoomID    primitive.ObjectID `json:"room_id,omitempty"`
	Message   string             `json:"msg,omitempty"`
	Timestamp time.Time          `json:"time,omitempty"`
}
