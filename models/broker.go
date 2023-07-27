package data

import "sync"

// struct representing  client connections
type Broker struct {
	// Registered Clients.
	Clients map[*Client]bool

	// messages from the Clients.
	Notification chan *ChatEvent

	// Register requests from the Clients.
	Join chan *Client

	// Unregister requests from Clients.
	Leave chan *Client

	Room *ChatRoom

	Mutex sync.Mutex
}

func NewBroker(room *ChatRoom) *Broker {
	return &Broker{
		Notification: make(chan *ChatEvent, 100),
		Join:         make(chan *Client),
		Leave:        make(chan *Client),
		Clients:      make(map[*Client]bool),
		Room:         room,
	}
}
