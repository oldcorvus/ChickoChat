package main

import (
	"chicko_chat/log"

	"github.com/gorilla/websocket"
)

// client represents a single chatting user.
type client struct {
	socket *websocket.Conn

	send chan []byte

	room *room

	logger logger.Logger
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
