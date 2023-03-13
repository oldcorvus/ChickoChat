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
	Send chan ChatEvent `json:"-"`
	// Broker for connection
	Broker *Broker
}

type UserData struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email  string             `json:"email" bson:"email" binding:"required"`
	Name   string             `json:"name" bson:"name"`
	Active bool               `json:"active" bson:"active"`
}

func  NewClient(conn *websocket.Conn, user *UserData, broker *Broker) *Client {

	client := &Client{
		User: *user,
		Conn: conn,
		Send: make(chan []byte, 256),
		Broker: broker ,

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
		var msg ChatEvent
		// Read in a new message as JSON and map it to a Message object
		msg, err := client.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		// handel message
	}

}


func (client *Client) disconnect() {
	close(client.Send)
	client.Conn.Close()
}

