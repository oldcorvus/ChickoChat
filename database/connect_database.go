package database

import (
	"chicko_chat/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDatabseTest() *ChatDatabase {

	// Create mongo client configuration
	co := options.Client().ApplyURI("mongodb://localhost:27017")

	// Establish database connection
	client, err := mongo.NewClient(co)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	db := &ChatDatabase{
		Users: client.Database("chicko_chat").Collection("users_test"),
		Rooms: client.Database("chicko_chat").Collection("rooms_test"),
		Messages : client.Database("chicko_chat").Collection("message_test"),

	}
	defer func() {
		if err = db.Users.Drop(context.TODO()); err != nil {
			log.Fatal(err)
		}
		if err = db.Rooms.Drop(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	room := &data.ChatRoom{
		Title: "Data For Test",
	}

	_ , err = db.Rooms.InsertOne(context.TODO(),room )


	if err != nil {
		log.Fatal(err)
	}
	return db

}
