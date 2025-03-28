package config

var App *Application

type Application struct {
	Config Config
}

type Config struct {
	Port uint16
	Name string
	Ip   string
}
