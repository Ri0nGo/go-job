package config

var App *Application

type Application struct {
	Server  Server
	Data    Data
	MySQL   MySQL
	Redis   Redis
	SMTP    map[string]SMTPProvider   `mapstructure:"smtp"`
	OAuth2  map[string]OAuth2Provider `mapstructure:"oauth2"`
	Metrics Metrics
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

type Redis struct {
	Auth string `mapstructure:"auth"`
	Addr string `mapstructure:"addr"`
	DB   int    `mapstructure:"db"`
}

//type SMTP struct {
//	Provider map[string]SMTPProvider `mapstructure:"smtp"`
//}

type SMTPProvider struct {
	Sender   string `mapstructure:"sender"`
	Key      string `mapstructure:"key"`
	SMTPHost string `mapstructure:"smtp_host"`
	SMTPPort int    `mapstructure:"smtp_port"`
}

type Metrics struct {
	Node NodeMetric `mapstructure:"node"`
}

type NodeMetric struct {
	Interval int `mapstructure:"interval"`
	Timeout  int `mapstructure:"timeout"`
}

type OAuth2Provider struct {
	ClientID         string `mapstructure:"client_id"`     // 泛指ID，可以是APPID
	ClientSecret     string `mapstructure:"client_secret"` // 泛指Secret, 可以是APPKey
	RedirectURL      string `mapstructure:"redirect_url"`
	RedirectFrontUrl string `mapstructure:"redirect_front_url"`
}
