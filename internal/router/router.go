package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/info", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {})

	return r
}
