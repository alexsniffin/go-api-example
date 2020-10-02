package todo

import (
	"errors"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"

	"github.com/alexsniffin/go-api-starter/internal/todo-api/clients/postgres"
	"github.com/alexsniffin/go-api-starter/internal/todo-api/models"
)

type TodoStore interface {
	GetTodo(ctx context.Context, id int) (models.TodoItem, bool, error)
	DeleteTodo(ctx context.Context, id int) (int, error)
	PostTodo(ctx context.Context, todo models.TodoItem) (int, error)
}

type Store struct {
	pgClient postgres.DatabaseClient
}

// NewStore creates a new Store
func NewStore(pgClient postgres.Client) Store {
	return Store{
		pgClient: &pgClient,
	}
}

// GetTodo gets a TodoItem from the database
func (s *Store) GetTodo(ctx context.Context, id int) (models.TodoItem, bool, error) {
	log.Ctx(ctx).Debug().Caller().Caller().Msg("get db request for todo")

	var result models.TodoItem
	err := s.pgClient.GetConnection().
		Model(&result).
		Context(ctx).
		Where("id = ?", id).
		Select(&result)
	if err != nil {
		if err.Error() == "pg: no rows in result set" {
			return models.TodoItem{}, false, nil
		}
		log.Ctx(ctx).Error().Err(err).Caller().Msg("failed to get todo from db")
		return result, false, err
	}

	log.Ctx(ctx).Debug().Caller().Msg("todo found from db")
	return result, true, nil
}

// DeleteTodo deletes a TodoItem from the database
func (s *Store) DeleteTodo(ctx context.Context, id int) (int, error) {
	log.Ctx(ctx).Debug().Caller().Msg("delete db request for todo")

	result, err := s.pgClient.GetConnection().
		Model((*models.TodoItem)(nil)).
		Context(ctx).
		Where("id = ?", id).
		Delete()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Caller().Msg("failed to delete todo from db")
		return 0, err
	}

	log.Ctx(ctx).Debug().Caller().Msgf("todo deleted from db")
	return result.RowsAffected(), nil
}

// PostTodo posts a TodoItem to the database
func (s *Store) PostTodo(ctx context.Context, todo models.TodoItem) (int, error) {
	log.Ctx(ctx).Debug().Caller().Msg("insert db request for todo")

	result, err := s.pgClient.GetConnection().
		Model(&todo).
		Context(ctx).
		Returning("id").
		Insert(&todo)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Caller().Msg("failed to insert todo into db")
		return 0, err
	}
	if result.RowsAffected() == 0 {
		iErr := errors.New("failed to insert record")
		log.Ctx(ctx).Error().Err(iErr).Caller().Msg("failed to insert todo into db")
		return 0, iErr
	}

	return todo.ID, err
}
