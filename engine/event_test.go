package engine

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type EventsTestSuit struct {
	suite.Suite
	testEvents []interface{}
}

func (suite *EventsTestSuit) SetupTest() {

	Logger = log.New(os.Stdout, "speedlog ", log.LstdFlags|log.Lshortfile)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, metric := range []string{"backendResponse", "frontendResponse"} {
		for i := 0; i < 2; i++ {
			suite.testEvents = append(suite.testEvents,
				struct {
					MetricName string  `json:"metricName"`
					DurationMs float64 `json:"durationMs"`
					Project    string  `json:"project"`
				}{
					metric,
					r.Float64(),
					"testProject", // todo: change name
				})
		}
	}

}

func (suite *EventsTestSuit) TestStoreEvents() {

	login, password := "admin10", "superpassword"

	router := mux.NewRouter()
	loc, _ := time.LoadLocation("Europe/Moscow")
	dbEngine, _ := mongo.New("speedlog", "127.0.0.1:27017")
	env := NewEnv(dbEngine, "1", loc)
	env.ExportEventRoutes(router)
	env.ExportUserRoutes(router)

	for _, event := range suite.testEvents {
		jsonStr, _ := json.Marshal(event)
		r, _ := http.NewRequest("PUT", "/event/", bytes.NewBuffer(jsonStr))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		suite.Equal(200, w.Code)
	}

	// login
	err := dbEngine.AddUser(login, password)
	defer func() {
		// delete user
		userId, err := dbEngine.GetUser(login)
		assert.Nil(suite.T(), err)
		err = dbEngine.UserDel(userId.Id.Hex())
		assert.Nil(suite.T(), err)
	}()

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
	suite.T().Log(w.Body)
	err = json.Unmarshal(w.Body.Bytes(), resp)
	assert.Nil(suite.T(), err)
	assert.Greater(suite.T(), len(resp.Token), 0)

	r, _ = http.NewRequest("GET", "/private/events/?metricName=backendResponse&metricTimeFrom=2019-08-20T01:10&metricTimeTo=2019-08-25T00:00&groupBy=minutes&project=pravoved.ru", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+resp.Token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	suite.Equal(200, w.Code)
	suite.Greater(len(w.Body.String()), 0)
}

func TestEvent(t *testing.T) {
	suite.Run(t, new(EventsTestSuit))
}
