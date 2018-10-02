package glutton

import "github.com/gin-gonic/gin"

type Settings struct {
	Host         string `env:"HOST"`
	Port         string `env:"PORT"`
	OutputFolder string `env:"OUTPUT_FOLDER"`
	OutputDB     string `env:"OUTPUT_DB"`
}

type Env struct {
	Settings *Settings
	Notifier PayloadNotifier
	Saver    PayloadSaver
	Parser   PayloadParser
	Server   *gin.Engine
}

var DefaultSettings = &Settings{}