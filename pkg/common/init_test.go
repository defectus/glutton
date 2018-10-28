package common

import (
	"io/ioutil"
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
	mc := MockConfigurable{}
	types := map[string]reflect.Type{"dummy": reflect.TypeOf(mc)}
	instance, err := createInstanceOf(types, "dummy", nil)
	assert.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(mc).Name(), reflect.TypeOf(instance).Elem().Name())
}

type MockSaver struct {
	mock.Mock
}

func (m *MockSaver) Save(*iface.PayloadRecord) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSaver) Configure(*iface.Settings) error {
	return nil
}

type MockParser struct {
	mock.Mock
}

func (m *MockParser) Configure(*iface.Settings) error {
	return nil
}

func (m *MockParser) Parse(*http.Request) (*iface.PayloadRecord, error) {
	args := m.Called()
	return args.Get(0).(*iface.PayloadRecord), args.Error(1)
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

func TestCreateHandler(t *testing.T) {
	mp := &MockParser{}
	mp.On("Parse").Return(&iface.PayloadRecord{}, nil)
	ms := &MockSaver{}
	ms.On("Save").Return(nil)
	mn := &MockNotifier{}
	mn.On("Notify").Return(nil)
	router := gin.Default()
	router.POST("test", createHandler("test", mp, mn, ms, false))
	req, _ := http.NewRequest("POST", "http://localhost/test", nil)
	testHTTPResponse(t, router, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(t, http.StatusOK, w.Code)
		p, _ := ioutil.ReadAll(w.Body)
		log.Printf("server reply: %s", string(p))
		mp.AssertExpectations(t)
		ms.AssertExpectations(t)
		mn.AssertExpectations(t)
		return true
	})
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
