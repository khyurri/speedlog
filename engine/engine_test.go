package engine

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/khyurri/speedlog/engine/mongo"
	"github.com/khyurri/speedlog/testutils"
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

// Test engine
func NewTestEnv(t testing.TB, allowOrigin string) (env *Env) {
	dbEngine := &mongo.DataStoreMock{
		FailMetricName:         failMetricName,
		ValidLogin:             validLogin,
		ValidPassword:          validPassword,
		DuplicatedProjectTitle: duplicatedProjectTitle,
	}
	location, err := time.LoadLocation(location)
	ok(t, err)
	env = NewEnv(dbEngine, singKey, location)
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
