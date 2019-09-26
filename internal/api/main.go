package main

import (
	"errors"
	"os"

	"github.com/alexsniffin/go-api-example/internal/api/server"
	
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	environment, exists := os.LookupEnv("GO_API_EXAMPLE_ENVRIONMENT")
	if !exists {
		panic(errors.New("Failed to initialize application, missing GO_API_EXAMPLE_ENVRIONMENT environment variable"))
	}
	
	if environment == "local" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Setting up server instance")
	server := server.NewServer(environment)

	server.Start()
}