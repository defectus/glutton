package common

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/defectus/glutton/pkg/iface"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

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

func TestValueFromEnvVar6(t *testing.T) {
	os.Setenv("TEST", "TEST1")
	os.Setenv("TESTBOOL", "true")
	os.Setenv("TESTINT", "1")
	os.Setenv("INNER", "inner")
	value := &struct {
		Test       string `env:"TEST"`
		TestBool   bool   `env:"TESTBOOL"`
		TestInt    int    `env:"TESTINT"`
		TestStruct *struct {
			Inner string `env:"INNER"`
		}
	}{TestStruct: &struct {
		Inner string `env:"INNER"`
	}{}}
	err := valueFromEnvVar(value)
	assert.NoError(t, err)
	assert.Equal(t, "inner", value.TestStruct.Inner)
	assert.Equal(t, "TEST1", value.Test)
	assert.Equal(t, 1, value.TestInt)
	assert.True(t, value.TestBool)
}

func TestValueFromEnvVar7(t *testing.T) {
	value := new(int)
	err := valueFromEnvVar(value)
	assert.Error(t, err)
}

type MockConfigurable struct {
	mock.Mock
}

func (m *MockConfigurable) Configure(*iface.Settings) error {
	return nil
}

func TestCreateInstanceOf(t *testing.T) {
	mc := &MockConfigurable{}
	types := map[string]reflect.Type{"dummy": reflect.TypeOf(mc)}
	instance, err := createInstanceOf(types, "dummy", nil)
	assert.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(mc).Name(), reflect.TypeOf(instance).Elem().Name())
}

type MockSaver struct {
	mock.Mock
}

func (m *MockSaver) Save(payload *iface.PayloadRecord) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSaver) Configure(*iface.Settings) error {
	return nil
}

type TestSaver struct {
}

func (t *TestSaver) Save(payload *iface.PayloadRecord) error {
	log.Printf("TestSaver_Save called with %+v", payload)
	return nil
}

func (t *TestSaver) Configure(settings *iface.Settings) error {
	log.Printf("TestSaver_Configure called with %+v", settings)
	return nil
}

type MockParser struct {
	mock.Mock
}

func (m *MockParser) Configure(*iface.Settings) error {
	return nil
}

func (m *MockParser) Parse(request *http.Request) (*iface.PayloadRecord, error) {
	args := m.Called()
	return args.Get(0).(*iface.PayloadRecord), args.Error(1)
}

type TestParser struct {
}

func (t *TestParser) Configure(settings *iface.Settings) error {
	log.Printf("TestParser_Configure called with %+v", settings)
	return nil
}

func (t *TestParser) Parse(request *http.Request) (*iface.PayloadRecord, error) {
	log.Printf("TestParser_Parse called with %+v", request)
	return &iface.PayloadRecord{}, nil
}

type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) Configure(*iface.Settings) error {
	return nil
}

func (m *MockNotifier) Notify(*iface.PayloadRecord) error {
	args := m.Called()
	return args.Error(0)
}

type TestNotifier struct {
}

func (t *TestNotifier) Configure(settings *iface.Settings) error {
	log.Printf("TestNotifier_Configure called with %+v", settings)
	return nil
}

func (t *TestNotifier) Notify(payload *iface.PayloadRecord) error {
	log.Printf("TestNotifier_Notify called with %+v", payload)
	return nil
}

func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if !f(w) {
		t.Fail()
	}
}

func TestCreateConfiguration1(t *testing.T) {
	os.Args = []string{"test", "-d"}
	config := CreateConfiguration(nil, true, nil)
	assert.True(t, config.Debug)
}

func TestCreateConfiguration2(t *testing.T) {
	os.Args = []string{"test", "-d"}
	yaml := `
settings:
  - name: test glutton
    redirect: /url
    parser: test`
	config := CreateConfiguration(nil, true, []byte(yaml))
	assert.Equal(t, "test glutton", config.Settings[0].Name)
	assert.Equal(t, "/url", config.Settings[0].Redirect)
	assert.Equal(t, "test", config.Settings[0].Parser)
}
