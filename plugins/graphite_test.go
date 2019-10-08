package plugins

import (
	"github.com/khyurri/speedlog/engine/mongo"
	"sync"
	"testing"
	"time"
)

func TestGraphite_Load(t *testing.T) {

	host := "127.0.0.1"
	sigStop := make(SigChan, 2)
	interval := time.Second
	testInterval := time.Millisecond * 2500 // 2.5 additional second

	var sigStopped sync.WaitGroup
	sigStopped.Add(1)
	gr := NewGraphite(host, interval)
	dbEngine := &mongo.DataStoreMock{}
	gr.Load(dbEngine, sigStop, &sigStopped)
	time.Sleep(testInterval)
	sigStop <- struct{}{}
	sigStopped.Wait()
	equals(t, 2, mongo.AllEventsCalledTimes)

}
