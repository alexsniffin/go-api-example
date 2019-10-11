package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexsniffin/go-api-example/internal/api/models"

	"github.com/unrolled/render"
)

type TestTodoStore struct {}

func (t *TestTodoStore) GetTodo(id int) (models.Todo, error) {
	return models.Todo{
		ID: 1,
		Todo: "test",
		CreatedOn: "test",
	}, nil
}

func (t *TestTodoStore) DeleteTodo(id int) (int64, error) {
	return 1, nil
}

func (t *TestTodoStore) PostTodo(todo models.Todo) (int, error) {
	return 1, nil
}

func TestHealthCheckHandler(t *testing.T) {
	todoHandler := NewTodoHandler(render.New(), &TestTodoStore{})

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
