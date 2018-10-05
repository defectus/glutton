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
	os.Setenv("TEST", "")
	value := &struct {
		Test string `env:"TEST" default:"override"`
	}{}
	err := valueFromEnvVar(value)
	assert.NoError(t, err)
	assert.Equal(t, "override", value.Test)
}

func TestValueFromEnvVar5(t *testing.T) {
	os.Setenv("TEST", "TEST1")
	os.Setenv("TESTBOOL", "true")
	os.Setenv("TESTINT", "1")
	value := &struct {
		Test     string `env:"TEST"`
		TestBool bool   `env:"TESTBOOL"`
		TestInt  int    `env:"TESTINT"`
	}{}
	err := valueFromEnvVar(value)
	assert.NoError(t, err)
	assert.Equal(t, "TEST1", value.Test)
	assert.Equal(t, 1, value.TestInt)
	assert.True(t, value.TestBool)
}

func TestValueFromEnvVar7(t *testing.T) {
	value := new(int)
	err := valueFromEnvVar(value)
	assert.Error(t, err)
}
