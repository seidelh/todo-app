package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
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
	collection := client.Database("todo").Collection("todos")
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var todo Todo
		err := cur.Decode(&todo)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todos)
}

func getTodo (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todo Todo
	collection := client.Database("todo").Collection("todos")
	params := mux.Vars(r)
	cur := collection.FindOne(context.Background(), bson.M{"_id": params["id"]})
	err := cur.Decode(&todo)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func createTodo (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todo Todo
	collection := client.Database("todo").Collection("todos")
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	todo.ID = uuid.New().String()
	_, err = collection.InsertOne(context.Background(), todo)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func updateTodo (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todo Todo
	collection := client.Database("todo").Collection("todos")
	params := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": params["id"]}, bson.M{"$set": todo})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
}

func deleteTodo (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := client.Database("todo").Collection("todos")
	params := mux.Vars(r)
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": params["id"]})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}