package main

import (
	"github.com/alexflint/go-arg"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine"
	"github.com/khyurri/speedlog/engine/mongo"
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

func runApp(cfg *config, env engine.AppEnvironment) {

	switch cfg.Mode {
	case "runserver":

		if len(cfg.JWTKey) == 0 {
			return
		}

		r := mux.NewRouter()

		env.ExportEventRoutes(r)
		env.ExportUserRoutes(r)
		env.ExportProjectRoutes(r)

		srv := &http.Server{
			Handler:      r,
			Addr:         ":8012",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.Fatal(srv.ListenAndServe())
		//case "adduser":
		//	//err := users.AddUser(cfg.Login, cfg.Password, env)
		//	if err != nil {
		//		log.Fatal(err)
		//	}
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

	dbEngine, err := mongo.New("speedlog", config.Mongo)
	defer dbEngine.Session.Close()

	if err != nil {
		logger.Fatalf("failed to initialize mongo: %v", err)
		return
	}

	eng := engine.New(dbEngine, logger, config.JWTKey)
	runApp(config, eng)

}
