package models

import (
	"github.com/alexsniffin/go-starter/pkg/models"
)

type Config struct {
	Environment string
	Logger      models.Logger
	HttpServer  HttpServer
	Database    Database
}

type HttpServer struct {
	Port int
}

type Database struct {
	Host         string
	Port         int
	User         string
	DbName       string
	Password     string
	SkipSSLCheck bool
	Tables       []string
}
