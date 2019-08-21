package mongo

import (
	"github.com/globalsign/mgo/bson"
)

type Project struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Title string        `bson:"title"`
}

func (mg *Mongo) ProjectExists(title string) (projectId bson.ObjectId, err error) {

	sess := mg.Clone()
	defer sess.Close()

	t := Project{}
	err = mg.Collection(projectCollection, sess).Find(bson.M{"title": title}).One(&t)
	if err != nil {
		return
	} else {
		projectId = t.ID
	}
	return
}

func (mg *Mongo) RegisterProject(title string) bool {

	sess := mg.Clone()
	defer sess.Close()

	t := &Project{Title: title}
	_ = mg.Collection(projectCollection, sess).Insert(t)
	return true
}
