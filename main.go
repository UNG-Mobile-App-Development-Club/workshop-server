package main

import (
	"encoding/json"
	"net/http"
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
