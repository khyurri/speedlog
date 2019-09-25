package main

import (
	"errors"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/gorilla/mux"
	"github.com/khyurri/speedlog/engine"
	"github.com/khyurri/speedlog/engine/mongo"
	"github.com/khyurri/speedlog/plugins"
	"github.com/khyurri/speedlog/utils"
	"log"
	"net/http"
	"sync"
	"time"
)

const defaultTimezone = "UTC-0"

type params struct {
	Mode        string `arg:"positional" help:"Available modes: runserver, adduser, addproject. Default: runserver"`
	Mongo       string `arg:"-d" help:"Mode runserver. Mongodb url. Default 127.0.0.1:27017"`
	JWTKey      string `arg:"-j" help:"Mode runserver. JWT secret key."`
	AllowOrigin string `arg:"-o" help:"Mode runserver. Add Access-Control-Allow-Origin header with passed by param value"`
	TZ          string `arg:"-t" help:"Mode runserver. Timezone. Default UTC±00:00."`
	Graphite    string `arg:"-g" help:"Mode runserver. Graphite host:port"`
	EventsTTL   int    `arg:"--ttl" help:"Mode runserver. Time in seconds after which events are deleted. Default 0 — never"`
	Project     string `arg:"-r" help:"Modes runserver, addproject. Project title."`
	Login       string `arg:"-l" help:"Mode adduser. Login for new user"`
	Password    string `arg:"-p" help:"Mode adduser. Password for new user"`
}

func parseTZ(timezone string) (*time.Location, error) {
	if timezone == defaultTimezone {
		return time.FixedZone(timezone, 0), nil
	}
	return time.LoadLocation(timezone)
}

func addProjectMode(cliParams *params, dbEngine mongo.DataStore) (err error) {
	if len(cliParams.Project) > 0 {
		err = dbEngine.AddProject(cliParams.Project)
	} else {
		err = errors.New("--project param not found")
	}
	return
}

// init logger
func initLogger() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
	utils.Level = utils.LG_DEBUG
}

func main() {

	cliParams := &params{}

	////////////////////////////////////////
	//
	// DEFAULTS
	cliParams.Mode = "runserver"
	cliParams.Mongo = "127.0.0.1:27017"
	cliParams.TZ = defaultTimezone
	//
	////////////////////////////////////////

	arg.MustParse(cliParams)

	initLogger()

	dbEngine, err := mongo.New("speedlog", cliParams.Mongo)
	utils.Ok(err)
	defer dbEngine.Session.Close()

	location, err := parseTZ(cliParams.TZ)
	utils.Ok(err)

	env := engine.NewEnv(dbEngine, cliParams.JWTKey, location)
	if len(cliParams.AllowOrigin) > 0 {
		env.AllowOrigin = cliParams.AllowOrigin
	}
	switch cliParams.Mode {
	case "runserver":

		if len(cliParams.JWTKey) == 0 {
			fmt.Println("Cannot start server. Required jwtkey")
			return
		}

		////////////////////////////////////////////////////////////////////////////////
		// LOAD PLUGINS

		var plgns []plugins.Plugin
		var stopped sync.WaitGroup
		sigStop := make(plugins.SigChan)

		if len(cliParams.Graphite) > 0 {
			graphite := plugins.NewGraphite(cliParams.Graphite, time.Minute*1)
			plgns = append(plgns, graphite)
		}

		if cliParams.EventsTTL > 0 {
			cleaner := plugins.NewCleaner(cliParams.EventsTTL, time.Minute*1)
			plgns = append(plgns, cleaner)
		}

		go plugins.LoadPlugins(plgns, sigStop, &stopped, dbEngine)

		// END LOAD PLUGINS
		////////////////////////////////////////////////////////////////////////////////

		if len(cliParams.Project) > 0 {
			_ = dbEngine.AddProject(cliParams.Project)
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
		err = srv.ListenAndServe()

		// UNLOAD PLUGINS
		sigStop <- struct{}{}
		stopped.Wait() // TODO: add timeout

	case "adduser":
		err := env.DBEngine.AddUser(cliParams.Login, cliParams.Password)
		utils.Ok(err)
	case "addproject":
		err = addProjectMode(cliParams, dbEngine)
		utils.Ok(err)
	default:
		utils.Ok(errors.New(fmt.Sprintf("unknown mode '%s'", cliParams.Mode)))
	}

}
