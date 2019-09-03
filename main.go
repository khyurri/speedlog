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

const defaultTimezone = "UTC-0"

func parseTZ(timezone string) (*time.Location, error) {
	if timezone == defaultTimezone {
		return time.FixedZone(timezone, 0), nil
	}
	return time.LoadLocation(timezone)
}

func main() {

	cliParams := &struct {
		Mode        string `arg:"-m" help:"Available modes: runserver, adduser"`
		Mongo       string `arg:"-d" help:"Mongodb url. Default 127.0.0.1:27017"`
		Login       string `arg:"-l" help:"Mode adduser. Login for new user"`
		Password    string `arg:"-p" help:"Mode adduser. Password for new user"`
		JWTKey      string `arg:"-j" help:"JWT secret key."`
		AllowOrigin string `arg:"-o" help:"Add Access-Control-Allow-Origin header with passed by param value"`
		TZ          string `arg:"-t" help:"Timezone. Default UTC±00:00."`
		Graphite    string `arg:"-g" help:"Graphite host:port"`
	}{}

	////////////////////////////////////////
	//
	// DEFAULTS
	cliParams.Mode = "runserver"
	cliParams.Mongo = "127.0.0.1:27017"
	cliParams.TZ = defaultTimezone
	//
	////////////////////////////////////////

	arg.MustParse(cliParams)
	cLogger := log.New(os.Stdout, "speedlog ", log.LstdFlags|log.Lshortfile)

	dbEngine, err := mongo.New("speedlog", cliParams.Mongo)
	defer dbEngine.Session.Close()

	if err != nil {
		cLogger.Fatalf("[error] failed to initialize mongo: %v", err)
		return
	}

	location, err := parseTZ(cliParams.TZ)
	if err != nil {
		cLogger.Fatalf("[error] failed to parse timezone: %v", err)
	}

	env := engine.NewEnv(dbEngine, cliParams.JWTKey, location)
	if len(cliParams.AllowOrigin) > 0 {
		env.AllowOrigin = cliParams.AllowOrigin
	}
	switch cliParams.Mode {
	case "runserver":

		if len(cliParams.JWTKey) == 0 {
			engine.Logger.Printf("[error] cannot start server. Required jwtkey")
			return
		}

		if len(cliParams.Graphite) > 0 {
			graphite := plugins.NewGraphite(cliParams.Graphite, location)
			graphite.Load(dbEngine)
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
	case "adduser":
		err := env.DBEngine.AddUser(cliParams.Login, cliParams.Password)
		if err != nil {
			log.Fatal(err)
		}
	}

}
