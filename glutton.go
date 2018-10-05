package glutton

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() error {
	env := createEnvironment(createSettings(new(Settings)))
	if env.Settings.Debug {
		log.Printf("current settings: %+v", env.Settings)
	}
	closing := make(chan struct{})
	hookOnExit(closing)
	go env.Server.Run(env.Settings.Host + ":" + env.Settings.Port)
	<-closing
	return nil
}

// hookOnExit listens for signal SIGHUP and once received it closes the provided `closing` channel.
func hookOnExit(closing chan struct{}) {
	go func(closing chan struct{}) {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Kill, os.Interrupt)
		<-signals
		log.Println("Initiating shutdown ...")
		close(closing)
	}(closing)
}

func initializeRoutes(router *gin.Engine, env *Env) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "HEAD", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	v1 := router.Group("v1")
	glutton := v1.Group("glutton")
	glutton.POST("/save", savePayload(env))
}

func savePayload(env *Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := env.Parser.Parse(c.Request)
		if err != nil {

		}
		err = env.Notifier.Notify(payload)
		if err != nil {

		}
		err = env.Saver.Save(payload)
		if err != nil {

		}
		c.Status(http.StatusOK)
	}
}

func renderError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error(), "detail": err})
}
