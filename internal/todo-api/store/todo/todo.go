package todo

import (
	"errors"

	"github.com/rs/zerolog"
	"golang.org/x/net/context"

	"github.com/alexsniffin/go-starter/internal/todo-api/clients/postgres"
	"github.com/alexsniffin/go-starter/internal/todo-api/models"
)

type TodoStore interface {
	GetTodo(ctx context.Context, id int) (models.Todo, bool, error)
	DeleteTodo(ctx context.Context, id int) (int, error)
	PostTodo(ctx context.Context, todo models.Todo) (int, error)
}

type Store struct {
	logger zerolog.Logger

	pgClient postgres.DatabaseClient
}

// Creates a new Store
func NewStore(logger zerolog.Logger, pgClient postgres.Client) Store {
	return Store{
		logger: logger,

		pgClient: &pgClient,
	}
}

// Gets a TodoItem from the database
func (s *Store) GetTodo(ctx context.Context, id int) (models.Todo, bool, error) {
	logFields := map[string]interface{}{
		"id": id,
	}
	s.logger.Debug().Caller().Fields(logFields).Caller().Msg("get db request for todo")

	var result models.Todo
	err := s.pgClient.GetConnection().
		Model(&result).
		Context(ctx).
		Where("id = ?", id).
		Select(&result)
	if err != nil {
		if err.Error() == "pg: no rows in result set" {
			return models.Todo{}, false, nil
		}
		s.logger.Error().Err(err).Fields(logFields).Caller().Msg("failed to get todo from db")
		return result, false, err
	}

	s.logger.Debug().Fields(logFields).Caller().Msg("todo found from db")
	return result, true, nil
}

// Deletes a TodoItem from the database
func (s *Store) DeleteTodo(ctx context.Context, id int) (int, error) {
	logFields := map[string]interface{}{
		"id": id,
	}
	s.logger.Debug().Caller().Fields(logFields).Msg("delete db request for todo")

	result, err := s.pgClient.GetConnection().
		Model((*models.Todo)(nil)).
		Context(ctx).
		Where("id = ?", id).
		Delete()
	if err != nil {
		s.logger.Error().Err(err).Fields(logFields).Caller().Msg("failed to delete todo from db")
		return 0, err
	}

	s.logger.Debug().Fields(logFields).Caller().Msgf("todo deleted from db")
	return result.RowsAffected(), nil
}

// Posts a TodoItem to the database
func (s *Store) PostTodo(ctx context.Context, todo models.Todo) (int, error) {
	logFields := map[string]interface{}{
		"id": todo.ID,
	}
	s.logger.Debug().Caller().Fields(logFields).Msg("insert db request for todo")

	result, err := s.pgClient.GetConnection().
		Model(&todo).
		Context(ctx).
		Returning("id").
		Insert(&todo)
	if err != nil {
		s.logger.Error().Err(err).Fields(logFields).Caller().Msg("failed to insert todo into db")
		return 0, err
	}
	if result.RowsAffected() == 0 {
		iErr := errors.New("failed to insert record")
		s.logger.Error().Err(iErr).Fields(logFields).Caller().Msg("failed to insert todo into db")
		return 0, iErr
	}

	return todo.ID, err
}
