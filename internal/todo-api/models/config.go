package models

import (
	"github.com/alexsniffin/go-api-starter/pkg/models"
)

type Config struct {
	Environment string
	Logger      models.Logger
	HTTPServer  HTTPServerConfig
	HTTPRouter  HTTPRouterConfig
	Database    DatabaseConfig
}

type HTTPServerConfig struct {
	Port int
}

type HTTPRouterConfig struct {
	TimeoutSec     int
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type DatabaseConfig struct {
	Host        string
	Port        int
	User        string
	DbName      string
	Password    string
	Tables      []string
	CreateTable bool
}
