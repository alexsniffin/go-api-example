package server

import (
	"github.com/alexsniffin/go-api-example/internal/api/store"
	"github.com/alexsniffin/go-api-example/internal/api/handlers"

	"github.com/go-chi/chi"
)

//Routes todo
func (s *Server) todoRoutes() {
	store := store.NewTodoStore(s.sqlClient)
	todoHandler := handlers.NewTodoHandler(s.render, store)

	s.router.Route("/todo", func(r chi.Router) {
		r.Get("/{id}", todoHandler.HandleGet)
		r.Delete("/{id}", todoHandler.HandleDelete)
		r.Post("/", todoHandler.HandlePost)
	})
}
