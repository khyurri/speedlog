package plugins

import (
	"sync"
	"testing"
	"time"
)

func TestCleaner_Load(t *testing.T) {
	sigStop := make(SigChan, 2)
	var sigStopped sync.WaitGroup
	interval := time.Second * 1
	testInterval := time.Second * 2 // 1 additional second
	ttl := 60 * 60 * 24 * 100       // delete all events older than 100 days

	sigStopped.Add(1)
	cl := NewCleaner(ttl, interval)
	dbEng := &dataStoreMock{}
	cl.Load(dbEng, sigStop, &sigStopped)
	time.Sleep(testInterval) // sleep 1 interval
	sigStop <- struct{}{}
	sigStopped.Wait()

	equals(t, 2, dbEng.DelEventsCalledTimes)
}
