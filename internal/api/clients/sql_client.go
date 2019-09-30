package clients

import (
	"database/sql"

	"github.com/alexsniffin/go-api-example/internal/api/config"
)

//SQLClient todo
type SQLClient interface {
	GetConnection() *sql.DB
	CreateConnection(config *config.Config) (*sql.DB, error)
	Shutdown() error
}
