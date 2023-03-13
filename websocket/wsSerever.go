package websocket

import (
	"chicko_chat/database"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

type WsServer struct {
	brokers      map[*data.Broker]bool
	
}

