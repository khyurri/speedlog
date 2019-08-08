package projects

import (
	"encoding/json"
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
	Title string        `bson:"title"`
}

type RegisterProjectRequest struct {
	Title string `json:"title"`
}

func ExportRoutes(router *mux.Router, app *rest.App) {
	private := router.PathPrefix("/private/").Subrouter()
	private.HandleFunc("/project/", app.MongoEngine(RegisterProjectHttp)).
		Methods("PUT")
	private.Use(rest.JWTMiddleware)
}

func RegisterProjectHttp(w http.ResponseWriter, r *http.Request, eng *engine.Engine) {
	var err error

	response := &rest.Resp{}
	defer response.Render(w)

	p := &RegisterProjectRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(p)
	if err != nil {
		response.Status = rest.StatusErr
		eng.Logger.Fatal(err)
		return
	}
	if len(p.Title) == 0 {
		response.Status = rest.StatusErr
		return
	}
	_, err = ProjectExists(p.Title, eng)
	// TODO: simplify this block
	if err != nil {
		if err.Error() == "not found" {
			registered := RegisterProject(p.Title, eng)
			if !registered {
				response.Status = rest.StatusIntErr
			} else {
				response.Status = rest.StatusOk
			}
		} else {
			response.Status = rest.StatusIntErr
		}
	} else {
		response.Status = rest.StatusExists
	}

}

func ProjectExists(title string, eng *engine.Engine) (projectId bson.ObjectId, err error) {
	t := Project{}
	dbEngine := eng.DBEngine
	err = dbEngine.Collection(collection).Find(bson.M{"title": title}).One(&t)
	if err != nil {
		return
	} else {
		projectId = t.ID
	}
	return
}

func RegisterProject(title string, eng *engine.Engine) bool {
	t := &Project{Title: title}
	col := eng.DBEngine.Collection("project")
	err := col.Insert(t)
	if err != nil {
		eng.Logger.Fatal(err)
	}
	return true
}
