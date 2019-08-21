package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/montanaflynn/stats"
	"sort"
	"time"
)

// Filter — mongodb request for filtering events
type Filter struct {
	MetricName     string        `bson:"metric_name"`
	ProjectId      bson.ObjectId `bson:"project_id"`
	MetricTimeFrom time.Time     `bson:"metric_time_from,omitempty"`
	MetricTimeTo   time.Time     `bson:"metric_time_to,omitempty"`
}

// todo: refactor
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

// todo: refactor
type FilteredEvents []*FilteredEvent

func (mg *Mongo) SaveEvent(event interface{}) (err error) {

	sess := mg.Clone()
	defer sess.Close()

	err = mg.Collection(eventCollection, sess).Insert(event)
	return
}

func (mg *Mongo) FilterEvents(req *Filter) (events []interface{}, err error) {

	sess := mg.Clone()
	defer sess.Close()

	events = make([]interface{}, 0)
	// TODO: check req can be sent as request
	err = mg.Collection(eventCollection, sess).
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

type keyFunc func(eventTime time.Time) time.Time

func average(items []float64) float64 {
	var acc float64 = 0
	for _, i := range items {
		acc += i
	}
	return acc / float64(len(items))
}

func (mg *Mongo) GroupBy(group string,
	events []struct {
		MetricName string
		ProjectId  bson.ObjectId
		MetricTime time.Time
		DurationMs float64
	}) (result FilteredEvents, err error) {
	m := map[string]keyFunc{
		"minutes": groupByMinutes,
		"hours":   groupByHours,
		"days":    groupByDays,
	}
	if nil == m[group] {
		return result, fmt.Errorf("unsupported group key %s", group)
	}
	result = mapEvent(events, m[group])
	ordered(result)
	return
}

func groupByMinutes(eventTime time.Time) time.Time {
	return time.Date(
		eventTime.Year(),
		eventTime.Month(),
		eventTime.Day(),
		eventTime.Hour(),
		eventTime.Minute(), 0, 0, time.UTC)
}

func groupByHours(eventTime time.Time) time.Time {
	return time.Date(
		eventTime.Year(),
		eventTime.Month(),
		eventTime.Day(),
		eventTime.Hour(),
		0, 0, 0, time.UTC)
}

func groupByDays(eventTime time.Time) time.Time {
	return time.Date(
		eventTime.Year(),
		eventTime.Month(),
		eventTime.Day(), 0, 0, 0, 0, time.UTC)
}

func mapEvent(
	event []struct {
		MetricName string
		ProjectId  bson.ObjectId
		MetricTime time.Time
		DurationMs float64
	}, keyFunc keyFunc) (result FilteredEvents) {

	if len(event) == 0 {
		return
	}

	m := map[time.Time]*FilteredEvent{}

	for _, e := range event {

		key := keyFunc(e.MetricTime)
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

func ordered(srcs FilteredEvents) {
	sort.Slice(srcs, func(i, j int) bool {
		return srcs[i].MetricTime.Before(srcs[j].MetricTime)
	})
}
