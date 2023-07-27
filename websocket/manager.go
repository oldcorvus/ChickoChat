package websocket

import (
	"chicko_chat/database"
	data "chicko_chat/models"
	"log"

	"fmt"
	"time"

	"github.com/gorilla/websocket"
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
	broker.Mutex.Lock()
	broker.Clients[client] = true
	broker.Mutex.Unlock()
	message := data.ChatEvent{
		EventType: data.Subscribe,
		RoomID:    broker.Room.ID,
		UserID:    client.User.ID,
	}
	broker.Notification <- &message

	log.Printf("Client added. %d registered Clients", len(broker.Clients))

}

func (manager *BrokerManager) unregisterClient(client *data.Client, broker *data.Broker) {
	broker.Mutex.Lock()
	if _, ok := broker.Clients[client]; ok {
		delete(broker.Clients, client)
		close(client.Send)
	}
	broker.Mutex.Unlock()

	message := data.ChatEvent{
		EventType: data.Unsubscribe,
		RoomID:    broker.Room.ID,
		UserID:    client.User.ID,
	}
	broker.Notification <- &message

	log.Printf("Removed client. %d registered Clients", len(broker.Clients))

}

func (manager *BrokerManager) broadcastToClients(message *data.ChatEvent, broker *data.Broker) {
	message.Timestamp = time.Now()
	msg, err := manager.DB.SaveMessage(message)
	if err != nil {
		log.Print("message not sent: " + msg.Hex())

	}
	broker.Mutex.Lock()

	for client := range broker.Clients {
		select {
		case client.Send <- message:
			log.Print("message sent to: " + client.User.ID.Hex())
		case <-time.After(patience):
			log.Print("Skipping client: " + client.User.ID.Hex())
		default:
			log.Print("Deleting client: " + client.User.ID.Hex())
			close(client.Send)
			delete(broker.Clients, client)
		}
	}
	broker.Mutex.Unlock()

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
	manager.client.Conn.Close()
}

func (manager *clientManager) handleNewMessage(message *data.ChatEvent) {
	fmt.Println(message)
	switch message.EventType {
	case data.Broadcast:
		manager.client.Broker.Notification <- message

	case data.Subscribe, data.Unsubscribe:
		manager.notifyJoinedLeft(message)

	}

}

func (manager *clientManager) notifyJoinedLeft(message *data.ChatEvent) {

	manager.client.Send <- message
}
