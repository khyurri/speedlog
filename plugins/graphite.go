package plugins

import (
	"fmt"
	"github.com/khyurri/speedlog/engine/mongo"
	"time"
)

type graphite struct {
	host   string
	ticker *time.Ticker
}

func NewGraphite(host string) *graphite {
	return &graphite{host: host}
}

func (gr *graphite) Load(dbEngine mongo.DataStore) {
	gr.ticker = time.NewTicker(10 * time.Second)
	go func() {
		for t := range gr.ticker.C {
			fmt.Println("Tick at", t)
		}
	}()
}

func (gr *graphite) Unload() {
	gr.ticker.Stop()
}
