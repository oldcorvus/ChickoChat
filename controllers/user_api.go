package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"chicko_chat/models"
)

func (c *Controller) GetUserRoomsApi(ctx *gin.Context) {
	// Validate input
	var user *data.UserData
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var rooms []data.ChatRoom
	// Search For rooms
	rooms, err := c.DB.GetHistoryOfUser(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	ctx.JSON(http.StatusOK, gin.H{"data": rooms})


}
