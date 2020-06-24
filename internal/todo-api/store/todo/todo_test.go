package todo

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/alexsniffin/go-starter/internal/todo-api/models"
	"github.com/alexsniffin/go-starter/mocks"
)

func unexpected(t *testing.T, err error) {
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
		t.FailNow()
	}
}

func createPgContainer(t *testing.T, user, pass, dbName string) testcontainers.Container {
	req := testcontainers.ContainerRequest{
		Image:        "frodenas/postgresql",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USERNAME": user,
			"POSTGRES_PASSWORD": pass,
			"POSTGRES_DBNAME":   dbName,
		},
		AlwaysPullImage: true,
		WaitingFor:      wait.ForLog("LOG:  database system is ready to accept connections"),
	}
	container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	unexpected(t, errors.Wrap(err, "failed to create pg container"))

	time.Sleep(2 * time.Second) // db isn't ready even when log says it is
	return container
}

func initDb(t *testing.T) (*pg.DB, testcontainers.Container) {
	cfg := models.Database{
		Host:     "localhost",
		User:     "test",
		DbName:   "tododb",
		Password: "pass123",
	}
	pgContainer := createPgContainer(t, cfg.User, cfg.Password, cfg.DbName)

	port, err := nat.NewPort("tcp", "5432")
	unexpected(t, errors.Wrap(err, "failed to get nat port"))

	exposedPort, err := pgContainer.MappedPort(context.Background(), port)
	unexpected(t, errors.Wrap(err, "failed to get mapped port"))

	pgClient := pg.Connect(&pg.Options{
		User:      cfg.User,
		Addr:      fmt.Sprint(cfg.Host, ":", exposedPort.Port()),
		Password:  cfg.Password,
		Database:  cfg.DbName,
		PoolSize:  20,
		TLSConfig: nil,
	})

	err = pgClient.CreateTable((*models.Todo)(nil), &orm.CreateTableOptions{
		Temp:          false,
		IfNotExists:   false,
		Varchar:       0,
		FKConstraints: false,
	})
	unexpected(t, errors.Wrap(err, "failed to create table"))

	return pgClient, pgContainer
}

// Example test using testcontainers
func TestGetTodo_ValidEmptyResponse(t *testing.T) {
	t.Parallel()

	db, container := initDb(t)
	defer container.Terminate(context.Background())

	dbMock := &mocks.DatabaseClient{}
	todoStore := Store{
		logger:   zerolog.New(os.Stdout),
		pgClient: dbMock,
	}

	dbMock.On("GetConnection").Return(db)

	emptyTodo, found, err := todoStore.GetTodo(context.Background(), 0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if found {
		t.Errorf("unexpected result: %v", emptyTodo)
	}

	dbMock.AssertNumberOfCalls(t, "GetConnection", 1)
	dbMock.AssertExpectations(t)
}
