package app

import (
	"../model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type ProjectItem struct {
	Title string
}

type LogItemRequest struct {
	MetricName string
	DurationMs int
}

type App struct {
	ds *model.DataStore
}

func NewApp() *App {
	dataStore := model.NewDataStore()
	return &App{dataStore}
}

func (a *App) dataStore() *model.DataStore {
	return a.ds.Clone()
}

// PUT /log/{key}/
func (a *App) WriteLogHandler(w http.ResponseWriter, r *http.Request) {
	ds := a.dataStore()
	defer ds.Close()

	////////////////////////////////////////////
	// CHECK REQUEST
	// TODO: check request
	//

	// Decode Request
	var logItem LogItemRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&logItem)

	if err != nil {
		panic(err)
	}

	project := mux.Vars(r)["key"]
	//
	// END OF CHECK REQUEST
	////////////////////////////////////////////

	projectId, err := ds.ProjectExists(project)
	if err != nil {
		_, err := fmt.Fprintf(w, "{'status': 'error', 'message': 'project not found'}")
		if err != nil {
			panic(err)
		}
	} else {
		ds.RegisterEvent(projectId, logItem.MetricName, logItem.DurationMs)
		_, err := fmt.Fprintf(w, "{'status': 'success'}")
		if err != nil {
			panic(err)
		}
	}
}

// PUT /project/
func (a *App) RegisterProjectHandler(w http.ResponseWriter, r *http.Request) {
	ds := a.dataStore()
	defer ds.Close()

	var t ProjectItem
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	ds.RegisterProject(t.Title)
}

/////////////////////////////////////////////////////////
//
// This struct casts request to ExportStats structure
// Usage:
//	values — slice of string values to cast
//	fns — slice of cast functions
//
// every item in values slice will be processed by cast
// function from fns slice in same order
//
type castReqToES struct {
	values []string
	fns    []func(string, *model.ExportStats) error
	err    error
}

func (rp *castReqToES) parse(target *model.ExportStats) {
	if len(rp.values) != len(rp.fns) {
		panic("count of values must mach fns")
	}
	for i, fn := range rp.fns {
		if rp.err != nil {
			return
		}
		rp.err = fn(rp.values[i], target)
	}
}

/////////////////////////////////////////////////////////
//
// cast functions
//

func castMetricName(value string, es *model.ExportStats) (err error) {
	es.MetricName = value
	return err
}

func castProject(value string, es *model.ExportStats) (err error) {

}

// GET 127.0.0.1:8012/log/stats/pravoved.ru/?time_from=2019-08-02T00:00:00&time_to=2019-08-03T00:00:00&group_by=minutes
func (a *App) StatsHandler(w http.ResponseWriter, r *http.Request) {
	ds := a.dataStore()
	defer ds.Close()
	vars := mux.Vars(r)
	project := vars["key"]

	projectId, err := ds.ProjectExists(project)
	if err != nil {
		_, err := fmt.Fprintf(w, "{'status': 'error', 'message': 'project not found'}")
		if err != nil {
			panic(err)
		}
	}

	dates, err := ds.StrToTime(vars["metricTimeFrom"], vars["metricTimeTo"])
	casted := &castReqToES{
		[]string{
			vars["metricName"],
		},
		[]func(string, *model.ExportStats) error{
			castMetricName,
		},
		err,
	}
	groupBy, err := model.StrToGroupBy(vars["GroupBy"])

	if err != nil {
		_, err := fmt.Fprintf(w, "{'status': 'error', 'message': ''}")
		if err != nil {
			panic(err)
		}
	} else {

		req := &model.ExportStats{
			Project:        projectId,
			MetricName:     vars["metricName"],
			MetricTimeFrom: dates[0],
			MetricTimeTo:   dates[1],
			GroupBy:        vars["GroupBy"],
		}
		fmt.Println(vars)
		fmt.Println(projectId)
		fmt.Println(project)
		//stats := ds.ExportStats()
	}
}
