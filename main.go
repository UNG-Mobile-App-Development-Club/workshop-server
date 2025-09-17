package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Todo struct {
	Id int `json:"id"`
	Task string `json:"task"`
	CreatedAt time.Time `json:"created_at"`
}

var todos = []Todo

func main() {

	todos = []Todo{
		{Id: 1, Task: "Learn Go", CreatedAt: time.Now()},
		{Id: 2, Task: "Build a REST API", CreatedAt: time.Now()},
	}
    router := chi.NewRouter()

    router.Use(middleware.Logger)

    http.ListenAndServe(":8080", router)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	// Mimics a database call
	raw, _ := json.Marshal(todos)

	w.Header().Set("Content-Type", "application/json")
	w.Write(raw)
}

func postTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	json.NewDecoder(r.Body).Decode(&newTodo)
	newTodo.Id = len(todos) + 1
	newTodo.CreatedAt = time.Now()
	todos = append(todos, newTodo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	raw, _ := json.Marshal(newTodo)
	w.Write(raw)
}