package main

import (
	"body-manipulation/api"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	router := chi.NewRouter()

	router.Post("/greetings", api.Greetings)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
