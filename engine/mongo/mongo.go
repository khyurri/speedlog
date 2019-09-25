package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/khyurri/speedlog/utils"
	"sync"
	"time"
)

var once sync.Once

const (
	userCollection    = "User"
	projectCollection = "project"
	eventCollection   = "event"
)

type DataStore interface {
	FilterEvents(from, to time.Time, metricName, project string) (events []Event, err error)
	AllEvents(from, to time.Time) (events []AllEvents, err error)
	SaveEvent(metricName, project string, durationMs float64) (err error)
	DelEvents(to time.Time) (err error)

	AddUser(login string, password string) (err error)
	GetUser(login string) (*User, error)
	UserDel(uid string) error

	AddProject(title string) error
	GetProject(title string) (project Project, err error)
	GetProjectById(id string) (project Project, err error)
	DelProject(id string) (err error)
}

type Mongo struct {
	Session *mgo.Session
	DbName  string
}

func New(db string, url string) (engine *Mongo, err error) {
	engine = &Mongo{DbName: db}
	engine.Session, err = mgo.Dial(url)
	if err != nil {
		return
	}
	err = engine.CreateIndexes()
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

func (mg *Mongo) CreateIndexes() (err error) {
	once.Do(func() {

		utils.Debug("starting index check")

		coll := mg.Collection(userCollection, nil)
		err = coll.EnsureIndex(mgo.Index{
			Key:    []string{"login"},
			Unique: true,
		})

		utils.Ok(err)
		utils.Debug("collection `project`")

		coll = mg.Collection(projectCollection, nil)
		err = coll.EnsureIndex(mgo.Index{
			Key:    []string{"title"},
			Unique: true,
		})
		utils.Ok(err)
		utils.Debug("index check complete")

	})
	return
}
