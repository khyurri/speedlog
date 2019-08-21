package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const (
	timeLayout = "2006-01-02T15:04:05"
)

const (
	GroupByMins   = iota
	GroupByHours  = iota
	GroupByDays   = iota
	GroupByWeeks  = iota
	GroupByMonths = iota
)

// Event — mongo event document
type Event struct {
	MetricName     string        `bson:"metric_name" json:"metric_name"`
	MetricTime     time.Time     `bson:"metric_time,omitempty" json:"metric_time"`
	ProjectId      bson.ObjectId `bson:"project_id" json:"project_id"`
	DurationMs     float64       `bson:"duration_ms,omitempty" json:"duration_ms,omitempty"`
	MetricTimeFrom time.Time     `bson:"metric_time_from,omitempty" json:"metric_time_from,omitempty"`
	MetricTimeTo   time.Time     `bson:"metric_time_to,omitempty" json:"metric_time_to,omitempty"`
	GroupBy        int           `bson:"group_by,omitempty" json:"group_by,omitempty"`
}

type SaveEventReq struct {
	MetricName string  `json:"metric_name"`
	DurationMs float64 `json:"duration_ms,omitempty"`
}

// CACFunc — check and cast function type
type CACFunc func(string, *Event, *Env) error

//////////////////////////////////////////////////////////////////////
//
// Events logic
//
//////////////////////////////////////////////////////////////////////

// SaveEventHttp PUT /pravoved.ru/event/
func (env *Env) SaveEventHttp(w http.ResponseWriter, r *http.Request) {
	var err error

	response := &Resp{}
	defer response.Render(w)

	event, err := MapSaveRequestToEvent(r, env)
	if err != nil {
		// TODO: make more readable error, please
		env.Logger.Println(err)
		response.Status = StatusIntErr
		return
	}

	err = env.DBEngine.SaveEvent(event)
	if err == nil {
		saved := struct {
			Saved bool `json:"saved"`
		}{true}
		response.Status = StatusOk
		response.JsonBody, err = json.Marshal(saved)
	}
	if err != nil {
		env.Logger.Fatal(err)
		response.Status = StatusIntErr
		return
	}

}

// GetEventsHttp GET /pravoved.ru/events/?metric_time_from=2019-08-02T00:00:00&metric_time_to=2019-08-03T00:00:00&group_by=minutes&metric_name=backend_response
// 	group_by : minutes, hours, days
func (env *Env) GetEventsHttp(w http.ResponseWriter, r *http.Request) {
	response := &Resp{}
	defer response.Render(w)

	engineRequest := &Event{}
	var err error
	vars := mux.Vars(r)
	env.Logger.Println(vars)
	// TODO: Simplify, create map function
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
		err, env,
	}
	cast.execute(engineRequest)
	env.Logger.Printf("[debug] request matched")
	if cast.err != nil {
		// TODO: return error
		env.Logger.Fatal(cast.err)
		return
	}
	filter := MapEventToFilter(engineRequest)

	// TODO: simplify
	events, err := env.DBEngine.FilterEvents(filter)
	if err != nil {
		env.Logger.Fatal(err)
		return
	}
	groupedEvents, err := env.DBEngine.GroupBy(vars["groupBy"], events)
	if err != nil {
		env.Logger.Fatal(err)
		return
	}
	response.JsonBody, err = json.Marshal(groupedEvents)
	if err != nil {
		env.Logger.Fatal(err)
		return
	}

	response.Status = StatusOk

}

///////////////////////////////////////////////////////////////
//
// Cast and check functions
//
///////////////////////////////////////////////////////////////

func CACTimeFrom(value string, fe *Event, env *Env) (err error) {
	fe.MetricTimeFrom, err = time.Parse(timeLayout, value)
	return
}

func CACTimeTo(value string, fe *Event, env *Env) (err error) {
	fe.MetricTimeTo, err = time.Parse(timeLayout, value)
	return
}

func CACMetricTimeNow(value string, fe *Event, env *Env) (err error) {
	fe.MetricTime = time.Now()
	return
}

func CACGroupBy(value string, fe *Event, env *Env) (err error) {
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

func CACMetricName(value string, fe *Event, env *Env) (err error) {
	env.Logger.Printf("[debug] save metric %s\n", value)
	fe.MetricName = value
	return
}

func CACProject(value string, fe *Event, env *Env) (err error) {
	env.Logger.Printf("[debug] looking for project %s\n", value)
	fe.ProjectId, err = env.DBEngine.ProjectExists(value)
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
	eng    *Env
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

func MapEventToFilter(e *Event) interface{} {
	return struct {
		MetricName     string
		ProjectId      bson.ObjectId
		MetricTimeFrom time.Time
		MetricTimeTo   time.Time
	}{
		e.MetricName,
		e.ProjectId,
		e.MetricTimeFrom,
		e.MetricTimeTo,
	}
}

func MapSaveRequestToEvent(r *http.Request, env *Env) (event *Event, err error) {
	event = &Event{}
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&event)
	env.Logger.Println(event.MetricName)
	if err != nil {
		return
	}
	// TODO: make CACFunc closure
	fns := []CACFunc{CACProject, CACMetricTimeNow}
	str := []string{vars["project"], ""}
	for i, fn := range fns {
		if err != nil {
			return
		}
		k := str[i]
		err = fn(k, event, env)
	}
	return
}
