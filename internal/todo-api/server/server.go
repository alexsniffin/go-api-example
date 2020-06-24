package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	negronimiddleware "github.com/slok/go-http-metrics/middleware/negroni"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"

	"github.com/alexsniffin/go-starter/internal/todo-api/clients/postgres"
	"github.com/alexsniffin/go-starter/internal/todo-api/handlers/logging"
	todoHandler "github.com/alexsniffin/go-starter/internal/todo-api/handlers/todo"
	"github.com/alexsniffin/go-starter/internal/todo-api/models"
	"github.com/alexsniffin/go-starter/internal/todo-api/router"
	"github.com/alexsniffin/go-starter/internal/todo-api/store/todo"
)

type Server struct {
	cfg    models.Config
	logger zerolog.Logger

	httpServer *http.Server
	pgClient   postgres.Client

	shutdown sync.Once
}

// Creates a new server instance with dependencies
func NewServer(cfg models.Config, logger zerolog.Logger) *Server {
	// set up pg client
	newPgClient, err := postgres.NewClient(logger, cfg.Database)
	if err != nil {
		logger.Panic().Caller().Err(err).Msg("failed to initialize pg client")
	}

	// set up store and handler
	newTodoStore := todo.NewStore(logger, newPgClient)
	newTodoHandler := todoHandler.NewHandler(logger, render.New(), newTodoStore)

	// set up router and middleware
	n := negroni.New()
	n.Use(negronimiddleware.Handler("", middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})))
	n.UseHandler(logging.NewHandler(logger))
	n.UseHandler(router.NewRouter(newTodoHandler))

	newHttpServer := &http.Server{
		Addr:    fmt.Sprint(":", cfg.HttpServer.Port),
		Handler: n,
	}

	return &Server{
		cfg:    cfg,
		logger: logger,

		httpServer: newHttpServer,
		pgClient:   newPgClient,
	}
}

// Starts all asynchronous processes the server
func (s *Server) Start() {
	go s.StartHTTPServer()
}

// Signals shutdown across the server
func (s *Server) Shutdown() {
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

		close(graceful)
	})
}

// Starts HTTP server
func (s *Server) StartHTTPServer() {
	s.logger.Info().Msg(fmt.Sprint("running server on 0.0.0.0:", s.cfg.HttpServer.Port))

	err := s.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		s.logger.Panic().Caller().Err(err).Msg("http server stopped unexpected")
		s.Shutdown()
	} else {
		s.logger.Info().Msg("http server stopped")
	}
}
