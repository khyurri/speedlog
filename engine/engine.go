package engine

import (
	"github.com/go-chi/jwtauth"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine/mongo"
)

type AppEnvironment interface {
	ExportUserRoutes(router *mux.Router)
	ExportEventRoutes(router *mux.Router)
	ExportProjectRoutes(router *mux.Router)
}

// Env - core struct for storing dependencies
type Env struct {
	DBEngine   mongo.DataStore
	SigningKey *jwtauth.JWTAuth
}

// NewEnv - create new env struct
func NewEnv(dbEngine mongo.DataStore, signingKey string) *Env {
	k := jwtauth.New("HS256", []byte(signingKey), nil)
	return &Env{dbEngine, k}
}

func (env *Env) ExportUserRoutes(router *mux.Router) {
	router.HandleFunc("/login/", env.authenticateHttp()).
		Methods("POST")
}

func (env *Env) ExportProjectRoutes(router *mux.Router) {
	private := router.PathPrefix("/private/").Subrouter()
	private.HandleFunc("/project/", env.addProjectHttp()).
		Methods("PUT")
	private.Use(env.JWTMiddleware)
}

func (env *Env) ExportEventRoutes(router *mux.Router) {
	router.HandleFunc("/event/", env.saveEventHttp()).
		Methods("PUT")

	private := router.PathPrefix("/private/").Subrouter()
	private.HandleFunc("/events/", env.getEventsHttp()).
		Methods("GET").
		Queries("metricName", "{metricName:.+}").
		Queries("metricTimeFrom", "{metricTimeFrom:.+}").
		Queries("metricTimeTo", "{metricTimeTo:.+}").
		Queries("project", "{project:.+}").
		Queries("groupBy", "{groupBy:.+}")
	private.Use(env.JWTMiddleware)
}
