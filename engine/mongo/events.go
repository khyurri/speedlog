package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/montanaflynn/stats"
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
	MinDurationMs    float64
	MaxDurationMs    float64
	MedianDurationMs float64
	MiddleDurationMs float64
	EventCount       int
	durationsMs      []float64
}

func (mg *Mongo) saveEventAtTime(metricName, project string, durationMs float64, eventTime time.Time) (err error) {
	sess := mg.Clone()
	defer sess.Close()

	projectId, err := mg.GetProject(project)

	if err != nil {
		fmt.Printf("[error] fetch project by name: %s", err)
		return
	}

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

func (mg *Mongo) FilterEvents(from, to time.Time, metricName, project string) (events []Event, err error) {

	sess := mg.Clone()
	defer sess.Close()

	projectId, err := mg.GetProject(project)

	if err != nil {
		fmt.Printf("[error] fetch project by name: %s", err)
		return
	}

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
	//ordered(result)
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
	src.MaxDurationMs, _ = stats.Max(src.durationsMs)
	src.MinDurationMs, _ = stats.Min(src.durationsMs)
	src.MedianDurationMs, _ = stats.Median(src.durationsMs)
	src.MiddleDurationMs = average(src.durationsMs)
	src.EventCount = len(src.durationsMs)
}

//func ordered(srcs []*AggregatedEvent) {
//	sort.Slice(srcs, func(i, j int) bool {
//		return srcs[i].MetricTime.Before(srcs[j].MetricTime)
//	})
//}
