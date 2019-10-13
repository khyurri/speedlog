package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine/mongo"
	"github.com/khyurri/speedlog/utils"
	"net/http"
	"time"
)

const (
	timeLayout = "2006-01-02T15:04:05"
)

var debug = utils.Debug

// badRequest returns true if StatusErr is set
func badRequest(err error, r *Resp) bool {
	if err != nil {
		utils.Ok(err)
		r.Status = StatusErr
		return true
	}
	return false
}

type eventRequest struct {
	Id         int     `json:"id"`
	MetricName string  `json:"metricName"`
	DurationMs float64 `json:"durationMs"`
	Project    string  `json:"project"`
}

type eventsRequest []eventRequest

func (env *Env) createEventsHttp() http.HandlerFunc {

	type respMessage struct {
		Id   int `json:"id"`
		Code int `json:"code"`
	}

	mapRequestToList := func(r *http.Request, target *eventsRequest) (err error) {
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(target)
		if err != nil {
			utils.Ok(fmt.Errorf(err.Error()+". Body: %+v", r.Body))
			return
		}
		return
	}

	return func(w http.ResponseWriter, r *http.Request) {
		response := &Resp{}
		defer response.Render(w)

		req := &eventsRequest{}
		err := mapRequestToList(r, req)
		if err != nil {
			utils.Ok(err)
			response.Status = StatusIntErr
			return
		}
		var respMessages []respMessage
		if len(*req) > 0 {
			for _, event := range *req {
				err := env.DBEngine.SaveEvent(
					event.MetricName, event.Project, event.DurationMs)
				if err != nil {
					respMessages = append(respMessages, respMessage{
						Id:   event.Id,
						Code: http.StatusBadRequest,
					})
				} else {
					respMessages = append(respMessages, respMessage{
						Id:   event.Id,
						Code: http.StatusOK,
					})
				}

			}
		}
		response.JsonBody, err = json.Marshal(respMessages)
		response.Status = StatusOk
	}
}

func (env *Env) createEventHttp() http.HandlerFunc {

	mapRequestToStruct := func(r *http.Request, target *eventRequest) (err error) {
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(target)
		if err != nil {
			utils.Ok(fmt.Errorf(err.Error()+". Body: %+v", r.Body))
			return
		}
		if len(target.MetricName) == 0 {
			return errors.New("empty metricName")
		}
		return
	}

	return func(w http.ResponseWriter, r *http.Request) {

		response := &Resp{}
		defer response.Render(w)

		req := &eventRequest{}
		err := mapRequestToStruct(r, req)

		if err != nil {
			utils.Ok(err)
			response.Status = StatusIntErr
			return
		}

		err = env.DBEngine.SaveEvent(req.MetricName, req.Project, req.DurationMs)
		if err != nil {
			utils.Ok(err)
			response.Status = StatusIntErr
			return
		}

		utils.Debug(fmt.Sprintf("requested params: %s", r.Body))
		saved := struct {
			Saved bool `json:"saved"`
		}{true}
		response.Status = StatusOk
		response.JsonBody, err = json.Marshal(saved)
		utils.Ok(err)
	}
}

func (env *Env) getEventsHttp() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		response := &Resp{}
		defer response.Render(w)

		// see eventRequest validation in env.ExportEventRoutes
		params := mux.Vars(r)

		metricTimeFrom, err := time.Parse(timeLayout, params["metricTimeFrom"])
		if badRequest(err, response) {
			return
		}

		metricTimeTo, err := time.Parse(timeLayout, params["metricTimeTo"])
		if badRequest(err, response) {
			return
		}

		utils.Debug(fmt.Sprintf("%s -> %s", metricTimeFrom, metricTimeTo))
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
