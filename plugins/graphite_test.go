package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine"
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

type GraphiteTestSuite struct {
	suite.Suite
	testEvents []interface{}
}

func (suite *GraphiteTestSuite) SetupTest() {
	engine.Logger = log.New(os.Stdout, "speedlog ", log.LstdFlags|log.Lshortfile)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, metric := range []string{"backendResponse", "frontendResponse"} {
		for i := 0; i < 10; i++ {
			suite.testEvents = append(suite.testEvents,
				struct {
					MetricName string  `json:"metricName"`
					DurationMs float64 `json:"durationMs"`
					Project    string  `json:"project"`
				}{
					metric,
					r.Float64(),
					"pravoved.ru", // todo: change name
				})
		}
	}
}

func (suite *GraphiteTestSuite) TestExportAll() {

	router := mux.NewRouter()
	dbEngine, _ := mongo.New("speedlog", "127.0.0.1:27017")
	env := engine.NewEnv(dbEngine, "1")
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
	now := time.Now()
	startTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), 0, 0, now.Location())

	endTime := startTime.Add(time.Minute)

	events, err := dbEngine.AllEvents(startTime, endTime)
	assert.Nil(suite.T(), err)
	for _, event := range events {
		fmt.Println(event)
	}

}

func TestGraphite(t *testing.T) {
	suite.Run(t, new(GraphiteTestSuite))
}
