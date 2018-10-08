package glutton

import (
	"log"
	"os"
	"reflect"
	"strconv"

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
	env.Notifier = &NilNotifier{}
	env.Saver = &SimpleFileSystemSaver{}
	env.Parser = &SimpleParser{}
	env.Notifier.Configure(env.Settings)
	env.Saver.Configure(env.Settings)
	env.Parser.Configure(env.Settings)
	return env
}

func valueFromEnvVar(value interface{}) error {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr {
		return errors.Errorf("valueFromEnvVar: only pointer type values are supported.")
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return errors.Errorf("valueFromEnvVar: only struct types are supported.")
	}
	for i := 0; i < val.NumField(); i++ {
		tag := val.Type().Field(i).Tag.Get("env")
		if len(tag) == 0 {
			tag = val.Type().Field(i).Name
		}
		v := os.Getenv(tag)
		if len(v) == 0 {
			v = val.Type().Field(i).Tag.Get("default")
		}
		switch val.Type().Field(i).Type.Kind() {
		case reflect.String:
			val.Field(i).SetString(v)
		case reflect.Int:
			in, _ := strconv.ParseInt(v, 10, 64)
			val.Field(i).SetInt(in)
		case reflect.Bool:
			bo, _ := strconv.ParseBool(v)
			val.Field(i).SetBool(bo)
		case reflect.Ptr:
			if val.Type().Field(i).Type.Elem().Kind() == reflect.Struct {
				err := valueFromEnvVar(val.Field(i).Interface())
				if err != nil {
					return errors.Wrapf(err, "error processing %s", val.Field(i).Type().Name())
				}
			}
		default:
			return errors.Errorf("valueFromEnvVar: unsupported kind %s.", val.Type().Field(i).Type.Kind())
		}
	}
	return nil
}
