package engine

import (
	"github.com/go-chi/jwtauth"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine/mongo"
	"time"
)

type AppEnvironment interface {
	ExportUserRoutes(router *mux.Router)
	ExportEventRoutes(router *mux.Router)
	ExportProjectRoutes(router *mux.Router)
}

// Env - core struct for storing dependencies
type Env struct {
	DBEngine    mongo.DataStore
	SigningKey  *jwtauth.JWTAuth
	AllowOrigin string
	Location    *time.Location
}

// NewEnv - create new env struct
func NewEnv(dbEngine mongo.DataStore, signingKey string, location *time.Location) *Env {
	k := jwtauth.New("HS256", []byte(signingKey), nil)
	return &Env{
		DBEngine:   dbEngine,
		SigningKey: k,
		Location:   location,
	}
}

func (env *Env) ExportUserRoutes(router *mux.Router) {
	router.HandleFunc("/login/", env.authenticateHttp()).
		Methods("POST", "OPTIONS")
	router.Use(env.corsMiddleware)
}

func (env *Env) ExportProjectRoutes(router *mux.Router) {
	private := router.PathPrefix("/private/").Subrouter()
	private.HandleFunc("/project/", env.createProjectHttp()).
		Methods("PUT", "OPTIONS")
	router.Use(env.corsMiddleware)
	private.Use(env.JWTMiddleware)
}

func (env *Env) ExportEventRoutes(router *mux.Router) {
	router.HandleFunc("/event/", env.createEventHttp()).
		Methods("PUT", "POST", "OPTIONS")

	router.HandleFunc("/events/", env.createEventsHttp()).
		Methods("PUT", "POST", "OPTIONS")

	private := router.PathPrefix("/private/").Subrouter()
	private.HandleFunc("/events/", env.getEventsHttp()).
		Methods("GET", "OPTIONS").
		Queries("metricName", "{metricName:.+}").
		Queries("metricTimeFrom", "{metricTimeFrom:.+}").
		Queries("metricTimeTo", "{metricTimeTo:.+}").
		Queries("project", "{project:.+}").
		Queries("groupBy", "{groupBy:.+}")
	router.Use(env.corsMiddleware)
	private.Use(env.JWTMiddleware)
}
