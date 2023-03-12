package controllers

import (
	"bytes"
	"chicko_chat/database"
	"chicko_chat/models"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestGetUserRoomsApi(t *testing.T) {
	db := database.ConnectDatabseTest()
	controller := Controller{
		DB: db,
	}
	r := SetUpRouter()
	r.POST("/user-rooms/", controller.GetUserRoomsApi)

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
	// Convert the JSON response to a map
	var response map[string][]data.ChatRoom
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	// Grab the value & whether or not it exists
	value, exists := response["data"]
	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, value, result)
}

func TestCreateRoomApi(t *testing.T) {
	db := database.ConnectDatabseTest()
	controller := Controller{
		DB: db,
	}
	r := SetUpRouter()
	r.POST("/create-room/", controller.CreateRoomApi)

	id, err := primitive.ObjectIDFromHex("640778694829658eebc2d55b")

	room := &data.ChatRoom{
		Title: "test",
		Clients : []primitive.ObjectID{id},
	}


	jsonValue, _ := json.Marshal(room)
	req, _ := http.NewRequest("POST", "/create-room/", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Convert the JSON response to a map
	var response map[string]data.ChatRoom
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	// Grab the value & whether or not it exists
	value, exists := response["data"]
	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, value.Clients, room.Clients)
}

func TestRoomHistoryApi(t *testing.T) {
	db := database.ConnectDatabseTest()
	controller := Controller{
		DB: db,
	}
	r := SetUpRouter()
	r.POST("/room-history/", controller.RoomHistoryApi)

	id, err := primitive.ObjectIDFromHex("640778694829658eebc2d55b")

	room := &data.ChatRoom{
		ID: id,
		Title: "test",
		Clients : []primitive.ObjectID{id},
	}

	var messages []interface{}
	for i := 1; i < 5; i++ {
		message := data.ChatEvent{
			EventType: data.Broadcast,
			ID:    primitive.NewObjectID(),
			RoomID:    id,
			Message:   "test message",
		}
		messages = append(messages, message)

	}
	_, err = db.Messages.InsertMany(context.TODO(), messages)
	if err != nil {
		t.Fatalf("failure in adding messages data to databse")
	}

	jsonValue, _ := json.Marshal(room)
	req, _ := http.NewRequest("POST", "/room-history/", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Convert the JSON response to a map
	var response map[string][]data.ChatEvent
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	// Grab the value & whether or not it exists
	value, exists := response["data"]
	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t , len(value), len(messages))
	assert.Equal(t , value[0], messages[0])

}
