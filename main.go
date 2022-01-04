package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Message struct (Model)
type Message struct {
	ID        string  `json:"id"`
	Content   string  `json:"content"`
	Author    *Author `json:"author"`
	Timestamp string  `json:""`
}

// Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Init servers var as a slice Message struct
var collection *mongo.Collection
var router *mux.Router

// Get all servers
func getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	servers, _ := GetAllMessages()
	json.NewEncoder(w).Encode(servers)
}

// Get single server
func getMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through messages and find one with the id from the params

	server, err := getMessageFromDb(params["id"])

	if err != nil {
		json.NewEncoder(w).Encode(&Message{})
		return
	}

	json.NewEncoder(w).Encode(server)
}

// Add new server
func createMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got Create Message Request!\n")
	w.Header().Set("Content-Type", "application/json")
	myuuid := uuid.NewV4()
	var serverData Message
	_ = json.NewDecoder(r.Body).Decode(&serverData)
	serverData.ID = myuuid.String()
	serverData.Timestamp = time.Now().UTC().String()
	insertResult, err := collection.InsertOne(context.TODO(), serverData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Stored Message info in database successfully: ", insertResult.InsertedID)
	json.NewEncoder(w).Encode(serverData)
}

// Delete message
func deleteMessage(w http.ResponseWriter, r *http.Request) {
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

	filter := bson.D{primitive.E{Key: "id", Value: message.ID}}

	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Removed message ID: %s, removed from database: %v\n", message.ID, res.DeletedCount)
	json.NewEncoder(w).Encode(message)
}

func updateMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Update Request received!\n")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var incoming Message
	_ = json.NewDecoder(r.Body).Decode(&incoming)
	message, err := getMessageFromDb(params["id"])

	messageData := Message{ID: message.ID, Content: incoming.Content, Author: incoming.Author, Timestamp: time.Now().UTC().String()}

	if err != nil {
		txtErr := fmt.Errorf("Unable to locate Message in db")
		fmt.Println(txtErr.Error())
		json.NewEncoder(w).Encode(nil)
		return
	}

	filter := bson.D{primitive.E{Key: "id", Value: messageData.ID}}

	collection.ReplaceOne(context.TODO(), filter, messageData)

	fmt.Println("Replaced Message info in database successfully: ", messageData.ID)
	json.NewEncoder(w).Encode(messageData)
}

// Main function
func main() {
	// Init router
	connect()

	collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"id": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	)

	router = mux.NewRouter()

	// Route handles & endpoints
	router.HandleFunc("/messages", getMessages).Methods("GET")
	router.HandleFunc("/messages/{id}", getMessage).Methods("GET")
	router.HandleFunc("/messages", createMessage).Methods("POST")
	router.HandleFunc("/messages/{id}", updateMessage).Methods("PUT")
	router.HandleFunc("/messages/{id}", deleteMessage).Methods("DELETE")

	fmt.Println("Starting server....")
	// Start server
	log.Println(http.ListenAndServe(":10000", router))
}

func getMessageFromDb(containerID string) (Message, error) {
	var server Message
	filter := bson.D{{Key: "id", Value: containerID}}
	err := collection.FindOne(context.TODO(), filter).Decode(&server)
	if err != nil {
		return server, err
	}

	fmt.Printf("Found a single document: %+v\n", server)
	return server, err
}

func connect() {
	fmt.Printf("IP Address is: %s\n", GetOutboundIP())
	val, present := os.LookupEnv("MONGODB")
	var commandString string
	var databaseString string

	if present {
		commandString = fmt.Sprintf("mongod --bind_ip %s", val)
		databaseString = fmt.Sprintf("mongodb://%s:27017", val)
	} else {
		commandString = fmt.Sprintf("mongod --bind_ip %s", GetOutboundIP())
		databaseString = fmt.Sprintf("mongodb://%s:27017", GetOutboundIP())
	}

	go func() {
		exec.Command("bash", "-c", commandString).Output()
	}()

	time.Sleep(2 * time.Second)
	fmt.Printf("Connecting to database: %s\n", databaseString)
	clientOptions := options.Client().ApplyURI(databaseString)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Println(err)
	}

	fmt.Println("Connected to MongoDB!")
	collection = client.Database("mydb").Collection("messages")
}

func GetAllMessages() ([]Message, error) {
	//Define filter query for fetching specific document from collection
	filter := bson.D{{}} //bson.D{{}} specifies 'all documents'
	servers := []Message{}

	//Perform Find operation & validate against the error.
	cur, findError := collection.Find(context.TODO(), filter)
	if findError != nil {
		return servers, findError
	}
	//Map result to slice
	for cur.Next(context.TODO()) {
		t := Message{}
		err := cur.Decode(&t)
		if err != nil {
			return servers, err
		}
		servers = append(servers, t)
	}
	// once exhausted, close the cursor
	cur.Close(context.TODO())
	if len(servers) == 0 {
		return servers, mongo.ErrNoDocuments
	}
	return servers, nil
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
