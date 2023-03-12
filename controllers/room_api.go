package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"chicko_chat/models"
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
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{"data": room})
		return

	}

}
