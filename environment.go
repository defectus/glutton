package glutton

import "github.com/gin-gonic/gin"

type Settings struct {
	Host         string `env:"HOST" default:"0.0.0.0"`
	Port         string `env:"PORT" default:"4354"`
	OutputFolder string `env:"OUTPUT_FOLDER" default:"glutton"`
	BaseName     string `env:"BASE_NAME" default:"glutton_%d"`
	OutputDB     string `env:"OUTPUT_DB"`
	Debug        string `env:"DEBUG" default:"true"`
}

type Env struct {
	Settings *Settings
	Notifier PayloadNotifier
	Saver    PayloadSaver
	Parser   PayloadParser
	Server   *gin.Engine
}
