package common

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/defectus/glutton/pkg/iface"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Run is the entry point to Glutton.
func Run() error {
	var (
		yamlConfiguration []byte
		err               error
	)
	// first see if we're configured by yaml
	file := flag.String("f", "", "configuration file path")
	debug := flag.Bool("d", false, "configuration file path")
	flag.Parse()

	if len(*file) > 0 {
		yamlConfiguration, err = ioutil.ReadFile(*file)
		if err != nil {
			log.Panicf("error reading configuration file %s %+v", *file, err)
		}
	}
	env := createEnvironment(CreateConfiguration(new(iface.Configuration), *debug, yamlConfiguration), nil)
	if env.Configuration.Debug {
		log.Printf("current settings: %+v", env.Configuration)
	}
	appContext, cancelFunc := context.WithCancel(context.Background())
	hookOnExit(cancelFunc)
	go env.Server.Run(env.Configuration.Host + ":" + env.Configuration.Port)
	<-appContext.Done()
	return nil
}

// hookOnExit listens for signal SIGHUP and once received it closes the provided `closing` channel.
func hookOnExit(cancelFunc context.CancelFunc) {
	go func(cancelFunc context.CancelFunc) {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Kill, os.Interrupt)
		<-signals
		log.Println("Initiating shutdown ...")
		cancelFunc()
	}(cancelFunc)
}

// initializeRoutes does the basic stuff needed to create a router.
func initializeRoutes(router *gin.Engine, env *iface.Env) *gin.RouterGroup {
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
