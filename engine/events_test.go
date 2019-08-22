package engine

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine/mongo"
	"github.com/stretchr/testify/suite"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var TestProject = "testproject"

type EventsTestSuit struct {
	suite.Suite
	testEvents []interface{}
}

func (suite *EventsTestSuit) SetupTest() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, metric := range []string{"backendResponse", "frontendResponse"} {
		for i := 0; i < 2; i++ {
			suite.testEvents = append(suite.testEvents,
				struct {
					MetricName string  `json:"metricName"`
					DurationMs float64 `json:"durationMs"`
				}{
					metric,
					r.Float64(),
				})
		}
	}

}

func (suite *EventsTestSuit) TestStoreEvents() {

	router := mux.NewRouter()
	logger := log.New(os.Stdout, "speedlog ", log.LstdFlags|log.Lshortfile)
	dbEngine, _ := mongo.New("speedlog", "127.0.0.1:27017")
	env := NewEnv(dbEngine, logger, "1")
	env.ExportEventRoutes(router)

	for _, event := range suite.testEvents {
		jsonStr, _ := json.Marshal(event)
		r, _ := http.NewRequest("PUT", "/pravoved.ru/event/", bytes.NewBuffer(jsonStr))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		suite.Equal(200, w.Code)
	}
}

func TestEvents(t *testing.T) {
	suite.Run(t, new(EventsTestSuit))
}
