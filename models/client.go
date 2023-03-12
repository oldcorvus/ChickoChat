package data

import (
	"log"
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
	Send chan []byte `json:"-"`
}

type UserData struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email  string             `json:"email" bson:"email" binding:"required"`
	Name   string             `json:"name" bson:"name"`
	Active bool               `json:"active" bson:"active"`
}

func (c *Client) newClient(conn *websocket.Conn, user *UserData) *Client {

	client := &Client{
		User: *user,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	return client
}

func (client *Client) Read() {
	defer func() {
		client.disconnect()
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error { client.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		_ = jsonMessage
		//handel messages
	}

}
func (client *Client) disconnect() {
	close(client.Send)
	client.Conn.Close()
}
