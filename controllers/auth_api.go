package controllers

import (
	"net/http"

	"chicko_chat/models"
	"github.com/gin-gonic/gin"
)



func (c *Controller) StartConversationApi(ctx *gin.Context) {
	// Validate input
	var user *data.UserData
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Search For existing user
	usr, err := c.DB.FindByEmail(user.Email)
	if err == nil{
		if usr.Name == "" {
			usr.Name = user.Name
			usr , err = c.DB.UpdateUserName(usr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			
		}
		
		ctx.JSON(http.StatusOK, gin.H{"data": usr})
		return

	}

	id, err := c.DB.AddUser(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	user.ID = id
	ctx.JSON(http.StatusOK, gin.H{"data": user})

}

