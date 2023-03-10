package controllers

import (
	"github.com/gin-gonic/gin"
	"chicko_chat/models"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"bytes"
	"encoding/json"
	"net/http"
	"chicko_chat/database"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine{
    router := gin.Default()
    return router
}



func TestGetUserRooms(t *testing.T) {
	db := database.ConnectDatabseTest()
	controller := Controller{
		DB: db,
	}
    r := SetUpRouter()
    r.POST("/user-rooms/",controller.GetUserRooms) 

	var rooms []interface{}
	id, err := primitive.ObjectIDFromHex("640778694829658eebc2d55b")

	user := &data.UserData{
		ID:    id,
		Email: "testuser@mail.com",
	}

	for i := 1; i < 5; i++ {
		room := data.ChatRoom{
			Title:   "test",
			Clients: []primitive.ObjectID{id, primitive.NewObjectID()},
		}
		rooms = append(rooms, room)

	}
	_, err = db.Rooms.InsertMany(context.TODO(), rooms)
	if err != nil {
		t.Fatalf("failure in adding rooms data to databse")
	}

	var result []data.ChatRoom
	result, err = db.GetHistoryOfUser(user)
	if err != nil || len(result) != 4 {
		t.Fatalf("failure in retriveing rooms data from databse")
	}


    jsonValue, _ := json.Marshal(user)
    req, _ := http.NewRequest("POST", "/user-rooms/", bytes.NewBuffer(jsonValue))

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
}