package test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type UsersTestSuit struct {
	suite.Suite
	SpeedLogTest
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
		err := suite.AddUser(k, v)
		if err != nil {
			suite.Logger.Panic(err)
		}
	}

}

func (suite *UsersTestSuit) TestAuthenticateHttp() {
	// test invalid creds
	errMsg := "{\"message\":\"invalid login or password\"}"
	req, _ := http.NewRequest("POST", "/login/", nil)
	_, res := suite.MakeRequest(req, suite.LoginHandler())
	assert.Equal(suite.T(), res, errMsg)

	// test valid creds
	for login, password := range suite.TestUsers {
		creds := &users.User{login, password}
		jsonStr, _ := json.Marshal(creds)
		req, _ := http.NewRequest("POST", "/login/", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		_, res := suite.MakeRequest(req, suite.LoginHandler())
		suite.T().Log(res)
		assert.Contains(suite.T(), res, "token")
	}
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(UsersTestSuit))
}
