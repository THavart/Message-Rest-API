package messaging

import (
	"encoding/json"
	"fmt"
	"main/utilities/database"
	structs "main/utilities/structs"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Get all messages
func GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	messages, _ := database.GetAll()
	json.NewEncoder(w).Encode(messages)
}

// Get single message
func GetMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through messages and find one with the id from the params

	message, err := getMessageFromDb(params["id"])

	if err != nil {
		json.NewEncoder(w).Encode(&structs.Message{})
		return
	}

	json.NewEncoder(w).Encode(message)
}

// Add new message
func CreateMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got Create Message Request!\n")
	w.Header().Set("Content-Type", "application/json")
	myuuid := uuid.NewV4()
	var messageData structs.Message
	_ = json.NewDecoder(r.Body).Decode(&messageData)
	messageData.ID = myuuid.String()
	messageData.Timestamp = time.Now().UTC().String()
	insertResult := database.InsertOne(messageData)
	fmt.Println("Stored Message info in database successfully: ", insertResult.InsertedID)
	json.NewEncoder(w).Encode(messageData)
}

// Delete message
func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Delete Message received!\n")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	message, err := getMessageFromDb(params["id"])

	if err != nil {
		txtErr := fmt.Errorf("Unable to locate Message in db")
		fmt.Println(txtErr.Error())
		json.NewEncoder(w).Encode(nil)
		return
	}

	fmt.Printf("Deleting.. please wait\n")

	res := database.DeleteOne(bson.D{primitive.E{Key: "id", Value: message.ID}})

	fmt.Printf("Removed message ID: %s, removed from database: %v\n", message.ID, res.DeletedCount)
	json.NewEncoder(w).Encode(message)
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Update Request received!\n")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var incoming structs.Message
	_ = json.NewDecoder(r.Body).Decode(&incoming)
	message, err := getMessageFromDb(params["id"])

	messageData := structs.Message{ID: message.ID, Content: incoming.Content, Author: incoming.Author, Timestamp: time.Now().UTC().String()}

	if err != nil {
		txtErr := fmt.Errorf("Unable to locate Message in db")
		fmt.Println(txtErr.Error())
		json.NewEncoder(w).Encode(nil)
		return
	}

	database.ReplaceOne(bson.D{primitive.E{Key: "id", Value: messageData.ID}}, messageData)
	fmt.Println("Replaced Message info in database successfully: ", messageData.ID)
	json.NewEncoder(w).Encode(messageData)
}

func getMessageFromDb(containerID string) (structs.Message, error) {
	message, err := database.FindOne(bson.D{{Key: "id", Value: containerID}})
	return message, err
}
