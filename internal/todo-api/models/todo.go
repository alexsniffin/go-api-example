package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// TodoItem model
type TodoItem struct {
	tableName struct{}  `pg:"todo"` // nolint:structcheck,unused
	ID        int       `json:"id" pg:"id,pk"`
	Todo      string    `json:"todo" pg:"todo"`
	CreatedOn time.Time `json:"created_on" pg:"created_on"`
}

// TodoPostResponse response model to POST
type TodoPostResponse struct {
	ID int `json:"id"`
}

// TodoPostRequest request model to POST
type TodoPostRequest struct {
	Todo string `json:"todo"`
}

func (tReq *TodoPostRequest) IsValid() error {
	return validation.ValidateStruct(tReq,
		validation.Field(&tReq.Todo, validation.Required),
	)
}
