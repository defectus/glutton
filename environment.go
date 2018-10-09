package glutton

import "github.com/gin-gonic/gin"

type Settings struct {
	Name         string `env:"NAME" default:"your friendly glutton"`
	Host         string `env:"HOST" default:"0.0.0.0"`
	Port         string `env:"PORT" default:"4354"`
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
}

type Env struct {
	Settings *Settings
	Notifier PayloadNotifier
	Saver    PayloadSaver
	Parser   PayloadParser
	Server   *gin.Engine
}
