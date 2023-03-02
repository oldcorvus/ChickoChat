package database

import (
	"chicko_chat/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestFindByEmail(t *testing.T) {
	db := ConnectDatabseTest()
	//test for not existing email
	res, err := db.FindByEmail("moelcrow@gmail.com")
	if err == nil {
		t.Fatalf("user found!")
	}
	_ = res

	//adding user to databse
	db.Users.InsertOne(context.TODO(), bson.M{"email": "moelcrow@gmail.com"})

	_, err = db.FindByEmail("moelcrow@gmail.com")
	if err != nil {
		t.Fatalf("error user not found")
	}
	// Delete record
	db.Users.DeleteOne(context.TODO(), bson.M{"email": "moelcrow@gmail.com"})

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
