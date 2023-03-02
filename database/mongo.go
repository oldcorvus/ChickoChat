package database

import (
	"chicko_chat/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatDatabase struct {
	Users *mongo.Collection

	Messages *mongo.Collection

	Rooms *mongo.Collection
}

// FindByEmail will be used to find a new user registry by email
func (c *ChatDatabase) FindByEmail(email string) (*data.Client, error) {

	// Find user by email
	var user = data.Client{}
	err := c.Users.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)

	if err != nil {
		// Checks if the user was not found
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ErrNoDocuments")
		}
		return nil, err
	}

	return &user, nil

}
