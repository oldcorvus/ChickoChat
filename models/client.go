package data

import (
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

// struct representing a user in a ChatRoom
type Client struct {
	User         UserData
	LastActivity time.Time `json:"last_activity"`
	// The websocket Connection.
	Conn *websocket.Conn `json:"-"`
	// Buffered channel of outbound messages.
	Send chan *ChatEvent `json:"-"`
	// Broker for connection
	Broker *Broker
}

type UserData struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email  string             `json:"email" bson:"email" binding:"required"`
	Name   string             `json:"name" bson:"name"`
	Active bool               `json:"active" bson:"active"`
}

func NewClient(conn *websocket.Conn, user *UserData, broker *Broker) *Client {

	client := &Client{
		User:   *user,
		Conn:   conn,
		Send:   make(chan *ChatEvent),
		Broker: broker,
	}
	return client
}
