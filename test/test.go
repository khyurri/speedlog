package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/khyurri/speedlog/engine"
	"github.com/khyurri/speedlog/engine/mongo"
	"github.com/khyurri/speedlog/engine/projects"
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

func (suite *SpeedLogTest) MakeRequest(req *http.Request, handler http.HandlerFunc) (result string) {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr.Body.String()
}

func (suite *SpeedLogTest) RegisterProject(title string) string {
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projects.RegisterProjectHttp(w, r, suite.Engine)
	})
	jsonStr, _ := json.Marshal(&projects.RegisterProjectRequest{Title: title})
	fmt.Println(string(jsonStr))
	req, _ := http.NewRequest("PUT", "/private/project/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	return suite.MakeRequest(req, f)
}
