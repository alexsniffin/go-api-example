package store

import (
	"github.com/alexsniffin/go-api-example/internal/api/models"
)

//GetTodo todo
func (p *Postgres) GetTodo(id int) (models.Todo, error) {
	var result models.Todo
	err := p.connection.QueryRow(`SELECT * FROM todo WHERE id = $1`, id).Scan(&result.ID, &result.Todo, &result.CreatedOn)
	if err != nil {
		return result, err
	}
	return result, nil
}

//DeleteTodo todo
func (p *Postgres) DeleteTodo(id int) (int64, error) {
	res, err := p.connection.Exec(`DELETE FROM todo WHERE id = $1`, id)
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
func (p *Postgres) PostTodo(todo models.Todo) (int, error) {
	var id int
	err := p.connection.QueryRow(`INSERT INTO todo(todo, created_on) VALUES($1, current_timestamp) RETURNING id`, todo.Todo).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}