package glutton

import (
	"log"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func createSettings(settings *Settings) *Settings {
	if settings == nil {
		settings = &Settings{}
	}
	err := valueFromEnvVar(settings)
	if err != nil {
		log.Panicf("createSettings: failed to read settings %+v", err)
	}
	return settings
}

func createEnvironment(settings *Settings) *Env {
	env := &Env{Settings: settings}
	env.Server = gin.Default()
	initializeRoutes(env.Server, env)
	return env
}

func valueFromEnvVar(value interface{}) error {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr {
		return errors.Errorf("valueFromEnvVar: only pointer type values are supported.")
	}
	val = val.Elem()
	for i := 0; i < val.NumField(); i++ {
		if val.Type().Field(i).Type.Kind() != reflect.String {
			return errors.Errorf("valueFromEnvVar: only string fields are supported, not %s.", val.Type().Field(i).Type.Kind().String())
		}
		tag := val.Type().Field(i).Tag.Get("env")
		if len(tag) == 0 {
			tag = val.Type().Field(i).Name
		}
		if len(val.Field(i).String()) == 0 {
			val.Field(i).SetString(os.Getenv(tag))
		}
	}
	return nil
}
