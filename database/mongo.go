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

func (c *ChatDatabase) ConvertId(result *mongo.InsertOneResult)(string, error) {
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		return id.Hex(), nil
	} else {
	
		return "" , errors.New("failed vonverting")
	}

}


// Add user to the databse
func (c *ChatDatabase) AddUser(user *data.UserData) (string, error) {

	res, err := c.Users.InsertOne(context.TODO(), user)

	if err != nil {
		return "" , err
	}

	return c.ConvertId(res)

}

// FindByEmail will be used to find a new user registry by email
func (c *ChatDatabase) FindByEmail(email string) (*data.Client, error) {

	// Find user by email
	var user = data.Client{}
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
