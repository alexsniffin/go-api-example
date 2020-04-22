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
	store  todo.Store
}

func NewHandler(logger zerolog.Logger, render *render.Render, store todo.Store) Handler {
	return Handler{
		logger: logger,

		render: render,
		store:  store,
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	todoIDStr := r.URL.Query().Get("id")
	if todoIDStr == "" {
		h.logger.Error().Caller().Msg("missing id in request")
		err := h.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Missing query parameter: id",
		})
		if err != nil {
			h.logger.Error().Err(err)
		}
		return
	}

	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		err := h.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "id must be an integer",
		})
		if err != nil {
			h.logger.Error().Err(err).Send()
		}
		return
	}

	todoResult, err := h.store.GetTodo(r.Context(), todoID)
	if err != nil {

		h.logger.Error().Caller().Err(err).Send()

		err := h.render.JSON(w, http.StatusInternalServerError, models.Error{
			Message: "Error retrieving record",
		})
		if err != nil {
			h.logger.Error().Caller().Err(err)
		}
		return
	}

	if todoResult.ID == 0 && todoResult.Todo == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = h.render.JSON(w, http.StatusOK, todoResult)
	if err != nil {
		h.logger.Error().Caller().Err(err)
	}
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	todoIDStr := r.URL.Query().Get("id")
	if todoIDStr == "" {
		h.logger.Error().Msg("missing id in request")
		err := h.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Missing query parameter: id",
		})
		if err != nil {
			h.logger.Error().Caller().Err(err).Send()
		}
		return
	}

	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		err := h.render.JSON(w, http.StatusInternalServerError, models.Error{
			Message: "Error decoding id",
		})
		if err != nil {
			h.logger.Error().Caller().Err(err).Send()
		}
		return
	}

	count, err := h.store.DeleteTodo(r.Context(), todoID)
	if err != nil {
		err := h.render.JSON(w, http.StatusInternalServerError, models.Error{
			Message: "Internal server error with request",
		})
		if err != nil {
			h.logger.Error().Caller().Err(err).Send()
		}
		return
	}

	if count == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	h.logger.Debug().Caller().Msg(fmt.Sprint(count, " rows deleted for ", todoID))

	w.WriteHeader(200)
}

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var todoRequest models.Todo
	err := unmarshalRequestBody(r, &todoRequest)
	if err != nil {
		h.logger.Error().Caller().Msgf("failed to decode todo body: %v", todoRequest)
		err := h.render.JSON(w, http.StatusBadRequest, models.Error{
			Message: "Error decoding body",
		})
		if err != nil {
			h.logger.Error().Caller().Err(err).Send()
		}
		return
	}

	todoRequest.CreatedOn = time.Now()

	id, err := h.store.PostTodo(r.Context(), todoRequest)
	if err != nil {
		h.logger.Error().Caller().Err(err).Msgf("failed to insert todo record: %v", todoRequest)
		err := h.render.JSON(w, http.StatusInternalServerError, models.Error{
			Message: "Internal server error with request",
		})
		if err != nil {
			h.logger.Error().Caller().Err(err).Send()
		}
		return
	}

	err = h.render.JSON(w, http.StatusOK, models.TodoPostResponse{ID: id})
	if err != nil {
		h.logger.Error().Caller().Err(err).Send()
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
