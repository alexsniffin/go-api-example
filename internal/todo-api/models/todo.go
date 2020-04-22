package models

import (
	"time"
)

type Todo struct {
	tableName struct{}  `pg:"todo"`
	ID        int       `json:"id" pg:"id,pk"`
	Todo      string    `json:"todo" pg:"todo"`
	CreatedOn time.Time `json:"created_on" pg:"created_on"`
}

type TodoPostResponse struct {
	ID int `json:"id"`
}
