package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/khyurri/speedlog/utils"
)

type Project struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Title string        `bson:"title"`
}

func (mg *Mongo) GetProject(title string) (project Project, err error) {

	sess := mg.Clone()
	defer sess.Close()

	err = mg.Collection(projectCollection, sess).Find(bson.M{"title": title}).One(&project)
	if err != nil {
		return
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
	utils.Ok(err)

	// Delete project
	err = mg.Collection(projectCollection, sess).Remove(bson.M{
		"_id": bson.ObjectIdHex(id),
	})
	utils.Ok(err)

	return
}

func (mg *Mongo) GetProjectById(id string) (project Project, err error) {
	sess := mg.Clone()
	defer sess.Close()

	err = mg.Collection(projectCollection, sess).Find(bson.M{
		"_id": bson.ObjectIdHex(id),
	}).One(&project)
	if err != nil {
		fmt.Printf("[error] %v\n", err)
	}
	return project, err
}
