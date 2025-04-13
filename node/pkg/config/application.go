package config

var App *Application

type Application struct {
	Server Server
	Data   Data
	Master Master
}

type Server struct {
	Port uint16
	Name string
	Ip   string
}

type Data struct {
	UploadJobDir string `mapstructure:"upload_job_dir"`
}

type Master struct {
	Address string `mapstructure:"address"`
	Key     string `mapstructure:"key"`
}
