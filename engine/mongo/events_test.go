package mongo

import (
	"testing"
	"time"
)

const (
	testDurationMinute = iota
	testDuration2Minutes
)

var testDuration = [...]time.Duration{
	time.Minute,
	2 * time.Minute,
}

func getDuration(d int) time.Duration {
	return testDuration[d]
}

func getDateFrom(t testing.TB, timeFrom string) (tm time.Time) {
	tm, err := time.Parse("2006-01-02 15:04", timeFrom)
	ok(t, err)
	return
}

type testRoundGroup struct {
	id               int // unique id to debug test round
	dateFrom         time.Time
	dateToInterval   time.Duration
	metricName       string
	groupBy          string // minutes, hours, days
	expectEvents     int    // count of all events at range from dateFrom to dateFrom + dateToInterval
	expectGroups     int    // count of project + event groups pairs. every pair is a new group
	aggregatedLen    int    // count of lists after GroupBy. depends on dateToInterval and groupBy values
	aggregatedEvents []testEventAggregated
}

type testEventAggregated struct {
	minDurationMs    float64
	maxDurationMs    float64
	medianDurationMs float64
	middleDurationMs float64
}

func TestGetAllEvents(t *testing.T) {
	mongo, err := New(testMongoDb, testMongoHost)
	ok(t, err)
	defer mongo.Session.Close()

	populateDb(t, mongo, eventsFixtures)

	// Setup testing values
	testRounds := []testRoundGroup{
		{
			id:       0,
			dateFrom: getDateFrom(t, "2019-08-30 00:01"),
			// all events gte:2019-08-30 00:01:00 lt:2019-08-30 00:02:00
			dateToInterval: getDuration(testDurationMinute),
			groupBy:        "minutes",
			expectEvents:   5,
			metricName:     "resp0",
			expectGroups:   1,
			aggregatedLen:  1, // duration 1 minute, group by minutes
			aggregatedEvents: []testEventAggregated{{
				minDurationMs:    0.1,
				maxDurationMs:    0.9,
				medianDurationMs: 0.3,
			}},
		},
		{
			id:             1,
			dateFrom:       getDateFrom(t, "2019-08-30 00:02"),
			dateToInterval: getDuration(testDurationMinute),
			groupBy:        "minutes",
			expectEvents:   6,
			metricName:     "resp0",
			expectGroups:   1,
			aggregatedLen:  1, // duration 1 minute, group by minutes
			aggregatedEvents: []testEventAggregated{{
				minDurationMs:    0.5,
				maxDurationMs:    2,
				medianDurationMs: 1.2000000000000002, // just hardcoded value
			}},
		},
	}

	// Run checks
	for _, round := range testRounds {
		dateFrom := round.dateFrom
		dateTo := dateFrom.Add(round.dateToInterval)
		// get all events
		events, err := mongo.AllEvents(dateFrom, dateTo)
		ok(t, err)
		for _, group := range events {
			t.Logf("[debug] running round %d", round.id)

			// check count of all events
			expect := round.expectEvents
			eventsCount := len(group.Events)
			equals(t, expect, eventsCount)

			// aggregation
			aggregated, err := GroupBy("minutes", group.Events)
			ok(t, err)
			equals(t, round.aggregatedLen, len(aggregated))
			for i, row := range aggregated {
				equals(t, round.metricName, row.MetricName)
				equals(t, round.aggregatedEvents[i].maxDurationMs, row.MaxDurationMs)
				equals(t, round.aggregatedEvents[i].minDurationMs, row.MinDurationMs)
				equals(t, round.aggregatedEvents[i].medianDurationMs, row.MedianDurationMs)
			}
		}
	}

}
