package models

import (
	"time"
)

// Model
type Todo struct {
	tableName struct{}  `pg:"todo"` // nolint:structcheck,unused
	ID        int       `json:"id" pg:"id,pk"`
	Todo      string    `json:"todo" pg:"todo"`
	CreatedOn time.Time `json:"created_on" pg:"created_on"`
}

// Response model to POST
type TodoPostResponse struct {
	ID int `json:"id"`
}
