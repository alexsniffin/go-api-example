package postgres

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/alexsniffin/go-starter/internal/todo-api/models"
)

type DatabaseClient interface {
	GetConnection() *pg.DB
	Shutdown() error
}

type Client struct {
	db *pg.DB
}

// Creates a postgres Client
func NewClient(logger zerolog.Logger, cfg models.Database) (Client, error) {
	db := pg.Connect(&pg.Options{
		User:     cfg.User,
		Addr:     fmt.Sprint(cfg.Host, ":", cfg.Port),
		Password: cfg.Password,
		Database: cfg.DbName,
		PoolSize: 20,
	})

	if cfg.CreateTable {
		err := db.CreateTable((*models.Todo)(nil), &orm.CreateTableOptions{
			Temp:          false,
			IfNotExists:   false,
			Varchar:       0,
			FKConstraints: false,
		})
		if err != nil {
			if err.Error()[:12] != "ERROR #42P07" {
				return Client{}, errors.Wrap(err, "failed to create todo table")
			}
		}
	}

	for i := 0; i < len(cfg.Tables); i++ {
		var check interface{}
		check, err := db.Exec(`SELECT to_regclass(?)`, cfg.Tables[i])
		if err != nil {
			return Client{}, errors.Wrap(err, "failed to execute pg init sql check")
		}
		if check == nil {
			return Client{}, errors.New(fmt.Sprintf("missing required table on pg db: host=%s dbname=%s", cfg.Host, cfg.DbName))
		}
	}

	logger.Info().Msg("connected to pg")

	return Client{
		db: db,
	}, nil
}

// Return the connection
func (p *Client) GetConnection() *pg.DB {
	return p.db
}

// Signals a shutdown to the client
func (p *Client) Shutdown() error {
	err := p.db.Close()
	if err != nil {
		return err
	}

	return nil
}
