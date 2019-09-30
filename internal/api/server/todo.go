package server

import (
	"github.com/alexsniffin/go-api-example/internal/api/handlers"

	"github.com/go-chi/chi"
)

//Routes todo
func (s *Server) todoRoutes() {
	todoHandler := handlers.NewTodoHandler(s.render, s.sqlClient)

	s.router.Route("/todo", func(r chi.Router) {
		r.Get("/{id}", todoHandler.HandleGet)
		r.Delete("/{id}", todoHandler.HandleDelete)
		r.Post("/", todoHandler.HandlePost)
	})
}