package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
	"time"
)

func TestFindByEmail(t *testing.T) {

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
		Users: client.Database("chicko_chat").Collection("users"),
	}
	//test for not existing email
	res, err := db.FindByEmail("moelcrow@gmail.com")
	if err == nil {
		t.Fatalf("user found!")
	}

	//adding user to databse
	db.Users.InsertOne(context.TODO(), bson.D{{"email", "moelcrow@gmail.com"}})

	res, err = db.FindByEmail("moelcrow@gmail.com")
	if err != nil {
		t.Fatalf("error user not found")
	}

}
