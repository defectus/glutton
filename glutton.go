package glutton

import "github.com/gin-gonic/gin"
import "github.com/gin-contrib/cors"

func Run() error {
	env := createEnvironment(createSettings(DefaultSettings))
	go env.Server.Run()
	return nil
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

	}
}
