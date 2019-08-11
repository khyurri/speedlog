package events

import (
	"errors"
	"fmt"
	"github.com/khyurri/speedlog/engine"
	"github.com/montanaflynn/stats"
	"sort"
	"time"
)

type keyFunc func(eventTime time.Time) time.Time

func average(items []float64) float64 {
	var acc float64 = 0
	for _, i := range items {
		acc += i
	}
	return acc / float64(len(items))
}

func GroupBy(group string, events []Event, eng *engine.Engine) (result FilteredEvents, err error) {
	m := map[string]keyFunc{
		"minutes": groupByMinutes,
		"hours":   groupByHours,
		"days":    groupByDays,
	}
	if nil == m[group] {
		return result, errors.New(fmt.Sprintf("unsupported group key %s", group))
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

func mapEvent(event []Event, keyFunc keyFunc) (result FilteredEvents) {

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
