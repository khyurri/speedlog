package rest

import (
	"github.com/khyurri/speedlog/engine"
	"github.com/khyurri/speedlog/engine/mongo"
	"log"
	"net/http"
)

type App struct {
	dbEngine *mongo.Engine
	Logger   *log.Logger
}

type AppHandlerFunc func(http.ResponseWriter, *http.Request, *engine.Engine)

func New(eng *engine.Engine) *App {
	return &App{eng.DBEngine, eng.Logger}
}

func (app *App) MongoEngine(next AppHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbEngine := app.dbEngine.Clone()
		defer dbEngine.Close()
		eng := &engine.Engine{
			DBEngine: dbEngine,
			Logger:   app.Logger,
		}
		next(w, r, eng)
	}
}
