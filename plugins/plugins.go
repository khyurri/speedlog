package plugins

import (
	"fmt"
	"github.com/khyurri/speedlog/engine/mongo"
	"sync"
)

type SigChan chan interface{}

type Plugin interface {
	Load(dbEngine mongo.DataStore, sigStop SigChan, sigStopped *sync.WaitGroup)
}

// LoadPlugins - loads plugins and controls their work
func LoadPlugins(plugins []Plugin, rcvSigStop SigChan, sigStopped *sync.WaitGroup, dbEngine mongo.DataStore) {
	if len(plugins) > 0 {
		sendSigStop := make(SigChan, len(plugins))
		for _, plugin := range plugins {
			sigStopped.Add(1)
			plugin.Load(dbEngine, sendSigStop, sigStopped)
		}
		fmt.Println("HERE!")
		<-rcvSigStop
		fmt.Println("RECEIVED")
		for range plugins {
			sendSigStop <- struct{}{}
		}
	}
}
