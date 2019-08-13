package events

import (
	"github.com/globalsign/mgo/bson"
	"github.com/khyurri/speedlog/engine"
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
