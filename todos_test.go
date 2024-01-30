package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetTodos(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/todos", nil)
	handler := http.HandlerFunc(getTodos)
	handler.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestCreateTodo(t *testing.T) {
	payload := `{"title": "Test Todo", "description": "Test Description", "due": "2024-01-30", "severity": "Medium", "completed": false}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/todos", strings.NewReader(payload))
	handler := http.HandlerFunc(createTodo)
	handler.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var createdTodo Todo
	err := json.Unmarshal(w.Body.Bytes(), &createdTodo)
	assert.NoError(t, err)
}

func TestUpdateTodo(t *testing.T) {
	payload := `{"title": "Updated Title", "description": "Updated Description", "due": "2024-01-31", "severity": "High", "completed": true}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/todos/some_todo_id", strings.NewReader(payload))
	handler := http.HandlerFunc(updateTodo)
	handler.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var updatedTodo Todo
	err := json.Unmarshal(w.Body.Bytes(), &updatedTodo)
	assert.NoError(t, err)
}

func TestDeleteTodo(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/todos/some_todo_id", nil)
	handler := http.HandlerFunc(deleteTodo)
	handler.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
}
