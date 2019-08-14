package test

import (
	"github.com/khyurri/speedlog/engine"
	"github.com/khyurri/speedlog/engine/mongo"
	"log"
	"os"
)

var DBName = "test_speedlog"

type SpeedLogTest struct {
	JWTTestKey string
	Mongo      string
	Logger     *log.Logger
	DBEngine   *mongo.Engine
	Engine     *engine.Engine
}

func (t *SpeedLogTest) Init() {
	t.JWTTestKey = "test_key"
	t.Mongo = "127.0.0.1:27017"
	t.Logger = log.New(os.Stdout, DBName+" ", log.LstdFlags|log.Lshortfile)
	t.DBEngine, _ = mongo.New(DBName, t.Mongo, t.Logger)
	t.Engine = engine.New(t.DBEngine, t.Logger, t.JWTTestKey)

	// clear mongo
	err := t.DBEngine.DropDatabase()
	if err != nil {
		t.Logger.Panic(err)
	}

}
