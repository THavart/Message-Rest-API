package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"

	database "main/utilities/database"
	messenger "main/utilities/messaging"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Init servers var as a slice Message struct

var router *mux.Router

// Main function
func main() {
	// Init router
	connect()

	database.CreateDB()

	router = mux.NewRouter()

	// Route handles & endpoints
	router.HandleFunc("/messages", messenger.GetMessages).Methods("GET")
	router.HandleFunc("/messages/{id}", messenger.GetMessage).Methods("GET")
	router.HandleFunc("/messages", messenger.CreateMessage).Methods("POST")
	router.HandleFunc("/messages/{id}", messenger.UpdateMessage).Methods("PUT")
	router.HandleFunc("/messages/{id}", messenger.DeleteMessage).Methods("DELETE")

	fmt.Println("Starting server....")
	// Start server
	log.Println(http.ListenAndServe(":10000", router))
}

func connect() {
	fmt.Printf("IP Address is: %s\n", getOutboundIP())
	val, present := os.LookupEnv("MONGODB")
	var command, databaseURL string

	if present {
		command = fmt.Sprintf("mongod --bind_ip %s", val)
		databaseURL = fmt.Sprintf("mongodb://%s:27017", val)
	} else {
		command = fmt.Sprintf("mongod --bind_ip %s", getOutboundIP())
		databaseURL = fmt.Sprintf("mongodb://%s:27017", getOutboundIP())
	}

	go func() {
		exec.Command("bash", "-c", command).Output()
	}()

	fmt.Printf("Connecting to database: %s\n", databaseURL)
	clientOptions := options.Client().ApplyURI(databaseURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	database.SetCollection(client.Database("mydb").Collection("messages"))
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
