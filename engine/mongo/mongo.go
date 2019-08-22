package mongo

import (
	"github.com/globalsign/mgo"
)

const (
	userCollection    = "user"
	projectCollection = "project"
	eventCollection   = "event"
)

type DataStore interface {
	AddUser(login string, password string) (err error)
	Authenticate(login string, password string) (err error)
	ProjectExists(title string) (projectId string, err error)
	RegisterProject(title string) bool
	FilterEvents(req *Filter) (events []*AggregatedEvent, err error)
	GroupBy(group string, events []*Event) (result []*AggregatedEvent, err error)
	SaveEvent(metricName, projectId string, durationMs float64) (err error)
}

type Mongo struct {
	Session *mgo.Session
	DbName  string
}

func New(db string, url string) (engine *Mongo, err error) {
	engine = &Mongo{DbName: db}
	engine.Session, err = mgo.Dial(url)
	return engine, err
}

func (mg *Mongo) Clone() *mgo.Session {
	return mg.Session.Clone()
}

func (mg *Mongo) Collection(collection string, sess *mgo.Session) *mgo.Collection {
	if sess == nil {
		sess = mg.Session
	}
	return sess.DB(mg.DbName).C(collection)
}

// DropDatabase - !run only for testing!
func (mg *Mongo) DropDatabase() error {
	return mg.Session.DB(mg.DbName).DropDatabase()
}
