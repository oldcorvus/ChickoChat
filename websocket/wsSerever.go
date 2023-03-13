package websocket

import (
	"chicko_chat/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gorilla/websocket"
	"log"

)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize}


type WsServer struct {
	brokers      map[*data.Broker]bool
	
}


func (server *WsServer) findBrokerbyRoomID(ID primitive.ObjectID) *Broker {
	for broker := range server.brokers {
		if broker.ChatRoom.ID == ID {
			return broker
		}
	}
	return nil 
}

func (server *WsServer) createBroker(room *ChatRoom) *Broker {
	broker := NewBroker(room)
	go broker.RunBroker()
	server.brokers[broker] = true

	return broker
}

func (server *WsServer)  ServeWs(w http.ResponseWriter, req *http.Request, roomId string, userId string ) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	roomId , err = data.ObjectIDFromHex()
	userId , err = data.ObjectIDFromHex()

	if err != nil {
		return
	}
	user := &data.UserData{
		ID : userId,
	}

	client := data.newClient(socket, user)
	room := &data.ChatRoom{
		ID :roomId,
	}
	broker := server.findBrokerbyRoomID(room)
	if broker == nil {
		server.createBroker(room)
	}
	broker.Clients[client] = true
	broker.join <- client
	defer func() { broker.leave <- client }()
	go client.write()
	go client.read()
}
