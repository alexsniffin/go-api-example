package server

import (
	"fmt"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexsniffin/go-api-example/internal/api/models"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

//Routes todo
func (s *Server) todoRoutes() {
	s.router.Route("/todo", func(r chi.Router) {
		r.Get("/{id}", s.handleGet)
		r.Delete("/{id}", s.handleDelete)
		r.Post("/", s.handlePost)
	})
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprint("Failed to decode todoID: " + chi.URLParam(r, "id")))
		s.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Error decoding TodoID",
		})
      	return
	}

	todo, err := s.postgresDb.GetTodo(todoID)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	s.render.JSON(w, http.StatusOK, todo)
	return
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprint("Failed to decode todoID: " + chi.URLParam(r, "id")))
		s.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Error decoding TodoID",
		})
      	return
	}

	count, err := s.postgresDb.DeleteTodo(todoID)
	if err != nil {
		s.render.JSON(w, http.StatusInternalServerError, models.Error{
			Message: "Error delete todo",
		})
		return
	}

	if count == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Debug().Msg(fmt.Sprint(count, " rows deleted for ", todoID))

	w.WriteHeader(200)
	return
}

func (s *Server) handlePost(w http.ResponseWriter, r *http.Request) {
    var todo models.Todo
    err := json.NewDecoder(r.Body).Decode(&todo)
    if err != nil {
		log.Error().Msg(fmt.Sprint("Failed to decode todo body: ", r.Body))
		s.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Error decoding body",
		})
		return
	}
	
	id, err := s.postgresDb.PostTodo(todo)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprint("Failed to insert todo record: ", r.Body))
		s.render.JSON(w, http.StatusInternalServerError, models.Error{
			Message: "Error inserting record to database",
		})
	}

	s.render.JSON(w, http.StatusOK, id)
	return
}