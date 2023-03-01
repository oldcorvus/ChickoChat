package main

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type chatDatabase struct {
	users *mongo.Collection

	messages *mongo.Collection

	rooms *mongo.Collection
}
