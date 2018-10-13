package glutton

import (
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func createConfiguration(configuration *Configuration) *Configuration {
	if configuration == nil {
		configuration = new(Configuration)
	}
	// first see if we're configured by yaml
	// second, try to use environment to configure at lease on setting
	settings := new(Settings)
	err := valueFromEnvVar(settings)
	if err != nil {
		log.Panicf("createConfigration: failed to read settings %+v", err)
	}
	configuration.Settings = append(configuration.Settings, *settings)
	return configuration
}

func createEnvironment(configuration *Configuration) *Env {
	env := &Env{Configuration: configuration}
	env.Server = gin.Default()
	gluttonRoute := initializeRoutes(env.Server, env)
	env.Notifiers = map[string]reflect.Type{"NilNotifier": reflect.TypeOf(NilNotifier{}), "SimpleFileSystemSaver": reflect.TypeOf(SimpleFileSystemSaver{})}
	env.Savers = map[string]reflect.Type{"SimpleFileSystemSaver": reflect.TypeOf(SimpleFileSystemSaver{})}
	env.Parsers = map[string]reflect.Type{"SimpleParser": reflect.TypeOf(SimpleParser{})}
	for _, settings := range env.Configuration.Settings {
		var (
			instance interface{}
			notifier PayloadNotifier
			saver    PayloadSaver
			parser   PayloadParser
			err      error
			ok       bool
		)
		if len(settings.Notifier) > 0 {
			instance, err = createInstanceOf(env.Notifiers, settings.Notifier, &settings)
			if err != nil {
				log.Panicf("error creating notifier %+v", err)
			}
			if notifier, ok = instance.(PayloadNotifier); !ok {
				log.Panicf("exptected notifier, got %s", reflect.TypeOf(instance))
			}
		}
		if len(settings.Saver) > 0 {
			instance, err = createInstanceOf(env.Savers, settings.Saver, &settings)
			if err != nil {
				log.Panicf("error creating saver %+v", err)
			}
			if saver, ok = instance.(PayloadSaver); !ok {
				log.Panicf("exptected saver, got %s", reflect.TypeOf(instance))
			}
		}
		if len(settings.Parser) > 0 {
			instance, err = createInstanceOf(env.Parsers, settings.Parser, &settings)
			if err != nil {
				log.Panicf("error creating parser %+v", err)
			}
			if parser, ok = instance.(PayloadParser); !ok {
				log.Panicf("exptected parser, got %s", reflect.TypeOf(instance))
			}
		}
		//TODO: ^^^^ default to nil or something one of the prevs. missing. ^^^^
		gluttonRoute.POST(settings.URI, createHandler(gluttonRoute, settings.URI, parser, notifier, saver))
	}
	return env
}

// createHandler appends a route to router and initialize the basic flow (request -> parser -> notifier -> saver)
func createHandler(baseRoute *gin.RouterGroup, URI string, parser PayloadParser, notifier PayloadNotifier, saver PayloadSaver) gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO: finish this
	}
}

// func savePayload(env *Env) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		payload, err := env.Parser.Parse(c.Request)
// 		if err != nil {

// 		}
// 		err = env.Notifier.Notify(payload)
// 		if err != nil {

// 		}
// 		err = env.Saver.Save(payload)
// 		if err != nil {

// 		}
// 		c.Status(http.StatusOK)
// 	}
// }

func createInstanceOf(types map[string]reflect.Type, name string, settings *Settings) (interface{}, error) {
	v := reflect.New(types[name]).Elem()
	if c, ok := v.Interface().(Configurable); ok {
		err := c.Configure(settings)
		if err != nil {
			return nil, errors.Wrapf(err, "error configuring instance %s with %+v", name, settings)
		}
	}
	return v.Interface(), nil
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
