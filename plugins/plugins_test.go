package plugins

import (
	"github.com/khyurri/speedlog/engine/mongo"
	"github.com/khyurri/speedlog/testutils"
	"sync"
	"testing"
	"time"
)

// wrappers to shorten name
var (
	//ok     = testutils.Ok
	//assert = testutils.Assert
	equals = testutils.Equals
)

type dataStoreMock struct {
	DelEventsCalledTimes int
	AllEventsCalledTimes int
}

func (d *dataStoreMock) FilterEvents(from, to time.Time, metricName, project string) (events []mongo.Event, err error) {
	panic("implement me")
}

func (d *dataStoreMock) AllEvents(from, to time.Time) (events []mongo.AllEvents, err error) {
	if d.AllEventsCalledTimes == 0 {
		d.AllEventsCalledTimes = 1
	} else {
		d.AllEventsCalledTimes++
	}
	return
}

func (d *dataStoreMock) SaveEvent(metricName, project string, durationMs float64) (err error) {
	panic("implement me")
}

func (d *dataStoreMock) DelEvents(to time.Time) (err error) {
	if d.DelEventsCalledTimes == 0 {
		d.DelEventsCalledTimes = 1
	} else {
		d.DelEventsCalledTimes++
	}
	return nil
}

func (d *dataStoreMock) AddUser(login string, password string) (err error) {
	panic("implement me")
}

func (d *dataStoreMock) GetUser(login string) (*mongo.User, error) {
	panic("implement me")
}

func (d *dataStoreMock) UserDel(uid string) error {
	panic("implement me")
}

func (d *dataStoreMock) AddProject(title string) error {
	panic("implement me")
}

func (d *dataStoreMock) GetProject(title string) (project mongo.Project, err error) {
	panic("implement me")
}

func (d *dataStoreMock) GetProjectById(id string) (project mongo.Project, err error) {
	panic("implement me")
}

func (d *dataStoreMock) DelProject(id string) (err error) {
	panic("implement me")
}

func TestLoadPlugins(t *testing.T) {

	var sigStopped sync.WaitGroup
	sendSigStop := make(SigChan)

	dbEng := &dataStoreMock{}
	var plugins []Plugin

	// Cleaner plugin
	interval := time.Second * 1
	ttl := 60 * 60 * 24 * 100 // delete all events older than 100 days
	cl := NewCleaner(ttl, interval)
	plugins = append(plugins, cl)

	// Graphite plugin
	gr := NewGraphite("127.0.0.1", interval)
	plugins = append(plugins, gr)

	go LoadPlugins(plugins, sendSigStop, &sigStopped, dbEng)

	// send stop signal
	sendSigStop <- struct{}{}
	sigStopped.Wait()
}
