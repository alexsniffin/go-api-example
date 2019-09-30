package clients

import (
	"database/sql"
	"fmt"

	"github.com/alexsniffin/go-api-example/internal/api/config"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

//PostgresClient todo
type PostgresClient struct {
	connection *sql.DB
}

//NewPostgresClient todo
func NewPostgresClient(config *config.Config, tables []string) *PostgresClient {
	postgresClient := &PostgresClient{}

	db, err := postgresClient.CreateConnection(config)
	if err != nil {
		log.Panic().Msg("Failed to initialize postgres connection")
	}

	for i := 0; i < len(tables); i++ {
		var check interface{}
		err = db.QueryRow("SELECT to_regclass($1)", tables[i]).Scan(&check)
		if err != nil {
			log.Panic().Err(err).Msgf("Failed to check if %s table exists", tables[i])
		}

		if check == nil {
			log.Panic().Err(err).Msg(fmt.Sprintf("Missing required table on postgresdb: host=%s dbname=%s", config.Cfg.Database.Host, config.Cfg.Database.DbName))
		} else {
			log.Info().Msg(fmt.Sprintf("Existing table found on postgresdb: host=%s dbname=%s", config.Cfg.Database.Host, config.Cfg.Database.DbName))
		}
	}

	postgresClient.connection = db

	return postgresClient
}

//GetConnection todo
func (p *PostgresClient) GetConnection() *sql.DB {
	return p.connection
}

//CreateConnection todo
func (p *PostgresClient) CreateConnection(config *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%d sslmode=disable",
		config.Cfg.Database.User,
		config.Cfg.Database.DbName,
		config.Cfg.Database.Password,
		config.Cfg.Database.Host,
		config.Cfg.Database.Port,
	)
	return sql.Open("postgres", connStr)
}

//Shutdown todo
func (p *PostgresClient) Shutdown() error {
	err := p.connection.Close()
	if err != nil {
		return err
	}

	return nil
}
