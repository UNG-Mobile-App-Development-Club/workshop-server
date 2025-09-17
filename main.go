package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Todo struct {
	Id        int       `json:"id"`
	Task      string    `json:"task"`
	CreatedAt time.Time `json:"created_at"`
}

var todos []Todo

func main() {

	todos = []Todo{
		{Id: 1, Task: "Gurt", CreatedAt: time.Now()},
		{Id: 2, Task: "Yo", CreatedAt: time.Now()},
	}
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	http.ListenAndServe(":8080", router)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	// Mimics a database call
	w.Header().Set("Content-Type", "application/json")
	raw, _ := json.Marshal(todos)
	w.Write(raw)
}

func postTodos(w http.ResponseWriter, r *http.Request) {
	// Deserializing
	var newTodo Todo
	json.NewDecoder(r.Body).Decode(&newTodo)

	// Data processing
	newTodo.Id = len(todos) + 1
	newTodo.CreatedAt = time.Now()

	// Append to database
	todos = append(todos, newTodo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	raw, _ := json.Marshal(newTodo)

	w.Write(raw)
}
