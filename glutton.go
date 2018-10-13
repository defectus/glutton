package glutton

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Run is the entry point to Glutton.
func Run() error {
	env := createEnvironment(createConfiguration(new(Configuration)))
	if env.Configuration.Debug {
		log.Printf("current settings: %+v", env.Configuration)
	}
	closing := make(chan struct{})
	hookOnExit(closing)
	go env.Server.Run(env.Configuration.Host + ":" + env.Configuration.Port)
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

// initializeRoutes does the basic stuff needed to create a router.
func initializeRoutes(router *gin.Engine, env *Env) *gin.RouterGroup {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "HEAD", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	v1 := router.Group("v1")
	glutton := v1.Group("glutton")
	return glutton
}

func renderError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error(), "detail": err})
}
