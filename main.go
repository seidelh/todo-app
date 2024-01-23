package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/todos", getTodos).Methods("GET")
	router.HandleFunc("/api/todos/{id}", getTodo).Methods("GET")
	router.HandleFunc("/api/todos", createTodo).Methods("POST")
	router.HandleFunc("/api/todos/{id}", updateTodo).Methods("PUT")
	router.HandleFunc("/api/todos/{id}", deleteTodo).Methods("DELETE")

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
