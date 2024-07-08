package router

import (
	"github.com/go-chi/chi/v5"
	"go-effective-mobile/internal/api"
	"net/http"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/info", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("/ping", api.Pong)

	return r
}
