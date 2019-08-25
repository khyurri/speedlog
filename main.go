package main

import (
	"github.com/alexflint/go-arg"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine"
	"github.com/khyurri/speedlog/engine/mongo"
	"github.com/khyurri/speedlog/plugins"
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
	engine.Logger = log.New(os.Stdout, "speedlog ", log.LstdFlags|log.Lshortfile)

	dbEngine, err := mongo.New("speedlog", config.Mongo)
	defer dbEngine.Session.Close()

	if err != nil {
		engine.Logger.Fatalf("failed to initialize mongo: %v", err)
		return
	}

	env := engine.NewEnv(dbEngine, config.JWTKey)
	switch config.Mode {
	case "runserver":

		if len(config.JWTKey) == 0 {
			engine.Logger.Printf("[error] cannot start server. Required jwtkey")
			return
		}

		r := mux.NewRouter()

		graphite := plugins.NewGraphite("1")
		graphite.Load(dbEngine)

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
	case "adduser":
		err := env.DBEngine.AddUser(config.Login, config.Password)
		if err != nil {
			log.Fatal(err)
		}
	}

}
