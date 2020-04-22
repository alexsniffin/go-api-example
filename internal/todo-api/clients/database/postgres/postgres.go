package postgres

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/rs/zerolog"

	"github.com/alexsniffin/go-starter/internal/todo-api/models"
)

type Client struct {
	db *pg.DB
}

func NewClient(logger zerolog.Logger, cfg models.Database) (Client, error) {
	db := pg.Connect(&pg.Options{
		User:     cfg.User,
		Addr:     fmt.Sprint(cfg.Host, ":", cfg.Port),
		Password: cfg.Password,
		Database: cfg.DbName,
		PoolSize: 20,
	})

	return Client{
		db: db,
	}, nil
}

func (p *Client) GetConnection() *pg.DB {
	return p.db
}

func (p *Client) Shutdown() error {
	err := p.db.Close()
	if err != nil {
		return err
	}

	return nil
}
