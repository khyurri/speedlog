package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/khyurri/speedlog/engine"
	"github.com/khyurri/speedlog/engine/mongo"
	"github.com/khyurri/speedlog/engine/projects"
	"github.com/khyurri/speedlog/engine/users"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

var DBName = "test_speedlog"

type SpeedLogTest struct {
	JWTTestKey string
	Mongo      string
	Logger     *log.Logger
	DBEngine   *mongo.Engine
	Engine     *engine.Engine
}

type JWTToken struct {
	Token string `json:"token"`
}

func (j *JWTToken) AuthHeader(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+j.Token)
}

func (suite *SpeedLogTest) Init() {
	suite.JWTTestKey = "test_key"
	suite.Mongo = "127.0.0.1:27017"
	suite.Logger = log.New(os.Stdout, DBName+" ", log.LstdFlags|log.Lshortfile)
	suite.DBEngine, _ = mongo.New(DBName, suite.Mongo, suite.Logger)
	suite.Engine = engine.New(suite.DBEngine, suite.Logger, suite.JWTTestKey)

	// clear mongo
	err := suite.DBEngine.DropDatabase()
	if err != nil {
		suite.Logger.Panic(err)
	}
}

func (suite *SpeedLogTest) MakeRequest(req *http.Request, handler http.HandlerFunc) (code int, result string) {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr.Code, rr.Body.String()
}

func (suite *SpeedLogTest) RegisterProject(title string, login string, password string) (result string, err error) {

	_, err = suite.Login(login, password)
	if err != nil {
		return
	}

	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projects.RegisterProjectHttp(w, r, suite.Engine)
	})
	jsonStr, _ := json.Marshal(&projects.RegisterProjectRequest{Title: title})
	req, _ := http.NewRequest("PUT", "/private/project/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	//authToken.AuthHeader(req)
	code, result := suite.MakeRequest(req, f)
	fmt.Println(code)
	return result, err
}

func (suite *SpeedLogTest) AddUser(login string, password string) error {
	return users.AddUser(login, password, suite.Engine)
}

func (suite *SpeedLogTest) LoginHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users.AuthenticateHttp(w, r, suite.Engine)
	})
}

func (suite *SpeedLogTest) Login(login string, password string) (result *JWTToken, err error) {
	jsonStr, _ := json.Marshal(&users.User{Login: login, Password: password})
	req, _ := http.NewRequest("POST", "/login/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	_, res := suite.MakeRequest(req, suite.LoginHandler())
	result = &JWTToken{}
	return result, json.Unmarshal([]byte(res), result)
}
