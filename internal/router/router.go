package router

import (
	"go-effective-mobile/internal/api/tasks"
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
	r.Get("/users/", user.GetAll)
	r.Post("/users", user.New)
	r.Get("/users/{userId}", user.Get)
	r.Patch("/users/{userId}", user.Update)
	r.Delete("/users/{userId}", user.Delete)

	// User's tasks
	r.Get("/users/{userId}/tasks", tasks.GetUserTasks)
	r.Get("/users/{userId}/tasks/{taskId}", tasks.Get)
	r.Post("/users/{userId}/tasks", tasks.New)
	r.Patch("/users/{userId}/tasks/{taskId}/start", tasks.Start)
	r.Patch("/users/{userId}/tasks/{taskId}/stop", tasks.Stop)

	// Tasks
	r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {})

	// Service
	r.Get("/ping", api.Pong)

	return r
}
