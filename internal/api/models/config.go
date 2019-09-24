package models

//Configuration todo
type Configuration struct {
	Server   Server
	Postgres Postgres
}

//Server todo
type Server struct {
	Port int
}

//Postgres todo
type Postgres struct {
	Host     string
	Port     int
	User     string
	DbName   string
	Password string
}