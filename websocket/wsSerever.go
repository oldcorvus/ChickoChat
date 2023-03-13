package websocket

import (
	"chicko_chat/database"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

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