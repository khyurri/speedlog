package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type trEventSave struct {
	MetricName string  `json:"metricName"`
	DurationMs float64 `json:"durationMs"`
	Project    string  `json:"project"`
	ExpCode    int     // Expected http code
}

type trBatchEventSave struct {
	Id         int     `json:"id"`
	MetricName string  `json:"metricName"`
	DurationMs float64 `json:"durationMs"`
	Project    string  `json:"project"`
	ExpCode    int     // Expected http code
}

type trEventsSave struct {
	Event   []trBatchEventSave
	ExpCode int
}

type trEventGet struct {
	MetricName     string
	Project        string
	MetricTimeFrom string // valid format ...
	MetricTimeTo   string
	GroupBy        string
	Login          string
	Password       string
	ExpCode        int // Expected http code
}

func TestCreateEventHttp(t *testing.T) {

	testRounds := []trEventSave{
		{
			ExpCode:    200,
			MetricName: "testMetric",
			DurationMs: 0.1,
			Project:    "testProject",
		},
		{
			ExpCode:    500,
			MetricName: failMetricName,
			DurationMs: 0.1,
			Project:    "testProject",
		},
	}

	env := NewTestEnv(t, "*")
	router := mux.NewRouter()
	env.ExportEventRoutes(router)

	for round, event := range testRounds {
		jsonStr, err := json.Marshal(event)
		ok(t, err)

		r, err := http.NewRequest("PUT", "/event/", bytes.NewBuffer(jsonStr))
		ok(t, err)

		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		assert(t, event.ExpCode == w.Code, fmt.Sprintf("wrong code `%d` at round %d", w.Code, round))
	}
}

func TestCreateEventsHttp(t *testing.T) {
	testRounds := []trEventsSave{{
		Event: []trBatchEventSave{
			{
				Id:         1,
				MetricName: "testMetric",
				DurationMs: 0.1,
				Project:    "testProject",
				ExpCode:    200,
			},
			{
				Id:         2,
				MetricName: failMetricName,
				DurationMs: 0.1,
				Project:    "testProject",
				ExpCode:    400,
			},
		},
		ExpCode: 200,
	}}

	env := NewTestEnv(t, "*")
	router := mux.NewRouter()
	env.ExportEventRoutes(router)

	for round, event := range testRounds {
		jsonStr, err := json.Marshal(event.Event)
		r, err := http.NewRequest("PUT", "/events/", bytes.NewBuffer(jsonStr))
		ok(t, err)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		assert(t, event.ExpCode == w.Code, fmt.Sprintf("wrong code `%d` at round %d", w.Code, round))
		debug(w.Body)

	}

}

func TestGetEventsHttp(t *testing.T) {
	testRounds := []trEventGet{
		{
			// Empty eventRequest
			ExpCode: 404,
		},
		{
			// not authorized
			ExpCode:        403,
			MetricName:     "metricName",
			Project:        "someProject",
			MetricTimeFrom: "2019-09-02T00:01:00",
			MetricTimeTo:   "2019-09-02T00:02:00",
			GroupBy:        "minutes",
		},
		{
			// valid eventRequest with authorization
			ExpCode:        200,
			Login:          validLogin,
			MetricName:     "metricName",
			Project:        "someProject",
			MetricTimeFrom: "2019-09-02T00:01:00",
			MetricTimeTo:   "2019-09-02T00:02:00",
			GroupBy:        "minutes",
		},
	}

	env := NewTestEnv(t, "*")
	router := mux.NewRouter()
	env.ExportEventRoutes(router)

	// check get
	for round, event := range testRounds {
		u, err := url.Parse("/private/events/")
		ok(t, err)

		v := url.Values{}
		if len(event.MetricName) > 0 {
			v.Add("metricName", event.MetricName)
		}
		if len(event.Project) > 0 {
			v.Add("project", event.Project)
		}
		if len(event.GroupBy) > 0 {
			v.Add("groupBy", event.GroupBy)
		}
		if len(event.MetricTimeFrom) > 0 {
			v.Add("metricTimeFrom", event.MetricTimeFrom)
		}
		if len(event.MetricTimeTo) > 0 {
			v.Add("metricTimeTo", event.MetricTimeTo)
		}

		u.RawQuery = v.Encode()
		r, err := http.NewRequest("GET", u.String(), nil)
		ok(t, err)

		// authorize
		if len(event.Login) > 0 {
			token := getToken(t, env, event.Login)
			r.Header.Set("Authorization", "Bearer "+token)
		}

		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)

		assert(t, event.ExpCode == w.Code,
			fmt.Sprintf("wrong code `%d` at round %d.\nurl: %s\n", w.Code, round, u.String()))

		if 200 == w.Code {
			options(t, u.String(), router)
		}
	}
}

// check OPTIONS header
func options(t testing.TB, url string, router *mux.Router) {
	r, err := http.NewRequest("OPTIONS", url, nil)
	ok(t, err)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	equals(t, "Content-Type", w.HeaderMap["Access-Control-Allow-Headers"][0])
	t.Log(w.HeaderMap["Access-Control-Allow-Origin"])
	equals(t, "*", w.HeaderMap["Access-Control-Allow-Origin"][0])
}
