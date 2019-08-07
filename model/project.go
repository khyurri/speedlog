package model

import (
	"errors"
	"github.com/globalsign/mgo/bson"
)

type Project struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Title string
}

func (ds *DataStore) RegisterProject(title string) bool {
	t := &Project{Title: title}
	col := ds.Collection("project")
	err := col.Insert(t)
	if err != nil {
		panic(err)
	}
	return true
}

func (ds *DataStore) ProjectExists(title string) (res bson.ObjectId, err error) {
	t := Project{}
	err = ds.Collection("project").Find(bson.M{"title": title}).One(&t)
	if err != nil {
		if err.Error() == "not found" {
			return res, errors.New("no s")
		} else {
			panic(err)
		}
	} else {
		res = t.ID
	}
	return t.ID, nil
}
