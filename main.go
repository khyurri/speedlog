package main

import (
	"flag"
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

func main() {
	//////////////////////////////////////////////////////

	var mode = flag.String("mode", "runserver", "modes: runserver, adduser")
	var login = flag.String("login", "", "login for new user")
	var password = flag.String("password", "", "password for new user")

	flag.Parse()

	//////////////////////////////////////////////////////

	logger := log.New(os.Stdout, "speedlog ", log.LstdFlags|log.Lshortfile)
	dbEngine, err := mongo.New("speedlog", "127.0.0.1:27017", logger)
	defer dbEngine.Close()

	if err != nil {
		logger.Fatalf("failed to initialize mongo: %v", err)
		return
	}

	eng := engine.New(dbEngine, logger)

	switch *mode {
	case "runserver":
		app := rest.New(eng)
		r := mux.NewRouter()

		events.ExportRoutes(r, app)
		users.ExportRoutes(r, app)
		projects.ExportRoutes(r, app)

		srv := &http.Server{
			Handler:      r,
			Addr:         "127.0.0.1:8012",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.Fatal(srv.ListenAndServe())
	case "adduser":
		err = users.AddUser(*login, *password, eng)
		if err != nil {
			log.Fatal(err)
		}
	}

}
