package router

import (
	"go-effective-mobile/internal/api/user"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"go-effective-mobile/internal/api"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)

	r.Get("/info", func(w http.ResponseWriter, r *http.Request) {})

	// Users
	r.Post("/users", user.New)
	r.Get("/users/{userId}", user.Get)

	// Tasks
	r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {})

	// Service
	r.Get("/ping", api.Pong)

	return r
}
