package engine

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type UserTestSuit struct {
	suite.Suite
}

func (suite *UserTestSuit) SetupTest() {
	Logger = log.New(os.Stdout, "speedlog ", log.LstdFlags|log.Lshortfile)
}

func (suite *UserTestSuit) TestAddUser() {

	login, password := "admin9", "hello"

	dbEngine, _ := mongo.New("speedlog", "127.0.0.1:27017")
	loc, _ := time.LoadLocation("Europe/Moscow")
	env := NewEnv(dbEngine, "1", loc)
	err := env.AddUser(login, password)
	assert.Nil(suite.T(), err)

	// check user exists
	user, err := env.DBEngine.GetUser(login)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), user.Id)

	// check auth
	router := mux.NewRouter()
	env.ExportUserRoutes(router)

	jsonStr, _ := json.Marshal(struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{login, password})

	r, _ := http.NewRequest("POST", "/login/", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	resp := &struct {
		Token string `json:"token"`
	}{}
	err = json.Unmarshal(w.Body.Bytes(), resp)
	assert.Nil(suite.T(), err)
	assert.Greater(suite.T(), len(resp.Token), 0)

	// check wrong password
	jsonStr, _ = json.Marshal(struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{login, "***"})

	r, _ = http.NewRequest("POST", "/login/", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(suite.T(), 403, w.Code)

	// del user
	err = env.DBEngine.UserDel(user.Id.Hex())
	assert.Nil(suite.T(), err)
}

func TestUser(t *testing.T) {
	suite.Run(t, new(UserTestSuit))
}
