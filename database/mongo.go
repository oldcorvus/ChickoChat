package database

import (
	"context"
	"errors"
	"chicko_chat/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatDatabase struct {
	Users *mongo.Collection

	Messages *mongo.Collection

	Rooms *mongo.Collection
}

func (c *ChatDatabase) ConvertId(result *mongo.InsertOneResult) (primitive.ObjectID, error) {
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		return id, nil
	} else {

		return primitive.NilObjectID, errors.New("failed converting")
	}

}

// Add user to the databse
func (c *ChatDatabase) AddUser(user *data.UserData) (primitive.ObjectID, error) {

	res, err := c.Users.InsertOne(context.TODO(), user)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return c.ConvertId(res)

}

// FindByEmail will be used to find a new user registry by email
func (c *ChatDatabase) FindByEmail(email string) (*data.UserData, error) {

	// Find user by email
	var user = data.UserData{}
	err := c.Users.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)

	if err != nil {
		// Checks if the user was not found
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil

}

func (c *ChatDatabase) CreateRoom(room *data.ChatRoom) (*data.ChatRoom, error) {

	res, err := c.Rooms.InsertOne(context.TODO(), room)

	if err != nil {
		return room, err
	}

	room.ID, err = c.ConvertId(res)

	if err != nil {
		return room, err
	}
	return room, nil

}

func (c *ChatDatabase) AddClientToRoom(room *data.ChatRoom, user *data.UserData) (*data.ChatRoom, error) {
	change := bson.M{
		"$push": bson.M{
			"users": user.ID,
		},
	}
	filter := bson.M{
		"_id": room.ID,
	}
	_, err := c.Rooms.UpdateOne(context.Background(), filter, change)
	if err != nil {
		return room, err
	}
	room.Clients = append(room.Clients, user.ID)
	return room, nil

}
