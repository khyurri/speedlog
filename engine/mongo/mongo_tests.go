package mongo

import (
	"encoding/json"
	"github.com/khyurri/speedlog/testutils"
	"io/ioutil"
	"path"
	"testing"
	"time"
)

// PACKAGE TEST CONFIGURATION
const (
	testMongoDb       = "speedlog_test"
	testMongoHost     = "localhost:27017"
	populateSourceDir = "testdata/"
	eventsFixtures    = "events.json"
	testTimeLayout    = "2006-01-02 15:04:05"
)

type eventsFixture struct {
	Title  string `json:"title"`
	Events []struct {
		MetricName string  `json:"metricName"`
		MetricTime string  `json:"metricTime"`
		DurationMs float64 `json:"durationMs"`
	} `json:"events"`
}

// wrappers to shorten name
var (
	ok     = testutils.Ok
	assert = testutils.Assert
	equals = testutils.Equals
)

func clearDb(t testing.TB, mongo *Mongo) error {
	return mongo.DropDatabase()
}

func populateDb(t testing.TB, mongo *Mongo, fixtureName string) {
	var jsonData []eventsFixture

	err := clearDb(t, mongo)
	ok(t, err)

	data, err := ioutil.ReadFile(path.Join(populateSourceDir, fixtureName))
	ok(t, err)

	err = json.Unmarshal(data, &jsonData)
	ok(t, err)

	for _, project := range jsonData {

		err = mongo.AddProject(project.Title)
		ok(t, err)

		for _, event := range project.Events {
			metricName := event.MetricName
			durationMs := event.DurationMs
			ts, err := time.Parse(testTimeLayout, event.MetricTime)
			ok(t, err)
			err = mongo.saveEventAtTime(metricName, project.Title, durationMs, ts)
		}

	}
}
