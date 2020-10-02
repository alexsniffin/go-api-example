package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexsniffin/go-api-starter/internal/todo-api/models"
	"github.com/alexsniffin/go-api-starter/internal/todo-api/server"
	"github.com/alexsniffin/go-api-starter/pkg/config"
	"github.com/alexsniffin/go-api-starter/pkg/logger"
)

const (
	configName = "todo-api"
	prefix     = "TODO"
)

// Entry point to the application.
//
// Exit status codes:
//    * 0 - success
//    * 1 - from fatal internal error
//    * 2 - shutdown timeout
func main() {
	newCfg := models.Config{}
	err := config.NewConfig(configName, prefix, &newCfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	newLogger, err := logger.NewLogger(newCfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	newLogger.Info().Msg("setting up todo api service")
	newServer := server.NewServer(newCfg, newLogger)
	go newServer.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	stopped := <-stop
	newLogger.Info().Msg(stopped.String() + " signal received, attempting to gracefully shutdown")
	newServer.Shutdown(false)

	newLogger.Info().Msg("exiting todo api service")
}
