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
	Notification chan []byte

	// Register requests from the Clients.
	Join chan *Client

	// Unregister requests from Clients.
	Leave chan *Client

	Room *ChatRoom
}

func newBroker(room *ChatRoom) *Broker {
	return &Broker{
		Notification: make(chan []byte),
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

func (br *Broker) unregisterClient(client *Client) {
	if _, ok := br.Clients[client]; ok {
		delete(br.Clients, client)
		close(client.Send)
	}

	log.Printf("Removed client. %d registered Clients", len(br.Clients))

}

func (br *Broker) broadcastToClients(message []byte) {

	for client := range br.Clients {
		select {
		case client.Send <- message:
		case <-time.After(patience):
			log.Print("Skipping client: " + client.User.Name)
		default:
			log.Print("Deleting client: " + client.User.Name)
			close(client.Send)
			delete(br.Clients, client)
		}
	}
}
