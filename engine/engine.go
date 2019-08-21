package engine

import (
	"github.com/go-chi/jwtauth"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine/mongo"
	"log"
)

type AppEnvironment interface {
	ExportUserRoutes(router *mux.Router)
	ExportEventRoutes(router *mux.Router)
	ExportProjectRoutes(router *mux.Router)
}

// Env - core struct for storing dependencies
type Env struct {
	DBEngine   mongo.DataStore
	Logger     *log.Logger
	SigningKey *jwtauth.JWTAuth
}

// New - create new env struct
func New(dbEngine mongo.DataStore, logger *log.Logger, signingKey string) *Env {
	k := jwtauth.New("HS256", []byte(signingKey), nil)
	return &Env{dbEngine, logger, k}
}

func (env *Env) ExportUserRoutes(router *mux.Router) {
	router.HandleFunc("/login/", env.AuthenticateHttp).
		Methods("POST")
}

func (env *Env) ExportProjectRoutes(router *mux.Router) {
	private := router.PathPrefix("/private/").Subrouter()
	private.HandleFunc("/project/", env.RegisterProjectHttp).
		Methods("PUT")
	private.Use(env.JWTMiddleware)
}

func (env *Env) ExportEventRoutes(router *mux.Router) {
	router.HandleFunc("/{project}/event/", env.SaveEventHttp).
		Methods("PUT")

	private := router.PathPrefix("/private/").Subrouter()
	private.HandleFunc("/{project}/events/", env.GetEventsHttp).
		Methods("GET").
		Queries("metric_name", "{metricName}").
		Queries("metric_time_from", "{metricTimeFrom}").
		Queries("metric_time_to", "{metricTimeTo}").
		Queries("group_by", "{groupBy}")
	private.Use(env.JWTMiddleware)
}
