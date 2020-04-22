package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/alexsniffin/go-starter/internal/todo-api/handlers/todo"
)

func NewRouter(todoHandler todo.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	buildRoutes(r, todoHandler)
	return r
}

func buildRoutes(r *chi.Mux, todoHandler todo.Handler) {
	r.Route("/api", func(r chi.Router) {
		r.Route("/todo", func(r chi.Router) {
			r.Get("/", todoHandler.HandleGet)
			r.Delete("/", todoHandler.HandleDelete)
			r.Post("/", todoHandler.HandlePost)
		})
		r.Get("/health", handleHealth)
	})
}

func handleHealth(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
}