package rest

import (
	"github.com/go-chi/jwtauth"
	"github.com/khyurri/speedlog/engine"
	"net/http"
)

type App struct {
	Eng *engine.Engine
}

type AppHandlerFunc func(http.ResponseWriter, *http.Request, *engine.Engine)

func New(eng *engine.Engine) *App {
	return &App{eng}
}

func (app *App) MongoEngine(next AppHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbEngine := app.Eng.DBEngine.Clone()
		defer dbEngine.Close()
		eng := &engine.Engine{
			DBEngine:   dbEngine,
			Logger:     app.Eng.Logger,
			SigningKey: app.Eng.SigningKey,
		}
		next(w, r, eng)
	}
}

func (app *App) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := jwtauth.TokenFromHeader(r)
		_, err := app.Eng.SigningKey.Decode(t)
		if err != nil {
			response := &Resp{}
			response.Status = StatusForbidden
			response.Render(w)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
