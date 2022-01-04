package database

import (
	"context"
	"fmt"
	"log"
	structs "main/utilities/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func CreateDB() {
	collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"id": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	)
}

func SetCollection(updatedCollection *mongo.Collection) {
	collection = updatedCollection
}

func InsertOne(messageData structs.Message) *mongo.InsertOneResult {
	insertResult, err := collection.InsertOne(context.TODO(), messageData)
	if err != nil {
		log.Fatal(err)
	}
	return insertResult
}

func DeleteOne(filter bson.D) *mongo.DeleteResult {
	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func ReplaceOne(filter bson.D, messageData structs.Message) *mongo.UpdateResult {
	res, err := collection.ReplaceOne(context.TODO(), filter, messageData)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func FindOne(filter bson.D) (structs.Message, error) {
	var message structs.Message
	err := collection.FindOne(context.TODO(), filter).Decode(&message)
	if err != nil {
		return message, err
	}
	fmt.Printf("Found a single document: %+v\n", message)
	return message, nil
}

func GetAll() ([]structs.Message, error) {
	//Define filter query for fetching specific document from collection
	filter := bson.D{{}} //bson.D{{}} specifies 'all documents'
	messages := []structs.Message{}

	//Perform Find operation & validate against the error.
	cur, findError := collection.Find(context.TODO(), filter)
	if findError != nil {
		return messages, findError
	}
	//Map result to slice
	for cur.Next(context.TODO()) {
		t := structs.Message{}
		err := cur.Decode(&t)
		if err != nil {
			return messages, err
		}
		messages = append(messages, t)
	}
	// once exhausted, close the cursor
	cur.Close(context.TODO())
	if len(messages) == 0 {
		return messages, mongo.ErrNoDocuments
	}
	return messages, nil
}
