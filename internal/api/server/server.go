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
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

var shutdownOnce sync.Once

//Instance todo
type Instance struct {
	environment string
	httpServer  *http.Server
	router      *chi.Mux
	config      *config.Config
	postgresDb  *store.Postgres
}

//NewInstance todo
func NewInstance(environment string) *Instance {
	s := &Instance{
		environment: environment,
		router: chi.NewRouter(),
	}

	return s
}

//Start todo
func (s *Instance) Start() {
	s.InitDependencies()

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
		log.Panic().Msg(stopped.String() + " signal received, attempting to gracefully shutdown")
		s.Shutdown()
	}
}

//Shutdown todo
func (s *Instance) Shutdown() {
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
			err := s.postgresDb.Connection.Close()
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
func (s *Instance) InitHTTPServer() {
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
func (s *Instance) InitDependencies() {
	config := config.InitConfig("config")
	postgresDb := store.InitPostgres(config, s.environment)

	s.config = config
	s.postgresDb = postgresDb
}