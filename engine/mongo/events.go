package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/montanaflynn/stats"
	"sort"
	"time"
)

// Event - mongo event document
type Event struct {
	MetricName string        `bson:"metricName"`
	MetricTime time.Time     `bson:"metricTime,omitempty"`
	ProjectId  bson.ObjectId `bson:"projectId"`
	DurationMs float64       `bson:"durationMs,omitempty"`
}

type AllEvents struct {
	Meta struct {
		ProjectId  string `bson:"projectId"`
		MetricName string `bson:"metricName"`
	} `bson:"_id"`
	Events []Event `bson:"events"`
}

// AggregatedEvent - aggregated mongo event document
type AggregatedEvent struct {
	MetricName       string        `bson:"metricName"`
	ProjectId        bson.ObjectId `bson:"projectId"`
	MetricTime       time.Time     // metric time up to a minute
	MinDurationMs    float64
	MaxDurationMs    float64
	MedianDurationMs float64
	MiddleDurationMs float64
	EventCount       int
	Percentile90     float64
	Percentile75     float64
	durationsMs      []float64
}

func (mg *Mongo) saveEventAtTime(metricName, projectTitle string, durationMs float64, eventTime time.Time) (err error) {
	sess := mg.Clone()
	defer sess.Close()

	project, err := mg.GetProject(projectTitle)

	if err != nil {
		fmt.Printf("[error] fetch project by name: %s", projectTitle)
		return
	}

	projectId := project.ID.Hex()

	err = mg.Collection(eventCollection, sess).Insert(struct {
		MetricName string    `bson:"metricName"`
		ProjectId  string    `bson:"projectId"`
		DurationMs float64   `bson:"durationMs"`
		MetricTime time.Time `bson:"metricTime"`
	}{
		metricName,
		projectId,
		durationMs,
		eventTime,
	})
	return
}

func (mg *Mongo) SaveEvent(metricName, project string, durationMs float64) (err error) {
	return mg.saveEventAtTime(metricName, project, durationMs, time.Now())
}

func (mg *Mongo) DelEvents(to time.Time) (err error) {
	return
}

func (mg *Mongo) FilterEvents(from, to time.Time, metricName, projectTitle string) (events []Event, err error) {

	sess := mg.Clone()
	defer sess.Close()

	project, err := mg.GetProject(projectTitle)

	if err != nil {
		return
	}

	projectId := project.ID.Hex()

	events = make([]Event, 0)
	err = mg.Collection(eventCollection, sess).
		Find(bson.M{
			"projectId": projectId,
			"metricTime": bson.M{
				"$gte": from,
				"$lt":  to,
			},
			"metricName": metricName,
		}).Sort("metricTime").All(&events)
	return
}

func (mg *Mongo) AllEvents(from, to time.Time) (events []AllEvents, err error) {
	sess := mg.Session.Clone()
	defer sess.Close()

	events = make([]AllEvents, 0)

	c := mg.Collection(eventCollection, sess)
	pipe := c.Pipe([]bson.M{{
		"$match": bson.M{
			"metricTime": bson.M{
				"$gte": from,
				"$lt":  to,
			},
		},
	}, {
		"$group": bson.M{
			"_id": bson.M{
				"projectId":  "$projectId",
				"metricName": "$metricName",
			},
			"events": bson.M{"$push": "$$ROOT"},
		},
	},
	})
	iter := pipe.Iter()
	defer func() {
		err = iter.Close()
		if err != nil {
			fmt.Printf("[error] close iter: %s\n", err)
		}
	}()
	var grouped AllEvents

	for iter.Next(&grouped) {
		events = append(events, grouped)
	}
	if iter.Err() != nil {
		fmt.Printf("[error] iter: %s\n", err)
	}
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

func GroupBy(group string, events []Event) (result []*AggregatedEvent, err error) {
	m := map[string]keyFunc{
		"minutes": groupByMinutes,
	}
	if nil == m[group] {
		return result, fmt.Errorf("unsupported group key %s", group)
	}
	result = mapEvent(events, m[group])
	ordered(result)
	return
}

func (mg *Mongo) delAllEvents(projectId string) (err error) {
	sess := mg.Clone()
	defer sess.Close()

	err = mg.Collection(eventCollection, sess).Remove(bson.M{
		"projectId": projectId,
	})

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

func upToAMinute(t time.Time) time.Time {
	res := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	return res
}

func mapEvent(event []Event, keyFunc keyFunc) (result []*AggregatedEvent) {

	if len(event) == 0 {
		return
	}

	m := map[time.Time]*AggregatedEvent{}

	for _, e := range event {

		key := keyFunc(e.MetricTime)
		if nil == m[key] {

			m[key] = &AggregatedEvent{
				MetricName: e.MetricName,
				ProjectId:  e.ProjectId,
				MetricTime: upToAMinute(e.MetricTime),
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

func collapse(src *AggregatedEvent) {

	var err error

	src.MaxDurationMs, _ = stats.Max(src.durationsMs)
	src.MinDurationMs, _ = stats.Min(src.durationsMs)
	src.MedianDurationMs, _ = stats.Median(src.durationsMs)
	src.MiddleDurationMs = average(src.durationsMs)
	src.EventCount = len(src.durationsMs)
	src.Percentile90, err = stats.Percentile(src.durationsMs, 90)
	src.Percentile75, err = stats.Percentile(src.durationsMs, 75)
	if err != nil {
		fmt.Printf("[error] persentile error:  %v\n", err)
	}
}

func ordered(srcs []*AggregatedEvent) {
	fmt.Println(srcs)
	sort.Slice(srcs, func(i, j int) bool {
		return srcs[i].MetricTime.Before(srcs[j].MetricTime)
	})
}
