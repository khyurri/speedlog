package plugins

import (
	"fmt"
	"github.com/khyurri/speedlog/engine/mongo"
	"sync"
	"time"
)

type cleaner struct {
	ttl      int
	interval time.Duration
	ticker   *time.Ticker
}

// NewCleaner return cleaner plugin
// ttl — deletes all events that arrived before the seconds in the ttl
// interval — run every interval
func NewCleaner(ttl int, interval time.Duration) *cleaner {
	return &cleaner{
		ttl:      ttl,
		interval: interval,
	}
}

// Load loads and runs plugin
// runs cleaner on load
func (c *cleaner) Load(dbEngine mongo.DataStore, sigStop SigChan, sigStopped *sync.WaitGroup) {

	go func() {
		c.clean(dbEngine)
		c.ticker = time.NewTicker(c.interval)
		for _ = range c.ticker.C {
			select {
			default:
				c.clean(dbEngine)
			case <-sigStop:
				c.ticker.Stop()
				sigStopped.Done()
			}
		}
	}()
}

func (c *cleaner) clean(dbEngine mongo.DataStore) {
	now := time.Now()
	then := now.Add(time.Duration(-1*c.ttl) * time.Second)
	err := dbEngine.DelEvents(then)
	if err != nil {
		fmt.Printf("[error] %v\n", err)
	}
}
