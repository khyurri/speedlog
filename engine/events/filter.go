package events

import (
	"github.com/globalsign/mgo/bson"
	"github.com/khyurri/speedlog/engine"
	"github.com/montanaflynn/stats"
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

func GroupBy(group string, events []Event, eng *engine.Engine) (result []*FilteredEvent, err error) {
	result = mapEvent(events)
	return
}

func mapEvent(event []Event) (result []*FilteredEvent) {
	// event should be ordered
	if len(event) == 0 {
		return
	}
	year := event[0].MetricTime.Year()
	month := event[0].MetricTime.Month()
	day := event[0].MetricTime.Day()
	fe := &FilteredEvent{
		MetricName: event[0].MetricName,
		MetricTime: time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
		ProjectId:  event[0].ProjectId,
	}
	// TODO: simplify
	for _, e := range event {
		if year != e.MetricTime.Year() &&
			month != e.MetricTime.Month() &&
			day != e.MetricTime.Day() {
			// TODO: bug here
			result = append(result, fe)
			collapse(fe)

			year := e.MetricTime.Year()
			month := e.MetricTime.Month()
			day := e.MetricTime.Day()

			fe = &FilteredEvent{
				MetricName: e.MetricName,
				MetricTime: time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
				ProjectId:  e.ProjectId,
			}
		}
		fe.durationsMs = append(fe.durationsMs, e.DurationMs)
	}
	return
}

func collapse(src *FilteredEvent) {
	src.MaxDurationMs, _ = stats.Max(src.durationsMs)
	src.MinDurationMs, _ = stats.Min(src.durationsMs)
	src.MedianDurationMs, _ = stats.Median(src.durationsMs)
	src.MiddleDurationMs = average(src.durationsMs)
}
