package models

//Configuration todo
type Configuration struct {
	Server   Server
	Database Database
}

//Server todo
type Server struct {
	Port int
}

//Database todo
type Database struct {
	Host     string
	Port     int
	User     string
	DbName   string
	Password string
	Tables   []string
}
