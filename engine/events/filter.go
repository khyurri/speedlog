package events

import (
	"github.com/globalsign/mgo/bson"
	"github.com/khyurri/speedlog/engine"
	"github.com/montanaflynn/stats"
	"sort"
	"time"
)

type Filter struct {
	MetricName     string        `bson:"metric_name"`
	ProjectId      bson.ObjectId `bson:"project_id"`
	MetricTimeFrom time.Time     `bson:"metric_time_from,omitempty"`
	MetricTimeTo   time.Time     `bson:"metric_time_to,omitempty"`
}

type FilteredEvent struct {
	MetricName       string        `bson:"metric_name"`
	MetricTime       time.Time     `bson:"metric_time,omitempty"`
	ProjectId        bson.ObjectId `bson:"project_id"`
	MaxDurationMs    float64       `bson:"max_duration_ms,omitempty"`
	MinDurationMs    float64       `bson:"min_duration_ms,omitempty"`
	MedianDurationMs float64       `bson:"median_duration_ms,omitempty"`
	MiddleDurationMs float64       `bson:"middle_duration_ms,omitempty"`
	EventCount       int           `bson:event_count,omitempty"`
	durationsMs      []float64
}

type FilteredEvents []*FilteredEvent

func (req *Filter) FilterEvents(eng *engine.Engine) (events []Event, err error) {
	dbEngine := eng.DBEngine
	events = make([]Event, 0)
	// TODO: check req can be sent as request
	err = dbEngine.Collection(collection).
		Find(bson.M{
			"project_id": req.ProjectId,
			"metric_time": bson.M{
				"$gte": req.MetricTimeFrom,
				"$lt":  req.MetricTimeTo,
			},
			"metric_name": req.MetricName,
		}).Sort("metric_time").All(&events)
	return
}

func average(items []float64) float64 {
	var acc float64 = 0
	for _, i := range items {
		acc += i
	}
	return acc / float64(len(items))
}

func GroupBy(group string, events []Event, eng *engine.Engine) (result FilteredEvents, err error) {
	result = mapEvent(events)
	ordered(result)
	return
}

func mapEvent(event []Event) (result FilteredEvents) {

	if len(event) == 0 {
		return
	}

	m := map[time.Time]*FilteredEvent{}

	for _, e := range event {

		key := time.Date(
			e.MetricTime.Year(),
			e.MetricTime.Month(),
			e.MetricTime.Day(), 0, 0, 0, 0, time.UTC)

		if nil == m[key] {
			m[key] = &FilteredEvent{
				MetricName: e.MetricName,
				MetricTime: key,
				ProjectId:  e.ProjectId,
			}
		}

		m[key].durationsMs = append(m[key].durationsMs, e.DurationMs)
	}
	if len(m) > 0 {
		for _, e := range m {
			collapse(e)
			result = append(result, e)
		}
	}
	return
}

func collapse(src *FilteredEvent) {
	src.MaxDurationMs, _ = stats.Max(src.durationsMs)
	src.MinDurationMs, _ = stats.Min(src.durationsMs)
	src.MedianDurationMs, _ = stats.Median(src.durationsMs)
	src.MiddleDurationMs = average(src.durationsMs)
	src.EventCount = len(src.durationsMs)
}

func (o FilteredEvents) Len() int {
	return len(o)
}

func (o FilteredEvents) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o FilteredEvents) Less(i, j int) bool {
	return o[i].MetricTime.Before(o[j].MetricTime)
}

func ordered(srcs FilteredEvents) {
	sort.Sort(srcs)
}
