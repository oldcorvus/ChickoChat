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

func (manager *BrokerManager) unregisterClient(client *data.Client, broker *data.Broker) {
	if _, ok := broker.Clients[client]; ok {
		delete(broker.Clients, client)
		close(client.Send)
	}

	log.Printf("Removed client. %d registered Clients", len(broker.Clients))

}

func (manager *BrokerManager) broadcastToClients(message *data.ChatEvent, broker *data.Broker) {
	msg, err := manager.DB.SaveMessage(message)
	if err != nil {
		log.Print("message not sent: " + msg.Hex())

	}
	for client := range broker.Clients {
		select {
		case client.Send <- message:
			log.Print("message sent to: " + client.User.Email)
		case <-time.After(patience):
			log.Print("Skipping client: " + client.User.Email)
		default:
			log.Print("Deleting client: " + client.User.Email)
			close(client.Send)
			delete(broker.Clients, client)
		}
	}
}

func (manager *clientManager) clientRead() {
	defer func() {
		manager.ClientDisconnect()
	}()

	manager.client.Conn.SetReadLimit(maxMessageSize)
	manager.client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	manager.client.Conn.SetPongHandler(func(string) error { manager.client.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		var msg data.ChatEvent
		// Read in a new message as JSON and map it to a Message object
		err := manager.client.Conn.ReadJSON(&msg)
		log.Printf("read message")
		fmt.Println(msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		// handel message
		manager.handleNewMessage(&msg)
	}

}

func (manager *clientManager) clientWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		manager.client.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-manager.client.Send:
			manager.client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				manager.client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			manager.client.Conn.WriteJSON(message)

		case <-ticker.C:
			manager.client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := manager.client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (manager *clientManager) ClientDisconnect() {
	close(manager.client.Send)
	manager.client.Conn.Close()
}

