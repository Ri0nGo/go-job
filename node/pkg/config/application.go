package config

var App *Application

type Application struct {
	Server Server
}

type Server struct {
	Port uint16
	Name string
	Ip   string
}
