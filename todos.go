package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

type Todo struct {
	ID          string `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string `json:"title,omitempty" bson:"title,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	DueDate     string `json:"due,omitempty" bson:"due,omitempty"`
	Severity    string `json:"severity,omitempty" bson:"severity,omitempty"`
	Completed   bool   `json:"completed,omitempty" bson:"completed,omitempty"`
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todos []Todo
	cursor, err := client.Database("todo").Collection("todos").Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	if err := cursor.Err(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todos)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	var todo Todo
	if err := client.Database("todo").Collection("todos").FindOne(context.Background(), bson.M{"_id": params["id"]}).Decode(&todo); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := client.Database("todo").Collection("todos").InsertOne(context.Background(), todo); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := client.Database("todo").Collection("todos").UpdateOne(context.Background(), bson.M{"_id": params["id"]}, bson.M{"$set": todo}); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	if _, err := client.Database("todo").Collection("todos").DeleteOne(context.Background(), bson.M{"_id": params["id"]}); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
