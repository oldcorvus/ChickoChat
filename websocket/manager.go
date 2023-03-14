package websocket

import (
	"chicko_chat/database"
	"log"
	"chicko_chat/models"

	"github.com/gorilla/websocket"
	"fmt"
	"time"
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

// the amount of time to wait when pushing a message to
// a slow client or a client that closed after `range Clients` started.
const patience time.Duration = time.Second * 1

type BrokerManager struct {
	Brokers map[*data.Broker]bool
	DB      *database.ChatDatabase
}

type clientManager struct {
	client *data.Client
}

// runs broker accepting various requests
func (manager *BrokerManager) RunBroker(broker *data.Broker) {
	for {
		select {
		case client := <-broker.Join:
			manager.registerClient(client, broker)

		case client := <-broker.Leave:
			manager.unregisterClient(client, broker)

		case message := <-broker.Notification:
			manager.broadcastToClients(message, broker)
		}

	}
}

func (manager *BrokerManager) registerClient(client *data.Client, broker *data.Broker) {
	broker.Clients[client] = true

	log.Printf("Client added. %d registered Clients", len(broker.Clients))

}

