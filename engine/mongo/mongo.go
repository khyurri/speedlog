package mongo

import (
	"github.com/globalsign/mgo"
	"log"
	"os"
	"sync"
	"time"
)

var (
	logger *log.Logger
	once   sync.Once
)

const (
	userCollection    = "user"
	projectCollection = "project"
	eventCollection   = "event"
)

type DataStore interface {
	FilterEvents(from, to time.Time, metricName, project string) (events []*Event, err error)
	AllEvents(from, to time.Time) (events []*Event, err error)
	SaveEvent(metricName, project string, durationMs float64) (err error)

	AddUser(login string, password string) (err error)
	GetUser(login string) (*user, error)
	UserDel(uid string) error

	AddProject(title string) error
	GetProject(title string) (projectId string, err error)
	DelProject(id string) (err error)
}

type Mongo struct {
	Session *mgo.Session
	DbName  string
}

func New(db string, url string) (engine *Mongo, err error) {
	logger = log.New(os.Stdout, "speedlog mongodb ", log.LstdFlags|log.Lshortfile)
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
		logger.Printf("[debug] index check")

		logger.Printf("[debug] collection `user`")
		coll := mg.Collection(userCollection, nil)
		err = coll.EnsureIndex(mgo.Index{
			Key:    []string{"login"},
			Unique: true,
		})
		if err != nil {
			logger.Printf("[error] cannot create index `user`: %s", err)
		}

		logger.Printf("[debug] collection `project`")
		coll = mg.Collection(projectCollection, nil)
		err = coll.EnsureIndex(mgo.Index{
			Key:    []string{"title"},
			Unique: true,
		})

		if err != nil {
			logger.Printf("[error] cannot create index `project`: %s", err)
		}

		logger.Printf("[debug] index check complete")
	})
	return
}
