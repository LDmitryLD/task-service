package router

import (
	"net/http"
	"projects/LDmitryLD/task-service/gateway/internal/modules"

	"github.com/go-chi/chi/v5"
)

func NewRouter(controllers *modules.Controllers) *chi.Mux {
	r := chi.NewRouter()

	setDefaultRoutes(r)

	r.Post("/api/users", controllers.User.Create)
	r.Get("/api/users/{id}", controllers.User.Profile)

	r.Post("/api/tasks", controllers.Task.Create)
	r.Get("/api/tasks/{userId}", controllers.Task.List)
	r.Delete("/api/users/{userId}/tasks/{taskId}", controllers.Task.Delete)

	return r
}

func setDefaultRoutes(r *chi.Mux) {
	r.Get("/swagger", swaggerUI)
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))).ServeHTTP(w, r)
	})
}
