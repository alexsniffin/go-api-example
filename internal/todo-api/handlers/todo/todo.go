package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"github.com/alexsniffin/go-starter/internal/todo-api/models"
	"github.com/alexsniffin/go-starter/internal/todo-api/store/todo"

	"github.com/unrolled/render"
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
func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	todoIDStr := r.URL.Query().Get("id")
	if todoIDStr == "" {
		h.logger.Error().Caller().Msg("missing id in request")
		err := h.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Missing query parameter: id",
		})
		if err != nil {
			h.logger.Error().Caller().Err(err).Msg("failed to marshal json todo get response")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		jErr := h.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "id must be an integer",
		})
		if jErr != nil {
			h.logger.Error().Caller().Err(err).Msg("failed to marshal json todo get response")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	todoResult, found, err := h.store.GetTodo(r.Context(), todoID)
	if err != nil {
		h.logger.Error().Caller().Err(err).Send()

		jErr := h.render.JSON(w, http.StatusInternalServerError, models.Error{
			Message: "Error retrieving record",
		})
		if jErr != nil {
			h.logger.Error().Caller().Err(err).Msg("failed to marshal json todo get response")
			w.WriteHeader(http.StatusInternalServerError)
		}
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

// Handle HTTP Del for TodoItem
func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	todoIDStr := r.URL.Query().Get("id")
	if todoIDStr == "" {
		h.logger.Error().Msg("missing id in request")
		err := h.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Missing query parameter: id",
		})
		if err != nil {
			h.logger.Error().Caller().Err(err).Msg("failed to marshal json todo delete response")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		rErr := h.render.JSON(w, http.StatusInternalServerError, models.Error{
			Message: "Error decoding id",
		})
		if rErr != nil {
			h.logger.Error().Caller().Err(rErr).Msg("failed to marshal json todo delete response")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	count, err := h.store.DeleteTodo(r.Context(), todoID)
	if err != nil {
		err := h.render.JSON(w, http.StatusInternalServerError, models.Error{
			Message: "Internal server error with request",
		})
		if err != nil {
			h.logger.Error().Caller().Err(err).Msg("failed to marshal json todo delete response")
			w.WriteHeader(http.StatusInternalServerError)
		}
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
func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var todoRequest models.Todo
	err := unmarshalRequestBody(r, &todoRequest)
	if err != nil {
		h.logger.Error().Caller().Msgf("failed to decode todo body: %v", todoRequest)
		jErr := h.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Error decoding body",
		})
		if jErr != nil {
			h.logger.Error().Caller().Err(err).Msg("failed to marshal json response")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	todoRequest.CreatedOn = time.Now()

	id, err := h.store.PostTodo(r.Context(), todoRequest)
	if err != nil {
		h.logger.Error().Caller().Err(err).Msgf("failed to insert todo record: %v", todoRequest)
		jErr := h.render.JSON(w, http.StatusInternalServerError, models.Error{
			Message: "Internal server error with request",
		})
		if jErr != nil {
			h.logger.Error().Caller().Err(err).Msg("failed to marshal json response")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = h.render.JSON(w, http.StatusOK, models.TodoPostResponse{ID: id})
	if err != nil {
		h.logger.Error().Caller().Err(err).Msg("failed to marshal json response")
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

	err = req.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &output)
	if err != nil {
		return err
	}

	return nil
}
