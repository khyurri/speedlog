package model

import (
	"github.com/globalsign/mgo"
)

const (
	Database = "speedlog"
	ColEvent = "event"
)

type DataStore struct {
	Session *mgo.Session
	db      string
}

func NewDataStore() *DataStore {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	return &DataStore{
		Session: session,
		db:      Database,
	}
}

func (ds *DataStore) Clone() *DataStore {
	session := ds.Session.Clone()
	return &DataStore{
		Session: session,
		db:      ds.db,
	}
}

func (ds *DataStore) Close() {
	ds.Session.Close()
}

func (ds *DataStore) Collection(collection string) *mgo.Collection {
	return ds.Session.DB(ds.db).C(collection)
}
