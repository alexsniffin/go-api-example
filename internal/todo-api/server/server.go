package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/unrolled/render"

	"golang.org/x/net/context"

	"github.com/alexsniffin/go-starter/internal/todo-api/clients/database/postgres"
	todo2 "github.com/alexsniffin/go-starter/internal/todo-api/handlers/todo"
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

func NewServer(cfg models.Config, logger zerolog.Logger) *Server {
	newPgClient, err := postgres.NewClient(logger, cfg.Database)
	if err != nil {
		logger.Panic().Caller().Err(err).Send()
	}
	newTodoStore := todo.NewStore(logger, newPgClient)
	newTodoHandler := todo2.NewHandler(logger, render.New(), newTodoStore)
	newRouter := router.NewRouter(newTodoHandler)
	newHttpServer := &http.Server{
		Addr:    fmt.Sprint(":", cfg.HttpServer.Port),
		Handler: newRouter,
	}

	return &Server{
		cfg:    cfg,
		logger: logger,

		httpServer: newHttpServer,
		pgClient:   newPgClient,
	}
}

func (s *Server) Start() {
	go s.StartHTTPServer()
}

func (s *Server) Shutdown() {
	s.shutdown.Do(func() {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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
	})
}

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
