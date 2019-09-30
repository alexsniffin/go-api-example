package handlers

import (
	"fmt"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexsniffin/go-api-example/internal/api/clients"
	"github.com/alexsniffin/go-api-example/internal/api/models"
	"github.com/alexsniffin/go-api-example/internal/api/store"

	"github.com/go-chi/chi"
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
	store  *store.TodoStore
}

//NewTodoHandler todo
func NewTodoHandler(render *render.Render, sqlClient clients.SQLClient) *TodoHandler {
	store := store.NewTodoStore(sqlClient)

	return &TodoHandler{
		render: render,
		store: store,
	}
}

//HandleGet todo
func (t *TodoHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprint("Failed to decode todoID: " + chi.URLParam(r, "id")))
		t.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Error decoding TodoID",
		})
      	return
	}

	todo, err := t.store.GetTodo(todoID)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	t.render.JSON(w, http.StatusOK, todo)
	return
}

//HandleDelete todo
func (t *TodoHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprint("Failed to decode todoID: " + chi.URLParam(r, "id")))
		t.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Error decoding TodoID",
		})
      	return
	}

	count, err := t.store.DeleteTodo(todoID)
	if err != nil {
		t.render.JSON(w, http.StatusInternalServerError, models.Error{
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

//HandlePost todo
func (t *TodoHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
    var todo models.Todo
    err := json.NewDecoder(r.Body).Decode(&todo)
    if err != nil {
		log.Error().Msg(fmt.Sprint("Failed to decode todo body: ", r.Body))
		t.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Error decoding body",
		})
		return
	}
	
	id, err := t.store.PostTodo(todo)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprint("Failed to insert todo record: ", r.Body))
		t.render.JSON(w, http.StatusInternalServerError, models.Error{
			Message: "Error inserting record to database",
		})
	}

	t.render.JSON(w, http.StatusOK, id)
	return
}