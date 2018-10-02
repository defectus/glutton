package glutton

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueFromEnvVar(t *testing.T) {
	os.Setenv("TEST", "TEST1")
	value := &struct {
		Test string `env:"TEST"`
	}{}
	err := valueFromEnvVar(value)
	assert.NoError(t, err)
	assert.Equal(t, "TEST1", value.Test)
}

func TestValueFromEnvVar2(t *testing.T) {
	value := struct {
		Test string `env:"TEST"`
	}{}
	err := valueFromEnvVar(value)
	assert.Error(t, err)
}

func TestValueFromEnvVar3(t *testing.T) {
	value := struct {
		Test int `env:"TEST"`
	}{}
	err := valueFromEnvVar(value)
	assert.Error(t, err)
}

func TestValueFromEnvVar4(t *testing.T) {
	os.Setenv("TEST", "value")
	value := &struct {
		Test string `env:"TEST"`
	}{Test: "override"}
	err := valueFromEnvVar(value)
	assert.NoError(t, err)
	assert.Equal(t, "override", value.Test)
}
