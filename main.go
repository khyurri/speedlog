package main

import (
	"github.com/alexflint/go-arg"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine"
	"github.com/khyurri/speedlog/engine/events"
	"github.com/khyurri/speedlog/engine/mongo"
	"github.com/khyurri/speedlog/engine/projects"
	"github.com/khyurri/speedlog/engine/users"
	"github.com/khyurri/speedlog/rest"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	Mode     string `arg:"-m" help:"Available modes: runserver, adduser"`
	Mongo    string `arg:"-d" help:"Mongodb url. Default 127.0.0.1:27017"`
	Login    string `arg:"-l" help:"Mode adduser. Login for new user "`
	Password string `arg:"-p" help:"Mode adduser. Password for new user"`
	JWTKey   string `arg:"-j" help:"JWT secret key."`
}

func runApp(cfg *config, eng *engine.Engine) {

	switch cfg.Mode {
	case "runserver":

		if len(cfg.JWTKey) == 0 {
			eng.Logger.Panic("missing jwtkey")
			return
		}

		app := rest.New(eng)
		r := mux.NewRouter()

		events.ExportRoutes(r, app)
		users.ExportRoutes(r, app)
		projects.ExportRoutes(r, app)

		srv := &http.Server{
			Handler:      r,
			Addr:         ":8012",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.Fatal(srv.ListenAndServe())
	case "adduser":
		err := users.AddUser(cfg.Login, cfg.Password, eng)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {

	config := &config{}

	////////////////////////////////////////
	//
	// DEFAULTS
	config.Mode = "runserver"
	config.Mongo = "127.0.0.1:27017"
	//
	////////////////////////////////////////

	arg.MustParse(config)
	logger := log.New(os.Stdout, "speedlog ", log.LstdFlags|log.Lshortfile)

	dbEngine, err := mongo.New("speedlog", config.Mongo, logger)

	defer dbEngine.Close()

	if err != nil {
		logger.Fatalf("failed to initialize mongo: %v", err)
		return
	}

	eng := engine.New(dbEngine, logger, config.JWTKey)
	runApp(config, eng)

}
