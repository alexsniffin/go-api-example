package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rs/zerolog"
	"github.com/unrolled/render"

	"github.com/alexsniffin/go-api-starter/internal/todo-api/models"
	"github.com/alexsniffin/go-api-starter/internal/todo-api/store/todo"
)

type Handler struct {
	logger zerolog.Logger

	render *render.Render
	store  todo.TodoStore
}

// Creates TodoItem handler
func NewHandler(logger zerolog.Logger, render *render.Render, store todo.Store) Handler {
	return Handler{
		logger: logger,

		render: render,
		store:  &store,
	}
}

// Handle HTTP Get for TodoItem
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	todoIDStr := chi.URLParam(r, "id")
	err := validation.Validate(todoIDStr, validation.Required, is.Int.Error("id must be an integer"))
	if err != nil {
		h.logger.Debug().Caller().Msg("missing id in request")
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		h.logger.Error().Caller().Err(err).Msg("failed to decode todoID")
		h.writeErrorResponse(w, http.StatusInternalServerError, "Error decoding id value")
		return
	}

	todoResult, found, err := h.store.GetTodo(r.Context(), todoID)
	if err != nil {
		h.logger.Error().Caller().Err(err).Msg("failed to get todoItem")
		h.writeErrorResponse(w, http.StatusBadRequest, "Error retrieving record")
		return
	}
	if !found {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = h.render.JSON(w, http.StatusOK, todoResult)
	if err != nil {
		h.logger.Error().Caller().Err(err).Msg("failed to marshal json todo get response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Handle HTTP Delete for TodoItem
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	todoIDStr := chi.URLParam(r, "id")
	err := validation.Validate(todoIDStr, validation.Required, is.Int.Error("id must be an integer"))
	if err != nil {
		h.logger.Debug().Caller().Msg("missing id in request")
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		h.logger.Error().Caller().Err(err).Msg("failed to decode todoID")
		h.writeErrorResponse(w, http.StatusInternalServerError, "Error decoding id value")
		return
	}

	count, err := h.store.DeleteTodo(r.Context(), todoID)
	if err != nil {
		h.logger.Error().Caller().Err(err).Msg("failed to delete todo")
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error with request")
		return
	}
	if count == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	h.logger.Debug().Caller().Msg(fmt.Sprint(count, " rows deleted for ", todoID))

	w.WriteHeader(http.StatusOK)
}

// Handle HTTP Post for TodoItem
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	var todoRequest models.TodoPostRequest
	if err := unmarshalRequestBody(r, &todoRequest); err != nil {
		h.logger.Error().Caller().Err(err).Msgf("failed to decode todo body: %v", todoRequest)
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := todoRequest.IsValid(); err != nil {
		h.logger.Debug().Caller().Err(err).Msg("invalid post")
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.store.PostTodo(r.Context(), models.TodoItem{
		Todo:      todoRequest.Todo,
		CreatedOn: time.Now(),
	})
	if err != nil {
		h.logger.Error().Caller().Err(err).Msgf("failed to insert todo record: %v", todoRequest)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Internal server error with request")
		return
	}

	if err = h.render.JSON(w, http.StatusOK, models.TodoPostResponse{ID: id}); err != nil {
		h.logger.Error().Caller().Err(err).Msg("failed to marshal json response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) writeErrorResponse(w http.ResponseWriter, statusCode int, responseMessage string) {
	if rErr := h.render.JSON(w, statusCode, models.Error{
		Message: responseMessage,
	}); rErr != nil {
		h.logger.Error().Caller().Err(rErr).Msg("failed to marshal json response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func unmarshalRequestBody(req *http.Request, output interface{}) error {
	if req.Body == nil {
		return errors.New("invalid body in request")
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	if err = req.Body.Close(); err != nil {
		return err
	}
	if err = json.Unmarshal(body, &output); err != nil {
		return err
	}

	return nil
}
