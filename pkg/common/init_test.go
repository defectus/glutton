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

func TestCreateRedirectHandlerNoRedirect(t *testing.T) {
	mp := &MockParser{}
	mp.On("Parse").Return(&iface.PayloadRecord{}, nil)
	ms := &MockSaver{}
	ms.On("Save").Return(nil)
	mn := &MockNotifier{}
	mn.On("Notify").Return(nil)
	router := gin.Default()
	router.POST("test", redirectHandler(createHandler("test", mp, mn, ms, false), http.StatusTemporaryRedirect, ""))
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

func TestCreateRedirectHandlerRedirect(t *testing.T) {
	mp := &MockParser{}
	mp.On("Parse").Return(&iface.PayloadRecord{}, nil)
	ms := &MockSaver{}
	ms.On("Save").Return(nil)
	mn := &MockNotifier{}
	mn.On("Notify").Return(nil)
	router := gin.Default()
	router.POST("test", redirectHandler(createHandler("test", mp, mn, ms, false), http.StatusTemporaryRedirect, "https://test.redirect"))
	req, _ := http.NewRequest("POST", "http://localhost/test", nil)
	testHTTPResponse(t, router, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		p, _ := ioutil.ReadAll(w.Body)
		log.Printf("server reply: %s", string(p))
		log.Printf("server header: %+v", w.HeaderMap)
		mp.AssertExpectations(t)
		ms.AssertExpectations(t)
		mn.AssertExpectations(t)
		return true
	})
}

func TestGenerateToken(t *testing.T) {
	env := createEnvironment(&iface.Configuration{
		Debug: true,
		Settings: []iface.Settings{
			{
				UseToken: true,
				URI:      "save",
			},
		},
	}, nil)
	req, _ := http.NewRequest("GET", "http://localhost/v1/glutton/save/token", nil)
	testHTTPResponse(t, env.Server, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(t, http.StatusOK, w.Code)
		p, _ := ioutil.ReadAll(w.Body)
		log.Printf("server reply: %s", string(p))
		log.Printf("server header: %+v", w.HeaderMap)
		return true
	})
}

func TestValidateToken1(t *testing.T) {
	// test fail token (token invalid)
	env := &iface.Env{
		Savers:    map[string]reflect.Type{},
		Notifiers: map[string]reflect.Type{},
		Parsers:   map[string]reflect.Type{},
	}
	env.Notifiers["TestNotifier"] = reflect.TypeOf(TestNotifier{})
	env.Savers["TestSaver"] = reflect.TypeOf(TestSaver{})
	env.Parsers["TestParser"] = reflect.TypeOf(TestParser{})
	env = createEnvironment(&iface.Configuration{
		Debug: true,
		Settings: []iface.Settings{
			{
				UseToken: true,
				URI:      "save",
				Saver:    "TestSaver",
				Parser:   "TestParser",
				Notifier: "TestNotifier",
			},
		},
	}, env)
	req, _ := http.NewRequest("POST", "http://localhost/v1/glutton/save", nil)
	req.Header.Add("token", "dummy")
	testHTTPResponse(t, env.Server, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(t, http.StatusPreconditionFailed, w.Code)
		p, _ := ioutil.ReadAll(w.Body)
		log.Printf("server reply: %s", string(p))
		log.Printf("server header: %+v", w.HeaderMap)
		return true
	})
}

func TestValidateToken2(t *testing.T) {
	// test get token and use token
	env := &iface.Env{
		Savers:    map[string]reflect.Type{},
		Notifiers: map[string]reflect.Type{},
		Parsers:   map[string]reflect.Type{},
	}
	env.Notifiers["TestNotifier"] = reflect.TypeOf(TestNotifier{})
	env.Savers["TestSaver"] = reflect.TypeOf(TestSaver{})
	env.Parsers["TestParser"] = reflect.TypeOf(TestParser{})
	env = createEnvironment(&iface.Configuration{
		Debug: true,
		Settings: []iface.Settings{
			{
				UseToken: true,
				URI:      "save",
				Saver:    "TestSaver",
				Parser:   "TestParser",
				Notifier: "TestNotifier",
			},
		},
	}, env)
	token := []byte{}
	req, _ := http.NewRequest("GET", "http://localhost/v1/glutton/save/token", nil)
	testHTTPResponse(t, env.Server, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(t, http.StatusOK, w.Code)
		token, _ = ioutil.ReadAll(w.Body)
		return true
	})
	req, _ = http.NewRequest("POST", "http://localhost/v1/glutton/save", nil)
	req.Header.Add("token", string(token))
	testHTTPResponse(t, env.Server, req, func(w *httptest.ResponseRecorder) bool {
		assert.Equal(t, http.StatusOK, w.Code)
		p, _ := ioutil.ReadAll(w.Body)
		log.Printf("server reply: %s", string(p))
		log.Printf("server header: %+v", w.HeaderMap)
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
