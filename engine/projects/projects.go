package projects

import (
	"errors"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine"
	"github.com/khyurri/speedlog/rest"
	"net/http"
)

const (
	collection = "project"
)

type Project struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Title string
}

func ExportRoutes(router *mux.Router, app *rest.App) {
	private := router.PathPrefix("/private/").Subrouter()
	private.HandleFunc("/project/", app.MongoEngine(RegisterProjectHttp)).
		Methods("PUT")
	private.Use(rest.JWTMiddleware)
}

func RegisterProjectHttp(http.ResponseWriter, *http.Request, *engine.Engine) {
	fmt.Println("HELLO!")
}

func ProjectExists(title string, eng *engine.Engine) (projectId bson.ObjectId, err error) {
	t := Project{}
	dbEngine := eng.DBEngine
	err = dbEngine.Collection(collection).Find(bson.M{"title": title}).One(&t)
	if err != nil {
		if err.Error() == "not found" {
			return projectId, errors.New(fmt.Sprintf("Project %s does not exists", title))
		} else {
			// TODO: just log it
			panic(err)
		}
	} else {
		projectId = t.ID
	}
	return t.ID, err
}
