package users

import (
	"bytes"
	"encoding/json"
	"github.com/khyurri/speedlog/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type UsersTestSuit struct {
	suite.Suite
	test.SpeedLogTest
	TestUsers map[string]string // login -> password
}

func (suite *UsersTestSuit) SetupTest() {

	suite.Init()

	suite.TestUsers = map[string]string{
		"admin": "password",
		"user1": "pas'\\sword1",
		"user2": "pa\"ssw0rd#~",
	}

	// add users
	for k, v := range suite.TestUsers {
		err := AddUser(k, v, suite.Engine)
		if err != nil {
			suite.Logger.Panic(err)
		}
	}

}

func (suite *UsersTestSuit) makeRequest(req *http.Request) (result string) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthenticateHttp(w, r, suite.Engine)
	})
	handler.ServeHTTP(rr, req)
	return rr.Body.String()
}

func (suite *UsersTestSuit) TestAuthenticateHttp() {
	// test invalid creds
	errMsg := "{\"message\":\"invalid login or password\"}"
	req, _ := http.NewRequest("POST", "/login/", nil)
	res := suite.makeRequest(req)
	assert.Equal(suite.T(), res, errMsg)

	// test valid creds
	for login, password := range suite.TestUsers {
		creds := &User{login, password}
		jsonStr, _ := json.Marshal(creds)
		req, _ := http.NewRequest("POST", "/login/", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		res := suite.makeRequest(req)
		suite.T().Log(res)
		assert.Contains(suite.T(), res, "token")
	}
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(UsersTestSuit))
}
