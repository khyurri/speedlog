package events

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/test"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var TestProject = "testproject"
var TestUser = "test_user"
var TestPassword = "test_password"

type EventsTestSuit struct {
	suite.Suite
	test.SpeedLogTest
	TestEvents []*SaveEventReq
}

func (suite *EventsTestSuit) SetupTest() {
	suite.Init()
	err := suite.AddUser(TestUser, TestPassword)
	if err != nil {
		suite.T().Log(err)
		suite.T().FailNow()
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, metric := range []string{"backend_response", "frontend_response"} {
		for i := 0; i < 2; i++ {
			suite.TestEvents = append(suite.TestEvents, &SaveEventReq{
				metric,
				r.Float64(),
			})
		}
	}

}

func (suite *EventsTestSuit) makeRequest(req *http.Request) (result string) {
	return
}

func (suite *EventsTestSuit) getStoreEventHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SaveEventHttp(w, r, suite.Engine)
	})
}

func (suite *EventsTestSuit) TestStoreEvents() {

	registered, err := suite.RegisterProject(TestProject, TestUser, TestPassword)
	suite.T().Log(err)
	suite.T().Log("[debug]" + registered)

	token, err := suite.Login(TestUser, TestPassword)
	if err != nil {
		suite.T().Log(err)
	}
	router := mux.NewRouter()
	app := rest.New(suite.Engine)
	ExportRoutes(router, app)
	for _, event := range suite.TestEvents {
		jsonStr, _ := json.Marshal(event)
		r, _ := http.NewRequest("PUT", "/"+TestProject+"/event/", bytes.NewBuffer(jsonStr))
		r.Header.Set("Content-Type", "application/json")
		token.AuthHeader(r)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		suite.Equal(200, w.Code)
	}
}

func TestEvents(t *testing.T) {
	suite.Run(t, new(EventsTestSuit))
}
