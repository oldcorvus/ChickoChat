package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"chicko_chat/models"
	"chicko_chat/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (c *Controller) CreateRoomApi(ctx *gin.Context) {
	// Validate input
	var room *data.ChatRoom
	if err := ctx.ShouldBindJSON(&room); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create room
	room.CreatedAt = time.Now()
	room, err := c.DB.CreateRoom(room)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": room})

}


func (c *Controller) RoomHistoryApi(ctx *gin.Context) {
	// Validate input
	var room *data.ChatRoom
	if err := ctx.ShouldBindJSON(&room); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var messages []data.ChatEvent
	messages, err := c.DB.GetHistoryOfRoom(room)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": messages})
}

func (c *Controller) AddUserToRoomApi(ctx *gin.Context) {
	// Validate input
	var room *data.ChatRoom
	if err := ctx.ShouldBindJSON(&room); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	room, err := c.DB.AddClientToRoom(room)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": room})
}



func (c *Controller) JoinRoom(ctx *gin.Context) {
	// Validate input
	var room *data.ChatRoom
	roomId := ctx.Query("roomId")
	userId := ctx.Query("userId")
	id , err := database.ObjectIDFromHex(roomId)
	userID, err := database.ObjectIDFromHex(userId)
	room , err = c.DB.FindRoomByID(id )
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "room not found"})
		return
	}	
	var  found bool = false
	for i := range room.Clients {
		if room.Clients[i] == userID {
			found = true
		}
	}
	if found != true {
		room , err = c.DB.AddClientToRoom(&data.ChatRoom{
			ID : id,
			Clients : []primitive.ObjectID{userID},
		})
	}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "room not found"})
		return
	}	

	 ctx.HTML(http.StatusOK, "chat.tmpl", gin.H{
		"title": "Sample Front",
		"name":  "Moel",
		"roomId" : roomId,
		"userId": userId,
	})
	
}



