package todo

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/unrolled/render"

	"github.com/alexsniffin/go-starter/internal/todo-api/models"
	"github.com/alexsniffin/go-starter/mocks"
)

func initTodoHandler() (Handler, *mocks.TodoStore) {
	todoStoreMock := mocks.TodoStore{}
	logger := zerolog.New(os.Stdout)
	todoHandler := Handler{
		logger: logger,
		render: render.New(),
		store:  &todoStoreMock,
	}
	return todoHandler, &todoStoreMock
}

func TestHealthCheckHandler(t *testing.T) {
	t.Run("foundTodo", func(t *testing.T) {
		todoHandler, todoStoreMock := initTodoHandler()
		id := 1
		todoStoreMock.On("GetTodo", mock.Anything, id).Return(models.Todo{
			ID:   1,
			Todo: "test",
		}, true, nil)

		req, err := http.NewRequest("GET", "/todo", nil)
		if err != nil {
			t.Fatal(err)
		}

		q := req.URL.Query()
		q.Add("id", strconv.Itoa(id))
		req.URL.RawQuery = q.Encode()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.HandleGet)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("unexpected status code: got %v want %v", status, http.StatusOK)
			t.FailNow()
		}

		expected := `{"id":1,"todo":"test","created_on":"0001-01-01T00:00:00Z"}`
		if rr.Body.String() != expected {
			t.Errorf("unexpected body: got %v want %v", rr.Body.String(), expected)
			t.FailNow()
		}

		todoStoreMock.AssertNumberOfCalls(t, "GetTodo", 1)
		todoStoreMock.AssertExpectations(t)
	})

	t.Run("noContent", func(t *testing.T) {
		todoHandler, todoStoreMock := initTodoHandler()
		id := 1
		todoStoreMock.On("GetTodo", mock.Anything, id).Return(models.Todo{}, false, nil)

		req, err := http.NewRequest("GET", "/todo", nil)
		if err != nil {
			t.Fatal(err)
		}

		q := req.URL.Query()
		q.Add("id", strconv.Itoa(id))
		req.URL.RawQuery = q.Encode()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.HandleGet)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("unexpected status code: got %v want %v", status, http.StatusNoContent)
			t.FailNow()
		}

		expected := ``
		if rr.Body.String() != expected {
			t.Errorf("unexpected body: got %v want %v", rr.Body.String(), expected)
			t.FailNow()
		}

		todoStoreMock.AssertNumberOfCalls(t, "GetTodo", 1)
		todoStoreMock.AssertExpectations(t)
	})

	t.Run("badParameter", func(t *testing.T) {
		todoHandler, _ := initTodoHandler()
		id := "bad"

		req, err := http.NewRequest("GET", "/todo", nil)
		if err != nil {
			t.Fatal(err)
		}

		q := req.URL.Query()
		q.Add("id", id)
		req.URL.RawQuery = q.Encode()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.HandleGet)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("unexpected status code: got %v want %v", status, http.StatusBadRequest)
			t.FailNow()
		}

		expected := `{"message":"id must be an integer"}`
		if rr.Body.String() != expected {
			t.Errorf("unexpected body: got %v want %v", rr.Body.String(), expected)
			t.Fail()
		}
	})
}
