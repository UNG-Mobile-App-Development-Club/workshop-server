package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Todo struct {
}

func main() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	http.ListenAndServe(":8080", router)
}
