package store

import (

	"database/sql"
	"testing"

	"github.com/alexsniffin/go-api-example/internal/api/config"

	"github.com/DATA-DOG/go-sqlmock"
)

type TestSQLClient struct {
	db   *sql.DB
}

func (t TestSQLClient) GetConnection() *sql.DB {
	return t.db
}

func (t TestSQLClient) CreateConnection(config *config.Config) (*sql.DB, error) {
	return t.db, nil
}

func (t TestSQLClient) Shutdown() error {
	return nil
}

func TestGetTodoValid(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	todoStore := TodoStore{
		sqlClient: TestSQLClient{
			db:   db,
		},
	}

	mock.ExpectQuery(`SELECT (.+) FROM todo`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "todo", "created_on"}).AddRow(1, "test", "time"))

	result, err := todoStore.GetTodo(1)
	if err != nil {
		t.Error(err)
	}
	if result.ID != 1 && result.Todo != "test" && result.CreatedOn != "time" {
		t.Error("Unexpected result")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
