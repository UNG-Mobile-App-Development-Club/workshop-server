package main

import (
	"net/http"
	"encoding/json"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Todo struct {
    Id int `json:"id"`
    Content string `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}

var todos []Todo

func main() {

	 todos = []Todo{
	{Id: 1, Content: "Learn Go", CreatedAt: time.Now()},
	{Id: 2, Content: "Build a web app", CreatedAt: time.Now()},
    }

    router := chi.NewRouter()

    router.Use(middleware.Logger)

	router.Get("/todos", getTodos)
    router.Post("/todos", postTodo)

    http.ListenAndServe(":8080", router)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
    raw, err := json.Marshal(todos)
    if err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(raw)
}

func postTodo(w http.ResponseWriter, r *http.Request) {
    var todo Todo
    if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
	http.Error(w, err.Error(), http.StatusBadRequest)
	return
    }

    todo.Id = len(todos) + 1
    todo.CreatedAt = time.Now()
    todos = append(todos, todo)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}
