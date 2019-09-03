package engine

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/khyurri/speedlog/engine/mongo"
	"github.com/khyurri/speedlog/testutils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"testing"
	"time"
)

// wrappers to shorten name
var (
	ok     = testutils.Ok
	assert = testutils.Assert
	equals = testutils.Equals
)

// PACKAGE TEST CONFIGURATION
const (
	singKey                = "60067d34056e7e4496127749d15b5dc6"
	location               = "Local"
	failMetricName         = "fmn"
	duplicatedProjectTitle = "fpt"
	validLogin             = "vLogin"
	validPassword          = "vPassword"
)

// Mongodb mock

type DataStoreMock struct {
}

func (d DataStoreMock) FilterEvents(from, to time.Time, metricName, project string) (events []mongo.Event, err error) {
	if metricName == failMetricName {
		err = errors.New("[testing] failed by metric name")
		return
	}
	return
}

func (d DataStoreMock) AllEvents(from, to time.Time) (events []mongo.AllEvents, err error) {
	panic("implement me")
}

func (d DataStoreMock) SaveEvent(metricName, project string, durationMs float64) (err error) {
	if metricName == failMetricName {
		return errors.New("[testing] failed by metric name")
	}
	return
}

func (d DataStoreMock) AddUser(login string, password string) (err error) {
	panic("implement me")
}

func (d DataStoreMock) GetUser(login string) (*mongo.User, error) {
	if login == validLogin {
		bytes, err := bcrypt.GenerateFromPassword([]byte(validPassword), 10)
		return &mongo.User{
			Login:    validLogin,
			Password: string(bytes)}, err
	}
	return nil, errors.New("[testing] user not found")
}

func (d DataStoreMock) UserDel(uid string) error {
	panic("implement me")
}

func (d DataStoreMock) AddProject(title string) (err error) {
	if title == duplicatedProjectTitle {
		return errors.New("[testing] failed by project title")
	}
	return
}

func (d DataStoreMock) GetProject(title string) (project mongo.Project, err error) {
	panic("implement me")
}

func (d DataStoreMock) GetProjectById(id string) (project mongo.Project, err error) {
	panic("implement me")
}

func (d DataStoreMock) DelProject(id string) (err error) {
	panic("implement me")
}

// Test engine
func NewTestEnv(t testing.TB, allowOrigin string) (env *Env) {
	dbEngine := &DataStoreMock{}
	location, err := time.LoadLocation(location)
	ok(t, err)
	env = NewEnv(dbEngine, singKey, location)
	env.Logger = log.New(os.Stdout, "testing ", log.LstdFlags|log.Lshortfile)
	if len(allowOrigin) > 0 {
		env.AllowOrigin = allowOrigin
	}
	return
}

func getToken(t *testing.T, env *Env, login string) (tokenString string) {
	if login != validLogin {
		return
	}
	_, tokenString, err := env.SigningKey.Encode(
		jwt.MapClaims{"source": "rest", "issuer": login})
	ok(t, err)
	return
}
