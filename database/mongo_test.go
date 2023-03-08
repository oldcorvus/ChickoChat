package database

import (
	"chicko_chat/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"

)

func TestFindByEmail(t *testing.T) {

	db := ConnectDatabseTest()
	//test for not existing email
	_, err := db.FindByEmail("moelcrow@gmail.com")
	if err == nil {
		t.Fatalf("user found!")
	}

	//adding user to databse
	db.Users.InsertOne(context.TODO(), bson.M{"email": "moelcrow@gmail.com"})

	_, err = db.FindByEmail("moelcrow@gmail.com")
	if err != nil {
		t.Fatalf("error user not found")
	}

}

func TestAddUser(t *testing.T) {

	db := ConnectDatabseTest()
	user := &data.UserData{
		Email:  "testregister@gmail.com",
		Name:   "test user",
		Active: true,
	}

	_, err := db.AddUser(user)

	if err != nil {
		t.Fatalf("failure in adding user data to databse")
	}
	var usr = data.UserData{}
	res := db.Users.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&usr)

	if res != nil {
		t.Fatalf("failure finding added user ")
	}

}

func TestCreateRoom(t *testing.T) {

	db := ConnectDatabseTest()
	room := &data.ChatRoom{
		Title: "Test",
	}

	_, err := db.CreateRoom(room)

	if err != nil {
		t.Fatalf("failure in adding room data to database")
	}
	var rm = data.ChatRoom{}
	res := db.Rooms.FindOne(context.TODO(), bson.M{"title": room.Title}).Decode(&rm)

	if res != nil {
		t.Fatalf("failure finding added room ")
	}

}

func TestAddClientToRoom(t *testing.T) {

	db := ConnectDatabseTest()
	room := &data.ChatRoom{
		Title:   "Data For Test",
		Clients: make([]primitive.ObjectID, 0),
	}

	_, err := db.Rooms.InsertOne(context.TODO(), room)

	err = db.Rooms.FindOne(context.TODO(), bson.M{"title": "Data For Test"}).Decode(&room)
	if err != nil {
		t.Fatalf("failure finding added room ")
	}

	user := &data.UserData{
		Email:  "testregister@gmail.com",
		Name:   "test user",
		Active: true,
	}
	db.Users.InsertOne(context.TODO(), user)
	err = db.Users.FindOne(context.TODO(), user).Decode(&user)
	if err != nil {
		t.Fatalf("failure finding added user ")

	}
	res, err := db.AddClientToRoom(room, user)
	if res.Clients[0] != user.ID {
		t.Fatalf("failure adding user into room ")

	}
}


func TestAddMessage(t *testing.T) {

	db := ConnectDatabseTest()
	message := &data.ChatEvent{
		EventType : data.Broadcast,
		UserID :    primitive.NewObjectID(),
		RoomID :    primitive.NewObjectID(),
		Message : "test message",

	}

	res , err := db.SaveMessage(message)

	if err != nil {
		t.Fatalf("failure in adding message data to databse")
	}
	msg := &data.ChatEvent{}
	err = db.Messages.FindOne(context.TODO(), bson.M{"_id": res}).Decode(&msg)

	if err != nil || msg.ID != res {
		t.Fatalf("failure finding added message ")
	}

}

func TestGetHitoryOfRoom(t *testing.T) {

	db := ConnectDatabseTest()
	var messages  []interface{} 
	id , err :=  primitive.ObjectIDFromHex("640778694829658eebc2d55b")

	room := &data.ChatRoom {
		ID : id,
		Title : "test",
	}

	for i := 1; i < 5; i++ {
		message := data.ChatEvent{
			EventType : data.Broadcast,
			UserID :    primitive.NewObjectID(),
			RoomID :   id ,
			Message : "test message" ,
	
		}
		message2 := data.ChatEvent{
			EventType : data.Broadcast,
			UserID :    primitive.NewObjectID(),
			RoomID :    primitive.NewObjectID() ,
			Message : "another room" ,
	
		}
		messages = append(messages, message)
		messages = append(messages, message2)

	}
	_, err = db.Messages.InsertMany(context.TODO(), messages)
	if err != nil {
		t.Fatalf("failure in adding messages data to databse")
	}
	var result []data.ChatEvent
	result , err = db.GetHistoryOfRoom(room)
	if err != nil || len(result) != 4 {
		t.Fatalf("failure in retriveing messages data from databse")
	}


}