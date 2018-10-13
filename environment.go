package glutton

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

type Configuration struct {
	Settings []Settings
	Debug    bool
	Host     string `env:"HOST" default:"0.0.0.0"`
	Port     string `env:"PORT" default:"4354"`
}

// Settings holds configuration of a single route.
type Settings struct {
	Name         string `env:"NAME" default:"your friendly glutton"`
	URI          string `env:URI default:save`
	Redirect     string `env:REDIRECT`
	OutputFolder string `env:"OUTPUT_FOLDER" default:"glutton"`
	BaseName     string `env:"BASE_NAME" default:"glutton_%d"`
	OutputDB     string `env:"OUTPUT_DB"`
	Debug        bool   `env:"DEBUG" default:"true"`
	SMTPServer   string `env:"SMTP_SERVER" default:"smtp.gmail.com"`
	SMTPPort     string `env:"SMTP_PORT" default:"25"`
	SMTPUseTLS   bool   `env:"SMTP_USE_TLS" default:"true"`
	SMTPFrom     string `env:"SMTP_FROM"`
	SMTPPassword string `env:"SMTP_PASSWORD"`
	SMTPTo       string `env:"SMTP_TO"`
	Parser       string `env:"PARSER"`
	Notifier     string `env:"NOTIFIER"`
	Saver        string `env:"SAVER"`
}

// Env holds references to almost all application resources.
type Env struct {
	Configuration *Configuration
	Notifiers     map[string]reflect.Type
	Savers        map[string]reflect.Type
	Parsers       map[string]reflect.Type
	Server        *gin.Engine
}
