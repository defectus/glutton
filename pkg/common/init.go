package common

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/defectus/glutton/pkg/notifier"
	"github.com/defectus/glutton/pkg/saver"
	"github.com/defectus/glutton/pkg/parser"
	"github.com/defectus/glutton/pkg/iface"
)

func createConfiguration(configuration *iface.Configuration, debug bool, yamlConfiguration []byte) *iface.Configuration {
	if configuration == nil {
		configuration = new(iface.Configuration)
	}
	valueFromEnvVar(configuration)
	configuration.Debug = debug
	if len(yamlConfiguration) > 0 {
		if err := yaml.Unmarshal(yamlConfiguration, configuration); err != nil {
			log.Printf("createConfigration: error parsing configuration %+v", err)
		}
	}
	// second, try to use environment to configure the app
	if len(configuration.Settings) == 0 {
		if configuration.Debug {
			log.Printf("configuration afer yaml contains no settings, using environment to configure")
		}
		settings := new(iface.Settings)
		err := valueFromEnvVar(settings)
		if err != nil {
			log.Panicf("createConfigration: failed to read settings %+v", err)
		}
		configuration.Settings = append(configuration.Settings, *settings)
	}
	return configuration
}

func createEnvironment(configuration *iface.Configuration) *iface.Env {
	env := &iface.Env{Configuration: configuration}
	if !configuration.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	env.Server = gin.Default()
	gluttonRoute := initializeRoutes(env.Server, env)
	env.Notifiers = map[string]reflect.Type{"NilNotifier": reflect.TypeOf(notifier.NilNotifier{}), "SMTPNotifier": reflect.TypeOf(notifier.SMTPNotifier{})}
	env.Savers = map[string]reflect.Type{"SimpleFileSystemSaver": reflect.TypeOf(saver.SimpleFileSystemSaver{})}
	env.Parsers = map[string]reflect.Type{"SimpleParser": reflect.TypeOf(parser.SimpleParser{})}
	for _, settings := range env.Configuration.Settings {
		var (
			instance interface{}
			notifier iface.PayloadNotifier
			saver    iface.PayloadSaver
			parser   iface.PayloadParser
			err      error
			ok       bool
		)
		if len(settings.Notifier) > 0 {
			instance, err = createInstanceOf(env.Notifiers, settings.Notifier, &settings)
			if err != nil {
				log.Panicf("error creating notifier %+v", err)
			}
			if notifier, ok = instance.(iface.PayloadNotifier); !ok {
				log.Panicf("exptected notifier, got %s", reflect.TypeOf(instance))
			}
		}
		if len(settings.Saver) > 0 {
			instance, err = createInstanceOf(env.Savers, settings.Saver, &settings)
			if err != nil {
				log.Panicf("error creating saver %+v", err)
			}
			if saver, ok = instance.(iface.PayloadSaver); !ok {
				log.Panicf("exptected saver, got %s", reflect.TypeOf(instance))
			}
		}
		if len(settings.Parser) > 0 {
			instance, err = createInstanceOf(env.Parsers, settings.Parser, &settings)
			if err != nil {
				log.Panicf("error creating parser %+v", err)
			}
			if parser, ok = instance.(iface.PayloadParser); !ok {
				log.Panicf("exptected parser, got %s", reflect.TypeOf(instance))
			}
		}
		gluttonRoute.POST(settings.URI, createHandler(settings.URI, parser, notifier, saver, settings.Debug))
	}
	return env
}

// createHandler appends a route to router and initialize the basic flow (request -> parser -> notifier -> saver)
func createHandler(URI string, parser iface.PayloadParser, notifier iface.PayloadNotifier, saver iface.PayloadSaver, debug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := parser.Parse(c.Request)
		if err != nil {
			log.Printf("%s: error parsing contents %+v", URI, err)
			log.Printf("%+v", c.Request)
		}
		err = notifier.Notify(payload)
		if err != nil {
			log.Printf("%s: error notifying of payload %+v", URI, err)
			log.Printf("%+v", payload)
		}
		err = saver.Save(payload)
		if err != nil {
			log.Printf("%s: error saving payload %+v", URI, err)
			log.Printf("%+v", payload)
		}
		c.Status(http.StatusOK)
	}
}

func createInstanceOf(types map[string]reflect.Type, name string, settings *iface.Settings) (interface{}, error) {
	v := reflect.New(types[name])
	if c, ok := v.Interface().(iface.Configurable); ok {
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
			log.Printf("valueFromEnvVar: unsupported kind %s.", val.Type().Field(i).Type.Kind())
		}
	}
	return nil
}
