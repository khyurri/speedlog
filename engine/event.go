package engine

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine/mongo"
	"net/http"
	"time"
)

const (
	timeLayout = "2006-01-02T15:04"
)

func (env *Env) saveEventHttp() http.HandlerFunc {

	type request struct {
		MetricName string  `json:"metricName"`
		DurationMs float64 `json:"durationMs"`
		Project    string  `json:"project"`
	}

	mapRequestToStruct := func(r *http.Request, target *request) (err error) {
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(target)
		if err != nil {
			Logger.Printf("[error] deconding request error: %s", err)
			Logger.Printf("[debug] request body: %s", r.Body)
			return
		}
		Logger.Printf("[debug] metricName: %s, durationMs: %f",
			target.MetricName, target.DurationMs)
		if len(target.MetricName) == 0 {
			return errors.New("empty metricName")
		}
		return
	}

	return func(w http.ResponseWriter, r *http.Request) {

		response := &Resp{}
		defer response.Render(w)

		req := &request{}
		err := mapRequestToStruct(r, req)

		if err != nil {
			Logger.Printf("[debug] internal error %s", err)
			response.Status = StatusIntErr
			return
		}

		err = env.DBEngine.SaveEvent(req.MetricName, req.Project, req.DurationMs)
		if err == nil {
			Logger.Printf("[debug] requested params: %s", r.Body)
			saved := struct {
				Saved bool `json:"saved"`
			}{true}
			response.Status = StatusOk
			response.JsonBody, err = json.Marshal(saved)
		}
		if err != nil {
			Logger.Printf("[error] %s", err)
			response.Status = StatusIntErr
			return
		}
	}

}

func (env *Env) getEventsHttp() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		response := &Resp{}
		defer response.Render(w)

		// see request validation in env.ExportEventRoutes
		params := mux.Vars(r)

		metricTimeFrom, err := time.Parse(timeLayout, params["metricTimeFrom"])
		if err != nil {
			Logger.Printf("[error] %s", err)
		}

		metricTimeTo, err := time.Parse(timeLayout, params["metricTimeTo"])
		if err != nil {
			Logger.Printf("[error] %s", err)
		}

		if err != nil {
			response.Status = StatusErr
		}

		Logger.Printf("[debug] %s -> %s", metricTimeFrom, metricTimeTo)
		events, err := env.DBEngine.FilterEvents(
			metricTimeFrom,
			metricTimeTo,
			params["metricName"],
			params["project"])

		grouped, err := mongo.GroupBy(params["groupBy"], events)
		if len(grouped) == 0 {
			response.Status = StatusOk
			return
		}
		response.JsonBody, err = json.Marshal(grouped)
	}
}
