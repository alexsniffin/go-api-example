package server

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/unrolled/render"

	"github.com/alexsniffin/go-api-starter/internal/todo-api/clients/postgres"
	todoHandler "github.com/alexsniffin/go-api-starter/internal/todo-api/handlers/todo"
	"github.com/alexsniffin/go-api-starter/internal/todo-api/models"
	"github.com/alexsniffin/go-api-starter/internal/todo-api/processes/http"
	"github.com/alexsniffin/go-api-starter/internal/todo-api/router"
	"github.com/alexsniffin/go-api-starter/internal/todo-api/store/todo"
)

// Server handles the runtime of the application.
type Server struct {
	cfg    models.Config
	logger zerolog.Logger

	httpServer *http.Server
	pgClient   postgres.Client

	fatalErrCh chan error
	shutdown   sync.Once
}

// NewServer creates a new server instance with dependencies.
func NewServer(cfg models.Config, logger zerolog.Logger) *Server {
	// set up pg client
	newPgClient, err := postgres.NewClient(logger, cfg.Database)
	if err != nil {
		logger.Panic().Caller().Err(err).Msg("failed to initialize pg client")
	}

	// set up store and handler
	newTodoStore := todo.NewStore(newPgClient)
	newTodoHandler := todoHandler.NewHandler(logger, render.New(), newTodoStore)

	// set up router and HTTP server
	newRouter := router.NewRouter(cfg.HTTPRouter, logger, newTodoHandler)
	newHTTPServer := http.NewServer(cfg.HTTPServer, logger, newRouter)

	return &Server{
		cfg:        cfg,
		logger:     logger,
		httpServer: newHTTPServer,
		pgClient:   newPgClient,
		fatalErrCh: make(chan error),
	}
}

// Start invokes all asynchronous server processes.
func (s *Server) Start() {
	go s.httpServer.Start(s.fatalErrCh)

	for err := range s.fatalErrCh {
		if err != nil {
			s.logger.Error().Caller().Err(err).Msg("fatal error received from process")
			s.Shutdown(true)
		}
	}
}

// Shutdown signals the shutdown process across all processes in the server.
func (s *Server) Shutdown(fromErr bool) {
	s.shutdown.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		graceful := make(chan bool)

		// terminate the server if the deadline is reached, regardless of shutdown processes and returns exit status 2
		go func(graceful <-chan bool) {
			for {
				select {
				case <-ctx.Done():
					s.logger.Panic().Msg("shutdown deadline reached, terminating remaining processes ungracefully")
				case <-graceful:
					return
				}
			}
		}(graceful)

		// shutdown http server first to prevent new requests
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			s.logger.Error().Caller().Err(err).Msg("failed to shutdown http server gracefully")
		} else {
			s.logger.Info().Msg("shutdown http server gracefully")
		}

		err = s.pgClient.Shutdown()
		if err != nil {
			s.logger.Error().Caller().Err(err).Msg("failed to shutdown postgres gracefully")
		} else {
			s.logger.Info().Msg("shutdown postgres gracefully")
		}

		close(s.fatalErrCh)
		close(graceful)

		if fromErr {
			s.logger.Info().Msg("graceful shutdown succeeded, an error was detected, exiting with status code 1")
			os.Exit(1)
		}
	})
}
