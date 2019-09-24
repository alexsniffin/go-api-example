package store

import (
	"database/sql"
	"fmt"

	"github.com/alexsniffin/go-api-example/internal/api/config"
	
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

//Postgres todo
type Postgres struct {
	Connection *sql.DB
}

//InitPostgres todo
func InitPostgres(config *config.Config, environment string) *Postgres {
	postgres := &Postgres{}
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%d sslmode=disable", config.Cfg.Postgres.User, config.Cfg.Postgres.DbName, config.Cfg.Postgres.Password, config.Cfg.Postgres.Host, config.Cfg.Postgres.Port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic().Msg("Failed to initialize postgres connection")
	}

	postgres.Connection = db

	var check interface{}
	err = db.QueryRow("SELECT to_regclass('todo')").Scan(&check)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to check if todo table exists")
	}

	// Table creation should occur in the CD process but for simplcity it will create if running local
	if check == nil && environment == "local" {
		log.Info().Msg("Todo table not found, attempting to create")

		stmt, err := db.Prepare("CREATE TABLE todo (id serial PRIMARY KEY, todo VARCHAR(255), created_on TIMESTAMP NOT NULL)")
		if err != nil {
			log.Panic().Err(err).Msg("Failed to create todo table")
		}
		defer stmt.Close()
	
		_, err = stmt.Exec()
		if err != nil {
			log.Panic().Err(err).Msg("Failed to create todo table")
		}
	} else if check == nil && environment != "local" {
		log.Panic().Err(err).Msg(fmt.Sprintf("Missing required table(s) on postgresdb: host=%s dbname=%s", config.Cfg.Postgres.Host, config.Cfg.Postgres.DbName))
	} else {
		log.Info().Msg(fmt.Sprintf("Existing table(s) found on postgresdb: host=%s dbname=%s", config.Cfg.Postgres.Host, config.Cfg.Postgres.DbName))
	}

	return postgres
}