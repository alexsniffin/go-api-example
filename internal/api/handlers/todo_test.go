package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexsniffin/go-api-example/internal/api/models"

	"github.com/unrolled/render"
)

type TestTodoStore struct {
	getTodoSuccess    models.Todo
	getTodoError      error
	deleteTodoSuccess int64
	deleteTodoError   error
	postTodoSuccess   int
	postTodoError     error
}

func (t *TestTodoStore) GetTodo(id int) (models.Todo, error) {
	return t.getTodoSuccess, t.getTodoError
}

func (t *TestTodoStore) DeleteTodo(id int) (int64, error) {
	return t.deleteTodoSuccess, t.deleteTodoError
}

func (t *TestTodoStore) PostTodo(todo models.Todo) (int, error) {
	return t.postTodoSuccess, t.postTodoError
}

func TestHealthCheckHandler_foundTodo(t *testing.T) {
	todoHandler := NewTodoHandler(render.New(), &TestTodoStore{
		getTodoSuccess: models.Todo{
			ID:        1,
			Todo:      "test",
			CreatedOn: "test",
		},
	})

	req, err := http.NewRequest("GET", "/todo", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("id", "1")
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.HandleGet)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":1,"todo":"test","created_on":"test"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHealthCheckHandler_noContent(t *testing.T) {
	todoHandler := NewTodoHandler(render.New(), &TestTodoStore{
		getTodoSuccess: models.Todo{},
		getTodoError: errors.New("Some error"),
	})

	req, err := http.NewRequest("GET", "/todo", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("id", "1")
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.HandleGet)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	expected := ``
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHealthCheckHandler_badParameter(t *testing.T) {
	todoHandler := NewTodoHandler(render.New(), &TestTodoStore{},)

	req, err := http.NewRequest("GET", "/todo", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("id", "bad")
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler.HandleGet)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"message":"id must be an integer"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
