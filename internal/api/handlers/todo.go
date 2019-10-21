package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alexsniffin/go-api-example/internal/api/models"
	"github.com/alexsniffin/go-api-example/internal/api/store"

	"github.com/rs/zerolog/log"
	"github.com/unrolled/render"
)

//Todo todo
type Todo interface {
	HandleGet(w http.ResponseWriter, r *http.Request)
	HandleDelete(w http.ResponseWriter, r *http.Request)
	HandlePost(w http.ResponseWriter, r *http.Request)
}

//TodoHandler todo
type TodoHandler struct {
	render *render.Render
	store  store.Todo
}

//NewTodoHandler todo
func NewTodoHandler(render *render.Render, store store.Todo) *TodoHandler {
	return &TodoHandler{
		render: render,
		store:  store,
	}
}

//HandleGet todo
func (t *TodoHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	todoIDStr := r.URL.Query().Get("id")
	if todoIDStr == "" {
		log.Error().Msg(fmt.Sprint("Missing todoID in request"))
		err := t.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Missing query parameter: id",
		})
		if err != nil {
			log.Error().Err(err)
		}
		return
	}

	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		err := t.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "id must be an integer",
		})
		if err != nil {
			log.Error().Err(err)
		}
		return
	}

	todo, err := t.store.GetTodo(r.Context(), todoID)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = t.render.JSON(w, http.StatusOK, todo)
	if err != nil {
		log.Error().Err(err)
	}
}

//HandleDelete todo
func (t *TodoHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	todoIDStr := r.URL.Query().Get("id")
	if todoIDStr == "" {
		log.Error().Msg(fmt.Sprint("Missing todoID in request"))
		err := t.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Missing query parameter: id",
		})
		if err != nil {
			log.Error().Err(err)
		}
		return
	}

	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		err := t.render.JSON(w, http.StatusInternalServerError, models.Error{
			Message: "Error decoding id to an integer",
		})
		if err != nil {
			log.Error().Err(err)
		}
		return
	}

	count, err := t.store.DeleteTodo(r.Context(), todoID)
	if err != nil {
		err := t.render.JSON(w, http.StatusInternalServerError, models.Error{
			Message: "Error delete todo",
		})
		if err != nil {
			log.Error().Err(err)
		}
		return
	}

	if count == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Debug().Msg(fmt.Sprint(count, " rows deleted for ", todoID))

	w.WriteHeader(200)
}

//HandlePost todo
func (t *TodoHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Error().Msg(fmt.Sprint("Failed to decode todo body: ", r.Body))
		err := t.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Error decoding body",
		})
		if err != nil {
			log.Error().Err(err)
		}
		return
	}

	id, err := t.store.PostTodo(r.Context(), todo)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprint("Failed to insert todo record: ", r.Body))
		err := t.render.JSON(w, http.StatusInternalServerError, models.Error{
			Message: "Error inserting record to database",
		})
		if err != nil {
			log.Error().Err(err)
		}
	}

	err = t.render.JSON(w, http.StatusOK, id)
	if err != nil {
		log.Error().Err(err)
	}
}
