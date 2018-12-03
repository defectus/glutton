package iface

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

// Configuration is the root of all configuration settings.
type Configuration struct {
	Settings []Settings `yaml:"settings"`
	Debug    bool       `env:"DEBUG" yaml:"debug"`
	Host     string     `env:"HOST" default:"0.0.0.0" yaml:"host"`
	Port     string     `env:"PORT" default:"4354" yaml:"port"`
}

// Settings holds configuration of a single route.
type Settings struct {
	Name                string `env:"NAME" default:"your friendly glutton" yaml:"name"`
	Debug               bool   `env:"DEBUG" yaml:"debug"`
	URI                 string `env:"URI" default:"save" yaml:"uri"`
	Redirect            string `env:"REDIRECT" yaml:"redirect" `
	OutputFolder        string `env:"OUTPUT_FOLDER" default:"glutton" yaml:"output_folder"`
	BaseName            string `env:"BASE_NAME" default:"glutton_%d" yaml:"base_name"`
	SMTPServer          string `env:"SMTP_SERVER" default:"smtp.gmail.com" yaml:"smtp_server"`
	SMTPPort            string `env:"SMTP_PORT" default:"25" yaml:"smtp_port"`
	SMTPUseTLS          bool   `env:"SMTP_USE_TLS" default:"true" yaml:"smtp_use_tls"`
	SMTPFrom            string `env:"SMTP_FROM" yaml:"smtp_from"`
	SMTPPassword        string `env:"SMTP_PASSWORD" yaml:"smtp_password"`
	SMTPTo              string `env:"SMTP_TO" yaml:"smtp_to"`
	Parser              string `env:"PARSER" default:"SimpleParser" yaml:"parser"`
	Notifier            string `env:"NOTIFIER" default:"NilNotifier" yaml:"notifier"`
	Saver               string `env:"SAVER" default:"SimpleFileSystemSaver" yaml:"saver"`
	UseToken            bool   `env:"USE_TOKEN" default:"false" yaml:"use_token"`
	TokenKey            string `env:"TOKEN_KEY" yaml:"token_key"`
	SQLDriver           string `env:"SQL_DRIVER" default:"postgres" yaml:"sql_driver"`
	SQLayout            string `env:"SQL_LAYOUT" yaml:"sql_layout"`
	SQLConnectionString string `env:"SQL_CONNECTION_STRING" default:"postgres://root:root@localhost/postgres?sslmode=disable" yaml:"sql_connection_string"`
}

// Env holds references to almost all application resources.
type Env struct {
	Configuration *Configuration
	Notifiers     map[string]reflect.Type
	Savers        map[string]reflect.Type
	Parsers       map[string]reflect.Type
	Server        *gin.Engine
}
