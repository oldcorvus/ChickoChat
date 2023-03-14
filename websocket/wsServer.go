package websocket

import (
	"chicko_chat/database"

	"log"
	"net/http"
	"chicko_chat/models"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize}

type WsServer struct {
	Manager *BrokerManager
}

func (server *WsServer) findBrokerbyRoomID(ID primitive.ObjectID) *data.Broker {
	for broker := range server.Manager.Brokers {
		if broker.Room.ID == ID {
			return broker
		}
	}
	return nil
}

func (server *WsServer) createBroker(room *data.ChatRoom) *data.Broker {
	broker := data.NewBroker(room)
	go server.Manager.RunBroker(broker)
	server.Manager.Brokers[broker] = true

	return broker
}

func (server *WsServer) ServeWs(w http.ResponseWriter, req *http.Request, roomId string, userId string) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	roomID, err := database.ObjectIDFromHex(roomId)
	userID, err := database.ObjectIDFromHex(userId)

	if err != nil {
		return
	}
	user := &data.UserData{
		ID: userID,
	}

	room := &data.ChatRoom{
		ID: roomID,
	}
	broker := server.findBrokerbyRoomID(room.ID)
	if broker == nil {
		broker = server.createBroker(room)
	}
	client := data.NewClient(socket, user, broker)
	clientManager := clientManager{
		client: client,
	}
	broker.Clients[client] = true
	broker.Join <- client
	defer func() { broker.Leave <- client }()
	go clientManager.clientWrite()
	 clientManager.clientRead()
}
