package store

import (
	"github.com/alexsniffin/go-api-example/internal/api/clients/database"
	"github.com/alexsniffin/go-api-example/internal/api/models"

	"golang.org/x/net/context"
)

//Todo todo
type Todo interface {
	GetTodo(ctx context.Context, id int) (models.Todo, error)
	DeleteTodo(ctx context.Context, id int) (int64, error)
	PostTodo(ctx context.Context, todo models.Todo) (int, error)
}

//TodoStore todo
type TodoStore struct {
	sqlClient database.SQLClient
}

//NewTodoStore todo
func NewTodoStore(sqlClient database.SQLClient) *TodoStore {
	return &TodoStore{
		sqlClient: sqlClient,
	}
}

//GetTodo todo
func (t *TodoStore) GetTodo(ctx context.Context, id int) (models.Todo, error) {
	var result models.Todo

	err := t.sqlClient.GetConnection().QueryRowContext(ctx, `SELECT * FROM todo WHERE id = $1`, id).Scan(&result.ID, &result.Todo, &result.CreatedOn)
	if err != nil {
		return result, err
	}
	return result, nil
}

//DeleteTodo todo
func (t *TodoStore) DeleteTodo(ctx context.Context, id int) (int64, error) {
	res, err := t.sqlClient.GetConnection().ExecContext(ctx, `DELETE FROM todo WHERE id = $1`, id)
	if err != nil {
		return 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, err
}

//PostTodo todo
func (t *TodoStore) PostTodo(ctx context.Context, todo models.Todo) (int, error) {
	var id int
	err := t.sqlClient.GetConnection().QueryRowContext(ctx, `INSERT INTO todo(todo, created_on) VALUES($1, current_timestamp) RETURNING id`, todo.Todo).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}
