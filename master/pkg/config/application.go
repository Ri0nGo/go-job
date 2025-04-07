package config

var App *Application

type Application struct {
	Server Server
	Data   Data
	MySQL  MySQL
}

type Server struct {
	Port uint16
	Name string
	Ip   string
}

type Data struct {
	UploadJobDir string `mapstructure:"upload_job_dir"`
}

type MySQL struct {
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Database    string `yaml:"database"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	MaxIdleConn int    `yaml:"max_idle_conn"`
	MaxOpenConn int    `yaml:"max_open_conn"`
	ShowSQL     bool   `yaml:"show_sql"`
}
