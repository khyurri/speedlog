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

type ProjectTestSuite struct {
	suite.Suite
}

func (suite *ProjectTestSuite) SetupTest() {
	Logger = log.New(os.Stdout, "speedlog ", log.LstdFlags|log.Lshortfile)
}

func (suite *ProjectTestSuite) TestCreateProject() {
	projectTitle := "test_project"
	dbEngine, _ := mongo.New("speedlog", "127.0.0.1:27017")
	err := dbEngine.AddProject(projectTitle)
	assert.Nil(suite.T(), err)

	project, err := dbEngine.GetProject(projectTitle)
	assert.Nil(suite.T(), err)
	projectId := project.ID.Hex()
	assert.Greater(suite.T(), len(projectId), 0)

	err = dbEngine.DelProject(projectId)
	assert.Nil(suite.T(), err)
}

func (suite *ProjectTestSuite) TestCreateProjectHTTP() {
	projectTitle := "test_http_project"
	login, password := "admin10", "superpassword"
	dbEngine, _ := mongo.New("speedlog", "127.0.0.1:27017")

	defer func() {
		// delete user
		userId, err := dbEngine.GetUser(login)
		assert.Nil(suite.T(), err)
		err = dbEngine.UserDel(userId.Id.Hex())
		assert.Nil(suite.T(), err)
		// delete project
		project, err := dbEngine.GetProject(projectTitle)
		assert.Nil(suite.T(), err)
		projectId := project.ID.Hex()
		err = dbEngine.DelProject(projectId)
		assert.Nil(suite.T(), err)
	}()
	loc, _ := time.LoadLocation("Europe/Moscow")
	router := mux.NewRouter()
	env := NewEnv(dbEngine, "1", loc)
	env.ExportProjectRoutes(router)
	env.ExportUserRoutes(router)

	// create user
	err := dbEngine.AddUser(login, password)

	// login

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

	// create project

	jsonStr, _ = json.Marshal(struct {
		Title string `json:"title"`
	}{projectTitle})
	r, _ = http.NewRequest("PUT", "/private/project/", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+resp.Token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	suite.T().Log(w.Code)

}

func TestProject(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}
