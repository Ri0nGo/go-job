package config

var App *Application

type Application struct {
	Server Server
	Data   Data
	MySQL  MySQL
	SMTP   SMTP
}

type Server struct {
	Port uint16
	Name string
	Ip   string
	Key  string
}

type Data struct {
	UploadJobDir string `mapstructure:"upload_job_dir"`
}

type MySQL struct {
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	Database    string `mapstructure:"database"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"`
	MaxOpenConn int    `mapstructure:"max_open_conn"`
	ShowSQL     bool   `mapstructure:"show_sql"`
}

type SMTP struct {
	Sender   string `mapstructure:"sender"`
	Key      string `mapstructure:"key"`
	SMTPHost string `mapstructure:"smtp_host"`
	SMTPPort int    `mapstructure:"smtp_port"`
}
