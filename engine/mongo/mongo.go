package mongo

import (
	"github.com/globalsign/mgo"
	"log"
)

type Engine struct {
	Session *mgo.Session
	DB      string
	logger  *log.Logger
}

func New(db string, url string, logger *log.Logger) (engine *Engine, err error) {
	engine = &Engine{DB: db, logger: logger}
	engine.Session, err = mgo.Dial(url)
	return engine, err
}

func (engine *Engine) Clone() *Engine {
	session := engine.Session.Clone()
	return &Engine{
		Session: session,
		DB:      engine.DB,
		logger:  engine.logger,
	}
}

func (engine *Engine) Close() {
	engine.Session.Close()
}

func (engine *Engine) Collection(collection string) *mgo.Collection {
	return engine.Session.DB(engine.DB).C(collection)
}
