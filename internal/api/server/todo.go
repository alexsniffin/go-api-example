package server

import (
	"github.com/alexsniffin/go-api-example/internal/api/handlers"
	"github.com/alexsniffin/go-api-example/internal/api/store"

	"github.com/go-chi/chi"
)

//Routes todo
func (s *Server) todoRoutes() {
	store := store.NewTodoStore(s.sqlClient)
	todoHandler := handlers.NewTodoHandler(s.render, store)

	s.router.Route("/todo", func(r chi.Router) {
		r.Get("/", todoHandler.HandleGet)
		r.Delete("/", todoHandler.HandleDelete)
		r.Post("/", todoHandler.HandlePost)
	})
}
