package mongo

import (
	"fmt"
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
	for i := 0; i < 10; i++ {
		err = mongo.SaveEvent("ev0", "myproject", 100+float64(i))
		err = mongo.SaveEvent("ev1", "myproject", 90+float64(i))
	}
	now := time.Now()
	loc, _ := time.LoadLocation("Europe/Moscow")
	from := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, loc)
	to := from.Add(time.Minute)
	events, err := mongo.AllEvents(from, to)
	if err != nil {
		t.Logf("[error] %s\n", err)
	}
	for _, group := range events {
		aggregated, _ := GroupBy("minutes", group.Events)
		fmt.Println(aggregated[0].Event.ProjectId)
		fmt.Println(aggregated[0].Event.MetricName)
	}
}
