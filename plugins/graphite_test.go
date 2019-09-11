package plugins

import (
	"sync"
	"testing"
	"time"
)

func TestGraphite_Load(t *testing.T) {

	host := "127.0.0.1"
	sigStop := make(SigChan, 2)
	interval := time.Second
	testInterval := time.Second * 2 // 2 additional second

	var sigStopped sync.WaitGroup
	sigStopped.Add(1)
	gr := NewGraphite(host, interval)
	dbEngine := &dataStoreMock{}
	gr.Load(dbEngine, sigStop, &sigStopped)
	time.Sleep(testInterval)
	sigStop <- struct{}{}
	sigStopped.Wait()
	equals(t, 2, dbEngine.AllEventsCalledTimes)
}
