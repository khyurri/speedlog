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
	equals = testutils.Equals
)

func TestLoadPlugins(t *testing.T) {

	var sigStopped sync.WaitGroup
	sendSigStop := make(SigChan)

	dbEng := &mongo.DataStoreMock{}
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
