package server

import (
	"fmt"
	"net/http"
	"sync"
	"syscall"
	"time"
	"os"
	"os/signal"

	"github.com/alexsniffin/go-api-example/internal/api/config"
	"github.com/alexsniffin/go-api-example/internal/api/store"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
	"github.com/unrolled/render"
	"golang.org/x/net/context"
)

var shutdownOnce sync.Once

//Server todo
type Server struct {
	environment string
	httpServer  *http.Server
	router      *chi.Mux
	render      *render.Render
	config      *config.Config
	postgresDb  *store.Postgres
}

//NewServer todo
func NewServer(environment string) *Server {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	return &Server {
		environment: environment,
		router: r,
		render: render.New(),
	}
}

//Start todo
func (s *Server) Start() {
	s.InitDependencies()
	s.InitRouting()

	go s.InitHTTPServer() // Run server on a seperate thread to not block input signals

	stop := make(chan os.Signal, 1)

	// If you’re using kubernetes, note that it sends SIGTERM signal to its pods for shutting down. Interrupt is normally sent.
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	stopped := <- stop

	switch stopped.String() {
	case "SIGTERM", "interrupt":
		log.Info().Msg(stopped.String() + " signal received, attempting to gracefully shutdown")
		s.Shutdown()
	default:
		log.Error().Msg(stopped.String() + " signal received, attempting to gracefully shutdown")
		s.Shutdown()
	}
}

//Shutdown todo
func (s *Server) Shutdown() {
	shutdownOnce.Do(func() {
		if s.httpServer != nil {
			ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
			err := s.httpServer.Shutdown(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Failed to shutdown http server gracefully")
			} else {
				log.Info().Msg("Shutdown http server gracefully")
				s.httpServer = nil
			}
		}
		if s.postgresDb != nil {
			err := s.postgresDb.Shutdown()
			if err != nil {
				log.Error().Err(err).Msg("Failed to shutdown postgres gracefully")
			} else {
				log.Info().Msg("Shutdown postgres gracefully")
				s.postgresDb = nil
			}
		}
	})
}

//InitHTTPServer todo
func (s *Server) InitHTTPServer() {
	log.Info().Msg(fmt.Sprint("Running server on 0.0.0.0:", s.config.Cfg.Server.Port))
	s.httpServer = &http.Server{Addr: fmt.Sprint(":", s.config.Cfg.Server.Port), Handler: s.router}

	err := s.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Panic().Err(err).Msg("Http server stopped unexpected")
		s.Shutdown()
	} else {
		log.Info().Msg("Http server stopped")
	}
}

//InitDependencies todo
func (s *Server) InitDependencies() {
	config := config.NewConfig("config")
	postgresDb := store.NewPostgres(config, s.environment)

	s.config = config
	s.postgresDb = postgresDb
}

//InitRouting todo
func (s *Server) InitRouting() {
	s.todoRoutes()
}