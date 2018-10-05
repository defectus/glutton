package glutton

import "github.com/gin-gonic/gin"

type Settings struct {
	Host         string `env:"HOST" default:"0.0.0.0"`
	Port         string `env:"PORT" default:"4354"`
	OutputFolder string `env:"OUTPUT_FOLDER" default:"glutton"`
	BaseName     string `env:"BASE_NAME" default:"glutton_%d"`
	OutputDB     string `env:"OUTPUT_DB"`
	Debug        bool   `env:"DEBUG" default:"true"`
	SMTPServer   string `env:"SMTP_SERVER" default:"smtp.gmail.com"`
	SMTPPort     string `env:"SMTP_PORT" default:"25"`
	SMTPUseTLS   bool   `env:"SMTP_USE_TLS" default:"true"`
	SMTPUser     string `env:"SMTP_USER"`
	SMTPPassword string `env:"SMTP_USER"`
}

type Env struct {
	Settings *Settings
	Notifier PayloadNotifier
	Saver    PayloadSaver
	Parser   PayloadParser
	Server   *gin.Engine
}
