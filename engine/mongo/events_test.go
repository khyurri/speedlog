package mongo

import (
	"testing"
	"time"
)

func TestAllEvents(t *testing.T) {
	mongo, err := New("speedlog", "localhost:27017")
	if err != nil {
		t.Logf("[error] %s\n", err)
		t.Fail()
	}
	// Create some events
	testEvents := [2]string{"ev0", "ev1"}
	for _, e := range testEvents {
		for i := 0; i < 10; i++ {
			_ = mongo.SaveEvent(e, "myproject", 100+float64(i))
		}
	}
	now := time.Now()
	loc, _ := time.LoadLocation("Europe/Moscow")
	from := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, loc)
	to := from.Add(time.Minute)
	events, err := mongo.AllEvents(from, to)
	if err != nil {
		t.Logf("[error] %s\n", err)
	}
	for i, group := range events {
		aggregated, _ := GroupBy("minutes", group.Events)
		// todo: aggregated order is not constant!
		if testEvents[i] != aggregated[0].Event.MetricName {
			t.Failed()
		}
	}
}
