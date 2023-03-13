package data

import (
	"log"
	"time"
)

// the amount of time to wait when pushing a message to
// a slow client or a client that closed after `range Clients` started.
const patience time.Duration = time.Second * 1

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
}

func NewBroker(room *ChatRoom) *Broker {
	return &Broker{
		Notification: make(chan *ChatEvent),
		Join:         make(chan *Client),
		Leave:        make(chan *Client),
		Clients:      make(map[*Client]bool),
		Room:         room,
	}
}

// runs broker accepting various requests
func (br *Broker) RunBroker() {
	for {
		select {
		case client := <-br.Join:
			br.registerClient(client)

		case client := <-br.Leave:
			br.unregisterClient(client)

		case message := <-br.Notification:
			br.broadcastToClients(message)
		}

	}
}
func (br *Broker) registerClient(client *Client) {
	br.Clients[client] = true

	log.Printf("Client added. %d registered Clients", len(br.Clients))

}

