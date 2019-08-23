package mongo

import (
	"github.com/globalsign/mgo/bson"
)

type Project struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Title string        `bson:"title"`
}

func (mg *Mongo) GetProject(title string) (projectId string, err error) {

	sess := mg.Clone()
	defer sess.Close()

	t := Project{}
	err = mg.Collection(projectCollection, sess).Find(bson.M{"title": title}).One(&t)
	if err != nil {
		return
	} else {
		projectId = t.ID.Hex()
	}
	return
}

func (mg *Mongo) AddProject(title string) (err error) {

	sess := mg.Clone()
	defer sess.Close()

	t := &Project{Title: title}
	err = mg.Collection(projectCollection, sess).Insert(t)

	return
}

func (mg *Mongo) DelProject(id string) (err error) {
	sess := mg.Clone()
	defer sess.Close()

	// Delete events
	err = mg.delAllEvents(id)
	if err != nil {
		logger.Printf("[error] cannot delete events by project id %s: %s", id, err)
	}
	// Delete project
	err = mg.Collection(projectCollection, sess).Remove(bson.M{
		"_id": bson.ObjectIdHex(id),
	})
	if err != nil {
		logger.Printf("[error] cannot delete project %s: %s", id, err)
	}

	return
}
