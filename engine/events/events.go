package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine"
	"github.com/khyurri/speedlog/engine/projects"
	"github.com/khyurri/speedlog/rest"
	"net/http"
	"time"
)

const (
	collection = "events"
	timeLayout = "2006-01-02T15:04:05"
)

const (
	GroupByMins   = iota
	GroupByHours  = iota
	GroupByDays   = iota
	GroupByWeeks  = iota
	GroupByMonths = iota
)

// Mongo event document
type Event struct {
	MetricName     string        `bson:"metric_name"`
	MetricTime     time.Time     `bson:"metric_time,omitempty"`
	ProjectId      bson.ObjectId `bson:"project_id"`
	DurationMs     int           `bson:"duration_ms,omitempty"`
	MetricTimeFrom time.Time     `bson:"metric_time_from,omitempty"`
	MetricTimeTo   time.Time     `bson:"metric_time_to,omitempty"`
	GroupBy        int           `bson:"group_by,omitempty"`
}

// Check and cast function type
type CACFunc func(string, *Event, *engine.Engine) error

// ROUTES

func ExportRoutes(router *mux.Router, app *rest.App) {
	router.HandleFunc("/{project}/event/", app.MongoEngine(SaveEventHttp)).
		Methods("PUT")

	private := router.PathPrefix("/private/").Subrouter()
	private.HandleFunc("/{project}/events/", app.MongoEngine(GetEventsHttp)).
		Methods("GET").
		Queries("metric_name", "{metricName}").
		Queries("metric_time_from", "{metricTimeFrom}").
		Queries("metric_time_to", "{metricTimeTo}").
		Queries("group_by", "{groupBy}")
	private.Use(rest.JWTMiddleware)
}

//////////////////////////////////////////////////////////////////////
//
// Events logic
//
//////////////////////////////////////////////////////////////////////

// PUT /pravoved.ru/event/
func SaveEventHttp(w http.ResponseWriter, r *http.Request, eng *engine.Engine) {
	var err error

	response := &rest.Resp{}
	defer response.Render(w)

	event := &Event{}
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&event)
	if err != nil {
		eng.Logger.Fatal(err)
		response.Status = rest.StatusIntErr
		return
	}

	cast := &CheckAndCast{
		[]string{
			vars["project"],
			"", // event.MetricTime = time.Now()
		},
		[]CACFunc{
			CACProject,
			CACMetricTimeNow,
		},
		err, eng,
	}
	cast.execute(event)

	err = SaveEvent(event, eng)
	if err == nil {
		saved := struct {
			Saved bool `json:"saved"`
		}{true}
		response.Status = rest.StatusOk
		response.JsonBody, err = json.Marshal(saved)
	}
	if err != nil {
		eng.Logger.Fatal(err)
		response.Status = rest.StatusIntErr
		return
	}

}

// GET /pravoved.ru/events/?metric_time_from=2019-08-02T00:00:00&metric_time_to=2019-08-03T00:00:00&group_by=minutes&metric_name=backend_response
// 	group_by : minutes, hours, days, weeks, months
func GetEventsHttp(w http.ResponseWriter, r *http.Request, eng *engine.Engine) {
	response := &rest.Resp{}
	defer response.Render(w)

	engineRequest := &Event{}
	var err error
	vars := mux.Vars(r)

	cast := &CheckAndCast{
		[]string{
			vars["metricTimeFrom"],
			vars["metricTimeTo"],
			vars["groupBy"],
			vars["metricName"],
			vars["project"],
		},
		[]CACFunc{
			CACTimeFrom,
			CACTimeTo,
			CACGroupBy,
			CACMetricName,
			CACProject,
		},
		err, eng,
	}
	cast.execute(engineRequest)
	eng.Logger.Printf("[debug] request matched")
	if cast.err != nil {
		// TODO: return error
		eng.Logger.Fatal(cast.err)
	} else {
		events, err := GetEvents(engineRequest, eng)
		if err != nil {
			eng.Logger.Fatal(err)
			return
		}
		response.JsonBody, err = json.Marshal(events)
		if err != nil {
			eng.Logger.Fatal(err)
			return
		}
		response.Status = rest.StatusOk
	}
}

func SaveEvent(event *Event, eng *engine.Engine) (err error) {
	dbEngine := eng.DBEngine
	err = dbEngine.Collection(collection).Insert(event)
	return
}

///////////////////////////////////////////////////////////////
//
// Cast and check functions
//
///////////////////////////////////////////////////////////////

func CACTimeFrom(value string, fe *Event, eng *engine.Engine) (err error) {
	fe.MetricTimeFrom, err = time.Parse(timeLayout, value)
	return
}

func CACTimeTo(value string, fe *Event, eng *engine.Engine) (err error) {
	fe.MetricTimeTo, err = time.Parse(timeLayout, value)
	return
}

func CACMetricTimeNow(value string, fe *Event, eng *engine.Engine) (err error) {
	fe.MetricTime = time.Now()
	return
}

func CACGroupBy(value string, fe *Event, eng *engine.Engine) (err error) {
	switch value {
	case "minutes":
		fe.GroupBy = GroupByMins
	case "hours":
		fe.GroupBy = GroupByHours
	case "days":
		fe.GroupBy = GroupByDays
	case "weeks":
		fe.GroupBy = GroupByWeeks
	case "months":
		fe.GroupBy = GroupByMonths
	default:
		// TODO: check escaping!
		err = errors.New(fmt.Sprintf("Range %s is not supported", value))
	}
	return
}

func CACMetricName(value string, fe *Event, eng *engine.Engine) (err error) {
	fe.MetricName = value
	return
}

func CACProject(value string, fe *Event, eng *engine.Engine) (err error) {
	fe.ProjectId, err = projects.ProjectExists(value, eng)
	return
}

/////////////////////////////////////////////////////////
//
// Usage:
//	values — slice of string values to cast
//	fns — slice of cast functions
//
// every item in values slice will be processed by cast
// function from fns slice in same order
//

type CheckAndCast struct {
	values []string
	fns    []CACFunc
	err    error
	eng    *engine.Engine
}

func (rp *CheckAndCast) execute(target *Event) {
	if len(rp.values) != len(rp.fns) {
		panic("count of values must mach fns")
	}
	for i, fn := range rp.fns {
		if rp.err != nil {
			return
		}
		rp.err = fn(rp.values[i], target, rp.eng)
	}
}
