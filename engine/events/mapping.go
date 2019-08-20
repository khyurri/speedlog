package events

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine"
	"net/http"
)

func MapEventToFilter(e *Event) (filter *Filter) {
	filter = &Filter{
		MetricName:     e.MetricName,
		ProjectId:      e.ProjectId,
		MetricTimeFrom: e.MetricTimeFrom,
		MetricTimeTo:   e.MetricTimeTo,
	}
	return
}

func MapSaveRequestToEvent(r *http.Request, eng *engine.Engine) (event *Event, err error) {
	event = &Event{}
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&event)
	eng.Logger.Println(event.MetricName)
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
		err = fn(k, event, eng)
	}
	return
}
