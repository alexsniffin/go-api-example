package logger

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/alexsniffin/go-starter/internal/todo-api/models"
)

// Creates a zerolog logger
func NewLogger(cfg models.Config) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.Logger.Level)
	if err != nil {
		return zerolog.Logger{}, err
	}

	logger := zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger()
	if cfg.Environment == "localhost" {
		logger = logger.Output(zerolog.ConsoleWriter{
			Out: os.Stderr,
		})
	}

	return logger, nil
}
