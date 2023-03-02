package data

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gorilla/websocket"
)

// struct representing a user in a ChatRoom
type Client struct {
	User         UserData
	LastActivity time.Time `json:"last_activity"`
	// The websocket Connection.
	Conn *websocket.Conn `json:"-"`
	// Buffered channel of outbound messages.
	Send chan []byte `json:"-"`
}

type UserData struct {
	ID     primitive.ObjectID  `json:"_id" bson:"_id"`
	Email  string `json:"email" bson:"email"`
	Name   string `json:"name" bson:"name"`
	Active bool   `json:"active" bson:"active"`
}
