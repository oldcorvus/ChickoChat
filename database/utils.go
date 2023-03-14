package database

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


// For converting result of insert document to  type primitive.ObjectID
func  convertId(result *mongo.InsertOneResult) (primitive.ObjectID, error) {
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		return id, nil
	} else {

		return primitive.NilObjectID, errors.New("failed converting")
	}

}

func ObjectIDFromHex(hexString string) (primitive.ObjectID, error){
	objID, err := primitive.ObjectIDFromHex(hexString)
	return objID, err
}